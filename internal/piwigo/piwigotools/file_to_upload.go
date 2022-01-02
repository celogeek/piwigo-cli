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
	md5       string
	size      int64
	ext       string
	createdAt *TimeResult
}

type FileToUpload struct {
	Dir        string
	Name       string
	CategoryId int

	info *FileInfo
}

func (f *FileToUpload) FullPath() string {
	return filepath.Join(f.Dir, f.Name)
}

func (f *FileToUpload) Info() *FileInfo {
	if f.info != nil {
		return f.info
	}

	file, err := os.Open(f.FullPath())
	if err != nil {
		return nil
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		return nil
	}

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return nil
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	info := FileInfo{
		size:      st.Size(),
		md5:       checksum,
		createdAt: f.exifCreatedAt(),
	}

	f.info = &info
	return f.info
}

func (f *FileToUpload) Checked() bool {
	return f.info != nil
}

func (f *FileToUpload) MD5() string {
	if info := f.Info(); info != nil {
		return info.md5
	}
	return ""
}

func (f *FileToUpload) Size() int64 {
	if info := f.Info(); info != nil {
		return info.size
	}
	return -1
}

func (f *FileToUpload) Ext() string {
	return strings.ToLower(filepath.Ext(f.Name)[1:])
}

func (f *FileToUpload) CreatedAt() *TimeResult {
	if info := f.Info(); info != nil {
		return info.createdAt
	}
	return nil
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

func (f *FileToUpload) Base64Chunker() (chan *FileToUploadChunk, error) {
	fh, err := os.Open(f.FullPath())
	if err != nil {
		return nil, err
	}

	out := make(chan *FileToUploadChunk, 8)
	chunker := func() {
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
	}

	go chunker()

	return out, nil
}

func (f *FileToUpload) exifCreatedAt() *TimeResult {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil
	}
	defer et.Close()

	var createdAt *time.Time
	var CreateDateFormat = "2006:01:02 15:04:05-07:00"

	fileInfos := et.ExtractMetadata(f.FullPath())
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
		return &result
	}
	return nil
}
