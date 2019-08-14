// Package kvlog provides a key-value based logging system primary targetting
// container based deployments.
package kvlog

// Level defines the valid log levels
type Level int

// String provides a string representation of the log level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}

const (
	// LevelDebug log level
	LevelDebug Level = iota
	// LevelInfo log level
	LevelInfo
	// LevelWarn log level
	LevelWarn
	// LevelError log level
	LevelError
)

// --

// Logger implements a logger component.
// The output is written to the given output.
type Logger struct {
	out       Output
	Threshold Level
}

// NewLogger constructs a new Logger and returns a pointer to it.
func NewLogger(out Output, threshold Level) *Logger {
	return &Logger{
		out:       out,
		Threshold: threshold,
	}
}

func (l *Logger) Log(m Message) {
	if m.Level() < l.Threshold {
		return
	}

	l.out.WriteLogMessage(m)
}

func (l *Logger) Debug(pairs ...KVPair) {
	l.Log(NewMessage(LevelDebug, pairs...))
}

func (l *Logger) Info(pairs ...KVPair) {
	l.Log(NewMessage(LevelInfo, pairs...))
}

func (l *Logger) Warn(pairs ...KVPair) {
	l.Log(NewMessage(LevelWarn, pairs...))
}

func (l *Logger) Error(pairs ...KVPair) {
	l.Log(NewMessage(LevelError, pairs...))
}
