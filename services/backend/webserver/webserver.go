package webserver

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thohui/youtubeaudio/services/backend/mq"
)

var (
	youtubeRegex = regexp.MustCompile("^(?:https?://)?(?:www\\.)?(?:youtu\\.be/|youtube\\.com(?:/embed/|/v/|/watch\\?v=|/watch\\?.+&v=))([\\w-]{11})(?:.+)?$")
)

type Webserver struct {
	fiber    *fiber.App
	mqClient *mq.Client
}

func New(mqClient *mq.Client) *Webserver {
	server := &Webserver{
		fiber:    fiber.New(),
		mqClient: mqClient,
	}
	server.setupRoutes()
	return server
}

func (s *Webserver) Start() {
	s.fiber.Listen(":80")
}

func (s *Webserver) setupRoutes() {
	s.fiber.Use(cors.New())
	s.fiber.Post("/convert", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		type body struct {
			URL string `json:"youtube_url"`
		}
		var b body
		if err := c.BodyParser(&b); err != nil {
			c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		match := youtubeRegex.MatchString(b.URL)
		if !match {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid youtube URL")
		}
		job := make(chan []byte)
		if err := s.mqClient.Publish(b.URL, job); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to publish message")
		}
		//TODO: timeout
		msg := <-job
		return c.Status(fiber.StatusOK).Send(msg)
	})
}
