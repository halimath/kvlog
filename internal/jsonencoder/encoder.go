// Package jsonencoder provides types and functions to encode JSON values.
package jsonencoder

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	_ = iota
	nestedObject
	nestedArray

	DefaultNestingDepth = 10
	DefaultBufferSize   = 1024
)

type nestedStructure struct {
	typ          int
	valueWritten bool
	keyWritten   bool
}

// Encoder implements JSON lines writing. A writer uses buffering to buffer
// creation of unfinished JSON values and flushes them to the underlying
// writer together with a new line ('\n').
type Encoder struct {
	buf []byte

	// nestingStack keeps track of opened nesting structures such as
	// JSON objects and arrays. Each entry represents a nested structure.
	// true is used to mark objects, false is used to mark arrays.
	// For efficiency reasons this stack is preallocated and
	// nestingStackPointer keeps track of the current position.
	nestingStack []nestedStructure

	nestingStackPointer int
}

// New creates a new Writer.
func New() *Encoder {
	return NewWithBufferSize(DefaultBufferSize)
}

func NewWithBufferSize(size uint) *Encoder {
	nestingStack := make([]nestedStructure, DefaultNestingDepth)

	return &Encoder{
		buf:                 make([]byte, 0, size),
		nestingStack:        nestingStack,
		nestingStackPointer: -1,
	}

}

// Reset resets the Writer to start with a fresh state.
func (w *Encoder) Reset() {
	w.buf = w.buf[0:0]
	w.nestingStackPointer = -1
}

// Bytes returns a byte slice containing the written bytes.
func (w *Encoder) Bytes() []byte {
	return w.buf
}

// String returns a string representation of the data written.
func (w *Encoder) String() string {
	return string(w.buf)
}

func (w *Encoder) isNested() bool {
	return w.nestingStackPointer > -1
}

func (w *Encoder) isNestedObject() bool {
	return w.nestingStackPointer > -1 && w.nestingStack[w.nestingStackPointer].typ == nestedObject
}

func (w *Encoder) increaseNesting() {
	if w.nestingStackPointer >= len(w.nestingStack) {
		w.nestingStack = append(w.nestingStack, nestedStructure{})
	}

	w.nestingStackPointer++
	w.nestingStack[w.nestingStackPointer].keyWritten = false
	w.nestingStack[w.nestingStackPointer].valueWritten = false
}

func (w *Encoder) beforeValue() {
	if !w.isNested() {
		return
	}

	if w.isNestedObject() {
		if w.nestingStack[w.nestingStackPointer].keyWritten {
			w.nestingStack[w.nestingStackPointer].keyWritten = false
		} else {
			panic("must write key before writing value")
		}
	} else {
		if w.nestingStack[w.nestingStackPointer].valueWritten {
			w.writeByte(',')
		} else {
			w.nestingStack[w.nestingStackPointer].valueWritten = true
		}
	}
}

func (w *Encoder) writeByte(bt byte) *Encoder {
	w.buf = append(w.buf, bt)
	return w
}

func (w *Encoder) writeString(s string) *Encoder {
	w.buf = append(w.buf, []byte(s)...)
	return w
}

// String outputs s formatted as a JSON string.
func (w *Encoder) Str(s string) *Encoder {
	w.beforeValue()
	w.writeJSONString(s)
	return w
}

// Int outputs i formatted as a JSON number.
func (w *Encoder) Int(i int64) *Encoder {
	w.beforeValue()
	w.buf = strconv.AppendInt(w.buf, i, 10)
	return w
}

// Float outputs f formatted as a JSON number.
func (w *Encoder) Float(f float64) *Encoder {
	w.beforeValue()
	w.buf = strconv.AppendFloat(w.buf, f, 'e', 8, 64)
	return w
}

// Bool outputs bol formatted as a JSON boolean.
func (w *Encoder) Bool(bol bool) *Encoder {
	w.beforeValue()
	w.buf = strconv.AppendBool(w.buf, bol)
	return w
}

// Null outputs a literal JSON "null".
func (w *Encoder) Null() *Encoder {
	w.beforeValue()
	w.writeString("null")
	return w
}

// StartObject starts a new JSON object.
func (w *Encoder) StartObject() *Encoder {
	w.beforeValue()
	w.writeByte('{')
	w.increaseNesting()
	w.nestingStack[w.nestingStackPointer].typ = nestedObject

	return w
}

// Key outputs k as a JSON object key.
func (w *Encoder) Key(k string) *Encoder {
	if !w.isNestedObject() {
		panic("current nested structure is not an object")
	}

	if w.nestingStack[w.nestingStackPointer].keyWritten {
		panic("must write value before the next key")
	}
	w.nestingStack[w.nestingStackPointer].keyWritten = true

	if w.nestingStack[w.nestingStackPointer].valueWritten {
		w.writeByte(',')
	} else {
		w.nestingStack[w.nestingStackPointer].valueWritten = true
	}

	w.writeJSONString(k)
	w.writeByte(':')

	return w
}

// EndObject ends the currently open JSON object.
func (w *Encoder) EndObject() *Encoder {
	w.writeByte('}')
	w.nestingStackPointer--
	return w
}

// StartArray starts a JSON array.
func (w *Encoder) StartArray() *Encoder {
	w.beforeValue()
	w.writeByte('[')
	w.increaseNesting()
	w.nestingStack[w.nestingStackPointer].typ = nestedArray

	return w
}

// EndArray end the currently open JSON array.
func (w *Encoder) EndArray() *Encoder {
	w.writeByte(']')
	w.nestingStackPointer--
	return w
}

func (w *Encoder) writeJSONString(s string) {
	w.buf = append(w.buf, '"')

	var start, i int

	for {
		if start+i >= len(s) {
			break
		}

		r, l := utf8.DecodeRuneInString(s[start+i:])

		if needsEscaping(r) {
			w.buf = append(w.buf, s[start:start+i]...)
			i += l
			start = i

			switch r {
			case '"':
				w.buf = append(w.buf, '\\', '"')
			case '\\':
				w.buf = append(w.buf, '\\', '\\')
			case '\b':
				w.buf = append(w.buf, '\\', 'b')
			case '\f':
				w.buf = append(w.buf, '\\', 'f')
			case '\n':
				w.buf = append(w.buf, '\\', 'n')
			case '\r':
				w.buf = append(w.buf, '\\', 'r')
			case '\t':
				w.buf = append(w.buf, '\\', 't')
			default:
				w.buf = append(w.buf, []byte(fmt.Sprintf("\\u%x", r))...)
			}
		} else {
			i += l
		}
	}

	w.buf = append(w.buf, s[start:]...)

	w.buf = append(w.buf, '"')
}

var (
	escapedRunes = []byte{
		'"',
		'\\',
		'\b',
		'\f',
		'\n',
		'\r',
		'\t',
	}
)

func needsEscaping(r rune) bool {
	return bytes.ContainsRune(escapedRunes, r) || unicode.IsControl(r) || r > 127
}
