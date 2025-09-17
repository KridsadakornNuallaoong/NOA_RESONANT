package db

import (
	"context"
	"errors"
	"log"
	"time"

	env "GOLANG_SERVER/configs/env"
	schema "GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// AuthenDevice function to authenticate device
func AuthenDevice(email string, pass string, deviceID string) (bool, error) {
	// Check if user exists
	// log.Println("Device authentication started")
	// log.Println("Email: ", email)
	// log.Println("Pass: ", pass)
	// log.Println("DeviceID: ", deviceID)

	// Check Email, Password, and id in the database
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION")) // Get collection user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)                               // Create a context with timeout
	defer cancel()                                                                                         // Defer cancel the context

	if email == "" || pass == "" || deviceID == "" {
		log.Println("Missing required fields")
		return false, errors.New("missing required fields")
	}

	// find the user by email
	user, err := FindUser(email)
	if err != nil {
		log.Println("Error finding user:", err)
		return false, err
	} else {
		log.Println("User ID is OK")
	}

	log.Println("User ID:", user.ID)

	filter := bson.M{"deviceID": deviceID, "userID": user.ID}
	var device schema.Device
	if err := collection.FindOne(ctx, filter).Decode(&device); err != nil {
		return false, err
	} else {
		log.Println("Device ID is OK")
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(device.Password), []byte(pass)); err != nil {
		log.Println("Error comparing password:", err)
		return false, err
	} else {
		log.Println("Password is OK")
	}

	log.Println("Device authentication successful")
	return true, nil
}
