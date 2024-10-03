package mongodb_movie

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"

	"tmdb-dump/internal/api_client"
	"tmdb-dump/internal/db/mongodb"
)

type MovieRepository struct {
	mongoDB *mongodb.MongoDB
}

func NewMovieRepository(mongodb *mongodb.MongoDB) *MovieRepository {
	return &MovieRepository{mongoDB: mongodb}
}

func (mr *MovieRepository) InsertMoviesPage(moviePage *api_client.GetMoviesResponse) (interface{}, error) {
	doc, err := bson.Marshal(moviePage)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to bson: %w", err)
	}

	insertRes, err := mr.mongoDB.Collection.InsertOne(context.Background(), doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert doc: %w", err)
	}

	return insertRes.InsertedID, nil
}
