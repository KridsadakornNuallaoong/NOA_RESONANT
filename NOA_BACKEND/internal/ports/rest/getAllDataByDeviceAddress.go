package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db"
)

// * get data use param
func HandleGetAllDataByDeviceAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get the device address from the URL
	deviceAddress := r.URL.Path[len("/data/"):]
	fmt.Println("Device Address:", deviceAddress)

	// Get the data from the database
	data, err := db.GetGyroDataByDeviceAddress(deviceAddress)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the data into JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
