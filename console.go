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
	"sort"
	"strings"
	"time"
)

// ConsoleFormatter is a formatter that outputs colorized log events to be
// used by developers sitting in front of a terminal.
func ConsoleFormatter() Formatter {
	start := time.Now()

	return FormatterFunc(func(w io.Writer, e *Event) (err error) {
		pairs := sorted(collectPairs(e))

		sort.Sort(pairs)

		for i, p := range pairs {
			if i > 0 {
				fmt.Fprint(w, " ")
			}

			if p.Key == KeyTime {
				t := p.Value.(time.Time).Sub(start)
				_, err = fmt.Fprintf(w, "\x1b[90m%s:%v\x1b[0m", p.Key, t)
				if err != nil {
					return
				}

			} else {
				var valueColor string
				switch p.Value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					valueColor = "1;34"
				case time.Duration:
					valueColor = "1;36"
				case error:
					valueColor = "1;31"
				default:
					valueColor = "97"
				}

				_, err = fmt.Fprintf(w, "\x1b[90m%s:\x1b[0m\x1b[%sm%v\x1b[0m", p.Key, valueColor, p.Value)
				if err != nil {
					return
				}
			}
		}

		_, err = w.Write([]byte("\n"))

		return
	})
}

func collectPairs(e *Event) []Pair {
	pairs := make([]Pair, 0, e.Len())
	e.EachPair(func(p Pair) {
		pairs = append(pairs, p)
	})
	return pairs
}

type sorted []Pair

func (s sorted) Len() int      { return len(s) }
func (s sorted) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sorted) Less(i, j int) bool {
	if s[i].Key == KeyTime {
		return true
	}
	if s[j].Key == KeyTime {
		return false
	}

	if s[i].Key == KeyMessage {
		return true
	}
	if s[j].Key == KeyMessage {
		return false
	}

	if s[i].Key == KeyError {
		return true
	}
	if s[j].Key == KeyError {
		return false
	}

	if s[i].Key == KeyDuration {
		return true
	}
	if s[j].Key == KeyDuration {
		return false
	}

	return strings.Compare(s[i].Key, s[j].Key) < 0
}

func formatTimestamp(w io.Writer, p Pair) (err error) {
	_, err = w.Write([]byte("\x1b[36m"))
	if err != nil {
		return
	}
	_, err = fmt.Fprint(w, p.Value)
	if err != nil {
		return
	}
	_, err = w.Write([]byte("\x1b[0m"))
	return
}
