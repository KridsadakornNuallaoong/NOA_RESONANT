package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"GOLANG_SERVER/components/db"

	"golang.org/x/crypto/bcrypt"
)

func HandleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the content type to JSON

	// Get the device address from the json request
	var userDetail map[string]string
	if err := json.NewDecoder(r.Body).Decode(&userDetail); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deviceName := userDetail["deviceName"]
	if deviceName == "" {
		http.Error(w, "Device name is required", http.StatusBadRequest)
		return
	}
	deviceID := userDetail["deviceID"]
	if deviceID == "" {
		http.Error(w, "Device ID is required", http.StatusBadRequest)
		return
	}
	userID := userDetail["userID"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	devicePassword := userDetail["password"]
	if devicePassword == "" {
		http.Error(w, "Device password is required", http.StatusBadRequest)
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(devicePassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	devicePassword = string(hashedPassword)

	// display device address
	// log.Println("Device ID:", deviceID)
	// log.Println("User ID:", userID)
	// log.Println("Device Password:", devicePassword)

	// Check Email
	user, err := db.FindUserID(userID)
	if err != nil {
		http.Error(w, "Invalid user ID.", http.StatusUnauthorized)
		return
	} else {
		log.Println("User ID is OK")
	}

	// Check if the device address already exists in the database
	exists, err := db.HandlercheckDeviceID(deviceID)
	if err != nil {
		http.Error(w, "Error checking device ID", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "Device ID already exists", http.StatusBadRequest)
		return
	} else {
		log.Println("Device ID is OK")
	}

	log.Printf("UserID: %s\n", user.ID)

	response := map[string]string{
		"message":  "Device created successfully",
		"deviceID": deviceID,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the deviceDetail to the database
	if err := db.SaveDevice(deviceName, deviceID, userID, devicePassword); err != nil {
		http.Error(w, "Error saving device details", http.StatusInternalServerError)
		return
	}
	log.Println("Device details saved successfully")

	// Time out
	elapsedTime := time.Since(startTime)
	log.Printf("Device Authentication time for  %s\n", elapsedTime)
}
