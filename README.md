# kvlog

`kvlog` is a logging library based on key-value logging implemented using go (golang).

## Description

`kvlog` provides types and functions to produce a log stream of key-value based log events. 
Key-value based log messages differ from conventional string-based log messages. They do not
contain a bare string message but any number of key-value-tuples which are encoded using a simple
to parse syntax. This allows log processor systems such as the [ELK-stack](https://www.elastic.co/de/what-is/elk-stack)
to analyze and index the log messages based on key-value-tuples.

### Components

`kvlog` is built from a set of _components_ that interact to implement logging functionality.

A `Message` is produced by the client. Every `Message` consists of `Pair`s each representing
a single key-value-pair. `kvlog` provides a convenient and idiomatic API to create `Pair`s and
`Message`s.

The `Message` is then given to a `Logger`. A `Logger` may augment the message with additional
`Pair`s. It's common for a `Logger` to add at least a `level` and a `ts` (timestamp) `Pair`, but
`Logger`s may add other `Pair`s.

Every `Logger` uses a set of `Handler`s. A `Handler` is responsible for
* formatting the `Message` using a `Formatter`
* delivering the `Message` using an `Output`

### Log Format

The format used by `kvlog` by default follows the defaults of the 
[logstash KV filter](https://www.elastic.co/guide/en/logstash/current/plugins-filters-kv.html). The following lines
show examples of the log output

```
ts=2019-08-16T12:58:22 level=info event=<started>
ts=2019-08-16T12:58:34 level=info event=<request> method=<GET> url=</> status=200 duration=0.001s
ts=2019-08-16T12:58:34 level=info event=<request> method=<GET> url=</favicon.ico> status=404 duration=0.000s
ts=2019-08-16T12:58:35 level=info event=<request> method=<POST> url=</pdf> status=200 duration=0.009s
```

## Installation

```
$ go get -u bitbucket.org/halimath/kvlog
```

## Usage

`kvlog` can be used in different ways. 

### Module functions

The most simple usage uses module functions.

```go
package main

import (
    "bitbucket.org/halimath/kvlog"
)

func main () {
    // ...

    kvlog.Info(kvlog.KV("event", "App started"))
}
```

The module provides functions for all log level (`Debug`, `Info`, `Warn`, `Error`) as well as a configuration function
for initializing the package level logger (i.e. configuring output and threshold as well as other filters). The default
is to log everything of level `Info` or above to `stdout` using the default log format.

### Logger instance

A more advanced usage involves a dedicated `Logger` instance which can be used in dependency injection scenarios.

```go
package main

import (
    "bitbucket.org/halimath/kvlog"
)

func main () {
    l := kvlog.NewLogger(kvlog.NewHandler(kvlog.KVFormatter, kvlog.Stdout(), kvlog.Threshold(kvlog.LevelWarn)))

    // ...

    l.Info(kvlog.KV("event", "App started"))
}
```

### HTTP Middleware

`kvlog` contains a HTTP middleware that generates an access log. It wraps another `http.Hander` allowing you to
log only requests on those handlers you are interested in.

```go
package main

import (
	"net/http"

	"bitbucket.org/halimath/kvlog"
)

func main() {
    mux := http.NewServeMux()
    // ...
	kvlog.Info(kvlog.KV("event", "started"))
	http.ListenAndServe(":8000", kvlog.Middleware(kvlog.L, mux))
}
```

# Changelog

## 0.4.0
* Export package level logger instance `L`

## 0.3.0
__Caution, breaking changes:__ This version provides a new API which is _not compatible_ to the 
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
