package gt

import (
	"database/sql/driver"
	"io"
)

/*
Creates a random UUID using `gt.ReadNullUuid` and "crypto/rand". Panics if
random bytes can't be read.
*/
func RandomNullUuid() NullUuid {
	return NullUuid(RandomUuid())
}

// Creates a UUID (version 4 variant 1) from bytes from the provided reader.
func ReadNullUuid(src io.Reader) (NullUuid, error) {
	val, err := ReadUuid(src)
	return NullUuid(val), err
}

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullUuid(src string) (val NullUuid) {
	try(val.Parse(src))
	return
}

/*
Variant of `gt.Uuid` where zero value is considered empty in text, and null in
JSON and SQL. Features:

	* Reversible encoding/decoding in text. Zero value is "".
	* Reversible encoding/decoding in JSON. Zero value is `null`.
	* Reversible encoding/decoding in SQL. Zero value is `null`.
	* Text encoding uses simplified format without dashes.
	* Text decoding supports formats with and without dashes, case-insensitive.

Differences from `"github.com/google/uuid".UUID`:

	* Text encoding uses simplified format without dashes.
	* Text decoding supports only simplified and canonical format.
	* Supports only version 4 (mostly-random).
	* Zero value is considered empty in text, and null in JSON and SQL.

Differences from `"github.com/google/uuid".NullUUID`:

	* Fewer states. Null is the same as "00000000000000000000000000000000".
	* Easier to use. `NullUuid` is a type alias of `Uuid`, not a wrapper.

For database columns, `NullUuid` is recommended over `Uuid`, even when columns
are non-nullable. It prevents you from accidentally using zero-initialized
"00000000000000000000000000000000" as a value in SQL or JSON, without the hassle
of dealing with pointers or additional fields.
*/
type NullUuid Uuid

var (
	_ = Encodable(NullUuid{})
	_ = Decodable((*NullUuid)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullUuid) IsZero() bool { return self == NullUuid{} }

// Implement `gt.Nullable`. True if zero.
func (self NullUuid) IsNull() bool { return self.IsZero() }

/*
Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `[16]byte`
understood by many DB drivers.
*/
func (self NullUuid) Get() interface{} {
	return nullGet(self.IsNull(), (*Uuid)(&self))
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullUuid) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullUuid) Zero() { (*Uuid)(self).Zero() }

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise returns a
simplified text representation: lowercase without dashes.
*/
func (self NullUuid) String() string {
	return nullStringStringer(self.IsNull(), (*Uuid)(&self))
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
requires a valid UUID representation. Supports both the short format without
dashes, and the canonical format with dashes. Parsing is case-insensitive.
*/
func (self *NullUuid) Parse(src string) (err error) {
	return nullParse(src, self, (*Uuid)(self))
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self NullUuid) Append(buf []byte) []byte {
	return nullAppend(buf, self.IsNull(), (*Uuid)(&self))
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullUuid) MarshalText() ([]byte, error) {
	return nullTextMarshal(self.IsNull(), (*Uuid)(&self))
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullUuid) UnmarshalText(src []byte) error {
	return nullTextUnmarshalParser(src, self)
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise returns bytes representing a JSON string with the same text as in
`.String`.
*/
func (self NullUuid) MarshalJSON() ([]byte, error) {
	return nullJsonMarshal(self.IsNull(), (*Uuid)(&self))
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise parses a JSON string, using the same algorithm
as `.Parse`.
*/
func (self *NullUuid) UnmarshalJSON(src []byte) error {
	return nullJsonUnmarshalString(src, self)
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullUuid) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullUuid` and
modifying the receiver. Acceptable inputs:

	* `nil`             -> use `.Zero`
	* `string`          -> use `.Parse`
	* `[16]byte`        -> assign
	* `*[16]byte`       -> use `.Zero` or assign
	* `gt.Uuid`         -> assign
	* `gt.NullUuid`     -> assign
	* `gt.Getter`       -> scan underlying value
*/
func (self *NullUuid) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case [UuidLen]byte:
		*self = NullUuid(src)
		return nil

	case *[UuidLen]byte:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUuid(*src)
		}
		return nil

	case Uuid:
		*self = NullUuid(src)
		return nil

	case NullUuid:
		*self = src
		return nil

	default:
		ok, err := scanGetter(src, self)
		if ok || err != nil {
			return err
		}
		return errScanType(self, src)
	}
}

// Equivalent to `a.String() < b.String()`. Useful for sorting.
func (self NullUuid) Less(other NullUuid) bool {
	return Uuid(self).Less(Uuid(other))
}

/*
Implement `fmt.GoStringer`, returning valid Go code that constructs this value.
The rendered code is biased for readability over performance: it parses a
string instead of using a literal constructor.
*/
func (self NullUuid) GoString() string {
	if self.IsNull() {
		return `gt.NullUuid{}`
	}

	const fun = `gt.ParseNullUuid`

	var arr [len(fun) + len("(`") + len(uuidStrZero) + len("`)")]byte

	buf := arr[:0]
	buf = append(buf, fun...)
	buf = append(buf, "(`"...)
	buf = self.Append(buf)
	buf = append(buf, "`)"...)

	return string(buf)
}
