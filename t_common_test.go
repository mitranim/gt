package gt_test

import (
	r "reflect"
	"testing"
	"time"

	"github.com/mitranim/gt"
)

// TODO: test various invalid inputs.
// TODO: verify that we can decode both RFC3339 timestamps and short dates.
func TestNullDate_common(t *testing.T) {
	var (
		primZero    = time.Time{}
		primNonZero = time.Date(1234, 5, 6, 0, 0, 0, 0, time.UTC)
		textZero    = ``
		textNonZero = `1234-05-06`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullDate{}
		nonZero     = gt.NullDateFrom(1234, 5, 6)
		dec         = new(gt.NullDate)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestNullFloat_common(t *testing.T) {
	var (
		primZero    = float64(0)
		primNonZero = float64(123)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullFloat(primZero)
		nonZero     = gt.NullFloat(primNonZero)
		dec         = new(gt.NullFloat)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestNullInt_common(t *testing.T) {
	var (
		primZero    = int64(0)
		primNonZero = int64(123)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullInt(primZero)
		nonZero     = gt.NullInt(primNonZero)
		dec         = new(gt.NullInt)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
// TODO: more tests for encoding and decoding.
func TestInterval_common(t *testing.T) {
	var (
		primZero    = `PT0S`
		primNonZero = `P1Y2M3DT4H5M6S`
		textZero    = `PT0S`
		textNonZero = primNonZero
		jsonZero    = jsonBytes(textZero)
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.Interval{}
		nonZero     = gt.IntervalFrom(1, 2, 3, 4, 5, 6)
		dec         = new(gt.Interval)
	)

	eq(false, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
// TODO: more tests for encoding and decoding.
func TestNullInterval_common(t *testing.T) {
	var (
		primZero    = ``
		primNonZero = `P1Y2M3DT4H5M6S`
		textZero    = ``
		textNonZero = primNonZero
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullInterval{}
		nonZero     = gt.NullIntervalFrom(1, 2, 3, 4, 5, 6)
		dec         = new(gt.NullInterval)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

func TestNullString_common(t *testing.T) {
	var (
		primZero    = string(``)
		primNonZero = string(`123`)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullString(primZero)
		nonZero     = gt.NullString(primNonZero)
		dec         = new(gt.NullString)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestNullTime_common(t *testing.T) {
	var (
		primZero    = time.Time{}
		primNonZero = time.Date(1234, 5, 6, 0, 0, 0, 0, time.UTC)
		textZero    = ``
		textNonZero = `1234-05-06T00:00:00Z`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullTime(primZero)
		nonZero     = gt.NullTime(primNonZero)
		dec         = new(gt.NullTime)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestNullUint_common(t *testing.T) {
	var (
		primZero    = uint64(0)
		primNonZero = uint64(123)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullUint(primZero)
		nonZero     = gt.NullUint(primNonZero)
		dec         = new(gt.NullUint)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestNullUrl_common(t *testing.T) {
	var (
		primZero    = ``
		primNonZero = `https://example.com`
		textZero    = primZero
		textNonZero = primNonZero
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullUrl{}
		nonZero     = gt.ParseNullUrl(textNonZero)
		dec         = new(gt.NullUrl)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestNullUuid_common(t *testing.T) {
	var (
		primZero    = ``
		primNonZero = [gt.UuidLen]byte{0xae, 0x68, 0xcc, 0xca, 0x87, 0xc3, 0x44, 0xaf, 0xa8, 0xa0, 0x20, 0x9c, 0xe, 0x20, 0x53, 0x43}
		textZero    = ``
		textNonZero = `ae68ccca87c344afa8a0209c0e205343`
		jsonZero    = bytesNull
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullUuid{}
		nonZero     = gt.ParseNullUuid(textNonZero)
		dec         = new(gt.NullUuid)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

// TODO: test various invalid inputs.
func TestUuid_common(t *testing.T) {
	var (
		primZero    = [gt.UuidLen]byte{}
		primNonZero = [gt.UuidLen]byte{0xae, 0x68, 0xcc, 0xca, 0x87, 0xc3, 0x44, 0xaf, 0xa8, 0xa0, 0x20, 0x9c, 0xe, 0x20, 0x53, 0x43}
		textZero    = `00000000000000000000000000000000`
		textNonZero = `ae68ccca87c344afa8a0209c0e205343`
		jsonZero    = jsonBytes(textZero)
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.Uuid{}
		nonZero     = gt.ParseUuid(textNonZero)
		dec         = new(gt.Uuid)
	)

	eq(false, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

/*
TODO:

  - Test various invalid inputs.
  - Ensure all three states are tested.
  - Test various auxiliary methods.
*/
func TestTer_common(t *testing.T) {
	var (
		primZero    = any(nil)
		primNonZero = true
		textZero    = ``
		textNonZero = `true`
		jsonZero    = bytesNull
		jsonNonZero = bytesTrue
		zero        = gt.TerNull
		nonZero     = gt.TerTrue
		dec         = new(gt.Ter)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

func TestRaw_common(t *testing.T) {
	var (
		primZero    = any(nil)
		primNonZero = []byte(`{"hello":"world"}`)
		textZero    = ``
		textNonZero = `{"hello":"world"}`
		jsonZero    = bytesNull
		jsonNonZero = primNonZero
		zero        = gt.Raw(nil)
		nonZero     = gt.Raw(primNonZero)
		dec         = new(gt.Raw)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

func testAny(
	t *testing.T,
	primZero, primNonZero any,
	textZero, textNonZero string,
	jsonZero, jsonNonZero []byte,
	zero, nonZero gt.Encodable,
	dec EncodableDecodable,
) {
	// Note: `primZero` is not tested here because it doesn't have to be Go zero.
	t.Run(`zeros`, func(t *testing.T) {
		eq(true, r.ValueOf(zero).IsZero())
		eq(false, r.ValueOf(primNonZero).IsZero())
		eq(false, r.ValueOf(nonZero).IsZero())
	})

	t.Run(`Encodable`, func(t *testing.T) {
		t.Run(`Zeroable`, func(t *testing.T) {
			eq(true, zero.IsZero())
			eq(false, nonZero.IsZero())
		})

		t.Run(`Nullable`, func(t *testing.T) {
			if zero.IsNull() {
				eq(true, zero.IsZero())
			}
			eq(false, nonZero.IsNull())
		})

		t.Run(`fmt.Stringer`, func(t *testing.T) {
			eq(textZero, zero.String())
			eq(textNonZero, nonZero.String())
		})

		t.Run(`encoding.TextMarshaler`, func(t *testing.T) {
			if zero.IsNull() {
				eq([]byte(nil), tryByteSlice(zero.MarshalText()))
			} else {
				eq([]byte(textZero), tryByteSlice(zero.MarshalText()))
			}
			eq([]byte(textNonZero), tryByteSlice(nonZero.MarshalText()))
		})

		t.Run(`json.Marshaler`, func(t *testing.T) {
			eq(jsonZero, jsonBytes(zero))
			eq(jsonNonZero, jsonBytes(nonZero))
		})

		t.Run(`driver.Valuer`, func(t *testing.T) {
			if zero.IsNull() {
				eq(nil, tryInterface(zero.Value()))
			} else {
				eq(primZero, tryInterface(zero.Value()))
			}
			eq(primNonZero, tryInterface(nonZero.Value()))
		})
	})

	t.Run(`Decodable`, func(t *testing.T) {
		t.Run(`Zeroer`, func(t *testing.T) {
			rval := r.ValueOf(dec).Elem()
			eq(true, rval.IsZero())

			set(dec, nonZero)
			eq(false, rval.IsZero())

			setZero(dec)
			eq(true, rval.IsZero())

			set(dec, nonZero)
			dec.Zero()

			/**
			Slice types such as `Raw` have more than one "zero" state. Going from
			non-zero to zero may secretly keep the capacity, and the resulting value
			is not equivalent to true zero.
			*/
			if rval.Kind() == r.Slice {
				eq(0, rval.Len())
			} else {
				eq(true, dec.IsZero())
			}
		})

		t.Run(`Parser`, func(t *testing.T) {
			t.Run(`empty`, func(t *testing.T) {
				set(dec, nonZero)
				eqDeref(nonZero, dec)

				if zero.IsNull() {
					try(dec.Parse(``))
					eq(true, dec.IsZero())
				} else {
					fail(dec.Parse(``))
				}
			})

			t.Run(`non-empty`, func(t *testing.T) {
				setZero(dec)
				eq(true, dec.IsZero())

				try(dec.Parse(textNonZero))
				eqDeref(nonZero, dec)
			})
		})

		t.Run(`encoding.TextUnmarshaler`, func(t *testing.T) {
			t.Run(`empty`, func(t *testing.T) {
				t.Run(`nil bytes`, func(t *testing.T) {
					set(dec, nonZero)
					eqDeref(nonZero, dec)

					if zero.IsNull() {
						try(dec.UnmarshalText(nil))
						eq(true, dec.IsZero())
					} else {
						fail(dec.UnmarshalText(nil))
					}
				})

				t.Run(`empty bytes`, func(t *testing.T) {
					set(dec, nonZero)
					eqDeref(nonZero, dec)

					if zero.IsNull() {
						try(dec.UnmarshalText([]byte{}))
						eq(true, dec.IsZero())
					} else {
						fail(dec.UnmarshalText([]byte{}))
					}
				})
			})

			t.Run(`non-empty`, func(t *testing.T) {
				setZero(dec)
				eq(true, dec.IsZero())

				try(dec.UnmarshalText([]byte(textNonZero)))
				eqDeref(nonZero, dec)
			})
		})

		t.Run(`json.Unmarshaler`, func(t *testing.T) {
			t.Run(`null`, func(t *testing.T) {
				set(dec, nonZero)
				eqDeref(nonZero, dec)

				if zero.IsNull() {
					try(dec.UnmarshalJSON(bytesNull))
					eq(true, dec.IsZero())
				} else {
					fail(dec.UnmarshalJSON(bytesNull))
				}
			})

			t.Run(`non-empty`, func(t *testing.T) {
				setZero(dec)
				eq(true, dec.IsZero())

				try(dec.UnmarshalJSON(jsonNonZero))
				eqDeref(nonZero, dec)
			})
		})

		t.Run(`sql.Scanner`, func(t *testing.T) {
			t.Run(`empty`, func(t *testing.T) {
				t.Run(`nil`, func(t *testing.T) {
					testScanEmpty(t, zero, nonZero, dec, nil)
				})

				t.Run(`nil bytes`, func(t *testing.T) {
					testScanEmpty(t, zero, nonZero, dec, []byte(nil))
				})

				t.Run(`empty bytes`, func(t *testing.T) {
					testScanEmpty(t, zero, nonZero, dec, []byte{})
				})

				t.Run(`empty string`, func(t *testing.T) {
					testScanEmpty(t, zero, nonZero, dec, ``)
				})
			})

			t.Run(`non-empty`, func(t *testing.T) {
				t.Run(`bytes`, func(t *testing.T) {
					testScanNonEmpty(t, zero, nonZero, dec, []byte(textNonZero))
				})

				t.Run(`string`, func(t *testing.T) {
					testScanNonEmpty(t, zero, nonZero, dec, textNonZero)
				})

				t.Run(`prim zero`, func(t *testing.T) {
					testScanNonEmpty(t, zero, zero, dec, primZero)
				})

				t.Run(`prim non-zero`, func(t *testing.T) {
					testScanNonEmpty(t, zero, nonZero, dec, primNonZero)
				})
			})
		})
	})
}

func testScanEmpty(
	t *testing.T,
	zero, nonZero gt.Encodable,
	dec EncodableDecodable,
	input any,
) {
	set(dec, nonZero)
	eqDeref(nonZero, dec)

	if zero.IsNull() {
		try(dec.Scan(input))
		eq(true, dec.IsZero())
	} else {
		fail(dec.Scan(input))
	}
}

func testScanNonEmpty(
	t *testing.T,
	zero, nonZero gt.Encodable,
	dec EncodableDecodable,
	input any,
) {
	setZero(dec)
	eq(true, dec.IsZero())

	try(dec.Scan(input))
	eqDeref(nonZero, dec)
}
