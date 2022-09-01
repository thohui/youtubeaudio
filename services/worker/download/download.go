package download

import (
	"os/exec"

	"github.com/lithammer/shortuuid"
)

func DownloadAudio(videoURL string) (string, error) {
	id := shortuuid.New()
	name := id + ".mp3"
	path := "/tmp/" + name
	cmd := exec.Command("youtube-dl", "--prefer-ffmpeg", "-o"+path, "--extract-audio", "--audio-format", "mp3", videoURL)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return path, nil
}
