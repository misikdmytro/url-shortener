package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/misikdmytro/url-shortener/internal/database"
	"github.com/misikdmytro/url-shortener/internal/helper"
	"github.com/misikdmytro/url-shortener/internal/model"
	"github.com/sethvargo/go-retry"
)

var (
	ErrURLNotFound = fmt.Errorf("URL not found")
)

type URLService interface {
	ShortenURL(ctx context.Context, url string, duration time.Duration) (string, error)
	GetURL(ctx context.Context, key string) (string, error)
}

type urlService struct {
	repository    database.Repository
	randomFactory helper.RandomGeneratorFactory
}

func NewURLService(repository database.Repository, randomFactory helper.RandomGeneratorFactory) URLService {
	return &urlService{
		repository:    repository,
		randomFactory: randomFactory,
	}
}

type zeroBackoff struct{}

func (b *zeroBackoff) Next() (time.Duration, bool) {
	return 0, false
}

func (s *urlService) ShortenURL(ctx context.Context, url string, duration time.Duration) (string, error) {
	randomGenerator := s.randomFactory.NewRandomGenerator()
	key := randomGenerator.NewKey(8)

	err := retry.Do(
		ctx,
		retry.WithMaxRetries(3, &zeroBackoff{}),
		func(ctx context.Context) error {
			err := s.repository.SaveShortURL(ctx, model.ShortURL{
				Key: key,
				URL: url,
				Ttl: time.Now().Add(duration).Unix(),
			})

			if errors.Is(err, database.ErrDuplicateKey) {
				key = randomGenerator.NewKey(8)
				return retry.RetryableError(err)
			}

			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *urlService) GetURL(ctx context.Context, key string) (string, error) {
	url, err := s.repository.GetShortURL(ctx, key)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", ErrURLNotFound
		}

		return "", err
	}

	return url.URL, nil
}
