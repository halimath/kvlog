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
	"os"
	"sort"
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

// --

// KVFormatter implements a Formatter that writes the default KV format.
var KVFormatter = FormatterFunc(formatMessageAsKV)

type byKey []KVPair

func (s byKey) Len() int      { return len(s) }
func (s byKey) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byKey) Less(i, j int) bool {
	if s[i].Key == KeyTimestamp {
		return true
	}
	if s[j].Key == KeyTimestamp {
		return false
	}
	if s[i].Key == KeyLevel {
		return true
	}
	if s[j].Key == KeyLevel {
		return false
	}
	if s[i].Key == KeyEvent {
		return true
	}
	if s[j].Key == KeyEvent {
		return false
	}

	return strings.Compare(s[i].Key, s[j].Key) < 0
}

func formatMessageAsKV(m Message, w io.Writer) error {
	sorted := byKey(m)
	sort.Sort(sorted)

	for i, p := range sorted {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		formatPair(p, w)
	}

	w.Write([]byte("\n"))

	return nil
}

func formatPair(k KVPair, w io.Writer) (err error) {
	if _, err := fmt.Fprintf(w, "%s=", k.Key); err != nil {
		return err
	}

	err = formatValue(k, w)

	return
}

func formatValue(k KVPair, w io.Writer) (err error) {
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

// --

// TerminalFormatter is a formatter that outputs colored log messages
// to be used on a terminal.
var TerminalFormatter = FormatterFunc(formatMessageForTerminal)

func formatMessageForTerminal(m Message, w io.Writer) (err error) {
	sorted := byKey(m)
	sort.Sort(sorted)

	for i, p := range sorted {
		if i > 0 {
			fmt.Fprint(w, " ")
		}

		if p.Key == KeyTimestamp {
			_, err = w.Write([]byte("\x1b[36m"))
			if err != nil {
				return
			}
			err = formatValue(p, w)
			if err != nil {
				return
			}
			_, err = w.Write([]byte("\x1b[0m"))
			if err != nil {
				return
			}

		} else if p.Key == KeyLevel {
			switch p.Value {
			case LevelDebug:
				_, err = w.Write([]byte("\x1b[90mDEBUG\x1b[0m"))
			case LevelInfo:
				_, err = w.Write([]byte(" \x1b[37mINFO\x1b[0m"))
			case LevelWarn:
				_, err = w.Write([]byte(" \x1b[30;103mWARN\x1b[0m"))
			case LevelError:
				_, err = w.Write([]byte("\x1b[37;41mERROR\x1b[0m"))
			default:
				panic(fmt.Sprintf("unexpected log level: %#v", p.Value))
			}

			if err != nil {
				return
			}

		} else {
			_, err = fmt.Fprintf(w, "\x1b[90m%s=\x1b[0m\x1b[97m", p.Key)
			if err != nil {
				return
			}
			err = formatValue(p, w)
			if err != nil {
				return
			}
			_, err = w.Write([]byte("\x1b[0m"))
			if err != nil {
				return
			}
		}
	}

	_, err = w.Write([]byte("\n"))

	return
}

func isTerminal() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
