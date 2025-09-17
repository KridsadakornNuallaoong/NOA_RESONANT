package db

import (
	"context"
	"errors"
	"log"
	"time"

	env "GOLANG_SERVER/configs/env"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterDevice(DeviceAddress string) (bool, error) {
	if len(DeviceAddress) == 0 {
		return false, errors.New("device address is empty")
	}
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if device already exists
	filter := bson.M{"deviceaddress": DeviceAddress}

	// log.Println("Querying database with filter:", filter)

	// Check if device already exists
	var result struct {
		DeviceAddress string `bson:"deviceaddress"`
	}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == nil { // Device already exists
		return false, errors.New("device already exists")
	} else if err != mongo.ErrNoDocuments {
		return false, err
	} else {
		_, err := collection.InsertOne(ctx, bson.M{"deviceaddress": DeviceAddress})
		if err != nil {
			return false, err
		}

		log.Println("Device registered successfully.")
	}

	return true, nil
}
