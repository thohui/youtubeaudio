package webserver

import (
	"encoding/json"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thohui/youtubeaudio/internal/structures"
	"github.com/thohui/youtubeaudio/services/backend/mq"
	"github.com/thohui/youtubeaudio/services/backend/youtube"
)

var (
	youtubeRegex = regexp.MustCompile("^(?:https?://)?(?:www\\.)?(?:youtu\\.be/|youtube\\.com(?:/embed/|/v/|/watch\\?v=|/watch\\?.+&v=))([\\w-]{11})(?:.+)?$")
)

type Webserver struct {
	fiber     *fiber.App
	mqClient  *mq.Client
	validator *youtube.YoutubeValidator
}

func New(mqClient *mq.Client, youtubeValidator *youtube.YoutubeValidator) *Webserver {
	server := &Webserver{
		fiber:     fiber.New(),
		mqClient:  mqClient,
		validator: youtubeValidator,
	}
	server.setupRoutes()
	return server
}

func (s *Webserver) Start() {
	err := s.fiber.Listen(":80")
	if err != nil {
		panic(err)
	}
}

type response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Location string `json:"location,omitempty"`
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
			c.Status(fiber.StatusBadRequest).JSON(response{
				Success: false,
				Message: "Invalid request body",
			})
		}
		match := youtubeRegex.FindStringSubmatch(b.URL)
		if len(match) < 2 {
			return c.Status(fiber.StatusBadRequest).JSON(response{
				Success: false,
				Message: "Invalid youtube url"},
			)
		}
		videoID := match[1]
		video, err := s.validator.ValidateURL(videoID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{
				Success: false,
				Message: "Invalid youtube url"},
			)
		}

		payload := structures.BackendPublishPayload{
			URL:   b.URL,
			Title: video.Snippet.Title,
		}
		job := make(chan []byte, 1)
		if err := s.mqClient.Publish(payload, job); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response{
				Success: false,
				Message: "Internal server error",
			})
		}
		//TODO: timeout
		msg := <-job
		r := &structures.WorkerResponse{}
		err = json.Unmarshal(msg, &r)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response{
				Success: false,
				Message: "Internal Server Error",
			})
		}
		if !r.Success {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response{
					Success: false,
					Message: "Internal Server Error",
				},
			)
		}
		return c.Status(fiber.StatusOK).JSON(response{
			Success:  true,
			Message:  "Success",
			Location: r.Location,
		})
	})
}
