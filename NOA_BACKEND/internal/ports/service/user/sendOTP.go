package user

import (
	"encoding/json"
	"net/http"
)

// SendOTP sends an OTP to the user's email and returns the OTP
func SendOTP(w http.ResponseWriter, r *http.Request) {

<<<<<<< Updated upstream:NOA_BACKEND/components/user/sendOTP.go
	if r.Method != http.MethodGet { // Allow only POST requests
=======
	if r.Method != http.MethodPost { // Allow only POST requests
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/sendOTP.go
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

	// Generate OTP
	otp := GenerateOTP()

	// Send OTP to user's email
	if err := SendOTPEmail(email, otp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		SaveOTP(email, otp)
	}

	// Send a response
	response := map[string]string{"message": "OTP sent successfully. Please check your email for the OTP.", "OTP": otp}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
