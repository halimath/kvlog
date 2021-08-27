// Package jsonl provides types and functions to writer JSON lines.
// See https://jsonlines.org/
package jsonl

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	_ = iota
	nestedObject
	nestedArray
)

type nestedStructure struct {
	typ          int
	valueWritten bool
	keyWritten   bool
}

// Writer implements JSON lines writing. A writer uses buffering to buffer
// creation of unfinished JSON values and flushes them to the underlying
// writer together with a new line ('\n').
type Writer struct {
	w   io.Writer
	buf bytes.Buffer

	// nestingStack keeps track of opened nesting structures such as
	// JSON objects and arrays. Each entry represents a nested structure.
	// true is used to mark objects, false is used to mark arrays
	nestingStack []nestedStructure
}

// New creates a new Writer that writes JSON lines to w.
func New(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) isNested() bool {
	return len(w.nestingStack) > 0
}

func (w *Writer) isNestedObject() bool {
	return len(w.nestingStack) > 0 && w.nestingStack[len(w.nestingStack)-1].typ == nestedObject
}

func (w *Writer) flush() {
	w.w.Write(w.buf.Bytes())
	w.buf.Reset()
}

func (w *Writer) afterValue() {
	if !w.isNested() {
		w.writeByte('\n')
		w.flush()
		return
	}
}

func (w *Writer) beforeValue() {
	if !w.isNested() {
		return
	}

	if w.isNestedObject() {
		if w.nestingStack[len(w.nestingStack)-1].keyWritten {
			w.nestingStack[len(w.nestingStack)-1].keyWritten = false
		} else {
			panic("must write key before writing value")
		}
	} else {
		if w.nestingStack[len(w.nestingStack)-1].valueWritten {
			w.writeByte(',')
		} else {
			w.nestingStack[len(w.nestingStack)-1].valueWritten = true
		}
	}
}

func (w *Writer) writeByte(bt byte) *Writer {
	w.buf.WriteByte(bt)
	return w
}

func (w *Writer) writeString(s string) *Writer {
	w.buf.WriteString(s)
	return w
}

// String outputs s formatted as a JSON string.
func (w *Writer) String(s string) *Writer {
	w.beforeValue()
	writeJSONString(&w.buf, s)
	w.afterValue()
	return w
}

// Int outputs i formatted as a JSON number.
func (w *Writer) Int(i int64) *Writer {
	w.beforeValue()
	w.writeString(strconv.FormatInt(i, 10))
	w.afterValue()
	return w
}

// Float outputs f formatted as a JSON number.
func (w *Writer) Float(f float64) *Writer {
	w.beforeValue()
	w.writeString(strconv.FormatFloat(f, 'e', 8, 64))
	w.afterValue()
	return w
}

// Bool outputs bol formatted as a JSON boolean.
func (w *Writer) Bool(bol bool) *Writer {
	w.beforeValue()
	if bol {
		w.writeString("true")
	} else {
		w.writeString("false")
	}
	w.afterValue()
	return w
}

// Null outputs a literal JSON "null".
func (w *Writer) Null() *Writer {
	w.beforeValue()
	w.writeString("null")
	w.afterValue()
	return w
}

// StartObject starts a new JSON object.
func (w *Writer) StartObject() *Writer {
	w.beforeValue()
	w.writeByte('{')
	w.nestingStack = append(w.nestingStack, nestedStructure{
		typ: nestedObject,
	})
	return w
}

// Key outputs k as a JSON object key.
func (w *Writer) Key(k string) *Writer {
	if !w.isNestedObject() {
		panic("current nested structure is not an object")
	}

	if w.nestingStack[len(w.nestingStack)-1].keyWritten {
		panic("must write value before the next key")
	}
	w.nestingStack[len(w.nestingStack)-1].keyWritten = true

	if w.nestingStack[len(w.nestingStack)-1].valueWritten {
		w.writeByte(',')
	} else {
		w.nestingStack[len(w.nestingStack)-1].valueWritten = true
	}

	writeJSONString(&w.buf, k)
	w.writeByte(':')

	return w
}

// EndObject ends the currently open JSON object.
func (w *Writer) EndObject() *Writer {
	w.writeByte('}')
	w.nestingStack = w.nestingStack[:len(w.nestingStack)-1]
	w.afterValue()
	return w
}

// StartArray starts a JSON array.
func (w *Writer) StartArray() *Writer {
	w.beforeValue()
	w.writeByte('[')
	w.nestingStack = append(w.nestingStack, nestedStructure{
		typ: nestedArray,
	})

	return w
}

// EndArray end the currently open JSON array.
func (w *Writer) EndArray() *Writer {
	w.writeByte(']')
	w.nestingStack = w.nestingStack[:len(w.nestingStack)-1]
	w.afterValue()
	return w
}

func writeJSONString(b *bytes.Buffer, s string) {
	b.Write([]byte{'"'})

	var start, i int

	for {
		if start+i >= len(s) {
			break
		}

		r, l := utf8.DecodeRuneInString(s[start+i:])

		if needsEscaping(r) {
			b.Write([]byte(s[start : start+i]))
			i += l
			start = i

			switch r {
			case '"':
				b.Write([]byte(`\"`))
			case '\\':
				b.Write([]byte(`\\`))
			case '\b':
				b.Write([]byte(`\b`))
			case '\f':
				b.Write([]byte(`\f`))
			case '\n':
				b.Write([]byte(`\n`))
			case '\r':
				b.Write([]byte(`\r`))
			case '\t':
				b.Write([]byte(`\t`))
			default:
				fmt.Fprintf(b, `\u%x`, r)
			}
		} else {
			i += l
		}
	}

	b.Write([]byte(s[start:]))

	b.Write([]byte{'"'})
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
