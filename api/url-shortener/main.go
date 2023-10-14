package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	svc, err := NewService()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}
	}
	log.Print(svc)
	return events.APIGatewayV2HTTPResponse{}
}
func main() {
	lambda.Start(handler)
}

// function to generate shortened url via sha1 algorithm
func generateShortUrl(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	sha1_hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	//retrun first 10 characters
	return sha1_hash[:10]
}
