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

// Package msgutil provides helper functions when working with msg.Messages.
package msgutil

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/halimath/kvlog/msg"
)

type SortByKey []msg.KVPair

func (s SortByKey) Len() int      { return len(s) }
func (s SortByKey) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByKey) Less(i, j int) bool {
	if s[i].Key == msg.KeyTimestamp {
		return true
	}
	if s[j].Key == msg.KeyTimestamp {
		return false
	}
	if s[i].Key == msg.KeyLevel {
		return true
	}
	if s[j].Key == msg.KeyLevel {
		return false
	}
	if s[i].Key == msg.KeyEvent {
		return true
	}
	if s[j].Key == msg.KeyEvent {
		return false
	}

	return strings.Compare(s[i].Key, s[j].Key) < 0
}

func FormatValue(k msg.KVPair, w io.Writer) (err error) {
	switch x := k.Value.(type) {
	case time.Time:
		_, err = w.Write([]byte(x.Format(time.RFC3339)))
	case time.Duration:
		_, err = fmt.Fprintf(w, "%.3fs", x.Seconds())
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		_, err = fmt.Fprintf(w, "%d", x)
	case float32, float64:
		_, err = fmt.Fprintf(w, "%.3f", x)
	case msg.Level:
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
