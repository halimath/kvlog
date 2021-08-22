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

package handler

import (
	"io"

	"github.com/halimath/kvlog/filter"
	"github.com/halimath/kvlog/formatter"
	"github.com/halimath/kvlog/msg"
	"github.com/halimath/kvlog/output"
)

// Handler implements a threshold
type Handler struct {
	formatter formatter.Interface
	output    output.Output
	filter    []filter.Interface
}

// Deliver performs the delivery of the given msg.
func (h *Handler) Deliver(m msg.Message) {
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

// New creates a new Handler using the provided values.
func New(formatter formatter.Interface, output output.Output, filters ...filter.Interface) *Handler {
	filterToUse := make([]filter.Interface, len(filters))
	copy(filterToUse, filters)

	return &Handler{
		formatter: formatter,
		output:    output,
		filter:    filterToUse,
	}
}
