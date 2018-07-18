package gorawio

import (
	"io"
)

// ReaderAt changes the reference of ReadAt to offset
type ReaderAt struct {
	offset int64
	r      io.ReaderAt
}

// NewReaderAt creates a new ReaderAt
func NewReaderAt(r io.ReaderAt, offset int64) *ReaderAt {
	return &ReaderAt{
		r:      r,
		offset: offset,
	}
}

// ReadAt implements io.ReadAt within the new reference space
func (r *ReaderAt) ReadAt(p []byte, off int64) (int, error) {
	return r.r.ReadAt(p, off+r.offset)
}
