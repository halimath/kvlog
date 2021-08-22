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

package kvlog

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/halimath/kvlog/formatter/kvformat"
	"github.com/halimath/kvlog/handler"
	"github.com/halimath/kvlog/logger"
)

func TestMiddleware(t *testing.T) {
	var out bytes.Buffer
	logger := logger.New(handler.New(kvformat.Formatter, &out))

	handler := Middleware(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Foo", "bar")
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("hello, world"))
	}))

	req := httptest.NewRequest("get", "/test/path", nil)
	w := httptest.NewRecorder()

	now := time.Now().Format(time.RFC3339)
	handler.ServeHTTP(w, req)

	logger.Close()

	expected := fmt.Sprintf("ts=%s lvl=info cat=http evt=request dur=0.000s method=get status=204 url=</test/path>\n", now)

	if expected != out.String() {
		t.Errorf("expected\n%s but got\n%s", expected, out.String())
	}
}
