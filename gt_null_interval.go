package gt

import (
	"database/sql/driver"
	"time"
)

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseNullInterval(src string) (val NullInterval) {
	try(val.Parse(src))
	return
}

// Simplified interval constructor without a time constituent.
func DateNullInterval(years int, months int, days int) NullInterval {
	return NullInterval{Years: years, Months: months, Days: days}
}

// Simplified interval constructor without a date constituent.
func TimeNullInterval(hours, mins, secs int) NullInterval {
	return NullInterval{Hours: hours, Minutes: mins, Seconds: secs}
}

// Simplified interval constructor.
func NullIntervalFrom(years int, months int, days, hours, mins, secs int) NullInterval {
	return NullInterval{years, months, days, hours, mins, secs}
}

// Uses `.SetDuration` and returns the resulting interval.
func DurationNullInterval(src time.Duration) (val NullInterval) {
	val.SetDuration(src)
	return
}

/*
Variant of `gt.Interval` where zero value is considered empty in text, and null
in JSON and SQL.
*/
type NullInterval Interval

var (
	_ = Encodable(NullInterval{})
	_ = Decodable((*NullInterval)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self NullInterval) IsZero() bool { return Interval(self).IsZero() }

// Implement `gt.Nullable`. True if zero.
func (self NullInterval) IsNull() bool { return self.IsZero() }

/*
Implement `gt.Getter`. If zero, returns `nil`, otherwise uses `.String` to
return a string representation.
*/
func (self NullInterval) Get() any {
	if self.IsNull() {
		return nil
	}
	return Interval(self).Get()
}

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *NullInterval) Set(src any) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *NullInterval) Zero() { (*Interval)(self).Zero() }

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise returns a
text representation in the standard machine-readable ISO 8601 format.
*/
func (self NullInterval) String() string {
	if self.IsNull() {
		return ``
	}
	return Interval(self).String()
}

/*
Implement `gt.Parser`. If the input is empty, zeroes the receiver. Otherwise
requires a valid machine-readable ISO 8601 representation.
*/
func (self *NullInterval) Parse(src string) error {
	if len(src) == 0 {
		self.Zero()
		return nil
	}
	return (*Interval)(self).Parse(src)
}

