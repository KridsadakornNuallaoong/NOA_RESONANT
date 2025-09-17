package ws

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type DeviceSession struct {
	UserConn *websocket.Conn
	HWConn   *websocket.Conn
	Mutex    sync.Mutex
	Ready    chan struct{}
}

var (
	waitingSession *DeviceSession
	waitingMutex   sync.Mutex
)

var (
	deviceSessions = make(map[string]*DeviceSession)
	sessionMutex   sync.Mutex
	upgrade        = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections
		},
	}
)

// HandleGetDeviceIDWebSocket handles WebSocket connections for both User and Hardware
func HandleGetDeviceIDWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// อ่าน role จาก client (user / hardware)
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read role:", err)
		conn.Close()
		return
	}
	role := string(message)
	log.Println("Connected role:", role)

	switch role {
	case "user":
		registerConnection("user", conn)
	case "hardware":
		registerConnection("hardware", conn)
	default:
		log.Println("Invalid role:", role)
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid role"))
		conn.Close()
	}
}

func registerConnection(role string, conn *websocket.Conn) {
	waitingMutex.Lock()
	defer waitingMutex.Unlock()

	if waitingSession == nil {
		// สร้าง session ใหม่ แล้วรออีกฝั่ง
		waitingSession = &DeviceSession{
			Ready: make(chan struct{}),
		}
		if role == "user" {
			waitingSession.UserConn = conn
		} else {
			waitingSession.HWConn = conn
		}
		log.Println(role, "is waiting for a pair...")
		return
	}

	// มีฝั่งหนึ่งรออยู่แล้ว ➜ จับคู่
	if role == "user" && waitingSession.UserConn == nil {
		waitingSession.UserConn = conn
	} else if role == "hardware" && waitingSession.HWConn == nil {
		waitingSession.HWConn = conn
	} else {
		log.Println("Duplicate role or both sides already connected.")
		conn.WriteMessage(websocket.TextMessage, []byte("Pairing error"))
		conn.Close()
		return
	}

	// ถ้าทั้งสองฝั่งพร้อมแล้ว ➜ เริ่มจับคู่
	if waitingSession.UserConn != nil && waitingSession.HWConn != nil {
		go handlePairedSession(waitingSession)
		waitingSession = nil
	}
}

func handlePairedSession(session *DeviceSession) {
	deviceID := strconv.Itoa(time.Now().Year()) + strconv.Itoa(rand.Intn(100000000))
	log.Println("Paired deviceID:", deviceID)

	data := map[string]string{"deviceID": deviceID}

	// ส่ง deviceID ให้ทั้งสองฝั่ง
	if err := session.UserConn.WriteJSON(data); err != nil {
		log.Println("Error sending to user:", err)
	}
	if err := session.HWConn.WriteJSON(data); err != nil {
		log.Println("Error sending to hardware:", err)
	}

	// ปิดการเชื่อมต่อ
	session.UserConn.Close()
	session.HWConn.Close()

	log.Println("Connections closed for deviceID:", deviceID)
}
