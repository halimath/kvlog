//
// This file is part of kvlog.
//
// Copyright 2019, 2020 Alexander Metzner.
//
// Copyright 2019, 2020 Alexander Metzner.
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
	"strings"
	"time"
)

// Formatter defines the interface implemented by all
// message formatters.
type Formatter interface {
	// Formats the given message into a slice of bytes.
	Format(m Message, w io.Writer) error
}

// FormatterFunc is a converter type that allows using
// a plain function as a Formatter.
type FormatterFunc func(m Message, w io.Writer) error

// Format simply calls ff.
func (ff FormatterFunc) Format(m Message, w io.Writer) error {
	return ff(m, w)
}

// KVFormatter implements a Formatter that writes the default KV format.
var KVFormatter = FormatterFunc(formatMessage)

func formatMessage(m Message, w io.Writer) error {
	for i, p := range m {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		formatPair(p, w)
	}

	w.Write([]byte("\n"))

	return nil
}

func formatPair(k KVPair, w io.Writer) error {
	var err error
	if _, err := fmt.Fprintf(w, "%s=", k.Key); err != nil {
		return err
	}

	switch x := k.Value.(type) {
	case time.Time:
		_, err = w.Write([]byte(x.Format("2006-01-02T15:04:05")))
	case time.Duration:
		_, err = fmt.Fprintf(w, "%.3fs", float64(x)/float64(time.Second))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		_, err = fmt.Fprintf(w, "%d", x)
	case float32, float64:
		_, err = fmt.Fprintf(w, "%.3f", x)
	case Level:
		_, err = w.Write([]byte(x.String()))
	case string:
		err = formatStringValue(w, x)
	default:
		_, err = fmt.Fprintf(w, "<%v>", x)
	}

	return err
}

func formatStringValue(w io.Writer, val string) error {
	var err error
	if strings.ContainsAny(val, "<> =\t\n\r") {
		_, err = fmt.Fprintf(w, "<%s>", val)
	} else {
		_, err = w.Write([]byte(val))
	}

	return err
}
