package gt_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/mitranim/gt"
)

// TODO: test various invalid inputs.
// TODO: verify that we can decode both RFC3339 timestamps and short dates.
func Test_NullDate(t *testing.T) {
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
func Test_NullFloat(t *testing.T) {
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
func Test_NullInt(t *testing.T) {
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
func Test_Interval(t *testing.T) {
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

func Test_DurationInterval(t *testing.T) {
	eq(gt.Interval{Seconds: 1}, gt.DurationInterval(time.Second))
	eq(gt.Interval{Minutes: 1}, gt.DurationInterval(time.Minute))
	eq(gt.Interval{Hours: 1}, gt.DurationInterval(time.Hour))
	eq(gt.Interval{Hours: 1, Minutes: 1, Seconds: 1}, gt.DurationInterval(time.Hour+time.Minute+time.Second))
	eq(gt.Interval{Hours: 0, Minutes: 1, Seconds: 1}, gt.DurationInterval(time.Minute+time.Second))
	eq(gt.Interval{Hours: 1, Minutes: 0, Seconds: 1}, gt.DurationInterval(time.Hour+time.Second))
	eq(gt.Interval{Hours: 1, Minutes: 1, Seconds: 0}, gt.DurationInterval(time.Hour+time.Minute))
	eq(gt.Interval{Hours: 12, Minutes: 34, Seconds: 56}, gt.DurationInterval(time.Hour*12+time.Minute*34+time.Second*56))
	eq(gt.Interval{Hours: 0, Minutes: 34, Seconds: 56}, gt.DurationInterval(time.Hour*0+time.Minute*34+time.Second*56))
	eq(gt.Interval{Hours: 12, Minutes: 0, Seconds: 56}, gt.DurationInterval(time.Hour*12+time.Minute*0+time.Second*56))
	eq(gt.Interval{Hours: 12, Minutes: 34, Seconds: 0}, gt.DurationInterval(time.Hour*12+time.Minute*34+time.Second*0))
}

// TODO: test various invalid inputs.
// TODO: more tests for encoding and decoding.
func Test_NullInterval(t *testing.T) {
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

func Test_NullString(t *testing.T) {
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
func Test_NullTime(t *testing.T) {
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
func Test_NullUint(t *testing.T) {
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
func Test_NullUrl(t *testing.T) {
	var (
		primZero    = ``
		primNonZero = `https://example.com`
		textZero    = primZero
		textNonZero = primNonZero
		jsonZero    = nullBytes
		jsonNonZero = jsonBytes(textNonZero)
		zero        = gt.NullUrl{}
		nonZero     = gt.ParseNullUrl(textNonZero)
		dec         = new(gt.NullUrl)
	)

	eq(true, zero.IsNull())
	eq(false, nonZero.IsNull())

	testAny(t, primZero, primNonZero, textZero, textNonZero, jsonZero, jsonNonZero, zero, nonZero, dec)
}

func Test_NullUrl_GoString(t *testing.T) {
	eq("gt.NullUrl{}", fmt.Sprintf(`%#v`, gt.NullUrl{}))
	eq("gt.ParseNullUrl(`one://two.three/four?five=six#seven`)", fmt.Sprintf(`%#v`, gt.ParseNullUrl(`one://two.three/four?five=six#seven`)))
}

// TODO: test various invalid inputs.
func Test_NullUuid(t *testing.T) {
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

func Test_ParseNullUuid(t *testing.T) {
	eq(``, gt.ParseNullUuid(``).String())
	eq(`ddf1bfce018c4bef898ba4f293946049`, gt.ParseNullUuid(`ddf1bfce018c4bef898ba4f293946049`).String())
	eq(`ddf1bfce018c4bef898ba4f293946049`, gt.ParseNullUuid(`ddf1bfce-018c-4bef-898b-a4f293946049`).String())
}

// TODO: test versioning.
func Test_RandomNullUuid(t *testing.T) {
	eq(false, gt.RandomNullUuid().IsZero())
	eq(false, gt.RandomNullUuid().IsZero())
	neq(gt.RandomNullUuid(), gt.RandomNullUuid())
	neq(gt.RandomNullUuid(), gt.RandomNullUuid())
}

func Test_NullUuid_GoString(t *testing.T) {
	eq("gt.NullUuid{}", fmt.Sprintf(`%#v`, gt.NullUuid{}))
	eq("gt.ParseNullUuid(`b85ae23dc3f4468995d688e1ee645501`)", fmt.Sprintf(`%#v`, gt.ParseNullUuid(`b85ae23dc3f4468995d688e1ee645501`)))
}

// TODO: test various invalid inputs.
func Test_Uuid(t *testing.T) {
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

func Test_Uuid_GoString(t *testing.T) {
	eq("gt.ParseUuid(`00000000000000000000000000000000`)", fmt.Sprintf(`%#v`, gt.Uuid{}))
	eq("gt.ParseUuid(`b85ae23dc3f4468995d688e1ee645501`)", fmt.Sprintf(`%#v`, gt.ParseUuid(`b85ae23dc3f4468995d688e1ee645501`)))
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

func Benchmark_ParseNullUuid_simple(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullUuid(`a915f35f0a3344bc8b9fb36bb650708d`)
	}
}

func Benchmark_ParseNullUuid_canon(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullUuid(`c230ed9a-e855-469c-8ebb-59c565aacaa7`)
	}
}

func Benchmark_NullUuid_string(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.String()
	}
}

func Benchmark_Uuid_GoString(b *testing.B) {
	val := gt.ParseUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.GoString()
	}
}

func Benchmark_NullUuid_GoString(b *testing.B) {
	val := gt.ParseNullUuid(`6b4c96c70bbc4d57a673de9620688f01`)

	for range counter(b.N) {
		_ = val.GoString()
	}
}

func Benchmark_ParseNullDate(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullDate(`1234-05-06`)
	}
}

func Benchmark_NullDate_String(b *testing.B) {
	val := gt.ParseNullDate(`1234-05-06`)

	for range counter(b.N) {
		_ = val.String()
	}
}

func Benchmark_ParseNullInterval(b *testing.B) {
	for range counter(b.N) {
		_ = gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)
	}
}

func Benchmark_NullInterval_String(b *testing.B) {
	val := gt.ParseNullInterval(`P12Y23M34DT45H56M67S`)

	for range counter(b.N) {
		_ = val.String()
	}
}
