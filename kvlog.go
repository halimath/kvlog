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

// Package kvlog provides a structured logging facility. It's structure is using key-value pairs.
package kvlog

import (
	"fmt"
	"os"
	"time"

	"github.com/halimath/kvlog/filter"
	"github.com/halimath/kvlog/formatter"
	"github.com/halimath/kvlog/formatter/kvformat"
	"github.com/halimath/kvlog/formatter/terminal"
	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/msg"
	"github.com/halimath/kvlog/output"
)

// L is the Logger instance used by package level functions.
// Use this logger as a convenience.
var L *Logger

func init() {
	var f formatter.Interface
	if isTerminal() {
		f = terminal.Formatter
	} else {
		f = kvformat.Formatter
	}
	Init(handler.New(f, output.Stdout(), filter.Threshold(msg.LevelInfo)))
}

// Init initializes the package global logger to a new logger
// using the given handler. The previous logger is closed if
// it had been set before.
func Init(handler ...*handler.Handler) {
	if L != nil {
		L.Close()
	}
	L = NewLogger(handler...)
}

// Debug emits a log message of level debug.
func Debug(pairs ...msg.KVPair) {
	L.Debug(pairs...)
}

// Info emits a log message of level info.
func Info(pairs ...msg.KVPair) {
	L.Info(pairs...)
}

// Warn emits a log message of level warn.
func Warn(pairs ...msg.KVPair) {
	L.Warn(pairs...)
}

// Error emits a log message of level error.
func Error(pairs ...msg.KVPair) {
	L.Error(pairs...)
}

// KV is a factory method for KVPair values.
func KV(key string, value interface{}) msg.KVPair {
	return msg.KVPair{
		Key:   key,
		Value: value,
	}
}

// Event creates a msg.KVPair for the "evt" key.
//
// Deprecated: Event is a deprecated alias for evt.
func Event(value interface{}) msg.KVPair {
	return Evt(value)
}

// Evt creates a KVpair with the "evt" key.
func Evt(value interface{}) msg.KVPair {
	return msg.KVPair{
		Key:   msg.KeyEvent,
		Value: value,
	}
}

// Err creates a msg.KVPair with the "err" key.
func Err(err error) msg.KVPair {
	return msg.KVPair{
		Key:   msg.KeyError,
		Value: err.Error(),
	}
}

// Msg creates a msg.KVPair with the "msg" key. It formats the message using
// fmt.Sprintf.
func Msg(format string, args ...interface{}) msg.KVPair {
	return msg.KVPair{
		Key:   msg.KeyMessage,
		Value: fmt.Sprintf(format, args...),
	}
}

// Dur creates a pair with the "dur" key containing an operation's duration.
func Dur(d time.Duration) msg.KVPair {
	return msg.KVPair{
		Key:   msg.KeyDuration,
		Value: d,
	}
}

func isTerminal() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
