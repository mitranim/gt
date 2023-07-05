package gt

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"path"
)

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullUrl(src string) (val NullUrl) {
	try(val.Parse(src))
	return
}

// Safe cast. Like `gt.NullUrl(*src)` but doesn't panic on nil pointer.
func ToNullUrl(src *url.URL) (val NullUrl) {
	val.Set(src)
	return
}

/*
Variant of `*url.URL` with a less-atrocious API.
Differences from `*url.URL`:

	* Used by value, not by pointer.
	* Full support for text, JSON, SQL encoding/decoding.
	* Zero value is considered empty in text, and null in JSON and SQL.
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
func (self *NullUrl) GetPtr() any { return (*url.URL)(self) }

/*
Implement `gt.Getter`. If zero, returns `nil`, otherwise uses `.String` to
return a string representation.
*/
func (self NullUrl) Get() any {
	if self.IsNull() {
		return nil
	}
	return self.String()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullUrl) Set(src any) { try(self.Scan(src)) }

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
	return self.Url().String()
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

// Implement `gt.AppenderTo`, using the same representation as `.String`.
func (self NullUrl) AppendTo(buf []byte) []byte {
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
	if self.IsNull() {
		return nil, nil
	}
	return self.AppendTo(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullUrl) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise returns bytes representing a JSON string with the same text as in
`.String`.
*/
func (self NullUrl) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return json.Marshal(self.String())
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise parses a JSON string, using the same algorithm
as `.Parse`.
*/
func (self *NullUrl) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}

	if isJsonStr(src) {
		return self.UnmarshalText(cutJsonStr(src))
	}

	return errJsonString(src, self)
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
func (self *NullUrl) Scan(src any) error {
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
		val, ok := get(src)
		if ok {
			return self.Scan(val)
		}
		return errScanType(self, src)
	}
}

// Converts to `*url.URL`. The returned pointer refers to new memory.
func (self NullUrl) Url() *url.URL { return (*url.URL)(&self) }

// Free cast to `*url.URL`.
func (self *NullUrl) UrlPtr() *url.URL { return (*url.URL)(self) }

// If zero, returns nil. Otherwise returns a non-nil `*url.URL`.
func (self NullUrl) Maybe() *url.URL {
	if self.IsNull() {
		return nil
	}
	return self.Url()
}

/*
Returns a modified variant where `.Path` is replaced by combining the segments
via `gt.Join`. See the docs on `gt.Join`. Also see `.AddPath` that appends to
the path instead of replacing it.
*/
func (self NullUrl) WithPath(vals ...string) NullUrl {
	self.Path = Join(vals...)
	return self
}

/*
Returns a modified variant where `.Path` is replaced by combining the existing
path with the segments via `gt.Join`. See the docs on `gt.Join`. Also see
`.WithPath` that replaces the path instead of appending.
*/
func (self NullUrl) AddPath(vals ...string) NullUrl {
	// Suboptimal, TODO tune.
	self.Path = Join(self.Path, Join(vals...))
	return self
}

// Returns a modified variant with replaced `.RawQuery`.
func (self NullUrl) WithRawQuery(val string) NullUrl {
	self.RawQuery = val
	return self
}

// Returns a modified variant with replaced `.RawQuery` encoded from input.
func (self NullUrl) WithQuery(val url.Values) NullUrl {
	return self.WithRawQuery(val.Encode())
}

// Returns a modified variant with replaced `.Fragment`.
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

// `gt.NullUrl` version of `(*url.URL).EscapedPath`.
func (self NullUrl) EscapedPath() string { return self.Url().EscapedPath() }

// `gt.NullUrl` version of `(*url.URL).EscapedFragment`.
func (self NullUrl) EscapedFragment() string { return self.Url().EscapedFragment() }

// `gt.NullUrl` version of `(*url.URL).Redacted`.
func (self NullUrl) Redacted() string { return self.Url().Redacted() }

// `gt.NullUrl` version of `(*url.URL).IsAbs`.
func (self NullUrl) IsAbs() bool { return self.Url().IsAbs() }

// `gt.NullUrl` version of `(*url.URL).Parse`.
func (self NullUrl) ParseIn(ref string) (NullUrl, error) {
	val, err := self.Url().Parse(ref)
	return ToNullUrl(val), err
}

// `gt.NullUrl` version of `(*url.URL).ResolveReference`.
func (self NullUrl) ResolveReference(ref NullUrl) NullUrl {
	return ToNullUrl(self.Url().ResolveReference(ref.Url()))
}

// `gt.NullUrl` version of `(*url.URL).Query`.
func (self NullUrl) Query() url.Values { return self.Url().Query() }

// `gt.NullUrl` version of `(*url.URL).RequestURI`.
func (self NullUrl) RequestURI() string { return self.Url().RequestURI() }

// `gt.NullUrl` version of `(*url.URL).Hostname`.
func (self NullUrl) Hostname() string { return self.Url().Hostname() }

// `gt.NullUrl` version of `(*url.URL).Port`.
func (self NullUrl) Port() string { return self.Url().Port() }

/*
Like `path.Join` but with safeguards. Used internally by `gr.NullUrl.WithPath`,
exported because it may be useful separately. Differences from `path.Join`:

	* More efficient if there's only 1 segment.

	* Panics if len > 1 and any segment is "".

	* Panics if any segment begins with ".." or "/..".

Combining segments of a URL path is usually done when building a URL for a
request. Accidentally calling the wrong endpoint can have consequences much
more annoying than a panic during request building.
*/
func Join(vals ...string) string {
	switch len(vals) {
	case 0:
		return ``

	case 1:
		val := vals[0]

		// `path.Clean` would return "." in this case.
		if val == `` {
			return val
		}

		noRelativeSegment(val)
		return path.Clean(val)

	default:
		for _, val := range vals {
			noEmptySegment(val)
			noRelativeSegment(val)
		}
		return path.Join(vals...)
	}
}
