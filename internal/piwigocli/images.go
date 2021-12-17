package piwigocli

type ImagesGroup struct {
	Details ImagesDetailsCommand `command:"details" description:"Details of the images"`
}

var imagesGroup ImagesGroup

func init() {
	parser.AddCommand("images", "Images management", "", &imagesGroup)
}
