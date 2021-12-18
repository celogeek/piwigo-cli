package piwigocli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ReflexionMethodDetailsCommand struct {
	MethodName string `short:"m" long:"method-name" description:"Method name to details"`
}

type ReflexionMethodDetailsParams struct {
	Name         string      `json:"name"`
	Optional     bool        `json:"optional"`
	Type         string      `json:"type"`
	AcceptArray  bool        `json:"acceptArray"`
	DefaultValue interface{} `json:"defaultValue"`
	MaxValue     interface{} `json:"maxValue"`
	Info         string      `json:"info"`
}

type ReflexionMethodDetailsOptions struct {
	Admin    bool `json:"admin_only"`
	PostOnly bool `json:"post_only"`
}

type ReflexionMethodDetailsResult struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Options     ReflexionMethodDetailsOptions  `json:"options"`
	Parameters  []ReflexionMethodDetailsParams `json:"params"`
}

func (j *ReflexionMethodDetailsOptions) UnmarshalJSON(data []byte) error {
	var r interface{}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	switch r := r.(type) {
	case map[string]interface{}:
		j.Admin, _ = r["admin_only"].(bool)
		j.PostOnly, _ = r["post_only"].(bool)
	}
	return nil
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

	t1 := table.NewWriter()
	t1.AppendRow(table.Row{"Name", result.Name})
	t1.AppendSeparator()
	t1.AppendRow(table.Row{"Description", strip.StripTags(result.Description)})
	t1.AppendRow(table.Row{"Admin", result.Options.Admin})
	t1.AppendRow(table.Row{"Post Only", result.Options.PostOnly})
	t1.SetOutputMirror(os.Stdout)
	t1.SetStyle(table.StyleLight)
	t1.Render()

	if len(result.Parameters) > 0 {
		fmt.Println("Parameters:")
		t2 := table.NewWriter()
		t2.AppendHeader(table.Row{"Name", "Type", "Optional", "Accept Array", "Default", "Max", "Info"})
		t2.AppendSeparator()
		for _, param := range result.Parameters {
			t2.AppendRow(table.Row{
				param.Name,
				param.Type,
				param.Optional,
				param.AcceptArray,
				fmt.Sprintf("%v", param.DefaultValue),
				fmt.Sprintf("%v", param.MaxValue),
				param.Info,
			})
		}
		t2.SetOutputMirror(os.Stdout)
		t2.SetStyle(table.StyleLight)
		t2.Style().Options.SeparateHeader = true
		t2.Style().Options.DrawBorder = true
		t2.Render()
	}

	return nil
}
