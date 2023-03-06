package dependency

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/url-shortener/internal/aws"
	"github.com/misikdmytro/url-shortener/internal/database"
	"github.com/misikdmytro/url-shortener/internal/handler"
	"github.com/misikdmytro/url-shortener/internal/helper"
	"github.com/misikdmytro/url-shortener/internal/server"
	"github.com/misikdmytro/url-shortener/internal/service"
)

type Dependencies struct {
	Engine *gin.Engine
}

func NewDependencies(ctx context.Context, baseURL, basePath, tableName string) (Dependencies, error) {
	db, err := aws.NewDynamoDBClient(ctx)
	if err != nil {
		return Dependencies{}, err
	}

	r := database.NewRepository(db, tableName)
	g := helper.NewRandomGeneratorFactory()
	s := service.NewURLService(r, g)
	h := handler.NewURLHandler(baseURL, s)
	e := server.NewEngine(basePath, h)

	d := Dependencies{
		Engine: e,
	}

	return d, nil
}
