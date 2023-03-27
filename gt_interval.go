package gt

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

/*
Shortcut: parses successfully or panics. Should be used only in root scope. When
error handling is relevant, use `.Parse`.
*/
func ParseInterval(src string) (val Interval) {
	try(val.Parse(src))
	return
}

// Simplified interval constructor without a time constituent.
func DateInterval(years, months, days int) Interval {
	return Interval{Years: years, Months: months, Days: days}
}

// Simplified interval constructor without a date constituent.
func TimeInterval(hours, mins, secs int) Interval {
	return Interval{Hours: hours, Minutes: mins, Seconds: secs}
}

// Simplified interval constructor.
func IntervalFrom(years, months, days, hours, mins, secs int) Interval {
	return Interval{years, months, days, hours, mins, secs}
}

// Uses `.SetDuration` and returns the resulting interval.
func DurationInterval(src time.Duration) (val Interval) {
	val.SetDuration(src)
	return
}

/*
Represents an ISO 8601 time interval that has only duration (no timestamps, no
range).

Features:

	* Reversible encoding/decoding in text.
	* Reversible encoding/decoding in JSON.
	* Reversible encoding/decoding in SQL.

Text encoding and decoding uses the standard ISO 8601 format:

	P0Y0M0DT0H0M0S

When interacting with a database, to make intervals parsable, configure your DB
to always output them in the standard ISO 8601 format.

Limitations:

	* Supports only the standard machine-readable format.
	* Doesn't support decimal fractions.

For a nullable variant, see `gt.NullInterval`.
*/
type Interval struct {
	Years   int `json:"years"   db:"years"`
	Months  int `json:"months"  db:"months"`
	Days    int `json:"days"    db:"days"`
	Hours   int `json:"hours"   db:"hours"`
	Minutes int `json:"minutes" db:"minutes"`
	Seconds int `json:"seconds" db:"seconds"`
}

var (
	_ = Encodable(Interval{})
	_ = Decodable((*Interval)(nil))
)

// Implement `gt.Zeroable`. Equivalent to `reflect.ValueOf(self).IsZero()`.
func (self Interval) IsZero() bool { return self == Interval{} }

// Implement `gt.Nullable`. Always `false`.
func (self Interval) IsNull() bool { return false }

// Implement `gt.Getter`, using `.String` to return a string representation.
func (self Interval) Get() any { return self.String() }

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *Interval) Set(src any) { try(self.Scan(src)) }

// Implement `gt.Zeroer`, zeroing the receiver.
func (self *Interval) Zero() {
	if self != nil {
		*self = Interval{}
	}
}

/*
Implement `fmt.Stringer`, returning a text representation in the standard
machine-readable ISO 8601 format.
*/
func (self Interval) String() string {
	if self.IsZero() {
		return zeroInterval
	}
	return bytesString(self.AppendTo(nil))
}

// Implement `gt.Parser`, parsing a valid machine-readable ISO 8601 representation.
func (self *Interval) Parse(src string) error { return self.parse(src) }

// Implement `gt.AppenderTo`, using the same representation as `.String`.
func (self Interval) AppendTo(buf []byte) []byte {
	if self.IsZero() {
		return append(buf, zeroInterval...)
	}

	buf = Raw(buf).Grow(self.bufLen())
	buf = append(buf, 'P')
	buf = appendIntervalPart(buf, self.Years, 'Y')
	buf = appendIntervalPart(buf, self.Months, 'M')
	buf = appendIntervalPart(buf, self.Days, 'D')

	if self.HasTime() {
		buf = append(buf, 'T')
		buf = appendIntervalPart(buf, self.Hours, 'H')
		buf = appendIntervalPart(buf, self.Minutes, 'M')
		buf = appendIntervalPart(buf, self.Seconds, 'S')
	}

	return buf
}

// Implement `encoding.TextMarhaler`, using the same representation as `.String`.
func (self Interval) MarshalText() ([]byte, error) {
	return self.AppendTo(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *Interval) UnmarshalText(src []byte) error {
	return self.Parse(bytesString(src))
}

// Implement `json.Marshaler`, using the same representation as `.String`.
func (self Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Get())
}

// Implement `json.Unmarshaler`, using the same algorithm as `.Parse`.
func (self *Interval) UnmarshalJSON(src []byte) error {
	if isJsonStr(src) {
		return self.UnmarshalText(cutJsonStr(src))
	}
	return errJsonString(src, self)
}

// Implement `driver.Valuer`, using `.Get`.
func (self Interval) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `gt.Interval` and
modifying the receiver. Acceptable inputs:

	* `string`          -> use `.Parse`
	* `[]byte`          -> use `.UnmarshalText`
	* `time.Duration`   -> use `.SetDuration`
	* `gt.Interval`     -> assign
	* `gt.NullInterval` -> assign
	* `gt.Getter`       -> scan underlying value
