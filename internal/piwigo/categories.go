package piwigo

import (
	"errors"
	"fmt"
	"net/url"
)

type CategoriesResult struct {
	Categories `json:"categories"`
}

type Categories []Category

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ImagesCount int    `json:"nb_images"`
	Url         string `json:"url"`
}

func (p *Piwigo) Categories() (map[int]Category, error) {
	var categories CategoriesResult

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

func (p *Piwigo) CategoriesId(catId int) (map[string]int, error) {
	var categories CategoriesResult

	err := p.Post("pwg.categories.getList", &url.Values{
		"cat_id": []string{fmt.Sprint(catId)},
	}, &categories)
	if err != nil {
		return nil, err
	}

	categoriesId := make(map[string]int)
	ok := false
	for _, category := range categories.Categories {
		switch category.Id {
		case catId:
			ok = true
		default:
			categoriesId[category.Name] = category.Id
		}
	}
	if !ok {
		return nil, errors.New("category doesn't exists")
	}
	return categoriesId, nil
}
