package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/misikdmytro/url-shortener/internal/model"
)

var (
	ErrDuplicateKey = fmt.Errorf("duplicate key")
	ErrItemNotFound = fmt.Errorf("item not found")
)

type Repository interface {
	SaveShortURL(ctx context.Context, shortURL model.ShortURL) error
	GetShortURL(ctx context.Context, key string) (model.ShortURL, error)
}

type repository struct {
	svc       *dynamodb.Client
	tableName string
}

func NewRepository(svc *dynamodb.Client, tableName string) Repository {
	return &repository{
		svc:       svc,
		tableName: tableName,
	}
}

func (r *repository) SaveShortURL(ctx context.Context, item model.ShortURL) error {
	data, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	condBuilder := expression.Name("Key").AttributeNotExists()
	expr, err := expression.NewBuilder().WithCondition(condBuilder).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                      data,
		TableName:                 aws.String(r.tableName),
		ReturnValues:              types.ReturnValueNone,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	_, err = r.svc.PutItem(ctx, input)
	if err != nil {
		var cerr *types.ConditionalCheckFailedException
		if errors.As(err, &cerr) {
			return ErrDuplicateKey
		}

		return err
	}

	return nil
}

func (r *repository) GetShortURL(ctx context.Context, key string) (model.ShortURL, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
		TableName:    aws.String(r.tableName),
		ReturnValues: types.ReturnValueAllOld,
	}

	res, err := r.svc.DeleteItem(ctx, input)
	if err != nil {
		return model.ShortURL{}, err
	}

	if res.Attributes == nil {
		return model.ShortURL{}, ErrItemNotFound
	}

	var result model.ShortURL
	err = attributevalue.UnmarshalMap(res.Attributes, &result)
	if err != nil {
		return model.ShortURL{}, err
	}

	return result, nil
}
