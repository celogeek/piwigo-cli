package piwigo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/text/unicode/norm"
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

func (p *Piwigo) Upload(file *FileToUpload, stat *FileToUploadStat, nbJobs int, hasVideoJS bool) error {
	if !file.Checked() {
		if file.MD5() == "" {
			stat.Fail()
			return errors.New("checksum error")
		}
		stat.Check()
	}

	if p.FileExists(file.MD5()) {
		stat.Skip()
		return errors.New("file already exists")
	}

	wg := &sync.WaitGroup{}
	chunks, err := Base64Chunker(file.FullPath())
	errout := make(chan error)
	if err != nil {
		return err
	}

	for j := 0; j < nbJobs; j++ {
		wg.Add(1)
		go p.UploadChunk(file.MD5(), chunks, errout, wg, stat)
	}
	go func() {
		wg.Wait()
		close(errout)
	}()

	var errstring string
	for err := range errout {
		errstring += err.Error() + "\n"
	}
	if errstring != "" {
		stat.Fail()
		return errors.New(errstring)
	}

	exif, _ := Exif(file.FullPath())
	var resp *FileUploadResult
	data := &url.Values{}
	data.Set("original_sum", file.MD5())
	data.Set("original_filename", file.Name)
	data.Set("check_uniqueness", "true")
	if exif != nil && exif.CreatedAt != nil {
		data.Set("date_creation", exif.CreatedAt.String())
	}
	if file.CategoryId > 0 {
		data.Set("categories", fmt.Sprint(file.CategoryId))
	}
	err = p.Post("pwg.images.add", data, &resp)
	if err != nil {
		stat.Fail()
		return err
	}

	if hasVideoJS {
		switch file.Ext() {
		case "ogg", "ogv", "mp4", "m4v", "webm", "webmv":
			p.VideoJSSync(resp.ImageId)
		}
	}

	stat.Done()
	return nil
}

func (p *Piwigo) UploadChunk(md5 string, chunks chan *Base64ChunkResult, errout chan error, wg *sync.WaitGroup, progress *FileToUploadStat) {
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
		progress.Commit(chunk.Size)
		if err != nil {
			errout <- fmt.Errorf("error on chunk %d: %v", chunk.Position, err)
			continue
		}
	}
}

func (p *Piwigo) UploadTree(rootPath string, parentCategoryId int, level int, filter UploadFileType) ([]FileToUpload, error) {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	categoriesId, err := p.CategoriesId(parentCategoryId)
	if err != nil {
		return nil, err
	}

	dirs, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	var files []FileToUpload

	levelStr := strings.Repeat("  ", level)
	for _, dir := range dirs {
		if !dir.IsDir() {
			ext := strings.ToLower(filepath.Ext(dir.Name())[1:])
			if !filter.Has(ext) {
				continue
			}
			filename := filepath.Join(rootPath, dir.Name())
			md5, err := Md5File(filename)
			if err != nil {
				return nil, err
			}
			status := "OK"
			if p.FileExists(md5) {
				status = "SKIP"
			}
			fmt.Printf("%s - %s %s - %s\n", levelStr, dir.Name(), md5, status)
			if status == "OK" {
				files = append(files, FileToUpload{
					Dir:        rootPath,
					Name:       dir.Name(),
					CategoryId: parentCategoryId,
				})
			}
			continue
		}
		dirname := norm.NFC.String(dir.Name())
		categoryId, ok := categoriesId[dirname]
		fmt.Printf("%s%s\n", levelStr, dirname)
		if !ok {
			var resp struct {
				Id int `json:"id"`
			}
			err = p.Post("pwg.categories.add", &url.Values{
				"name":   []string{strings.ReplaceAll(dirname, "'", `\'`)},
				"parent": []string{fmt.Sprint(parentCategoryId)},
			}, &resp)
			if err != nil {
				return nil, err
			}
			categoryId = resp.Id
		}
		newFiles, err := p.UploadTree(filepath.Join(rootPath, dirname), categoryId, level+1, filter)
		if err != nil {
			return nil, err
		}
		files = append(files, newFiles...)
	}

	return files, nil
}
