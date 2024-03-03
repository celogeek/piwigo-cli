package piwigotools

import (
	"fmt"
	"os"
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

func NewFileToUploadStat() *FileToUploadStat {
	bar := progressbar.DefaultBytes(1, "...")
	progressbar.OptionOnCompletion(func() { _, _ = os.Stderr.WriteString("\n") })(bar)
	return &FileToUploadStat{
		Progress: bar,
	}
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
	defer s.mu.Unlock()

	s.Checked++
	s.Refresh()
}

func (s *FileToUploadStat) AddBytes(filesize int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalBytes += filesize
	s.Progress.ChangeMax64(s.TotalBytes + 1)
	s.Refresh()
}

func (s *FileToUploadStat) Add() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Total++
	s.Refresh()
}

func (s *FileToUploadStat) Commit(fileread int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.UploadedBytes += fileread
	_ = s.Progress.Set64(s.UploadedBytes)
}

func (s *FileToUploadStat) Done() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Uploaded++
	s.Refresh()
}

func (s *FileToUploadStat) Close() {
	_ = s.Progress.Close()
}

func (s *FileToUploadStat) Fail() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Failed++
	s.Refresh()
}

func (s *FileToUploadStat) Skip() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Skipped++
	s.Refresh()
}

func (s *FileToUploadStat) Error(origin string, filename string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = s.Progress.Clear()
	fmt.Printf("[%s] %s: %s\n", origin, filename, err)
	_ = s.Progress.RenderBlank()
}
