package gt_test

import (
	"testing"

	"github.com/mitranim/gt"
)

// TODO: test various invalid inputs.
func TestNullUint(t *testing.T) {
	var (
		primZero    = uint64(0)
		primNonZero = uint64(123)
		zero        = gt.NullUint(primZero)
		nonZero     = gt.NullUint(primNonZero)
		dec         = new(gt.NullUint)
	)

	t.Run(`Decodable/sql.Scanner`, func(t *testing.T) {
		t.Run(`uint`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, uint(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, uint(primNonZero))
		})
		t.Run(`uint8`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, uint8(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, uint8(primNonZero))
		})
		t.Run(`uint16`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, uint16(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, uint16(primNonZero))
		})
		t.Run(`uint32`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, uint32(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, uint32(primNonZero))
		})
	})
}
