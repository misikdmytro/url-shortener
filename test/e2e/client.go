package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/misikdmytro/url-shortener/pkg/model"
)

type Client interface {
	ShortenURL(ctx context.Context, url string, duration time.Duration) (model.ShortenURLResponse, error)
	GetURL(ctx context.Context, key string) (*url.URL, error)
}

type client struct {
	c       http.Client
	baseURL string
}

type APIError struct {
	Err  string
	Code int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("unexpected status code: %d. details: %s", e.Code, e.Err)
}

func NewClient(baseURL string) Client {
	return &client{
		c: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		baseURL: baseURL,
	}
}

func (c *client) ShortenURL(ctx context.Context, url string, duration time.Duration) (model.ShortenURLResponse, error) {
	input := model.ShortenURLRequest{
		URL:      url,
		Duration: int64(duration.Seconds()),
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return model.ShortenURLResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/shorten", c.baseURL), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return model.ShortenURLResponse{}, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return model.ShortenURLResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var err model.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
			return model.ShortenURLResponse{}, fmt.Errorf("failed to decode response: %w", err)
		}

		return model.ShortenURLResponse{}, &APIError{
			Err:  err.Error,
			Code: resp.StatusCode,
		}
	}

	var output model.ShortenURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return model.ShortenURLResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return output, nil
}

func (c *client) GetURL(ctx context.Context, key string) (*url.URL, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.baseURL, key), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMovedPermanently {
		var err model.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
			return nil, err
		}

		return nil, &APIError{
			Err:  err.Error,
			Code: resp.StatusCode,
		}
	}

	return resp.Location()
}
