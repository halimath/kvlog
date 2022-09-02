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
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// The default key used to identify an event's time stamp.
	KeyTime = "time"

	// The default key used to identify an event's error.
	KeyError = "err"

	// The default key used to identify an event's message.
	KeyMessage = "msg"

	// The default key used to identify an event's duration value.
	KeyDuration = "dur"

	// The default size Events created from an Event pool.
	DefaultEventSize = 16

	// Number of events to allocate for a new Event pool.
	InitialEventPoolSize = 128
)

// Pair defines a single key-value pair as part of a logging event.
type Pair struct {
	Key   string
	Value interface{}
}

// Event defines the type for a single logging event. key-value Pairs are given in descending priority order.
type Event struct {
	pairs []Pair
	len   int
	l     *logger
}

// Len returns the number of Pairs contained in e.
func (e *Event) Len() int {
	return e.len
}

// KV adds a new pair to e. The pair is constructed using key and value. If either key is an empty string
// or value is nil, no pair is added.
func (e *Event) KV(key string, value interface{}) *Event {
	if key == "" || value == nil {
		return e
	}

	if e.len < len(e.pairs) {
		e.pairs[e.len].Key = key
		e.pairs[e.len].Value = value
	} else {
		e.pairs = append(e.pairs, Pair{Key: key, Value: value})
	}

	e.len++

	return e
}

// Err adds a pair with KeyError and err.
func (e *Event) Err(err error) *Event {
	return e.KV(KeyError, err)
}

// Dur adds a pair with KeyDuration and d.
func (e *Event) Dur(d time.Duration) *Event {
	return e.KV(KeyDuration, d)
}

// Pairs defines a map of key-value-pairs to be added to an Event.
type Pairs map[string]interface{}

// Pairs adds all key-value pairs from p to e and returns e.
func (e *Event) Pairs(p Pairs) *Event {
	for k, v := range p {
		e.KV(k, v)
	}
	return e
}

// Log emits a log event.
func (e *Event) Log(v ...string) {
	if len(v) > 0 {
		e.KV(KeyMessage, strings.Join(v, ", "))
	}

	e.l.deliver(e)
}

// Logf emits a log event with a message produced by formatting args according to format.
func (e *Event) Logf(format string, args ...interface{}) {
	e.l.Log(fmt.Sprintf(format, args...))
}

// Logger create a nested logger utilizing the key-value pairs added to this builder which will be
// appended to all events produced from the resulting logger.
// Calling Logger exhausts e so e must not be used afterwards.
func (e *Event) Logger() Logger {
	return &logger{
		deliverFunc: e.l.deliverFunc,
		newEventFunc: func() *Event {
			n := e.l.newEventFunc()
			for i := 0; i < e.len; i++ {
				n.KV(e.pairs[i].Key, e.pairs[i].Value)
			}
			n.l = e.l
			return n
		},
	}
}

// EachPair applies f to each Pair in e in priority order (most important first).
func (e *Event) EachPair(f func(Pair)) {
	for i := e.len - 1; i >= 0; i-- {
		f(e.pairs[i])
	}
}

// A Hook is some kind of processing that gets applied to every log event going through some logger. A typical
// example is the time field that gets added via a hook from the root logger.
type Hook interface {
	// ApplyHook performs the hook's logic on e.
	ApplyHook(e *Event)
}

// HookFunc is a convenience used to implement hooks as a simple function.
type HookFunc func(e *Event)

func (h HookFunc) ApplyHook(e *Event) { h(e) }

// Logger defines the interface for all types that allow client code to emit log events.
type Logger interface {
	// AddHook adds a Hook to the Logger and returns it (or some derived Logger).
	AddHook(Hook) Logger

	// With starts a new Event used to configure either a log event or a sub-logger.
	With() *Event

	// Calling l.Log() is effectivly equivalent to calling
	//
	//   l.With().Log()
	//
	// It is provided for convenient logging of message-only events.
	Log(v ...string)

	// Calling l.Logf("hello %s", "foo") is effectivly equivalent to calling
	//
	//   l.With().Logf("hello %s", "foo")
	//
	// It is provided for convenient logging of message-only events.
	Logf(format string, args ...interface{})
}

// Formatter defines the interface implemented by all event formatters.
type Formatter interface {
	// Formats e on w.
	Format(w io.Writer, e *Event) error
}

// FormatterFunc is a converter type that allows using a plain function as a
// Formatter.
type FormatterFunc func(io.Writer, *Event) error

// Format simply calls ff.
func (ff FormatterFunc) Format(w io.Writer, e *Event) error {
	return ff(w, e)
}

// A Handler is used to deliver events to a given sink.
type Handler interface {
	Close()
	deliver(*Event)
	// TODO: Add filter
}

func newEvent() *Event {
	return &Event{
		pairs: make([]Pair, DefaultEventSize),
		len:   0,
	}
}

// New creates a new root Logger. It sends the events to all given handlers.
func New(handler ...Handler) Logger {
	eventPool := &sync.Pool{
		New: func() interface{} {
			return newEvent()
		},
	}

	for i := 0; i < InitialEventPoolSize; i++ {
		eventPool.Put(newEvent())
	}

	l := &logger{}

	l.newEventFunc = func() *Event {
		e := eventPool.Get().(*Event)
		e.len = 0
		e.l = l
		return e
	}

	l.deliverFunc = func(e *Event) {
		for _, h := range l.hooks {
			h.ApplyHook(e)
		}

		for _, h := range handler {
			h.deliver(e)
		}

		e.l = nil
		eventPool.Put(e)
	}

	return l
}

type logger struct {
	hooks        []Hook
	deliverFunc  func(e *Event)
	newEventFunc func() *Event
}

func (l *logger) AddHook(h Hook) Logger {
	l.hooks = append(l.hooks, h)

	return l
}

func (l *logger) With() *Event                            { return l.newEventFunc() }
func (l *logger) Log(v ...string)                         { l.With().Log(v...) }
func (l *logger) Logf(format string, args ...interface{}) { l.With().Logf(format, args...) }
func (l *logger) deliver(e *Event)                        { l.deliverFunc(e) }

// L is a default Logger which is initialized based on the environment the app is running in. When os.Stdout
// is a terminal device, the TerminalFormatter is used otherwise a JSONLFormatter is used. A TimeHook is
// applied to L.
var L Logger

func init() {
	var f Formatter
	if isTerminal() {
		f = ConsoleFormatter()
	} else {
		f = JSONLFormatter()
	}

	L = New(NewSyncHandler(os.Stdout, f)).AddHook(TimeHook)
}

func isTerminal() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
