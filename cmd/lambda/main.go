package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/misikdmytro/url-shortener/internal/dependency"
)

var ginLambda *ginadapter.GinLambda

func init() {
	stageName := os.Getenv("STAGE_NAME")
	tableName := os.Getenv("TABLE_NAME")
	baseURL := os.Getenv("BASE_URL")

	d, err := dependency.NewDependencies(context.Background(), baseURL, stageName, tableName)
	if err != nil {
		panic(err)
	}

	ginLambda = ginadapter.New(d.Engine)
}

func LambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, event)
}

func main() {
	lambda.Start(LambdaHandler)
}
