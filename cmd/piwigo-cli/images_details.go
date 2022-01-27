/*
Images Details Command

Get details of an image. It supports the birthday plugin.

    $ piwigo-cli images details -i 38062

    ┌───────────────┬───────────────────────────────────────────────────────────────────────────┐
    │ KEY           │ VALUE                                                                     │
    ├───────────────┼───────────────────────────────────────────────────────────────────────────┤
    │ Id            │ 38062                                                                     │
    │ Md5           │ 6ad2abade6d5460181890e2bad671002                                          │
    │ Name          │ 2006 04 14 015                                                            │
    │ DateAvailable │ 2021-11-25 20:25:05                                                       │
    │ DateCreation  │ 2006-04-14 04:14:00                                                       │
    │ LastModified  │ 2022-01-01 23:11:48                                                       │
    │ Width         │ 1984                                                                      │
    │ Height        │ 1488                                                                      │
    │ Url           │ https://yourphotos/picture.php?/38062                                     │
    │ ImageUrl      │ https://yourphotos/upload/2021/11/25/20211125202505-6ad2abad.jpg          │
    │ Filename      │ 2006_04_14_015.jpeg                                                       │
    │ Filesize      │ 513                                                                       │
    │ Categories    │ 2007                                                                      │
    │ Tags          │ User Tag 1 (46 years old)                                                 │
    │               │ User Tag 2 (8 months old)                                                 │
    │               │ User Tag 3 (48 years old)                                                 │
    └───────────────┴───────────────────────────────────────────────────────────────────────────┘
    Derivatives:
    ┌─────────┬───────┬────────┬──────────────────────────────────────────────────────────────────────────────────────┐
    │ NAME    │ WIDTH │ HEIGHT │ URL                                                                                  │
    ├─────────┼───────┼────────┼──────────────────────────────────────────────────────────────────────────────────────┤
    │ thumb   │   144 │    108 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-th.jpg           │
    │ xsmall  │   432 │    324 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-xs.jpg           │
    │ xxlarge │  1656 │   1242 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-xx.jpg          │
    │ square  │   120 │    120 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-sq.jpg          │
    │ small   │   576 │    432 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-sm.jpg          │
    │ medium  │   792 │    594 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-me.jpg          │
    │ large   │  1008 │    756 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-la.jpg           │
    │ xlarge  │  1224 │    918 │ https://yourphotos/_data/i/upload/2021/11/25/20211125202505-6ad2abad-xl.jpg          │
    │ 2small  │   240 │    180 │ https://yourphotos/i.php?/upload/2021/11/25/20211125202505-6ad2abad-2s.jpg           │
    └─────────┴───────┴────────┴──────────────────────────────────────────────────────────────────────────────────────┘

*/
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

func (c *ImageDetailsCommand) Execute(args []string) error {
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
