package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db"

	"go.mongodb.org/mongo-driver/mongo"
)

func HandleGetDeviceAddressByDeviceAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the content type to JSON

	// Get the device address from the URL
	deviceAddress := r.URL.Path[len("/checkdeviceaddresses/"):]
	log.Println("Received request for device address:", deviceAddress)

	// Get the data from the database
	deviceAddresses, err := db.GetDeviceAddressByDeviceAddress(deviceAddress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No documents found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if the result is empty
	if len(deviceAddresses) == 0 {
		log.Println("No device addresses found for:", deviceAddress)
		http.Error(w, "No device addresses found", http.StatusNotFound)
		return
	}

	// Encode the data into JSON
	if err := json.NewEncoder(w).Encode(deviceAddresses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Successfully retrieved device addresses for:", deviceAddress)
}
