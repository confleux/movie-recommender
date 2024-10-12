package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tmdb-dump/internal/api_client"
)

type MovieRepository struct {
	Pool *pgxpool.Pool
}

func NewMovieRepository(pool *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{Pool: pool}
}

func (mr *MovieRepository) InsertMovie(ctx context.Context, movie api_client.Movie) (int, error) {
	query := `
		INSERT INTO movies_info (id, title, release_date, vote_average, vote_count, is_adult, poster_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int

	if err := mr.Pool.QueryRow(ctx, query, movie.Id, movie.Title, movie.ReleaseDate, movie.VoteAverage, movie.VoteCount, movie.Adult, movie.PosterPath).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert movie: %w", err)
	}

	for _, v := range movie.GenreIds {
		query = `
			INSERT INTO movies_genre_map (movie_id, genre_id)
			VALUES ($1, $2)
		`

		if _, err := mr.Pool.Exec(ctx, query, movie.Id, v); err != nil {
			return 0, fmt.Errorf("insert movie genre: %w", err)
		}
	}

	return id, nil
}
