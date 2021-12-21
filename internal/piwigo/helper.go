package piwigo

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"strings"
)

var CHUNK_SIZE int64 = int64(math.Pow(1024, 2))

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
		params.Add(r[0], r[1])
	}
	return params, nil
}

func Md5File(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func Base64Chunk(file *os.File, position int64) (int, string, error) {
	b := make([]byte, CHUNK_SIZE)
	n, err := file.ReadAt(b, position*CHUNK_SIZE)
	if err != nil && err != io.EOF {
		return 0, "", err
	}
	if n == 0 {
		return 0, "", errors.New("position out of bound")
	}
	return n, base64.StdEncoding.EncodeToString(b[:n]), nil
}
