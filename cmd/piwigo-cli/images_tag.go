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

type ImagesTagCommand struct {
	Id int `short:"i" long:"id" description:"image id to tag" required:"true"`
}

func (c *ImagesTagCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var imgDetails piwigotools.ImageDetails
	if err := p.Post("pwg.images.getInfo", &url.Values{
		"image_id": []string{fmt.Sprint(c.Id)},
	}, &imgDetails); err != nil {
		return err
	}

	var tags struct {
		Tags piwigotools.Tags `json:"tags"`
	}
	if err := p.Post("pwg.tags.getAdminList", &url.Values{
		"image_id": []string{fmt.Sprint(c.Id)},
	}, &tags); err != nil {
		return err
	}

	img, err := imgDetails.Preview(25)
	if err != nil {
		return err
	}
	fmt.Println(img)

	t := table.NewWriter()
	t.AppendRows([]table.Row{
		{"Name", imgDetails.Name},
		{"Url", imgDetails.Url},
		{"CreatedAt", imgDetails.DateCreation},
		{"Size", fmt.Sprintf("%d x %d", imgDetails.Width, imgDetails.Height)},
		{"Categories", strings.Join(imgDetails.Categories.Names(), "\n")},
		{"Tags", strings.Join(imgDetails.Tags.NamesWithAgeAt(&imgDetails.DateCreation), "\n")},
	})

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()
	return nil
}
