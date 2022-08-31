package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/thohui/youtubeaudio/services/worker/handler"
	"github.com/thohui/youtubeaudio/services/worker/mq"
	"github.com/thohui/youtubeaudio/services/worker/s3"
)

func main() {
	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("S3_REGION")
	keyId := os.Getenv("S3_KEY_ID")
	applicationKey := os.Getenv("S3_ACCESS_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")

	s3client, err := s3.New(endpoint, region, keyId, applicationKey, "", bucketName)
	if err != nil {
		panic(err)
	}
	client, err := mq.New(os.Getenv("RABBITMQ_URI"), os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		panic(err)
	}
	worker := handler.New(client, s3client)
	worker.Start()
}
