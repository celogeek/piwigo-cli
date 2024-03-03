package piwigo

import (
	"fmt"
	"net/http"
	"net/url"
)

func (p *Piwigo) VideoJSSync(imageId int) error {
	Url, err := url.Parse(p.Url)
	if err != nil {
		return err
	}
	Url.Path += "/admin.php"
	q := Url.Query()
	q.Set("page", "plugin")
	q.Set("section", "piwigo-videojs/admin/admin_photo.php")
	q.Set("sync_metadata", "1")
	q.Set("image_id", fmt.Sprint(imageId))
	Url.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", Url.String(), nil)
	if err != nil {
		return err
	}
	if p.Token != "" {
		req.AddCookie(&http.Cookie{Name: "pwg_id", Value: p.Token, HttpOnly: true})
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	return nil
}
