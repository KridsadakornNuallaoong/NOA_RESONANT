package db

import (
	"context"
	"fmt"
	"time"

	env "GOLANG_SERVER/configs/env"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

// * Connect to mongo db
func Connect() (bool, error) {
	clientOptions := options.Client().ApplyURI(env.GetEnv("MONGO_URI"))
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Can't connect to mongo db:", err)
		return false, err
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Can't ping mongo db:", err)
		return false, err
	}

	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_COLLECTION"))
	return true, nil
}
