package kvlog

import (
	"bytes"
	"io"
	"os"
)

// Output describes the interface to be implemented by
// log output streams
type Output interface {
	WriteLogMessage(m Message)
}

// WriterLogOutput implements a Output that writes to
// an io.Writer
type WriterLogOutput struct {
	W io.Writer
}

// WriteLogMessage writes the bytes to the writer
func (w *WriterLogOutput) WriteLogMessage(m Message) {
	var buf bytes.Buffer
	m.WriteTo(&buf)
	buf.WriteString("\n")

	w.W.Write(buf.Bytes())
}

// Stdout returns an Output that sends log messages to STDOUT.
func Stdout() Output {
	return &WriterLogOutput{
		W: os.Stdout,
	}
}

// Stderr returns an Output that sends log messages to STDERR.
func Stderr() Output {
	return &WriterLogOutput{
		W: os.Stderr,
	}
}
