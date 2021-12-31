package main

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type MethodListCommand struct {
	Filter string `short:"x" long:"filter" description:"Regexp filter"`
}

type MethodListResult struct {
	Methods Methods `json:"methods"`
}

type Methods []string

type MethodParams []MethodParam

type MethodParam struct {
	Name         string      `json:"name"`
	Optional     bool        `json:"optional"`
	Type         string      `json:"type"`
	AcceptArray  bool        `json:"acceptArray"`
	DefaultValue interface{} `json:"defaultValue"`
	MaxValue     interface{} `json:"maxValue"`
	Info         string      `json:"info"`
}

type MethodOptions struct {
	Admin    bool `json:"admin_only"`
	PostOnly bool `json:"post_only"`
}

func (j *MethodOptions) UnmarshalJSON(data []byte) error {
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

type MethodDetails struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Options     MethodOptions `json:"options"`
	Parameters  MethodParams  `json:"params"`
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