*/
func (self *Interval) Scan(src any) error {
	switch src := src.(type) {
	case string:
		return self.Parse(src)

	case []byte:
		return self.UnmarshalText(src)

	case time.Duration:
		self.SetDuration(src)
		return nil

	case Interval:
		*self = src
		return nil

	case NullInterval:
		*self = Interval(src)
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
Sets the interval to an approximate value of the given duration, expressed in
hours, minutes, seconds, truncating any fractions that don't fit.
*/
func (self *Interval) SetDuration(val time.Duration) {
	const minSecs = 60
	const hourMins = 60

	// TODO simpler math.
	hours := int(val.Hours())
	minutes := int(val.Minutes()) - (hours * hourMins)
	seconds := int(val.Seconds()) - (minutes * minSecs) - (hours * hourMins * minSecs)

	*self = Interval{Hours: hours, Minutes: minutes, Seconds: seconds}
}

// Returns the date portion of the interval, disregarding the time portion. The
// result can be passed to `time.Time.AddDate` and `gt.NullTime.AddDate`.
func (self Interval) Date() (years int, months int, days int) {
	return self.Years, self.Months, self.Days
}

// Returns only the date portion of this interval, with other fields set to 0.
func (self Interval) OnlyDate() Interval {
	return Interval{Years: self.Years, Months: self.Months, Days: self.Days}
}

// Returns only the time portion of this interval, with other fields set to 0.
func (self Interval) OnlyTime() Interval {
	return Interval{Hours: self.Hours, Minutes: self.Minutes, Seconds: self.Seconds}
}

// True if the interval has years, months, or days.
func (self Interval) HasDate() bool {
	return self.Years != 0 || self.Months != 0 || self.Days != 0
}

// True if the interval has hours, minutes, or seconds.
func (self Interval) HasTime() bool {
	return self.Hours != 0 || self.Minutes != 0 || self.Seconds != 0
}

/*
Returns the duration of ONLY the time portion of this interval. Panics if the
interval has a date constituent. To make it clear that you're explicitly
disregarding the date part, call `.OnlyTime` first. Warning: there are no
overflow checks. Usage example:

	someInterval.OnlyTime().Duration()
*/
func (self Interval) Duration() time.Duration {
	if self.HasDate() {
		panic(fmt.Errorf(`[gt] failed to convert interval %q to duration: days/months/years can't be converted to nanoseconds`, &self))
	}
	return time.Duration(self.Hours)*time.Hour +
		time.Duration(self.Minutes)*time.Minute +
		time.Duration(self.Seconds)*time.Second
}

// Returns a version of this interval with `.Years = val`.
func (self Interval) WithYears(val int) Interval {
	self.Years = val
	return self
}

// Returns a version of this interval with `.Months = val`.
func (self Interval) WithMonths(val int) Interval {
	self.Months = val
	return self
}

// Returns a version of this interval with `.Days = val`.
func (self Interval) WithDays(val int) Interval {
	self.Days = val
	return self
}

// Returns a version of this interval with `.Hours = val`.
func (self Interval) WithHours(val int) Interval {
	self.Hours = val
	return self
}

// Returns a version of this interval with `.Minutes = val`.
func (self Interval) WithMinutes(val int) Interval {
	self.Minutes = val
	return self
}

// Returns a version of this interval with `.Seconds = val`.
func (self Interval) WithSeconds(val int) Interval {
	self.Seconds = val
	return self
}

// Returns a version of this interval with `.Years += val`.
func (self Interval) AddYears(val int) Interval {
	self.Years += val
	return self
}

// Returns a version of this interval with `.Months += val`.
func (self Interval) AddMonths(val int) Interval {
	self.Months += val
	return self
}

// Returns a version of this interval with `.Days += val`.
func (self Interval) AddDays(val int) Interval {
	self.Days += val
	return self
}

// Returns a version of this interval with `.Hours += val`.
func (self Interval) AddHours(val int) Interval {
	self.Hours += val
	return self
}

// Returns a version of this interval with `.Minutes += val`.
func (self Interval) AddMinutes(val int) Interval {
	self.Minutes += val
	return self
}

// Returns a version of this interval with `.Seconds += val`.
func (self Interval) AddSeconds(val int) Interval {
	self.Seconds += val
	return self
}

/*
Adds every field of one interval to every field of another interval, returning
the sum. Does NOT convert different time units, such as seconds to minutes or
vice versa.
*/
func (self Interval) Add(val Interval) Interval {
	return Interval{
		Years:   self.Years + val.Years,
		Months:  self.Months + val.Months,
		Days:    self.Days + val.Days,
		Hours:   self.Hours + val.Hours,
		Minutes: self.Minutes + val.Minutes,
		Seconds: self.Seconds + val.Seconds,
	}
}

/*
Subtracts every field of one interval from every corresponding field of another
interval, returning the difference. Does NOT convert different time units, such
as seconds to minutes or vice versa.
*/
func (self Interval) Sub(val Interval) Interval {
	return self.Add(val.Neg())
}

/*
Returns a version of this interval with every field inverted: positive fields
become negative, and negative fields become positive.
*/
func (self Interval) Neg() Interval {
	return Interval{
		Years:   -self.Years,
		Months:  -self.Months,
		Days:    -self.Days,
		Hours:   -self.Hours,
		Minutes: -self.Minutes,
		Seconds: -self.Seconds,
	}
}

func (self Interval) bufLen() (num int) {
	if self.IsZero() {
		return len(zeroInterval)
	}

	num += len(`P`)
	addIntervalPartLen(&num, self.Years)
	addIntervalPartLen(&num, self.Months)
	addIntervalPartLen(&num, self.Days)

	if self.HasTime() {
		num += len(`T`)
	}

	addIntervalPartLen(&num, self.Hours)
	addIntervalPartLen(&num, self.Minutes)
	addIntervalPartLen(&num, self.Seconds)
	return
}
