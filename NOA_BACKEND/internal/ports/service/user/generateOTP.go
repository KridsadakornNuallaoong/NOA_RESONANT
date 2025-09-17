package user

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a random 6-digit OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(10000)
	return fmt.Sprintf("%04d", otp)
}
