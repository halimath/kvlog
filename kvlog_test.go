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
	"fmt"
	"testing"
	"time"

	"github.com/halimath/kvlog"
)

func TestPackage_noTimeHook(t *testing.T) {
	var buf bytes.Buffer

	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter()))

	l.Log("hello")
	l.Logf("hello, %s", "world")

	nl := l.With().KV("tracing_id", "1234").Logger()
	fmt.Printf("%#v\n", nl)
	nl.Log("got request")

	exp := `{"msg":"hello"}
{"msg":"hello, world"}
{"msg":"got request","tracing_id":"1234"}
`

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func TestPackage_withTimeHook(t *testing.T) {
	var buf bytes.Buffer

	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter())).
		AddHook(kvlog.TimeHook)
	now := time.Now()

	l.Log("hello")
	l.Logf("hello, %s", "world")

	exp := fmt.Sprintf(`{"time":"%s","msg":"hello"}
{"time":"%s","msg":"hello, world"}
`, now.Format(time.RFC3339), now.Format(time.RFC3339))

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}
