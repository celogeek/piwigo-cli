package piwigotools

type Categories []*Category

type Category struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ImagesCount      int    `json:"nb_images"`
	TotalImagesCount int    `json:"total_nb_images"`
	Url              string `json:"url"`
}

func (c *Categories) Names() []string {
	names := make([]string, len(*c))
	for i, category := range *c {
		names[i] = category.Name
	}
	return names
}
