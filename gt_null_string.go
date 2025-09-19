package gt

import (
	"database/sql/driver"
	"encoding/json"
)

/*
Variant of `string` where zero value is considered empty in text, and null in
JSON and SQL. Use this for fields where an empty string is not allowed, such as
enums or text foreign keys.

Unlike `string`, encoding/decoding is not always reversible:

	JSON "" → Go "" → JSON null
	SQL  '' → Go "" → SQL  null

Differences from `"database/sql".NullString`:

  - Much easier to use.
  - Supports text.
  - Supports JSON.
  - Fewer states: null and empty string are one.

In your data model, text fields should be either:

  - Non-nullable, zero value = empty string -> use `string`.
  - Nullable, zero value = `null`, empty string is not allowed -> use `gt.NullString`.

Avoid `*string` or `sql.NullString`.
*/
type NullString string

var (
	_ = Encodable(NullString(``))
	_ = Decodable((*NullString)(nil))
)

// Implement `gt.Zeroable`. True if empty.
func (self NullString) IsZero() bool { return self == `` }

// Implement `gt.Nullable`. True if empty.
func (self NullString) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*string`.
func (self *NullString) GetPtr() any { return (*string)(self) }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `string`.
func (self NullString) Get() any {
	if self.IsNull() {
		return nil
	}
	return string(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullString) Set(src any) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullString) Zero() {
	if self != nil {
		*self = ``
	}
}

// Implement `fmt.Stringer`, returning the string as-is.
func (self NullString) String() string {
	return string(self)
}

// Implement `gt.Parser`, assigning the string as-is.
func (self *NullString) Parse(src string) error {
	*self = NullString(src)
	return nil
}

// Implement `gt.AppenderTo`, appending the string as-is.
func (self NullString) AppendTo(buf []byte) []byte {
	return append(buf, self...)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
string as-is.
*/
func (self NullString) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.AppendTo(nil), nil
}

// Implement `encoding.TextUnmarshaler`, assigning the string as-is.
func (self *NullString) UnmarshalText(src []byte) error {
	// This makes a copy, which is intentional because streaming decoders tend to
	// reuse one buffer for different content.
	*self = NullString(src)
	return nil
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `string`.
*/
func (self NullString) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}

	buf := make([]byte, 0, len(self)+2)
	buf = append(buf, '"')
	buf = self.AppendTo(buf)
	buf = append(buf, '"')
	return buf, nil
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*string`.
*/
func (self *NullString) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	return json.Unmarshal(src, self.GetPtr())
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullString) Value() (driver.Value, error) { return self.Get(), nil }

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullString` and
modifying the receiver. Acceptable inputs:

  - `nil`         -> use `.Zero`
  - `string`      -> use `.Parse`
  - `[]byte`      -> use `.UnmarshalText`
  - `NullString`  -> assign
  - `gt.Getter`   -> scan underlying value
*/
func (self *NullString) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case NullString:
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

// Same as `len(self)`. Sometimes handy when embedding `gt.NullString` in
// single-field structs.
func (self NullString) Len() int { return len(self) }
