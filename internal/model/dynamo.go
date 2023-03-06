package model

type ShortURL struct {
	Key string `dynamodbav:"Key"`
	URL string `dynamodbav:"Url"`
	Ttl int64  `dynamodbav:"Ttl"`
}
