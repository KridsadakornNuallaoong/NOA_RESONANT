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

// FindUserID retrieves a user from the database by their UserID
func FindUserID(userID string) (*schema.User, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_USERCOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query the database for the user
	filter := bson.M{"userID": userID}
	var user schema.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println("Error finding user:", err)
		return nil, errors.New("user not found")
	}

	return &user, nil
}
