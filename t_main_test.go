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

func eq(exp, act any) {
	if !r.DeepEqual(exp, act) {
		panic(fmt.Errorf(`
expected (detailed):
	%#[1]v
actual (detailed):
	%#[2]v
expected (simple):
	%[1]v
actual (simple):
	%[2]v
`, exp, act))
	}
}

func neq(one, other any) {
	if r.DeepEqual(one, other) {
		panic(fmt.Errorf(`
unexpected identical values (detailed):
	%#[1]v
unexpected identical values (simple):
	%[1]v
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

func funcName(val any) string {
	return runtime.FuncForPC(r.ValueOf(val).Pointer()).Name()
}

func catchAny(fun func()) (val any) {
	defer recAny(&val)
	fun()
	return
}

func recAny(ptr *any) { *ptr = recover() }

func eqDeref(exp, ptr any) {
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

func tryInterface(val any, err error) any {
	try(err)
	return val
}

func set(ptr, val any) {
	r.ValueOf(ptr).Elem().Set(r.ValueOf(val))
}

func setZero(ptr gt.Zeroable) {
	set(ptr, r.Zero(r.TypeOf(ptr).Elem()).Interface())
	eq(true, ptr.IsZero())
}

func jsonBytes(val any) []byte {
	return tryByteSlice(json.Marshal(val))
}

func list(vals ...any) []any { return vals }

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
