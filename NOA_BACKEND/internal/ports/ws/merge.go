package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	schema "GOLANG_SERVER/pkg/schema"
)

const FrameSize = 50 // Define the size of the sliding frame

var (
	currentClient *websocket.Conn // Store the current WebSocket connection
	clientMutex   sync.Mutex      // Mutex to protect the currentClient
)

func HandleWebSocketMerge(w http.ResponseWriter, r *http.Request) {
	// Extract UserID or DeviceID from query parameters
	userID := r.URL.Query().Get("DeviceID")
	if userID == "" {
		http.Error(w, "Missing DeviceID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}

	// Ensure only one client is connected at a time
	clientMutex.Lock()
	if currentClient != nil {
		log.Println("Closing previous connection")
		currentClient.Close() // Close the previous connection
	}
	currentClient = conn
	clientMutex.Unlock()

	defer func() {
		clientMutex.Lock()
		if currentClient == conn {
			currentClient = nil
		}
		clientMutex.Unlock()
		conn.Close()
	}()

	// Create unique frames for this WebSocket connection
	frameX := make([]float64, FrameSize)
	frameY := make([]float64, FrameSize)
	frameZ := make([]float64, FrameSize)

	// Create a unique channel to receive data for this WebSocket connection
	dataChan := make(chan []byte)
	//go SubscribeMQTTTopic(dataChan)

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

		// Merge all data into a single JSON structure
		mergedData := struct {
			Timestamp int64     `json:"timestamp"`
			FrameX    []float64 `json:"x"`
			FrameY    []float64 `json:"y"`
			FrameZ    []float64 `json:"z"`
		}{
			Timestamp: mqttData.Timestamp,
			FrameX:    frameX,
			FrameY:    frameY,
			FrameZ:    frameZ,
		}

		// Send the updated frames to the client
		clientMutex.Lock()
		if currentClient == conn {
			if err := conn.WriteJSON(mergedData); err != nil {
				log.Println("Error sending data to client:", err)
				clientMutex.Unlock()
				break // Exit the loop if there's an error sending data
			}
		}
		clientMutex.Unlock()
	}
}

func updateFrames(frameX, frameY, frameZ []float64, data schema.MQTTData) {
	// Shift the frames to the left
	copy(frameX, frameX[1:])
	copy(frameY, frameY[1:])
	copy(frameZ, frameZ[1:])

	// Append the new acceleration values to the end of the frames
	frameX[FrameSize-1] = data.X.Acceleration
	frameY[FrameSize-1] = data.Y.Acceleration
	frameZ[FrameSize-1] = data.Z.Acceleration
}
