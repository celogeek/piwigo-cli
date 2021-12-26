package piwigocli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type ImagesUploadTreeCommand struct {
	Dirname    string `short:"d" long:"dirname" description:"Directory to upload" required:"true"`
	NbJobs     int    `short:"j" long:"jobs" description:"Number of jobs" default:"1"`
	CategoryId int    `short:"c" long:"category" description:"Category to upload the file" required:"true"`
}

func (c *ImagesUploadTreeCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	rootPath, err := filepath.Abs(c.Dirname)
	if err != nil {
		return err
	}

	categoriesId, err := p.CategoriesId(c.CategoryId)
	if err != nil {
		return err
	}

	dirs, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		if _, ok := categoriesId[dir.Name()]; !ok {
			fmt.Println("Creating", dir.Name(), "...")
		}
	}

	return nil
}
