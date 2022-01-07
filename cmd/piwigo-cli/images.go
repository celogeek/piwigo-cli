/*
Images Group Commands

This is a group of commands to manipulate the images:

	piwigo-cli images list
	piwigo-cli images details
	piwigo-cli images upload
	piwigo-cli images upload-tree
	piwigo-cli images tag

*/
package main

type ImagesGroup struct {
	List       ImagesListCommand       `command:"list" description:"List of images"`
	Details    ImageDetailsCommand     `command:"details" description:"Details of the images"`
	Upload     ImagesUploadCommand     `command:"upload" description:"Upload of an images"`
	UploadTree ImagesUploadTreeCommand `command:"upload-tree" description:"Upload of a directory of images"`
	Tag        ImagesTagCommand        `command:"tag" description:"Tag an image"`
}

var imagesGroup ImagesGroup

func init() {
	parser.AddCommand("images", "Images management", "", &imagesGroup)
}
