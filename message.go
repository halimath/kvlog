package kvlog

import (
	"time"
)

// Level defines the valid log levels.
type Level int

// String provides a string representation of the log level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}

const (
	// LevelDebug log level
	LevelDebug Level = iota
	// LevelInfo log level
	LevelInfo
	// LevelWarn log level
	LevelWarn
	// LevelError log level
	LevelError
)

// --

// KVPair implements a key-value pair
type KVPair struct {
	// Key stores the key of the pair
	Key string

	// Value stores the value
	Value interface{}
}

// KV is a factory method for KVPair values.
func KV(key string, value interface{}) KVPair {
	return KVPair{
		Key:   key,
		Value: value,
	}
}

// Event is a factory method for a KVPair that uses the default event key.
func Event(value interface{}) KVPair {
	return KVPair{
		Key:   KeyEvent,
		Value: value,
	}
}

// --

const (
	// KeyLevel defines the message key containing the message's level.
	KeyLevel = "level"

	// KeyTimestamp defines the message key containing the message's timestamp.
	KeyTimestamp = "ts"

	// KeyEvent defines the default message key containing the message's event.
	KeyEvent = "evt"
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

// NewMessage creates a new message from the given log level and key-value pairs
func NewMessage(l Level, pairs ...KVPair) Message {
	m := make(Message, len(pairs)+2)
	m[0] = KV(KeyTimestamp, time.Now())
	m[1] = KV(KeyLevel, l)

	copy(m[2:], pairs)

	return m
}
