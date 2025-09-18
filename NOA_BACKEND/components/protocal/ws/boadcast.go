package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"GOLANG_SERVER/components/env"
	"GOLANG_SERVER/components/schema"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

// WebSocket variables
var (
	// clients is a package-level, concurrency-safe container that holds active WebSocket connections
	// keyed by user ID. It embeds a sync.Mutex to serialize access to its maps and is initialized with
	// an empty connections map.
	//
	// Fields:
	//   - connections: primary map[string]*websocket.Conn that maps a user ID to its active WebSocket
	//     connection.
	//   - make: auxiliary map[string]*websocket.Conn (appears to mirror the connections map; preserve or
	//     remove only with care).
	//
	// Usage:
	//   Always acquire the embedded mutex with clients.Lock() before reading or modifying either map and
	//   release it with clients.Unlock() afterwards to avoid data races.
	clients = struct {
		sync.Mutex
		connections map[string]*websocket.Conn
		aux         map[string]*websocket.Conn // renamed from `make` to avoid confusion with builtin
	}{
		connections: make(map[string]*websocket.Conn),
		aux:         make(map[string]*websocket.Conn),
	}

	// deviceClients: maps deviceID -> *websocket.Conn (used by predict/merge handlers)
	deviceClients = struct {
		sync.Mutex
		connections map[string]*websocket.Conn
		aux         map[string]*websocket.Conn
	}{
		connections: make(map[string]*websocket.Conn),
		aux:         make(map[string]*websocket.Conn),
	}

	// single upgrader for the package
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var connectionRegistry = struct {
	sync.Mutex
	activeKeys map[string]string // key: clientID or clientKey, value: "merge" or "broadcast"
}{
	activeKeys: make(map[string]string),
}

// WebSocket handler for multiple userIDs and deviceIDs
func HandleWebSocketBoadcast(w http.ResponseWriter, r *http.Request) {
	// Extract userID and deviceID from query parameters
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	deviceID := r.URL.Query().Get("deviceID")
	if deviceID == "" {
		http.Error(w, "Missing deviceID", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Create a unique key for this client
	clientKey := userID + ":" + deviceID

	// Register the client in connectionRegistry
	connectionRegistry.Lock()
	if _, exists := connectionRegistry.activeKeys[clientKey]; exists {
		connectionRegistry.Unlock()
		http.Error(w, "Client is already connected on another handler", http.StatusConflict)
		return
	}
	connectionRegistry.activeKeys[clientKey] = "broadcast"
	connectionRegistry.Unlock()

	defer func() {
		// Remove the client from connectionRegistry when the connection is closed
		connectionRegistry.Lock()
		delete(connectionRegistry.activeKeys, clientKey)
		connectionRegistry.Unlock()
	}()

	// Register the client
	clients.Lock()
	if existingConn, exists := clients.connections[clientKey]; exists { // Access the `connections` field
		log.Printf("Closing existing connection for clientKey=%s\n", clientKey)
		existingConn.Close()                   // Close the existing connection
		delete(clients.connections, clientKey) // Access the `connections` field
	}
	clients.connections[clientKey] = conn // Access the `connections` field
	clients.Unlock()

	// Log the new connection
	log.Printf("New WebSocket client connected: userID=%s, deviceID=%s\n", userID, deviceID)

	// Block until the connection is closed
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// Cleanup when the connection is closed
	clients.Lock()
	delete(clients.connections, clientKey) // Access the `connections` field
	clients.Unlock()
	log.Printf("WebSocket client disconnected: userID=%s, deviceID=%s\n", userID, deviceID)
}

// Send message to a specific userID and deviceID
func sendMessageToUser(clientKey string, data []byte) {
	clients.Lock()
	defer clients.Unlock()

	conn, exists := clients.connections[clientKey] // Access the `connections` field
	if exists {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error writing message to WebSocket client (key=%s): %v\n", clientKey, err)
			conn.Close()
			delete(clients.connections, clientKey) // Access the `connections` field
		}
	}
}

// SubscribeMQTTTopicForClient subscribes to MQTT for a specific client
func SubscribeMQTTTopicForClient(clientKey string, dataChan chan<- []byte) {
	// MQTT topic
	opts := mqtt.NewClientOptions().AddBroker(env.GetEnv("MQTT_BROKER"))
	opts.SetClientID(env.GetEnv("MQTT_CLIENT_ID") + "_" + clientKey) // Unique client ID for each subscription
	opts.SetUsername(env.GetEnv("MQTT_USERNAME"))
	opts.SetPassword(env.GetEnv("MQTT_PASSWORD"))
	client := mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error connecting to MQTT broker for clientKey=%s: %v\n", clientKey, token.Error())
		return
	}

	// Subscribe to the topic
	if token := client.Subscribe("vibration", 1, func(client mqtt.Client, msg mqtt.Message) {
		// Parse the incoming message
		var payload schema.Data
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Printf("Error unmarshaling MQTT message for clientKey=%s: %v\n", clientKey, err)
			return
		}

		// Check if the message is for this client
		if clientKey == payload.UserID+":"+payload.DeviceID {
			dataChan <- msg.Payload() // Send the message to the client's channel
		}
	}); token.Wait() && token.Error() != nil {
		log.Printf("Error subscribing to topic for clientKey=%s: %v\n", clientKey, token.Error())
		return
	}

	log.Printf("MQTT client subscribed to topic for clientKey=%s\n", clientKey)
}

func StartGlobalMQTTSubscriber() {
	opts := mqtt.NewClientOptions().AddBroker(env.GetEnv("MQTT_BROKER"))
	opts.SetClientID(env.GetEnv("MQTT_CLIENT_ID") + "_global") // ใช้ client ID แบบ global
	opts.SetUsername(env.GetEnv("MQTT_USERNAME"))
	opts.SetPassword(env.GetEnv("MQTT_PASSWORD"))
	client := mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Error connecting to MQTT broker:", token.Error())
	}

	// Subscribe to the topic
	topic := "vibration"
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		// Parse the incoming message
		var payload schema.Data
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Println("Failed to unmarshal payload:", err)
			return
		}

		// Create clientKey from userID and deviceID
		clientKey := payload.UserID + ":" + payload.DeviceID
		// Send the message to the WebSocket client with the matching clientKey
		sendMessageToUser(clientKey, msg.Payload())

	})
	if token.Wait() && token.Error() != nil {
		log.Fatal("Error subscribing to MQTT topic:", token.Error())
	}

	log.Println("Global MQTT subscriber started and subscribed to topic:", topic)
}
