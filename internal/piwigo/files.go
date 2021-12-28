package piwigo

import (
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

func (p *Piwigo) CheckUploadFile(file *FileToUpload, stat *FileToUploadStat) (err error) {
	if !file.Checked() {
		if file.MD5() == "" {
			stat.Fail()
			stat.Check()
			err = fmt.Errorf("%s: checksum error", file.FullPath())
			stat.Error("CheckUploadFile", file.FullPath(), err)
			return
		}

		if p.FileExists(file.MD5()) {
			stat.Skip()
			stat.Check()
			err = fmt.Errorf("%s: file already exists", file.FullPath())
			return
		}

		stat.Check()
		stat.AddBytes(file.Size())
	}
	return nil
}

func (p *Piwigo) Upload(file *FileToUpload, stat *FileToUploadStat, nbJobs int, hasVideoJS bool) {
	err := p.CheckUploadFile(file, stat)
	if err != nil {
		return
	}
	wg := &sync.WaitGroup{}
	chunks, err := Base64Chunker(file.FullPath())
	if err != nil {
		stat.Error("Base64Chunker", file.FullPath(), err)
		return
	}

	ok := true
	for j := 0; j < nbJobs; j++ {
		wg.Add(1)
		go p.UploadChunk(file, chunks, wg, stat, &ok)
	}
	wg.Wait()
	if !ok {
		return
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
		stat.Error("Upload", file.FullPath(), err)
		return
	}

	if hasVideoJS {
		switch file.Ext() {
		case "ogg", "ogv", "mp4", "m4v", "webm", "webmv":
			p.VideoJSSync(resp.ImageId)
		}
	}

	stat.Done()
}

func (p *Piwigo) UploadChunk(file *FileToUpload, chunks chan *Base64ChunkResult, wg *sync.WaitGroup, stat *FileToUploadStat, ok *bool) {
	defer wg.Done()
	for chunk := range chunks {
		var err error
		data := &url.Values{
			"original_sum": []string{file.MD5()},
			"position":     []string{fmt.Sprint(chunk.Position)},
			"type":         []string{"file"},
			"data":         []string{chunk.Buffer.String()},
		}
		for i := 0; i < 3; i++ {
			err = p.Post("pwg.images.addChunk", data, nil)
			if err == nil {
				break
			}
			stat.Error(fmt.Sprintf("UploadChunk %d", i), file.FullPath(), err)
		}
		stat.Commit(chunk.Size)
		if err != nil {
			stat.Fail()
			*ok = false
			return
		}
	}
}

func (p *Piwigo) ScanTree(
	rootPath string,
	parentCategoryId int,
	level int,
	filter *UploadFileType,
	stat *FileToUploadStat,
	files chan *FileToUpload,
) {
	if level == 0 {
		defer close(files)
	}
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		stat.Error("ScanTree Abs", rootPath, err)
		return
	}

	categoriesId, err := p.CategoriesId(parentCategoryId)
	if err != nil {
		stat.Error("ScanTree CategoriesId", rootPath, err)
		return
	}

	dirs, err := ioutil.ReadDir(rootPath)
	if err != nil {
		stat.Error("ScanTree Dir", rootPath, err)
		return
	}

	for _, dir := range dirs {
		switch dir.IsDir() {
		case true: // Directory
			dirname := norm.NFC.String(dir.Name())
			categoryId, ok := categoriesId[dirname]
			if !ok {
				var resp struct {
					Id int `json:"id"`
				}
				err = p.Post("pwg.categories.add", &url.Values{
					"name":   []string{strings.ReplaceAll(dirname, "'", `\'`)},
					"parent": []string{fmt.Sprint(parentCategoryId)},
				}, &resp)
				if err != nil {
					stat.Error("ScanTree Categories Add", rootPath, err)
					return
				}
				categoryId = resp.Id
			}
			p.ScanTree(filepath.Join(rootPath, dirname), categoryId, level+1, filter, stat, files)
		case false: // File
			file := &FileToUpload{
				Dir:        rootPath,
				Name:       dir.Name(),
				CategoryId: parentCategoryId,
			}
			if !filter.Has(file.Ext()) {
				continue
			}
			stat.Add()
			files <- file
		}
	}

}

func (p *Piwigo) CheckFiles(filesToCheck chan *FileToUpload, files chan *FileToUpload, stat *FileToUploadStat, nbJobs int) {
	defer close(files)

	wg := &sync.WaitGroup{}
	for i := 0; i < nbJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesToCheck {
				err := p.CheckUploadFile(file, stat)
				if err != nil {
					continue
				}
				files <- file
			}
		}()
	}

	wg.Wait()
}

func (p *Piwigo) UploadFiles(files chan *FileToUpload, stat *FileToUploadStat, hasVideoJS bool, nbJobs int) {
	defer stat.Close()

	wg := &sync.WaitGroup{}
	for i := 0; i < nbJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range files {
				p.Upload(file, stat, 2, hasVideoJS)
			}
		}()
	}
	wg.Wait()
}
