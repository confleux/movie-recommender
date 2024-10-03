package mongodb

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoDB(uri string, databaseName string, collectionName string) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	collection := client.Database(databaseName).Collection(collectionName)

	return &MongoDB{Client: client, Collection: collection}, nil
}
