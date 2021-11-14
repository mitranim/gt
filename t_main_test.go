package gt_test

import (
	"encoding/json"
	"fmt"
	r "reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"

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

func panics(t testing.TB, msg string, fun func()) {
	t.Helper()
	val := catchAny(fun)

	if val == nil {
		t.Fatalf(`expected %v to panic, found no panic`, funcName(fun))
	}

	str := fmt.Sprint(val)
	if !strings.Contains(str, msg) {
		t.Fatalf(`
expected %v to panic with a message containing:
	%v
found the following message:
	%v
`, funcName(fun), msg, str)
	}
}

func funcName(val interface{}) string {
	return runtime.FuncForPC(r.ValueOf(val).Pointer()).Name()
}

func catchAny(fun func()) (val interface{}) {
	defer recAny(&val)
	fun()
	return
}

func recAny(ptr *interface{}) { *ptr = recover() }

func eqDeref(exp, ptr interface{}) {
	eq(exp, r.ValueOf(ptr).Elem().Interface())
}

// nolint:structcheck
type slice struct {
	dat uintptr
	len int
	cap int
}

func sliceShared(prev, next []byte) {
	prevSlice := *(*slice)(unsafe.Pointer(&prev))
	nextSlice := *(*slice)(unsafe.Pointer(&next))

	if prevSlice.dat != nextSlice.dat {
		panic(errSlices(`two slices unexpectedly don't share the backing array`, prev, next))
	}
}

func sliceNotShared(prev, next []byte) {
	prevSlice := *(*slice)(unsafe.Pointer(&prev))
	nextSlice := *(*slice)(unsafe.Pointer(&next))

	if prevSlice.dat == nextSlice.dat {
		panic(errSlices(`two slices unexpectedly share the backing array`, prev, next))
	}
}

func errSlices(prefix string, prev, next []byte) error {
	return fmt.Errorf(`
%[1]v;
prev (detailed):
	%#[2]v
next (detailed):
	%#[3]v
prev (simple):
	%[2]q
next (simple):
	%[3]q
`, prefix, prev, next)
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

func pathSegments(val string) []string {
	return []string{
		val,
		`/` + val,
		val + `/`,
		`/` + val + `/`,
	}
}

var (
	pathSegmentsOne   = pathSegments(`one`)
	pathSegmentsTwo   = pathSegments(`two`)
	pathSegmentsThree = pathSegments(`three`)
)
