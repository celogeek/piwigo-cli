package piwigocli

type ImagesGroup struct {
	List       ImagesListCommand       `command:"list" description:"List of images"`
	Details    ImageDetailsCommand     `command:"details" description:"Details of the images"`
	Upload     ImagesUploadCommand     `command:"upload" description:"Upload of an images"`
	UploadTree ImagesUploadTreeCommand `command:"upload-tree" description:"Upload of a directory of images"`
}

var imagesGroup ImagesGroup

func init() {
	parser.AddCommand("images", "Images management", "", &imagesGroup)
}
