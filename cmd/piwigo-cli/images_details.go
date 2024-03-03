package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ImageDetailsCommand struct {
	Id int `short:"i" long:"id" description:"ID of the images" required:"true"`
}

func (c *ImageDetailsCommand) Execute([]string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var resp piwigotools.ImageDetails
	if err := p.Post("pwg.images.getInfo", &url.Values{
		"image_id": []string{fmt.Sprint(c.Id)},
	}, &resp); err != nil {
		return err
	}

	categories, err := p.CategoryFromId()
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
		{"Tags", strings.Join(resp.Tags.NamesWithAgeAt(&resp.DateCreation), "\n")},
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
