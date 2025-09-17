package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CleanData() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	return true, nil
}
