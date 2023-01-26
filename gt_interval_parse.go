package gt

import "io"

/*
Moved to a separate file due to length, to keep the main file browsable.
Equivalent to the following regexp but faster:

	^P(?:(-?\d+)Y)?(?:(-?\d+)M)?(?:(-?\d+)D)?(?:T(?:(-?\d+)H)?(?:(-?\d+)M)?(?:(-?\d+)S)?)?$

Benchmarks at the time of writing (Go 1.17):

	regexp:

		766.1 ns/op  224 B/op  2 allocs/op

	manual:

		36.76 ns/op  0 B/op  0 allocs/op
*/
func (self *Interval) parse(src string) (err error) {
	defer errParse(&err, src, `interval`)
	defer rec(&err)

	var buf Interval
	var pos int
	var num int

	if !(pos < len(src)) {
		panic(io.EOF)
	}

	switch src[pos] {
	case 'P':
		pos++
		goto years
	default:
		panic(errFormatMismatch)
	}

years:
	{
		if !(pos < len(src)) {
			goto done
		}
		if src[pos] == 'T' {
			pos++
			goto hours
		}

		num, pos = popPrefixInt(src, pos)
		if !(pos < len(src)) {
			panic(io.EOF)
		}

		switch src[pos] {
		case 'Y':
			pos++
			buf.Years = num
			goto months
		case 'M':
			pos++
			buf.Months = num
			goto days
		case 'D':
			pos++
			buf.Days = num
			goto time
		default:
			panic(errFormatMismatch)
		}
	}

months:
	{
		if !(pos < len(src)) {
			goto done
		}
		if src[pos] == 'T' {
			pos++
			goto hours
		}

		num, pos = popPrefixInt(src, pos)
		if !(pos < len(src)) {
			panic(io.EOF)
		}

		switch src[pos] {
		case 'M':
			pos++
			buf.Months = num
			goto days
		case 'D':
			pos++
			buf.Days = num
			goto time
		default:
			panic(errFormatMismatch)
		}
	}

days:
	{
		if !(pos < len(src)) {
			goto done
		}
		if src[pos] == 'T' {
			pos++
			goto hours
		}

		num, pos = popPrefixInt(src, pos)
		if !(pos < len(src)) {
			panic(io.EOF)
		}

		switch src[pos] {
		case 'D':
			pos++
			buf.Days = num
			goto time
		default:
			panic(errFormatMismatch)
		}
	}

time:
	if !(pos < len(src)) {
		goto done
	}
	if src[pos] == 'T' {
		pos++
		goto hours
	}
	panic(errFormatMismatch)

hours:
	{
		if !(pos < len(src)) {
			goto done
		}
		num, pos = popPrefixInt(src, pos)
		if !(pos < len(src)) {
			panic(io.EOF)
		}

		switch src[pos] {
		case 'H':
			pos++
			buf.Hours = num
			goto minutes
		case 'M':
			pos++
			buf.Minutes = num
			goto seconds
		case 'S':
			pos++
			buf.Seconds = num
			goto eof
		default:
			panic(errFormatMismatch)
		}
	}

minutes:
	{
		if !(pos < len(src)) {
			goto done
		}
		num, pos = popPrefixInt(src, pos)
		if !(pos < len(src)) {
			panic(io.EOF)
		}

		switch src[pos] {
		case 'M':
			pos++
			buf.Minutes = num
			goto seconds
		case 'S':
			pos++
			buf.Seconds = num
			goto eof
		default:
			panic(errFormatMismatch)
		}
	}

seconds:
	{
		if !(pos < len(src)) {
			goto done
		}
		num, pos = popPrefixInt(src, pos)
		if !(pos < len(src)) {
			panic(io.EOF)
		}

		switch src[pos] {
		case 'S':
			pos++
			buf.Seconds = num
			goto eof
		default:
			panic(errFormatMismatch)
		}
	}

eof:
	if pos < len(src) {
		panic(errFormatMismatch)
	}

done:
	*self = buf
	return nil
}

func popPrefixInt(src string, pos int) (int, int) {
	var sig int
	var num int

	if pos < len(src) && src[pos] == '-' {
		pos++
		sig = -1
	} else {
		sig = 1
	}

	if pos < len(src) && charsetDigitDec.has(src[pos]) {
		num = undigit(src[pos])
		pos++
	} else {
		panic(errDigitEof)
	}

	for pos < len(src) && charsetDigitDec.has(src[pos]) {
		num = inc(num, src[pos])
		pos++
	}

	return (sig * num), pos
}

// Short for "increment".
func inc(num int, char byte) int { return (num * 10) + undigit(char) }

func undigit(char byte) int { return int(char - '0') }
