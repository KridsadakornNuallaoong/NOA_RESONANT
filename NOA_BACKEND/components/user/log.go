package user

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"GOLANG_SERVER/components/db"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Allow only POST requests
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse the request body to get user details
	var userDetails map[string]string
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle both lowercase and uppercase keys
	email := userDetails["email"]
	if email == "" {
		email = userDetails["Email"]
	}
	password := userDetails["password"]
	if password == "" {
		password = userDetails["Password"]
	}

	// Check if user exists
	user, err := db.Login(email, password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	log.Println("User logged in successfully:", user.Username, user.ID)

	// Generate a JWT token
	token, err := GenerateJWT(user.Username, user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// log.Println("Generated JWT token:", token)

	// Save ssession data to the database

	// Send a response
	response := map[string]string{
		"message": "Login successful",
		"token":   token,
	}
	log.Println("User logged in successfully.")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GenerateJWT(username, userID string) (string, error) {
	// Define the secret key (use a secure key in production)
	secretKey := []byte("User Login") // Replace with your actual secret key

	// Define the claims
	claims := jwt.MapClaims{
		"username": username,
		"userID":   userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"message":  "Login successfully",
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
