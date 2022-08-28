package kvlog_test

import (
	"fmt"
	"os"
	"time"

	"github.com/halimath/kvlog"
)

func ExampleL() {
	kvlog.L.With().
		KV("tracing_id", 1234).
		Dur(time.Second).
		Err(fmt.Errorf("some error")).
		Log("hello, world")
}

func Example_customLogger() {
	logger := kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter())).
		AddHook(kvlog.TimeHook)

	logger.With().
		KV("tracing_id", 1234).
		Dur(time.Second).
		Err(fmt.Errorf("some error")).
		Log("hello, world")
}

func ExampleNewAsyncHandler() {
	h := kvlog.NewAsyncHandler(os.Stdout, kvlog.JSONLFormatter())
	logger := kvlog.New(h)
	logger.Log("test")
	h.Close()

	// Output: {"msg":"test"}

}

func ExampleCustomHook() {
	extractTracingID := func() string {
		// some real implementation here
		return "1234"
	}

	logger := kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter())).
		AddHook(kvlog.HookFunc(func(e *kvlog.Event) {
			e.KV("tracing_id", extractTracingID())
		}))

	logger.Log("request")

	// Output: {"tracing_id":"1234","msg":"request"}
}
