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
		l.Logs("some message",
			kvlog.WithKV("spam", "eggs"),
			kvlog.WithKV("foo", 17),
			kvlog.WithKV("enabled", true),
			kvlog.WithDur(time.Second),
		)
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
		l.Logs("some message",
			kvlog.WithKV("foo", 17),
			kvlog.WithKV("enabled", true),
			// Dur(time.Second).
			kvlog.WithKV(kvlog.KeyDuration, time.Second),
		)
	}
}
