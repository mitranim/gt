package gt_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/mitranim/gt"
)

// TODO: test various invalid inputs.
// TODO: verify that we can decode both RFC3339 timestamps and short dates.
func TestNullDate(t *testing.T) {
	var (
		primZero    = time.Time{}
		primNonZero = time.Date(1234, 5, 6, 0, 0, 0, 0, time.UTC)
		textZero    = ``
		textNonZero = `1234-05-06`
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullDate{}
		nonZero     = gt.NullDateFrom(1234, 5, 6)
		dec         = new(gt.NullDate)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)

	t.Run(`fmt.Stringer`, func(t *testing.T) {
		eq(``, gt.NullDateFrom(0, 0, 0).String())
		eq(`0001-01-01`, gt.NullDateFrom(1, 1, 1).String())
		eq(`0000-12-31`, gt.NullDateFrom(1, 1, 0).String())
	})
}

// TODO: test various invalid inputs.
func TestNullFloat(t *testing.T) {
	var (
		primZero    = float64(0)
		primNonZero = float64(123)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullFloat(primZero)
		nonZero     = gt.NullFloat(primNonZero)
		dec         = new(gt.NullFloat)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)

	t.Run(`Decodable/sql.Scanner`, func(t *testing.T) {
		t.Run(`float32`, func(t *testing.T) {
			testScanEmpty(t, zero, nonZero, dec, float32(primZero))
			testScanNonEmpty(t, zero, nonZero, dec, float32(primNonZero))
		})
	})
}

