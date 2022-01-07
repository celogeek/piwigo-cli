/*
Images List Command

List the images of a category.

Recursive list:

	$ piwigo-cli images list -r

	Category1/SubCategory1/IMG_00001.jpeg
	Category1/SubCategory1/IMG_00002.jpeg
	Category1/SubCategory1/IMG_00003.jpeg
	Category1/SubCategory1/IMG_00004.jpeg
	Category1/SubCategory2/IMG_00005.jpeg
	Category1/SubCategory2/IMG_00006.jpeg
	Category2/SubCategory1/IMG_00007.jpeg

Specify a category:

	$ piwigo-cli images list -r -c 2

	Category2/SubCategory1/IMG_00007.jpeg

Tree view:
	$ piwigo-cli images list -r -c 1 -t

	.
	├── SubCategory1
	│   ├── IMG_00001.jpeg
	│   ├── IMG_00002.jpeg
	│   ├── IMG_00003.jpeg
	│   └── IMG_00004.jpeg
	└── SubCategory2
		├── IMG_00005.jpeg
		└── IMG_00006.jpeg

*/
package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
	"github.com/celogeek/piwigo-cli/internal/tree"
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
		rootCatName = rootCat.Name
	}

	filter := regexp.MustCompile("")
	if c.Filter != "" {
		filter = regexp.MustCompile("(?i)" + c.Filter)
	}

	rootTree := tree.New()

	bar := progressbar.Default(1, "listing")
	progressbar.OptionOnCompletion(func() {
		os.Stderr.WriteString("\n")
	})(bar)
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
				filename, _ := filepath.Rel(rootCatName,
					filepath.Join(
						categories[cat.Id].Name,
						image.Filename,
					),
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
