package main

import (
	"context"
	"encoding/json"
	"errors"
	"go-url-shortener/lib/constants"
	s3service "go-url-shortener/lib/s3-service"
	"log"
	"sort"
	"strconv"

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
	response, err := getMetricsData(svc, req)
	if err != nil {
		return ErrorResponse(err), nil
	}
	return SuccessResponse(response), nil
}
func getMetricsData(svc service, req events.APIGatewayV2HTTPRequest) (mostShortenedUrls []Response, err error) {
	isFileExist, err := s3service.KeyExists(constants.Bucket, constants.Key, svc.s3Client)
	if err != nil {
		log.Printf("Failed to check file existence: %v", err)
		return nil, err
	}

	if isFileExist {
		//get file from s3 bucket
		existingInfo, err := s3service.GetS3Object(svc.s3Client, constants.Bucket, constants.Key)
		if err != nil {
			return nil, err
		}
		var urlDetails []UrlInfo
		if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
			log.Printf("Failed to unmarshal:%v", err)
			return nil, err
		}
		// check if query param is specified . If not then use default value for limit
		if limitParam, ok := req.QueryStringParameters["limit"]; ok {
			limit, err := strconv.ParseInt(limitParam, 10, 32)
			if err != nil {
				log.Printf("%v", err)
				return nil, errors.New("specify a valid integer value")
			}
			mostShortenedUrls = getMostShortenedUrls(urlDetails, int(limit))
		} else {
			mostShortenedUrls = getMostShortenedUrls(urlDetails, 3)
		}

		return mostShortenedUrls, nil
	} else {
		return nil, errors.New("metrics data does not exist")
	}
}

// function to get n most shortened url
func getMostShortenedUrls(urlDetails []UrlInfo, limit int) (mostShortenedUrls []Response) {
	metricMap := make(map[string]int)
	// create a map of domain and its frequency
	for _, v := range urlDetails {
		metricMap[v.Domain] = metricMap[v.Domain] + 1
	}
	keys := make([]string, 0, len(metricMap))
	for key := range metricMap {
		keys = append(keys, key)
	}
	// sort the map in descending order of frequency of domain
	sort.Slice(keys, func(i, j int) bool { return metricMap[keys[i]] > metricMap[keys[j]] })
	count := 0
	for _, key := range keys {
		if count >= limit {
			break
		}
		mostShortenedUrls = append(mostShortenedUrls, Response{Domain: key, Count: metricMap[key]})
		count++
	}

	return mostShortenedUrls
}
func main() {
	lambda.Start(handler)
}
