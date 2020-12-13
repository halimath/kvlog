package kvlog

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	table := map[Level]string{
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
	}

	for level, expected := range table {
		actual := level.String()
		if actual != expected {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	}
}

func TestMessageLevel(t *testing.T) {
	m := NewMessage(LevelInfo, KV("foo", "bar"))
	if m.Level() != LevelInfo {
		t.Errorf("expected info but got %s", m.Level())
	}
}
