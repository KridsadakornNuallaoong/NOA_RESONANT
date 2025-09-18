package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	schema "GOLANG_SERVER/components/schema"
)

const MQTT_TOPIC = "your/mqtt/topic"
const MergeFrameSize = 50 // Define the size of the sliding frame

var (
	clientMerge = make(map[*websocket.Conn]string) // Map of clients with identifiers
	frameMutex  sync.Mutex                         // Mutex to protect frame updates
)

func SendToSpecificClients(data []float64, targetID string) {
	frameMutex.Lock()
	defer frameMutex.Unlock()

	for client, clientID := range clientMerge {
		if clientID == targetID { // Send data to specific client with matching ID
			log.Printf("Sending data to client: %s, Data: %v\n", targetID, data)
			err := client.WriteJSON(data)
			if err != nil {
				log.Println("Error sending data to client:", err)
				client.Close()
				delete(clientMerge, client)
			}
		}
	}
}

func HandleWebSocketMerge(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// use frameMutex (declared above) to protect clientMerge
	frameMutex.Lock()
	clientMerge[conn] = "client1" // Assign a unique identifier to the client
	frameMutex.Unlock()

	defer func() {
		frameMutex.Lock()
		delete(clientMerge, conn)
		frameMutex.Unlock()
	}()

	// Create unique frames for this WebSocket connection
	frameX := make([]float64, MergeFrameSize)
	frameY := make([]float64, MergeFrameSize)
	frameZ := make([]float64, MergeFrameSize)

	// Create a unique channel to receive data for this WebSocket connection
	dataChan := make(chan []byte)
	go SubscribeMQTTTopic(dataChan)

	for {
		// Receive data from the channel
		data := <-dataChan
		log.Println("Received data from MQTT channel")

		// Decode the incoming MQTT data
		var mqttData schema.MQTTData
		err := json.Unmarshal(data, &mqttData)
		if err != nil {
			log.Println("Error decoding MQTT data:", err)
			continue
		}

		// Update the sliding frames for this connection
		updateFrames(frameX, frameY, frameZ, mqttData)

		// Log the updated frames
		log.Printf("Timestamp: %d, FrameX: %v, FrameY: %v, FrameZ: %v\n", mqttData.Timestamp, frameX, frameY, frameZ)

		// Send the updated frames to the client
		SendToSpecificClients(frameX, "client1")
	}
}

func updateFrames(frameX, frameY, frameZ []float64, data schema.MQTTData) {
	// Shift the frames to the left
	copy(frameX, frameX[1:])
	copy(frameY, frameY[1:])
	copy(frameZ, frameZ[1:])

	// Append the new acceleration values to the end of the frames (cast from float32 -> float64)
	frameX[MergeFrameSize-1] = float64(data.X.Acceleration)
	frameY[MergeFrameSize-1] = float64(data.Y.Acceleration)
	frameZ[MergeFrameSize-1] = float64(data.Z.Acceleration)
}
