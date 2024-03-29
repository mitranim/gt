package gt

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// `gt.NullTime` version of `time.Date`.
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
Creates a UTC timestamp with the given time of day for the first day of the
Gregorian calendar.
*/
func ClockNullTime(hour, min, sec int) NullTime {
	return NullTimeUTC(0, 1, 1, hour, min, sec, 0)
}

/*
Shortcut for `gt.NullDateIn` in UTC.

Note: due to peculiarities of civil time, `gt.NullDateUTC(1, 1, 1)` returns a
zero value, while `gt.NullDateUTC(0, 0, 0)` returns a "negative" time.
*/
func NullDateUTC(year int, month time.Month, day int) NullTime {
	return NullDateIn(year, month, day, time.UTC)
}

// `gt.NullTime` version of `time.Now`.
func NullTimeNow() NullTime { return NullTime(time.Now()) }

// `gt.NullTime` version of `time.Since`.
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
True if the timestamps are ordered like this: A < B < C < ....
Also see `gt.NullTimeLessOrEqual` which uses `<=`.
*/
func NullTimeLess(val ...NullTime) bool {
	return NullTimeOrder(val, NullTime.Less)
}

/*
True if the timestamps are ordered like this: A <= B <= C <= ....
Also see `gt.NullTimeLess` which uses `<`.
*/
func NullTimeLessOrEqual(val ...NullTime) bool {
	return NullTimeOrder(val, NullTime.LessOrEqual)
}

/*
Variant of `time.Time` where zero value is considered empty in text, and null in
JSON and SQL. Prevents you from accidentally inserting nonsense times like
0001-01-01 or 1970-01-01 into date/time columns, without the hassle of pointers
such as `*time.Time` or unusable types such as `sql.NullTime`.

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

/*
Implement `gt.Zeroable`. Same as `self.Time().IsZero()`. Unlike most
implementations of `gt.Zeroable` in this package, this is NOT equivalent to
`reflect.ValueOf(self).IsZero()`, but rather a superset of it.
*/
func (self NullTime) IsZero() bool { return self.Time().IsZero() }

// Implement `gt.Nullable`. True if zero.
func (self NullTime) IsNull() bool { return self.IsZero() }

// Implement `gt.PtrGetter`, returning `*time.Time`.
func (self *NullTime) GetPtr() any { return self.TimePtr() }

// Implement `gt.Getter`. If zero, returns `nil`, otherwise returns `time.Time`.
func (self NullTime) Get() any {
	if self.IsNull() {
		return nil
	}
	return self.Time()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullTime) Set(src any) { try(self.Scan(src)) }

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
requires either an RFC3339 timestamp (default time parsing format in Go),
or a numeric timestamp in milliseconds.
*/
func (self *NullTime) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}

	if isIntString(src) {
		milli, err := strconv.ParseInt(src, 10, 64)
		*self = NullTime(time.UnixMilli(milli).In(time.UTC))
		return err
	}

	val, err := time.Parse(timeFormat, src)
	if err != nil {
		return err
	}

	*self = NullTime(val)
	return nil
}

