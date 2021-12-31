package piwigotools

import (
	"fmt"
)

type Tags []*Tag

type Tag struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	LastModified TimeResult `json:"lastmodified"`
	Birthdate    TimeResult `json:"birthdate"`
	Url          string     `json:"url"`
	UrlName      string     `json:"url_name"`
	ImageUrl     string     `json:"page_url"`
}

func (c Tags) NamesWithAgeAt(createdAt TimeResult) []string {
	names := make([]string, len(c))
	for i, category := range c {
		bd := category.Birthdate.AgeAt(createdAt)
		if bd != "" {
			names[i] = fmt.Sprintf("%s (%s)", category.Name, bd)
		} else {
			names[i] = category.Name
		}
	}
	return names
}
