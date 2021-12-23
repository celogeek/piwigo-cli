package piwigocli

import (
	"os"
	"regexp"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type MethodListCommand struct {
	Filter string `short:"x" long:"filter" description:"Regexp filter"`
}

type MethodListResult struct {
	Methods piwigo.Methods `json:"methods"`
}

func (c *MethodListCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var result MethodListResult

	if err := p.Post("reflection.getMethodList", nil, &result); err != nil {
		return err
	}

	filter := regexp.MustCompile("")
	if c.Filter != "" {
		filter = regexp.MustCompile("(?i)" + c.Filter)
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Methods"})
	for _, method := range result.Methods {
		if filter.MatchString(method) {
			t.AppendRow(table.Row{
				method,
			})
		}
	}
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
