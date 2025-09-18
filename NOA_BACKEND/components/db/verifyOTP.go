package db

import (
	"context"
	"log"
	"strings"
	"time"

	env "GOLANG_SERVER/components/env"

	"go.mongodb.org/mongo-driver/bson"
)

// VerifyOTP verifies the OTP of the user with token in 1 minute
func VerifyOTP(userID string, otp string) string {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_AUTHCOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Defer cancel the context

	// Check if user already exists
	filter := bson.M{"userID": userID}

	// Check if OTP already exists
	var result struct {
		UserID string `bson:"userID"`
		OTP    string `bson:"otp"`
	}
	// Verify the OTP
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	// Trim and return the OTP
	result.OTP = strings.TrimSpace(result.OTP)
	return result.OTP
}
