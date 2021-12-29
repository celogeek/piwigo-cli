package piwigotools

import (
	"os"
	"path/filepath"
	"strings"
)

type FileToUpload struct {
	Dir        string
	Name       string
	CategoryId int

	md5  *string
	size *int64
	ext  *string
}

func (f *FileToUpload) FullPath() string {
	return filepath.Join(f.Dir, f.Name)
}

func (f *FileToUpload) Checked() bool {
	return f.md5 != nil
}

func (f *FileToUpload) MD5() string {
	if f.md5 == nil {
		md5, err := Md5File(f.FullPath())
		if err != nil {
			return ""
		}
		f.md5 = &md5
	}
	return *f.md5
}

func (f *FileToUpload) Size() int64 {
	if f.size == nil {
		st, err := os.Stat(f.FullPath())
		if err != nil {
			return -1
		}
		size := st.Size()
		f.size = &size
	}
	return *f.size
}

func (f *FileToUpload) Ext() string {
	if f.ext == nil {
		ext := strings.ToLower(filepath.Ext(f.Name)[1:])
		f.ext = &ext
	}
	return *f.ext
}
