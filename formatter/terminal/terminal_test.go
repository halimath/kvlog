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

package terminal

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/halimath/kvlog/msg"
)

func TestTerminalFormatter(t *testing.T) {
	now := time.Now()

	table := map[*msg.Message]string{
		m(msg.LevelDebug, msg.KV("spam", "eggs"), msg.KV("foo", "bar")): fmt.Sprintf("\x1b[36m%s\x1b[0m \x1b[90mDEBUG\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format(time.RFC3339)),
		m(msg.LevelInfo, msg.KV("spam", "eggs"), msg.KV("foo", "bar")):  fmt.Sprintf("\x1b[36m%s\x1b[0m  \x1b[37mINFO\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format(time.RFC3339)),
		m(msg.LevelWarn, msg.KV("spam", "eggs"), msg.KV("foo", "bar")):  fmt.Sprintf("\x1b[36m%s\x1b[0m  \x1b[30;103mWARN\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format(time.RFC3339)),
		m(msg.LevelError, msg.KV("spam", "eggs"), msg.KV("foo", "bar")): fmt.Sprintf("\x1b[36m%s\x1b[0m \x1b[37;41mERROR\x1b[0m \x1b[90mfoo=\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam=\x1b[0m\x1b[97meggs\x1b[0m\n", now.Format(time.RFC3339)),
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

func m(l msg.Level, pairs ...msg.KVPair) *msg.Message {
	m := msg.NewMessage(l, pairs...)
	return &m
}
