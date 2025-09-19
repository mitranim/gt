package gt_test

import (
	"math"
	"testing"

	"github.com/mitranim/gt"
)

func TestInt64ToHexUint(t *testing.T) {
	eq(gt.HexUint(0), gt.Int64ToHexUint(0))

	eq(
		gt.HexUint(0b0111111111111111111111111111111111111111111111111111111111111111),
		gt.Int64ToHexUint(math.MaxInt64),
	)

	eq(
		gt.HexUint(0b0111111111111111111111111111111111111111111111111111111111111111),
		gt.Int64ToHexUint(0b0111111111111111111111111111111111111111111111111111111111111111),
	)

	eq(
		gt.HexUint(0b1000000000000000000000000000000000000000000000000000000000000000),
		gt.Int64ToHexUint(math.MinInt64),
	)

	eq(
		gt.HexUint(0b1100000000000000000000000000000000000000000000000000000000000000),
		gt.Int64ToHexUint(-0b100000000000000000000000000000000000000000000000000000000000000),
	)

	eq(gt.HexUint(0xffffffffffffffff), gt.Int64ToHexUint(-1))
}

func TestHexUint_String(t *testing.T) {
	eq(``, gt.HexUint(0).String())
	eq(`0000000000000001`, gt.HexUint(0x1).String())
	eq(`0000000000000002`, gt.HexUint(0x2).String())
	eq(`000000000000000f`, gt.HexUint(0xf).String())
	eq(`00000000000000ff`, gt.HexUint(0xff).String())
	eq(`f000000000000000`, gt.HexUint(0xf000000000000000).String())
	eq(`ff00000000000000`, gt.HexUint(0xff00000000000000).String())
	eq(`ffffffffffffffff`, gt.HexUint(math.MaxUint64).String())
}

func TestHexUint_common(t *testing.T) {
	var (
		primZero    = int64(0)
		primNonZero = int64(0x0123456789abcdef)
		textZero    = ``
		textNonZero = `0123456789abcdef`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.HexUint(primZero)
		nonZero     = gt.HexUint(primNonZero)
		dec         = new(gt.HexUint)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}
