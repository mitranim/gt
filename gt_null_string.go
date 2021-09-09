package gt

import (
	"database/sql/driver"
)

/*
Variant of `string` where zero value is considered empty in text, and null in
JSON and SQL. Use this for fields where an empty string is not allowed, such as
enums or text foreign keys.

Unlike `string`, encoding/decoding is not always reversible:

	JSON "" → Go "" → JSON null
	SQL  '' → Go "" → SQL  null

Differences from `"database/sql".NullString`:

	* Much easier to use.
	* Supports text.
	* Supports JSON.
	* Fewer states: null and empty string are one.

In your data model, text fields should be either:

	* Non-nullable, zero value = empty string -> use `string`.
	* Nullable, zero value = `null`, empty string is not allowed -> use `gt.NullString`.

Avoid `*string` or `sql.NullString`.
*/
type NullString string

var (
	_ = Encodable(NullString(``))
	_ = Decodable((*NullString)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullString) IsZero() bool { return self == `` }

// Implement `gt.Nullable`. True if zero.
func (self NullString) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*string`.
func (self *NullString) GetPtr() interface{} { return (*string)(self) }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `string`.
func (self NullString) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return string(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullString) Set(src interface{}) { try(self.Scan(src)) }

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

// Implement `gt.Appender`, appending the string as-is.
func (self NullString) Append(buf []byte) []byte {
	return append(buf, self...)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
string as-is.
*/
func (self NullString) MarshalText() ([]byte, error) {
	return nullNilAppend(&self), nil
}

// Implement `encoding.TextUnmarshaler`, assigning the string as-is.
func (self *NullString) UnmarshalText(src []byte) error {
	*self = NullString(src)
	return nil
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `string`.
*/
func (self NullString) MarshalJSON() ([]byte, error) {
	return nullJsonMarshalGetter(&self)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*string`.
*/
func (self *NullString) UnmarshalJSON(src []byte) error {
	return nullJsonUnmarshalGetter(src, self)
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullString) Value() (driver.Value, error) { return self.Get(), nil }

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullString` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `NullString`  -> assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *NullString) Scan(src interface{}) error {
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
		ok, err := scanGetter(src, self)
		if ok || err != nil {
			return err
		}
		return errScanType(self, src)
	}
}
