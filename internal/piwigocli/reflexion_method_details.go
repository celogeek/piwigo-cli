package piwigocli

import (
	"net/url"
	"os"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ReflexionMethodDetailsCommand struct {
	MethodName string `short:"m" long:"method-name" description:"Method name to details"`
}

type ReflexionMethodDetailsResult struct {
	Description string `json:"description"`
}

func (c *ReflexionMethodDetailsCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var result ReflexionMethodDetailsResult

	if err := p.Post("reflection.getMethodDetails", &url.Values{
		"methodName": []string{c.MethodName},
	}, &result); err != nil {
		return err
	}

	desc := strip.StripTags(result.Description)

	t := table.NewWriter()
	t.AppendRow(table.Row{"Method", c.MethodName})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Description", desc})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Parameters"})
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}
