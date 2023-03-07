package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
	c := NewClient(BaseAddr)

	shorten, err := c.ShortenURL(context.Background(), "https://google.com", 60*time.Second)
	require.NoError(t, err)

	resp, err := c.GetURL(context.Background(), shorten.Key)
	require.NoError(t, err)

	assert.Equal(t, "https://google.com", resp.String())
}

func TestGetURLNotFound(t *testing.T) {
	c := NewClient(BaseAddr)
	key := uuid.NewString()

	_, err := c.GetURL(context.Background(), key)
	assert.Error(t, err)

	apiErr, ok := err.(*APIError)
	assert.True(t, ok)
	assert.Equal(t, 404, apiErr.Code)
	assert.Equal(t, "URL not found", apiErr.Err)
}
