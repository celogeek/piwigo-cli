package piwigocli

import (
	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type ImagesUploadCommand struct {
	Filename string `short:"f" long:"filename" description:"File to upload"`
}

func (c *ImagesUploadCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	err = p.UploadChunks(c.Filename)
	if err != nil {
		return err
	}

	return nil
}
