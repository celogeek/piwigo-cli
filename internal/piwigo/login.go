package piwigo

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type UploadFileType map[string]bool

type StatusResponse struct {
	User           string         `json:"username"`
	Role           string         `json:"status"`
	Version        string         `json:"version"`
	Token          string         `json:"pwg_token"`
	UploadFileType UploadFileType `json:"upload_file_types"`
}

func (uft *UploadFileType) UnmarshalJSON(data []byte) error {
	var r string
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*uft = UploadFileType{}
	for _, v := range strings.Split(r, ",") {
		(*uft)[v] = true
	}
	return nil
}

func (uft UploadFileType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + uft.String() + `"`), nil
}

func (uft UploadFileType) String() string {
	keys := make([]string, 0, len(uft))
	for k, _ := range uft {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
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
