package youtube

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeValidator struct {
	youtubeService *youtube.Service
}

func New(youtubeApiKey string) (*YoutubeValidator, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(youtubeApiKey))
	if err != nil {
		return &YoutubeValidator{}, err
	}
	return &YoutubeValidator{youtubeService: service}, nil
}

func (v *YoutubeValidator) ValidateURL(videoURL string) bool {
	res, err := v.youtubeService.Videos.List([]string{"id"}).Id(videoURL).Do()
	if err != nil {
		return false
	}
	return len(res.Items) > 0
}
