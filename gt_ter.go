package gt

import (
	"bytes"
	"database/sql/driver"
	"fmt"
)

/*
Valid representations of `gt.Ter`. Other values are considered invalid and will
cause panics.
*/
const (
	TerNull  Ter = 0
	TerFalse Ter = 1
	TerTrue  Ter = 2
)

/*
Shortcut: parses successfully or panics. Provided only for consistency with
other types. Prefer constants such as `gt.TerNull`.
*/
func ParseTer(src string) (val Ter) {
	try(val.Parse(src))
	return
}

/*
Converts boolean to ternary:

	* false = gt.TerFalse
	* true  = gt.TerTrue

For inverse conversion, use `gt.Ter.LaxBool` or `gt.Ter.TryBool`.
*/
func BoolTer(val bool) Ter {
	if val {
		return TerTrue
	}
	return TerFalse
}

/*
Converts boolean pointer to ternary:

	* nil    = gt.TerNull
	* &false = gt.TerFalse
	* &true  = gt.TerTrue

For inverse conversion, use `gt.Ter.BoolPtr`.
*/
func BoolPtrTer(val *bool) Ter {
	if val == nil {
		return TerNull
	}
	return BoolTer(*val)
}

/*
Ternary type / nullable boolean type. Similar to `*bool`, with various
advantages. Has three states with the following representations:

	TerNull  | 0 | ""      in text | null  in JSON | null  in SQL
	TerFalse | 1 | "false" in text | false in JSON | false in SQL
	TerTrue  | 2 | "true"  in text | true  in JSON | true  in SQL

Differences from `bool`:

	* 3 states rather than 2.
	* Nullable in JSON and SQL.
	* Zero value is empty/null rather than false.

Differences from `*bool`:

	* More efficient: 1 byte, no heap indirection, no added GC pressure.
	* Safer: no nil pointer panics.
	* Zero value is considered empty in text.
	* Text encoding/decoding is reversible.

Differences from `sql.NullBool`:

	* More efficient: 1 byte rather than 2.
	* Much easier to use.
	* Supports text.
	* Supports JSON.
*/
type Ter byte

var (
	_ = Encodable(Ter(0))
	_ = Decodable((*Ter)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self Ter) IsZero() bool { return self == TerNull }

// Implement `gt.Nullable`. True if zero.
func (self Ter) IsNull() bool { return self.IsZero() }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `bool`.
func (self Ter) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return self.LaxBool()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *Ter) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *Ter) Zero() {
	if self != nil {
		*self = TerNull
	}
}

/*
Implement `fmt.Stringer`, using the following representations:

	* gt.TerNull  = ""
	* gt.TerFalse = "false"
	* gt.TerTrue  = "true"
*/
func (self Ter) String() string {
	switch self {
	case TerNull:
		return ``
	case TerFalse:
		return `false`
	case TerTrue:
		return `true`
	default:
		panic(self.invalid())
	}
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
expects the input to be "false" or "true".
*/
func (self *Ter) Parse(src string) (err error) {
	defer errParse(&err, src, `ternary`)

	switch src {
	case ``:
		*self = TerNull
		return nil
	case `false`:
		*self = TerFalse
		return nil
	case `true`:
		*self = TerTrue
		return nil
	default:
		return fmt.Errorf(`[gt] failed to parse ternary: expected empty string, "false", or "true", got %q`, src)
	}
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self Ter) Append(buf []byte) []byte {
	return append(buf, self.String()...)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self Ter) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.Append(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *Ter) UnmarshalText(src []byte) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`, using the following representations:

	* gt.TerNull  = []byte("null")
	* gt.TerFalse = []byte("false")
	* gt.TerTrue  = []byte("true")
*/
func (self Ter) MarshalJSON() ([]byte, error) {
	switch self {
	case TerNull:
		return bytesNull, nil
	case TerFalse:
		return bytesFalse, nil
	case TerTrue:
		return bytesTrue, nil
	default:
		return nil, self.invalid()
	}
}

/*
Implement `json.Unmarshaler`, using the following representations:

	* []byte(nil)     = gt.TerNull
	* []byte("")      = gt.TerNull
	* []byte("null")  = gt.TerNull
	* []byte("false") = gt.TerFalse
	* []byte("true")  = gt.TerTrue
*/
func (self *Ter) UnmarshalJSON(src []byte) error {
	if bytes.Equal(src, bytesNull) {
		self.Zero()
		return nil
	}
	return self.UnmarshalText(src)
}

// Implement `driver.Valuer`, using `.Get`.
func (self Ter) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.Ter` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `bool`        -> use `.SetBool`
	* `*bool`       -> use `.SetBoolPtr`
	* `Ter`         -> assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *Ter) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case bool:
		self.SetBool(src)
		return nil

	case *bool:
		self.SetBoolPtr(src)
		return nil

	case Ter:
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

func (self Ter) invalid() error {
	return fmt.Errorf(`[gt] unrecognized value of %[1]T: %[1]v`, self)
}

// Sets the receiver to the result of `gt.BoolTer`.
func (self *Ter) SetBool(val bool) {
	*self = BoolTer(val)
}

// Sets the receiver to the result of `gt.BoolPtrTer`.
func (self *Ter) SetBoolPtr(val *bool) {
	*self = BoolPtrTer(val)
}

/*
Semi-inverse of `gt.BoolTer`. Permissive conversion, where anything untrue is
considered false. Equivalent to `.IsTrue()`.
*/
func (self Ter) LaxBool() bool {
	return self.IsTrue()
}

/*
Exact inverse of `gt.BoolTer`. If true or false, converts to a boolean,
otherwise panics.
*/
func (self Ter) TryBool() bool {
	switch self {
	case TerNull:
		panic(errTerNullBool)
	case TerFalse:
		return false
	case TerTrue:
		return true
	default:
		panic(self.invalid())
	}
}

/*
Inverse of `gt.BoolPtrTer`. Converts to a boolean pointer:

	* gt.TerNull  = nil
	* gt.TerFalse = &false
	* gt.TerTrue  = &true

The returned values are statically allocated and must never be modified.
*/
func (self Ter) BoolPtr() *bool {
	switch self {
	case TerNull:
		return nil
	case TerFalse:
		return ptrFalse
	case TerTrue:
		return ptrTrue
	default:
		panic(self.invalid())
	}
}

/*
Exact boolean equality. If the receiver is not true or false, this returns false
regardless of the input.
*/
func (self Ter) EqBool(val bool) bool {
	if val {
		return self.IsTrue()
	}
	return self.IsFalse()
}

// Same as `== gt.TerTrue`.
func (self Ter) IsTrue() bool { return self == TerTrue }

// Same as `== gt.TerFalse`.
func (self Ter) IsFalse() bool { return self == TerFalse }

// Implement `fmt.GoStringer`, returning valid Go code representing this value.
func (self Ter) GoString() string {
	switch self {
	case TerNull:
		return `gt.TerNull`
	case TerFalse:
		return `gt.TerFalse`
	case TerTrue:
		return `gt.TerTrue`
	default:
		return fmt.Sprintf(`gt.Ter(%v)`, byte(self))
	}
}
