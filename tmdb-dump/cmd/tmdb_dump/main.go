package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"tmdb-dump/internal/api_client"
	"tmdb-dump/internal/config"
	"tmdb-dump/internal/repository/postgres"
)

const (
	initialBackoff     = 100 * time.Millisecond
	maxRequestRetries  = 10
	backoffCoefficient = 1.5
)

func main() {
	cfg := config.MustLoad()

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.Uri)
	if err != nil {
		log.Fatalf("postges connection: %v", err)
	}
	defer pool.Close()

	apiClient := api_client.NewApiClient(cfg.TmdbApi.BaseUrl, cfg.TmdbApi.Token, &http.Client{})

	movieRepository := postgres.NewMovieRepository(pool)

	start := time.Now()

	for page := 1; page <= cfg.PagesCount; page++ {
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
						log.Fatalf("max request retries reached: %v", err)
					}

					backoff = time.Duration(float64(backoff) * backoffCoefficient)
					continue backoffLoop
				}

				retryCount = 0
				backoff = initialBackoff

				for _, v := range result.Results {
					id, err := movieRepository.InsertMovie(context.Background(), v)

					if err != nil {
						log.Fatalf("Failed to insert movie: %v", err)
					}

					log.Printf("Inserted movie with id: %d\n", id)
				}
				log.Printf("Processed page: %d (%d movies)\n", page, len(result.Results))

				break backoffLoop
			}
		}
	}

	elapsed := time.Since(start)

	log.Printf("TMDB Dump took: %s", elapsed)
}
