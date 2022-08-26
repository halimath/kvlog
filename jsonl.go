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
	"fmt"
	"io"
	"time"

	"github.com/halimath/kvlog/internal/jsonencoder"
)

type jsonlFormatter struct {
	enc *jsonencoder.Encoder
}

// JSONLFormatter creates a new Formatter that formats JSON lines.
func JSONLFormatter() Formatter {
	return &jsonlFormatter{
		enc: jsonencoder.New(),
	}
}

func (f *jsonlFormatter) Format(w io.Writer, e *Event) error {
	f.enc.Reset()

	f.enc.StartObject()

	e.EachPair(func(p Pair) {
		f.enc.Key(p.Key)

		switch x := p.Value.(type) {
		case time.Time:
			f.enc.Str(x.Format(time.RFC3339))
		case time.Duration:
			f.enc.Str(fmt.Sprintf("%.3fs", x.Seconds()))
		case int:
			f.enc.Int(int64(x))
		case int8:
			f.enc.Int(int64(x))
		case int16:
			f.enc.Int(int64(x))
		case int32:
			f.enc.Int(int64(x))
		case int64:
			f.enc.Int(x)
		case uint:
			f.enc.Int(int64(x))
		case uint8:
			f.enc.Int(int64(x))
		case uint16:
			f.enc.Int(int64(x))
		case uint32:
			f.enc.Int(int64(x))
		case uint64:
			f.enc.Int(int64(x))
		case float32:
			f.enc.Float(float64(x))
		case float64:
			f.enc.Float(x)
		case string:
			f.enc.Str(x)
		default:
			f.enc.Str(fmt.Sprintf("%s", x))
		}
	})

	f.enc.EndObject()

	w.Write(f.enc.Bytes())
	w.Write([]byte{'\n'})

	return nil
}
