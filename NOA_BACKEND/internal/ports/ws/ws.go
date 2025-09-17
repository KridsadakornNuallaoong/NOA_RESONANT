package ws

import (
	"log"
	"net/http"
	"sync"
	"time"

	"GOLANG_SERVER/configs/env"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

var client mqtt.Client // MQTT client

// WebSocket variables
var (
	clients      = make(map[string]*websocket.Conn) // Map of userID to WebSocket connections
	clientsMutex = sync.Mutex{}                     // Mutex to protect the clients map
	upgrader     = websocket.Upgrader{              // Upgrader for WebSocket connections
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin function to allow all connections
			return true // Allow all connections by default
		},
	}
)

// WebSocket handler for multiple userIDs
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract userID from query parameters
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Ensure only one connection per userID
	clientsMutex.Lock()
	if existingConn, exists := clients[userID]; exists {
		log.Printf("Closing existing connection for userID=%s\n", userID)
		existingConn.Close() // Close the existing connection
		delete(clients, userID)
	}
	clients[userID] = conn
	clientsMutex.Unlock()

	// Set a ping handler to keep the connection alive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Log the new connection
	log.Printf("New WebSocket client connected: userID=%s\n", userID)

	// Create a channel to receive data from MQTT
	dataChan := make(chan []byte)

	// Start the MQTT subscription in a separate goroutine
	go SubscribeMQTTTopic(dataChan)

	// Wait for messages from the MQTT client
	for {
		// Receive data from the channel
		data := <-dataChan

		// Send the message to the specific userID
		sendMessageToUser(userID, data)
	}
}

// Send message to a specific userID
func sendMessageToUser(userID string, data []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	conn, exists := clients[userID]
	if exists {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error writing message to WebSocket client (userID=%s): %v\n", userID, err)
			conn.Close()
			delete(clients, userID)
		}
	}
}

func SubscribeMQTTTopic(dataChan chan<- []byte) {
	// MQTT topic
	opts := mqtt.NewClientOptions().AddBroker(env.GetEnv("MQTT_BROKER"))
	opts.SetClientID(env.GetEnv("MQTT_CLIENT_ID"))
	opts.SetUsername(env.GetEnv("MQTT_USERNAME"))
	opts.SetPassword(env.GetEnv("MQTT_PASSWORD"))
	client = mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Error connecting to MQTT broker:", token.Error())
	}

	// Subscribe to the topic
	if token := client.Subscribe("vibration", 1, func(client mqtt.Client, msg mqtt.Message) {
		// Send the message payload to the channel
		dataChan <- msg.Payload()
	}); token.Wait() && token.Error() != nil {
		log.Fatal("Error subscribing to topic:", token.Error())
	}

	log.Println("MQTT client connected and subscribed to topic.")
}
