package server

import (
	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/url-shortener/internal/handler"
)

func NewEngine(base string, h handler.URLHandler) *gin.Engine {
	r := gin.Default()

	g := r.Group("/" + base)
	{
		g.GET("/:key", h.ToURL)
		g.PUT("/shorten", h.ShortenURL)
	}

	return r
}
