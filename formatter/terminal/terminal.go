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

// Package terminal provides a formatter.Interface writing messages in a colored terminal format.
package terminal

import (
	"fmt"
	"io"
	"sort"

	"github.com/halimath/kvlog/formatter"
	"github.com/halimath/kvlog/internal/msgutil"
	"github.com/halimath/kvlog/msg"
)

// Formatter is a formatter that outputs colored log messages
// to be used on a terminal.
var Formatter = formatter.FormatterFunc(formatMessageForTerminal)

func formatMessageForTerminal(m msg.Message, w io.Writer) (err error) {
	sorted := msgutil.SortByKey(m)
	sort.Sort(sorted)

	for i, p := range sorted {
		if i > 0 {
			fmt.Fprint(w, " ")
		}

		if p.Key == msg.KeyTimestamp {
			err = formatTimestamp(p, w)
			if err != nil {
				return
			}

		} else if p.Key == msg.KeyLevel {
			err = formatLevel(p, w)
			if err != nil {
				return
			}

		} else {
			_, err = fmt.Fprintf(w, "\x1b[90m%s=\x1b[0m\x1b[97m", p.Key)
			if err != nil {
				return
			}
			err = msgutil.FormatValue(p, w)
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

func formatLevel(p msg.KVPair, w io.Writer) (err error) {
	switch p.Value {
	case msg.LevelDebug:
		_, err = w.Write([]byte("\x1b[90mDEBUG\x1b[0m"))
	case msg.LevelInfo:
		_, err = w.Write([]byte(" \x1b[37mINFO\x1b[0m"))
	case msg.LevelWarn:
		_, err = w.Write([]byte(" \x1b[30;103mWARN\x1b[0m"))
	case msg.LevelError:
		_, err = w.Write([]byte("\x1b[37;41mERROR\x1b[0m"))
	default:
		panic(fmt.Sprintf("unexpected log level: %#v", p.Value))
	}
	return
}

func formatTimestamp(p msg.KVPair, w io.Writer) (err error) {
	_, err = w.Write([]byte("\x1b[36m"))
	if err != nil {
		return
	}
	err = msgutil.FormatValue(p, w)
	if err != nil {
		return
	}
	_, err = w.Write([]byte("\x1b[0m"))
	return
}
