package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	svc, err := NewService()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body: err.Error(),
		}, nil
	}
	return shortenUrl(svc, req)
}

func shortenUrl(svc service, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if req.Body == "" {
		return events.APIGatewayV2HTTPResponse{
			Body: errors.New("Empty Request Body").Error(),
		}, nil
	}
	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}
	domain, err := parseUrl(request.Url)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body: "Invalid URL",
		}, nil
	}
	log.Print(domain)
	return events.APIGatewayV2HTTPResponse{}, nil

}

func parseUrl(inputUrl string) (string, error) {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}
	replacer := strings.NewReplacer("www.", "", ".com", "")
	domain := replacer.Replace(parsedUrl.Hostname())
	return domain, nil
}

// function to generate shortened url via sha1 algorithm
func generateShortUrl(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	sha1_hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	//retrun first 10 characters
	return sha1_hash[:10]
}

func main() {
	lambda.Start(handler)
}
