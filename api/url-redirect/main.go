package main

import (
	"encoding/json"
	"errors"
	"go-url-shortener/lib/constants"
	s3service "go-url-shortener/lib/s3-service"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	svc, err := NewService()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	redirectUrl, err := redirect(svc, req)
	if err != nil {
		return ErrorResponse(err), nil
	}
	// redirect the short url to original url
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusTemporaryRedirect,
		Headers: map[string]string{
			"Location": redirectUrl,
		},
		IsBase64Encoded: false,
	}, nil
}

func redirect(svc service, req events.APIGatewayV2HTTPRequest) (string, error) {
	shortUrl, ok := req.PathParameters["shortUrl"]
	if !ok {
		return "", errors.New("shortURL is empty. Please provide valid shortURL")
	}
	// if file exists in s3 bucket
	isFileExist, err := s3service.KeyExists(constants.Bucket, constants.Key, svc.s3Client)
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}
	// if file exists then check if  shortUrl also exists
	if isFileExist {
		existingInfo, err := s3service.GetS3Object(svc.s3Client, constants.Bucket, constants.Key)
		if err != nil {
			return "", err
		}
		var urlDetails []UrlInfo
		if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
			log.Printf("%v", err)
			return "", err
		}
		//if short url exist then return it otherwise return error
		redirectUrl, isExist := isUrlExist(urlDetails, shortUrl)
		if !isExist {
			return "", errors.New("short URL does not exist. Please specify a valid shortUrl")
		}
		return redirectUrl, nil
	} else {
		//return error if file does not exist in s3 bucket
		return "", errors.New("short URL does not exist.Please specify a valid shortUr")
	}
}

// function to check if url has already been shortened
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
