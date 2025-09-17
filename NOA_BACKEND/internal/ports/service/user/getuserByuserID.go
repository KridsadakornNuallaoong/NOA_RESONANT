package user

import (
	"encoding/json"
	"log"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db" // Import the db package
)

// UserResponse defines the structure of the user response without the password
type UserResponse struct {
	ID       string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetUserByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Allow only POST requests
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse the request body to get user ID
	var userDetails map[string]string
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle both lowercase and uppercase keys
	userID := userDetails["userID"]
	if userID == "" {
		userID = userDetails["UserID"]
	}

	// Check if user ID is provided
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Retrieve user details from the database using userID
	user, err := db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Create a response struct excluding the password
	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// Convert user details to JSON and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to convert user details to JSON", http.StatusInternalServerError)
		return
	}

	// Log the user details
	log.Println("User details retrieved successfully:")
}
