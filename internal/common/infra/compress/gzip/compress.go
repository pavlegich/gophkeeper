// Package compress contains compress writer and reader implementation.
package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// compressWriter implements http.ResponseWriter interface
// for compressing the data.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter returns new compress writer.
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns response header.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write returns implements method Write of gzip.Writer
// and returns it's result.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sets Content-Encoding header and status code.
func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Close closes gzip.Writer and sends all data from the buffer.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader implements io.ReadCloser interface for
// decompressing the data obtained from the client.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader returns new gzip reader.
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("NewCompressReader: gzip new reader %w", err)
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read implements io.Reader from gzip.Reader, reading uncompressed bytes
// from its underlying Reader and returns it's result.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close return the result of closing io.ReadCloser and gzip.Reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("Close: reader close error %w", err)
	}
	return c.zr.Close()
}
