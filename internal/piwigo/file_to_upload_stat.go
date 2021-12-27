package piwigo

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type FileToUploadStat struct {
	Checked       int64
	Total         int64
	TotalBytes    int64
	Uploaded      int64
	UploadedBytes int64
	Skipped       int64
	Failed        int64
	Progress      *progressbar.ProgressBar
	mu            sync.Mutex
}

func (s *FileToUploadStat) Refresh() {
	s.Progress.Describe(fmt.Sprintf("check:%d, upload:%d, skip:%d, failed:%d, total:%d", s.Checked, s.Uploaded, s.Skipped, s.Failed, s.Total))
}

func (s *FileToUploadStat) Check() {
	s.mu.Lock()
	s.Checked++
	s.Refresh()
	s.mu.Unlock()
}

func (s *FileToUploadStat) Add(filesize int64) {
	s.mu.Lock()
	s.Total++
	s.TotalBytes += filesize
	s.Progress.ChangeMax64(s.TotalBytes)
	s.Refresh()
	s.mu.Unlock()
}

func (s *FileToUploadStat) Commit(filereaded int64) {
	s.mu.Lock()
	s.UploadedBytes += filereaded
	s.Progress.Set64(s.UploadedBytes)
	s.mu.Unlock()
}

func (s *FileToUploadStat) Done() {
	s.mu.Lock()
	s.Uploaded++
	s.Refresh()
	s.mu.Unlock()
}

func (s *FileToUploadStat) Close() {
	s.Progress.Close()
}

func (s *FileToUploadStat) Fail() {
	s.mu.Lock()
	s.Failed++
	s.Refresh()
	s.mu.Unlock()
}

func (s *FileToUploadStat) Skip() {
	s.mu.Lock()
	s.Skipped++
	s.Refresh()
	s.mu.Unlock()
}
