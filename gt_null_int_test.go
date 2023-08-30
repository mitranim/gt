package gt_test

import (
	"testing"

	"github.com/mitranim/gt"
)

// TODO: test various invalid inputs.
func TestNullInt(t *testing.T) {
	t.Run(`Decodable/sql.Scanner`, func(t *testing.T) {
		var (
			primZero    = int64(0)
			primNonZero = int64(123)
			zero        = gt.NullInt(primZero)
			nonZero     = gt.NullInt(primNonZero)
			dec         = new(gt.NullInt)
		)

		t.Run(`int`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, int(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, int(primNonZero))
		})
		t.Run(`int8`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, int8(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, int8(primNonZero))
		})
		t.Run(`int16`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, int16(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, int16(primNonZero))
		})
		t.Run(`int32`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, int32(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, int32(primNonZero))
		})
	})
}
