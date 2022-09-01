package handler

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/thohui/youtubeaudio/internal/structures"
	"github.com/thohui/youtubeaudio/services/worker/download"
	"github.com/thohui/youtubeaudio/services/worker/mq"
	"github.com/thohui/youtubeaudio/services/worker/s3"
)

type Handler struct {
	client   *mq.Client
	s3client *s3.Client
}

func New(client *mq.Client, s3Client *s3.Client) *Handler {
	return &Handler{
		client:   client,
		s3client: s3Client,
	}
}

func (h *Handler) Start() {
	messages, err := h.client.Consume()
	if err != nil {
		panic(err)
	}
	fmt.Println("Waiting for jobs...")
	for msg := range messages {
		fmt.Println("Got job", string(msg.Body))
		go h.handle(msg)
	}
}

func (h *Handler) handle(msg amqp.Delivery) {
	var pl structures.BackendPublishPayload
	err := json.Unmarshal(msg.Body, &pl)
	if err != nil {
		return
	}
	path, err := download.DownloadAudio(pl.URL)
	response := &structures.WorkerResponse{}
	if err == nil {
		location, err := h.s3client.Upload(pl.Title+".mp3", path)
		if err == nil {
			response.Success = true
			response.Location = location
		}
	}
	data, _ := json.Marshal(response)
	h.client.Channel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Body:          data,
		ReplyTo:       msg.ReplyTo,
		CorrelationId: msg.CorrelationId,
	})
}