// Implement `gt.AppenderTo`, using the same representation as `.String`.
func (self NullTime) AppendTo(buf []byte) []byte {
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
	if self.IsNull() {
		return nil, nil
	}
	return self.AppendTo(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullTime) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise uses the default `json.Marshal` behavior for `time.Time`.
*/
func (self NullTime) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return json.Marshal(self.Get())
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise uses the default `json.Unmarshal` behavior
for `*time.Time`.
*/
func (self *NullTime) UnmarshalJSON(src []byte) error {
	if isJsonEmpty(src) {
		self.Zero()
		return nil
	}
	return json.Unmarshal(src, self.GetPtr())
}

// Implement `driver.Valuer`, using `.Get`.
func (self NullTime) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullTime` and
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
func (self *NullTime) Scan(src any) error {
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
		val, ok := get(src)
		if ok {
			return self.Scan(val)
		}
		return errScanType(self, src)
	}
}

/*
Implement `fmt.GoStringer`, returning Go code that constructs this value. For
UTC, the resulting code is valid. For non-UTC, the resulting code is invalid,
because `*time.Location` doesn't implement `fmt.GoStringer`.
*/
func (self NullTime) GoString() string {
	year, month, day := self.Date()
	hour, min, sec, nsec := self.Hour(), self.Minute(), self.Second(), self.Nanosecond()
	loc := self.Location()

	if hour == 0 && min == 0 && sec == 0 && nsec == 0 {
		if loc == time.UTC {
			return fmt.Sprintf(`gt.NullDateUTC(%v, %v, %v)`, year, int(month), day)
		}
		return fmt.Sprintf(`gt.NullDateIn(%v, %v, %v, %q)`, year, int(month), day, loc)
	}

	if loc == time.UTC {
		return fmt.Sprintf(`gt.NullTimeUTC(%v, %v, %v, %v, %v, %v, %v)`, year, int(month), day, hour, min, sec, nsec)
	}
	return fmt.Sprintf(`gt.NullTimeIn(%v, %v, %v, %v, %v, %v, %v, %q)`, year, int(month), day, hour, min, sec, nsec, loc)
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

// Shortcut for `gt.NullDateFrom(self.Date())`.
func (self NullTime) NullDate() NullDate {
	return NullDateFrom(self.Date())
}

/*
Adds the interval to the time, returning the modified time. If the interval is a
zero value, the resulting time should be identical to the source.
*/
func (self NullTime) AddInterval(val Interval) NullTime {
	return self.AddDate(val.Date()).Add(val.OnlyTime().Duration())
}

/*
Subtracts the given interval, returning the modified time.
Inverse of `NullTime.AddInterval`.
*/
func (self NullTime) SubInterval(val Interval) NullTime {
	return self.AddInterval(val.Neg())
}

// Same as `gt.NullTime.AddInterval` but for `gt.NullInterval`.
func (self NullTime) AddNullInterval(val NullInterval) NullTime {
	return self.AddInterval(Interval(val))
}

// Same as `gt.NullTime.SubInterval` but for `gt.NullInterval`.
func (self NullTime) SubNullInterval(val NullInterval) NullTime {
	return self.AddNullInterval(val.Neg())
}

/*
Alias for `time.Time.Before`. Also see `gt.NullTime.Before` which is variadic.
*/
func (self NullTime) Less(val NullTime) bool {
	return self.Time().Before(val.Time())
}

// Equivalent to `self.Equal(val) || self.Less(val)`.
func (self NullTime) LessOrEqual(val NullTime) bool {
	return self.Equal(val) || self.Less(val)
}

// Same as `time.Time.Before`.
func (self NullTime) Before(val NullTime) bool {
	return self.Time().Before(val.Time())
}

// Same as `time.Time.After`.
func (self NullTime) After(val NullTime) bool {
	return self.Time().After(val.Time())
}

// `gt.NullTime` version of `time.Time.Equal`.
func (self NullTime) Equal(val NullTime) bool { return self.Time().Equal(val.Time()) }

// `gt.NullTime` version of `time.Time.Date`.
func (self NullTime) Date() (int, time.Month, int) { return self.Time().Date() }

// `gt.NullTime` version of `time.Time.Year`.
func (self NullTime) Year() int { return self.Time().Year() }

// `gt.NullTime` version of `time.Time.Month`.
func (self NullTime) Month() time.Month { return self.Time().Month() }

// `gt.NullTime` version of `time.Time.Day`.
func (self NullTime) Day() int { return self.Time().Day() }

// `gt.NullTime` version of `time.Time.Weekday`.
func (self NullTime) Weekday() time.Weekday { return self.Time().Weekday() }

// `gt.NullTime` version of `time.Time.ISOWeek`.
func (self NullTime) ISOWeek() (int, int) { return self.Time().ISOWeek() }

// `gt.NullTime` version of `time.Time.Clock`.
func (self NullTime) Clock() (int, int, int) { return self.Time().Clock() }

// `gt.NullTime` version of `time.Time.Hour`.
func (self NullTime) Hour() int { return self.Time().Hour() }

// `gt.NullTime` version of `time.Time.Minute`.
func (self NullTime) Minute() int { return self.Time().Minute() }

// `gt.NullTime` version of `time.Time.Second`.
func (self NullTime) Second() int { return self.Time().Second() }

// `gt.NullTime` version of `time.Time.Nanosecond`.
func (self NullTime) Nanosecond() int { return self.Time().Nanosecond() }

// `gt.NullTime` version of `time.Time.YearDay`.
func (self NullTime) YearDay() int { return self.Time().YearDay() }

// `gt.NullTime` version of `time.Time.Add`.
func (self NullTime) Add(val time.Duration) NullTime { return NullTime(self.Time().Add(val)) }

// `gt.NullTime` version of `time.Time.Sub`.
func (self NullTime) Sub(val NullTime) time.Duration { return self.Time().Sub(val.Time()) }

// `gt.NullTime` version of `time.Time.AddDate`.
func (self NullTime) AddDate(y, m, d int) NullTime { return NullTime(self.Time().AddDate(y, m, d)) }

// `gt.NullTime` version of `time.Time.UTC`.
func (self NullTime) UTC() NullTime { return NullTime(self.Time().UTC()) }

// `gt.NullTime` version of `time.Time.Local`.
func (self NullTime) Local() NullTime { return NullTime(self.Time().Local()) }

// `gt.NullTime` version of `time.Time.In`.
func (self NullTime) In(loc *time.Location) NullTime { return NullTime(self.Time().In(loc)) }

// `gt.NullTime` version of `time.Time.Location`.
func (self NullTime) Location() *time.Location { return self.Time().Location() }

// `gt.NullTime` version of `time.Time.Zone`.
func (self NullTime) Zone() (string, int) { return self.Time().Zone() }

// `gt.NullTime` version of `time.Time.Unix`.
func (self NullTime) Unix() int64 { return self.Time().Unix() }

// `gt.NullTime` version of `time.Time.UnixMilli`.
func (self NullTime) UnixMilli() int64 { return self.Time().UnixMilli() }

// `gt.NullTime` version of `time.Time.UnixMicro`.
func (self NullTime) UnixMicro() int64 { return self.Time().UnixMicro() }

// `gt.NullTime` version of `time.Time.UnixNano`.
func (self NullTime) UnixNano() int64 { return self.Time().UnixNano() }

// `gt.NullTime` version of `time.Time.IsDST`.
func (self NullTime) IsDST() bool { return self.Time().IsDST() }

// `gt.NullTime` version of `time.Time.Truncate`.
func (self NullTime) Truncate(val time.Duration) NullTime { return NullTime(self.Time().Truncate(val)) }

// `gt.NullTime` version of `time.Time.Round`.
func (self NullTime) Round(val time.Duration) NullTime { return NullTime(self.Time().Round(val)) }

// `gt.NullTime` version of `time.Time.Format`.
func (self NullTime) Format(layout string) string { return self.Time().Format(layout) }

// `gt.NullTime` version of `time.Time.AppendFormat`.
func (self NullTime) AppendFormat(a []byte, b string) []byte { return self.Time().AppendFormat(a, b) }

/*
True if the timestamps are ordered in such a way that the given function returns
true for each subsequent pair. If the function is nil, returns false.
Otherwise, if length is 0 or 1, returns true.
*/
func NullTimeOrder(src []NullTime, fun func(NullTime, NullTime) bool) bool {
	if fun == nil {
		return false
	}

	for len(src) > 1 {
		if !fun(src[0], src[1]) {
			return false
		}
		src = src[1:]
	}
	return true
}
