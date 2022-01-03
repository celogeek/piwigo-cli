package main

import (
	"os"
	"regexp"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CategoriesListCommand struct {
	Filter string `short:"x" long:"filter" description:"Regexp filter"`
	Empty  bool   `short:"e" long:"empty" description:"Find empty album without any photo or sub album"`
}

func (c *CategoriesListCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	categories, err := p.Categories()
	if err != nil {
		return err
	}

	var filter *regexp.Regexp
	if c.Filter != "" {
		filter = regexp.MustCompile("(?i)" + c.Filter)
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Id", "Name", "Images", "Total Images", "Url"})
	for _, category := range categories {
		if filter != nil && !filter.MatchString(category.Name) {
			continue
		}
		if c.Empty && category.TotalImagesCount != 0 {
			continue
		}
		t.AppendRow(table.Row{
			category.Id,
			category.Name,
			category.ImagesCount,
			category.TotalImagesCount,
			category.Url,
		})
	}

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
