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
	"github.com/halimath/kvlog/formatter/kvformat"
	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/msg"
)

func TestPackage(t *testing.T) {
	var buf bytes.Buffer

	kvlog.Init(handler.New(kvformat.Formatter, &buf, handler.Threshold(msg.LevelWarn)))

	now := time.Now()
	kvlog.Debug(kvlog.Evt("test"), kvlog.KV("foo", "bar"))
	kvlog.Info(kvlog.Evt("test"), kvlog.KV("foo", "bar"))
	kvlog.Warn(kvlog.Evt("test"), kvlog.KV("foo", "bar"))
	kvlog.Error(kvlog.Evt("test"), kvlog.KV("foo", "bar"))

	// Run Init again to close the old handler
	kvlog.Init(handler.New(kvformat.Formatter, &buf, handler.Threshold(msg.LevelWarn)))

	exp := fmt.Sprintf("ts=%s lvl=warn evt=test foo=bar\nts=%s lvl=error evt=test foo=bar\n", now.Format("2006-01-02T15:04:05"), now.Format("2006-01-02T15:04:05"))

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func Example() {
	kvlog.Debug(kvlog.Evt("test"), kvlog.KV("foo", "bar"))
	kvlog.Info(kvlog.Evt("test"), kvlog.KV("foo", "bar"))
}
