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
	"io"
)

// Filter defines the interface for types that filter
// messages.
type Filter interface {
	// Filter filters the given message m and returns
	// either a message (which may be m) to be handled
	// or nil if the given message should be dropped.
	Filter(m Message) Message
}

// FilterFunc is a wrapper type implementing Filter
// that wraps a plain function.
type FilterFunc func(m Message) Message

// Filter just calls f to perform filtering.
func (f FilterFunc) Filter(m Message) Message {
	return f(m)
}

// Threshold is a factory for a Filter that
// drops messages if their level is less
// then the given threshold.
func Threshold(threshold Level) Filter {
	return FilterFunc(func(m Message) Message {
		if m.Level() >= threshold {
			return m
		}
		return nil
	})
}

// --

// Handler implements a threshold
type Handler struct {
	formatter Formatter
	output    Output
	filter    []Filter
}

// Deliver performs the delivery of the given message.
func (h *Handler) Deliver(m Message) {
	for _, f := range h.filter {
		m = f.Filter(m)
		if m == nil {
			return
		}
	}
	h.formatter.Format(m, h.output)
}

// Close closes the handler terminating its service.
// The underlying output is also closed.
func (h *Handler) Close() {
	c, ok := h.output.(io.Closer)
	if ok {
		c.Close()
	}
}

// NewHandler creates a new Handler using the provided values.
func NewHandler(formatter Formatter, output Output, filter ...Filter) *Handler {
	filterToUse := make([]Filter, len(filter))
	copy(filterToUse, filter)

	return &Handler{
		formatter: formatter,
		output:    output,
		filter:    filterToUse,
	}
}
