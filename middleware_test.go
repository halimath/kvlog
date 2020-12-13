//
// This file is part of kvlog.
//
// Copyright 2019, 2020 Alexander Metzner.
//
// Copyright 2019, 2020 Alexander Metzner.
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
)

func TestMiddleware(t *testing.T) {
	var out bytes.Buffer
	logger := NewLogger(NewHandler(KVFormatter, &out))

	handler := Middleware(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Foo", "bar")
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("hello, world"))
	}))

	req := httptest.NewRequest("get", "/test/path", nil)
	w := httptest.NewRecorder()

	now := time.Now().Format("2006-01-02T15:04:05")
	handler.ServeHTTP(w, req)

	logger.Close()

	expected := fmt.Sprintf("ts=%s level=info event=request method=get url=</test/path> status=204 duration=0.000s\n", now)

	if expected != out.String() {
		t.Errorf("expected\n%s but got\n%s", expected, out.String())
	}
}
