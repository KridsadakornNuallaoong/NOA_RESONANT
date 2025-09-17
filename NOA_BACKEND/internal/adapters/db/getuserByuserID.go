// filepath: c:\Users\toong\Documents\GitHub\NOA_Backend\components\db\user.go
package db

import (
	"context"
	"errors"
	"time"

	env "GOLANG_SERVER/configs/env"
	"GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByID(userID string) (*schema.User, error) {
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_USERCOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userID": userID}
	var user schema.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
