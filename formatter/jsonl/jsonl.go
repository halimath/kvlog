// Package jsonl provides a formatter to format messages in JSON lines format.
package jsonl

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/halimath/kvlog/formatter"
	"github.com/halimath/kvlog/internal/jsonl"
	"github.com/halimath/kvlog/internal/msgutil"
	"github.com/halimath/kvlog/msg"
)

type jsonlFormatter struct {
	w  io.Writer
	jw *jsonl.Writer
}

// New creates a new formatter.Interface that formats JSON lines.
func New() formatter.Interface {
	return &jsonlFormatter{}
}

func (f *jsonlFormatter) Format(m msg.Message, w io.Writer) error {
	if f.jw == nil || f.w != w {
		f.w = w
		f.jw = jsonl.New(w)
	}

	sorted := msgutil.SortByKey(m)
	sort.Sort(sorted)

	f.jw.StartObject()

	for _, p := range sorted {
		f.jw.Key(p.Key)

		switch x := p.Value.(type) {
		case time.Time:
			f.jw.String(x.Format(time.RFC3339))
		case time.Duration:
			f.jw.String(fmt.Sprintf("%.3fs", x.Seconds()))
		case int:
			f.jw.Int(int64(x))
		case int8:
			f.jw.Int(int64(x))
		case int16:
			f.jw.Int(int64(x))
		case int32:
			f.jw.Int(int64(x))
		case int64:
			f.jw.Int(x)
		case uint:
			f.jw.Int(int64(x))
		case uint8:
			f.jw.Int(int64(x))
		case uint16:
			f.jw.Int(int64(x))
		case uint32:
			f.jw.Int(int64(x))
		case uint64:
			f.jw.Int(int64(x))
		case float32:
			f.jw.Float(float64(x))
		case float64:
			f.jw.Float(x)
		case msg.Level:
			f.jw.String(x.String())
		case string:
			f.jw.String(x)
		default:
			f.jw.String(fmt.Sprintf("%s", x))
		}
	}

	f.jw.EndObject()

	return nil
}
