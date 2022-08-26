package jsonencoder

import (
	"encoding/json"
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
	w := New()
	for i := 0; i < b.N; i++ {
		w.Reset()
		w.StartObject()
		w.Key("foo").Str("bar")
		w.Key("spam").Int(12)
		w.Key("eggs").StartObject()
		w.Key("foo").Str("bar")
		w.Key("spam").Str("eggs")
		w.EndObject()
		w.EndObject()
	}
}
