package kvlog_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/halimath/kvlog"
)

func BenchmarkKVLog_syncHandler_JSONLFormatter(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	l := kvlog.New(kvlog.NewSyncHandler(out, kvlog.JSONLFormatter()))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.With().
			KV("spam", "eggs").
			KV("foo", 17).
			KV("enabled", true).
			Dur(time.Second).
			Log("some message")
	}
}

func BenchmarkKVLog_syncHandler_noopFormatter(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	l := kvlog.New(kvlog.NewSyncHandler(out, kvlog.FormatterFunc(func(io.Writer, *kvlog.Event) error { return nil })))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.With().
			KV("spam", "eggs").
			KV("foo", 17).
			KV("enabled", true).
			// Dur(time.Second).
			KV(kvlog.KeyDuration, time.Second).
			Log("some message")
	}
}
