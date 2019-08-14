package kvlog

var L *Logger

func init() {
	L = NewLogger(Stdout(), LevelInfo)
}

func ConfigureOutput(out Output) {
	L = NewLogger(out, L.Threshold)
}

func ConfigureThreshold(t Level) {
	L = NewLogger(L.out, t)
}

func Debug(pairs ...KVPair) {
	L.Debug(pairs...)
}

func Info(pairs ...KVPair) {
	L.Info(pairs...)
}

func Warn(pairs ...KVPair) {
	L.Warn(pairs...)
}

func Error(pairs ...KVPair) {
	L.Error(pairs...)
}
