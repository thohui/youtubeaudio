package main

import (
	"os"

	"github.com/thohui/youtubeaudio/services/backend/mq"
	"github.com/thohui/youtubeaudio/services/backend/webserver"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	client, err := mq.New(os.Getenv("RABBITMQ_URI"), os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		panic(err)
	}
	go client.HandleMessages()
	webserver := webserver.New(client)
	webserver.Start()
}
