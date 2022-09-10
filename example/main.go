package main

import (
	"fmt"
	"time"

	"github.com/halimath/kvlog"
)

func main() {
	kvlog.L.Logs("hello, world",
		kvlog.WithKV("tracing_id", 1234),
		kvlog.WithDur(time.Second),
		kvlog.WithErr(fmt.Errorf("some error")),
	)

	kvlog.L.Logs("some event",
		kvlog.WithPairs(kvlog.Pairs{
			"user_name": "foobar",
			"user_id":   "1234",
		})...)
}
