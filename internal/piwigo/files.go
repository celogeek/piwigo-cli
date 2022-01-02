package piwigo

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/celogeek/piwigo-cli/internal/piwigo/piwigotools"
	"golang.org/x/text/unicode/norm"
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

func (p *Piwigo) CheckUploadFile(file *piwigotools.FileToUpload, stat *piwigotools.FileToUploadStat) (err error) {
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

func (p *Piwigo) Upload(file *piwigotools.FileToUpload, stat *piwigotools.FileToUploadStat, nbJobs int, hasVideoJS bool) {
	err := p.CheckUploadFile(file, stat)
	if err != nil {
		return
	}
	wg := &sync.WaitGroup{}
	chunks, err := file.Base64Chunker()
	if err != nil {
		stat.Error("Base64Chunker", file.FullPath(), err)
		return
	}

	ok := true
	wg.Add(nbJobs)
	for j := 0; j < nbJobs; j++ {
		go p.UploadChunk(file, chunks, wg, stat, &ok)
	}
	wg.Wait()
	if !ok {
		return
	}

	// lock this process for committing the file

	var resp *FileUploadResult
	data := &url.Values{}
	data.Set("original_sum", file.MD5())
	data.Set("original_filename", file.Name)
	data.Set("check_uniqueness", "true")
	if file.CreatedAt() != nil {
		data.Set("date_creation", file.CreatedAt().String())
	}
	if file.CategoryId > 0 {
		data.Set("categories", fmt.Sprint(file.CategoryId))
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < 3; i++ {
		err = p.Post("pwg.images.add", data, &resp)
		if err == nil || err.Error() == "[Error 500] file already exists" {
			err = nil
			break
		}
		stat.Error(fmt.Sprintf("Upload %d", i), file.FullPath(), err)
	}

	if err != nil {
		stat.Fail()
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

func (p *Piwigo) UploadChunk(file *piwigotools.FileToUpload, chunks chan *piwigotools.FileToUploadChunk, wg *sync.WaitGroup, stat *piwigotools.FileToUploadStat, ok *bool) {
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
	filter *piwigotools.UploadFileType,
	stat *piwigotools.FileToUploadStat,
	files chan *piwigotools.FileToUpload,
) {
	if level == 0 {
		defer close(files)
	}
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		stat.Error("ScanTree Abs", rootPath, err)
		return
	}

	categoryFromName, err := p.CategoryFromName(parentCategoryId)
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
			category, ok := categoryFromName[dirname]
			if !ok {
				category = &piwigotools.Category{}
				p.mu.Lock()
				err = p.Post("pwg.categories.add", &url.Values{
					"name":   []string{strings.ReplaceAll(dirname, "'", `\'`)},
					"parent": []string{fmt.Sprint(parentCategoryId)},
				}, &category)
				p.mu.Unlock()
				if err != nil {
					stat.Error("ScanTree Categories Add", rootPath, err)
					return
				}
			}
			p.ScanTree(filepath.Join(rootPath, dirname), category.Id, level+1, filter, stat, files)
		case false: // File
			file := &piwigotools.FileToUpload{
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

func (p *Piwigo) CheckFiles(filesToCheck chan *piwigotools.FileToUpload, files chan *piwigotools.FileToUpload, stat *piwigotools.FileToUploadStat, nbJobs int) {
	defer close(files)

	wg := &sync.WaitGroup{}
	wg.Add(nbJobs)
	for i := 0; i < nbJobs; i++ {
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

func (p *Piwigo) UploadFiles(
	files chan *piwigotools.FileToUpload,
	stat *piwigotools.FileToUploadStat,
	hasVideoJS bool,
	nbJobs int,
	nbJobsChunk int,
) {
	defer stat.Close()

	wg := &sync.WaitGroup{}
	wg.Add(nbJobs)
	for i := 0; i < nbJobs; i++ {
		go func() {
			defer wg.Done()
			for file := range files {
				p.Upload(file, stat, nbJobsChunk, hasVideoJS)
			}
		}()
	}
	wg.Wait()
}
