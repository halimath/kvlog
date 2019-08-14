package kvlog

import (
	"fmt"
	"io"
	"time"
)

// KVPair implements a key-value pair
type KVPair struct {
	// Key stores the key of the pair
	Key string

	// Value stores the value
	Value interface{}
}

// WriteTo writes the pair to the given writer
func (k KVPair) WriteTo(w io.Writer) error {
	var err error
	if _, err := fmt.Fprintf(w, "%s=", k.Key); err != nil {
		return err
	}

	switch x := k.Value.(type) {
	case time.Time:
		_, err = w.Write([]byte(x.Format("2006-01-02T15:04:05")))
	case time.Duration:
		_, err = fmt.Fprintf(w, "%.3fs", float64(x)/float64(time.Second))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		_, err = fmt.Fprintf(w, "%d", x)
	case float32, float64:
		_, err = fmt.Fprintf(w, "%.3f", x)
	case Level:
		_, err = w.Write([]byte(x.String()))
	default:
		_, err = fmt.Fprintf(w, "<%v>", x)
	}

	return err
}

// KV is a factory method for KVPair objects
func KV(key string, value interface{}) KVPair {
	return KVPair{
		Key:   key,
		Value: value,
	}
}
