package main

import (
	"os"

	"github.com/thohui/youtubeaudio/services/backend/mq"
	"github.com/thohui/youtubeaudio/services/backend/webserver"
	"github.com/thohui/youtubeaudio/services/backend/youtube"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	client, err := mq.New(os.Getenv("RABBITMQ_URI"), os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		panic(err)
	}
	validator, err := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	if err != nil {
		panic(err)
	}
	webserver := webserver.New(client, validator)
	go client.HandleMessages()
	webserver.Start()
}
