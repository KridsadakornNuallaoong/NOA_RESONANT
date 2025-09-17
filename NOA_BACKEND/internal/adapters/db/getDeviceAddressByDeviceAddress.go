package db

import (
	"context"
	"log"
	"time"

	env "GOLANG_SERVER/configs/env"

	"go.mongodb.org/mongo-driver/bson"
)

func GetDeviceAddressByDeviceAddress(deviceAddress string) ([]string, error) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Create a context with timeout
	defer cancel()

	filter := bson.M{"deviceaddress": deviceAddress}
	log.Println("Querying database with filter:", filter)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deviceAddresses []string
	for cursor.Next(ctx) {
		var result struct {
			DeviceAddress string `bson:"deviceaddress"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		deviceAddresses = append(deviceAddresses, result.DeviceAddress)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	log.Println("Found device addresses:", deviceAddresses)
	return deviceAddresses, nil
}
