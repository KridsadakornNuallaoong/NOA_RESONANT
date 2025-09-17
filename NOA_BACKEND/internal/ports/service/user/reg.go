package user

import (
	"encoding/json"
	"log"
	"net/http"
<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
=======
	"regexp"
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go

	"GOLANG_SERVER/internal/adapters/db"
	"GOLANG_SERVER/pkg/schema"

	"golang.org/x/crypto/bcrypt"
)

<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
=======
// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	// Regular expression for validating email format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidatePassword(password string) bool {
	// Check if the password is at least 8 characters long
	if len(password) < 8 {
		return false
	}
	// Check if the password contains at least one letter
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)

	// Check if the password contains at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	// The password is valid if it has both a letter and a digit
	return hasLetter && hasDigit
}

>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go
// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
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
<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
=======
	username := userDetails["username"]
	if username == "" {
		username = userDetails["Username"]
	}
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go
	email := userDetails["email"]
	if email == "" {
		email = userDetails["Email"]
	}
	password := userDetails["password"]
	if password == "" {
		password = userDetails["Password"]
	}

<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
	// Declare otp variable outside the if block
	var otp string
=======
	if !ValidateEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !ValidatePassword(password) {
		http.Error(w, "Password must be at least 8 characters long and contain at least one letter and one number", http.StatusBadRequest)
		return
	}
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the user details to a User struct
<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
	user := schema.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save user details to database
	if _, err := db.StoreUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Println("User registered successfully.")
		// Generate OTP and assign it to the otp variable
		otp = GenerateOTP()
		SendOTPEmail(email, otp)
		// Save the OTP in the database
		SaveOTP(email, otp)
=======
	// user := schema.User{
	//     Username: username,
	//     Email:    email,
	//     Password: string(hashedPassword),
	// }

	// Use = instead of := since err is already declared
	_, err = db.StoreEmail(email) // Check if email already exists
	if err != nil {
		if err.Error() == "email already exists" {
			http.Error(w, "Email already exists", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else {
		log.Println("Email registered successfully.")
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go
	}

	// Send a response
	response := map[string]string{
<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
		"message": "User registered successfully. Please check your email for the OTP.",
		"otp":     otp,
=======
		"message":  "User registered successfully. Please check your email for the OTP.",
		"username": username,
		"email":    email,
		"password": string(hashedPassword),
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go
	}
	log.Println("User registered successfully.")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
<<<<<<< Updated upstream:NOA_BACKEND/components/user/reg.go
=======

// function StoreUser to database after verifying the OTP
func StoreUser(username string, email string, password string) error {

	log.Println("Storing user:", username, email, password)

	// Check emtpy email and password
	if email == "" || password == "" || username == "" {
		log.Println("Email, password, or username is empty.")
		return nil
	}

	// Convert the user details to a User struct
	user := schema.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Store the user in the database
	if _, err := db.StoreUser(user); err != nil {
		log.Println("Error storing user:", err)
	}
	log.Println("User stored successfully.")
	return nil

}
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/reg.go
