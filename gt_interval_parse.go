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
	var cur int
	var num int

	if !(cur < len(src)) {
		panic(io.EOF)
	}

	switch src[cur] {
	case 'P':
		cur++
		goto years
	default:
		panic(errFormatMismatch)
	}

years:
	{
		if !(cur < len(src)) {
			goto done
		}
		if src[cur] == 'T' {
			cur++
			goto hours
		}

		num, cur = popPrefixInt(src, cur)
		if !(cur < len(src)) {
			panic(io.EOF)
		}

		switch src[cur] {
		case 'Y':
			cur++
			buf.Years = num
			goto months
		case 'M':
			cur++
			buf.Months = num
			goto days
		case 'D':
			cur++
			buf.Days = num
			goto time
		default:
			panic(errFormatMismatch)
		}
	}

months:
	{
		if !(cur < len(src)) {
			goto done
		}
		if src[cur] == 'T' {
			cur++
			goto hours
		}

		num, cur = popPrefixInt(src, cur)
		if !(cur < len(src)) {
			panic(io.EOF)
		}

		switch src[cur] {
		case 'M':
			cur++
			buf.Months = num
			goto days
		case 'D':
			cur++
			buf.Days = num
			goto time
		default:
			panic(errFormatMismatch)
		}
	}

days:
	{
		if !(cur < len(src)) {
			goto done
		}
		if src[cur] == 'T' {
			cur++
			goto hours
		}

		num, cur = popPrefixInt(src, cur)
		if !(cur < len(src)) {
			panic(io.EOF)
		}

		switch src[cur] {
		case 'D':
			cur++
			buf.Days = num
			goto time
		default:
			panic(errFormatMismatch)
		}
	}

time:
	if !(cur < len(src)) {
		goto done
	}
	if src[cur] == 'T' {
		cur++
		goto hours
	}
	panic(errFormatMismatch)

hours:
	{
		if !(cur < len(src)) {
			goto done
		}
		num, cur = popPrefixInt(src, cur)
		if !(cur < len(src)) {
			panic(io.EOF)
		}

		switch src[cur] {
		case 'H':
			cur++
			buf.Hours = num
			goto minutes
		case 'M':
			cur++
			buf.Minutes = num
			goto seconds
		case 'S':
			cur++
			buf.Seconds = num
			goto eof
		default:
			panic(errFormatMismatch)
		}
	}

minutes:
	{
		if !(cur < len(src)) {
			goto done
		}
		num, cur = popPrefixInt(src, cur)
		if !(cur < len(src)) {
			panic(io.EOF)
		}

		switch src[cur] {
		case 'M':
			cur++
			buf.Minutes = num
			goto seconds
		case 'S':
			cur++
			buf.Seconds = num
			goto eof
		default:
			panic(errFormatMismatch)
		}
	}

seconds:
	{
		if !(cur < len(src)) {
			goto done
		}
		num, cur = popPrefixInt(src, cur)
		if !(cur < len(src)) {
			panic(io.EOF)
		}

		switch src[cur] {
		case 'S':
			cur++
			buf.Seconds = num
			goto eof
		default:
			panic(errFormatMismatch)
		}
	}

eof:
	if cur < len(src) {
		panic(errFormatMismatch)
	}

done:
	*self = buf
	return nil
}

func popPrefixInt(src string, cur int) (int, int) {
	var sig int
	var num int

	if cur < len(src) && src[cur] == '-' {
		cur++
		sig = -1
	} else {
		sig = 1
	}

	if cur < len(src) && charsetDigitDec.has(src[cur]) {
		num = undigit(src[cur])
		cur++
	} else {
		panic(errDigitEof)
	}

	for cur < len(src) && charsetDigitDec.has(src[cur]) {
		num = inc(num, src[cur])
		cur++
	}

	return (sig * num), cur
}

func inc(num int, char byte) int { return (num * 10) + undigit(char) }
func undigit(char byte) int      { return int(char - '0') }
