package gt_test

import (
	"encoding/json"
	"fmt"
	r "reflect"

	"github.com/mitranim/gt"
)

type EncodableDecodable interface {
	gt.Encodable
	gt.Decodable
}

var (
	bytesNull = []byte(`null`)
	bytesTrue = []byte(`true`)
)

func eq(exp, act interface{}) {
	if !r.DeepEqual(exp, act) {
		panic(fmt.Errorf(`
expected (detailed):
	%#[1]v
actual (detailed):
	%#[2]v
expected (simple):
	%[1]s
actual (simple):
	%[2]s
`, exp, act))
	}
}

func neq(one, other interface{}) {
	if r.DeepEqual(one, other) {
		panic(fmt.Errorf(`
unexpected identical values (detailed):
	%#[1]v
unexpected identical values (simple):
	%[1]s
`, one))
	}
}

func eqPtr(exp, ptr interface{}) {
	eq(exp, r.ValueOf(ptr).Elem().Interface())
}

func fail(err error) {
	if err == nil {
		panic(`expected error, got none`)
	}
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func tryByteSlice(val []byte, err error) []byte {
	try(err)
	return val
}

func tryInterface(val interface{}, err error) interface{} {
	try(err)
	return val
}

func set(ptr, val interface{}) {
	r.ValueOf(ptr).Elem().Set(r.ValueOf(val))
}

func setZero(ptr gt.Zeroable) {
	set(ptr, r.Zero(r.TypeOf(ptr).Elem()).Interface())
	eq(true, ptr.IsZero())
}

func jsonBytes(val interface{}) []byte {
	return tryByteSlice(json.Marshal(val))
}

func counter(count int) []struct{} { return make([]struct{}, count) }

func list(vals ...interface{}) []interface{} { return vals }

type intervalPart struct {
	int
	string
}

func intervalParts(suffix string) []intervalPart {
	return []intervalPart{
		{0, ``},
		{0, `0` + suffix},
		{0, `-0` + suffix},
		{1, `1` + suffix},
		{-1, `-1` + suffix},
		{19, `19` + suffix},
		{19, `019` + suffix},
		{-19, `-19` + suffix},
		{-19, `-019` + suffix},
	}
}
