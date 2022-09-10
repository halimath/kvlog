package kvlogbenchmarks

import (
	"os"
	"testing"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/halimath/kvlog"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
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

func BenchmarkKVLog_syncHandler_KVFormatter(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	l := kvlog.New(kvlog.NewSyncHandler(out, kvlog.KVFormatter))
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

func BenchmarkKVLog_asyncHandler(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	h := kvlog.NewAsyncHandler(out, kvlog.JSONLFormatter())
	l := kvlog.New(h)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Logs("some message",
			kvlog.WithKV("spam", "eggs"),
			kvlog.WithKV("foo", 17),
			kvlog.WithKV("enabled", true),
			kvlog.WithDur(time.Second),
		)
	}

	h.Close()
}

func BenchmarkZerolog(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	l := zerolog.New(out)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Log().
			Str("spam", "eggs").
			Int("foo", 17).
			Bool("enabled", true).
			Dur("duration", time.Second).
			Msg("some message")
	}
}

func BenchmarkLogrus(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(out)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logrus.WithFields(logrus.Fields{
			"spam":     "eggs",
			"foo":      17,
			"enabled":  true,
			"duration": time.Second,
		}).Info("some message")
	}
}

func BenchmarkGoKitLog(b *testing.B) {
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer out.Close()

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(out))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Log("spam", "eggs", "foo", 17, "enabled", true, "duration", time.Second, "msg", "some message")
	}
}
