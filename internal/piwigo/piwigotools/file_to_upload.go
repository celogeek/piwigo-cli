package piwigotools

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

type FileInfo struct {
	fullpath  *string
	md5       *string
	size      *int64
	ext       *string
	createdAt *TimeResult
}

type FileToUpload struct {
	Dir        string
	Name       string
	CategoryId int

	info FileInfo
}

func (f *FileToUpload) FullPath() *string {
	if f.info.fullpath != nil {
		return f.info.fullpath
	}
	fp := filepath.Join(f.Dir, f.Name)
	f.info.fullpath = &fp

	return f.info.fullpath
}

func (f *FileToUpload) Size() *int64 {
	if f.info.size != nil {
		return f.info.size
	}

	st, err := os.Stat(*f.FullPath())
	if err != nil {
		return nil
	}
	size := st.Size()

	f.info.size = &size

	return f.info.size
}

func (f *FileToUpload) Ext() *string {
	if f.info.ext != nil {
		return f.info.ext
	}

	ext := strings.ToLower(filepath.Ext(f.Name)[1:])
	f.info.ext = &ext

	return f.info.ext
}

func (f *FileToUpload) MD5() *string {
	if f.info.md5 != nil {
		return f.info.md5
	}
	file, err := os.Open(*f.FullPath())
	if err != nil {
		return nil
	}
	defer file.Close()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return nil
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	f.info.md5 = &checksum

	return f.info.md5
}

func (f *FileToUpload) CreatedAt() *TimeResult {
	if f.info.createdAt != nil {
		return f.info.createdAt
	}

	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil
	}
	defer et.Close()

	var createdAt *time.Time
	var CreateDateFormat = "2006:01:02 15:04:05-07:00"

	fileInfos := et.ExtractMetadata(*f.FullPath())
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
			if createdAt == nil || createdAt.After(t) {
				createdAt = &t
			}
		}
	}

	if createdAt != nil {
		result := TimeResult(*createdAt)
		f.info.createdAt = &result
	}

	return f.info.createdAt
}

func (f *FileToUpload) Checked() bool {
	return f.info.md5 != nil
}

var (
	CHUNK_SIZE       int64 = 1 * 1024 * 1024
	CHUNK_BUFF_SIZE  int64 = 32 * 1024
	CHUNK_BUFF_COUNT       = CHUNK_SIZE / CHUNK_BUFF_SIZE
)

type FileToUploadChunk struct {
	Position int64
	Size     int64
	Buffer   bytes.Buffer
}

func (f *FileToUpload) Base64BuildChunk() (chan *FileToUploadChunk, error) {
	fh, err := os.Open(*f.FullPath())
	if err != nil {
		return nil, err
	}

	out := make(chan *FileToUploadChunk, 8)
	go func() {
		b := make([]byte, CHUNK_BUFF_SIZE)
		defer fh.Close()
		defer close(out)
		ok := false
		for position := int64(0); !ok; position += 1 {
			bf := &FileToUploadChunk{
				Position: position,
			}
			b64 := base64.NewEncoder(base64.StdEncoding, &bf.Buffer)
			for i := int64(0); i < CHUNK_BUFF_COUNT; i++ {
				n, _ := fh.Read(b)
				if n == 0 {
					ok = true
					break
				}
				bf.Size += int64(n)
				b64.Write(b[:n])
			}
			b64.Close()
			if bf.Size > 0 {
				out <- bf
			}
		}
	}()

	return out, nil
}
