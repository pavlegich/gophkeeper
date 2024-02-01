// Package rwmanager contains object and methods for
// reading from the input and writing into the output.
package rwmanager

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

// RWManager contains reader and writer
// for interacting with input and output.
type RWManager struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// RWService describes methods for reading data from the input
// and writing data to the output.
type RWService interface {
	Read(ctx context.Context) (string, error)
	Write(ctx context.Context, out string) error
	WriteString(ctx context.Context, out string) error
}

// NewRWManager creates and returns new RWManager object.
func NewRWManager(ctx context.Context, in *os.File, out *os.File) RWService {
	return &RWManager{
		reader: bufio.NewReader(in),
		writer: bufio.NewWriter(out),
	}
}

// Read reads data from the input and returns it.
func (m *RWManager) Read(ctx context.Context) (string, error) {
	var in string
	_, err := fmt.Fscanln(m.reader, &in)
	if err != nil && len(in) != 0 {
		return "", fmt.Errorf("Read: read string from input failed %w", err)
	}
	return in, nil
}

// Write writes the requested text into the output.
func (m *RWManager) Write(ctx context.Context, out string) error {
	_, err := fmt.Fprintf(m.writer, "%s", out)
	if err != nil {
		return fmt.Errorf("Write: print into the output failed %w", err)
	}
	m.writer.Flush()
	return nil
}

// WriteString writes the requested text into the output from the new line.
func (m *RWManager) WriteString(ctx context.Context, out string) error {
	_, err := fmt.Fprintf(m.writer, "%s\n", out)
	if err != nil {
		return fmt.Errorf("WriteString: print into the output failed %w", err)
	}
	m.writer.Flush()
	return nil
}
