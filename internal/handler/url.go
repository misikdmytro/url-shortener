package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/url-shortener/internal/service"
	"github.com/misikdmytro/url-shortener/pkg/model"
)

type URLHandler interface {
	ShortenURL(ctx *gin.Context)
	ToURL(ctx *gin.Context)
}

type urlHandler struct {
	base string
	srvc service.URLService
}

func NewURLHandler(base string, srvc service.URLService) URLHandler {
	return &urlHandler{
		base: base,
		srvc: srvc,
	}
}

func (h *urlHandler) ShortenURL(ctx *gin.Context) {
	var input model.ShortenURLRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	key, err := h.srvc.ShortenURL(ctx, input.URL, time.Duration(input.Duration)*time.Second)
	if err != nil {
		log.Printf("failed to shorten url: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, model.ShortenURLResponse{Key: fmt.Sprintf("%s/%s", h.base, key)})
}

func (h *urlHandler) ToURL(ctx *gin.Context) {
	key := ctx.Param("key")

	url, err := h.srvc.GetURL(ctx, key)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "URL not found"})
		} else {
			log.Printf("failed to get url: %v", err)
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	ctx.Redirect(http.StatusMovedPermanently, url)
}
