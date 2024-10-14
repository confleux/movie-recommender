package model

import "time"

type Movie struct {
	Id          int
	Title       string
	VoteAverage float32
	ReleaseDate time.Time
	PosterPath  string
}
