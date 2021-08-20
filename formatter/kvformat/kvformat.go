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

// Package kvformat provides a formatter.Interface writing messages in the KV format.
package kvformat

import (
	"fmt"
	"io"
	"sort"

	"github.com/halimath/kvlog/formatter"
	"github.com/halimath/kvlog/internal/msgutil"
	"github.com/halimath/kvlog/msg"
)

// Formatter implements a formatter.Interface that writes the default KV format.
var Formatter = formatter.FormatterFunc(formatMessageAsKV)

func formatMessageAsKV(m msg.Message, w io.Writer) error {
	sorted := msgutil.SortByKey(m)
	sort.Sort(sorted)

	for i, p := range sorted {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		formatPair(p, w)
	}

	w.Write([]byte("\n"))

	return nil
}

func formatPair(k msg.KVPair, w io.Writer) (err error) {
	if _, err := fmt.Fprintf(w, "%s=", k.Key); err != nil {
		return err
	}

	return msgutil.FormatValue(k, w)
}
