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
func ParseNullFloat(src string) (val NullFloat) {
	try(val.Parse(src))
	return
}

/*
Variant of `float64` where zero value is considered empty in text, and null in
JSON and SQL.

Unlike `float64`, encoding/decoding is not always reversible:

	JSON 0 → Go 0 → JSON null
	SQL  0 → Go 0 → SQL  null

Also unlike `float64`, this type doesn't use the scientific notation when
encoding to a string.

Differences from `"database/sql".NullFloat64`:

	* Much easier to use.
	* Supports text.
	* Supports JSON.
	* Fewer states: zero and null are the same.

Caution: like any floating point number, this should not be used for financial
columns. Store money as integers or use a specialized decimal type.
*/
type NullFloat float64

var (
	_ = Encodable(NullFloat(0))
	_ = Decodable((*NullFloat)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullFloat) IsZero() bool { return self == 0 }

// Implement `gt.Nullable`. True if zero.
func (self NullFloat) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*float64`.
func (self *NullFloat) GetPtr() interface{} { return (*float64)(self) }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `float64`.
func (self NullFloat) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return float64(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullFloat) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullFloat) Zero() {
	if self != nil {
		*self = 0
	}
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise formats
using `strconv.FormatFloat`.
*/
func (self NullFloat) String() string {
	if self.IsNull() {
		return ``
	}
	return strconv.FormatFloat(float64(self), 'f', -1, 64)
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
parses the input using `strconv.ParseFloat`.
*/
func (self *NullFloat) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	val, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return err
	}

	*self = NullFloat(val)
	return nil
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self NullFloat) Append(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return strconv.AppendFloat(buf, float64(self), 'f', -1, 64)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullFloat) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.Append(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullFloat) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `float64`.
*/
func (self NullFloat) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return json.Marshal(self.Get())
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*float64`.
*/
func (self *NullFloat) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	return json.Unmarshal(src, self.GetPtr())
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullFloat) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullFloat` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `intN`        -> convert and assign
	* `*intN`       -> use `.Zero` or convert and assign
	* `floatN`      -> convert and assign
	* `*floatN`     -> use `.Zero` or convert and assign
	* `NullFloat`   -> assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *NullFloat) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case int:
		*self = NullFloat(src)
		return nil

	case *int:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case int8:
		*self = NullFloat(src)
		return nil

	case *int8:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case int16:
		*self = NullFloat(src)
		return nil

	case *int16:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case int32:
		*self = NullFloat(src)
		return nil

	case *int32:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case int64:
		*self = NullFloat(src)
		return nil

	case *int64:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case float32:
		*self = NullFloat(src)
		return nil

	case *float32:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case float64:
		*self = NullFloat(src)
		return nil

	case *float64:
		if src == nil {
			self.Zero()
		} else {
			*self = NullFloat(*src)
		}
		return nil

	case NullFloat:
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
Free cast to the underlying `float64`. Sometimes handy when this type is
embedded in a struct.
*/
func (self NullFloat) Uint64() float64 { return float64(self) }
