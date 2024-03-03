package piwigo

import (
	"errors"
	"net/url"

	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
)

type StatusResponse struct {
	User           string                     `json:"username"`
	Role           string                     `json:"status"`
	Version        string                     `json:"version"`
	Token          string                     `json:"pwg_token"`
	UploadFileType piwigotools.UploadFileType `json:"upload_file_types"`
	Plugins        piwigotools.ActivePlugin   `json:"plugins"`
}

func (p *Piwigo) GetStatus() (*StatusResponse, error) {
	if p.Url == "" || p.Username == "" || p.Password == "" {
		return nil, errors.New("missing configuration url or token")
	}

	resp := &StatusResponse{}

	err := p.Post("pwg.session.getStatus", nil, resp)
	if err != nil {
		return nil, err
	}

	if err := p.Post("pwg.plugins.getList", nil, &resp.Plugins); err != nil {
		return nil, err
	}

	if resp.User == p.Username {
		return resp, nil
	}
	return nil, errors.New("you are a guest")
}

func (p *Piwigo) Login() (*StatusResponse, error) {
	resp, err := p.GetStatus()
	if err != nil && err.Error() != "you are a guest" {
		return nil, err
	}
	if resp != nil {
		return resp, nil
	}

	err = p.Post("pwg.session.login", &url.Values{
		"username": []string{p.Username},
		"password": []string{p.Password},
	}, nil)
	if err != nil {
		return nil, err
	}

	err = p.SaveConfig()
	if err != nil {
		return nil, err
	}

	return p.GetStatus()
}
