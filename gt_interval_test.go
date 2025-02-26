package gt_test

import (
	"testing"
	"time"

	"github.com/mitranim/gt"
)

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
	t.Run(`Parse_invalid`, func(t *testing.T) {
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
		if testing.Short() {
			t.Skip(`takes too long (half a second)`)
		}

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
