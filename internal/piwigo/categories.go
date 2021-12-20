package piwigo

import (
	"net/url"
)

type Categories []Category

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ImagesCount int    `json:"nb_images"`
	Url         string `json:"url"`
}

func (p *Piwigo) Categories() (map[int]Category, error) {
	var categories struct {
		Categories `json:"categories"`
	}

	err := p.Post("pwg.categories.getList", &url.Values{
		"fullname":  []string{"true"},
		"recursive": []string{"true"},
	}, &categories)
	if err != nil {
		return nil, err
	}

	result := map[int]Category{}

	for _, category := range categories.Categories {
		result[category.Id] = category
	}
	return result, nil
}

func (c Categories) Names() []string {
	names := []string{}
	for _, category := range c {
		names = append(names, category.Name)
	}
	return names
}
