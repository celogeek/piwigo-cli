package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LoginCommand struct {
	Url      string `short:"u" long:"url" description:"Url of the instance"`
	Login    string `short:"l" long:"login" description:"Login"`
	Password string `short:"p" long:"password" description:"Password"`
}

var loginCommand LoginCommand

type Result struct {
	Stat   string `json:"stat"`
	Result bool   `json:"result"`
}

func (c *LoginCommand) Execute(args []string) error {
	fmt.Printf("Login on %s...\n", c.Url)

	Url, err := url.Parse(c.Url)
	if err != nil {
		return err
	}
	Url.Path = "ws.php"
	q := Url.Query()
	q.Set("format", "json")
	q.Set("method", "pwg.session.login")
	Url.RawQuery = q.Encode()
	fmt.Println(Url.String())

	Form := url.Values{}
	Form.Set("username", c.Login)
	Form.Set("password", c.Password)

	r, err := http.PostForm(Url.String(), Form)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	result := Result{}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	if !result.Result {
		return errors.New("can't login with the credential provided")
	}

	for _, c := range r.Cookies() {
		if c.Name == "pwg_id" {
			fmt.Println("Token:", c.Value)
			break
		}
	}

	return nil
}

func init() {
	parser.AddCommand("login",
		"Initialize a connection to a piwigo instance",
		"Initialize a connection to a piwigo instance",
		&loginCommand)
}
