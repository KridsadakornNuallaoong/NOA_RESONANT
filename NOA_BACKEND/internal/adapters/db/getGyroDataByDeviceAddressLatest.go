package db

import (
	"context"
	"errors"
	"strings"
	"time"

	env "GOLANG_SERVER/configs/env"
	schema "GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetGyroDataByDeviceAddressLatest(DeviceAddress string) ([]schema.GyroData, error) {
	if len(DeviceAddress) == 0 {
		return nil, errors.New("device address is empty")
	}
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_COLLECTION"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var gyroData []schema.GyroData
	cursor, err := collection.Find(ctx, bson.M{strings.ToLower("deviceaddress"): DeviceAddress}, options.Find().SetSort(bson.D{{Key: strings.ToLower("timestamp"), Value: -1}}).SetLimit(50))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &gyroData); err != nil {
		return nil, err
	}

	if len(gyroData) == 0 {
		return nil, errors.New("no data found")
	}
	return gyroData, nil
}
