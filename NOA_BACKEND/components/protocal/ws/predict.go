package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"GOLANG_SERVER/components/schema"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

const PredictFrameSize = 200

type SlidingWindow struct {
	X []float32
	Y []float32
	Z []float32
}

type PredictionResult struct {
	Prediction     [][]float32 `json:"prediction"`
	PredictedClass []int       `json:"predicted_class"`
	Timestamp      int64       `json:"timestamp"`
}

var (
	deviceFrames = struct {
		sync.Mutex
		frames map[string]*SlidingWindow
	}{frames: make(map[string]*SlidingWindow)}
	cooldownMap sync.Map
)

func HandleWebSocketPredict(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	deviceID := r.URL.Query().Get("deviceID")
	userID := r.URL.Query().Get("userID")
	if deviceID == "" || userID == "" {
		log.Println("[ERROR] Missing deviceID or userID")
		return
	}

	clients.Lock()
	if _, exists := clients.connections[deviceID]; exists {
		clients.Unlock()
		http.Error(w, "Device already connected", http.StatusConflict)
		return
	}
	clients.connections[deviceID] = conn
	clients.Unlock()

	defer func() {
		clients.Lock()
		delete(clients.connections, deviceID)
		clients.Unlock()
	}()

	startMQTTSubscribe(userID, deviceID)

	for {
		if _, _, err := conn.NextReader(); err != nil {
			log.Printf("[INFO] WebSocket closed: %s", deviceID)
			break
		}
	}
}

func startMQTTSubscribe(userID, deviceID string) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("mqtt_subscriber_" + deviceID)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("[ERROR] MQTT connect:", token.Error())
		return
	}

	if token := client.Subscribe("vibration", 1, func(_ mqtt.Client, msg mqtt.Message) {
		var data schema.GyroData
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Println("[ERROR] MQTT message unmarshal:", err)
			return
		}
		if data.DeviceID != deviceID {
			return // สนใจเฉพาะ device นี้
		}

		if ready, frame := updateSlidingWindow(deviceID, data.Data); ready {
			if _, ok := cooldownMap.Load(deviceID); !ok {
				go predictAndSend(userID, deviceID, frame)
				cooldownMap.Store(deviceID, true)
				time.AfterFunc(3*time.Second, func() {
					cooldownMap.Delete(deviceID)
				})
			}
		}
	}); token.Wait() && token.Error() != nil {
		log.Println("[ERROR] MQTT subscribe:", token.Error())
		return
	}
}

func updateSlidingWindow(deviceID string, data schema.GyroDataDetail) (bool, *SlidingWindow) {
	deviceFrames.Lock()
	defer deviceFrames.Unlock()

	f, exists := deviceFrames.frames[deviceID]
	if !exists {
		f = &SlidingWindow{
			X: make([]float32, PredictFrameSize),
			Y: make([]float32, PredictFrameSize),
			Z: make([]float32, PredictFrameSize),
		}
		deviceFrames.frames[deviceID] = f
	}

	f.X = append(f.X[1:], float32(data.X.Acceleration))
	f.Y = append(f.Y[1:], float32(data.Y.Acceleration))
	f.Z = append(f.Z[1:], float32(data.Z.Acceleration))

	return true, &SlidingWindow{
		X: append([]float32(nil), f.X...),
		Y: append([]float32(nil), f.Y...),
		Z: append([]float32(nil), f.Z...),
	}
}

func predictAndSend(userID, deviceID string, frame *SlidingWindow) {
	input := map[string]interface{}{
		"inputs": []interface{}{
			append([]interface{}{time.Now().UnixMilli()}, append(toInterface(frame.X), append(toInterface(frame.Y), toInterface(frame.Z)...)...)...),
		},
	}
	payload, _ := json.Marshal(input)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws/predict", nil)
	if err != nil {
		log.Println("[ERROR] Connect to Python WebSocket:", err)
		return
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		log.Println("[ERROR] Send to Python WebSocket:", err)
		return
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("[ERROR] Read from Python WebSocket:", err)
		return
	}

	var result PredictionResult
	if err := json.Unmarshal(message, &result); err != nil {
		log.Println("[ERROR] Unmarshal prediction result:", err)
		return
	}

	sendToClient(deviceID, &result)
	saveResult(userID, deviceID, &result)
}

func sendToClient(deviceID string, result *PredictionResult) {
	clients.Lock()
	conn, exists := clients.connections[deviceID]
	clients.Unlock()
	if !exists {
		log.Println("[WARN] No WebSocket client for deviceID:", deviceID)
		return
	}

	if err := conn.WriteJSON(result); err != nil {
		log.Println("[ERROR] Send prediction to client:", err)
		conn.Close()

		clients.Lock()
		delete(clients.connections, deviceID)
		clients.Unlock()
	}
}

func saveResult(userID, deviceID string, result *PredictionResult) {
	record := map[string]interface{}{
		"userID":         userID,
		"deviceID":       deviceID,
		"predictedClass": classLabel(result.PredictedClass),
		"result":         toPercent(result.Prediction),
	}

	file, err := os.OpenFile("notification.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[ERROR] Open file:", err)
		return
	}
	defer file.Close()

	data, _ := json.Marshal(record)
	file.WriteString(string(data) + "\n")
}

func classLabel(classes []int) string {
	if len(classes) == 0 {
		return "Unknown"
	}
	switch classes[0] {
	case 0:
		return "Close"
	case 1:
		return "Normal"
	case 2:
		return "Fault"
	default:
		return "Unknown"
	}
}

func toPercent(prediction [][]float32) [][]float32 {
	for i := range prediction {
		for j := range prediction[i] {
			prediction[i][j] *= 100
		}
	}
	return prediction
}

func toInterface(floats []float32) []interface{} {
	interfaces := make([]interface{}, len(floats))
	for i, v := range floats {
		interfaces[i] = v
	}
	return interfaces
}
