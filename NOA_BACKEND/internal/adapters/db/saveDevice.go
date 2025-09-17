package db

import (
	env "GOLANG_SERVER/configs/env"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// SaveDevice saves a new device to the database
func SaveDevice(deviceID, userID, devicePassword string) error {
	if len(deviceID) == 0 {
		return errors.New("device ID is empty")
	}
	if len(userID) == 0 {
		return errors.New("device email is empty")
	}
	if len(devicePassword) == 0 {
		return errors.New("device password is empty")
	}

	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the new device into the database
	_, err := collection.InsertOne(ctx, bson.M{
		"deviceID": deviceID,
		"userID":   userID,
		"password": devicePassword,
	})
	if err != nil {
		log.Println("Error saving device:", err)
		return err
	}

	log.Println("Device saved successfully:", deviceID)
	return nil
}
