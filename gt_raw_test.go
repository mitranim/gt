package gt_test

import (
	"fmt"
	"testing"

	"github.com/mitranim/gt"
)

func TestRaw(t *testing.T) {
	t.Run(`Grow`, func(t *testing.T) {
		prev := gt.Raw(nil)

		eq(0, len(prev))
		eq(0, cap(prev))

		next := prev.Grow(1)
		sliceNotShared(prev, next)
		eq(0, len(next))
		eq(1, cap(next))
		prev = next

		next = prev.Grow(11)
		sliceNotShared(prev, next)
		eq(0, len(next))
		eq(13, cap(next))
		prev = next

		next = prev.Grow(7)
		sliceShared(prev, next)
		eq(0, len(next))
		eq(13, cap(next))
		prev = next

		next = prev[:1][1:]
		sliceNotShared(prev, next)
		eq(0, len(next))
		eq(12, cap(next))
		prev = next

		next = prev.Grow(7)
		sliceShared(prev, next)
		eq(0, len(next))
		eq(12, cap(next))
		prev = next

		next = prev.Grow(11)
		sliceShared(prev, next)
		eq(0, len(next))
		eq(12, cap(next))
		prev = next

		next = prev.Grow(13)
		sliceNotShared(prev, next)
		eq(0, len(next))
		eq(37, cap(next))
		prev = next

		next = prev.Grow(0)
		sliceShared(prev, next)
		eq(prev, next)
		eq(0, len(next))
		eq(37, cap(next))
		prev = next
	})

	t.Run(`GoString`, func(t *testing.T) {
		test := func(exp string, val gt.Raw) {
			t.Helper()
			eq(exp, fmt.Sprintf(`%#v`, val))
		}

		test("gt.Raw(nil)", gt.Raw(nil))
		test("gt.Raw(nil)", gt.Raw{})
		test("gt.Raw(nil)", gt.Raw(``))
		test("gt.Raw(`one`)", gt.Raw(`one`))
		test("gt.Raw(`one two`)", gt.Raw(`one two`))
		test("gt.Raw(`\t`)", gt.Raw("\t"))
		test(`gt.Raw("\v")`, gt.Raw("\v"))
		test(`gt.Raw("\r")`, gt.Raw("\r"))
		test("gt.Raw(`{\"hello\": \"world\"}`)", gt.Raw(`{"hello": "world"}`))

		// Result of `strconv.CanBackquote` which only allows single lines.
		// Kind of unfortunate since backquoted strings can include newlines.
		test(`gt.Raw("\n")`, gt.Raw("\n"))
	})
}
