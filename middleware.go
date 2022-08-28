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
	"net/http"
	"time"
)

type accessLogMW struct {
	logger   Logger
	delegate http.Handler
}

type responseWriterWrapper struct {
	w          http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	return w.w.Write(data)
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (m *accessLogMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	wrapper := &responseWriterWrapper{
		w:          w,
		statusCode: 200,
	}

	m.delegate.ServeHTTP(wrapper, r)

	requestTime := time.Since(startTime)
	m.logger.With().
		KV("method", r.Method).
		KV("url", r.URL).
		KV("status", wrapper.statusCode).
		Dur(requestTime).
		Log("request")
}

// Middleware returns a http.Handler that acts as an access log middleware.
func Middleware(l Logger, h http.Handler) http.Handler {
	return &accessLogMW{
		logger:   l,
		delegate: h,
	}
}
