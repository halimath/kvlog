package kvlog

import (
	"testing"
)

type mockOutput struct {
	messages []Message
}

func (m *mockOutput) WriteLogMessage(msg Message) {
	m.messages = append(m.messages, msg)
}

func TestLogger(t *testing.T) {
	o := mockOutput{}
	l := &Logger{
		out:       &o,
		Threshold: LevelWarn,
	}

	l.Debug(KV("foo", "bar"))
	l.Info(KV("foo", "bar"))
	l.Warn(KV("foo", "bar"))
	l.Error(KV("foo", "bar"))

	if len(o.messages) != 2 {
		t.Errorf("expected 2 messages but got %d", len(o.messages))
	}
}
