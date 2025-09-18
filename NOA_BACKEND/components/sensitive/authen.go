package sensitive

import (
	"GOLANG_SERVER/components/db"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// AuthenRequest represents the structure of the request body
type AuthenRequest struct {
	Email    string `json:"email"`
	DeviceID string `json:"DeviceID"`
	Pass     string `json:"pass"`
}

// AuthenDevice function to authenticate device
func AuthenDevice(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	if r.Method != http.MethodPost { // Allow only POST requests
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse the request body to get user details
	var deviceDetails map[string]string
	if err := json.NewDecoder(r.Body).Decode(&deviceDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle both lowercase and uppercase keys
	email := deviceDetails["email"]
	if email == "" {
		email = deviceDetails["Email"]
	}
	pass := deviceDetails["password"]
	if pass == "" {
		pass = deviceDetails["Password"]
	}
	deviceID := deviceDetails["deviceID"]
	if deviceID == "" {
		deviceID = deviceDetails["DeviceID"]
	}

	// Check if user exists
	log.Println("Device authentication started")
	// log.Println("Email: ", email)
	// log.Println("Password: ", pass)
	// log.Println("DeviceID: ", deviceID)

	if email == "" || pass == "" || deviceID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := db.FindUser(email)
	if err != nil {
		log.Println("Error finding user:", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Authenticate the device using the provided email, password, and deviceID
	if _, err := db.AuthenDevice(email, pass, deviceID); err != nil {
		log.Println("Error authenticating device:", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// deviceHash := bson.M{
	// 	"email":  email,
	// 	"pass":   pass,
	// 	"userID": userID,
	// }

	// log.Println("Device hash: ", deviceHash)
	// Calculate the elapsed time

	response := map[string]string{
		"message": "Device authenticated successfully",
		"userID":  user.ID,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the elapsed time
	elapsedTime := time.Since(startTime)
	log.Printf("Device Authentication time for  %s\n", elapsedTime)
}
