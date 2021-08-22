package jsonl

import (
	"encoding/json"
	"strings"
	"testing"
)

func Benchmark_MarshalMap(b *testing.B) {
	data := map[string]interface{}{
		"foo":  "bar",
		"spam": 1,
		"eggs": map[string]interface{}{
			"foo":  "bar",
			"spam": "eggs",
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_MarshalStruct(b *testing.B) {
	type nested struct {
		Foo  string `json:"foo"`
		Spam string `json:"spam"`
	}

	type dataType struct {
		Foo  string `json:"foo"`
		Spam int    `json:"spam"`
		Eggs nested `json:"eggs"`
	}

	// b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data := dataType{
			Foo:  "bar",
			Spam: i,
			Eggs: nested{
				Foo:  "bar",
				Spam: "eggs",
			},
		}

		_, err := json.Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_jsonl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b strings.Builder
		w := New(&b)
		w.StartObject()
		w.Key("foo").String("bar")
		w.Key("spam").Int(int64(i))
		w.Key("eggs").StartObject()
		w.Key("foo").String("bar")
		w.Key("spam").String("eggs")
		w.EndObject()
		w.EndObject()
	}
}
