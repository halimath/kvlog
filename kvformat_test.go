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
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestKVFormatter(t *testing.T) {
	table := map[*Event]string{
		newEvent().KV("spam", "eggs").KV("foo", "hello world"): "foo=<hello world> spam=eggs\n",
	}

	for evt, exp := range table {
		var buf bytes.Buffer
		if err := KVFormatter.Format(&buf, evt); err != nil {
			t.Errorf("failed to format message: %s", err)
		} else if exp != buf.String() {
			t.Errorf("expected '%s' but got '%s'", exp, buf.String())
		}
	}
}

func TestFormatPair(t *testing.T) {
	ts := time.Now()
	tab := map[Pair]string{
		{Key: "foo", Value: "bar"}:           "foo=bar",
		{Key: "foo", Value: 19}:              "foo=19",
		{Key: "foo", Value: 19.3}:            "foo=19.300",
		{Key: "foo", Value: "Hello world"}:   "foo=<Hello world>",
		{Key: "foo", Value: ts}:              fmt.Sprintf("foo=%s", ts.Format(time.RFC3339)),
		{Key: "foo", Value: 2 * time.Second}: "foo=2.000s",
	}

	for p, exp := range tab {
		var buf bytes.Buffer
		if err := formatPair(&buf, p); err != nil {
			t.Errorf("failed to format %#v: %s", p, err)
		} else if buf.String() != exp {
			t.Errorf("failed to write %#v: expected '%s' but got '%s'", p, exp, buf.String())
		}
	}
}
