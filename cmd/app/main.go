package main

import (
	"context"

	"github.com/misikdmytro/url-shortener/internal/aws"
	"github.com/misikdmytro/url-shortener/internal/database"
	"github.com/misikdmytro/url-shortener/internal/handler"
	"github.com/misikdmytro/url-shortener/internal/helper"
	"github.com/misikdmytro/url-shortener/internal/server"
	"github.com/misikdmytro/url-shortener/internal/service"
)

func main() {
	db, err := aws.NewDynamoDBClient(context.Background())
	if err != nil {
		panic(err)
	}

	r := database.NewRepository(db, "url-shortener-table")
	g := helper.NewRandomGeneratorFactory()
	s := service.NewURLService(r, g)
	h := handler.NewURLHandler("http://localhost:4000", s)
	e := server.NewEngine("", h)

	if err := e.Run(":4000"); err != nil {
		panic(err)
	}
}
