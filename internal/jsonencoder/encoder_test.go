package jsonencoder

import (
	"strings"
	"testing"
)

func TestWriter_String(t *testing.T) {
	tab := map[string]string{
		"foo":              `"foo"`,
		"hello, world":     `"hello, world"`,
		"hello\nworld":     `"hello\nworld"`,
		"hello\tworld":     `"hello\tworld"`,
		"hello\rworld":     `"hello\rworld"`,
		"hello\fworld":     `"hello\fworld"`,
		"hello\bworld":     `"hello\bworld"`,
		"hello\\world":     `"hello\\world"`,
		"hello\"world":     `"hello\"world"`,
		"hello\u1234world": `"hello\u1234world"`,
	}

	for in, exp := range tab {
		w := New()
		w.Str(in)
		act := strings.TrimSpace(w.String())

		if exp != act {
			t.Errorf("'%s': expected '%s' got '%s'", in, exp, act)
		}
	}
}

func TestWriter_Int(t *testing.T) {
	tab := map[int64]string{
		1:    "1",
		5e12: "5000000000000",
	}

	for in, exp := range tab {
		w := New()
		w.Int(in)
		act := strings.TrimSpace(w.String())

		if exp != act {
			t.Errorf("%d: expected '%s' got '%s'", in, exp, act)
		}

	}
}

func TestWriter_Bool(t *testing.T) {
	tab := map[bool]string{
		true:  "true",
		false: "false",
	}

	for in, exp := range tab {
		w := New()
		w.Bool(in)
		act := strings.TrimSpace(w.String())

		if exp != act {
			t.Errorf("%v: expected '%s' got '%s'", in, exp, act)
		}
	}
}

func TestWriter_Null(t *testing.T) {
	w := New()
	w.Null()
	act := strings.TrimSpace(w.String())

	if act != "null" {
		t.Errorf("expected 'null' got '%s'", act)
	}
}

func TestWriter_Float(t *testing.T) {
	w := New()
	w.Float(1.2345)
	act := strings.TrimSpace(w.String())

	if act != "1.23450000e+00" {
		t.Errorf("expected '1.23450000e+00' got '%s'", act)
	}
}

func TestWriter_EmptyArray(t *testing.T) {
	w := New()
	w.StartArray()
	w.EndArray()
	act := strings.TrimSpace(w.String())

	if act != "[]" {
		t.Errorf("expected '[]' got '%s'", act)
	}
}

func TestWriter_ArrayWithSingleElement(t *testing.T) {
	w := New()
	w.StartArray()
	w.Int(2)
	w.EndArray()
	act := strings.TrimSpace(w.String())

	if act != "[2]" {
		t.Errorf("expected '[2]' got '%s'", act)
	}
}

func TestWriter_ArrayWithMultipleElements(t *testing.T) {
	w := New()
	w.StartArray()
	w.Int(2)
	w.Int(3)
	w.Int(4)
	w.EndArray()
	act := strings.TrimSpace(w.String())

	if act != "[2,3,4]" {
		t.Errorf("expected '[2,3,4]' got '%s'", act)
	}
}

func TestWriter_EmptyObject(t *testing.T) {
	w := New()
	w.StartObject()
	w.EndObject()
	act := strings.TrimSpace(w.String())

	if act != "{}" {
		t.Errorf("expected '{}' got '%s'", act)
	}
}
func TestWriter_ObjectWithSingleKey(t *testing.T) {
	w := New()
	w.StartObject()
	w.Key("foo")
	w.Int(1)
	w.EndObject()
	act := strings.TrimSpace(w.String())

	if act != `{"foo":1}` {
		t.Errorf(`expected '{"foo":1}' got '%s'`, act)
	}
}

func TestWriter_ObjectWithMultipleKeys(t *testing.T) {
	w := New()
	w.StartObject()
	w.Key("foo")
	w.Int(1)
	w.Key("bar")
	w.Int(2)
	w.EndObject()
	act := strings.TrimSpace(w.String())

	if act != `{"foo":1,"bar":2}` {
		t.Errorf(`expected '{"foo":1,"bar":2}' got '%s'`, act)
	}
}

func TestWriter_InvalidKeyUsage(t *testing.T) {
	type testCase func(*Encoder)
	expectPanic := func(t *testing.T, tc testCase) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()

		w := New()
		tc(w)
	}

	t.Run("key w/o object", func(t *testing.T) {
		expectPanic(t, func(w *Encoder) {
			w.Key("foo")
		})
	})

	t.Run("key two times", func(t *testing.T) {
		expectPanic(t, func(w *Encoder) {
			w.StartObject()
			w.Key("foo")
			w.Key("foo")
		})
	})

	t.Run("value w/o key", func(t *testing.T) {
		expectPanic(t, func(w *Encoder) {
			w.StartObject()
			w.Str("foo")
		})
	})
}

// func TestWriter_Array(t *testing.T) {
// 	b := NewBuffer()

// 	a := b.Array()

// 	a.String("foo")
// 	a.Int(1)
// 	a.Float(1.2)
// 	a.Null()
// 	a.Bool(true)
// 	o := a.Object()
// 	o.Key("foo").String("bar")
// 	o.End()
// 	a.End()

// 	exp := `["foo",1,1.20000000e+00,null,true,{"foo":"bar"}]`
// 	act := b.String()
// 	if act != exp {
// 		t.Errorf("expected '%s' but got '%s'", exp, act)
// 	}
// }
