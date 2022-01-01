package piwigotools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type ImageDetails struct {
	Id            int        `json:"id"`
	Md5           string     `json:"md5sum"`
	Name          string     `json:"name"`
	DateAvailable TimeResult `json:"date_available"`
	DateCreation  TimeResult `json:"date_creation"`
	LastModified  TimeResult `json:"lastmodified"`
	Width         int        `json:"width"`
	Height        int        `json:"height"`
	Url           string     `json:"page_url"`
	ImageUrl      string     `json:"element_url"`
	Filename      string     `json:"file"`
	Filesize      int64      `json:"filesize"`
	Categories    Categories `json:"categories"`
	Tags          Tags       `json:"tags"`
	Derivatives   map[string]struct {
		Height int    `json:"height"`
		Width  int    `json:"width"`
		Url    string `json:"url"`
	} `json:"derivatives"`
}

func (img *ImageDetails) Preview(height int) (string, error) {
	url := img.ImageUrl
	if der, ok := img.Derivatives["medium"]; ok {
		url = der.Url
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("[error %d] failed to get image", resp.StatusCode)
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer([]byte{})

	buf.WriteString("\033]1337")
	buf.WriteString(fmt.Sprintf(";File=%s", img.Filename))
	buf.WriteString(";inline=1")
	buf.WriteString(fmt.Sprintf(";size=%d;", resp.ContentLength))
	if height > 0 {
		buf.WriteString(fmt.Sprintf(";height=%d", height))
	}
	buf.WriteString(":")

	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	defer encoder.Close()
	if _, err := io.Copy(encoder, resp.Body); err != nil {
		return "", err
	}
	buf.WriteString("\a")

	return buf.String(), nil
}
