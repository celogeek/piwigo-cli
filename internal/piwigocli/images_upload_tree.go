package piwigocli

import (
	"fmt"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
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

	files, err := p.UploadTree(c.Dirname, c.CategoryId, 0, status.UploadFileType)
	if err != nil {
		return err
	}
	fmt.Println("Total", len(files))

	return nil
}
