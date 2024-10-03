package api_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	getMoviesEndpoint = "/3/discover/movie"
)

type ApiClient struct {
	baseUrl    string
	token      string
	httpClient *http.Client
}

type GetMoviesResponse struct {
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
}

type Movie struct {
	Id               int     `json:"id"`
	VoteAverage      float64 `json:"vote_average"`
	GenreIds         []int   `json:"genre_ids"`
	OriginalTitle    string  `json:"original_title"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteCount        int     `json:"vote_count"`
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	OriginalLanguage string  `json:"original_language"`
	Popularity       float64 `json:"popularity"`
	Title            string  `json:"title"`
}

func NewApiClient(baseUrl string, token string, httpClient *http.Client) *ApiClient {
	return &ApiClient{baseUrl: baseUrl, token: token, httpClient: httpClient}
}

func (ac *ApiClient) GetMovies(page int) (*GetMoviesResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, getMoviesEndpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))

	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ac.token))

	res, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status is not ok: %d", res.StatusCode)
	}

	var response GetMoviesResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy to discard: %w", err)
	}

	return &response, nil
}
