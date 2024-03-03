package main

type CategoriesGroup struct {
	List CategoriesListCommand `command:"list" description:"List categories"`
}

var categoriesGroup CategoriesGroup

func init() {
	_, err := parser.AddCommand("categories", "Categories management", "", &categoriesGroup)
	if err != nil {
		panic(err)
	}
}
