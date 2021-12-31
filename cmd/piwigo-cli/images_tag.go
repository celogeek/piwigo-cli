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

	var resp piwigotools.ImageDetails
	if err := p.Post("pwg.images.getInfo", &url.Values{
		"image_id": []string{fmt.Sprint(c.Id)},
	}, &resp); err != nil {
		return err
	}

	img, err := resp.Preview(25)
	if err != nil {
		return err
	}

	fmt.Println(img)
	t := table.NewWriter()
	t.AppendRows([]table.Row{
		{"Name", resp.Name},
		{"Url", resp.Url},
		{"CreatedAt", resp.DateCreation},
		{"Size", fmt.Sprintf("%d x %d", resp.Width, resp.Height)},
		{"Categories", strings.Join(resp.Categories.Names(), "\n")},
		{"Tags", strings.Join(resp.Tags.NamesWithAgeAt(resp.DateCreation), "\n")},
	})

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
