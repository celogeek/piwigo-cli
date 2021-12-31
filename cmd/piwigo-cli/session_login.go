package main

import (
	"fmt"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type SessionLoginCommand struct {
	Url      string `short:"u" long:"url" description:"Url of the instance" required:"true"`
	Login    string `short:"l" long:"login" description:"Login" required:"true"`
	Password string `short:"p" long:"password" description:"Password" required:"true"`
}

func (c *SessionLoginCommand) Execute(args []string) error {
	fmt.Printf("Login on %s...\n", c.Url)

	p := piwigo.Piwigo{
		Url:      c.Url,
		Username: c.Login,
		Password: c.Password,
	}

	_, err := p.Login()
	if err != nil {
		return err
	}

	fmt.Println("Login succeed!")

	return nil
}
