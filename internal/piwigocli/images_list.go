package piwigocli

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
	"github.com/schollz/progressbar/v3"
)

type ImagesListCommand struct {
	Recursive  bool   `short:"r" long:"recursive" description:"recursive listing"`
	CategoryId int    `short:"c" long:"category" description:"list for this category"`
	Filter     string `short:"x" long:"filter" description:"Regexp filter"`
	Tree       bool   `short:"t" long:"tree" description:"Tree view"`
}

type ImagesListResult struct {
	Images []*piwigotools.ImageDetails `json:"images"`
	Paging struct {
		Count   int `json:"count"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Total   int `json:"total_count,string"`
	} `json:"paging"`
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

	categories, err := p.CategoryFromId()
	if err != nil {
		return err
	}

	var rootCatName string
	if c.CategoryId > 0 {
		rootCat, ok := categories[c.CategoryId]
		if !ok {
			return fmt.Errorf("category doesn't exists")
		}
		rootCatName = rootCat.Name + " / "
	}

	filter := regexp.MustCompile("")
	if c.Filter != "" {
		filter = regexp.MustCompile("(?i)" + c.Filter)
	}

	rootTree := piwigotools.NewTree()

	bar := progressbar.Default(1, "listing")
	for page := 0; ; page++ {
		var resp ImagesListResult
		data := &url.Values{}
		if c.CategoryId > 0 {
			data.Set("cat_id", fmt.Sprint(c.CategoryId))
		}
		data.Set("recursive", fmt.Sprintf("%v", c.Recursive))
		data.Set("page", fmt.Sprint(page))
		if err := p.Post("pwg.categories.getImages", data, &resp); err != nil {
			return err
		}

		if page == 0 {
			bar.ChangeMax(resp.Paging.Total)
		}

		for _, image := range resp.Images {
			for _, cat := range image.Categories {
				filename := filepath.Join(
					strings.ReplaceAll(categories[cat.Id].Name[len(rootCatName):], " / ", "/"),
					image.Filename,
				)
				if !filter.MatchString(filename) {
					continue
				}
				rootTree.AddPath(filename)
			}
			bar.Add(1)
		}

		if resp.Paging.Count < resp.Paging.PerPage {
			break
		}
	}
	bar.Close()

	var results chan string
	if c.Tree {
		results = rootTree.TreeView()
	} else {
		results = rootTree.FlatView()
	}
	for filename := range results {
		fmt.Println(filename)
	}

	return nil
}
