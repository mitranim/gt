package gt

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"io"
)

/*
Creates a random UUID using `gt.ReadUuid` and "crypto/rand". Panics if random
bytes can't be read.
*/
func RandomUuid() Uuid {
	val, err := ReadUuid(rand.Reader)
	try(err)
	return val
}

// Creates a UUID (version 4 variant 1) from bytes from the provided reader.
func ReadUuid(src io.Reader) (val Uuid, err error) {
	_, err = io.ReadFull(src, val[:])
	if err != nil {
		err = fmt.Errorf(`[gt] failed to read random bytes for UUID: %w`, err)
		return
	}

	val.setVersion()
	return
}

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseUuid(src string) (val Uuid) {
	try(val.Parse(src))
	return
}

/*
Simple UUID implementation. Features:

	* Reversible encoding/decoding in text.
	* Reversible encoding/decoding in JSON.
	* Reversible encoding/decoding in SQL.
	* Text encoding uses simplified format without dashes.
	* Text decoding supports formats with and without dashes, case-insensitive.

Differences from "github.com/google/uuid".UUID:

	* Text encoding uses simplified format without dashes.
	* Text decoding supports only simplified and canonical format.
	* Supports only version 4 (mostly-random).

When dealing with databases, it's highly recommended to use `NullUuid` instead.
*/
type Uuid [UuidLen]byte

var (
	_ = Encodable(Uuid{})
	_ = Decodable((*Uuid)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self Uuid) IsZero() bool { return self == Uuid{} }

// Implement `gt.Nullable`. Always `false`.
func (self Uuid) IsNull() bool { return false }

// Implement `gt.Getter`, returning `[16]byte` understood by many DB drivers.
func (self Uuid) Get() interface{} { return [UuidLen]byte(self) }

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *Uuid) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *Uuid) Zero() {
	if self != nil {
		*self = Uuid{}
	}
}

/*
Implement `fmt.Stringer`, returning a simplified text representation: lowercase
without dashes.
*/
func (self Uuid) String() string {
	return bytesString(self.Append(nil))
}

/*
Implement `gt.Parser`, parsing a valid UUID representation. Supports both
the short format without dashes, and the canonical format with dashes. Parsing
is case-insensitive.
*/
func (self *Uuid) Parse(src string) (err error) {
	defer errParse(&err, src, `UUID`)

	switch len(src) {
	case 32:
		return self.maybeSet(uuidParseSimple(src))
	case 36:
		return self.maybeSet(uuidParseCanon(src))
	default:
		return errUnrecLength
	}
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self Uuid) Append(buf []byte) []byte {
	buf = append(buf, uuidStrZero[:]...)
	hex.Encode(buf[len(buf)-len(uuidStrZero):], self[:])
	return buf
}

// Implement `encoding.TextMarhaler`, using the same representation as `.String`.
func (self Uuid) MarshalText() ([]byte, error) {
	return self.Append(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *Uuid) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

// Implement `json.Marshaler`, using the same representation as `.String`.
func (self Uuid) MarshalJSON() ([]byte, error) {
	var buf [UuidStrLen + 2]byte
	buf[0] = '"'
	hex.Encode(buf[1:len(buf)-1], self[:])
	buf[len(buf)-1] = '"'
	return buf[:], nil
}

// Implement `json.Unmarshaler`, using the same algorithm as `.Parse`.
func (self *Uuid) UnmarshalJSON(src []byte) error {
	if isJsonStr(src) {
		return self.UnmarshalText(cutJsonStr(src))
	}
	return errJsonString(src, self)
}

// Implement `driver.Valuer`, using `.Get`.
func (self Uuid) Value() (driver.Value, error) { return self.Get(), nil }

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.Uuid` and
modifying the receiver. Acceptable inputs:

	* `string`          -> use `.Parse`
	* `[]byte`          -> use `.UnmarshalText`
	* `[16]byte`        -> assign
	* `gt.Uuid`         -> assign
	* `gt.NullUuid`     -> assign
	* `gt.Getter`       -> scan underlying value
*/
func (self *Uuid) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case [UuidLen]byte:
		*self = Uuid(src)
		return nil

	case Uuid:
		*self = src
		return nil

	case NullUuid:
		*self = Uuid(src)
		return nil

	default:
		val, ok := get(src)
		if ok {
			return self.Scan(val)
		}
		return errScanType(self, src)
	}
}

// Equivalent to `a.String() < b.String()`. Useful for sorting.
func (self Uuid) Less(other Uuid) bool {
	for i := range self {
		if self[i] < other[i] {
			return true
		}
		if self[i] > other[i] {
			return false
		}
	}
	return false
}

// Reminder: https://en.wikipedia.org/wiki/Universally_unique_identifier
func (self *Uuid) setVersion() {
	// Version 4.
	(*self)[6] = ((*self)[6] & 0b00001111) | 0b01000000

	// Variant 1.
	(*self)[8] = ((*self)[8] & 0b00111111) | 0b10000000
}

func (self *Uuid) maybeSet(val Uuid, err error) error {
	if err == nil {
		*self = val
	}
	return err
}

func uuidParseSimple(src string) (val Uuid, err error) {
	if len(src) != 32 {
		err = errLengthMismatch
		return
	}
	_, err = hex.Decode(val[:], stringBytesUnsafe(src))
	return
}

func uuidParseCanon(src string) (val Uuid, err error) {
	if len(src) != 36 {
		err = errLengthMismatch
		return
	}

	if !(src[8] == '-' && src[13] == '-' && src[18] == '-' && src[23] == '-') {
		err = errFormatMismatch
		return
	}

	for i, pair := range uuidGroups {
		char, ok := hexDecode(src[pair[0]], src[pair[1]])
		if !ok {
			err = errInvalidCharAt(src, pair[0])
			return
		}
		val[i] = char
	}
	return
}

/*
Implement `fmt.GoStringer`, returning valid Go code that constructs this value.
The rendered code is biased for readability over performance: it parses a
string instead of using a literal constructor.
*/
func (self Uuid) GoString() string {
	const fun = `gt.ParseUuid`

	var arr [len(fun) + len("(`") + len(uuidStrZero) + len("`)")]byte

	buf := arr[:0]
	buf = append(buf, fun...)
	buf = append(buf, "(`"...)
	buf = self.Append(buf)
	buf = append(buf, "`)"...)

	return string(buf)
}
