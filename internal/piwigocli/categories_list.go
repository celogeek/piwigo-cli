package piwigocli

import (
	"net/url"
	"os"
	"regexp"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CategoriesListCommand struct {
	Filter string `short:"x" long:"filter" description:"Regexp filter"`
}

type GetCategoriesListResponse struct {
	Categories piwigo.Categories `json:"categories"`
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

	var resp GetCategoriesListResponse

	if err := p.Post("pwg.categories.getList", &url.Values{
		"recursive": []string{"true"},
		"fullname":  []string{"true"},
	}, &resp); err != nil {
		return err
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Id", "Name", "Images", "Url"})
	filter := regexp.MustCompile("")
	if c.Filter != "" {
		filter = regexp.MustCompile("(?i)" + c.Filter)
	}
	for _, category := range resp.Categories {
		if filter.MatchString(category.Name) {
			t.AppendRow(table.Row{
				category.Id,
				category.Name,
				category.ImagesCount,
				category.Url,
			})
		}
	}

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
