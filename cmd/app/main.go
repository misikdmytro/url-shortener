package main

import (
	"context"

	"github.com/misikdmytro/url-shortener/internal/dependency"
)

func main() {
	d, err := dependency.NewDependencies(context.Background(), "http://localhost:4000", "", "url-shortener-table")
	if err != nil {
		panic(err)
	}

	if err := d.Engine.Run(":4000"); err != nil {
		panic(err)
	}
}
