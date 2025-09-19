package gt

import (
	"database/sql/driver"
	r "reflect"
	"strconv"
)

/*
Similar to `json.RawMessage` but supports text, JSON, SQL. In all contexts,
stores and returns self as-is, with no encoding or decoding.
*/
type Raw []byte

var (
	_ = Encodable(Raw(nil))
	_ = Decodable((*Raw)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `len(self)` <= 0.
// NOT equivalent to `self == nil` or `reflect.ValueOf(self).IsZero()`.
func (self Raw) IsZero() bool { return len(self) <= 0 }

// Implement `gt.Nullable`. True if empty.
func (self Raw) IsNull() bool { return self.IsZero() }

/*
Implement `gt.Getter`. If empty, returns `nil`, otherwise returns self as
`[]byte`.
*/
func (self Raw) Get() any {
	if self.IsNull() {
		return nil
	}
	return []byte(self)
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *Raw) Set(src any) { try(self.Scan(src)) }

/*
Implement `gt.Zeroer`, emptying the receiver. If the receiver was non-nil, its
length is reduced to 0 while keeping any capacity, and it remains non-nil.
*/
func (self *Raw) Zero() {
	if self != nil && *self != nil {
		*self = self.empty()
	}
}

// Implement `fmt.Stringer`. Returns self as-is, performing an unsafe cast.
func (self Raw) String() string { return bytesString(self) }

/*
Implement `gt.Parser`. If the input is empty, empties the receiver while keeping
any capacity. Otherwise stores the input as-is, copying it for safety. If the
receiver had enough capacity, its backing array may be mutated by this.
*/
func (self *Raw) Parse(src string) error {
	if len(src) <= 0 {
		self.Zero()
		return nil
	}
	*self = append(self.empty(), src...)
	return nil
}

// Implement `gt.AppenderTo`, appending self to the buffer as-is.
func (self Raw) AppendTo(buf []byte) []byte { return append(buf, self...) }

// Implement `encoding.TextMarhaler`, returning self as-is.
func (self Raw) MarshalText() ([]byte, error) { return self, nil }

/*
Implement `encoding.TextUnmarshaler`. If the input is empty, empties the
receiver while keeping any capacity. Otherwise stores the input as-is, copying
it for safety. If the receiver had enough capacity, its backing array may be
mutated by this.
*/
func (self *Raw) UnmarshalText(src []byte) error {
	if len(src) <= 0 {
		self.Zero()
		return nil
	}
	*self = append(self.empty(), src...)
	return nil
}

/*
Implement `json.Marshaler`. If empty, returns bytes representing `null`.
Otherwise returns self as-is, assuming that the source representation was
originally valid JSON.
*/
func (self Raw) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return self, nil
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
empties the receiver while keeping any capacity. Otherwise stores the input
as-is, copying it for safety. If the receiver had enough capacity, its backing
array may be mutated by this.
*/
func (self *Raw) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	return self.UnmarshalText(src)
}

// Implement `driver.Valuer`, using `.Get`.
func (self Raw) Value() (driver.Value, error) { return self.Get(), nil }

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.Raw` and
modifying the receiver. Acceptable inputs:

  - `nil`                   -> use `.Zero`
  - `string`                -> use `.Parse`
  - `[]byte`                -> use `.UnmarshalText`
  - convertible to `string` -> use `.Parse`
  - convertible to `[]byte` -> use `.UnmarshalText`
  - `gt.Raw`                -> assign, replacing the receiver
  - `gt.Getter`             -> scan underlying value
*/
func (self *Raw) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case Raw:
		*self = src
		return nil

	default:
		val := r.ValueOf(src)
		if val.Kind() == r.String {
			return self.Parse(val.String())
		}
		if val.Kind() == r.Slice && val.Type().Elem().Kind() == r.Uint8 {
			return self.UnmarshalText(val.Bytes())
		}

		got, ok := get(src)
		if ok {
			return self.Scan(got)
		}
		return errScanType(self, src)
	}
}

// Same as `len(self)`. Handy in edge case scenarios involving embedding.
func (self Raw) Len() int { return len(self) }

/*
Missing feature of the language / standard library. Grows the slice to ensure at
least this much additional capacity (not total capacity), returning a modified
version of the slice. The returned slice always has the same length as the
original, but its capacity and backing array may have changed. This doesn't
ensure EXACTLY the given additional capacity. It follows the usual hidden Go
rules for slice growth, and may allocate significantly more than asked. Similar
to `(*bytes.Buffer).Grow` but without wrapping, unwrapping, or spurious escapes
to the heap.
*/
func (self Raw) Grow(size int) Raw {
	len, cap := len(self), cap(self)
	if cap-len >= size {
		return self
	}

	next := make(Raw, len, 2*cap+size)
	copy(next, self)
	return next
}

/*
Implement `fmt.GoStringer`, returning valid Go code that constructs this value.
Assumes that the contents are UTF-8 text that can be represented with a Go
string.
*/
func (self Raw) GoString() string {
	if self.IsNull() {
		return `gt.Raw(nil)`
	}
	if strconv.CanBackquote(self.String()) {
		return self.goStringBackquote()
	}
	return self.goStringDoubleQuote()
}

func (self Raw) goStringBackquote() string {
	const (
		pre = "gt.Raw(`"
		suf = "`)"
	)
	buf := make([]byte, 0, len(pre)+len(self)+len(suf))
	buf = append(buf, pre...)
	buf = append(buf, self...)
	buf = append(buf, suf...)
	return bytesString(buf)
}

func (self Raw) goStringDoubleQuote() string {
	const (
		pre = `gt.Raw(`
		suf = `)`
	)

	// Guaranteed to be not enough. TODO better solution.
	buf := make([]byte, 0, len(pre)+len(self)+len(suf))

	buf = append(buf, pre...)
	buf = strconv.AppendQuoteToGraphic(buf, self.String())
	buf = append(buf, suf...)
	return bytesString(buf)
}

func (self Raw) empty() []byte { return self[:0] }
