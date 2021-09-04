package gt

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
)

/*
Implemented by all types in this package. Returns the underlying value as
a "primitive" / "well-known" type, such as `int64`, `string`, `time.Time`
depending on the type. All types in this package use `.Get` to implement
`sql.Valuer`.
*/
type Getter interface{ Get() interface{} }

/*
Implemented by all types in this package. Same as `.Scan`, but panics on error.
*/
type Setter interface{ Set(interface{}) }

/*
Implemented by all types in this package, as well as some stdlib types.
Equivalent to `reflect.ValueOf(val).IsZero()`, but also works on pointer
receivers.
*/
type Zeroable interface{ IsZero() bool }

/*
Implemented by all types in this package. For all "null" types, this is
equivalent to `gt.Zeroable`. For all non-"null" types, this always returns
`false`.
*/
type Nullable interface{ IsNull() bool }

/*
Zeroes the receiver. Implemented by all types in this package, as well as some
stdlib types.
*/
type Zeroer interface{ Zero() }

/*
Missing counterpart to `encoding.TextUnmarshaler`. Parses a string, rather than
a byte slice.
*/
type Parser interface{ Parse(string) error }

/*
Implemented by all types in this package, as well as some stdlib types. Appends
the default text representation of the receiver to the provided buffer.
*/
type Appender interface{ Append([]byte) []byte }

/*
Mutable counterpart to `gt.Getter`. Where `.Get` returns an underlying primitive
value as a copy, `.GetPtr` returns an underlying primitive value as a pointer.
Decoding into the underlying value by using `json.Unmarshal`, SQL decoding, or
any other decoding mechanism must mutate the target, and the resulting state
must be valid for that type.

Unlike other interfaces, not every type in this package implements this. This is
implemented only by types whose underlying value is built-in (strings and
numbers) or also supports decoding (`time.Time`).
*/
type PtrGetter interface{ GetPtr() interface{} }

/*
Implemented by all types in this package. Various methods implemented on value
types, rather than pointer types, for converting the value to another
representation.
*/
type Encodable interface {
	Getter
	Zeroable
	Nullable
	fmt.Stringer
	encoding.TextMarshaler
	json.Marshaler
	driver.Valuer
}

/*
Implemented by all types in this package. Various methods implemented on pointer
types, rather than value types, for mutating the underlying value by decoding
or zeroing.
*/
type Decodable interface {
	Setter
	Zeroer
	Parser
	Appender
	encoding.TextUnmarshaler
	json.Unmarshaler
	sql.Scanner
}
