package main

import (
	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
)

type ImagesUploadTreeCommand struct {
	Dirname    string `short:"d" long:"dirname" description:"Directory to upload" required:"true"`
	NbJobs     int    `short:"j" long:"jobs" description:"Number of jobs" default:"1"`
	CategoryId int    `short:"c" long:"category" description:"Category to upload the file" required:"true"`
}

func (c *ImagesUploadTreeCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	status, err := p.Login()
	if err != nil {
		return err
	}
	_, hasVideoJS := status.Plugins["piwigo-videojs"]

	stat := piwigotools.NewFileToUploadStat()

	defer stat.Close()
	filesToCheck := make(chan *piwigotools.FileToUpload, 10000)
	files := make(chan *piwigotools.FileToUpload, 10000)

	go p.ScanTree(c.Dirname, c.CategoryId, 0, &status.UploadFileType, stat, filesToCheck)
	go p.CheckFiles(filesToCheck, files, stat, 8)
	p.UploadFiles(files, stat, hasVideoJS, c.NbJobs, 2)

	return nil
}
