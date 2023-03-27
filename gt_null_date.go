package gt

import (
	"database/sql/driver"
	"fmt"
	"time"
)

/*
Shortcut for making a date from a time:

	inst := time.Now()
	date := gt.NullDateFrom(inst.Date())

Reversible:

	date == gt.NullDateFrom(date.Date())

Note that `gt.NullDateFrom(0, 0, 0)` returns a zero value which is considered
empty/null, but NOT equivalent to `time.Time{}`. The equivalent of zero time is
`gt.NullDateFrom(1, 1, 1)`.
*/
func NullDateFrom(year int, month time.Month, day int) NullDate {
	return NullDate{year, month, day}
}

// Shortcut for `gt.NullTimeNow().NullDate()`.
func NullDateNow() NullDate {
	return NullTimeNow().NullDate()
}

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullDate(src string) (val NullDate) {
	try(val.Parse(src))
	return
}

/*
Civil date without time. Corresponds to SQL type `date` and HTML input with
`type="date"`. Zero value is considered empty in text, and null in JSON and
SQL. Features:

	* Reversible encoding/decoding in text. Zero value is "".
	* Reversible encoding/decoding in JSON. Zero value is `null`.
	* Reversible encoding/decoding in SQL. Zero value is `null`.
	* Text encoding uses the ISO 8601 extended calendar date format: "0001-02-03".
	* Text decoding supports date-only strings and full RFC3339 timestamps.
	* Convertible to and from `gt.NullTime`.

Caution: `gt.NullDate{}` or `gt.NullDate{0, 0, 0}` is considered empty/null, but
when converted to `time.Time` or `gt.NullTime`, it's NOT equivalent to the zero
time. The equivalent of zero time is `gt.NullDate{1, 1, 1}`.
*/
type NullDate struct {
	Year  int        `json:"year"  db:"year"`
	Month time.Month `json:"month" db:"month"`
	Day   int        `json:"day"   db:"day"`
}

var (
	_ = Encodable(NullDate{})
	_ = Decodable((*NullDate)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullDate) IsZero() bool { return self == NullDate{} }

// Implement `gt.Nullable`. True if zero.
func (self NullDate) IsNull() bool { return self.IsZero() }

/*
Implement `gt.Getter`. If zero, returns `nil`, otherwise uses `.TimeUTC` to
return a timestamp suitable for SQL encoding.
*/
func (self NullDate) Get() any {
	if self.IsNull() {
		return nil
	}
	return self.TimeUTC()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullDate) Set(src any) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullDate) Zero() {
	if self != nil {
		*self = NullDate{}
	}
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise returns a
text representation in the standard machine-readable ISO 8601 format.
*/
func (self NullDate) String() string {
	if self.IsNull() {
		return ``
	}
	return bytesString(self.AppendTo(nil))
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
requires an ISO 8601 date representation, one of:

	* Extended calendar date: "2006-01-02"
	* RFC3339 (default Go timestamp format): "2006-01-02T15:04:05Z07:00"
*/
func (self *NullDate) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	var val time.Time
	var err error

	// Too restrictive. TODO fuzzier detection.
	if len(src) == len(dateFormat) {
		val, err = time.Parse(dateFormat, src)
	} else {
		val, err = time.Parse(timeFormat, src)
	}
	if err != nil {
		return err
	}

	self.SetTime(val)
	return nil
}

// Implement `gt.AppenderTo`, using the same representation as `.String`.
func (self NullDate) AppendTo(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}

	// `time.Time.AppendFormat` doesn't seem to do this.
	buf = Raw(buf).Grow(dateStrLen)

	return self.TimeUTC().AppendFormat(buf, dateFormat)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullDate) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.AppendTo(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullDate) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise returns bytes representing a JSON string with the same text as in
`.String`.
*/
func (self NullDate) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}

	var arr [dateStrLen + 2]byte
	buf := arr[:0]
	buf = append(buf, '"')
	buf = self.AppendTo(buf)
	buf = append(buf, '"')
	return buf, nil
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise parses a JSON string, using the same algorithm
as `.Parse`.
*/
func (self *NullDate) UnmarshalJSON(src []byte) error {
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
func (self NullDate) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullDate` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `time.Time`   -> use `.SetTime`
	* `*time.Time`  -> use `.Zero` or `.SetTime`
	* `gt.NullTime` -> use `.SetTime`
	* `gt.NullDate` -> assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *NullDate) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case time.Time:
		self.SetTime(src)
		return nil

	case *time.Time:
		if src == nil {
			self.Zero()
		} else {
			self.SetTime(*src)
		}
		return nil

	case NullTime:
		self.SetTime(src.Time())
		return nil

	case NullDate:
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

// Implement `fmt.GoStringer`, returning valid Go code that constructs this value.
func (self NullDate) GoString() string {
	year, month, day := self.Date()
	return fmt.Sprintf(`gt.NullDateFrom(%v, %v, %v)`, year, int(month), day)
}

/*
If the input is zero, zeroes the receiver. Otherwise uses `time.Time.Date` and
assigns the resulting year, month, day to the receiver, ignoring smaller
constituents such as hour.
*/
func (self *NullDate) SetTime(src time.Time) {
	// Note: `time.Time.Date()` "normalizes" zeros into 1 even when `.IsZero()`.
	if src.IsZero() {
		self.Zero()
	} else {
		*self = NullDateFrom(src.Date())
	}
}

// Same as `time.Time.Date`. Returns a tuple of the underlying year, month, day.
func (self NullDate) Date() (year int, month time.Month, day int) {
	return self.Year, self.Month, self.Day
}

// Converts to `gt.NullTime` with `T00:00:00` in the provided timezone.
func (self NullDate) NullTimeIn(loc *time.Location) NullTime {
	return NullTime(time.Date(self.Year, self.Month, self.Day, 0, 0, 0, 0, loc))
}

// Converts to `gt.NullTime` with `T00:00:00` in UTC.
func (self NullDate) NullTimeUTC() NullTime {
	return self.NullTimeIn(time.UTC)
}

// Converts to `time.Time` with `T00:00:00` in the provided timezone.
func (self NullDate) TimeIn(loc *time.Location) time.Time {
	return self.NullTimeIn(loc).Time()
}

// Converts to `time.Time` with `T00:00:00` in UTC.
func (self NullDate) TimeUTC() time.Time {
	return self.NullTimeUTC().Time()
}

/*
Similar to `time.Time.AddDate`. Returns a modified version of the current value,
with the year, month, day deltas added to the corresponding fields. The deltas
may be negative. Note that `time.Time` and all time-related types in this
package have a convenient `.Date` method that returns this tuple. The
calculations are performed for the UTC timezone.

As a special case, because the zero value is considered null, calling this on a
zero date ALWAYS returns the same zero date. This matches general SQL semantics
of operations involving nulls. Note that the equivalent of zero TIME is not
`gt.NullDateFrom(0, 0, 0)`, but rather `gt.NullDateFrom(1, 1, 1)`.
*/
func (self NullDate) AddDate(years int, months int, days int) NullDate {
	if self.IsZero() {
		return self
	}
	return NullDateFrom(self.NullTimeUTC().AddDate(years, months, days).Date())
}
