package gt

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullUint(src string) (val NullUint) {
	try(val.Parse(src))
	return
}

/*
Variant of `uint64` where zero value is considered empty in text, and null in
JSON and SQL. Use this for fields where 0 is not allowed, such as primary and
foreign keys, or unique bigserials.

Unlike `uint64`, encoding/decoding is not always reversible:

	JSON 0 → Go 0 → JSON null
	SQL  0 → Go 0 → SQL  null

In your data model, positive numeric fields should be either:

	* Non-nullable; zero value = 0; use `uint64`.
	* Nullable; zero value = `null`; 0 is not allowed; use `gt.NullUint`.

Avoid `*uintN`.
*/
type NullUint uint64

var (
	_ = Encodable(NullUint(0))
	_ = Decodable((*NullUint)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullUint) IsZero() bool { return self == 0 }

// Implement `gt.Nullable`. True if zero.
func (self NullUint) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*uint64`.
func (self *NullUint) GetPtr() interface{} { return (*uint64)(self) }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `uint64`.
func (self NullUint) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return uint64(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullUint) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullUint) Zero() {
	if self != nil {
		*self = 0
	}
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise formats
using `strconv.FormatUint`.
*/
func (self NullUint) String() string {
	if self.IsNull() {
		return ``
	}
	return strconv.FormatUint(uint64(self), 10)
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
parses the input using `strconv.ParseUint`.
*/
func (self *NullUint) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	val, err := strconv.ParseUint(src, 10, 64)
	if err != nil {
		return err
	}

	*self = NullUint(val)
	return nil
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self NullUint) Append(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return strconv.AppendUint(buf, uint64(self), 10)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullUint) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.Append(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullUint) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `uint64`.
*/
func (self NullUint) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return json.Marshal(self.Get())
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*uint64`.
*/
func (self *NullUint) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	return json.Unmarshal(src, self.GetPtr())
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullUint) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullUint` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `uintN`       -> convert and assign
	* `*uintN`      -> use `.Zero` or convert and assign
	* `NullUint`    -> assign
	* `gt.Getter`   -> scan underlying value

TODO also support signed ints.
*/
func (self *NullUint) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case uint:
		*self = NullUint(src)
		return nil

	case *uint:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUint(*src)
		}
		return nil

	case uint8:
		*self = NullUint(src)
		return nil

	case *uint8:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUint(*src)
		}
		return nil

	case uint16:
		*self = NullUint(src)
		return nil

	case *uint16:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUint(*src)
		}
		return nil

	case uint32:
		*self = NullUint(src)
		return nil

	case *uint32:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUint(*src)
		}
		return nil

	case uint64:
		*self = NullUint(src)
		return nil

	case *uint64:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUint(*src)
		}
		return nil

	case NullUint:
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

/*
Free cast to the underlying `uint64`. Sometimes handy when this type is embedded
in a struct.
*/
func (self NullUint) Uint64() uint64 { return uint64(self) }
