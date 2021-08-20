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

package formatter

import (
	"io"

	"github.com/halimath/kvlog/msg"
)

// Interface defines the interface implemented by all
// message formatters.
type Interface interface {
	// Formats the given message into a slice of bytes.
	Format(m msg.Message, w io.Writer) error
}

// FormatterFunc is a converter type that allows using
// a plain function as a Formatter.
type FormatterFunc func(m msg.Message, w io.Writer) error

// Format simply calls ff.
func (ff FormatterFunc) Format(m msg.Message, w io.Writer) error {
	return ff(m, w)
}
