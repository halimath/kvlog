package main

import (
	"fmt"
	"time"

	"github.com/halimath/kvlog"
)

func main() {
	kvlog.L.With().
		KV("tracing_id", 1234).
		Dur(time.Second).
		Err(fmt.Errorf("some error")).
		Log("hello, world")

	kvlog.L.With().
		Pairs(kvlog.Pairs{
			"user_name": "foobar",
			"user_id":   "1234",
		}).
		Log("some event")
}
