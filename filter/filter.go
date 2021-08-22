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

// Package filter provides types and functions that filter messages.
package filter

import "github.com/halimath/kvlog/msg"

// Interface defines the interface for types that filter
// messages.
type Interface interface {
	// Filter filters the given message m and returns
	// either a message (which may be m) to be handled
	// or nil if the given message should be dropped.
	Filter(m msg.Message) msg.Message
}

// FilterFunc is a wrapper type implementing Filter
// that wraps a plain function.
type FilterFunc func(m msg.Message) msg.Message

// Filter just calls f to perform filtering.
func (f FilterFunc) Filter(m msg.Message) msg.Message {
	return f(m)
}

// Threshold is a factory for a Filter that
// drops messages if their level is less
// then the given threshold.
func Threshold(threshold msg.Level) Interface {
	return FilterFunc(func(m msg.Message) msg.Message {
		if m.Level() >= threshold {
			return m
		}
		return nil
	})
}
