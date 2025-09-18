package ws

import (
	"log"
	"net/http"
	"sync"
	"time"

	"GOLANG_SERVER/components/env"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

var client mqtt.Client // MQTT client

// WebSocket variables
var (
	// renamed from `clients` to `broadcastClients` to avoid redeclaration
	broadcastClients = struct {
		sync.Mutex
		conns map[*websocket.Conn]bool
	}{
		conns: make(map[*websocket.Conn]bool),
	}

	wsUpgrader = websocket.Upgrader{ // renamed upgrader -> wsUpgrader
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin function to allow all connections
			return true // Allow all connections by default
		},
	}
)

// WebSocket Multi cliend handler
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Set a ping handler to keep the connection alive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Add the client to the clients map
	broadcastClients.Lock()
	broadcastClients.conns[conn] = true
	broadcastClients.Unlock()

	// Log the new connection
	log.Println("New WebSocket client connected.")

	// Create a channel to receive data from MQTT
	dataChan := make(chan []byte)

	// Start the MQTT subscription in a separate goroutine
	go SubscribeMQTTTopic(dataChan)

	// Wait for messages from the MQTT client
	for {
		// Receive data from the channel
		data := <-dataChan

		// Broadcast the message to all WebSocket clients
		broadcastMessage(data)
	}
}

// Broadcast message to all WebSocket clients
func broadcastMessage(data []byte) {
	// Iterate over all clients
	broadcastClients.Lock()
	for conn := range broadcastClients.conns {

		// Process the message and store it in the database
		// Write the message to the client
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Error writing message to WebSocket client:", err)
			conn.Close()
			delete(broadcastClients.conns, conn)
		}
	}
	broadcastClients.Unlock()
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
		// log.Printf("Sub topic: %s\n", msg.Topic())

		// Send the message payload to the channel
		dataChan <- msg.Payload()

	}); token.Wait() && token.Error() != nil {
		log.Fatal("Error subscribing to topic:", token.Error())
	}

	log.Println("MQTT client connected and subscribed to topic.")
}

func HandleWsTest(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Set a ping handler to keep the connection alive
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Log the new connection
	log.Println("New WebSocket client connected.")

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("Error reading message from WebSocket client:", err)
			break
		}

		// Send a message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from server!")); err != nil {
			log.Println("Error writing message to WebSocket client:", err)
			break
		}

	}
}
