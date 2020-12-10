package kvlog_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"bitbucket.org/halimath/kvlog"
)

func TestPackage(t *testing.T) {
	var buf bytes.Buffer
	now := time.Now()
	kvlog.ConfigureOutput(&kvlog.WriterLogOutput{
		W: &buf,
	})
	kvlog.ConfigureThreshold(kvlog.LevelWarn)

	kvlog.Debug(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Info(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Warn(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Error(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))

	exp := fmt.Sprintf("ts=%s level=warn event=test foo=bar\nts=%s level=error event=test foo=bar\n", now.Format("2006-01-02T15:04:05"), now.Format("2006-01-02T15:04:05"))

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func Example() {
	kvlog.Debug(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Info(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
}
