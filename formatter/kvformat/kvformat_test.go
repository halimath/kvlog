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

package kvformat

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/halimath/kvlog/msg"
)

func TestKVFormatter(t *testing.T) {
	now := time.Now()

	table := map[*msg.Message]string{
		m(msg.LevelInfo, msg.KV("spam", "eggs"), msg.KV("foo", "bar")): fmt.Sprintf("ts=%s lvl=info foo=bar spam=eggs\n", now.Format("2006-01-02T15:04:05")),
	}

	for msg, exp := range table {
		var buf bytes.Buffer
		if err := Formatter.Format(*msg, &buf); err != nil {
			t.Errorf("failed to format message: %s", err)
		} else if exp != buf.String() {
			t.Errorf("expected '%s' but got '%s'", exp, buf.String())
		}
	}
}

func TestFormatPair(t *testing.T) {
	ts := time.Now()
	tab := map[msg.KVPair]string{
		msg.KV("foo", "bar"):         "foo=bar",
		msg.KV("foo", 19):            "foo=19",
		msg.KV("foo", 19.3):          "foo=19.300",
		msg.KV("foo", msg.LevelInfo): "foo=info",
		msg.KV("foo", "Hello world"): "foo=<Hello world>",
		msg.KV("foo", ts):            fmt.Sprintf("foo=%s", ts.Format("2006-01-02T15:04:05")),
		msg.KV("foo", 2*time.Second): "foo=2.000s",
	}

	for kv, exp := range tab {
		var buf bytes.Buffer
		if err := formatPair(kv, &buf); err != nil {
			t.Errorf("failed to format %#v: %s", kv, err)
		} else if buf.String() != exp {
			t.Errorf("failed to write %#v: expected '%s' but got '%s'", kv, exp, buf.String())
		}
	}
}

func m(l msg.Level, pairs ...msg.KVPair) *msg.Message {
	m := msg.NewMessage(l, pairs...)
	return &m
}
