package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
	c := NewClient(BaseAddr())

	resp, err := c.ShortenURL(context.Background(), "https://google.com", 60*time.Second)
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Key)
	assert.Equal(t, fmt.Sprintf("%s/%s", BaseAddr(), resp.Key), resp.URL)
}

func TestShortenURLTwice(t *testing.T) {
	c := NewClient(BaseAddr())

	resp1, err1 := c.ShortenURL(context.Background(), "https://google.com", 60*time.Second)
	resp2, err2 := c.ShortenURL(context.Background(), "https://google.com", 60*time.Second)
	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.NotEqual(t, resp1.Key, resp2.Key)
	assert.NotEqual(t, resp1.URL, resp2.URL)
}

func TestShortenURLValidationURL(t *testing.T) {
	input := []string{"test", "h", ".com", ""}

	c := NewClient(BaseAddr())
	for _, tc := range input {
		t.Run(fmt.Sprintf("validation test '%s'", tc), func(t *testing.T) {
			_, err := c.ShortenURL(context.Background(), tc, 60*time.Second)
			assert.Error(t, err)

			apiErr, ok := err.(*APIError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, apiErr.Code)
		})
	}
}

func TestShortenURLValidationDuration(t *testing.T) {
	input := []time.Duration{0 * time.Second, -1 * time.Second, 604801 * time.Second}

	c := NewClient(BaseAddr())
	for _, tc := range input {
		t.Run(fmt.Sprintf("validation test '%s'", tc), func(t *testing.T) {
			_, err := c.ShortenURL(context.Background(), "https://google.com", tc)
			assert.Error(t, err)

			apiErr, ok := err.(*APIError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, apiErr.Code)
		})
	}
}
