package piwigo

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/schollz/progressbar/v3"
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

func (p *Piwigo) UploadChunks(filename string, nbJobs int) error {
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

	in := make(chan int64)
	out := make(chan error)
	wg := &sync.WaitGroup{}
	bar := progressbar.DefaultBytes(
		st.Size(),
		"uploading",
	)

	for j := 0; j < nbJobs; j++ {
		wg.Add(1)
		go p.UploadChunk(md5, f, in, out, wg, bar)
	}

	go func() {
		nbChunks := st.Size()/CHUNK_SIZE + 1
		for position := int64(0); position < nbChunks; position++ {
			in <- position
		}
		close(in)
		wg.Wait()
		close(out)
		bar.Close()
	}()

	var errString string
	for err := range out {
		errString += err.Error() + "\n"
	}
	if errString != "" {
		return errors.New(errString[:len(errString)-1])
	}

	fmt.Println(md5)

	return nil
}

func (p *Piwigo) UploadChunk(md5 string, f *os.File, in chan int64, out chan error, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for position := range in {
		n, b64, err := Base64Chunk(f, position)
		if err != nil {
			out <- fmt.Errorf("error on chunk %d: %v", position, err)
			continue
		}

		err = p.Post("pwg.images.addChunk", &url.Values{
			"original_sum": []string{md5},
			"position":     []string{fmt.Sprint(position)},
			"type":         []string{"file"},
			"data":         []string{b64},
		}, nil)
		if err != nil {
			out <- fmt.Errorf("error on chunk %d: %v", position, err)
			continue
		}
		bar.Add(n)
	}
}
