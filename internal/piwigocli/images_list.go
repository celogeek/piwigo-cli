package piwigocli

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
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

	var results []string
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
			for _, category := range image.Categories {
				cat, ok := categories[category.Id]
				if !ok {
					continue
				}
				catName := strings.ReplaceAll(cat.Name[len(rootCatName):], " / ", "/")
				filename := fmt.Sprintf("%s/%s", catName, image.Filename)
				if filter.MatchString(filename) {
					results = append(results, filename)
				}
			}
			bar.Add(1)
		}

		if resp.Paging.Count < resp.Paging.PerPage {
			break
		}
	}
	bar.Close()

	sort.Strings(results)

	if !c.Tree {
		for _, r := range results {
			fmt.Println(r)
		}
		return nil
	}

	type Tree struct {
		Name     string
		Children []*Tree
	}

	treeMap := make(map[string]*Tree)
	treeMap[""] = &Tree{Name: "."}

	for _, r := range results {
		parentpath := ""
		fullpath := ""
		for _, s := range strings.Split(r, "/") {
			parentpath = fullpath
			fullpath += s + "/"
			if _, ok := treeMap[fullpath]; ok {
				continue
			}
			treeMap[fullpath] = &Tree{Name: s}
			treeMap[parentpath].Children = append(treeMap[parentpath].Children, treeMap[fullpath])
		}
	}

	var treeView func(*Tree, string)
	treeLinkChar := "│   "
	treeMidChar := "├── "
	treeEndChar := "└── "
	treeAfterEndChar := "    "

	treeView = func(t *Tree, prefix string) {
		for i, st := range t.Children {
			switch i {
			case len(t.Children) - 1:
				fmt.Println(prefix + treeEndChar + st.Name)
				treeView(st, prefix+treeAfterEndChar)
			case 0:
				fmt.Println(prefix + treeMidChar + st.Name)
				treeView(st, prefix+treeLinkChar)
			default:
				fmt.Println(prefix + treeMidChar + st.Name)
				treeView(st, prefix+treeLinkChar)
			}
		}
	}

	fmt.Println(treeMap[""].Name)
	treeView(treeMap[""], "")

	return nil
}
