package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type ShortenURLRequest struct {
	URL      string `json:"url" binding:"required,max=2048,url"`
	Duration int64  `json:"duration" binding:"required,min=1"`
}

type ShortenURLResponse struct {
	Key string `json:"key"`
}
