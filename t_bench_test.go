package gt_test

import (
	"path"
	"testing"

	"github.com/mitranim/gt"
)

func Benchmark_ParseNullUuid_simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gt.ParseNullUuid(`a915f35f0a3344bc8b9fb36bb650708d`)
	}
}

func Benchmark_ParseNullUuid_canon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gt.ParseNullUuid(`c230ed9a-e855-469c-8ebb-59c565aacaa7`)
	}
}

func Benchmark_NullUuid_string(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = val.String()
	}
}

func Benchmark_Uuid_GoString(b *testing.B) {
	val := gt.ParseUuid(`6b4c96c70bbc4d57a673de9620688f01`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = val.GoString()
	}
}

func Benchmark_NullUuid_GoString(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = val.GoString()
	}
}

func Benchmark_RandomUuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gt.RandomUuid()
	}
}

func Benchmark_RandomUuid_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gt.RandomUuid().String()
	}
}

func Benchmark_ParseNullDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gt.ParseNullDate(`1234-05-06`)
	}
}

func Benchmark_NullDate_String(b *testing.B) {
	val := gt.ParseNullDate(`1234-05-06`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = val.String()
	}
}

func Benchmark_ParseNullInterval(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)
	}
}

func Benchmark_NullInterval_String(b *testing.B) {
	val := gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = val.String()
	}
}

func Benchmark_path_Join_one(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path.Join(`one`)
	}
}

func Benchmark_gt_Join_one(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gt.Join(`one`)
	}
}

func Benchmark_path_Join_many(b *testing.B) {
	for i := 0; i < b.N; i++ {
		path.Join(`one`, `two`, `three`)
	}
}

func Benchmark_gt_Join_many(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gt.Join(`one`, `two`, `three`)
	}
}

func Benchmark_ParseNullUrl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gt.ParseNullUrl(`https://user:pass@one.two.three/four/five/six?one=two&one=three&four=five&five=six&seven=eight&nine=ten&nine=eleven#hash`)
	}
}
