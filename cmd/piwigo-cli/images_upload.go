package main

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
)

type ImagesUploadCommand struct {
	Filename   string `short:"f" long:"filename" description:"File to upload" required:"true"`
	NbJobs     int    `short:"j" long:"jobs" description:"Number of jobs" default:"1"`
	CategoryId int    `short:"c" long:"category" description:"Category to upload the file"`
}

func (c *ImagesUploadCommand) Execute([]string) error {
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

	file := &piwigotools.FileToUpload{
		Dir:        filepath.Dir(c.Filename),
		Name:       filepath.Base(c.Filename),
		CategoryId: c.CategoryId,
	}

	stat := piwigotools.NewFileToUploadStat()
	defer stat.Close()
	stat.Add()
	p.Upload(file, stat, c.NbJobs, hasVideoJS)

	return nil
}
