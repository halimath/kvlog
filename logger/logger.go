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

package logger

import (
	"io"
	"sync"

	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/msg"
)

// Interface defines the interface for loggers.
type Interface interface {
	io.Closer

	// Debug logs a message with level Debug.
	Debug(pairs ...msg.KVPair)
	// Info logs a message with level Info.
	Info(pairs ...msg.KVPair)
	// Warn logs a message with level Warn.
	Warn(pairs ...msg.KVPair)
	// Error logs a message with level Error.
	Error(pairs ...msg.KVPair)
}

// root implements a root logger.
type root struct {
	handlers []chan msg.Message
	wg       sync.WaitGroup
}

// New constructs a new root logger and returns it.
func New(handlers ...*handler.Handler) Interface {
	l := root{
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

// log logs the given msg.
func (r *root) log(m msg.Message) {
	for _, c := range r.handlers {
		c <- m
	}
}

// Close closes the handlers registered to this logger and waits for the goroutines
// to finish.
func (r *root) Close() error {
	for _, c := range r.handlers {
		close(c)
	}
	r.wg.Wait()

	return nil
}

// Debug logs a message with level Debug.
func (r *root) Debug(pairs ...msg.KVPair) {
	r.log(msg.NewMessage(msg.LevelDebug, pairs...))
}

// Info logs a message with level Info.
func (r *root) Info(pairs ...msg.KVPair) {
	r.log(msg.NewMessage(msg.LevelInfo, pairs...))
}

// Warn logs a message with level Warn.
func (r *root) Warn(pairs ...msg.KVPair) {
	r.log(msg.NewMessage(msg.LevelWarn, pairs...))
}

// Error logs a message with level Error.
func (r *root) Error(pairs ...msg.KVPair) {
	r.log(msg.NewMessage(msg.LevelError, pairs...))
}

// --

type nested struct {
	parent Interface
	pairs  []msg.KVPair
}

// With creates a new logger nested under l which hosts pairs
// and adds them to all messages emitted via the returned logger.
func With(l Interface, pairs ...msg.KVPair) Interface {
	return &nested{
		parent: l,
		pairs:  pairs,
	}
}

// WithCategory uses With to create a nested logger with a fixed category.
func WithCategory(l Interface, cat string) Interface {
	return With(l, msg.Cat(cat))
}

// Close does nothing.
func (n *nested) Close() error {
	return nil
}

// Debug logs a message with level Debug.
func (n *nested) Debug(pairs ...msg.KVPair) {
	p := make([]msg.KVPair, len(pairs)+len(n.pairs))
	copy(p, n.pairs)
	copy(p[len(n.pairs):], pairs)

	n.parent.Debug(p...)
}

// Info logs a message with level Info.
func (n *nested) Info(pairs ...msg.KVPair) {
	p := make([]msg.KVPair, len(pairs)+len(n.pairs))
	copy(p, n.pairs)
	copy(p[len(n.pairs):], pairs)

	n.parent.Info(p...)
}

// Warn logs a message with level Warn.
func (n *nested) Warn(pairs ...msg.KVPair) {
	p := make([]msg.KVPair, len(pairs)+len(n.pairs))
	copy(p, n.pairs)
	copy(p[len(n.pairs):], pairs)

	n.parent.Warn(p...)
}

// Error logs a message with level Error.
func (n *nested) Error(pairs ...msg.KVPair) {
	p := make([]msg.KVPair, len(pairs)+len(n.pairs))
	copy(p, n.pairs)
	copy(p[len(n.pairs):], pairs)

	n.parent.Error(p...)
}
