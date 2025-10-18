package allanime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	allanimeBaseUrl   = "https://allanime.day"
	allanimeApiUrl    = "https://api.allanime.day/api"
	allanimeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
	allanimeReferer   = "https://allanime.to"
)

type AllanimeClient struct {
	client *http.Client
}

func NewAllanimeClient() *AllanimeClient {
	return &AllanimeClient{
		client: &http.Client{},
	}
}

func (c AllanimeClient) do(
	ctx context.Context,
	query string,
	variables, v any,
) error {
	data, err := json.Marshal(variables)
	if err != nil {
		return err
	}

	urlParams := url.Values{}
	urlParams.Set("query", query)
	urlParams.Set("variables", string(data))

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s?%s", allanimeApiUrl, urlParams.Encode()),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", allanimeUserAgent)
	req.Header.Set("Referer", allanimeReferer)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("non-OK response (%d): %s", res.StatusCode, body)
	}

	return json.NewDecoder(res.Body).Decode(&v)
}
