package piwigotools

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type Tags []*Tag

type Tag struct {
	Id           int        `json:"id,string"`
	Name         string     `json:"name"`
	LastModified TimeResult `json:"lastmodified"`
	Birthdate    TimeResult `json:"birthdate"`
	Url          string     `json:"url"`
	UrlName      string     `json:"url_name"`
	ImageUrl     string     `json:"page_url"`
}

func (t Tags) NamesWithAgeAt(createdAt *TimeResult) []string {
	names := make([]string, len(t))
	for i, tag := range t {
		names[i] = tag.NameWithAgeAt(createdAt)
	}
	return names
}

func (t *Tag) NameWithAgeAt(createdAt *TimeResult) string {
	bd := t.Birthdate.AgeAt(createdAt)
	if bd != "" {
		return fmt.Sprintf("%s (%s)", t.Name, bd)
	}
	return t.Name
}

func (t Tags) Ids() []int {
	ids := make([]int, len(t))
	for i, tag := range t {
		ids[i] = tag.Id
	}
	return ids
}

func (t Tags) JoinIds(sep string) string {
	ids := make([]string, len(t))
	for i, tag := range t {
		ids[i] = fmt.Sprint(tag.Id)
	}
	return strings.Join(ids, sep)
}

func (t Tags) Selector(exclude *regexp.Regexp, keepFilter bool) func() Tags {
	options := make([]string, 1, len(t))
	options[0] = "Same as before"
	tags := map[string]*Tag{}
	tags[options[0]] = nil
	for _, tag := range t {
		if exclude != nil && exclude.MatchString(tag.Name) {
			continue
		}
		options = append(options, tag.Name)
		tags[tag.Name] = tag
	}

	var previousTags Tags
	return func() Tags {
		answer := []string{}

		survey.AskOne(&survey.MultiSelect{
			Message:  "Tags:",
			Options:  options,
			PageSize: 20,
		}, &answer, survey.WithKeepFilter(keepFilter))

		result := make([]*Tag, 0, len(answer))
		alreadySelected := map[int]bool{}
		for _, a := range answer {
			if tags[a] == nil {
				for _, p := range previousTags {
					if _, ok := alreadySelected[p.Id]; !ok {
						result = append(result, p)
					}
					alreadySelected[p.Id] = true
				}
			} else {
				if _, ok := alreadySelected[tags[a].Id]; !ok {
					result = append(result, tags[a])
				}
				alreadySelected[tags[a].Id] = true
			}
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})

		previousTags = result
		return result
	}

}
