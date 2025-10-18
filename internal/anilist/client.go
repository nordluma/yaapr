package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiURL = "https://graphql.anilist.co"

type Client struct {
	http  *http.Client
	token string // optional for auth
}

func NewClient(token string) *Client {
	return &Client{
		http:  &http.Client{},
		token: token,
	}
}

func (c *Client) do(
	ctx context.Context,
	query string,
	variables map[string]any,
	v any,
) error {
	body := map[string]any{
		"query":     query,
		"variables": variables,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		apiURL,
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// TODO: more descriptive error
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return json.NewDecoder(res.Body).Decode(&v)
}
