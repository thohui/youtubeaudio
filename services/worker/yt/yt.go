package yt

import (
	"os"
	"os/exec"

	"github.com/lithammer/shortuuid"
)

type File struct {
	Name string
	Path string
}

func (f File) Delete() error {
	return os.Remove(f.Path)
}

func DownloadAudio(videoURL string) (File, error) {
	id := shortuuid.New()
	name := id + ".mp3"
	path := "/tmp/" + name
	cmd := exec.Command("youtube-dl", "--prefer-ffmpeg", "-o"+path, "--extract-audio", "--audio-format", "mp3", videoURL)
	err := cmd.Run()
	return File{name, path}, err
}
