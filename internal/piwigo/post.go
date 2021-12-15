package piwigo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	req, err := http.NewRequest("POST", Url, strings.NewReader(encodedForm))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if p.Token != nil {
		req.AddCookie(p.Token)
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	Result := PiwigoResult{
		Result: resp,
	}

	if os.Getenv("DEBUG") == "1" {
		newBody := &bytes.Buffer{}
		tee := io.TeeReader(r.Body, newBody)

		var RawResult map[string]interface{}
		err = json.NewDecoder(tee).Decode(&RawResult)
		if err != nil {
			return err
		}
		DumpResponse(RawResult)

		err = json.NewDecoder(newBody).Decode(&Result)
		if err != nil {
			return err
		}
	} else {
		err = json.NewDecoder(r.Body).Decode(&Result)
		if err != nil {
			return err
		}
	}

	if Result.Stat != "ok" {
		return fmt.Errorf("[Error %d] %s", Result.Err, Result.ErrMessage)
	}

	for _, c := range r.Cookies() {
		if c.Name == "pwg_id" {
			p.Token = c
			break
		}
	}

	return nil
}
