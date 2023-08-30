package gt_test

import (
	"testing"

	"github.com/mitranim/gt"
)

func TestFloat(t *testing.T) {
	t.Run(`Decodable/sql.Scanner`, func(t *testing.T) {
		var (
			primZero    = float64(0)
			primNonZero = float64(123)
			zero        = gt.NullFloat(primZero)
			nonZero     = gt.NullFloat(primNonZero)
			dec         = new(gt.NullFloat)
		)

		t.Run(`float32`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, float32(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, float32(primNonZero))
		})
	})
}
