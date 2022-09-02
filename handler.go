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
	"bytes"
	"io"
	"sync"
)

var (
	// Defines the size of an async handler's buffer that is preallocated.
	AsyncHandlerBufferSize = 2048

	// Defines the number of preallocated buffers in a pool of buffers.
	AsyncHandlerPoolSize = 64

	// Number of log events to buffer in an async handler's channel.
	AsyncHandlerChannelSize = 1024
)

type syncHandler struct {
	out       io.Writer
	formatter Formatter
}

// NewSyncHandler creates a new Handler that works synchronously by writing log events formatted with f to o.
// f works by directly writing bytes to o; no buffering is done inbetween those.
func NewSyncHandler(o io.Writer, f Formatter) Handler {
	return &syncHandler{
		out:       o,
		formatter: f,
	}
}

func (h *syncHandler) Close() {}

func (h *syncHandler) deliver(e *Event) {
	h.formatter.Format(h.out, e)
}

type asyncHandler struct {
	formatter    Formatter
	pool         *sync.Pool
	bufferChan   chan *bytes.Buffer
	finishedChan chan struct{}
}

// NewAsyncHandler creates a new Handler that works asynchronously. f is applied on every event writing to a
// bytes.Buffer. This happens in the same goroutine as emitting the log. The resulting bytes are them send
// over a channel to a dedicated goroutine which copies the bytes onto o.
func NewAsyncHandler(o io.Writer, f Formatter) Handler {
	bufferChan := make(chan *bytes.Buffer, AsyncHandlerChannelSize)
	finishedChan := make(chan struct{})

	pool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, AsyncHandlerBufferSize))
		},
	}

	for i := 0; i < AsyncHandlerPoolSize; i++ {
		pool.Put(bytes.NewBuffer(make([]byte, 0, AsyncHandlerBufferSize)))
	}

	go func() {
		defer close(finishedChan)

		for buf := range bufferChan {
			// TODO: Handle error
			o.Write(buf.Bytes())
			buf.Reset()
			pool.Put(buf)
		}
	}()

	return &asyncHandler{
		formatter:    f,
		bufferChan:   bufferChan,
		finishedChan: finishedChan,
		pool:         pool,
	}
}

func (h *asyncHandler) Close() {
	close(h.bufferChan)
	<-h.finishedChan
}

func (h *asyncHandler) deliver(e *Event) {
	buf := h.pool.Get().(*bytes.Buffer)
	h.formatter.Format(buf, e)
	h.bufferChan <- buf
}

type noopHandler struct{}

func (*noopHandler) Close()         {}
func (*noopHandler) deliver(*Event) {}

// NoOpHandler creates a no-operation handler that simply discards every event. Use this handler to silence
// logging output completely.
func NoOpHandler() Handler {
	return &noopHandler{}
}
