// This file is part of kvlog.
//
// Copyright 2019, 2020, 2021 Alexander Metzner.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package msg contains types and functions that implement log messages.
package msg

import (
	"fmt"
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

// Event creates a KVPair for the "evt" key.
//
// Deprecated: Event is a deprecated alias for evt.
func Event(value interface{}) KVPair {
	return Evt(value)
}

// Evt creates a KVpair with the "evt" key.
func Evt(value interface{}) KVPair {
	return KVPair{
		Key:   KeyEvent,
		Value: value,
	}
}

// Err creates a KVPair with the "err" key.
func Err(err error) KVPair {
	return KVPair{
		Key:   KeyError,
		Value: err.Error(),
	}
}

// Msg creates a KVPair with the "msg" key. It formats the message using
// fmt.Sprintf.
func Msg(format string, args ...interface{}) KVPair {
	return KVPair{
		Key:   KeyMessage,
		Value: fmt.Sprintf(format, args...),
	}
}

// Dur creates a pair with the "dur" key containing an operation's duration.
func Dur(d time.Duration) KVPair {
	return KVPair{
		Key:   KeyDuration,
		Value: d,
	}
}

// Cat creates a pair with the "cat" key containing a message's category.
func Cat(cat string) KVPair {
	return KVPair{
		Key:   KeyCategory,
		Value: cat,
	}
}

// --

const (
	// KeyLevel defines the message key containing the message's level.
	KeyLevel = "lvl"

	// KeyTimestamp defines the message key containing the message's timestamp.
	KeyTimestamp = "ts"

	// KeyEvent defines the default message key containing the message's event.
	KeyEvent = "evt"

	// KeyError defines the default message key containing an error.
	KeyError = "err"

	// KeyMessage defines the default message key containing a textual msg.
	KeyMessage = "msg"

	// KeyDuration defines the default message key containing a duration measured in seconds.
	KeyDuration = "dur"

	// KeyCategory defines the default message key containing a message's category.
	KeyCategory = "cat"
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
