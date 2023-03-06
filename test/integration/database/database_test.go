package databaase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	internalaws "github.com/misikdmytro/url-shortener/internal/aws"
	"github.com/misikdmytro/url-shortener/internal/database"
	"github.com/misikdmytro/url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tableName = "url-shortener-table"

func TestSaveShortURLOK(t *testing.T) {
	const key = "_test"

	svc, err := internalaws.NewDynamoDBClient(context.Background())
	require.NoError(t, err)

	delete := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(context.Background(), delete)
	require.NoError(t, err)

	db := database.NewRepository(svc, tableName)
	err = db.SaveShortURL(context.Background(), model.ShortURL{
		Key: key,
		URL: "https://google.com",
		Ttl: time.Now().Add(time.Minute).Unix(),
	})
	require.NoError(t, err)

	get := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName: aws.String(tableName),
	}

	res, err := svc.GetItem(context.Background(), get)
	require.NoError(t, err)

	var result model.ShortURL
	err = attributevalue.UnmarshalMap(res.Item, &result)
	require.NoError(t, err)

	require.Equal(t, key, result.Key)
}

func TestSaveShortURLFailIfKeyExist(t *testing.T) {
	const key = "_test"

	svc, err := internalaws.NewDynamoDBClient(context.Background())
	require.NoError(t, err)

	delete := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(context.Background(), delete)
	require.NoError(t, err)

	db := database.NewRepository(svc, tableName)
	err = db.SaveShortURL(context.Background(), model.ShortURL{
		Key: key,
		URL: "https://google.com",
		Ttl: time.Now().Add(time.Minute).Unix(),
	})
	require.NoError(t, err)

	err = db.SaveShortURL(context.Background(), model.ShortURL{
		Key: key,
		URL: "https://google.com",
		Ttl: time.Now().Add(time.Minute).Unix(),
	})

	assert.Equal(t, database.ErrDuplicateKey, err)
}

func TestGetShortURLAbsent(t *testing.T) {
	const key = "_test"

	svc, err := internalaws.NewDynamoDBClient(context.Background())
	require.NoError(t, err)

	delete := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(context.Background(), delete)
	require.NoError(t, err)

	db := database.NewRepository(svc, tableName)
	item, err := db.GetShortURL(context.Background(), key)

	assert.Equal(t, database.ErrItemNotFound, err)
	assert.Equal(t, model.ShortURL{}, item)
}

func TestGetShortURLDeletedAfterRead(t *testing.T) {
	const key = "_test"

	svc, err := internalaws.NewDynamoDBClient(context.Background())
	require.NoError(t, err)

	delete := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(context.Background(), delete)
	require.NoError(t, err)

	shortURL := model.ShortURL{
		Key: key,
		URL: "https://google.com",
		Ttl: time.Now().Add(time.Minute).Unix(),
	}

	db := database.NewRepository(svc, tableName)
	err = db.SaveShortURL(context.Background(), shortURL)
	require.NoError(t, err)

	item, err := db.GetShortURL(context.Background(), key)
	require.NoError(t, err)
	require.Equal(t, shortURL, item)

	get := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName: aws.String(tableName),
	}

	res, err := svc.GetItem(context.Background(), get)
	require.NoError(t, err)
	assert.Len(t, res.Item, 0)
}
