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

var (
	pairPool sync.Pool
)

const (
	pairPoolInitialSize = 128
)

func init() {
	pairPool = sync.Pool{
		New: func() interface{} {
			return &Pair{}
		},
	}

	for i := 0; i < pairPoolInitialSize; i++ {
		pairPool.Put(&Pair{})
	}
}

// WithKV creates a new Pair to be added to either an Event or a Logger. Pairs are pulled from a pool and are
// returned once the Event has been created.
func WithKV(key string, value interface{}) *Pair {
	p := pairPool.Get().(*Pair)
	p.Key = key
	p.Value = value
	return p
}

// WithErr creates a Pair with KeyError and err.
func WithErr(err error) *Pair {
	return WithKV(KeyError, err)
}

// WithDur creates a Pair with KeyDuration and d.
func WithDur(d time.Duration) *Pair {
	return WithKV(KeyDuration, d)
}

// Pairs defines a map of key-value-pairs to be added to an Event.
type Pairs map[string]interface{}

// WithPairs adds all key-value pairs from p to e and returns e.
func WithPairs(p Pairs) []*Pair {
	pairs := make([]*Pair, 0, len(p))
	for k, v := range p {
		pairs = append(pairs, WithKV(k, v))
	}
	return pairs
}

// Event defines the type for a single logging event. key-value Pairs are given in descending priority order.
type Event struct {
	pairs []Pair
	len   int
}

// Len returns the number of Pairs contained in e.
func (e *Event) Len() int {
	return e.len
}

// EachPair applies f to each Pair in e in priority order (most important first).
func (e *Event) EachPair(f func(Pair)) {
	for i := e.len - 1; i >= 0; i-- {
		f(e.pairs[i])
	}
}

func (e *Event) AddPair(p *Pair) {
	if e.len < len(e.pairs) {
		e.pairs[e.len].Key = p.Key
		e.pairs[e.len].Value = p.Value
	} else {
		e.pairs = append(e.pairs, Pair{Key: p.Key, Value: p.Value})
	}

	e.len++
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

	// Log logs an Event consisting of pairs merged with this logger's pairs.
	Log(pairs ...*Pair)

	// Logs logs an Event consisting of paris merged with this logger's pairs. In addition, a message pair
	// is added with msg being the value.
	Logs(msg string, pairs ...*Pair)

	// Logf works similar to Logs. It formats the message according to format. args is filtered before being
	// passed to fmt.Sprintf; all Pair values are removed and are passed to Logs separately.
	//
	// Example
	//
	//   l.Logf("hello, %s", "world", kvlog.WithKV("foo", "bar"))
	//
	// is identical to
	//
	//   l.Logs("hello, world", kvlog.WithKV("foo", "bar"))
	//
	Logf(format string, args ...interface{})

	// Sub creates a sub-logger using pairs for every event.
	Sub(pairs ...*Pair) Logger
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
		return e
	}

	l.deliverFunc = func(e *Event) {
		for _, h := range l.hooks {
			h.ApplyHook(e)
		}

		for _, h := range handler {
			h.deliver(e)
		}

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

func (l *logger) Log(pairs ...*Pair) {
	evt := l.newEventFunc()
	for _, p := range pairs {
		evt.AddPair(p)
		pairPool.Put(p)
	}
	l.deliverFunc(evt)
}

func (l *logger) Logs(msg string, pairs ...*Pair) {
	if len(pairs) == 0 {
		l.Log(WithKV(KeyMessage, msg))
		return
	}

	pairs = append(pairs, WithKV(KeyMessage, msg))
	l.Log(pairs...)
}

func (l *logger) Logf(format string, args ...interface{}) {
	formatArgs := make([]interface{}, 0, len(args))
	pairs := make([]*Pair, 0, len(args)+1)

	for _, arg := range args {
		if p, ok := arg.(*Pair); ok {
			pairs = append(pairs, p)
		} else {
			formatArgs = append(formatArgs, arg)
		}
	}

	pairs = append(pairs, WithKV(KeyMessage, fmt.Sprintf(format, formatArgs...)))

	l.Log(pairs...)
}

func (l *logger) Sub(pairs ...*Pair) Logger {
	h := HookFunc(func(e *Event) {
		for _, p := range pairs {
			e.AddPair(p)
		}
	})

	sub := &logger{
		hooks:        []Hook{h},
		newEventFunc: l.newEventFunc,
	}

	sub.deliverFunc = func(e *Event) {
		for _, h := range sub.hooks {
			h.ApplyHook(e)
		}
		l.deliverFunc(e)
	}

	return sub
}

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

type noOpLogger struct{}

func (l *noOpLogger) AddHook(Hook) Logger {
	return l
}
func (*noOpLogger) Log(pairs ...*Pair)                      {}
func (*noOpLogger) Logs(msg string, pairs ...*Pair)         {}
func (*noOpLogger) Logf(format string, args ...interface{}) {}
func (l *noOpLogger) Sub(pairs ...*Pair) Logger {
	return l
}

var noOpLoggerValue = &noOpLogger{}

// NoOPLogger creates a new no operation logger that discards all log messages.
func NoOpLogger() Logger {
	return noOpLoggerValue
}
