package piwigo

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
)

type CategoriesResult struct {
	Categories piwigotools.Categories `json:"categories"`
}

func (p *Piwigo) Categories() (piwigotools.Categories, error) {
	var result CategoriesResult

	if err := p.Post("pwg.categories.getList", &url.Values{
		"fullname":  []string{"true"},
		"recursive": []string{"true"},
	}, &result); err != nil {
		return nil, err
	}
	return result.Categories, nil
}

func (p *Piwigo) CategoryFromId() (map[int]piwigotools.Category, error) {
	categories, err := p.Categories()
	if err != nil {
		return nil, err
	}
	result := map[int]piwigotools.Category{}
	for _, category := range categories {
		result[category.Id] = category
	}
	return result, nil
}

func (p *Piwigo) CategoryFromName(catId int) (map[string]piwigotools.Category, error) {
	var results CategoriesResult

	err := p.Post("pwg.categories.getList", &url.Values{
		"cat_id": []string{fmt.Sprint(catId)},
	}, &results)
	if err != nil {
		return nil, err
	}

	categoriesId := map[string]piwigotools.Category{}
	ok := false
	for _, category := range results.Categories {
		switch category.Id {
		case catId:
			ok = true
		default:
			categoriesId[category.Name] = category
		}
	}
	if !ok {
		return nil, errors.New("category doesn't exists")
	}
	return categoriesId, nil
}
