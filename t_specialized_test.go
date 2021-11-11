package gt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mitranim/gt"
)

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

// This test might fail if invoked precisely at midnight.
// That would only validate our assumptions.
func TestNullDateNow(t *testing.T) {
	eq(list(time.Now().Date()), list(gt.NullDateNow().Date()))
}

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

func TestDurationInterval(t *testing.T) {
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

func TestInterval(t *testing.T) {
	t.Run(`Parse invalid`, func(t *testing.T) {
		test := func(src string) {
			t.Helper()
			fail(new(gt.Interval).Parse(src))
		}

		test(``)
		test(` `)
		test(`1Y`)
		test(`1H`)
		test(`T1H`)
		test(`P0`)
		test(`PT0`)
		test(`P-`)
		test(`P-0`)
		test(`P-1`)
		test(`PT-`)
		test(`PT-0`)
		test(`PT-1`)
		test(`P1Y-`)
		test(`P1Y-0`)
		test(`P1Y-1`)
		test(`P1YT-`)
		test(`P1YT-0`)
		test(`P1YT-1`)
		test(`P--0Y`)
		test(`PT--0H`)
		test(`P+0Y`)
		test(`PT+0H`)
	})

	t.Run(`Parse`, func(t *testing.T) {
		test := func(exp gt.Interval, src string) {
			t.Helper()
			tar := gt.ParseInterval(src)
			if exp != tar {
				t.Logf(`failure when parsing %q`, src)
			}
			eq(exp, tar)
		}

		for _, year := range intervalParts(`Y`) {
			for _, month := range intervalParts(`M`) {
				for _, day := range intervalParts(`D`) {
					YMD := year.string + month.string + day.string

					test(
						gt.DateInterval(year.int, month.int, day.int),
						`P`+YMD,
					)

					test(
						gt.DateInterval(year.int, month.int, day.int),
						`P`+YMD+`T`,
					)

					for _, hour := range intervalParts(`H`) {
						for _, minute := range intervalParts(`M`) {
							for _, second := range intervalParts(`S`) {
								test(
									gt.IntervalFrom(year.int, month.int, day.int, hour.int, minute.int, second.int),
									`P`+YMD+`T`+hour.string+minute.string+second.string,
								)
							}
						}
					}
				}
			}
		}
	})

	t.Run(`Date`, func(t *testing.T) {
		eq(list(0, 0, 0), list(gt.Interval{}.Date()))

		eq(list(0, 0, 0), list(gt.TimeInterval(1, 2, 3).Date()))
		eq(list(0, 0, 0), list(gt.TimeInterval(-1, -2, -3).Date()))
		eq(list(1, 2, 3), list(gt.DateInterval(1, 2, 3).Date()))
		eq(list(-1, -2, -3), list(gt.DateInterval(-1, -2, -3).Date()))

		eq(
			list(1, 2, 3),
			list(gt.IntervalFrom(1, 2, 3, 4, 5, 6).Date()),
		)
	})

	t.Run(`OnlyDate`, func(t *testing.T) {
		eq(gt.Interval{}, gt.Interval{}.OnlyDate())
		eq(gt.DateInterval(0, 0, 0), gt.TimeInterval(1, 2, 3).OnlyDate())
		eq(gt.DateInterval(1, 2, 3), gt.DateInterval(1, 2, 3).OnlyDate())
		eq(gt.DateInterval(-1, -2, -3), gt.DateInterval(-1, -2, -3).OnlyDate())
	})

	t.Run(`OnlyTime`, func(t *testing.T) {
		eq(gt.Interval{}, gt.Interval{}.OnlyTime())
		eq(gt.TimeInterval(0, 0, 0), gt.DateInterval(1, 2, 3).OnlyTime())
		eq(gt.TimeInterval(1, 2, 3), gt.TimeInterval(1, 2, 3).OnlyTime())
		eq(gt.TimeInterval(-1, -2, -3), gt.TimeInterval(-1, -2, -3).OnlyTime())
	})

	// TODO also test panics in case of date constituent.
	t.Run(`Duration`, func(t *testing.T) {
		eq(time.Hour*3, gt.Interval{Hours: 3}.Duration())
		eq(time.Minute*3, gt.Interval{Minutes: 3}.Duration())
		eq(time.Second*3, gt.Interval{Seconds: 3}.Duration())

		eq(
			time.Hour*3+time.Minute*5+time.Second*7,
			gt.TimeInterval(3, 5, 7).Duration(),
		)
	})
}

func TestNullTime(t *testing.T) {
	t.Run(`Before`, func(t *testing.T) {
		eq(false, gt.NullTime{}.Before(gt.NullTime{}))
		eq(false, gt.NullDateUTC(1, 2, 3).Before(gt.NullTime{}))
		eq(false, gt.NullTime{}.Before(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 3).Before(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 4).Before(gt.NullDateUTC(1, 2, 3)))
		eq(true, gt.NullDateUTC(1, 2, 3).Before(gt.NullDateUTC(1, 2, 4)))
	})

	t.Run(`After`, func(t *testing.T) {
		eq(false, gt.NullTime{}.After(gt.NullTime{}))
		eq(false, gt.NullDateUTC(1, 2, 3).After(gt.NullTime{}))
		eq(false, gt.NullTime{}.After(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 3).After(gt.NullDateUTC(1, 2, 3)))
		eq(false, gt.NullDateUTC(1, 2, 3).After(gt.NullDateUTC(1, 2, 4)))
		eq(true, gt.NullDateUTC(1, 2, 4).After(gt.NullDateUTC(1, 2, 3)))
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
}

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

func TestNullUrl_GoString(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq(`gt.NullUrl{}`, fmt.Sprintf(`%#v`, gt.NullUrl{}))
		eq("gt.ParseNullUrl(`one://two.three/four?five=six#seven`)", fmt.Sprintf(`%#v`, gt.ParseNullUrl(`one://two.three/four?five=six#seven`)))
	})
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

func TestNullUuid(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq("gt.NullUuid{}", fmt.Sprintf(`%#v`, gt.NullUuid{}))
		eq("gt.ParseNullUuid(`b85ae23dc3f4468995d688e1ee645501`)", fmt.Sprintf(`%#v`, gt.ParseNullUuid(`b85ae23dc3f4468995d688e1ee645501`)))
	})
}

func TestUuid(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq("gt.ParseUuid(`00000000000000000000000000000000`)", fmt.Sprintf(`%#v`, gt.Uuid{}))
		eq("gt.ParseUuid(`b85ae23dc3f4468995d688e1ee645501`)", fmt.Sprintf(`%#v`, gt.ParseUuid(`b85ae23dc3f4468995d688e1ee645501`)))
	})
}

func TestTer(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq(`gt.TerNull`, fmt.Sprintf(`%#v`, gt.Ter(0)))
		eq(`gt.TerNull`, fmt.Sprintf(`%#v`, gt.TerNull))
		eq(`gt.TerFalse`, fmt.Sprintf(`%#v`, gt.TerFalse))
		eq(`gt.TerTrue`, fmt.Sprintf(`%#v`, gt.TerTrue))
		eq(`gt.Ter(3)`, fmt.Sprintf(`%#v`, gt.Ter(3)))
		eq(`gt.Ter(255)`, fmt.Sprintf(`%#v`, gt.Ter(255)))
	})
}