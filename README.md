# kvlog

![CI Status][ci-img-url] [![Go Report Card][go-report-card-img-url]][go-report-card-url]
[![Package Doc][package-doc-img-url]][package-doc-url] [![Releases][release-img-url]][release-url]

[ci-img-url]: https://github.com/halimath/kvlog/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/kvlog
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/kvlog
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/kvlog
[release-img-url]: https://img.shields.io/github/v/release/halimath/kvlog.svg
[release-url]: https://github.com/halimath/kvlog/releases

`kvlog` provides a structured logging facility. The underlying structure is based on key-value pairs.
key-value pairs are rendered as [JSON lines] but other Formatters can be used to provide different outputs
including custom ones.

[JSON lines]: https://jsonlines.org/

# Why another logging lib?

`kvlog` tries to find a balance between highly performance optimized libs such as `zap` or `zerolog` and those
that define a "dead simple" API (such as `go-kit` or `logrus`). It does not perform as well as the first two
but significantly better than the last two. `kvlog` tries to provide an easy-to-use API for producing
structured log events while keeping up a good performance.

# Usage

## Installation

`kvlog` uses go modules and requires Go 1.16 or greater.

```
$ go get -u github.com/halimath/kvlog
```
## Creating a Logger

To emit log events, you need a `Logger`. `kvlog` provides a ready-to-use `Logger` via the `L` variable.
You can create a custom logger giving you more flexibility on the logger's target and format. 
Creating a new `Logger` is done via the `kvlog.New` function. 
It accepts any number of `Handler`s. Each `Handler` pairs an `io.Writer` as well as a `Formatter`.

```go
logger := kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter())).
	AddHook(kvlog.TimeHook)
```

`Handler`s can be synchronous as well as asynchronous. 
Synchronous Handlers execute the Formatter as well as writing the output in the same goroutine that invoked
the `Logger`. 
Asynchronous `Handler`s dispatch the log event to a different goroutine via a channel. 
Thus, asynchronous Handlers must be closed before shutdown in order to flush the channel and emit all log 
events.

## Emitting Events

The easiest way to emit a simple log message is to use a `Logger`'s `Log`, `Log` or `Logf` method.

```go
kvlog.L.Logs("hello, world")
kvlog.L.Logf("hello, %s", "world)
```

`Log` will log all given key-value pairs while `Logs` and `Logf` will format a message with additional 
key-value pairs. With the default JSONL formatter, this produces:

```json
{"msg":"hello"}
{"msg":"hello, world"}
```

If you want to add more key-value pairs - which is the primary use case for a structured logger - you can pass
additional arguments to any of the log methods. key-value pairs are best created using one of the `With...`
functions from `kvlog`.

```go
kvlog.L.Logs("hello, world",
	kvlog.WithKV("tracing_id", 1234),
	kvlog.WithDur(time.Second),
	kvlog.WithErr(fmt.Errorf("some error")),
)
```

## Deriving Loggers

Logger's can be derived from another Logger. This enables to configure a set of key-value-pairs to be added
to every event emmitted via the deriverd logger. The syntax works similar to emitting log messages this time
only invoking the `Sub` method instead of `Log`.

```go
dl := l.Sub(
	kvlog.WithKV("tracing_id", "1234"),
)
```

## Hooks

In addition to deriving loggers, any number of `Hook`s may be added to a logger. The hook's callback function
is invoked everytime an `Event` is emitted via this logger or any of its derived loggers. Hooks are useful
to add dynamic values, such as timestamps or anything else read from the surrounding context. Adding a
timestamp to every log event is realized via the `TimeHook`.

```go
l := kvlog.New(kvlog.NewSyncHandler(&buf, kvlog.JSONLFormatter())).
	AddHook(kvlog.TimeHook)
```

You can write your own hook by implement the `kvlog.Hook` interface or using the `kvlog.HookFunc` convenience
type for a simple function. 

```go
// This is an example for some function that determines a dynamic value.
extractTracingID := func() string {
	// some real implementation here
	return "1234"
}

// Create a logger and add the hook
logger := kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter())).
	AddHook(kvlog.HookFunc(func(e *kvlog.Event) {
		e.AddPair(kvlog.WithKV("tracing_id", extractTracingID()))
	}))

// Emit some event
logger.Logs("request")
```

This example produces:

```json
{"tracing_id":"1234","msg":"request"}
```

## Formatters

The kvlog package comes with three Formatters out of the box:
- `JSONLFormatter` formats events as JSON line values
- `ConsoleFormatter` formats events for output on a terminal which includes colorizing the event
- `KVFormatter` formats events in the legacy KV-Format

The `JSONLFormatter` features a lot of optimizations to improve time and memory behavior. The other two have a
less optimized performance. While the `ConsoleFormatter` is intended for dev use the `KVFormatter` is only
provided for compatibility reasons and should be considered deprecated. Use `JSONLFormatter` for production
systems.

Custom formatters may be created by implementing the `kvlog.Formatter` interface or using the 
`kvlog.FormatterFunc` convenience type.

## HTTP Middleware

`kvlog` contains a HTTP middleware that generates an access log. It wraps another `http.Hander` allowing you 
to log only requests on those handlers you are interested in.

