package rest

import (
	"encoding/json"
	"net/http"
)

// 		//TODO--------------------------------------------------------------------------------------------------------------------------||

func Notification(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the content type to JSON

	response := map[string]string{
		"Message": "Notification sent successfully",
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 		//TODO--------------------------------------------------------------------------------------------------------------------------||
