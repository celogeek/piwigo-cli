package piwigocli

import (
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type MethodTryCommand struct {
	MethodName   string     `short:"m" long:"method-name" description:"Method name to test"`
	MethodParams url.Values `short:"p" long:"params" description:"Parameter for the method" env-delim:","`
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

	if err := p.Post(c.MethodName, &c.MethodParams, &result); err != nil {
		return err
	}

	piwigo.DumpResponse(result)
	piwigo.DumpResponse(c.MethodParams)

	return nil
}
