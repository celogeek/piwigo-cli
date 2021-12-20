package piwigocli

import (
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type ImagesDetailsCommand struct {
	Id string `short:"i" long:"id" description:"ID of the images" required:"true"`
}

type GetImagesDetailsResponse struct {
	Categories    []piwigo.Category `json:"categories"`
	DateAvailable piwigo.TimeResult `json:"date_available"`
	DateCreation  piwigo.TimeResult `json:"date_creation"`
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

	categories, err := p.Categories()
	if err != nil {
		return err
	}

	for i, category := range resp.Categories {
		resp.Categories[i] = categories[category.Id]
	}

	piwigo.DumpResponse(resp)
	return nil
}
