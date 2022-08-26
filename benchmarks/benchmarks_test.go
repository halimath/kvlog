package kvlogbenchmarks

import (
	"bytes"
	"io"
	"testing"

	"github.com/halimath/kvlog"
	"github.com/rs/zerolog"
)

var noopFormatter kvlog.Formatter = kvlog.FormatterFunc(func(io.Writer, *kvlog.Event) error { return nil })

func BenchmarkKVLog(b *testing.B) {
	var buf bytes.Buffer
	l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter()))
	// l = kvlog.New(kvlog.NewSyncHandler(&buf, noopFormatter))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.With().
			// KV("foo", 17).
			Log("some message")
		// l.With("a", "f").Log("1")
	}
}

func BenchmarkZerolog(b *testing.B) {
	var buf bytes.Buffer
	l := zerolog.New(&buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Log().
			Str("spam", "eggs").
			Int("foo", 17).
			Msg("some message")
	}

}
