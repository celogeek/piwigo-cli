package piwigocli

import (
	"os"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ReflexionMethodListCommand struct{}

type ReflexionMethodListResult struct {
	Methods []string `json:"methods"`
}

func (c *ReflexionMethodListCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var result ReflexionMethodListResult

	if err := p.Post("reflection.getMethodList", nil, &result); err != nil {
		return err
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Methods"})
	for _, method := range result.Methods {
		t.AppendRow(table.Row{
			method,
		})
	}
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
