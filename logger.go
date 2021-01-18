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

package kvlog

import "sync"

// Logger implements a logger component.
type Logger struct {
	handler []chan Message
	wg      sync.WaitGroup
}

// NewLogger constructs a new Logger and returns a pointer to it.
func NewLogger(handler ...*Handler) *Logger {
	l := Logger{
		handler: make([]chan Message, 0, len(handler)),
	}

	for _, h := range handler {
		c := make(chan Message, 10)

		l.wg.Add(1)
		go func(c chan Message, h *Handler) {
			defer l.wg.Done()
			for m := range c {
				h.Deliver(m)
			}
			h.Close()
		}(c, h)

		l.handler = append(l.handler, c)
	}

	return &l
}

// Log logs the given message.
func (l *Logger) Log(m Message) {
	for _, c := range l.handler {
		c <- m
	}
}

// Close closes the handlers registered to this logger and waits for the goroutines
// to finish.
func (l *Logger) Close() {
	for _, c := range l.handler {
		close(c)
	}
	l.wg.Wait()
}

// Debug logs a message with level Debug.
func (l *Logger) Debug(pairs ...KVPair) {
	l.Log(NewMessage(LevelDebug, pairs...))
}

// Info logs a message with level Info.
func (l *Logger) Info(pairs ...KVPair) {
	l.Log(NewMessage(LevelInfo, pairs...))
}

// Warn logs a message with level Warn.
func (l *Logger) Warn(pairs ...KVPair) {
	l.Log(NewMessage(LevelWarn, pairs...))
}

// Error logs a message with level Error.
func (l *Logger) Error(pairs ...KVPair) {
	l.Log(NewMessage(LevelError, pairs...))
}
