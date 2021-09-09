package gt

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
)

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullUrl(src string) (val NullUrl) {
	try(val.Parse(src))
	return
}

// Safe cast. Like `gt.NullUrl(*src)` but also accepts nil.
func ToNullUrl(src *url.URL) (val NullUrl) {
	val.Set(src)
	return
}

/*
Variant of `*url.URL` with a less-atrocious API.

Differences from `*url.URL`:

	* Used by value, not by pointer.
	* Zero value is considered empty in text, and null in JSON and SQL.
	* Full support for text, JSON, SQL encoding/decoding.
	* Easier to use.
	* Fewer invalid states.
*/
type NullUrl url.URL

var (
	_ = Encodable(NullUrl{})
	_ = Decodable((*NullUrl)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullUrl) IsZero() bool { return self == NullUrl{} }

// Implement `gt.Nullable`. True if zero.
func (self NullUrl) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*url.URL`.
func (self *NullUrl) GetPtr() interface{} { return (*url.URL)(self) }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `*url.URL`.
func (self NullUrl) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return self.UrlPtr()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullUrl) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullUrl) Zero() {
	if self != nil {
		*self = NullUrl{}
	}
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise formats
using `(*url.URL).String`.
*/
func (self NullUrl) String() string {
	if self.IsNull() {
		return ``
	}
	return self.UrlPtr().String()
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
parses the input using `url.Parse`.
*/
func (self *NullUrl) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	val, err := url.Parse(src)
	if err != nil {
		return err
	}

	*self = NullUrl(*val)
	return nil
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self NullUrl) Append(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return append(buf, self.String()...)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullUrl) MarshalText() ([]byte, error) {
	return nullNilAppend(&self), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullUrl) UnmarshalText(src []byte) error {
	return nullTextUnmarshalParser(src, self)
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise returns bytes representing a JSON string with the same text as in
`.String`.
*/
func (self NullUrl) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return nullBytes, nil
	}
	return json.Marshal(self.String())
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise parses a JSON string, using the same algorithm
as `.Parse`.
*/
func (self *NullUrl) UnmarshalJSON(src []byte) error {
	return nullJsonUnmarshalString(src, self)
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullUrl) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullUrl` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `url.URL`     -> convert and assign
	* `*url.URL`    -> use `.Zero` or convert and assign
	* `NullUrl`     -> assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *NullUrl) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case url.URL:
		*self = NullUrl(src)
		return nil

	case *url.URL:
		if src == nil {
			self.Zero()
		} else {
			*self = NullUrl(*src)
		}
		return nil

	case NullUrl:
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

// Free cast to `*url.URL`.
func (self *NullUrl) UrlPtr() *url.URL { return (*url.URL)(self) }

// If zero, returns nil. Otherwise returns a non-nil `*url.URL`.
func (self NullUrl) MaybeUrl() *url.URL {
	if self.IsNull() {
		return nil
	}
	return self.UrlPtr()
}

// Returns modified variant with replaced `.Path`.
func (self NullUrl) WithPath(val string) NullUrl {
	self.Path = val
	return self
}

// Returns modified variant with replaced `.RawQuery`.
func (self NullUrl) WithRawQuery(val string) NullUrl {
	self.RawQuery = val
	return self
}

// Returns modified variant with replaced `.RawQuery` encoded from input.
func (self NullUrl) WithQuery(val url.Values) NullUrl {
	return self.WithRawQuery(val.Encode())
}

// Returns modified variant with replaced `.Fragment`.
func (self NullUrl) WithFragment(val string) NullUrl {
	self.Fragment = val
	return self
}

// Implement `fmt.GoStringer`, returning valid Go code that constructs this value.
func (self NullUrl) GoString() string {
	if self.IsNull() {
		return `gt.NullUrl{}`
	}
	return "gt.ParseNullUrl(`" + self.String() + "`)"
}

// `NullUrl` version of `*url.URL.EscapedPath`.
func (self NullUrl) EscapedPath() string { return self.UrlPtr().EscapedPath() }

// `NullUrl` version of `*url.URL.EscapedFragment`.
func (self NullUrl) EscapedFragment() string { return self.UrlPtr().EscapedFragment() }

// `NullUrl` version of `*url.URL.Redacted`.
func (self NullUrl) Redacted() string { return self.UrlPtr().Redacted() }

// `NullUrl` version of `*url.URL.IsAbs`.
func (self NullUrl) IsAbs() bool { return self.UrlPtr().IsAbs() }

// `NullUrl` version of `*url.URL.Parse`.
func (self NullUrl) ParseIn(ref string) (NullUrl, error) {
	val, err := self.UrlPtr().Parse(ref)
	return ToNullUrl(val), err
}

// `NullUrl` version of `*url.URL.ResolveReference`.
func (self NullUrl) ResolveReference(ref NullUrl) NullUrl {
	return ToNullUrl(self.UrlPtr().ResolveReference(ref.UrlPtr()))
}

// `NullUrl` version of `*url.URL.Query`.
func (self NullUrl) Query() url.Values { return self.UrlPtr().Query() }

// `NullUrl` version of `*url.URL.RequestURI`.
func (self NullUrl) RequestURI() string { return self.UrlPtr().RequestURI() }

// `NullUrl` version of `*url.URL.Hostname`.
func (self NullUrl) Hostname() string { return self.UrlPtr().Hostname() }

// `NullUrl` version of `*url.URL.Port`.
func (self NullUrl) Port() string { return self.UrlPtr().Port() }
