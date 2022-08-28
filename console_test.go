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
	"testing"
)

func TestConsoleFormatter(t *testing.T) {
	table := map[*Event]string{
		newEvent().KV("spam", "eggs").KV("foo", "bar"): "\x1b[90mfoo:\x1b[0m\x1b[97mbar\x1b[0m \x1b[90mspam:\x1b[0m\x1b[97meggs\x1b[0m\n",
	}

	for evt, exp := range table {
		var buf bytes.Buffer
		if err := ConsoleFormatter().Format(&buf, evt); err != nil {
			t.Errorf("failed to format message: %s", err)
		} else if exp != buf.String() {
			t.Errorf("expected '%s' but got '%s'", exp, buf.String())
		}
	}
}
