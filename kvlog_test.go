package kvlog

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestPackage(t *testing.T) {
	var buf bytes.Buffer
	now := time.Now()
	ConfigureOutput(&WriterLogOutput{
		w: &buf,
	})
	ConfigureThreshold(LevelWarn)

	Debug(KV("event", "test"), KV("foo", "bar"))
	Info(KV("event", "test"), KV("foo", "bar"))
	Warn(KV("event", "test"), KV("foo", "bar"))
	Error(KV("event", "test"), KV("foo", "bar"))

	exp := fmt.Sprintf("ts=%s level=warn event=<test> foo=<bar>\nts=%s level=error event=<test> foo=<bar>\n", now.Format("2006-01-02T15:04:05"), now.Format("2006-01-02T15:04:05"))

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func Example() {
	Debug(KV("event", "test"), KV("foo", "bar"))
	Info(KV("event", "test"), KV("foo", "bar"))

}
