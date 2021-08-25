# kvlog

![CI Status][ci-img-url] [![Go Report Card][go-report-card-img-url]][go-report-card-url] [![Package Doc][package-doc-img-url]][package-doc-url] [![Releases][release-img-url]][release-url]

`kvlog` is a library which provides a structured logging facility for the go programming language(golang). 
`kvlog`'s structure is based on key-value pairs.

## Description

`kvlog` provides types and functions to produce a log stream of log events. Each event consists
of key-value pairs which may be encoded in different ways. The library provides multiple formatters
and allows for custom formatters to be used.

Structured log messages differ from conventional string-based log messages. They do not
contain a bare string message but its information is structured in a way which allows other
systems to interpret and use the data. `kvlog` uses key-value-pairs as its underlying structure.

### Components

`kvlog` is built from a set of _components_ that interact to implement logging functionality.

A `Message` is produced by the client. Every `Message` consists of `Pair`s each representing
a single key-value-pair. `kvlog` provides a convenient and idiomatic API to create `Pair`s and
`Message`s.

The `Message` is then given to a logger, an instance of `logger.Interface`. A logger may augment the message with additional
`Pair`s. It's common for a `Logger` to add at least a `level` and a `ts` (timestamp) `Pair`, but loggers may add other `Pair`s.

Every logger uses a set of `Handler`s. A `Handler` is responsible for
* filtering `Message`s i.e. discarding them based on a level threshold
* formatting the `Message` using a formatter - an instance of `formatter.Interface`
* delivering the `Message` using an output - an instance of `output.Interface`

Loggers may be nested. Nesting a logger allows setting some default pairs to be added to every message emitted by a
nested logger. A typical scenario sets the _category_ pair for a nested logger. A nested logger has no handler on
its own but resorts to the _parent logger_ for actually delivering the message. Loggers may be nested at any depth.

### Log Format

`kvlog` supports several log formats out of the box. Others can be added by implementing a 
`formatter.Interface`.

#### kvformat

The default format used by `kvlog` follows the defaults of the 
[logstash KV filter](https://www.elastic.co/guide/en/logstash/current/plugins-filters-kv.html). The following lines
show examples of the log output

```
ts=2019-08-16T12:58:22+02:00 lvl=info evt=started
ts=2019-08-16T12:58:34+02:00 lvl=info evt=request dur=0.001s method=GET status=200 url=/
ts=2019-08-16T12:58:34+02:00 lvl=info evt=request dur=0.000s method=GET status=404 url=/favicon.ico
ts=2019-08-16T12:58:35+02:00 lvl=info evt=request dur=0.009s method=POST status=200 url=/pdf 
```

#### Terminal

`kvlog` provides a terminal format which looks similar to the kvformat described above but uses
ANSI color escape sequences to produce colored output suitable to be read by humans.

#### JSON Lines

`kvlog` provides another formatter that renders JSON lines (one line for each message). This format is
especially useful when using log collectors such as logstash, GrayLog or fluentd. 

The above example would look like the following in jsonl.

```json
{"ts":"2019-08-16T12:58:22+02:00","lvl":"info","evt":"started"}
{"ts":"2019-08-16T12:58:34+02:00","lvl":"info","evt":"request","dur":"0.001s","method":"GET","status":200,"url":"/"}
{"ts":"2019-08-16T12:58:34+02:00","lvl":"info","evt":"request","dur":"0.000s","method":"GET","status":404,"url":"/favicon.ico"}
{"ts":"2019-08-16T12:58:35+02:00","lvl":"info","evt":"request","dur":"0.009s","method":"POST","status":200,"url":"/pdf"}
```

## Installation

```
$ go get -u github.com/halimath/kvlog
```

`kvlog` requires Go >= 1.11 and has no dependencies except for the standard library.

## Usage

`kvlog` can be used in different ways. 

### Module functions

The most simple usage uses module functions.

```go
package main

import (
    "github.com/halimath/kvlog"
)

func main () {
    // ...

    kvlog.Info(kvlog.Evt("started"), kvlog.KV("port", 8080))
}
```

The module provides functions for all log level (`Debug`, `Info`, `Warn`, `Error`) as well as 
a configuration function for initializing the package level logger (i.e. configuring output 
and threshold as well as other filters). The default is to log everything of level `Info` or 
above to `stdout` using the console formatter (when connected to a terminal) or the JSON lines
formatter (when connected to a file/stream).

### Logger instance

A more advanced usage involves a dedicated `Logger` instance which can be used in dependency injection scenarios.

```go
package main

import (
	"os"

	"github.com/halimath/kvlog"
	"github.com/halimath/kvlog/filter"
	"github.com/halimath/kvlog/formatter/kvformat"
	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/logger"
	"github.com/halimath/kvlog/msg"
)

func main () {
	l := logger.New(handler.New(kvformat.Formatter, os.Stdout, filter.Threshold(msg.LevelInfo)))

	name, _ := os.Hostname()
	l.Info(kvlog.Evt("appStarted"), kvlog.KV("hostname", name))

    // ...
}
```

### HTTP Middleware

`kvlog` contains a HTTP middleware that generates an access log. It wraps another `http.Hander` allowing you to
log only requests on those handlers you are interested in.

```go
package main

import (
	"net/http"

	"github.com/halimath/kvlog"
)

func main() {
    mux := http.NewServeMux()
    // ...
	kvlog.Info(kvlog.Evt("started"))
	http.ListenAndServe(":8000", kvlog.Middleware(kvlog.L, mux))
}
```

### Default Keys

The following table lists the default keys used by `kvlog`.

Key | Value Type | Description
-- | -- | --
`level` | `debug`; `info`; `warn`; `error`; `unknown` | The log level used when issuing the message
`ts` | Timestamp in RFC3339 format| The timestamp when the message was created formatted as a string in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339).
`evt` | case sensitive string | A descriptive token of the event that caused this log msg.
`err` | string | A free form string containing a textual description of an error.
`msg` | string | A free form string containing a human readable msg.
`dur` | float | A duration measured in seconds as as floating point number.
`cat` | string | A string describing the message's category, such as `http`, `service`, `persistence`, ...

# Changelog

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

[ci-img-url]: https://github.com/halimath/kvlog/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/kvlog
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/kvlog
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/kvlog
[release-img-url]: https://img.shields.io/github/v/release/halimath/kvlog.svg
[release-url]: https://github.com/halimath/kvlog/releases



