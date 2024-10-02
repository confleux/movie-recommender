package main

import (
	"fmt"
	"log"
	"time"
	"tmdb-dump/internal/api_client"
	"tmdb-dump/internal/config"
	"tmdb-dump/internal/db/mongodb"
	mongodb_movie "tmdb-dump/internal/repository/mongodb"
	"tmdb-dump/internal/service"
)

func main() {
	cfg := config.MustLoad()

	MongoDB, err := mongodb.NewMongoDB(cfg.Mongo.Uri, cfg.Mongo.Database, cfg.Mongo.Collection)
	if err != nil {
		log.Fatalf("failed to create MongoDB: %w", err)
	}

	ApiClient := api_client.NewApiClient(cfg.TmdbApi.BaseUrl, cfg.TmdbApi.Token)

	MovieRepository := mongodb_movie.NewMovieRepository(MongoDB)

	MovieService := service.NewMovieService(ApiClient)

	for page := 1; page < cfg.PagesCount; page++ {
		result, err := MovieService.FetchMovies(page)
		if err != nil {
			log.Fatalf("failed to fetch movies: %v", err)
		}

		insertedId, err := MovieRepository.InsertMoviesPage(result)
		if err != nil {
			log.Fatalf("failed to insert movie page: %v", err)
		}

		fmt.Printf("Successfully added %s (page: %d)", insertedId, page)

		time.Sleep(2 * time.Second) // We avoid 429
	}
}
