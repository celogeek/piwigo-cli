package main

import (
	"os"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type SessionStatusCommand struct {
}

func (c *SessionStatusCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	resp, err := p.Login()
	if err != nil {
		return err
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"", "Value"})
	t.AppendRows([]table.Row{
		{"Version", resp.Version},
		{"User", resp.User},
		{"Role", resp.Role},
		{"Admin Token", resp.Token},
		{"Supported formats", resp.UploadFileType},
		{"Supported plugins", resp.Plugins},
	})

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
