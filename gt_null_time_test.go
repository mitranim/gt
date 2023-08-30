package gt_test

import (
	"testing"
	"time"

	"github.com/mitranim/gt"
)

func TestNullTime(t *testing.T) {
	t.Run(`Parse`, func(t *testing.T) {
		test := func(exp gt.NullTime, src string) {
			t.Helper()
			eq(exp, gt.ParseNullTime(src))
		}

		test(gt.NullDateUTC(1970, 1, 1), `-0`)
		test(gt.NullDateUTC(1970, 1, 1), `+0`)
		test(gt.NullDateUTC(1970, 1, 1), `0`)

		test(gt.NullTime(time.UnixMilli(-1234567890123).In(time.UTC)), `-1234567890123`)
		test(gt.NullTimeUTC(1930, 11, 18, 0, 28, 29, 877000000), `-1234567890123`)

		test(gt.NullTime(time.UnixMilli(+1234567890123).In(time.UTC)), `+1234567890123`)
		test(gt.NullTimeUTC(2009, 2, 13, 23, 31, 30, 123000000), `+1234567890123`)

		test(gt.NullTime(time.UnixMilli(1234567890123).In(time.UTC)), `1234567890123`)
		test(gt.NullTimeUTC(2009, 2, 13, 23, 31, 30, 123000000), `1234567890123`)

		test(gt.NullTimeUTC(1234, 5, 6, 7, 8, 9, 0), `1234-05-06T07:08:09Z`)
	})

	t.Run(`AddInterval`, func(t *testing.T) {
		eq(
			gt.NullDateUTC(5, 7, 9),
			gt.NullDateUTC(1, 2, 3).AddInterval(gt.DateInterval(4, 5, 6)),
		)

		eq(
			gt.NullDateUTC(3, 2, 1),
			gt.NullDateUTC(4, 5, 6).AddInterval(gt.DateInterval(-1, -3, -5)),
		)

		eq(
			gt.NullTimeUTC(1, 2, 3, 4, 5, 6, 0),
			gt.NullDateUTC(1, 2, 3).AddInterval(gt.TimeInterval(4, 5, 6)),
		)

		eq(
			gt.NullTimeUTC(1, 2, 3, 3, 2, 1, 0),
			gt.NullTimeUTC(1, 2, 3, 4, 5, 6, 0).AddInterval(gt.TimeInterval(-1, -3, -5)),
		)

		eq(
			gt.NullTimeUTC(1+2, 2+3, 3+4, 4+5, 5+6, 6+7, 0),
			gt.NullTimeUTC(1, 2, 3, 4, 5, 6, 0).AddInterval(
				gt.IntervalFrom(2, 3, 4, 5, 6, 7),
			),
		)
	})

	t.Run(`NullDate`, func(t *testing.T) {
		eq(
			gt.NullDateFrom(1, 1, 1),
			gt.NullTime{}.NullDate(),
		)

		eq(
			gt.NullDateFrom(2, 3, 4),
			gt.NullDateUTC(2, 3, 4).NullDate(),
		)

		// Might fail at midnight.
		eq(
			gt.NullDateNow(),
			gt.NullTimeNow().NullDate(),
		)
	})

	t.Run(`Before`, func(t *testing.T) {
		eq(false, gt.NullTime{}.Before(gt.NullTime{}))
		eq(false, gt.NullDateUTC(1, 2, 3).Before(gt.NullTime{}))
		eq(true, gt.NullTime{}.Before(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 3).Before(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 4).Before(gt.NullDateUTC(1, 2, 3)))
		eq(true, gt.NullDateUTC(1, 2, 3).Before(gt.NullDateUTC(1, 2, 4)))
	})

	t.Run(`After`, func(t *testing.T) {
		eq(false, gt.NullTime{}.After(gt.NullTime{}))
		eq(true, gt.NullDateUTC(1, 2, 3).After(gt.NullTime{}))
		eq(false, gt.NullTime{}.After(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 3).After(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 3).After(gt.NullDateUTC(1, 2, 4)))
		eq(true, gt.NullDateUTC(1, 2, 4).After(gt.NullDateUTC(1, 2, 3)))
	})
}

func TestNullTimeLess(t *testing.T) {
	test := func(exp bool, src ...gt.NullTime) {
		t.Helper()
		eq(exp, gt.NullTimeLess(src...))
	}

	test(true, gt.NullTime{})
	test(false, gt.NullTime{}, gt.NullTime{})
	test(true, gt.NullDateUTC(1, 1, 1))
	test(true, gt.NullDateUTC(1, 1, 2))
	test(true, gt.NullTime{}, gt.NullDateUTC(1, 1, 2))
	test(false, gt.NullDateUTC(1, 1, 2), gt.NullDateUTC(1, 1, 2))
	test(true, gt.NullDateUTC(1, 1, 2), gt.NullDateUTC(1, 1, 3))
	test(false, gt.NullDateUTC(1, 1, 2), gt.NullDateUTC(1, 1, 3), gt.NullDateUTC(1, 1, 3))
	test(true, gt.NullDateUTC(1, 1, 2), gt.NullDateUTC(1, 1, 3), gt.NullDateUTC(1, 1, 4))
	test(false, gt.NullDateUTC(1, 1, 2), gt.NullTime{}, gt.NullDateUTC(1, 1, 3))
	test(false, gt.NullDateUTC(1, 1, 2), gt.NullDateUTC(1, 1, 3), gt.NullTime{})
	test(true, gt.NullTime{}, gt.NullDateUTC(1, 1, 3), gt.NullDateUTC(1, 1, 4))
}
