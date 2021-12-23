package piwigo

import (
	"fmt"
	"time"

	"github.com/barasher/go-exiftool"
)

type ExifResult struct {
	CreatedAt *TimeResult
}

func Exif(filename string) (*ExifResult, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil, err
	}
	defer et.Close()

	resp := &ExifResult{}
	fileInfos := et.ExtractMetadata(filename)
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			switch k {
			case "CreateDate", "CreationDate":
				switch v := v.(type) {
				case string:
					t, err := time.Parse("2006:01:02 15:04:05-07:00", v)
					if err == nil {
						if resp.CreatedAt == nil || time.Time(*resp.CreatedAt).After(t) {
							r := TimeResult(t)
							resp.CreatedAt = &r
						}
					}
				}
			}
		}
	}
	return resp, nil
}
