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
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	w            http.ResponseWriter
	statusCode   int
	bytesWritten uint64
}

var (
	_ http.ResponseWriter = &responseWriterWrapper{}
	_ http.Flusher        = &responseWriterWrapper{}
	_ http.Hijacker       = &responseWriterWrapper{}
	_ http.Pusher         = &responseWriterWrapper{}
)

func (w *responseWriterWrapper) Header() http.Header {
	return w.w.Header()
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		// Default to 200 if WriteHeader hasn't been called yet
		rw.statusCode = http.StatusOK
	}
	n, err := rw.w.Write(b)
	rw.bytesWritten += uint64(n)
	return n, err
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

// satisfy http.Flusher
func (rw *responseWriterWrapper) Flush() {
	if f, ok := rw.w.(http.Flusher); ok {
		f.Flush()
	}
}

// satisfy http.Hijacker
func (rw *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.w.(http.Hijacker); ok {
		return h.Hijack()
	}

	return nil, nil, fmt.Errorf("http.Hijacker not supported")
}

// satisfy http.Pusher
func (rw *responseWriterWrapper) Push(target string, opts *http.PushOptions) error {
	if p, ok := rw.w.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Middleware returns a middleware function that enables logging on l.
// If addToContext is true, l will be added to every request's context.
func Middleware(l Logger, addToContext bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			wrapper := &responseWriterWrapper{
				w:          w,
				statusCode: 200,
			}

			l := l.Sub(
				WithKV("method", r.Method),
				WithKV("url", r.URL),
			)

			if addToContext {
				r = r.WithContext(ContextWithLogger(r.Context(), l))
			}

			h.ServeHTTP(wrapper, r)

			requestTime := time.Since(startTime)
			l.Logs("request",
				WithKV("status", wrapper.statusCode),
				WithDur(requestTime),
			)
		})
	}
}
