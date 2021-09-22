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
func (self Interval) Get() interface{} { return self.String() }

// Implement `gt.Setter`, using `.Scan`. Panics on error.
func (self *Interval) Set(src interface{}) { try(self.Scan(src)) }

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
	return bytesToMutableString(self.Append(nil))
}

// Implement `gt.Parser`, parsing a valid machine-readable ISO 8601 representation.
func (self *Interval) Parse(src string) (err error) {
	/**
	Regexp-based approach: easy to implement but about 6-10 times slower than
	decent manual parsing. TODO optimize.
	*/

	defer errParse(&err, src, `interval`)

	match := reInterval.FindStringSubmatch(src)
	if match == nil {
		return fmt.Errorf(`format mismatch`)
	}

	var buf Interval

	buf.Years, err = parseIntOpt(match[1])
	if err != nil {
		return err
	}

	buf.Months, err = parseIntOpt(match[2])
	if err != nil {
		return err
	}

	buf.Days, err = parseIntOpt(match[3])
	if err != nil {
		return err
	}

	buf.Hours, err = parseIntOpt(match[4])
	if err != nil {
		return err
	}

	buf.Minutes, err = parseIntOpt(match[5])
	if err != nil {
		return err
	}

	buf.Seconds, err = parseIntOpt(match[6])
	if err != nil {
		return err
	}

	*self = buf
	return nil
}

// Implement `gt.Appender`, using the same representation as `.String`.
func (self Interval) Append(buf []byte) []byte {
	if self.IsZero() {
		return append(buf, zeroInterval...)
	}

	buf = growBytes(buf, self.bufLen())
	buf = append(buf, 'P')
	buf = appendIntervalPart(buf, self.Years, 'Y')
	buf = appendIntervalPart(buf, self.Months, 'M')
	buf = appendIntervalPart(buf, self.Days, 'D')

	if self.HasTime() {
		buf = append(buf, 'T')
	}

	buf = appendIntervalPart(buf, self.Hours, 'H')
	buf = appendIntervalPart(buf, self.Minutes, 'M')
	buf = appendIntervalPart(buf, self.Seconds, 'S')
	return buf
}

// Implement `encoding.TextMarhaler`, using the same representation as `.String`.
func (self Interval) MarshalText() ([]byte, error) {
	return self.Append(nil), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *Interval) UnmarshalText(src []byte) error {
	return self.Parse(bytesToMutableString(src))
}

// Implement `json.Marshaler`, using the same representation as `.String`.
func (self Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Get())
}

// Implement `json.Unmarshaler`, using the same algorithm as `.Parse`.
func (self *Interval) UnmarshalJSON(src []byte) error {
	return jsonUnmarshalString(src, self)
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
func (self *Interval) Scan(src interface{}) error {
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
		ok, err := scanGetter(src, self)
		if ok || err != nil {
			return err
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

	hours := int(val.Hours())
	minutes := int(val.Minutes()) - (hours * hourMins)
	seconds := int(val.Seconds()) - (minutes * minSecs) - (hours * hourMins * minSecs)

	*self = Interval{Hours: hours, Minutes: minutes, Seconds: seconds}
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
Converts to `time.Duration` if possible. Returns an error if the interval has a
date constituent. Warning: there are no overflow checks.
*/
func (self Interval) Duration() (val time.Duration, err error) {
	if self.HasDate() {
		return 0, fmt.Errorf(`failed to convert interval %q to duration: civil time can't be converted to nanoseconds`, &self)
	}

	return time.Duration(self.Hours)*time.Hour +
		time.Duration(self.Minutes)*time.Minute +
		time.Duration(self.Seconds)*time.Second, nil
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
