package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"go-url-shortener/lib/constants"
	s3service "go-url-shortener/lib/s3-service"
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
	response, err := shortenUrl(svc, req)
	if err != nil {
		return ErrorResponse(err), nil
	}
	return SuccessResponse(response), nil
}

func shortenUrl(svc service, req events.APIGatewayV2HTTPRequest) (UrlInfo, error) {
	if req.Body == "" {
		return UrlInfo{}, errors.New("URL cannot be empty. Please provide a valid URL")
	}
	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		return UrlInfo{}, err
	}
	isFileExist, err := s3service.KeyExists(constants.Bucket, constants.Key, svc.s3Client)
	if err != nil {
		return UrlInfo{}, err
	}
	if isFileExist {
		return getExistingFileInfo(svc, request.Url)
	}
	// if file does not exist , create a new file and add url information
	urlInfoBytes, urlInformation, err := generateUrlFileOutput(nil, request.Url)
	if err != nil {
		return UrlInfo{}, err
	}
	err = s3service.PutS3Object(svc.s3Client, urlInfoBytes, constants.Bucket, constants.Key)
	if err != nil {
		return UrlInfo{}, err
	}
	return urlInformation, nil

}
func parseUrl(inputUrl string) (string, error) {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		log.Fatal(err)
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
func isUrlExist(svc service, urlDetails []UrlInfo, url string) (UrlInfo, bool, error) {

	for _, val := range urlDetails {
		if val.Url == url {
			return val, true, nil
		}
	}
	return UrlInfo{}, false, nil
}

func generateUrlFileOutput(urlDetails []UrlInfo, url string) ([]byte, UrlInfo, error) {
	domain, err := parseUrl(url)
	if err != nil {
		return nil, UrlInfo{}, errors.New("invalid URL")
	}
	shortUrl := generateShortUrl(url)
	urlInformation := UrlInfo{

		Url:      url,
		ShortUrl: shortUrl,
		Domain:   domain,
	}
	urlDetails = append(urlDetails, urlInformation)
	allInfoBytes, _ := json.Marshal(urlDetails)
	return allInfoBytes, urlInformation, nil
}

func getExistingFileInfo(svc service, url string) (UrlInfo, error) {
	existingInfo, err := s3service.GetS3Object(svc.s3Client, constants.Bucket, constants.Key)
	if err != nil {
		return UrlInfo{}, err
	}
	var urlDetails []UrlInfo
	if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
		return UrlInfo{}, err
	}
	shortUrl, ok, err := isUrlExist(svc, urlDetails, url)
	if err != nil {
		return UrlInfo{}, err
	}
	if ok {
		return shortUrl, nil
	} else {
		allInfoBytes, urlInformation, err := generateUrlFileOutput(urlDetails, url)
		if err != nil {
			return UrlInfo{}, err
		}
		err = s3service.PutS3Object(svc.s3Client, allInfoBytes, constants.Bucket, constants.Key)
		if err != nil {
			return UrlInfo{}, err
		}
		return urlInformation, nil
	}
}

func main() {
	lambda.Start(handler)
}
