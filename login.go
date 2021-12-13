package main

import (
	"fmt"
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type LoginCommand struct {
	Url      string `short:"u" long:"url" description:"Url of the instance"`
	Login    string `short:"l" long:"login" description:"Login"`
	Password string `short:"p" long:"password" description:"Password"`
}

var loginCommand LoginCommand

func (c *LoginCommand) Execute(args []string) error {
	fmt.Printf("Login on %s...\n", c.Url)

	Piwigo := piwigo.Piwigo{
		Url:    c.Url,
		Method: "pwg.session.login",
	}

	result := false

	if Err := Piwigo.Post(&url.Values{
		"username": []string{c.Login},
		"password": []string{c.Password},
	}, &result); Err != nil {
		return Err
	}

	fmt.Printf("Token: %s\n", Piwigo.Token)

	return nil
}

func init() {
	parser.AddCommand("login",
		"Initialize a connection to a piwigo instance",
		"Initialize a connection to a piwigo instance",
		&loginCommand)
}
