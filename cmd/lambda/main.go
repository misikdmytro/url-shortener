package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/misikdmytro/url-shortener/internal/aws"
	"github.com/misikdmytro/url-shortener/internal/database"
	"github.com/misikdmytro/url-shortener/internal/handler"
	"github.com/misikdmytro/url-shortener/internal/helper"
	"github.com/misikdmytro/url-shortener/internal/server"
	"github.com/misikdmytro/url-shortener/internal/service"
)

var ginLambda *ginadapter.GinLambda

func init() {
	stageName := os.Getenv("STAGE_NAME")
	tableName := os.Getenv("TABLE_NAME")
	baseURL := os.Getenv("BASE_URL")

	db, err := aws.NewDynamoDBClient(context.Background())
	if err != nil {
		panic(err)
	}

	r := database.NewRepository(db, tableName)
	g := helper.NewRandomGeneratorFactory()
	s := service.NewURLService(r, g)
	h := handler.NewURLHandler(baseURL, s)
	e := server.NewEngine(stageName, h)

	ginLambda = ginadapter.New(e)
}

func LambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, event)
}

func main() {
	lambda.Start(LambdaHandler)
}
