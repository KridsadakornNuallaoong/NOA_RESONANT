package sensitive

import (
	"GOLANG_SERVER/internal/adapters/db"
	"encoding/json"
	"log"
	"net/http"
	"time"
<<<<<<< Updated upstream:NOA_BACKEND/components/sensitive/authen.go

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
=======
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/sensitive/authen.go
)

// AuthenRequest represents the structure of the request body
type AuthenRequest struct {
<<<<<<< Updated upstream:NOA_BACKEND/components/sensitive/authen.go
	Email  string `json:"email"`
	UserID string `json:"userID"`
	Pass   string `json:"pass"`
=======
	Email    string `json:"email"`
	DeviceID string `json:"DeviceID"`
	Pass     string `json:"pass"`
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/sensitive/authen.go
}

// AuthenDevice function to authenticate device
func AuthenDevice(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	if r.Method != http.MethodPost { // Allow only POST requests
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse the request body to get user details
	var deviceDetails map[string]string
	if err := json.NewDecoder(r.Body).Decode(&deviceDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle both lowercase and uppercase keys
	email := deviceDetails["email"]
	if email == "" {
		email = deviceDetails["Email"]
	}
<<<<<<< Updated upstream:NOA_BACKEND/components/sensitive/authen.go
	pass := deviceDetails["pass"]
	if pass == "" {
		pass = deviceDetails["Pass"]
	}
	userID := deviceDetails["userID"]
	if userID == "" {
		userID = deviceDetails["UserID"]
=======
	pass := deviceDetails["password"]
	if pass == "" {
		pass = deviceDetails["Password"]
	}
	deviceID := deviceDetails["deviceID"]
	if deviceID == "" {
		deviceID = deviceDetails["DeviceID"]
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/sensitive/authen.go
	}

	// Check if user exists
	log.Println("Device authentication started")
<<<<<<< Updated upstream:NOA_BACKEND/components/sensitive/authen.go
	log.Println("Email: ", email)
	log.Println("Pass: ", pass)
	log.Println("UserID: ", userID)

	// TODO use Email find in the database
	_, err := db.FindUser(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err == nil {
		log.Println("User found")
	}

	// Hash device details
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	deviceHash := bson.M{
		"email":  email,
		"pass":   string(hashedPass),
		"userID": userID,
	}

	log.Println("Device hash: ", deviceHash)
	// Calculate the elapsed time
=======
	// log.Println("Email: ", email)
	// log.Println("Password: ", pass)
	// log.Println("DeviceID: ", deviceID)

	if email == "" || pass == "" || deviceID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := db.FindUser(email)
	if err != nil {
		log.Println("Error finding user:", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Authenticate the device using the provided email, password, and deviceID
	if _, err := db.AuthenDevice(email, pass, deviceID); err != nil {
		log.Println("Error authenticating device:", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// deviceHash := bson.M{
	// 	"email":  email,
	// 	"pass":   pass,
	// 	"userID": userID,
	// }

	// log.Println("Device hash: ", deviceHash)
	// Calculate the elapsed time

	response := map[string]string{
		"message": "Device authenticated successfully",
		"userID":  user.ID,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the elapsed time
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/sensitive/authen.go
	elapsedTime := time.Since(startTime)
	log.Printf("Device Authentication time for  %s\n", elapsedTime)
}
