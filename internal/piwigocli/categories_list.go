package piwigocli

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CategoriesListCommand struct {
}

type Category struct {
	Id          int
	Name        string
	FullName    string
	ImagesCount int
}

func getInt(n interface{}) (r int) {
	switch n := n.(type) {
	case string:
		r, _ = strconv.Atoi(n)
	case int:
		r = n
	}
	return
}

func (b *Category) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	b.Id = getInt(v["id"])
	b.Name = v["name"].(string)
	b.FullName = v["fullname"].(string)
	b.ImagesCount = getInt(v["nb_images"])
	return nil
}

type GetCategoriesListResponse struct {
	Categories []Category `json:"categories"`
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
	if err := p.Post("pwg.categories.getAdminList", &url.Values{}, &resp); err != nil {
		return err
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Id", "FullName", "ImagesCount"})
	for _, category := range resp.Categories {
		t.AppendRow(table.Row{
			category.Id,
			category.FullName,
			category.ImagesCount,
		})
	}

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
