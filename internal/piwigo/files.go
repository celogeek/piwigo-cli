package piwigo

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type FileUploadResult struct {
	ImageId int    `json:"image_id"`
	Url     string `json:"url"`
}

func (p *Piwigo) FileExists(md5 string) bool {
	var resp map[string]*string

	if err := p.Post("pwg.images.exist", &url.Values{
		"md5sum_list": []string{md5},
	}, &resp); err != nil {
		return false
	}

	return resp[md5] != nil
}

func (p *Piwigo) UploadChunks(filename string, nbJobs int, categoryId int) (*FileUploadResult, error) {
	md5, err := Md5File(filename)
	if err != nil {
		return nil, err
	}

	if p.FileExists(md5) {
		return nil, errors.New("file already exists")
	}

	st, _ := os.Stat(filename)
	wg := &sync.WaitGroup{}
	chunks, err := Base64Chunker(filename)
	errout := make(chan error)
	bar := progressbar.DefaultBytes(
		st.Size(),
		"uploading",
	)
	if err != nil {
		return nil, err
	}

	for j := 0; j < nbJobs; j++ {
		wg.Add(1)
		go p.UploadChunk(md5, chunks, errout, wg, bar)
	}
	go func() {
		wg.Wait()
		bar.Close()
		close(errout)
	}()

	var errstring string
	for err := range errout {
		errstring += err.Error() + "\n"
	}
	if errstring != "" {
		return nil, errors.New(errstring)
	}

	exif, _ := Exif(filename)
	var resp *FileUploadResult
	data := &url.Values{}
	data.Set("original_sum", md5)
	data.Set("original_filename", filepath.Base(filename))
	data.Set("check_uniqueness", "true")
	if exif != nil && exif.CreatedAt != nil {
		data.Set("date_creation", exif.CreatedAt.String())
	}
	if categoryId > 0 {
		data.Set("categories", fmt.Sprint(categoryId))
	}
	err = p.Post("pwg.images.add", data, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *Piwigo) UploadChunk(md5 string, chunks chan *Base64ChunkResult, errout chan error, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for chunk := range chunks {
		var err error
		data := &url.Values{
			"original_sum": []string{md5},
			"position":     []string{fmt.Sprint(chunk.Position)},
			"type":         []string{"file"},
			"data":         []string{chunk.Buffer.String()},
		}
		for i := 0; i < 3; i++ {
			err = p.Post("pwg.images.addChunk", data, nil)
			if err == nil {
				break
			}
		}
		bar.Add64(chunk.Size)
		if err != nil {
			errout <- fmt.Errorf("error on chunk %d: %v", chunk.Position, err)
			continue
		}
	}
}
