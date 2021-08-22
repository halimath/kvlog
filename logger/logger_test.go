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

package logger

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/halimath/kvlog/filter"
	"github.com/halimath/kvlog/formatter/kvformat"
	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/msg"
)

func TestRoot(t *testing.T) {
	var output1, output2 bytes.Buffer

	l := New(
		handler.New(kvformat.Formatter, &output1, filter.Threshold(msg.LevelDebug)),
		handler.New(kvformat.Formatter, &output2, filter.Threshold(msg.LevelError)),
	)

	now := time.Now().Format(time.RFC3339)

	l.Debug(msg.KV("foo", "bar"))
	l.Info(msg.KV("foo", "bar"))
	l.Warn(msg.KV("foo", "bar"))
	l.Error(msg.KV("foo", "bar"))

	l.Close()

	exp1 := fmt.Sprintf(`ts=%[1]s lvl=debug foo=bar
ts=%[1]s lvl=info foo=bar
ts=%[1]s lvl=warn foo=bar
ts=%[1]s lvl=error foo=bar
`, now)

	if output1.String() != exp1 {
		t.Errorf("expected '%s' but got '%s'", exp1, output1.String())
	}

	exp2 := fmt.Sprintf("ts=%[1]s lvl=error foo=bar\n", now)
	if output2.String() != exp2 {
		t.Errorf("expected '%s' but got '%s'", exp2, output2.String())
	}
}

func TestNested(t *testing.T) {
	var output1 bytes.Buffer

	root := New(handler.New(kvformat.Formatter, &output1, filter.Threshold(msg.LevelDebug)))

	l := WithCategory(root, "test")

	now := time.Now().Format(time.RFC3339)

	l.Debug(msg.KV("foo", "bar"))
	l.Info(msg.KV("foo", "bar"))
	l.Warn(msg.KV("foo", "bar"))
	l.Error(msg.KV("foo", "bar"))

	root.Close()

	exp := fmt.Sprintf(`ts=%[1]s lvl=debug cat=test foo=bar
ts=%[1]s lvl=info cat=test foo=bar
ts=%[1]s lvl=warn cat=test foo=bar
ts=%[1]s lvl=error cat=test foo=bar
`, now)

	if output1.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, output1.String())
	}
}
