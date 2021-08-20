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

package output

import (
	"io"
	"os"
)

// Output defines the interface that must be implemented by types
// that handle output.
type Output interface {
	io.Writer
}

type nonClosingWriterOutput struct {
	io.Writer
}

func (n *nonClosingWriterOutput) Close() error {
	return nil
}

// Stdout returns an Output that writes the STDOUT but ignores any request to close the stream.
func Stdout() Output {
	return &nonClosingWriterOutput{
		Writer: os.Stdout,
	}
}

// Stderr returns an Output that writes the STDOUT but ignores any request to close the stream.
func Stderr() Output {
	return &nonClosingWriterOutput{
		Writer: os.Stderr,
	}
}