// TODO: test various invalid inputs.
func TestNullInt(t *testing.T) {
	var (
		primZero    = int64(0)
		primNonZero = int64(123)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullInt(primZero)
		nonZero     = gt.NullInt(primNonZero)
		dec         = new(gt.NullInt)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)

	t.Run(`Decodable/sql.Scanner`, func(t *testing.T) {
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

// TODO: test various invalid inputs.
// TODO: more tests for encoding and decoding.
func TestInterval(t *testing.T) {
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
func TestNullInterval(t *testing.T) {
	var (
		primZero    = ``
		primNonZero = `P1Y2M3DT4H5M6S`
		textZero    = ``
		textNonZero = primNonZero
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullInterval{}
		nonZero     = gt.NullIntervalFrom(1, 2, 3, 4, 5, 6)
		dec         = new(gt.NullInterval)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

func TestNullString(t *testing.T) {
	var (
		primZero    = string(``)
		primNonZero = string(`123`)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = nullBytes
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
func TestNullTime(t *testing.T) {
	var (
		primZero    = time.Time{}
		primNonZero = time.Date(1234, 5, 6, 0, 0, 0, 0, time.UTC)
		textZero    = ``
		textNonZero = `1234-05-06T00:00:00Z`
		jsonZero    = nullBytes
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
func TestNullUint(t *testing.T) {
	var (
		primZero    = uint64(0)
		primNonZero = uint64(123)
		textZero    = ``
		textNonZero = `123`
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(primNonZero)
		zero        = gt.NullUint(primZero)
		nonZero     = gt.NullUint(primNonZero)
		dec         = new(gt.NullUint)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)

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

// TODO: test various invalid inputs.
func TestNullUuid(t *testing.T) {
	var (
		primZero    = ``
		primNonZero = [gt.UuidLen]byte{0xae, 0x68, 0xcc, 0xca, 0x87, 0xc3, 0x44, 0xaf, 0xa8, 0xa0, 0x20, 0x9c, 0xe, 0x20, 0x53, 0x43}
		textZero    = ``
		textNonZero = `ae68ccca87c344afa8a0209c0e205343`
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullUuid{}
		nonZero     = gt.ParseNullUuid(textNonZero)
		dec         = new(gt.NullUuid)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

func TestParseNullUuid(t *testing.T) {
	eq(``, gt.ParseNullUuid(``).String())
	eq(`ddf1bfce018c4bef898ba4f293946049`, gt.ParseNullUuid(`ddf1bfce018c4bef898ba4f293946049`).String())
	eq(`ddf1bfce018c4bef898ba4f293946049`, gt.ParseNullUuid(`ddf1bfce-018c-4bef-898b-a4f293946049`).String())
}

// TODO: test versioning.
func TestRandomNullUuid(t *testing.T) {
	eq(false, gt.RandomNullUuid().IsZero())
	eq(false, gt.RandomNullUuid().IsZero())
	neq(gt.RandomNullUuid(), gt.RandomNullUuid())
	neq(gt.RandomNullUuid(), gt.RandomNullUuid())
}

// TODO: test various invalid inputs.
func TestUuid(t *testing.T) {
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

func testAny(
	t *testing.T,
	primZero, primNonZero interface{},
	textZero, textNonZero string,
	jsonZero, jsonNonZero []byte,
	zero, nonZero gt.Encodable,
	dec EncodableDecodable,
) {
	// Note: `primZero` doesn't have to be Go zero.
	t.Run(`zeros`, func(t *testing.T) {
		eq(true, reflect.ValueOf(zero).IsZero())
		eq(false, reflect.ValueOf(primNonZero).IsZero())
		eq(false, reflect.ValueOf(nonZero).IsZero())
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
			rval := reflect.ValueOf(dec).Elem()
			eq(true, rval.IsZero())

			set(dec, nonZero)
			eq(false, rval.IsZero())

			setZero(dec)
			eq(true, rval.IsZero())

			set(dec, nonZero)
			dec.Zero()
			eq(true, rval.IsZero())
		})

		t.Run(`Parser`, func(t *testing.T) {
			t.Run(`empty`, func(t *testing.T) {
				set(dec, nonZero)
				eqPtr(nonZero, dec)

				if zero.IsNull() {
					try(dec.Parse(``))
					eqPtr(zero, dec)
				} else {
					fail(dec.Parse(``))
				}
			})

			t.Run(`non-empty`, func(t *testing.T) {
				setZero(dec)
				eqPtr(zero, dec)

				try(dec.Parse(textNonZero))
				eqPtr(nonZero, dec)
			})
		})

		t.Run(`encoding.TextUnmarshaler`, func(t *testing.T) {
			t.Run(`empty`, func(t *testing.T) {
				t.Run(`nil bytes`, func(t *testing.T) {
					set(dec, nonZero)
					eqPtr(nonZero, dec)

					if zero.IsNull() {
						try(dec.UnmarshalText(nil))
						eqPtr(zero, dec)
					} else {
						fail(dec.UnmarshalText(nil))
					}
				})

				t.Run(`empty bytes`, func(t *testing.T) {
					set(dec, nonZero)
					eqPtr(nonZero, dec)

					if zero.IsNull() {
						try(dec.UnmarshalText([]byte{}))
						eqPtr(zero, dec)
					} else {
						fail(dec.UnmarshalText([]byte{}))
					}
				})
			})

			t.Run(`non-empty`, func(t *testing.T) {
				setZero(dec)
				eqPtr(zero, dec)

				try(dec.UnmarshalText([]byte(textNonZero)))
				eqPtr(nonZero, dec)
			})
		})

		t.Run(`json.Unmarshaler`, func(t *testing.T) {
			t.Run(`null`, func(t *testing.T) {
				set(dec, nonZero)
				eqPtr(nonZero, dec)

				if zero.IsNull() {
					try(dec.UnmarshalJSON(nullBytes))
					eqPtr(zero, dec)
				} else {
					fail(dec.UnmarshalJSON(nullBytes))
				}
			})

			t.Run(`non-empty`, func(t *testing.T) {
				setZero(dec)
				eqPtr(zero, dec)

				try(dec.UnmarshalJSON(jsonNonZero))
				eqPtr(nonZero, dec)
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
	input interface{},
) {
	set(dec, nonZero)
	eqPtr(nonZero, dec)

	if zero.IsNull() {
		try(dec.Scan(input))
		eqPtr(zero, dec)
	} else {
		fail(dec.Scan(input))
	}
}

func testScanNonEmpty(
	t *testing.T,
	zero, nonZero gt.Encodable,
	dec EncodableDecodable,
	input interface{},
) {
	setZero(dec)
	eqPtr(zero, dec)

	try(dec.Scan(input))
	eqPtr(nonZero, dec)
}

func Benchmark_uuid_parse_simple(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullUuid(`a915f35f0a3344bc8b9fb36bb650708d`)
	}
}

func Benchmark_uuid_parse_canon(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullUuid(`c230ed9a-e855-469c-8ebb-59c565aacaa7`)
	}
}

func Benchmark_uuid_string(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.String()
	}
}

func Benchmark_date_parse(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullDate(`1234-05-06`)
	}
}

func Benchmark_date_string(b *testing.B) {
	val := gt.ParseNullDate(`1234-05-06`)

	for range counter(b.N) {
		_ = val.String()
	}
}

func Benchmark_interval_parse(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)
	}
}

func Benchmark_interval_string(b *testing.B) {
	val := gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)

	for range counter(b.N) {
		_ = val.String()
	}
}
