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

// L is the Logger instance used by package level functions.
// Use this logger as a convenience.
var L *Logger

func init() {
	var f Formatter
	if isTerminal() {
		f = TerminalFormatter
	} else {
		f = KVFormatter
	}
	Init(NewHandler(f, Stdout(), Threshold(LevelInfo)))
}

// Init initializes the package global logger to a new logger
// using the given handler. The previous logger is closed if
// it had been set before.
func Init(handler ...*Handler) {
	if L != nil {
		L.Close()
	}
	L = NewLogger(handler...)
}

// Debug emits a log message of level debug.
func Debug(pairs ...KVPair) {
	L.Debug(pairs...)
}

// Info emits a log message of level info.
func Info(pairs ...KVPair) {
	L.Info(pairs...)
}

// Warn emits a log message of level warn.
func Warn(pairs ...KVPair) {
	L.Warn(pairs...)
}

// Error emits a log message of level error.
func Error(pairs ...KVPair) {
	L.Error(pairs...)
}
