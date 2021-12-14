package piwigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (p *Piwigo) BuildUrl(method string) (string, error) {

	Url, err := url.Parse(p.Url)
	if err != nil {
		return "", err
	}
	Url.Path += "/ws.php"
	q := Url.Query()
	q.Set("format", "json")
	q.Set("method", method)
	Url.RawQuery = q.Encode()
	return Url.String(), nil
}

func (p *Piwigo) Post(method string, req *url.Values, resp interface{}) error {
	Url, err := p.BuildUrl(method)
	if err != nil {
		return err
	}

	r, err := http.PostForm(Url, *req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	Result := PiwigoResult{
		Result: resp,
	}

	err = json.Unmarshal(b, &Result)
	if err != nil {
		return err
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
