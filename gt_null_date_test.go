package gt_test

import (
	"testing"
	"time"

	"github.com/mitranim/gt"
)

// This test might fail if invoked precisely at midnight.
// That would only validate our assumptions.
func TestNullDateNow(t *testing.T) {
	eq(list(time.Now().Date()), list(gt.NullDateNow().Date()))
}

func TestNullDate(t *testing.T) {
	t.Run(`String`, func(t *testing.T) {
		eq(``, gt.NullDateFrom(0, 0, 0).String())
		eq(`0001-01-01`, gt.NullDateFrom(1, 1, 1).String())
		eq(`0000-12-31`, gt.NullDateFrom(1, 1, 0).String())
	})

	t.Run(`AddDate`, func(t *testing.T) {
		eq(gt.NullDateFrom(0, 0, 0), gt.NullDateFrom(0, 0, 0).AddDate(0, 0, 0))
		eq(gt.NullDateFrom(0, 0, 0), gt.NullDateFrom(0, 0, 0).AddDate(1, 2, 3))
		eq(gt.NullDateFrom(0, 0, 0), gt.NullDateFrom(0, 0, 0).AddDate(-1, -2, -3))
		eq(gt.NullDateFrom(1, 2, 3), gt.NullDateFrom(1, 2, 3).AddDate(0, 0, 0))
		eq(gt.NullDateFrom(5, 7, 9), gt.NullDateFrom(1, 2, 3).AddDate(4, 5, 6))
		eq(gt.NullDateFrom(5, 7, 9), gt.NullDateFrom(4, 5, 6).AddDate(1, 2, 3))
		eq(gt.NullDateFrom(3, 2, 1), gt.NullDateFrom(4, 5, 6).AddDate(-1, -3, -5))
	})
}
