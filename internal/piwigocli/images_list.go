package piwigocli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type ImagesListCommand struct {
	Recursive  bool `short:"r" long:"recursive" description:"recursive listing"`
	CategoryId int  `short:"c" long:"category" description:"list for this category" required:"true"`
}

type ImagesListResult struct {
	Images     []*piwigo.ImagesDetails `json:"images"`
	Pagination struct {
		Count   int `json:"count"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Total   int `json:"total_count"`
	} `json:"page"`
}

func (c *ImagesListCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var resp ImagesListResult
	data := &url.Values{}
	data.Set("cat_id", fmt.Sprint(c.CategoryId))
	data.Set("recursive", fmt.Sprintf("%v", c.Recursive))
	if err := p.Post("pwg.categories.getImages", data, &resp); err != nil {
		return err
	}

	categories, err := p.Categories()
	if err != nil {
		return err
	}

	for _, image := range resp.Images {
		for _, category := range image.Categories {
			fmt.Printf("%s/%s\n", strings.ReplaceAll(categories[category.Id].Name, " / ", "/"), image.Filename)
		}
	}

	// piwigo.DumpResponse(resp)

	return nil
}
