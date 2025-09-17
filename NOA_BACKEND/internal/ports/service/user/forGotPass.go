package user

import (
	"encoding/json"
	"net/http"

	"GOLANG_SERVER/internal/adapters/db"
)

// ForgotPassword sends an OTP to the user's email
func ForgotPasswordReq(w http.ResponseWriter, r *http.Request) {
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

	// Check if user exists
	_, err := db.ForgotpasswordCheck(email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Send OTP to the user's email
	otp := GenerateOTP()
	SendOTPEmail(email, otp)

	// save OTP
	SaveOTP(email, otp)

	// Send a response
	response := map[string]string{"message": "Login successful", "otp": otp}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
