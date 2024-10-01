package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"tmdb-dump/internal/api_client"
)

const (
	MOVIES_PATH = "/3/discover/movie"
)

type MovieService struct {
	apiClient *api_client.ApiClient
}

func NewMovieService(apiClient *api_client.ApiClient) *MovieService {
	return &MovieService{apiClient: apiClient}
}

func (ms *MovieService) FetchMovies(page int) (map[string]interface{}, error) {
	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))

	res, err := ms.apiClient.Get(MOVIES_PATH, &queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("response status is not ok: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	if result["status_code"] != nil {
		return nil, fmt.Errorf("invalid status_code: %v", result["status_code"])
	}

	return result, nil
}
