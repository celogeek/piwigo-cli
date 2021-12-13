package piwigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Piwigo struct {
	Url    string
	Token  string
	Method string
}

type PiwigoResult struct {
	Stat       string      `json:"stat"`
	Err        int         `json:"err"`
	ErrMessage string      `json:"message"`
	Result     interface{} `json:"result"`
}

func (p *Piwigo) BuildUrl() (string, error) {

	Url, Err := url.Parse(p.Url)
	if Err != nil {
		return "", Err
	}
	Url.Path += "/ws.php"
	q := Url.Query()
	q.Set("format", "json")
	q.Set("method", p.Method)
	Url.RawQuery = q.Encode()
	return Url.String(), nil
}

func (p *Piwigo) Post(req *url.Values, resp interface{}) error {
	Url, Err := p.BuildUrl()
	if Err != nil {
		return Err
	}

	r, Err := http.PostForm(Url, *req)
	if Err != nil {
		return Err
	}

	defer r.Body.Close()

	b, Err := ioutil.ReadAll(r.Body)
	if Err != nil {
		return Err
	}

	Result := PiwigoResult{
		Result: resp,
	}

	Err = json.Unmarshal(b, &Result)
	if Err != nil {
		return Err
	}

	if Result.Stat != "ok" {
		return fmt.Errorf("[Error %d] %s", Result.Err, Result.ErrMessage)
	}

	for _, c := range r.Cookies() {
		if c.Name == "pwg_id" {
			p.Token = c.Value
			break
		}
	}

	return nil
}
