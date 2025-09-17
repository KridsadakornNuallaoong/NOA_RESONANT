// save OTP in  1 minute
package user

import (
	"log"

	"GOLANG_SERVER/internal/adapters/db"
)

// SaveOTP saves the OTP in a for 1 minute
func SaveOTP(email string, otp string) {
	// Find the user by email
	// Don't get password
	user, err := db.FindUser(email)
	if err != nil {
		log.Println("User not found:", email)
		return
	}

	log.Println("User found:", user.ID+" "+otp)
	// Save the OTP in the database
	db.SaveOTP(user.ID, otp)
}
