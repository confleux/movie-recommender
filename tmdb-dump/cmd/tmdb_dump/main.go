package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tmdb-dump/internal/api_client"
	"tmdb-dump/internal/config"
	"tmdb-dump/internal/db/mongodb"
	mongodb_movie "tmdb-dump/internal/repository/mongodb"
)

const (
	initialBackoff     = 100 * time.Millisecond
	maxRequestRetries  = 5
	backoffCoefficient = 1.5
)

func main() {
	cfg := config.MustLoad()

	mongoDB, err := mongodb.NewMongoDB(cfg.Mongo.Uri, cfg.Mongo.Database, cfg.Mongo.Collection)
	if err != nil {
		log.Fatalf("failed to create MongoDB: %v", err)
	}

	apiClient := api_client.NewApiClient(cfg.TmdbApi.BaseUrl, cfg.TmdbApi.Token, &http.Client{})

	movieRepository := mongodb_movie.NewMovieRepository(mongoDB)

	for page := 1; page < cfg.PagesCount; page++ {
		retryCount := 0
		backoff := initialBackoff

	backoffLoop:
		for {
			select {
			case <-time.After(backoff):
				result, err := apiClient.GetMovies(page)
				if err != nil {
					retryCount++

					if retryCount > maxRequestRetries {
						log.Fatalf("%v", err)
					}

					backoff = time.Duration(float64(backoff) * backoffCoefficient)
					continue backoffLoop
				}

				retryCount = 0
				backoff = initialBackoff

				insertedId, err := movieRepository.InsertMoviesPage(result)
				if err != nil {
					log.Fatalf("failed to insert movie page: %v", err)
				}

				fmt.Printf("Successfully added %s (page: %d)\n", insertedId, page)

				break backoffLoop
			}
		}
	}
}
