package gt

import (
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

/*
Exactly the same as `HexUint(src)` but makes it clear to readers that the
conversion is done the "approved" way. A reader may be unsure about how a
regular `uint64(int64)` conversion, or vice versa, treats the sign bit,
which is: the conversion is simply a cast; the sign bit becomes the high
bit, and vice versa, with no special treatment.
*/
func Int64ToHexUint(src int64) HexUint { return HexUint(src) }

/*
Creates a random [HexUint] using [gt.ReadHexUint] and "crypto/rand".
Panics if random bytes can't be read.
*/
func RandomHexUint() HexUint {
	val, err := ReadHexUint(rand.Reader)
	try(err)
	return val
}

// Creates a [HexUint] from bytes from the provided reader.
func ReadHexUint(src io.Reader) (val HexUint, err error) {
	buf := unsafe.Slice((*byte)(unsafe.Pointer(&val)), unsafe.Sizeof(val))
	_, err = io.ReadFull(src, bufNoEscape(buf))
	if err != nil {
		err = fmt.Errorf(`[gt] unable to read random bytes for HexUint: %w`, err)
	}
	return
}

/*
Shortcut: parses successfully or panics. Should be used only in root scope.
When error handling is relevant, use `.Parse`.
*/
func ParseHexUint(src string) (val HexUint) {
	try(val.Parse(src))
	return
}

/*
Typedef of `uint64` with the following features:

  - Text format: zero is ""; otherwise: 16-character hex string.
  - JSON format: zero is `null`; otherwise: 16-character hex string.
  - SQL interop: zero is `null`; otherwise: Go `int64` / SQL "bigint".

Sometimes useful for random keys where 64-bit entropy is sufficient,
and 128-bit [Uuid] / [NullUuid] is overkill.

When interoping with SQL, this is represented with `int64`. The high bit of
`uint64` is the sign of `int64`. In SQL tables, use randomly-generated `bigint`
for primary keys, and in Go, represent them with `HexUint`.

When copy-pasting ints from SQL into Go, convert with [Int64ToHexUint].

Also see [NullInt], [NullUint], [NullUuid].
*/
type HexUint uint64

var (
	_ = Encodable(HexUint(0))
	_ = Decodable((*HexUint)(nil))
)

// Implement `gt.Zeroable`. True if 0.
func (self HexUint) IsZero() bool { return self == 0 }

// Implement `gt.Nullable`. True if 0.
func (self HexUint) IsNull() bool { return self.IsZero() }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `int64`.
func (self HexUint) Get() any {
	if self.IsNull() {
		return nil
	}
	// This conversion is simply a cast. The sign bit becomes the high bit.
	return int64(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *HexUint) Set(src any) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *HexUint) Zero() {
	if self != nil {
		*self = 0
	}
}

// Implement `fmt.Stringer`, returning a hex representation.
func (self HexUint) String() string {
	if self == 0 {
		return ``
	}
	return bytesString(self.AppendTo(nil))
}

/*
Implement `gt.Parser`. Input must be either empty or a hex string that fits
into `uint64` (≤16 characters). Hex decoding is case-insensitive.
*/
func (self *HexUint) Parse(src string) error {
	if len(src) <= 0 {
		*self = 0
		return nil
	}

	val, err := strconv.ParseUint(src, 16, 64)
	if err == nil {
		*self = HexUint(val)
		return nil
	}

	return fmt.Errorf(`[gt] unable to parse %q into HexUint: %w`, src, err)
}

// Implement `gt.AppenderTo`, using the same representation as `.String`.
func (self HexUint) AppendTo(buf []byte) []byte {
	if self == 0 {
		return buf
	}

	buf = Raw(buf).Grow(hexUintStrLen)
	buf = append(buf, hexUintZeros[self.strconvLen():]...)
	return strconv.AppendUint(buf, uint64(self), 16)
}

// Implement `encoding.TextMarhaler`, using the same representation as `.String`.
func (self HexUint) MarshalText() ([]byte, error) {
	return self.AppendTo(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *HexUint) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise returns bytes representing a JSON string with the same text
as in `.String`.
*/
func (self HexUint) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}

	var buf [hexUintStrLen + 2]byte
	size := self.strconvLen()

	buf[0] = '"'
	copy(buf[1:1+size], hexUintZeros[:])
	strconv.AppendUint(buf[1:1+hexUintStrLen-size], uint64(self), 16)
	buf[len(buf)-1] = '"'

	return buf[:], nil
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise requires a JSON string and parses it via
`.UnmarshalText`.
*/
func (self *HexUint) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	if isJsonStr(src) {
		return self.UnmarshalText(cutJsonStr(src))
	}
	return fmt.Errorf(`[gt] unable to unmarshal %q into HexUint: input must be a string`, src)
}

// Implement `driver.Valuer`, using `.Get`.
func (self HexUint) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement [sql.Scanner], converting the input to [HexUint]
and modifying the receiver. Acceptable inputs:

  - `nil`         -> use [HexUint.Zero]
  - `string`      -> use `.Parse`
  - `[]byte`      -> use `.UnmarshalText`
  - `int64`       -> convert via [Int64ToHexUint] and assign
  - `*int64`      -> use [HexUint.Zero] or convert via [Int64ToHexUint] and assign
  - `uint64`      -> convert and assign
  - `*uint64`     -> use [HexUint.Zero] or convert and assign
  - `HexUint`     -> assign
  - `gt.Getter`   -> scan underlying value
*/
func (self *HexUint) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case int64:
		*self = Int64ToHexUint(src)
		return nil

	case *int64:
		if src == nil {
			self.Zero()
		} else {
			*self = Int64ToHexUint(*src)
		}
		return nil

	case uint64:
		*self = HexUint(src)
		return nil

	case *uint64:
		if src == nil {
			self.Zero()
		} else {
			*self = HexUint(*src)
		}
		return nil

	case HexUint:
		*self = src
		return nil

	default:
		val, ok := get(src)
		if ok {
			return self.Scan(val)
		}
		return errScanType(self, src)
	}
}

// Also see `intStrLen`.
func (self HexUint) strconvLen() (out int) {
	if self == 0 {
		return 1
	}
	for self != 0 {
		self /= 16
		out++
	}
	return
}
