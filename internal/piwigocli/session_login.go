package piwigocli

import (
	"fmt"
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type LoginCommand struct {
	Url      string `short:"u" long:"url" description:"Url of the instance" required:"true"`
	Login    string `short:"l" long:"login" description:"Login" required:"true"`
	Password string `short:"p" long:"password" description:"Password" required:"true"`
}

func (c *LoginCommand) Execute(args []string) error {
	fmt.Printf("Login on %s...\n", c.Url)

	p := piwigo.Piwigo{
		Url: c.Url,
	}

	result := false

	err := p.Post("pwg.session.login", &url.Values{
		"username": []string{c.Login},
		"password": []string{c.Password},
	}, &result)
	if err != nil {
		return err
	}

	err = p.SaveConfig()
	if err != nil {
		return err
	}

	fmt.Println("Login succeed!")

	return nil
}
