package rest

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"GOLANG_SERVER/internal/adapters/db" // Import the db package
)

func HandleGenerateDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json") // Set the content type to JSON

	var deviceID string
	for {
		// Generate a new device ID
		currentYear := time.Now().Year()                                                // Get the current year
		randomNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000000) // Generate a random 8-digit number
		deviceID = strconv.Itoa(currentYear) + strconv.Itoa(randomNumber)

		// Check if the DeviceID already exists in the database
		exists, err := db.HandlercheckDeviceID(deviceID)
		if err != nil {
			http.Error(w, "Error checking device ID in database", http.StatusInternalServerError)
			return
		}
		if !exists {
			break // If the DeviceID does not exist, exit the loop
		}
	}

	// Create a response
	response := map[string]string{
		"deviceID": deviceID,
	}

	// Send the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to generate device ID", http.StatusInternalServerError)
		return
	}
}
