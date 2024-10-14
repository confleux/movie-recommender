package service

import (
	"context"
	"fmt"

	"tg-bot/internal/model"
	"tg-bot/internal/repository"
)

type MovieService struct {
	movieRepository *repository.MovieRepository
}

func NewMovieService(movieRepository *repository.MovieRepository) *MovieService {
	return &MovieService{movieRepository: movieRepository}
}

func (ms *MovieService) GetRandomMovies(count int) ([]model.Movie, error) {
	if movies, err := ms.movieRepository.GetRandomMovies(context.Background(), count); err != nil {
		return nil, fmt.Errorf("get random movies: %w", err)
	} else {
		return movies, nil
	}
}
