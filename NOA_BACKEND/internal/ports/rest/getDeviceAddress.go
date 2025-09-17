package rest

import (
	"encoding/json"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db"
)

func HandleGetDeviceAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the content type to JSON

	// * get device address from database
	deviceAddresses, err := db.GetDeviceAddress()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// * send device addresses .json to client
	response := map[string][]string{"deviceAddresses": deviceAddresses}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
