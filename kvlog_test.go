package kvlog_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/halimath/kvlog"
)

func TestPackage(t *testing.T) {
	var buf bytes.Buffer

	kvlog.Init(kvlog.NewHandler(kvlog.KVFormatter, &buf, kvlog.Threshold(kvlog.LevelWarn)))

	now := time.Now()
	kvlog.Debug(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Info(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Warn(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Error(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))

	// Run Init again to close the old handler
	kvlog.Init(kvlog.NewHandler(kvlog.KVFormatter, &buf, kvlog.Threshold(kvlog.LevelWarn)))

	exp := fmt.Sprintf("ts=%s level=warn event=test foo=bar\nts=%s level=error event=test foo=bar\n", now.Format("2006-01-02T15:04:05"), now.Format("2006-01-02T15:04:05"))

	if buf.String() != exp {
		t.Errorf("expected '%s' but got '%s'", exp, buf.String())
	}
}

func Example() {
	kvlog.Debug(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
	kvlog.Info(kvlog.KV("event", "test"), kvlog.KV("foo", "bar"))
}
