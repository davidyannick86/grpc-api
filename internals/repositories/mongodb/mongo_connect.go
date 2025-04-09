package mongodb

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient() (*mongo.Client, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to connect to MongoDB")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to ping MongoDB")
	}

	utils.PrintHanlder("Connected to MongoDB successfully")

	return client, nil
}
