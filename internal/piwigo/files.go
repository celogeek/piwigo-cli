package piwigo

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

func (p *Piwigo) FileExists(md5 string) bool {
	var resp map[string]*string

	if err := p.Post("pwg.images.exist", &url.Values{
		"md5sum_list": []string{md5},
	}, &resp); err != nil {
		return false
	}

	return resp[md5] != nil
}

func (p *Piwigo) UploadChunks(filename string) error {
	md5, err := Md5File(filename)
	if err != nil {
		return err
	}

	if p.FileExists(md5) {
		return errors.New("file already exists")
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return err
	}

	nbChunks := st.Size()/CHUNK_SIZE + 1
	for position := int64(0); position < nbChunks; position++ {
		b64, err := Base64Chunk(f, int64(position))
		if err != nil {
			return err
		}

		err = p.Post("pwg.images.addChunk", &url.Values{
			"original_sum": []string{md5},
			"position":     []string{fmt.Sprint(position)},
			"type":         []string{"file"},
			"data":         []string{b64},
		}, nil)
		if err != nil {
			return err
		}
		fmt.Printf("Upload %d/%d ok\n", position+1, nbChunks)
	}
	fmt.Println(md5)

	return nil
}
