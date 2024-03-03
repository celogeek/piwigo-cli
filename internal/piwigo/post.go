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
	"time"

	"github.com/celogeek/piwigo-cli/internal/debug"
)

type PostResult struct {
	Stat       string      `json:"stat"`
	Err        int         `json:"err"`
	ErrMessage string      `json:"message"`
	Result     interface{} `json:"result"`
}

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

	result := PostResult{
		Result: resp,
	}

	raw := bytes.NewBuffer([]byte{})

	for i := range 3 {
		if i > 0 {
			time.Sleep(time.Second) // wait 1 sec before retry
		}

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

		_, err = io.Copy(raw, r.Body)
		if err != nil {
			_ = r.Body.Close()
			continue
		}

		err = json.Unmarshal(raw.Bytes(), &result)
		if err != nil {
			_ = r.Body.Close()
			continue
		}

		for _, c := range r.Cookies() {
			if c.Name == "pwg_id" {
				p.Token = c.Value
				break
			}
		}

		_ = r.Body.Close()
		break
	}

	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") == "1" {
		var RawResult interface{}
		err = json.Unmarshal(raw.Bytes(), &RawResult)
		if err != nil {
			return err
		}

		fmt.Println(debug.Dump(RawResult))
	}

	if result.Stat != "ok" {
		return fmt.Errorf("[Error %d] %s", result.Err, result.ErrMessage)
	}

	return nil
}
