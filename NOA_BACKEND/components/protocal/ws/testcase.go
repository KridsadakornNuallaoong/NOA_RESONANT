package ws

import (
	"log"
	"net/http"
)

// Testcase use WebSocket to read data in MongoDB for testing purpose
func testcase(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "[ERROR] Failed to upgrade to WebSocket", http.StatusInternalServerError)
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
	for {
		if _, _, err := conn.NextReader(); err != nil {
			log.Printf("[INFO] WebSocket closed: %s", deviceID)
			break
		}
	}

}
