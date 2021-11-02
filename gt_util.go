package gt

import (
	"bytes"
	"database/sql"
	"encoding"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

type zeroerParser interface {
	Zeroer
	Parser
}

type nullableAppender interface {
	Nullable
	Appender
}

type nullableGetter interface {
	Nullable
	Getter
}

type zeroerPtrGetter interface {
	Zeroer
	PtrGetter
}

type zeroerTextUnmarshaler interface {
	Zeroer
	encoding.TextUnmarshaler
}

const (
	timeFormat   = time.RFC3339
	dateFormat   = "2006-01-02"
	zeroInterval = `PT0S`
	UuidLen      = 16
	UuidStrLen   = UuidLen * 2

	/**
	Should be enough for any date <= year 9999, which is the most common use case.
	TODO compute length dynamically, like we do for intervals.
	*/
	dateStrLen = len(dateFormat)
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
func bytesToMutableString(input []byte) string {
	// return string(input)
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

func errScanType(tar, inp interface{}) error {
	return fmt.Errorf(`[gt] unrecognized input for type %T: type %T, value %v`, tar, inp, inp)
}

// Original must be passed by pointer to avoid copying.
func nullGet(isNull bool, enc Getter) interface{} {
	if isNull {
		return nil
	}
	return enc.Get()
}

// Original must be passed by pointer to avoid copying.
func nullStringAppend(isNull bool, enc Appender) string {
	if isNull {
		return ``
	}
	return bytesToMutableString(enc.Append(nil))
}

// Original must be passed by pointer to avoid copying.
func nullStringStringer(isNull bool, enc fmt.Stringer) string {
	if isNull {
		return ``
	}
	return enc.String()
}

// Originals must be passed by pointers.
func nullParse(src string, zeroer Zeroer, parser Parser) error {
	if len(src) == 0 {
		zeroer.Zero()
		return nil
	}
	return parser.Parse(src)
}

// Original should be passed by pointer to avoid copying.
func nullAppend(buf []byte, isNull bool, enc Appender) []byte {
	if isNull {
		return buf
	}
	return enc.Append(buf)
}

// Original should be passed by pointer to avoid copying.
func nullNilAppend(enc nullableAppender) []byte {
	return nullAppend(nil, enc.IsNull(), enc)
}

// Original should be passed by pointer to avoid copying.
func nullTextMarshal(isNull bool, enc encoding.TextMarshaler) ([]byte, error) {
	if isNull {
		return nil, nil
	}
	return enc.MarshalText()
}

// Original should be passed by pointer to avoid copying.
func nullTextUnmarshalParser(src []byte, dec zeroerParser) error {
	if len(src) == 0 {
		dec.Zero()
		return nil
	}
	return dec.Parse(bytesToMutableString(src))
}

// Original should be passed by pointer to avoid copying.
func nullJsonMarshal(isNull bool, enc json.Marshaler) ([]byte, error) {
	if isNull {
		return bytesNull, nil
	}
	return enc.MarshalJSON()
}

// Original should be passed by pointer to avoid copying.
func nullJsonMarshalGetter(enc nullableGetter) ([]byte, error) {
	if enc == nil || enc.IsNull() {
		return bytesNull, nil
	}
	return json.Marshal(enc.Get())
}

// Original must be passed by pointer.
func nullJsonUnmarshalGetter(src []byte, dec zeroerPtrGetter) error {
	if isJsonEmpty(src) {
		dec.Zero()
		return nil
	}
	return json.Unmarshal(src, dec.GetPtr())
}

func jsonUnmarshalString(src []byte, dec encoding.TextUnmarshaler) error {
	if isJsonStr(src) {
		return dec.UnmarshalText(cutJsonStr(src))
	}
	return fmt.Errorf(`[gt] can't decode %q into %T: expected string`, src, dec)
}

func nullJsonUnmarshalString(src []byte, dec zeroerTextUnmarshaler) error {
	if isJsonEmpty(src) {
		dec.Zero()
		return nil
	}
	return jsonUnmarshalString(src, dec)
}

func scanGetter(src interface{}, tar sql.Scanner) (bool, error) {
	impl, _ := src.(Getter)
	if impl != nil {
		val := impl.Get()
		if val != src {
			return true, tar.Scan(val)
		}
	}
	return false, nil
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

func errParse(ptr *error, src string, typ string) {
	if *ptr != nil {
		*ptr = fmt.Errorf(`[gt] failed to parse %q into %v: %w`, src, typ, *ptr)
	}
}

func errInvalidChar(src string, i int) error {
	for _, char := range src[i:] {
		return fmt.Errorf(`[gt] invalid character %q in position %v`, char, i)
	}
	return fmt.Errorf(`[gt] invalid character`)
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
var decBools = [255]bool{
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
}

func isDecDigit(val byte) bool { return decBools[val] }

func prefixNumLen(src string) int {
	var hasDigit bool

	for i := 0; i < len(src); i++ {
		char := src[i]
		if i == 0 && (char == '-' || char == '+') {
			continue
		}

		if isDecDigit(char) {
			hasDigit = true
			continue
		}

		if hasDigit {
			return i
		}
	}

	return 0
}

func parseInt(src string) (val int, size int) {
	size = prefixNumLen(src)
	if size > 0 {
		var err error
		val, err = strconv.Atoi(src[:size])
		try(err)
	}
	return
}
*/

func parseIntOpt(src string) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	return strconv.Atoi(src)
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

var reInterval = regexp.MustCompile(
	`^P(?:(-?\d+)Y)?(?:(-?\d+)M)?(?:(-?\d+)D)?(?:T(?:(-?\d+)H)?(?:(-?\d+)M)?(?:(-?\d+)S)?)?$`,
)

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
