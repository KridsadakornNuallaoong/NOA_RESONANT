package db

import (
	env "GOLANG_SERVER/configs/env"
	schema "GOLANG_SERVER/pkg/schema"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckDeviceID checks if a device ID is already registered in the database
func HandlercheckDeviceID(deviceID string) (bool, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Checking device ID:", deviceID)

	if deviceID == "" {
		return false, errors.New("device ID is empty")
	}

	// Check if the device ID exists in the database
	filter := bson.M{"deviceID": deviceID}
	var result schema.Device
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Device ID " + deviceID + " does not exist:")
			return false, nil // Device ID does not exist
		}
		log.Println("Error checking device ID:", err)
		return false, err // Error occurred while checking device ID

	}

	log.Println("Device ID " + deviceID + " exists.")
	return true, nil // Device ID exists

}
