package kvlog_test

import (
	"fmt"
	"os"
	"time"

	"github.com/halimath/kvlog"
)

func ExampleL() {
	kvlog.L.Logs("hello, world",
		kvlog.WithKV("tracing_id", 1234),
		kvlog.WithDur(time.Second),
		kvlog.WithErr(fmt.Errorf("some error")),
	)
}

func Example_customLogger() {
	logger := kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter())).
		AddHook(kvlog.TimeHook)

	logger.Logs("hello, world",
		kvlog.WithKV("tracing_id", 1234),
		kvlog.WithDur(time.Second),
		kvlog.WithErr(fmt.Errorf("some error")),
	)

}

func ExampleNewAsyncHandler() {
	h := kvlog.NewAsyncHandler(os.Stdout, kvlog.JSONLFormatter())
	logger := kvlog.New(h)
	logger.Logs("test")
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
			e.AddPair(kvlog.WithKV("tracing_id", extractTracingID()))
		}))

	logger.Logs("request")

	// Output: {"tracing_id":"1234","msg":"request"}
}
