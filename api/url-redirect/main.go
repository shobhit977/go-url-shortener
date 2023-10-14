package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-url-shortener/lib/constants"
	s3service "go-url-shortener/lib/s3-service"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	svc, err := NewService()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body: err.Error(),
		}, nil
	}
	redirectUrl, err := redirect(svc, req)
	if err != nil {
		return ErrorResponse(err), nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": redirectUrl,
		},
		IsBase64Encoded: false,
	}, nil
}

func redirect(svc service, req events.APIGatewayV2HTTPRequest) (string, error) {
	shortUrl := strings.Split(req.RawPath, "/redirect")[1]
	fmt.Println(shortUrl)
	isFileExist, err := s3service.KeyExists(constants.Bucket, constants.Key, svc.s3Client)
	if err != nil {
		return "", err
	}
	if isFileExist {
		existingInfo, err := s3service.GetS3Object(svc.s3Client, constants.Bucket, constants.Key)
		if err != nil {
			return "", err
		}
		var urlDetails []UrlInfo
		if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
			return "", err
		}
		redirectUrl, isExist := isUrlExist(urlDetails, shortUrl)
		if !isExist {
			return "", errors.New("short URL does not exist. Please specify a valid shortUrl")
		}
		return redirectUrl, nil
	} else {
		return "", errors.New("short URL does not exist.Please specify a valid shortUr")
	}
}
func isUrlExist(urlDetails []UrlInfo, shortUrl string) (string, bool) {
	for _, val := range urlDetails {
		if val.ShortUrl == shortUrl {
			return val.Url, true
		}
	}
	return "", false
}
func main() {
	lambda.Start(handler)
}
