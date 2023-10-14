package s3service

import (
	"bytes"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func PutS3Object(s3Client *s3.S3, byteData []byte, bucket string, key string) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(bytes.NewReader(byteData)),
	}
	_, err := s3Client.PutObject(params)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func GetS3Object(s3Client *s3.S3, bucket string, key string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := s3Client.GetObject(params)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	defer result.Body.Close()

	// capture all bytes from upload
	byteData, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return byteData, nil

}

func KeyExists(bucket string, key string, s3Client *s3.S3) (bool, error) {
	_, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}
