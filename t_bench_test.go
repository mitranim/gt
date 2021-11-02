package gt_test

import (
	"testing"

	"github.com/mitranim/gt"
)

func Benchmark_ParseNullUuid_simple(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullUuid(`a915f35f0a3344bc8b9fb36bb650708d`)
	}
}

func Benchmark_ParseNullUuid_canon(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullUuid(`c230ed9a-e855-469c-8ebb-59c565aacaa7`)
	}
}

func Benchmark_NullUuid_string(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.String()
	}
}

func Benchmark_Uuid_GoString(b *testing.B) {
	val := gt.ParseUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.GoString()
	}
}

func Benchmark_NullUuid_GoString(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.GoString()
	}
}

func Benchmark_ParseNullDate(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullDate(`1234-05-06`)
	}
}

func Benchmark_NullDate_String(b *testing.B) {
	val := gt.ParseNullDate(`1234-05-06`)

	for range counter(b.N) {
		_ = val.String()
	}
}

func Benchmark_ParseNullInterval(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)
	}
}

func Benchmark_NullInterval_String(b *testing.B) {
	val := gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)

	for range counter(b.N) {
		_ = val.String()
	}
}
