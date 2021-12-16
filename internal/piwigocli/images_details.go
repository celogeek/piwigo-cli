package piwigocli

import (
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type ImagesDetailsCommand struct {
	Id string `short:"i" long:"id" description:"ID of the images" required:"true"`
}

type GetImagesDetailsResponse struct {
}

func (c *ImagesDetailsCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var resp GetImagesDetailsResponse
	if err := p.Post("pwg.images.getInfo", &url.Values{
		"image_id": []string{c.Id},
	}, &resp); err != nil {
		return err
	}
	return nil
}
