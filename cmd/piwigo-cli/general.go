package main

import (
	"net/url"
	"os"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type GetInfosCommand struct {
}

type GetInfosResponse struct {
	Infos Infos `json:"infos"`
}

type Infos []Info

type Info struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

var getInfosCommand GetInfosCommand

func (c *GetInfosCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var resp GetInfosResponse

	if err := p.Post("pwg.getInfos", &url.Values{}, &resp); err != nil {
		return err
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Key", "Value"})
	for _, info := range resp.Infos {
		t.AppendRow(table.Row{info.Name, info.Value})
	}

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}

func init() {
	parser.AddCommand("getinfos", "Get general information", "", &getInfosCommand)
}
