package piwigocli

type CategoriesGroup struct {
	List CategoriesListCommand `command:"list" description:"List categories"`
}

var categoriesGroup CategoriesGroup

func init() {
	parser.AddCommand("categories", "Categories management", "", &categoriesGroup)
}
