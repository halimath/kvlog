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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var out bytes.Buffer
	logger := New(NewSyncHandler(&out, JSONLFormatter()))

	handler := Middleware(logger, true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		FromContext(r.Context()).Logs("from context")
		w.Header().Add("X-Foo", "bar")
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("hello, world"))
	}))

	req := httptest.NewRequest("get", "/test/path", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	req = httptest.NewRequest("delete", "/test/anotherpath", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	expected := `{"url":"/test/path","method":"get","msg":"from context"}
{"url":"/test/path","method":"get","msg":"request","dur":"0.000s","status":204}
{"url":"/test/anotherpath","method":"delete","msg":"from context"}
{"url":"/test/anotherpath","method":"delete","msg":"request","dur":"0.000s","status":204}
`

	if expected != out.String() {
		t.Errorf("expected\n%s but got\n%s", expected, out.String())
	}
}
