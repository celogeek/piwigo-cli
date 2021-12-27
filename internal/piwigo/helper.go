package piwigo

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

var CHUNK_BUFF_SIZE int64 = 32 * 1024
var CHUNK_BUFF_COUNT int = 32
var CHUNK_PRECOMPUTE_SIZE int = 8

func DumpResponse(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func ArgsToForm(args []string) (*url.Values, error) {
	params := &url.Values{}
	for _, arg := range args {
		r := strings.SplitN(arg, "=", 2)
		if len(r) != 2 {
			return nil, errors.New("args should be key=value")
		}
		params.Add(r[0], strings.ReplaceAll(r[1], "'", `\'`))
	}
	return params, nil
}

func Md5File(filename string) (result string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}
	result = fmt.Sprintf("%x", hash.Sum(nil))
	return
}

type Base64ChunkResult struct {
	Position int64
	Size     int64
	Buffer   bytes.Buffer
}

func Base64Chunker(filename string) (out chan *Base64ChunkResult, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	out = make(chan *Base64ChunkResult, CHUNK_PRECOMPUTE_SIZE)
	go func() {
		b := make([]byte, CHUNK_BUFF_SIZE)
		defer f.Close()
		defer close(out)
		ok := false
		for position := int64(0); !ok; position += 1 {
			bf := &Base64ChunkResult{
				Position: position,
			}
			b64 := base64.NewEncoder(base64.StdEncoding, &bf.Buffer)
			for i := 0; i < CHUNK_BUFF_COUNT; i++ {
				n, _ := f.Read(b)
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

	return
}
