package api_client

import (
	"fmt"
	"net/http"
	"net/url"
)

type ApiClient struct {
	baseUrl    string
	token      string
	httpClient *http.Client
}

func NewApiClient(baseUrl string, token string, httpClient *http.Client) *ApiClient {
	return &ApiClient{baseUrl: baseUrl, token: token, httpClient: httpClient}
}

func (ac *ApiClient) Get(endpoint string, queryParams *url.Values) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s", ac.baseUrl, endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

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

	return res, nil
}
