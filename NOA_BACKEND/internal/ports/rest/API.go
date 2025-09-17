package rest

import (
	"encoding/json"
	"net/http"
)

// Response represents the structure of the API response
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// HandleAPI is a REST API handler
func HandleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check the HTTP method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  http.StatusMethodNotAllowed,
		})
		return
	}

	// Example response
	response := Response{
		Message: "REST API is running on the server",
		Status:  http.StatusOK,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
