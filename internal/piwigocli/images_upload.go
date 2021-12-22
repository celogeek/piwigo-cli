package piwigocli

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type ImagesUploadCommand struct {
	Filename string `short:"f" long:"filename" description:"File to upload"`
	NBJobs   int    `short:"j" long:"jobs" description:"Number of jobs" default:"1"`
}

func (c *ImagesUploadCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	status, err := p.Login()
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(c.Filename)[1:])
	if _, ok := status.UploadFileType[ext]; !ok {
		return errors.New("unsupported file extension")
	}

	resp, err := p.UploadChunks(c.Filename, c.NBJobs)
	if err != nil {
		return err
	}

	switch ext {
	case "ogg", "ogv", "mp4", "m4v", "webm", "webmv":
		err = p.VideoJSSync(resp.ImageId)
		if err != nil {
			return err
		}
	}

	return nil
}
