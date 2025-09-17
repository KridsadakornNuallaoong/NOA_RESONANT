package rest

import (
	"encoding/json"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db"
)

// Handle a request for the schema
func HandleGetAllData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get the data from the database
	data, err := db.GetGyroData()
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
