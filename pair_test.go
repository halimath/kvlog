package kvlog

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestWriteTo(t *testing.T) {
	ts := time.Now()
	tab := map[KVPair]string{
		KV("foo", "bar"):         "foo=<bar>",
		KV("foo", 19):            "foo=19",
		KV("foo", 19.3):          "foo=19.300",
		KV("foo", LevelInfo):     "foo=info",
		KV("foo", "Hello world"): "foo=<Hello world>",
		KV("foo", ts):            fmt.Sprintf("foo=%s", ts.Format("2006-01-02T15:04:05")),
		KV("foo", 2*time.Second): "foo=2.000s",
	}

	for kv, exp := range tab {
		var buf bytes.Buffer
		kv.WriteTo(&buf)

		if buf.String() != exp {
			t.Errorf("failed to write %#v: expected '%s' but got '%s'", kv, exp, buf.String())
		}
	}
}
