//
// This file is part of kvlog.
//
// Copyright 2019, 2020 Alexander Metzner.
//
// Copyright 2019, 2020 Alexander Metzner.
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

// Package kvlog provides a key-value based logging system.
package kvlog

var l *Logger

func init() {
	Init(NewHandler(KVFormatter, Stdout(), Threshold(LevelInfo)))
}

// Init initializes the package global logger to a new logger
// using the given handler. The previous logger is closed if
// it had been set before.
func Init(handler ...*Handler) {
	if l != nil {
		l.Close()
	}
	l = NewLogger(handler...)
}

// Debug emits a log message of level debug.
func Debug(pairs ...KVPair) {
	l.Debug(pairs...)
}

// Info emits a log message of level info.
func Info(pairs ...KVPair) {
	l.Info(pairs...)
}

// Warn emits a log message of level warn.
func Warn(pairs ...KVPair) {
	l.Warn(pairs...)
}

// Error emits a log message of level error.
func Error(pairs ...KVPair) {
	l.Error(pairs...)
}
