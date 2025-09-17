package db

import (
	"context"
	"errors"
	"log"
	"time"

	env "GOLANG_SERVER/configs/env"
	"GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
)

// FindDevice retrieves a device from the database by its DeviceID
func FindDevice(deviceID string) (*schema.Device, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query the database for the device
	filter := bson.M{"deviceID": deviceID}
	var device schema.Device
	err := collection.FindOne(ctx, filter).Decode(&device)
	if err != nil {
		log.Println("Error finding device:", err)
		return nil, errors.New("device not found")
	}

	return &device, nil
}
