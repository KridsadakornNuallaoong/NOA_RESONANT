package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db"
	schema "GOLANG_SERVER/pkg/schema"
)

// Handle a request to store data
func HandleStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into a GyroData struct
	var data schema.GyroData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// * print the data
	fmt.Printf("Data: %+v\n", data) // Print the data

	// Store the data in the database
	db.StoreGyroData(data)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Data stored!"}`)
}
