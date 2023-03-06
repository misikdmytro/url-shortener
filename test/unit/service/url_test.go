package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/misikdmytro/url-shortener/internal/database"
	"github.com/misikdmytro/url-shortener/internal/helper"
	"github.com/misikdmytro/url-shortener/internal/model"
	"github.com/misikdmytro/url-shortener/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) SaveShortURL(ctx context.Context, shortURL model.ShortURL) error {
	args := r.Called(ctx, shortURL)
	return args.Error(0)
}

func (r *repositoryMock) GetShortURL(ctx context.Context, key string) (model.ShortURL, error) {
	args := r.Called(ctx, key)
	return args.Get(0).(model.ShortURL), args.Error(1)
}

type randomGeneratorMock struct {
	mock.Mock
}

type randomGeneratorFactory struct {
	g helper.RandomGenerator
}

func (r *randomGeneratorFactory) NewRandomGenerator() helper.RandomGenerator {
	return r.g
}

func (r *randomGeneratorMock) NewKey(length int) string {
	args := r.Called(length)
	return args.String(0)
}

var (
	_ (database.Repository)           = (*repositoryMock)(nil)
	_ (helper.RandomGenerator)        = (*randomGeneratorMock)(nil)
	_ (helper.RandomGeneratorFactory) = (*randomGeneratorFactory)(nil)
)

func TestShortenURL(t *testing.T) {
	input := []struct {
		name        string
		key         string
		dbErr       error
		expectedKey string
		expectedErr error
	}{
		{
			"success",
			"12345678",
			nil,
			"12345678",
			nil,
		},
		{
			"db error",
			"12345678",
			fmt.Errorf("some error"),
			"",
			fmt.Errorf("some error"),
		},
	}

	for _, tc := range input {
		t.Run(tc.name, func(t *testing.T) {
			repository := &repositoryMock{}
			repository.On("SaveShortURL", mock.Anything, mock.MatchedBy(func(s model.ShortURL) bool { return s.Key == tc.key })).Return(tc.dbErr)

			randomGenerator := &randomGeneratorMock{}
			randomGeneratorFactory := &randomGeneratorFactory{
				g: randomGenerator,
			}
			randomGenerator.On("NewKey", mock.Anything).Return(tc.key)

			service := service.NewURLService(repository, randomGeneratorFactory)
			key, err := service.ShortenURL(context.Background(), "https://google.com", 0)

			assert.Equal(t, key, tc.expectedKey)
			assert.Equal(t, err, tc.expectedErr)

			repository.AssertNumberOfCalls(t, "SaveShortURL", 1)
		})
	}
}

func TestShortenURLRetryOnce(t *testing.T) {
	key1 := "12345678"
	key2 := "87654321"

	repository := &repositoryMock{}
	repository.On("SaveShortURL", mock.Anything, mock.Anything).Return(database.ErrDuplicateKey).Times(1)
	repository.On("SaveShortURL", mock.Anything, mock.Anything).Return(nil).Times(1)

	randomGenerator := &randomGeneratorMock{}
	randomGeneratorFactory := &randomGeneratorFactory{
		g: randomGenerator,
	}

	randomGenerator.On("NewKey", mock.Anything).Return(key1).Times(1)
	randomGenerator.On("NewKey", mock.Anything).Return(key2).Times(1)

	service := service.NewURLService(repository, randomGeneratorFactory)
	key, err := service.ShortenURL(context.Background(), "https://google.com", 0)
	require.NoError(t, err)
	assert.Equal(t, key, key2)
}

func TestGetURL(t *testing.T) {
	input := []struct {
		name        string
		key         string
		dbResult    model.ShortURL
		dbErr       error
		expectedURL string
		expectedErr error
	}{
		{
			"success",
			"12345678",
			model.ShortURL{
				Key: "12345678",
				URL: "https://google.com",
			},
			nil,
			"https://google.com",
			nil,
		},
		{
			"db error",
			"12345678",
			model.ShortURL{},
			fmt.Errorf("some error"),
			"",
			fmt.Errorf("some error"),
		},
		{
			"not found",
			"12345678",
			model.ShortURL{},
			database.ErrItemNotFound,
			"",
			service.ErrURLNotFound,
		},
	}

	for _, tc := range input {
		t.Run(tc.name, func(t *testing.T) {
			repository := &repositoryMock{}
			repository.On("GetShortURL", mock.Anything, tc.key).Return(tc.dbResult, tc.dbErr)

			service := service.NewURLService(repository, nil)
			url, err := service.GetURL(context.Background(), tc.key)

			assert.Equal(t, url, tc.expectedURL)
			assert.Equal(t, err, tc.expectedErr)

			repository.AssertNumberOfCalls(t, "GetShortURL", 1)
		})
	}
}
