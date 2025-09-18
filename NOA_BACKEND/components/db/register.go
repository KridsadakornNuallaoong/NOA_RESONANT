package db

import (
	"context"
	"errors"
	"time"

	env "GOLANG_SERVER/components/env"
	schema "GOLANG_SERVER/components/schema"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Store Email and Password to mongoDB collection user
func StoreUser(user schema.User) (bool, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_USERCOLLECTION")) // Get collection user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)                             // Create a context with timeout
	defer cancel()                                                                                       // Defer cancel the context

	// Generate user ID
	user.ID = generateUserID()

	// Prepare document to insert
	doc := bson.M{
		"id":       user.ID,
		"email":    user.Email,
		"password": user.Password, // consider hashing (bcrypt) before storing
	}

	// Check if user already exists
	filter := bson.M{"email": user.Email}
	var existing schema.User
	err := collection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		// Found an existing user with same email
		return false, errors.New("email already exists")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		// Some unexpected DB error
		return false, err
	}

	// Insert new user (no existing document)
	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return false, err
	}

	return true, nil
}

// generateUserID generates a unique user ID
func generateUserID() string {
	return uuid.New().String() // Generate a new UUID
}
