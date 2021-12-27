package piwigocli

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/schollz/progressbar/v3"
)

type ImagesUploadCommand struct {
	Filename   string `short:"f" long:"filename" description:"File to upload" required:"true"`
	NbJobs     int    `short:"j" long:"jobs" description:"Number of jobs" default:"1"`
	CategoryId int    `short:"c" long:"category" description:"Category to upload the file"`
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
	if !status.UploadFileType.Has(ext) {
		return errors.New("unsupported file extension")
	}

	_, hasVideoJS := status.Plugins["piwigo-videojs"]

	file := &piwigo.FileToUpload{
		Dir:        filepath.Dir(c.Filename),
		Name:       filepath.Base(c.Filename),
		CategoryId: c.CategoryId,
	}

	stat := &piwigo.FileToUploadStat{
		Progress: progressbar.DefaultBytes(1, "..."),
	}
	defer stat.Close()
	stat.Add()
	err = p.Upload(file, stat, c.NbJobs, hasVideoJS)
	if err != nil {
		return err
	}

	return nil
}
