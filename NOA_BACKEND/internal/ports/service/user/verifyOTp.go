// TODO: Description: This file includes the functions that are used to verify the OTP of the user.
package user

import (
	"GOLANG_SERVER/internal/adapters/db"
	"encoding/json"
	"log"
	"net/http"
)

// VerifyOTP verifies the OTP of the user with token in 1 minute
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
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
<<<<<<< Updated upstream:NOA_BACKEND/components/user/verifyOTp.go
=======
	username := userDetails["username"]
	if username == "" {
		username = userDetails["Username"]
	}
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/verifyOTp.go
	email := userDetails["email"]
	if email == "" {
		email = userDetails["Email"]
	}
<<<<<<< Updated upstream:NOA_BACKEND/components/user/verifyOTp.go
=======
	password := userDetails["password"]
	if password == "" {
		password = userDetails["Password"]
	}
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/verifyOTp.go
	otp := userDetails["otp"]
	if otp == "" {
		otp = userDetails["OTP"]
	}

	// Check if user exists
	user, err := db.FindUser(email)
	if err != nil {
		http.Error(w, "Invalid email.", http.StatusUnauthorized)
		return
	}

<<<<<<< Updated upstream:NOA_BACKEND/components/user/verifyOTp.go
	log.Println("User:", user.ID+" "+"Forget Password")
=======
	log.Println("User:", user.ID+" "+"Verifying OTP")
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/verifyOTp.go

	// Declare checkOTP and verify the OTP
	checkOTP := db.VerifyOTP(user.ID, otp)

	if checkOTP == "" {
		http.Error(w, "Invalid OTP.", http.StatusUnauthorized)
		return
	} else if checkOTP != otp {
		http.Error(w, "Invalid OTP.", http.StatusUnauthorized)
		return
	} else if checkOTP == otp {
		log.Println("OTP Verified")
<<<<<<< Updated upstream:NOA_BACKEND/components/user/verifyOTp.go
=======
		// Update the user status to verified
		StoreUser(username, email, password)
	} else {
		http.Error(w, "Invalid OTP.", http.StatusUnauthorized)
		return
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/verifyOTp.go
	}

	// Send a response
	response := map[string]string{"message": "OTP verified"}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
<<<<<<< Updated upstream:NOA_BACKEND/components/user/verifyOTp.go
	// Verify the OTP

=======
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/service/user/verifyOTp.go
}
