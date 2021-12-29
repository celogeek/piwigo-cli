package piwigotools

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

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
