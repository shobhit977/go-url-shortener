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
	//return error if url is empty
	if req.Body == "" {
		return UrlInfo{}, errors.New("URL cannot be empty. Please provide a valid URL")
	}
	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Printf("%v", err)
		return UrlInfo{}, err
	}
	// check if file already exists in s3 bucket
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
	// add the file to s3 bucket
	err = s3service.PutS3Object(svc.s3Client, urlInfoBytes, constants.Bucket, constants.Key)
	if err != nil {
		return UrlInfo{}, err
	}
	return urlInformation, nil

}

/*
function to get the domain of the URL

Example - URL : www.google.com
Domain : google
*/
func getDomain(inputUrl string) (string, error) {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		log.Print(err)
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

// check if the is url already shortened
func isUrlExist(svc service, urlDetails []UrlInfo, url string) (UrlInfo, bool) {
	for _, val := range urlDetails {
		if val.Url == url {
			return val, true
		}
	}
	return UrlInfo{}, false
}

// generate short url and append to existing file
func generateUrlFileOutput(urlDetails []UrlInfo, url string) ([]byte, UrlInfo, error) {
	domain, err := getDomain(url)
	if err != nil {
		log.Printf("%v", err)
		return nil, UrlInfo{}, errors.New("invalid URL. Please specify a valid URL")
	}
	shortUrl := generateShortUrl(url)
	urlInformation := UrlInfo{
		Url:      url,
		ShortUrl: shortUrl,
		Domain:   domain,
	}
	urlDetails = append(urlDetails, urlInformation)
	allUrlInfoBytes, _ := json.Marshal(urlDetails)
	return allUrlInfoBytes, urlInformation, nil
}

// get existing url details from s3 bucket
func getExistingFileInfo(svc service, url string) (UrlInfo, error) {
	existingInfo, err := s3service.GetS3Object(svc.s3Client, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return UrlInfo{}, err
	}
	var urlDetails []UrlInfo
	if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
		log.Printf("%v", err)
		return UrlInfo{}, err
	}
	// if short url already shortened then return short url instead of generating again
	shortUrl, ok := isUrlExist(svc, urlDetails, url)
	if ok {
		return shortUrl, nil
	} else {
		allUrlInfoBytes, urlInformation, err := generateUrlFileOutput(urlDetails, url)
		if err != nil {
			log.Printf("%v", err)
			return UrlInfo{}, err
		}
		err = s3service.PutS3Object(svc.s3Client, allUrlInfoBytes, constants.Bucket, constants.Key)
		if err != nil {
			log.Printf("%v", err)
			return UrlInfo{}, err
		}
		return urlInformation, nil
	}
}

func main() {
	lambda.Start(handler)
}
