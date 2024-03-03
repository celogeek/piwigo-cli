package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ImagesTagCommand struct {
	Id                 int    `short:"i" long:"id" description:"image id to tag"`
	TagId              int    `short:"t" long:"tag-id" description:"look up for the first image of this tagId"`
	TagName            string `short:"T" long:"tag" description:"look up for the first image of this tagName"`
	ExcludeTags        string `short:"x" long:"exclude" description:"exclude tag from selection"`
	MaxImages          int    `short:"m" long:"max" description:"loop on a maximum number of images" default:"1"`
	KeepSurveyFilter   bool   `short:"k" long:"keep" description:"keep survey filter"`
	KeepPreviousAnswer bool   `short:"K" long:"keep-previous-answer" description:"Preserve previous answer"`
}

func (c *ImagesTagCommand) Execute(args []string) error {
	if c.MaxImages < 0 || c.MaxImages > 100 {
		return fmt.Errorf("maxImages should be between 1 and 100")
	}
	if c.Id > 0 {
		c.MaxImages = 1
	}

	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var tags struct {
		Tags piwigotools.Tags `json:"tags"`
	}
	if err := p.Post("pwg.tags.getAdminList", nil, &tags); err != nil {
		return err
	}

	sort.Slice(tags.Tags, func(i, j int) bool {
		return tags.Tags[i].Name < tags.Tags[j].Name
	})

	var exclude *regexp.Regexp
	if c.ExcludeTags != "" {
		exclude = regexp.MustCompile(c.ExcludeTags)
	}
	selectTags := tags.Tags.Selector(exclude, c.KeepSurveyFilter, c.KeepPreviousAnswer)

	imagesToTags := make([]int, 0, c.MaxImages)
	if c.Id > 0 {
		imagesToTags = append(imagesToTags, c.Id)
	} else {
		data := &url.Values{}

		switch {
		case c.TagId > 0:
			data.Set("tag_id", fmt.Sprint(c.TagId))
		case c.TagName != "":
			data.Set("tag_name", c.TagName)
		default:
			return fmt.Errorf("id or tagId or tagName are required")
		}

		data.Set("order", "date_creation")
		data.Set("per_page", fmt.Sprint(c.MaxImages))
		var results struct {
			Images []piwigotools.ImageDetails `json:"images"`
		}
		if err := p.Post("pwg.tags.getImages", data, &results); err != nil {
			return err
		}
		for _, img := range results.Images {
			imagesToTags = append(imagesToTags, img.Id)
		}
	}

	if len(imagesToTags) == 0 {
		return fmt.Errorf("no image to tag")
	}

	for _, imgId := range imagesToTags {
		for {
			var imgDetails piwigotools.ImageDetails
			if err := p.Post("pwg.images.getInfo", &url.Values{
				"image_id": []string{fmt.Sprint(imgId)},
			}, &imgDetails); err != nil {
				return err
			}

			img, err := imgDetails.Preview(25)
			if err != nil {
				return err
			}

			fmt.Println("\033[2J") // clear screen
			fmt.Println(img)

			t := table.NewWriter()
			t.AppendRows([]table.Row{
				{"Id", imgDetails.Id},
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

			fmt.Println()

			sel := selectTags()

			fmt.Println("Selection:")
			for _, name := range sel.NamesWithAgeAt(&imgDetails.DateCreation) {
				fmt.Printf("  - %s\n", name)
			}

			if len(sel) == 0 {
				exit := false
				_ = survey.AskOne(&survey.Confirm{
					Message: "Selection is empty, exit:",
					Default: false,
				}, &exit)
				if exit {
					return nil
				}
			}

			confirmSel := false
			_ = survey.AskOne(&survey.Confirm{
				Message: "Confirm:",
				Default: true,
			}, &confirmSel)

			if !confirmSel {
				continue
			}

			fmt.Println("Applying changes...")
			data := &url.Values{}
			data.Set("image_id", fmt.Sprint(imgId))
			data.Set("multiple_value_mode", "replace")
			data.Set("tag_ids", sel.JoinIds(","))

			if err := p.Post("pwg.images.setInfo", data, nil); err != nil {
				return err
			}
			fmt.Println("Done!")
			break
		}
	}
	return nil
}
