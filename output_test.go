package kvlog

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestWriterLogOutput_WriteMessage(t *testing.T) {
	var buf bytes.Buffer
	o := &WriterLogOutput{
		w: &buf,
	}

	now := time.Now()
	m := NewMessage(LevelInfo, KV("foo", "bar"), KV("spam", "eggs"))
	o.WriteLogMessage(m)

	exp := fmt.Sprintf("ts=%s level=info foo=<bar> spam=<eggs>\n", now.Format("2006-01-02T15:04:05"))
	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}
