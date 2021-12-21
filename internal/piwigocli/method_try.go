package piwigocli

import (
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

	var result interface{}
	params, err := piwigo.ArgsToForm(args)
	if err != nil {
		return err
	}

	if err := p.Post(c.MethodName, params, &result); err != nil {
		piwigo.DumpResponse(params)
		return err
	}

	piwigo.DumpResponse(map[string]interface{}{
		"params": params,
		"result": result,
	})
	return nil
}
