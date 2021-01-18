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

func TestLogger(t *testing.T) {
	var output1, output2 bytes.Buffer

	l := NewLogger(
		NewHandler(KVFormatter, &output1, Threshold(LevelDebug)),
		NewHandler(KVFormatter, &output2, Threshold(LevelError)),
	)

	now := time.Now().Format("2006-01-02T15:04:05")

	l.Debug(KV("foo", "bar"))
	l.Info(KV("foo", "bar"))
	l.Warn(KV("foo", "bar"))
	l.Error(KV("foo", "bar"))

	l.Close()

	exp1 := fmt.Sprintf(`ts=%[1]s level=debug foo=bar
ts=%[1]s level=info foo=bar
ts=%[1]s level=warn foo=bar
ts=%[1]s level=error foo=bar
`, now)

	if output1.String() != exp1 {
		t.Errorf("expected '%s' but got '%s'", exp1, output1.String())
	}

	exp2 := fmt.Sprintf("ts=%[1]s level=error foo=bar\n", now)
	if output2.String() != exp2 {
		t.Errorf("expected '%s' but got '%s'", exp2, output2.String())
	}

}
