package jsonl

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/halimath/kvlog/msg"
)

func TestJSONLFormatter(t *testing.T) {
	now := time.Now()

	table := map[*msg.Message]string{
		m(msg.LevelInfo, msg.KV("spam", "eggs"), msg.KV("foo", "bar")): fmt.Sprintf(`{"ts":"%s","lvl":"info","foo":"bar","spam":"eggs"}`, now.Format(time.RFC3339)),
		m(msg.LevelInfo, msg.Dur(2*time.Second)):                       fmt.Sprintf(`{"ts":"%s","lvl":"info","dur":"2.000s"}`, now.Format(time.RFC3339)),
		m(msg.LevelInfo, msg.KV("foo", 17)):                            fmt.Sprintf(`{"ts":"%s","lvl":"info","foo":17}`, now.Format(time.RFC3339)),
		m(msg.LevelInfo, msg.KV("foo", 17.2)):                          fmt.Sprintf(`{"ts":"%s","lvl":"info","foo":1.72000000e+01}`, now.Format(time.RFC3339)),
	}

	f := New()

	for msg, exp := range table {
		var buf bytes.Buffer
		if err := f.Format(*msg, &buf); err != nil {
			t.Errorf("failed to format message: %s", err)
		} else if exp != strings.TrimSpace(buf.String()) {
			t.Errorf("expected '%s' but got '%s'", exp, strings.TrimSpace(buf.String()))
		}
	}
}

func m(l msg.Level, pairs ...msg.KVPair) *msg.Message {
	m := msg.NewMessage(l, pairs...)
	return &m
}
