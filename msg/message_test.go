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

package msg

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	table := map[Level]string{
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
	}

	for level, expected := range table {
		actual := level.String()
		if actual != expected {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	}
}

func TestMessageLevel(t *testing.T) {
	m := NewMessage(LevelInfo, KV("foo", "bar"))
	if m.Level() != LevelInfo {
		t.Errorf("expected info but got %s", m.Level())
	}
}