```go
package main

import (
	"net/http"

	"github.com/halimath/kvlog"
)

func main() {
    mux := http.NewServeMux()
    // ...
	kvlog.L.Log("started")
	http.ListenAndServe(":8000", kvlog.Middleware(kvlog.L, mux))
}
```

## Default Keys

The following table lists the default keys used by `kvlog`. You can customize these by setting a module-level
variable. These are also given in the table below.

Key | Used with | Variable to change | Description
-- | -- | -- | --
`time` | `TimeHook` | `KeyTime` | The default key used to identify an event's time stamp.
`err` | `Event.Err` | `KeyError` | The default key used to identify an event's error.
`msg` | `Event.Log` or `Event.Logf` | `KeyMessage` | The default key used to identify an event's message.
`dur` | `Event.Dur` | `KeyDuration` | The default key used to identify an event's duration value.

## Customizing memory behavior

`kvlog` includes a set of performance optimizations. Most of them work by pre-allocating memory for data
structures to be reused in order to avoid allocations for each log event. These pre-allocations can be tuned
in case an application has a specific load profile. This tuning can be performed by setting module-global
variables. 

There are two primary areas, where pre-allocation is used.

1. New `Event`s are not created every time a logger's `With` method is called. Instead, most of the time a
   pre-existing `Event` is pulled from a `sync.Pool` and put back after the event has been formatted. Such a
   pool exists for _each root_ logger - that is _every_ logger created with `kvlog.New`. The pool is 
   pre-filled with a number of events. All the events pre-filled into the pool also have a pre-allocated
   number of key-value-pair slots which are re-used by overwriting them whenever a `KV` method is called.
   Both numbers - initial pool size and pre-allocated number of pairs - can be changed.
1. When using an asynchronous handler, the handler's formatter is invoked synchronously. The output is written
   to a `bytes.Buffer`. This buffer comes from a `sync.Pool` and has a pre-allocated bytes slice. After the
   event has been formatted, the buffer is sent over a bufferd channel. The channel is consumed by another
   goroutine, which copies the buffer's bytes on the output writer. After that, the buffer is put back into
   the pool. Pool size, buffer size and channel buffer size can be customized.

Changes to these variables only take effect for loggers/handlers created after the variable have been 
assigned. Use at your own risk.

Variable | Default Value | Description
:-- | --: | --
`DefaultEventSize` | 16 | The default size Events created from an Event pool.
`InitialEventPoolSize` | 128 | Number of events to allocate for a new Event pool.
`AsyncHandlerBufferSize` | 2048 | Defines the size of an async handler's buffer that is preallocated.
`AsyncHandlerPoolSize` | 64 | Defines the number of preallocated buffers in a pool of buffers.
`AsyncHandlerChannelSize` | 1024 | Number of log events to buffer in an async handler's channel.

# Benchmarks

The [`benchmarks`](./benchmarks) directory contains some benchmarking tests which compare different logging
libs with a structured log event consisting of several key-value-pairs of different types. These are the
results run on a developers laptop.

Library | ns/op | B/op | allocs/op
-- | --: | --: | --:
kvlog (sync handler) | 1527 | 152 | 8
kvlog (async handler) | 1392 | 153 | 8
zerolog | 352.9 | 0 | 0
logrus | 4430 | 2192 |34
go-kit/log  | 2201 | 970 | 18


# Changelog

## 0.9.0

__:warning: breaking change:__ This version provides a new API which is _not compatible_ to the API exposed
before. This involves the way loggers and other components are configured (which normally only affects a small
portion of the using code) as well as the way log events are emitted. The interface emitting log messages
remains similar.

* New API
* Performance improvements

## 0.8.1 - retracted

* Fix: add `sync.Mutex` to lock `Handler`

## 0.8.0 - retracted

* added `NoOpHandler` to easily silence logging output

## 0.7.0 - retracted

__:warning: breaking change:__ This version provides a new API which is _not compatible_ to the API exposed
before. This involves the way loggers and other components are configured as well as how log events are
emitted.

* New chaining API to create messages
* Performance optimization

## 0.6.0

__:warning: breaking change:__ This version provides a new API which is in parts
_not compatible_ to the API exposed before. This basically involves the way loggers
and other components are configured (which normally only affects a small portion of
the using code). The interface to creating and emitting log messages stays the same.

* Nested loggers allow adding default key value pairs to add to all logger (i.e. for use with a _category_)
* Reorganization into several packages; root package `kvlog` acts as a facade
* Renamed some types (mostly interfaces) to better match the new package name (i.e. `handler.Interface` instead of `handler.Handler`)
* Added `jsonl` formatter which creates [JSON lines](https://jsonlines.org/) output

## 0.5.0
* `KVFormatter` sorts pairs based on key
* New `TerminalFormatter` providing colored output on terminals
* Moved to github

## 0.4.0
* Export package level logger instance `L`

## 0.3.0
__:warning: breaking changes:__ This version provides a new API which is _not compatible_ to the 
API exposed before.
* Introduction of new component structure (see description above)

## 0.2.0
* Improve log message rendering

## 0.1.0
* Initial release

# License

```
Copyright 2019, 2020 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
