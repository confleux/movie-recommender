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

func NewApiClient(baseUrl string, token string, httpClient *http.Client) *ApiClient {
	return &ApiClient{baseUrl: baseUrl, token: token, httpClient: httpClient}
}

func (ac *ApiClient) GetMovies(page int) (map[string]interface{}, error) {
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy to discard: %w", err)
	}

	var moviesPage map[string]interface{}
	err = json.Unmarshal(body, &moviesPage)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	if moviesPage["status_code"] != nil {
		return nil, fmt.Errorf("invalid status_code: %v", moviesPage["status_code"])
	}

	return moviesPage, nil
}
