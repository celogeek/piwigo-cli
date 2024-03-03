package piwigotools

import (
	"fmt"
	"github.com/apoorvam/goterminal"
	ct "github.com/daviddengcn/go-colortext"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	megabytes = 1 << 20
	kilobytes = 1 << 10
)

type FileToUploadStat struct {
	Checked  uint32
	Total    uint32
	Uploaded uint32
	Skipped  uint32
	Failed   uint32

	UploadedBytes int64
	TotalBytes    int64

	progress    *goterminal.Writer
	mu          sync.Mutex
	lastRefresh time.Time
}

func NewFileToUploadStat() *FileToUploadStat {
	return &FileToUploadStat{
		progress: goterminal.New(os.Stdout),
	}
}

func (s *FileToUploadStat) formatBytes(b int64) string {
	if b > megabytes {
		return fmt.Sprintf("%.2fMB", float64(b)/(1<<20))
	} else if b > kilobytes {
		return fmt.Sprintf("%.2fKB", float64(b)/(1<<10))
	} else {
		return fmt.Sprintf("%dB", b)
	}
}

func (s *FileToUploadStat) refreshIfNeeded() {
	if time.Since(s.lastRefresh) > 200*time.Millisecond {
		s.refresh()
	}

}

func (s *FileToUploadStat) refresh() {
	s.progress.Clear()
	_, _ = s.progress.Write([]byte("Statistics:\n"))
	_, _ = fmt.Fprintf(
		s.progress,
		strings.Repeat("%20s: %d\n", 5)+strings.Repeat("%20s: %s\n", 2),
		"Checked", s.Checked,
		"Uploaded", s.Uploaded,
		"Skipped", s.Skipped,
		"Failed", s.Failed,
		"Total", s.Total,
		"Uploaded Bytes", s.formatBytes(s.UploadedBytes),
		"Total Bytes", s.formatBytes(s.TotalBytes),
	)
	_ = s.progress.Print()
	s.lastRefresh = time.Now()
}

func (s *FileToUploadStat) Check() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Checked++
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Add() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Total++
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) AddBytes(filesize int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalBytes += filesize
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Commit(fileread int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.UploadedBytes += fileread
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Done() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Uploaded++
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Fail() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Failed++
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Skip() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Skipped++
	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Error(origin string, filename string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.progress.Clear()
	ct.Foreground(ct.Red, false)
	_, _ = fmt.Fprintf(s.progress, "[%s] %s: %s\n", origin, filename, err)
	_ = s.progress.Print()
	ct.ResetColor()
	s.progress.Reset()

	s.refreshIfNeeded()
}

func (s *FileToUploadStat) Close() {
	s.refresh()
}
