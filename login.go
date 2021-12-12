package main

import "fmt"

type LoginCommand struct {
	Url      string `short:"u" long:"url" description:"Url of the instance"`
	Login    string `short:"l" long:"login" description:"Login"`
	Password string `short:"p" long:"password" description:"Password"`
}

var loginCommand LoginCommand

func (c *LoginCommand) Execute(args []string) error {
	fmt.Printf("Login on %s...\n", c.Url)

	return nil
}

func init() {
	parser.AddCommand("login",
		"Initialize a connection to a piwigo instance",
		"Initialize a connection to a piwigo instance",
		&loginCommand)
}
