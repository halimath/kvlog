package kvlog

import (
	"fmt"
	"io"
	"time"
)

const (
	KeyLevel     = "level"
	KeyTimestamp = "ts"
)

// Message represents a single log message expressed as an ordered list of key value pairs
type Message []KVPair

// Level returns the message's level
func (m Message) Level() Level {
	for _, p := range m {
		switch l := p.Value.(type) {
		case Level:
			return l
		}
	}

	return LevelDebug
}

func (m Message) WriteTo(w io.Writer) error {
	for i, p := range m {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		p.WriteTo(w)
	}

	return nil
}

// NewMessage creates a new message from the given log level and key-value pairs
func NewMessage(l Level, pairs ...KVPair) Message {
	m := make(Message, len(pairs)+2)
	m[0] = KV(KeyTimestamp, time.Now())
	m[1] = KV(KeyLevel, l)

	copy(m[2:], pairs)

	return m
}
