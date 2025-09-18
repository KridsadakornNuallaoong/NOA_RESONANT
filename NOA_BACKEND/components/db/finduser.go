package db

import (
	"context"
	"time"

	env "GOLANG_SERVER/components/env"
	schema "GOLANG_SERVER/components/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// find user by email
func FindUser(email string) (schema.User, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_USERCOLLECTION")) // Get collection user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)                             // Create a context with timeout
	defer cancel()                                                                                       // Defer cancel the context

	// Check if user exists
	var result schema.User
	filter := bson.M{"email": email}
	projection := bson.M{"password": 0} // Exclude the password field

	// Use options.FindOne() to set the projection
	findOptions := options.FindOne().SetProjection(projection)

	// Find the user by email
	err := collection.FindOne(ctx, filter, findOptions).Decode(&result)
	if err != nil {
		return schema.User{}, err
	}

	return result, nil
}
