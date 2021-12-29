package piwigotools

type Categories []Category

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ImagesCount int    `json:"nb_images"`
	Url         string `json:"url"`
}

func (c *Categories) Names() []string {
	names := []string{}
	for _, category := range *c {
		names = append(names, category.Name)
	}
	return names
}
