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
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestKVFormatter(t *testing.T) {
	now := time.Now()

	table := map[*Message]string{
		m(LevelInfo, KV("spam", "eggs"), KV("foo", "bar")): fmt.Sprintf("ts=%s level=info foo=bar spam=eggs\n", now.Format("2006-01-02T15:04:05")),
	}

	for msg, exp := range table {
		var buf bytes.Buffer
		if err := KVFormatter.Format(*msg, &buf); err != nil {
			t.Errorf("failed to format message: %s", err)
		} else if exp != buf.String() {
			t.Errorf("expected '%s' but got '%s'", exp, buf.String())
		}
	}
}

func TestFormatPair(t *testing.T) {
	ts := time.Now()
	tab := map[KVPair]string{
		KV("foo", "bar"):         "foo=bar",
		KV("foo", 19):            "foo=19",
		KV("foo", 19.3):          "foo=19.300",
		KV("foo", LevelInfo):     "foo=info",
		KV("foo", "Hello world"): "foo=<Hello world>",
		KV("foo", ts):            fmt.Sprintf("foo=%s", ts.Format("2006-01-02T15:04:05")),
		KV("foo", 2*time.Second): "foo=2.000s",
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

func TestTerminalFormatter(t *testing.T) {
	now := time.Now()

	table := map[*Message]string{
		m(LevelDebug, KV("spam", "eggs"), KV("foo", "bar")): fmt.Sprintf("\x1b[36m%s\x1b[0m \x1b[90mDEBUG\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format("2006-01-02T15:04:05")),
		m(LevelInfo, KV("spam", "eggs"), KV("foo", "bar")):  fmt.Sprintf("\x1b[36m%s\x1b[0m  \x1b[37mINFO\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format("2006-01-02T15:04:05")),
		m(LevelWarn, KV("spam", "eggs"), KV("foo", "bar")):  fmt.Sprintf("\x1b[36m%s\x1b[0m  \x1b[30;103mWARN\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format("2006-01-02T15:04:05")),
		m(LevelError, KV("spam", "eggs"), KV("foo", "bar")): fmt.Sprintf("\x1b[36m%s\x1b[0m \x1b[37;41mERROR\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format("2006-01-02T15:04:05")),
	}

	for msg, exp := range table {
		var buf bytes.Buffer
		if err := TerminalFormatter.Format(*msg, &buf); err != nil {
			t.Errorf("failed to format message: %s", err)
		} else if exp != buf.String() {
			t.Errorf("expected '%s' but got '%s'", exp, buf.String())
		}
	}
}

func m(l Level, pairs ...KVPair) *Message {
	m := NewMessage(l, pairs...)
	return &m
}
