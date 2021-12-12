package main

import "fmt"

type TagsCommand struct {
	All bool `short:"a" long:"all" description:"Get all available tags"`
}

var tagsCommand TagsCommand

func (c *TagsCommand) Execute(args []string) error {
	fmt.Printf("List tags %v\n", c.All)

	return nil
}

func init() {
	parser.AddCommand("tags",
		"List tags",
		"List tags",
		&tagsCommand)
}
