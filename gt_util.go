package gt

import (
	"bytes"
	"strconv"
	"time"
	"unsafe"
)

const (
	timeFormat   = time.RFC3339
	dateFormat   = `2006-01-02`
	zeroInterval = `PT0S`
	UuidLen      = 16
	UuidStrLen   = UuidLen * 2

	/**
	Should be enough for any date <= year 999999.
	TODO compute length dynamically, like we do for intervals.
	*/
	dateStrLen = len(dateFormat) + 2
)

var (
	bytesNull  = stringBytesUnsafe(`null`)
	bytesFalse = stringBytesUnsafe(`false`)
	bytesTrue  = stringBytesUnsafe(`true`)

	uuidStrZero [UuidStrLen]byte

	charsetDigitDec  = new(charset).add(`0123456789`)
	charsetDigitSign = new(charset).add(`+-`)
)

func try(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Should not be used when the underlying byte array is
volatile, for example when it's part of a scratch buffer during SQL scanning.
*/
func bytesString(src []byte) string {
	return *(*string)(unsafe.Pointer(&src))
}

/*
Allocation-free conversion. Returns a byte slice backed by the provided string.
If the string is backed by read-only memory, attempting to mutate the output
may cause a segfault panic.

In Go 1.20 this can be written in a marginally "safer" way.
TODO update if we ever raise the required language version:

	unsafe.Slice(unsafe.StringData(src), len(src))
*/
func stringBytesUnsafe(src string) []byte {
	return *(*[]byte)(unsafe.Pointer(&src))
}

/*
Empty input = edge case of calling `.UnmarshalJSON` directly and passing nil or
`[]byte{}`. Not sure if we care.
*/
func isJsonEmpty(val []byte) bool {
	return len(val) == 0 || bytes.Equal(val, bytesNull)
}

func isJsonStr(val []byte) bool {
	return len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"'
}

func cutJsonStr(val []byte) []byte {
	return val[1 : len(val)-1]
}

func hexDecode(one, two byte) (byte, bool) {
	return ((hexDigits[one] << 4) | hexDigits[two]), (hexBools[one] && hexBools[two])
}

var hexDigits = [256]byte{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'a': 10,
	'b': 11,
	'c': 12,
	'd': 13,
	'e': 14,
	'f': 15,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
}

var hexBools = [256]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
}

func intStrLen(val int) (out int) {
	if val == 0 {
		return 1
	}
	if val < 0 {
		out += 1
	}
	for val != 0 {
		val /= 10
		out++
	}
	return
}

/*
See https://en.wikipedia.org/wiki/Universally_unique_identifier

Array indexes correspond to UUID bytes, values correspond to characters in a
source string encoded in the canonical format.
*/
var uuidGroups = [16][2]int{
	{0, 1}, {2, 3}, {4, 5}, {6, 7},
	{9, 10}, {11, 12},
	{14, 15}, {16, 17},
	{19, 20}, {21, 22},
	{24, 25}, {26, 27}, {28, 29}, {30, 31}, {32, 33}, {34, 35},
}

func addIntervalPartLen(ptr *int, val int) {
	if val != 0 {
		*ptr += 1 + intStrLen(val)
	}
}

func appendIntervalPart(buf []byte, val int, delim byte) []byte {
	if val == 0 {
		return buf
	}
	buf = strconv.AppendInt(buf, int64(val), 10)
	buf = append(buf, delim)
	return buf
}

func get(src any) (any, bool) {
	impl, _ := src.(Getter)
	if impl != nil {
		val := impl.Get()
		if val != src {
			return val, true
		}
	}
	return nil, false
}

type charset [256]bool

func (self *charset) has(val byte) bool { return self[val] }

func (self *charset) add(val string) *charset {
	for _, val := range val {
		self[val] = true
	}
	return self
}

// Must be deferred.
func rec(ptr *error) {
	val := recover()
	if val == nil {
		return
	}

	err, _ := val.(error)
	if err != nil {
		*ptr = err
		return
	}

	panic(val)
}

func noEmptySegment(val string) {
	if val == `` {
		panic(errEmptySegment)
	}
}

func noRelativeSegment(val string) {
	if (len(val) >= 2 && val[0] == '.' && val[1] == '.') ||
		len(val) >= 3 && val[0] == '/' && val[1] == '.' && val[2] == '.' {
		panic(errInvalidSegment(val))
	}
}

func isIntString(val string) bool {
	if len(val) == 0 {
		return false
	}

	if len(val) > 0 && charsetDigitSign.has(val[0]) {
		val = val[1:]
	}

	if len(val) == 0 {
		return false
	}
	for ind := 0; ind < len(val); ind++ {
		if !charsetDigitDec.has(val[ind]) {
			return false
		}
	}
	return true
}
