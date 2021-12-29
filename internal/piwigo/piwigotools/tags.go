package piwigotools

import (
	"fmt"
)

type Tags []Tag

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
	names := []string{}
	for _, category := range c {
		bd := category.Birthdate.AgeAt(createdAt)
		if bd != "" {
			names = append(names, fmt.Sprintf("%s (%s)", category.Name, bd))
		} else {
			names = append(names, category.Name)
		}
	}
	return names
}
