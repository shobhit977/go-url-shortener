package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"go-url-shortener/lib/constants"
	"go-url-shortener/lib/models"
	s3service "go-url-shortener/lib/s3-service"
	"go-url-shortener/lib/service"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	svc, err := service.NewService()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body: err.Error(),
		}, nil
	}
	response, err := shortenUrl(svc, req)
	if err != nil {
		return service.ErrorResponse(err), nil
	}
	return service.UrlShortenerSuccessResponse(response), nil
}

func shortenUrl(svc service.Service, req events.APIGatewayV2HTTPRequest) (models.UrlInfo, error) {
	//return error if url is empty
	if req.Body == "" {
		return models.UrlInfo{}, errors.New("URL cannot be empty. Please provide a valid URL")
	}
	var request models.UrlShortenerRequest
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Printf("%v", err)
		return models.UrlInfo{}, err
	}
	// check if file already exists in s3 bucket
	isFileExist, err := s3service.KeyExists(constants.Bucket, constants.Key, svc.S3Client)
	if err != nil {
		return models.UrlInfo{}, err
	}
	if isFileExist {
		return getExistingFileInfo(svc, request.Url)
	}
	// if file does not exist , create a new file and add url information
	urlInfoBytes, urlInformation, err := generateUrlFileOutput(nil, request.Url)
	if err != nil {
		return models.UrlInfo{}, err
	}
	// add the file to s3 bucket
	err = s3service.PutS3Object(svc.S3Client, urlInfoBytes, constants.Bucket, constants.Key)
	if err != nil {
		return models.UrlInfo{}, err
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
func isUrlExist(svc service.Service, urlDetails []models.UrlInfo, url string) (models.UrlInfo, bool) {
	for _, val := range urlDetails {
		if val.Url == url {
			return val, true
		}
	}
	return models.UrlInfo{}, false
}

// generate short url and append to existing file
func generateUrlFileOutput(urlDetails []models.UrlInfo, url string) ([]byte, models.UrlInfo, error) {
	domain, err := getDomain(url)
	if err != nil {
		log.Printf("%v", err)
		return nil, models.UrlInfo{}, errors.New("invalid URL. Please specify a valid URL")
	}
	shortUrl := generateShortUrl(url)
	urlInformation := models.UrlInfo{
		Url:      url,
		ShortUrl: shortUrl,
		Domain:   domain,
	}
	urlDetails = append(urlDetails, urlInformation)
	allUrlInfoBytes, _ := json.Marshal(urlDetails)
	return allUrlInfoBytes, urlInformation, nil
}

// get existing url details from s3 bucket
func getExistingFileInfo(svc service.Service, url string) (models.UrlInfo, error) {
	existingInfo, err := s3service.GetS3Object(svc.S3Client, constants.Bucket, constants.Key)
	if err != nil {
		log.Printf("%v", err)
		return models.UrlInfo{}, err
	}
	var urlDetails []models.UrlInfo
	if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
		log.Printf("%v", err)
		return models.UrlInfo{}, err
	}
	// if short url already shortened then return short url instead of generating again
	shortUrl, ok := isUrlExist(svc, urlDetails, url)
	if ok {
		return shortUrl, nil
	} else {
		allUrlInfoBytes, urlInformation, err := generateUrlFileOutput(urlDetails, url)
		if err != nil {
			log.Printf("%v", err)
			return models.UrlInfo{}, err
		}
		err = s3service.PutS3Object(svc.S3Client, allUrlInfoBytes, constants.Bucket, constants.Key)
		if err != nil {
			log.Printf("%v", err)
			return models.UrlInfo{}, err
		}
		return urlInformation, nil
	}
}

func main() {
	lambda.Start(handler)
}
