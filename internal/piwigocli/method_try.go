package piwigocli

import (
	"errors"
	"net/url"
	"strings"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
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

	var result map[string]interface{}
	params := &url.Values{}
	for _, arg := range args {
		r := strings.SplitN(arg, "=", 2)
		if len(r) != 2 {
			return errors.New("args should be key=value")
		}
		params.Add(r[0], r[1])
	}

	if err := p.Post(c.MethodName, params, &result); err != nil {
		piwigo.DumpResponse(params)
		return err
	}

	piwigo.DumpResponse(result)
	piwigo.DumpResponse(params)
	return nil
}
