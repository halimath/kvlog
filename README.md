# kvlog

`kvlog` is a logging library based on key-value logging implemented using go (golang).

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
    // Optionally configure threshold
    kvlog.ConfigureThreshold(kvlog.LevelWarn)

    // ...

    kvlog.Info(kvlog.KV("event", "App started"))
}
```

The module provides methods for all log level (`Debug`, `Info`, `Warn`, `Error`) as well as configuration methods
for the threshold (`ConfigureThreshold`) which defaults to `info` and for the output (`ConfigureOutput`) which 
defaults to `stdout`.

### Logger instance

A more advanced usage involves a dedicated `Logger` instance which can be used in dependency injection scenarios.

```go
package main

import (
    "bitbucket.org/halimath/kvlog"
)

func main () {
    l := kvlog.NewLogger(kvlog.Stdout(), kvlog.LevelInfo)

    // ...

    l.Info(kvlog.KV("event", "App started"))
}
```

### HTTP handler

`kvlog` contains an HTTP access log handler, that can be used to wrap other `http.Hander`s.

```go

```

# Changelog

## 0.1.0

* Initial release

# License

```
Copyright 2019 Alexander Metzner.

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
