package db

import (
	"context"
	"time"

	schema "GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
)

func GetGyroData() ([]schema.GyroData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
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
