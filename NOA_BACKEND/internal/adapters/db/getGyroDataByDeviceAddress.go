package db

import (
	"context"
	"errors"
	"strings"
	"time"

	schema "GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
)

func GetGyroDataByDeviceAddress(DeviceAddress string) ([]schema.GyroData, error) {
	if len(DeviceAddress) == 0 {
		return []schema.GyroData{}, errors.New("device address is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{strings.ToLower("DeviceAddress"): DeviceAddress})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var gyroData []schema.GyroData
	if err = cursor.All(ctx, &gyroData); err != nil {
		return nil, err
	}
	return gyroData, nil
}
