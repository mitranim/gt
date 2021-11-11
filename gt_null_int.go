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
func ParseNullInt(src string) (val NullInt) {
	try(val.Parse(src))
	return
}

/*
Variant of `int64` where zero value is considered empty in text, and null in
JSON and SQL. Use this for fields where 0 is not allowed, such as primary and
foreign keys, or unique bigserials.

Unlike `int64`, encoding/decoding is not always reversible:

	JSON 0 → Go 0 → JSON null
	SQL  0 → Go 0 → SQL  null

Differences from `"database/sql".NullInt64`:

	* Much easier to use.
	* Supports text.
	* Supports JSON.
	* Fewer states: null and zero are one.

In your data model, numeric fields should be either:

	* Non-nullable; zero value = 0; use `int64`.
	* Nullable; zero value = `null`; 0 is not allowed; use `gt.NullInt`.

Avoid `*intN` or `sql.NullIntN`.
*/
type NullInt int64

var (
	_ = Encodable(NullInt(0))
	_ = Decodable((*NullInt)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullInt) IsZero() bool { return self == 0 }

// Implement `gt.Nullable`. True if zero.
func (self NullInt) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*int64`.
func (self *NullInt) GetPtr() interface{} { return (*int64)(self) }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `int64`.
func (self NullInt) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return int64(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullInt) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullInt) Zero() {
	if self != nil {
		*self = 0
	}
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise formats
using `strconv.FormatInt`.
*/
func (self NullInt) String() string {
	if self.IsNull() {
		return ``
	}
	return strconv.FormatInt(int64(self), 10)
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
parses the input using `strconv.ParseInt`.
*/
func (self *NullInt) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	val, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return err
	}

	*self = NullInt(val)
	return nil
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self NullInt) Append(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return strconv.AppendInt(buf, int64(self), 10)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullInt) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.Append(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullInt) UnmarshalText(src []byte) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `int64`.
*/
func (self NullInt) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return json.Marshal(self.Get())
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*int64`.
*/
func (self *NullInt) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	return json.Unmarshal(src, self.GetPtr())
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullInt) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullInt` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `intN`        -> convert and assign
	* `*intN`       -> use `.Zero` or convert and assign
	* `NullInt`     -> assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *NullInt) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case int:
		*self = NullInt(src)
		return nil

	case *int:
		if src == nil {
			self.Zero()
		} else {
			*self = NullInt(*src)
		}
		return nil

	case int8:
		*self = NullInt(src)
		return nil

	case *int8:
		if src == nil {
			self.Zero()
		} else {
			*self = NullInt(*src)
		}
		return nil

	case int16:
		*self = NullInt(src)
		return nil

	case *int16:
		if src == nil {
			self.Zero()
		} else {
			*self = NullInt(*src)
		}
		return nil

	case int32:
		*self = NullInt(src)
		return nil

	case *int32:
		if src == nil {
			self.Zero()
		} else {
			*self = NullInt(*src)
		}
		return nil

	case int64:
		*self = NullInt(src)
		return nil

	case *int64:
		if src == nil {
			self.Zero()
		} else {
			*self = NullInt(*src)
		}
		return nil

	case NullInt:
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
