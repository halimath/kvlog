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

package filter

import (
	"testing"

	"github.com/halimath/kvlog/msg"
)

func TestThreshold(t *testing.T) {
	f := Threshold(msg.LevelWarn)

	tab := map[*msg.Message]bool{
		m(msg.LevelDebug): false,
		m(msg.LevelInfo):  false,
		m(msg.LevelWarn):  true,
		m(msg.LevelError): true,
	}

	for in, exp := range tab {
		act := f.Filter(*in) != nil

		if exp != act {
			t.Errorf("%s: expected %v but got %v", in.Level(), exp, act)
		}
	}
}

func m(l msg.Level) *msg.Message {
	m := msg.NewMessage(l)
	return &m
}
