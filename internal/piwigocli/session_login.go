package piwigocli

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

func (c *LoginCommand) Execute(args []string) error {
	fmt.Printf("Login on %s...\n", c.Url)

	Piwigo := piwigo.Piwigo{
		Url: c.Url,
	}

	result := false

	err := Piwigo.Post("pwg.session.login", &url.Values{
		"username": []string{c.Login},
		"password": []string{c.Password},
	}, &result)
	if err != nil {
		return err
	}

	err = Piwigo.SaveConfig()
	if err != nil {
		return err
	}

	fmt.Println("Login succeed!")

	return nil
}
