package exif

import (
	"fmt"
	"time"

	"github.com/barasher/go-exiftool"
)

type Info struct {
	CreatedAt *time.Time
}

var (
	CreateDateFormat = "2006:01:02 15:04:05-07:00"
)

func Extract(filename string) (*Info, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil, err
	}
	defer et.Close()

	var resp *Info = &Info{}
	fileInfos := et.ExtractMetadata(filename)
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			continue
		}

		var t time.Time
		for k, v := range fileInfo.Fields {
			switch k {
			case "CreateDate":
				offset, ok := fileInfo.Fields["OffsetTime"]
				if !ok {
					offset = "+00:00"
				}
				v := fmt.Sprintf("%s%s", v, offset)
				t, err = time.Parse(CreateDateFormat, v)
			case "CreationDate":
				t, err = time.Parse(CreateDateFormat, fmt.Sprint(v))
			default:
				continue
			}
			if err != nil {
				continue
			}
			if resp.CreatedAt == nil || resp.CreatedAt.After(t) {
				resp.CreatedAt = &t
			}
		}
	}
	return resp, nil
}
