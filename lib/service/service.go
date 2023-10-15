package service

import (
	"encoding/json"
	"go-url-shortener/lib/constants"
	errorlib "go-url-shortener/lib/errorLib"
	"go-url-shortener/lib/models"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Service struct {
	Sess     *session.Session
	S3Client *s3.S3
}

func NewService() (Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.Region)},
	)
	if err != nil {
		return Service{}, err
	}
	s3Client := s3.New(sess)
	return Service{
		Sess:     sess,
		S3Client: s3Client,
	}, nil
}
func UrlShortenerSuccessResponse(resp models.UrlInfo) events.APIGatewayV2HTTPResponse {
	respBytes, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusOK,
	}
}

func MetricsSuccessResponse(resp []models.MetricsResponse) events.APIGatewayV2HTTPResponse {
	respBytes, _ := json.Marshal(resp)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusOK,
	}
}

func ErrorResponse(err error) events.APIGatewayV2HTTPResponse {
	errRes := errorlib.ErrorLib{
		Message: err.Error(),
		Code:    http.StatusBadRequest,
	}
	respBytes, _ := json.Marshal(errRes)
	return events.APIGatewayV2HTTPResponse{
		Body:       string(respBytes),
		StatusCode: http.StatusBadRequest,
	}
}
