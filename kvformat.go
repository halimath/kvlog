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
	"strings"
	"time"
)

// Formatter implements a formatter.Interface that writes the default KV format.
var KVFormatter = FormatterFunc(formatMessageAsKV)

func formatMessageAsKV(w io.Writer, e *Event) error {
	var pairWritten bool

	e.EachPair(func(p Pair) {
		if pairWritten {
			fmt.Fprint(w, " ")
		}
		formatPair(w, p)
		pairWritten = true
	})

	_, err := w.Write([]byte("\n"))
	return err
}

func formatPair(w io.Writer, p Pair) (err error) {
	if _, err := fmt.Fprintf(w, "%s=", p.Key); err != nil {
		return err
	}

	return formatValue(w, p)
}

func formatValue(w io.Writer, p Pair) (err error) {
	switch x := p.Value.(type) {
	case time.Time:
		_, err = w.Write([]byte(x.Format(time.RFC3339)))
	case time.Duration:
		_, err = fmt.Fprintf(w, "%.3fs", x.Seconds())
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		_, err = fmt.Fprintf(w, "%d", x)
	case float32, float64:
		_, err = fmt.Fprintf(w, "%.3f", x)
	case string:
		err = formatStringValue(w, x)
	default:
		_, err = fmt.Fprintf(w, "<%v>", x)
	}

	return
}

func formatStringValue(w io.Writer, val string) (err error) {
	if strings.ContainsAny(val, "<> =\t\n\r") {
		_, err = fmt.Fprintf(w, "<%s>", val)
	} else {
		_, err = w.Write([]byte(val))
	}

	return
}
