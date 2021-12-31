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
