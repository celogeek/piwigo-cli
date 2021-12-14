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

type StatusCommand struct {
}

type SessionGroup struct {
	Login  LoginCommand  `command:"login" description:"Initialize a connection to a piwigo instance"`
	Status StatusCommand `command:"status" description:"Get the status of your session"`
}

var sessionGroup SessionGroup

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

type StatusResponse struct {
	User    string `json:"username"`
	Role    string `json:"status"`
	Version string `json:"version"`
}

func (c *StatusCommand) Execute(args []string) error {
	fmt.Println("Status:")

	Piwigo := piwigo.Piwigo{}
	if err := Piwigo.LoadConfig(); err != nil {
		return err
	}

	resp := &StatusResponse{}

	if err := Piwigo.Post("pwg.session.getStatus", &url.Values{}, &resp); err != nil {
		return err
	}
	fmt.Printf("  Version: %s\n", resp.Version)
	fmt.Printf("  User   : %s\n", resp.User)
	fmt.Printf("  Role   : %s\n", resp.Role)
	return nil
}

func init() {
	parser.AddCommand("session", "Session management", "", &sessionGroup)
}
