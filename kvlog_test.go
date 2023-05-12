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

package kvlog_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/halimath/kvlog"
)

func TestLogger_noTimeHook(t *testing.T) {
	var buf bytes.Buffer

	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter()))

	l.Logs("hello")
	l.Logf("hello, %s", "world")

	nl := l.Sub(kvlog.WithKV("tracing_id", "1234"))
	nl.Logs("got request")

	l.Logs("goodbye")

	exp := `{"msg":"hello"}
{"msg":"hello, world"}
{"tracing_id":"1234","msg":"got request"}
{"msg":"goodbye"}
`

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func TestLogger_pairs(t *testing.T) {
	var buf bytes.Buffer

	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter()))

	l.Logs("pairs", kvlog.WithPairs(kvlog.Pairs{
		"foo":  "bar",
		"spam": "eggs",
	})...)

	t.Log(buf.String())

	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	foo := got["foo"]
	if foo != "bar" {
		t.Errorf("expected foo == \"bar\" but got \"%s\"", foo)
	}

	spam := got["spam"]
	if spam != "eggs" {
		t.Errorf("expected spam == \"eggs\" but got \"%s\"", spam)
	}

	msg := got[kvlog.KeyMessage]
	if msg != "pairs" {
		t.Errorf("expected msg == \"pairs\" but got \"%s\"", msg)
	}
}

func TestLogger_withTimeHook(t *testing.T) {
	var buf bytes.Buffer

	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter())).
		AddHook(kvlog.TimeHook)
	now := time.Now()

	l.Logs("hello")
	l.Logf("hello, %s", "world")

	exp := fmt.Sprintf(`{"time":"%s","msg":"hello"}
{"time":"%s","msg":"hello, world"}
`, now.Format(time.RFC3339), now.Format(time.RFC3339))

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func TestLogger_concurrentTest(t *testing.T) {
	var buf bytes.Buffer

	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter()))

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				l.Logs("msg",
					kvlog.WithKV("i", i),
					kvlog.WithKV("j", j),
				)

				time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	if len(lines) != 100*100 {
		t.Fatalf("unexpected number of log lines: %d", len(lines))
	}

}
