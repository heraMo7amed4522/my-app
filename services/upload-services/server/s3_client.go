package server

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	session    *session.Session
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	s3Service  *s3.S3
	bucket     string
	region     string
}

func NewS3Client() (*S3Client, error) {
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if region == "" {
		region = "us-east-1"
	}
	if bucket == "" {
		return nil, fmt.Errorf("AWS_S3_BUCKET environment variable is required")
	}

	var sess *session.Session
	var err error

	if accessKey != "" && secretKey != "" {
		// Use explicit credentials
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"", // token
			),
		})
	} else {
		// Use default credential chain (IAM roles, etc.)
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &S3Client{
		session:    sess,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		s3Service:  s3.New(sess),
		bucket:     bucket,
		region:     region,
	}, nil
}

func (s3c *S3Client) UploadFile(fileData []byte, key, contentType string) (string, error) {
	input := &s3manager.UploadInput{
		Bucket:      aws.String(s3c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(contentType),
		ACL:         aws.String("private"), // Change to "public-read" if you want public access
	}

	result, err := s3c.uploader.Upload(input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return result.Location, nil
}

func (s3c *S3Client) DeleteFile(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3c.bucket),
		Key:    aws.String(key),
	}

	_, err := s3c.s3Service.DeleteObject(input)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	return nil
}

func (s3c *S3Client) GeneratePresignedURL(key string, expiresIn time.Duration) (string, error) {
	req, _ := s3c.s3Service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3c.bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(expiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return urlStr, nil
}