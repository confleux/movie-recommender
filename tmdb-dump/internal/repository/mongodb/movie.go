package mongodb_movie

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"

	"tmdb-dump/internal/api_client"
	"tmdb-dump/internal/db/mongodb"
)

type MovieRepository struct {
	mongoConnection *mongodb.MongoConnection
}

func NewMovieRepository(mongoConnection *mongodb.MongoConnection) *MovieRepository {
	return &MovieRepository{mongoConnection: mongoConnection}
}

func (mr *MovieRepository) InsertMoviesPage(moviePage *api_client.GetMoviesResponse) (interface{}, error) {
	doc, err := bson.Marshal(moviePage)
	if err != nil {
		return "", fmt.Errorf("marshal to bson: %w", err)
	}

	insertRes, err := mr.mongoConnection.Collection.InsertOne(context.Background(), doc)
	if err != nil {
		return "", fmt.Errorf("insert doc: %w", err)
	}

	return insertRes.InsertedID, nil
}
