package gt

import (
	"bytes"
	"strconv"
	"time"
	"unsafe"
)

const (
	timeFormat   = time.RFC3339
	dateFormat   = "2006-01-02"
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
	bytesNull  = []byte(`null`)
	bytesFalse = []byte(`false`)
	bytesTrue  = []byte(`true`)

	staticFalse = false
	staticTrue  = true
	ptrFalse    = &staticFalse
	ptrTrue     = &staticTrue

	uuidStrZero [UuidStrLen]byte

	charsetDigitDec = new(charset).add(`0123456789`)
)

func try(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Reasonably safe. Should not be used when the underlying
byte array is volatile, for example when it's part of a scratch buffer during
SQL scanning.
*/
func bytesString(input []byte) string {
	return *(*string)(unsafe.Pointer(&input))
}

/*
Allocation-free conversion. Returns a byte slice backed by the provided string.
Violates memory safety: the resulting value includes one word of arbitrary
memory following the original string.
*/
func stringToBytesUnsafe(input string) []byte {
	return *(*[]byte)(unsafe.Pointer(&input))
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

func hexDecode(a, b byte) (byte, bool) {
	return ((hexDigits[a] << 4) | hexDigits[b]), (hexBools[a] && hexBools[b])
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
	if val < 0 {
		out += 1
	}
	for val > 0 {
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

// This should be part of the language/standard library...
func growBytes(chunk []byte, size int) []byte {
	if cap(chunk)-len(chunk) >= size {
		return chunk
	}

	buf := bytes.NewBuffer(chunk)
	buf.Grow(size)
	return buf.Bytes()
}

func get(src interface{}) (interface{}, bool) {
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

func (self *charset) add(vals string) *charset {
	for _, val := range vals {
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
