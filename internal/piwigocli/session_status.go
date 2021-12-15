package piwigocli

import (
	"os"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type StatusCommand struct {
}

type StatusResponse struct {
	User    string `json:"username"`
	Role    string `json:"status"`
	Version string `json:"version"`
}

func (c *StatusCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	resp := &StatusResponse{}

	if err := p.Post("pwg.session.getStatus", nil, &resp); err != nil {
		return err
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"", "Value"})
	t.AppendRows([]table.Row{
		{"Version", resp.Version},
		{"User", resp.User},
		{"Role", resp.Role},
	})

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
