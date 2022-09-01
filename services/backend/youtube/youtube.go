package youtube

import (
	"context"
	"errors"

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

func (v *YoutubeValidator) ValidateURL(videoID string) (*youtube.Video, error) {
	res, err := v.youtubeService.Videos.List([]string{"id", "snippet"}).Id(videoID).Do()
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, errors.New("could not find video with id " + videoID)
	}
	return res.Items[0], nil
}