// Implement `gt.AppenderTo`, using the same representation as `.String`.
func (self NullInterval) AppendTo(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return Interval(self).AppendTo(buf)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self NullInterval) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return Interval(self).MarshalText()
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *NullInterval) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise returns bytes representing a JSON string with the same text as in
`.String`.
*/
func (self NullInterval) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return bytesNull, nil
	}
	return Interval(self).MarshalJSON()
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. Otherwise parses a JSON string, using the same algorithm
as `.Parse`.
*/
func (self *NullInterval) UnmarshalJSON(src []byte) error {
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
func (self NullInterval) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.NullInterval` and
modifying the receiver. Acceptable inputs:

	* `nil`             -> use `.Zero`
	* `string`          -> use `.Parse`
	* `[]byte`          -> use `.UnmarshalText`
	* `time.Duration`   -> use `.SetDuration`
	* `*time.Duration`  -> use `.Zero` or `.SetDuration`
	* `gt.Interval`     -> assign
	* `*gt.Interval`    -> use `.Zero` or assign
	* `gt.NullInterval` -> assign
	* `gt.Getter`       -> scan underlying value
*/
func (self *NullInterval) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		self.Zero()
		return nil

	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case time.Duration:
		self.SetDuration(src)
		return nil

	case *time.Duration:
		if src == nil {
			self.Zero()
		} else {
			self.SetDuration(*src)
		}
		return nil

	case Interval:
		*self = NullInterval(src)
		return nil

	case *Interval:
		if src == nil {
			self.Zero()
		} else {
			*self = NullInterval(*src)
		}
		return nil

	case NullInterval:
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

// Same as `(*gt.Interval).SetDuration`.
func (self *NullInterval) SetDuration(val time.Duration) {
	(*Interval)(self).SetDuration(val)
}

// Same as `gt.Interval.Date`.
func (self NullInterval) Date() (years int, months int, days int) {
	return Interval(self).Date()
}

// Same as `gt.Interval.OnlyDate`.
func (self NullInterval) OnlyDate() NullInterval {
	return NullInterval(Interval(self).OnlyDate())
}

// Same as `gt.Interval.OnlyTime`.
func (self NullInterval) OnlyTime() NullInterval {
	return NullInterval(Interval(self).OnlyTime())
}

// Same as `gt.Interval.HasDate`.
func (self NullInterval) HasDate() bool {
	return Interval(self).HasDate()
}

// Same as `gt.Interval.HasTime`.
func (self NullInterval) HasTime() bool {
	return Interval(self).HasTime()
}

// Same as `gt.Interval.Duration`.
func (self NullInterval) Duration() time.Duration {
	return Interval(self).Duration()
}

// Returns a version of this interval with `.Years = val`.
func (self NullInterval) WithYears(val int) NullInterval {
	return NullInterval(Interval(self).WithYears(val))
}

// Returns a version of this interval with `.Months = val`.
func (self NullInterval) WithMonths(val int) NullInterval {
	return NullInterval(Interval(self).WithMonths(val))
}

// Returns a version of this interval with `.Days = val`.
func (self NullInterval) WithDays(val int) NullInterval {
	return NullInterval(Interval(self).WithDays(val))
}

// Returns a version of this interval with `.Hours = val`.
func (self NullInterval) WithHours(val int) NullInterval {
	return NullInterval(Interval(self).WithHours(val))
}

// Returns a version of this interval with `.Minutes = val`.
func (self NullInterval) WithMinutes(val int) NullInterval {
	return NullInterval(Interval(self).WithMinutes(val))
}

// Returns a version of this interval with `.Seconds = val`.
func (self NullInterval) WithSeconds(val int) NullInterval {
	return NullInterval(Interval(self).WithSeconds(val))
}

// Returns a version of this interval with `.Years += val`.
func (self NullInterval) AddYears(val int) NullInterval {
	return NullInterval(Interval(self).AddYears(val))
}

// Returns a version of this interval with `.Months += val`.
func (self NullInterval) AddMonths(val int) NullInterval {
	return NullInterval(Interval(self).AddMonths(val))
}

// Returns a version of this interval with `.Days += val`.
func (self NullInterval) AddDays(val int) NullInterval {
	return NullInterval(Interval(self).AddDays(val))
}

// Returns a version of this interval with `.Hours += val`.
func (self NullInterval) AddHours(val int) NullInterval {
	return NullInterval(Interval(self).AddHours(val))
}

// Returns a version of this interval with `.Minutes += val`.
func (self NullInterval) AddMinutes(val int) NullInterval {
	return NullInterval(Interval(self).AddMinutes(val))
}

// Returns a version of this interval with `.Seconds += val`.
func (self NullInterval) AddSeconds(val int) NullInterval {
	return NullInterval(Interval(self).AddSeconds(val))
}

/*
Adds every field of one interval to every field of another interval, returning
the sum. Does NOT convert different time units, such as seconds to minutes or
vice versa.
*/
func (self NullInterval) Add(val NullInterval) NullInterval {
	return NullInterval(Interval(self).Add(Interval(val)))
}

/*
Subtracts every field of one interval from every corresponding field of another
interval, returning the difference. Does NOT convert different time units, such
as seconds to minutes or vice versa.
*/
func (self NullInterval) Sub(val NullInterval) NullInterval {
	return NullInterval(Interval(self).Sub(Interval(val)))
}

/*
Returns a version of this interval with every field inverted: positive fields
become negative, and negative fields become positive.
*/
func (self NullInterval) Neg() NullInterval {
	return NullInterval(Interval(self).Neg())
}
