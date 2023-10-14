package main

import (
	"encoding/json"
	errorlib "go-url-shortener/lib/errorLib"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	Url string `json:"url"`
}

type UrlInfo struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shorturl"`
	Domain   string `json:"domain"`
}
type service struct {
	sess     *session.Session
	s3Client *s3.S3
}

func NewService() (service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)
	if err != nil {
		return service{}, err
	}
	s3Client := s3.New(sess)
	return service{
		sess:     sess,
		s3Client: s3Client,
	}, nil
}

func SuccessResponse(resp UrlInfo) events.APIGatewayV2HTTPResponse {
	respBytes, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: 200,
	}
}

func ErrorResponse(err error) events.APIGatewayV2HTTPResponse {
	errRes := errorlib.ErrorLib{
		Message: err.Error(),
		Code:    400,
	}
	respBytes, _ := json.Marshal(errRes)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: 400,
	}
}
