package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tg-bot/internal/model"
)

type MovieRepository struct {
	Pool *pgxpool.Pool
}

func NewMovieRepository(pool *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{Pool: pool}
}

func (mr *MovieRepository) GetRandomMovies(ctx context.Context, count int) ([]model.Movie, error) {
	query := `
		SELECT title, vote_average from movies_info
		ORDER BY RANDOM()
		LIMIT $1
  `

	res := make([]model.Movie, 0)

	rows, err := mr.Pool.Query(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("query movies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := model.Movie{}

		if err := rows.Scan(&tmp.Title, &tmp.VoteAverage); err != nil {
			return nil, fmt.Errorf("scan rows: %w", err)
		}

		res = append(res, tmp)
	}

	return res, nil
}
