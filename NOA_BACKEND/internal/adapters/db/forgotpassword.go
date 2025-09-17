package db

import (
	"context"
	"errors"
	"log"
	"time"

	env "GOLANG_SERVER/configs/env"
	schema "GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Login checks if the user exists and returns the user object and an error
func ForgotpasswordCheck(email string) (schema.User, error) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_USERCOLLECTION")) // Get collection user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)                            // Create a context with timeout
	defer cancel()                                                                                      // Defer cancel the context

	// Check if user exists
	filter := bson.M{"email": email}
	var result schema.User
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return schema.User{}, errors.New("user not found")
		}
		return schema.User{}, err
	}

	log.Println("ForgotPassword User found:", email)

	return result, nil
}

func ForgotpasswordNewPassword(email string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Create a context with timeout
	defer cancel()                                                           // Defer cancel the context

	// Check if user exists
	filter := bson.M{"email": email}
	var result schema.User
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	log.Println("User found:", email)

	// In Database Password is hashed
	// Compare the stored password with the input password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
