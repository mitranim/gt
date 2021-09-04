package gt

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// `NullTime` version of `time.Date`.
func NullTimeIn(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) NullTime {
	return NullTime(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

// Shortcut for `gt.NullTimeIn` with `T00:00:00`.
func NullDateIn(year int, month time.Month, day int, loc *time.Location) NullTime {
	return NullTime(time.Date(year, month, day, 0, 0, 0, 0, loc))
}

/*
Shortcut for `gt.NullTimeIn` in UTC.

Note: due to peculiarities of civil time, `gt.NullTimeUTC(1, 1, 1, 0, 0, 0, 0)`
returns a zero value, while `gt.NullTimeUTC(0, 0, 0, 0, 0, 0, 0)` returns
a "negative" time.
*/
func NullTimeUTC(year int, month time.Month, day, hour, min, sec, nsec int) NullTime {
	return NullTimeIn(year, month, day, hour, min, sec, nsec, time.UTC)
}

/*
Shortcut for `gt.NullDateIn` in UTC.

Note: due to peculiarities of civil time, `gt.NullDateUTC(1, 1, 1)` returns a
zero value, while `gt.NullDateUTC(0, 0, 0)` returns a "negative" time.
*/
func NullDateUTC(year int, month time.Month, day int) NullTime {
	return NullDateIn(year, month, day, time.UTC)
}

// `NullTime` version of `time.Now`.
func NullTimeNow() NullTime { return NullTime(time.Now()) }

// `NullTime` version of `time.Since`.
func NullTimeSince(val NullTime) time.Duration { return NullTimeNow().Sub(val) }

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullTime(src string) (val NullTime) {
	try(val.Parse(src))
	return
}

/*
Variant of `time.Time` where zero value is considered empty in text, and null in
JSON and SQL. Prevents you from accidentally inserting nonsense times like
0001-01-01 or 1970-01-01 into date/time columns, without the hassle of dealing
with pointers such as `*time.Time` or unusable types such as `sql.NullTime`.

Differences from `time.Time`:

	* Zero value is "" in text.
	* Zero value is `null` in JSON.
	* Zero value is `null` in SQL.
	* Default text encoding is RFC3339.
	* Text encoding/decoding is automatically reversible.

Differences from `"database/sql".NullTime`:

	* Much easier to use.
	* Supports text; zero value is "".
	* Supports JSON; zero value is `null`.
	* Fewer states: null and zero are one.

In your data model, `time.Time` is often the wrong choice, because the zero
value of `time.Time` is considered "non-empty". It leads to accidentally
inserting junk data. `*time.Time` is a better choice, but it introduces nil
pointer hazards without eliminating the invalid state `&time.Time{}`.
`sql.NullTime` is unusable due to lack of support for text and JSON encoding.
`gt.NullTime` avoids all of those issues.

For civil dates without time, use `gt.NullDate`.
*/
type NullTime time.Time

var (
	_ = Encodable(NullTime{})
	_ = Decodable((*NullTime)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullTime) IsZero() bool { return self.Time().IsZero() }

// Implement `gt.Nullable`. True if zero.
func (self NullTime) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*time.Time`.
func (self *NullTime) GetPtr() interface{} { return self.TimePtr() }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `time.Time`.
func (self NullTime) Get() interface{} {
	if self.IsNull() {
		return nil
	}
	return self.Time()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullTime) Set(src interface{}) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullTime) Zero() {
	if self != nil {
		*self = NullTime{}
	}
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise returns a
text representation in the RFC3339 format.
*/
func (self NullTime) String() string {
	if self.IsNull() {
		return ``
	}
	return self.Time().Format(timeFormat)
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
requires an RFC3339 timestamp (default time parsing format in Go).
*/
func (self *NullTime) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	val, err := time.Parse(timeFormat, src)
	if err != nil {
		return err
	}

	*self = NullTime(val)
	return nil
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self NullTime) Append(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return self.Time().AppendFormat(buf, timeFormat)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullTime) MarshalText() ([]byte, error) {
	return nullNilAppend(&self), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullTime) UnmarshalText(src []byte) error {
	return nullTextUnmarshalParser(src, self)
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `time.Time`.
*/
func (self NullTime) MarshalJSON() ([]byte, error) {
	return nullJsonMarshalGetter(&self)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*time.Time`.
*/
func (self *NullTime) UnmarshalJSON(src []byte) error {
	return nullJsonUnmarshalGetter(src, self)
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullTime) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.Date` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Zero`
	* `string`      -> use `.Parse`
	* `[]byte`      -> use `.UnmarshalText`
	* `time.Time`   -> assign
	* `*time.Time`  -> use `.Zero` or assign
	* `gt.NullTime` -> assign
	* `gt.NullDate` -> assume UTC, convert, assign
	* `gt.Getter`   -> scan underlying value
*/
func (self *NullTime) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case time.Time:
		*self = NullTime(src)
		return nil

	case *time.Time:
		if src == nil {
			self.Zero()
		} else {
			*self = NullTime(*src)
		}
		return nil

	case NullTime:
		*self = src
		return nil

	case NullDate:
		*self = src.NullTimeUTC()
		return nil

	default:
		ok, err := scanGetter(src, self)
		if ok || err != nil {
			return err
		}
		return errScanType(self, src)
	}
}

// Implement `fmt.GoStringer`, returning valid Go code that constructs this value.
func (self NullTime) GoString() string {
	year, month, day := self.Date()
	hour, min, sec, nsec := self.Hour(), self.Minute(), self.Second(), self.Nanosecond()
	loc := self.Location()

	if hour == 0 && min == 0 && sec == 0 && nsec == 0 {
		if loc == time.UTC {
			return fmt.Sprintf(`gt.NullDateUTC(%v, %v, %v)`, year, int(month), day)
		}
		return fmt.Sprintf(`gt.NullDateIn(%v, %v, %v, %v)`, year, int(month), day, loc)
	}

	if loc == time.UTC {
		return fmt.Sprintf(`gt.NullTimeUTC(%v, %v, %v, %v, %v, %v, %v)`, year, int(month), day, hour, min, sec, nsec)
	}
	return fmt.Sprintf(`gt.NullTimeIn(%v, %v, %v, %v, %v, %v, %v, %v)`, year, int(month), day, hour, min, sec, nsec, loc)
}

// Free cast to `time.Time`.
func (self NullTime) Time() time.Time { return time.Time(self) }

// Free cast to `*time.Time`.
func (self *NullTime) TimePtr() *time.Time { return (*time.Time)(self) }

// If zero, returns nil. Otherwise returns a non-nil `*time.Time`.
func (self NullTime) MaybeTime() *time.Time {
	if self.IsNull() {
		return nil
	}
	return self.TimePtr()
}

// `NullTime` version of `time.Time.After`.
func (self NullTime) After(val NullTime) bool { return self.Time().After(val.Time()) }

// `NullTime` version of `time.Time.Before`.
func (self NullTime) Before(val NullTime) bool { return self.Time().Before(val.Time()) }

// `NullTime` version of `time.Time.Equal`.
func (self NullTime) Equal(val NullTime) bool { return self.Time().Equal(val.Time()) }

// `NullTime` version of `time.Time.Date`.
func (self NullTime) Date() (int, time.Month, int) { return self.Time().Date() }

// `NullTime` version of `time.Time.Year`.
func (self NullTime) Year() int { return self.Time().Year() }

// `NullTime` version of `time.Time.Month`.
func (self NullTime) Month() time.Month { return self.Time().Month() }

// `NullTime` version of `time.Time.Day`.
func (self NullTime) Day() int { return self.Time().Day() }

// `NullTime` version of `time.Time.Weekday`.
func (self NullTime) Weekday() time.Weekday { return self.Time().Weekday() }

// `NullTime` version of `time.Time.ISOWeek`.
func (self NullTime) ISOWeek() (int, int) { return self.Time().ISOWeek() }

// `NullTime` version of `time.Time.Clock`.
func (self NullTime) Clock() (int, int, int) { return self.Time().Clock() }

// `NullTime` version of `time.Time.Hour`.
func (self NullTime) Hour() int { return self.Time().Hour() }

// `NullTime` version of `time.Time.Minute`.
func (self NullTime) Minute() int { return self.Time().Minute() }

// `NullTime` version of `time.Time.Second`.
func (self NullTime) Second() int { return self.Time().Second() }

// `NullTime` version of `time.Time.Nanosecond`.
func (self NullTime) Nanosecond() int { return self.Time().Nanosecond() }

// `NullTime` version of `time.Time.YearDay`.
func (self NullTime) YearDay() int { return self.Time().YearDay() }

// `NullTime` version of `time.Time.Add`.
func (self NullTime) Add(val time.Duration) NullTime { return NullTime(self.Time().Add(val)) }

// `NullTime` version of `time.Time.Sub`.
func (self NullTime) Sub(val NullTime) time.Duration { return self.Time().Sub(val.Time()) }

// `NullTime` version of `time.Time.AddDate`.
func (self NullTime) AddDate(y, m, d int) NullTime { return NullTime(self.Time().AddDate(y, m, d)) }

// `NullTime` version of `time.Time.UTC`.
func (self NullTime) UTC() NullTime { return NullTime(self.Time().UTC()) }

// `NullTime` version of `time.Time.Local`.
func (self NullTime) Local() NullTime { return NullTime(self.Time().Local()) }

// `NullTime` version of `time.Time.In`.
func (self NullTime) In(loc *time.Location) NullTime { return NullTime(self.Time().In(loc)) }

// `NullTime` version of `time.Time.Location`.
func (self NullTime) Location() *time.Location { return self.Time().Location() }

// `NullTime` version of `time.Time.Zone`.
func (self NullTime) Zone() (string, int) { return self.Time().Zone() }

// `NullTime` version of `time.Time.Unix`.
func (self NullTime) Unix() int64 { return self.Time().Unix() }

// `NullTime` version of `time.Time.UnixMilli`.
func (self NullTime) UnixMilli() int64 { return self.Time().UnixMilli() }

// `NullTime` version of `time.Time.UnixMicro`.
func (self NullTime) UnixMicro() int64 { return self.Time().UnixMicro() }

// `NullTime` version of `time.Time.UnixNano`.
func (self NullTime) UnixNano() int64 { return self.Time().UnixNano() }

// `NullTime` version of `time.Time.IsDST`.
func (self NullTime) IsDST() bool { return self.Time().IsDST() }

// `NullTime` version of `time.Time.Truncate`.
func (self NullTime) Truncate(val time.Duration) NullTime { return NullTime(self.Time().Truncate(val)) }

// `NullTime` version of `time.Time.Round`.
func (self NullTime) Round(val time.Duration) NullTime { return NullTime(self.Time().Round(val)) }

// `NullTime` version of `time.Time.Format`.
func (self NullTime) Format(layout string) string { return self.Time().Format(layout) }

// `NullTime` version of `time.Time.AppendFormat`.
func (self NullTime) AppendFormat(a []byte, b string) []byte { return self.Time().AppendFormat(a, b) }
