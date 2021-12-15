package piwigocli

import (
	"fmt"
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo"
)

type StatusCommand struct {
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
