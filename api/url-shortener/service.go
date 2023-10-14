package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	Url string `json:"url"`
}
type Response struct {
	Status  int    `json:"statusCode"`
	Message string `json:"message"`
}

type UrlInfo []struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shorturl"`
	Domain   string `json:"domain"`
}
type service struct {
	sess     *session.Session
	s3Client *s3.S3
}

func NewService() (service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)
	if err != nil {
		return service{}, err
	}
	s3Client := s3.New(sess)
	return service{
		sess:     sess,
		s3Client: s3Client,
	}, nil
}
