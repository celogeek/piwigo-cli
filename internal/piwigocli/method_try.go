package piwigocli

import (
	"errors"
	"net/url"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
)

type MethodTryCommand struct {
	MethodName string `short:"m" long:"method-name" description:"Method name to test"`
}

func (c *MethodTryCommand) Execute(args []string) error {
	p := piwigo.Piwigo{}
	if err := p.LoadConfig(); err != nil {
		return err
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	var result interface{}
	params, err := ArgsToForm(args)
	if err != nil {
		return err
	}

	if err := p.Post(c.MethodName, params, &result); err != nil {
		piwigotools.DumpResponse(params)
		return err
	}

	piwigotools.DumpResponse(map[string]interface{}{
		"params": params,
		"result": result,
	})
	return nil
}

func ArgsToForm(args []string) (*url.Values, error) {
	params := &url.Values{}
	for _, arg := range args {
		r := strings.SplitN(arg, "=", 2)
		if len(r) != 2 {
			return nil, errors.New("args should be key=value")
		}
		params.Add(r[0], strings.ReplaceAll(r[1], "'", `\'`))
	}
	return params, nil
}
