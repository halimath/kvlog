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

import (
	"sync"

	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/msg"
)

// Logger implements a logger component.
type Logger struct {
	handlers []chan msg.Message
	wg       sync.WaitGroup
}

// NewLogger constructs a new Logger and returns a pointer to it.
func NewLogger(handlers ...*handler.Handler) *Logger {
	l := Logger{
		handlers: make([]chan msg.Message, 0, len(handlers)),
	}

	for _, h := range handlers {
		c := make(chan msg.Message, 10)

		l.wg.Add(1)
		go func(c chan msg.Message, h *handler.Handler) {
			defer l.wg.Done()
			for m := range c {
				h.Deliver(m)
			}
			h.Close()
		}(c, h)

		l.handlers = append(l.handlers, c)
	}

	return &l
}

// Log logs the given msg.
func (l *Logger) Log(m msg.Message) {
	for _, c := range l.handlers {
		c <- m
	}
}

// Close closes the handlers registered to this logger and waits for the goroutines
// to finish.
func (l *Logger) Close() {
	for _, c := range l.handlers {
		close(c)
	}
	l.wg.Wait()
}

// Debug logs a message with level Debug.
func (l *Logger) Debug(pairs ...msg.KVPair) {
	l.Log(msg.NewMessage(msg.LevelDebug, pairs...))
}

// Info logs a message with level Info.
func (l *Logger) Info(pairs ...msg.KVPair) {
	l.Log(msg.NewMessage(msg.LevelInfo, pairs...))
}

// Warn logs a message with level Warn.
func (l *Logger) Warn(pairs ...msg.KVPair) {
	l.Log(msg.NewMessage(msg.LevelWarn, pairs...))
}

// Error logs a message with level Error.
func (l *Logger) Error(pairs ...msg.KVPair) {
	l.Log(msg.NewMessage(msg.LevelError, pairs...))
}
