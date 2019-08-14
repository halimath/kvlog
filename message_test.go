package kvlog

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestMessageLevel(t *testing.T) {
	m := NewMessage(LevelInfo, KV("foo", "bar"))
	if m.Level() != LevelInfo {
		t.Errorf("expected info but got %s", m.Level())
	}
}

func TestMessageWriteTo(t *testing.T) {
	now := time.Now()
	m := NewMessage(LevelInfo, KV("foo", "bar"), KV("spam", "eggs"))

	var buf bytes.Buffer
	m.WriteTo(&buf)

	exp := fmt.Sprintf("ts=%s level=info foo=<bar> spam=<eggs>", now.Format("2006-01-02T15:04:05"))
	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}
