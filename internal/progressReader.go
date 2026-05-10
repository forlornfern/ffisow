package internal

import (
	"fmt"
	"io"
	"os"
)

type ProgressReader struct {
	Reader  io.Reader
	Total   int64
	Written int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Written += int64(n)
	percent := float64(pr.Written) / float64(pr.Total) * 100
	mb := pr.Written / 1024 / 1024
	totalMb := pr.Total / 1024 / 1024

	fmt.Fprintf(os.Stderr, "\r%.1f%%  %d / %d MB", percent, mb, totalMb)

	if pr.Written >= pr.Total && err == io.EOF {
		fmt.Fprintln(os.Stderr)
	}
	return n, err
}
