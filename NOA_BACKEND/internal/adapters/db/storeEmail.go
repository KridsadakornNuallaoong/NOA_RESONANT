package db

import (
	env "GOLANG_SERVER/configs/env"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// StoreEmail checks if an email exists in the database
func StoreEmail(email string) (bool, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_USERCOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use a BSON document for the filter
	filter := bson.M{"email": email}

	// Check if the email already exists
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, errors.New("email already exists")
	}

	// Generate user ID
	var ID = generateUserID()

	// Insert the email into the database (if needed)
	_, err = collection.InsertOne(ctx, bson.M{"userID": ID, "email": email})
	if err != nil {
		return false, err
	}

	return true, nil
}
