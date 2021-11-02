package gt_test

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	if !reflect.DeepEqual(exp, act) {
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
	if reflect.DeepEqual(one, other) {
		panic(fmt.Errorf(`
unexpected identical values (detailed):
	%#[1]v
unexpected identical values (simple):
	%[1]s
`, one))
	}
}

func eqPtr(exp, ptr interface{}) {
	eq(exp, reflect.ValueOf(ptr).Elem().Interface())
}

func fail(err error) {
	if err == nil {
		panic("expected failure")
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
	reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(val))
}

func setZero(ptr gt.Zeroable) {
	set(ptr, reflect.Zero(reflect.TypeOf(ptr).Elem()).Interface())
	eq(true, ptr.IsZero())
}

func jsonBytes(val interface{}) []byte {
	return tryByteSlice(json.Marshal(val))
}

func counter(n int) []struct{} { return make([]struct{}, n) }

type DateTuple [3]int

func DateTupleFrom(a, b, c int) DateTuple { return DateTuple{a, b, c} }
