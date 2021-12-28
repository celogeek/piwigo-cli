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
	s.Progress.Describe(fmt.Sprintf(
		"%d / %d - check:%d, upload:%d, skip:%d, fail:%d",
		s.Uploaded+s.Skipped+s.Failed,
		s.Total,
		s.Checked,
		s.Uploaded,
		s.Skipped,
		s.Failed,
	),
	)
}

func (s *FileToUploadStat) Check() {
	s.mu.Lock()
	s.Checked++
	s.Refresh()
	s.mu.Unlock()
}

func (s *FileToUploadStat) AddBytes(filesize int64) {
	s.mu.Lock()
	s.TotalBytes += filesize
	s.Progress.ChangeMax64(s.TotalBytes + 1)
	s.Refresh()
	s.mu.Unlock()
}

func (s *FileToUploadStat) Add() {
	s.mu.Lock()
	s.Total++
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

func (s *FileToUploadStat) Error(origin string, filename string, err error) error {
	s.mu.Lock()
	s.Progress.Clear()
	fmt.Printf("[%s] %s: %s\n", origin, filename, err)
	s.Progress.RenderBlank()
	s.mu.Unlock()
	return err
}
