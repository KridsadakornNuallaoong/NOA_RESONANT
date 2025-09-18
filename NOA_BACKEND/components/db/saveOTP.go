package db

import (
	"context"
	"log"
	"time"

	env "GOLANG_SERVER/components/env"

	"go.mongodb.org/mongo-driver/bson"
)

// SaveOTP saves the OTP in the database
func SaveOTP(userID string, otp string) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_AUTHCOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Defer cancel the context

	// Check if user already exists
	filter := bson.M{"userID": userID}

	//log.Println("Querying database with filter:", filter)

	// Check if OTP already exists
	var result struct {
		UserID string `bson:"userID"`
		OTP    string `bson:"otp"`
	}
	if err := collection.FindOne(ctx, filter).Decode(&result); err == nil {
		// If no error, the user exists, so update the OTP
		collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"otp": otp, "expireAt": time.Now().Add(time.Minute)}})
		log.Println("OTP updated:", userID+" "+otp)

		// Schedule OTP deletion after 1 minute
		time.AfterFunc(1*time.Minute, func() {
			deleteOTP(userID)
		})
		return
	}

	// Save the OTP in the database
	_, err := collection.InsertOne(ctx, bson.M{"userID": userID, "otp": otp, "expireAt": time.Now().Add(time.Minute)})
	if err != nil {
		log.Println("Error saving OTP:", err)
		return
	}

	// Schedule OTP deletion after 1 minute
	time.AfterFunc(1*time.Minute, func() {
		deleteOTP(userID)
	})

	log.Println("OTP saved:", userID+" "+otp)
}

// deleteOTP deletes the OTP from the database
func deleteOTP(userID string) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_AUTHCOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Defer cancel the context

	// Check if user already exists
	filter := bson.M{"userID": userID}

	log.Println("Querying database with filter:", filter)

	// Check if OTP already exists
	var result struct {
		UserID string `bson:"userID"`
		OTP    string `bson:"otp"`
	}
	if collection.FindOne(ctx, filter).Decode(&result) == nil {
		// DELETE OTP only
		collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"otp": "", "expireAt": time.Now()}})
		log.Println("OTP deleted:", userID)
		return
	}

	log.Println("OTP not found:", userID)
}
