package piwigocli

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ImagesDetailsCommand struct {
	Id string `short:"i" long:"id" description:"ID of the images" required:"true"`
}

type Derivative struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type GetImagesDetailsResponse struct {
	Id            int                   `json:"id"`
	Md5           string                `json:"md5sum"`
	Name          string                `json:"name"`
	DateAvailable piwigo.TimeResult     `json:"date_available"`
	DateCreation  piwigo.TimeResult     `json:"date_creation"`
	LastModified  piwigo.TimeResult     `json:"lastmodified"`
	Width         int                   `json:"width"`
	Height        int                   `json:"height"`
	Url           string                `json:"page_url"`
	ImageUrl      string                `json:"element_url"`
	Filename      string                `json:"file"`
	Filesize      int64                 `json:"filesize"`
	Categories    piwigo.Categories     `json:"categories"`
	Tags          piwigo.Tags           `json:"tags"`
	Derivatives   map[string]Derivative `json:"derivatives"`
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

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Key", "Value"})
	t.AppendRows([]table.Row{
		{"Id", resp.Id},
		{"Md5", resp.Md5},
		{"Name", resp.Name},
		{"DateAvailable", resp.DateAvailable},
		{"DateCreation", resp.DateCreation},
		{"LastModified", resp.LastModified},
		{"Width", resp.Width},
		{"Height", resp.Height},
		{"Url", resp.Url},
		{"ImageUrl", resp.ImageUrl},
		{"Filename", resp.Filename},
		{"Filesize", resp.Filesize},
		{"Categories", strings.Join(resp.Categories.Names(), "\n")},
		{"Tags", strings.Join(resp.Tags.NamesWithAgeAt(resp.DateCreation), "\n")},
		// {"Derivatives", resp.Derivatives},
	})
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	fmt.Println("Derivatives:")
	t = table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Width", "Height", "Url"})
	for k, v := range resp.Derivatives {
		t.AppendRow(table.Row{
			k,
			v.Width,
			v.Height,
			v.Url,
		})
	}
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}
