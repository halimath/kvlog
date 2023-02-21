package kvlog

import (
	"context"
	"strings"
	"testing"
)

func TestContext(t *testing.T) {
	var buf strings.Builder
	lo := New(NewSyncHandler(&buf, JSONLFormatter()))

	ctx := context.Background()

	l := FromContext(ctx)
	l.Logs("1")

	ctx = ContextWithLogger(ctx, lo)

	l = FromContext(ctx)
	l.Logs("2")

	if strings.TrimSpace(buf.String()) != `{"msg":"2"}` {
		t.Errorf("unexpected log output %q", buf.String())
	}
}
