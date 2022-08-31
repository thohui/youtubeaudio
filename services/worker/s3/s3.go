package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client struct {
	bucket string
	s3     *s3.S3
}

func New(endpoint, region, keyId, applicationKey, token, bucket string) (*Client, error) {
	conf := &aws.Config{
		Endpoint:    &endpoint,
		Region:      &region,
		Credentials: credentials.NewStaticCredentials(keyId, applicationKey, token),
	}

	awsSession, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(awsSession)

	return &Client{
		s3:     s3Client,
		bucket: bucket,
	}, nil
}

func (b *Client) Upload(fileName string, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploader := s3manager.NewUploaderWithClient(b.s3)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &b.bucket,
		Key:    &fileName,
		Body:   file,
	})
	return result.Location, err
}
