package main

import (
	"context"
	"encoding/json"
	"errors"
	"go-url-shortener/lib/constants"
	s3service "go-url-shortener/lib/s3-service"
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
func getMetricsData(svc service, req events.APIGatewayV2HTTPRequest) (top3Urls []Response, err error) {
	isFileExist, err := s3service.KeyExists(constants.Bucket, constants.Key, svc.s3Client)
	if err != nil {
		return nil, err
	}

	if isFileExist {
		existingInfo, err := s3service.GetS3Object(svc.s3Client, constants.Bucket, constants.Key)
		if err != nil {
			return nil, err
		}
		var urlDetails []UrlInfo
		if err := json.Unmarshal(existingInfo, &urlDetails); err != nil {
			return nil, err
		}
		if limitParam, ok := req.QueryStringParameters["limit"]; ok {
			limit, err := strconv.ParseInt(limitParam, 10, 32)
			if err != nil {
				return nil, errors.New("Specify a valid integer value")
			}
			top3Urls = getTopthreeUrls(urlDetails, int(limit))
		} else {
			top3Urls = getTopthreeUrls(urlDetails, 3)
		}

		return top3Urls, nil
	} else {
		return nil, errors.New("Metrics data does not exist")
	}
}

func getTopthreeUrls(urlDetails []UrlInfo, limit int) []Response {
	metricMap := make(map[string]int)
	for _, v := range urlDetails {
		metricMap[v.Domain] = metricMap[v.Domain] + 1
	}
	keys := make([]string, 0, len(metricMap))
	for key := range metricMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return metricMap[keys[i]] > metricMap[keys[j]] })
	count := 0
	top3Urls := []Response{}
	for _, key := range keys {
		if count >= limit {
			break
		}
		top3Urls = append(top3Urls, Response{Domain: key, Count: metricMap[key]})
		count++
	}

	return top3Urls
}
func main() {
	lambda.Start(handler)
}
