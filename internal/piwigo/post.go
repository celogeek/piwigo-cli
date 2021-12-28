package piwigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
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

func (p *Piwigo) Post(method string, form *url.Values, resp interface{}) error {
	Url, err := p.BuildUrl(method)
	if err != nil {
		return err
	}

	var encodedForm string
	if form != nil {
		encodedForm = form.Encode()
	}

	Result := PiwigoResult{
		Result: resp,
	}

	var raw []byte

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", Url, strings.NewReader(encodedForm))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if p.Token != "" {
			req.AddCookie(&http.Cookie{Name: "pwg_id", Value: p.Token, HttpOnly: true})
		}

		r, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}

		raw, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			continue
		}

		err = json.Unmarshal(raw, &Result)
		if err != nil {
			continue
		}

		for _, c := range r.Cookies() {
			if c.Name == "pwg_id" {
				p.Token = c.Value
				break
			}
		}

		break
	}

	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") == "1" {
		var RawResult interface{}
		err = json.Unmarshal(raw, RawResult)
		if err != nil {
			return err
		}

		DumpResponse(RawResult)
	}

	if Result.Stat != "ok" {
		return fmt.Errorf("[Error %d] %s", Result.Err, Result.ErrMessage)
	}

	return nil
}
