// This file is part of kvlog.
//
// Copyright 2019, 2020, 2021 Alexander Metzner.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package kvlog provides a structured logging facility. The underlying structure is based on key-value pairs.
// key-value pairs are rendered as [JSON lines] but other Formatters can be used to provide different outputs
// including custom ones.
//
// # Creating a Logger
//
// While package kvlog provides a ready-to-use Logger via the L variable, creating a custom logger gives you
// more flexibility on the logger's target and format. Creating a new Logger is done via the New function.
// It accepts any number of Handlers. Each Handler pairs an io.Writer as well as a Formatter.
//
// Handlers can be synchronous as well as asynchronous. Synchronous Handlers execute the Formatter as well as
// writing the output in the same goroutine that invoked the Logger. Asynchronous Handlers dispatch the log
// event to a different goroutine via a channel. Thus, asynchronous Handlers must be closed before shutdown
// in order to flush the channel and emit all log events.
//
// # Emitting Events
//
// Events are constructed by calling a Logger's With() method. This creates an Event that can be extended by
// calling its KV or delegating methods (such as Err or Dur). Eventually, submit the Event by invoking
// Log or Logf. Anything passed to Log is under under the Event's "msg" key.
//
// # Deriving Loggers
//
// Loggers can be derived from other Loggers. This enables to configure a set of key-value-pairs to be added
// to every Event emmitted via a deriverd logger. The syntax works similar to emitting log messages this time
// only invoking Logger() instead of Log.
//
// # Hooks
//
// In addition to deriving Loggers, any number of Hooks may be added to a Logger. The Hook's callback function
// is invoked everytime an Event is emitted via this Logger or any Logger's derived from it. Hooks are useful
// to add dynamic values, such as timestamps or anything else read from the surrounding context. Adding a
// timestamp to every log Event is realized via the TimeHook.
//
// # Formatters
//
// The kvlog package comes with three Formatters out of the box:
//   - JSONLFormatter formats events as JSON line values
//   - TerminalFormatter formats events for output on a terminal which includes colorizing the event
//   - KVFormatter formats events in the legacy KV-Format
//
// Custom formatters may be created by implementing the Formatter interface or using the FormatterFunc
// convenience type.
//
// [JSON lines]: https://jsonlines.org/
package kvlog
