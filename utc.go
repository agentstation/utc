// Package utc provides a time.Time wrapper that ensures all times are in UTC.
//
// The package offers enhanced safety by gracefully handling nil receivers instead
// of panicking, making it more suitable for production environments. When compiled
// with the debug build tag (-tags debug), it provides additional logging for nil
// receiver method calls to help identify potential bugs during development.
//
// Key features:
//   - All times are automatically converted to and stored in UTC
//   - JSON marshaling/unmarshaling with flexible parsing
//   - SQL database compatibility
//   - Timezone conversion helpers with automatic DST handling
//   - Extensive formatting options for US and EU date formats
//   - Nil-safe operations that return errors instead of panicking
//
// Debug mode:
//
//	To enable debug logging, compile with: go build -tags debug
//	This will log warnings when methods are called on nil receivers.
package utc

import (
	"database/sql/driver"
	"encoding"
	"errors"
	"fmt"
	"sync"
	"time"
)

type TimeLayout string

// Add layout constants at package level
const (
	TimeLayoutUSDateShort  TimeLayout = "01/02/2006"
	TimeLayoutUSDateLong   TimeLayout = "January 2, 2006"
	TimeLayoutUSDateTime12 TimeLayout = "01/02/2006 03:04:05 PM"
	TimeLayoutUSDateTime24 TimeLayout = "01/02/2006 15:04:05"
	TimeLayoutUSTime12     TimeLayout = "3:04 PM"
	TimeLayoutUSTime24     TimeLayout = "15:04"

	TimeLayoutEUDateShort  TimeLayout = "02/01/2006"
	TimeLayoutEUDateLong   TimeLayout = "2 January 2006"
	TimeLayoutEUDateTime12 TimeLayout = "02/01/2006 03:04:05 PM"
	TimeLayoutEUDateTime24 TimeLayout = "02/01/2006 15:04:05"
	TimeLayoutEUTime12     TimeLayout = "3:04 PM"
	TimeLayoutEUTime24     TimeLayout = "15:04"

	TimeLayoutDateOnly     TimeLayout = "2006-01-02"
	TimeLayoutTimeOnly     TimeLayout = "15:04:05"
	TimeLayoutWeekdayLong  TimeLayout = "Monday"
	TimeLayoutWeekdayShort TimeLayout = "Mon"
	TimeLayoutMonthLong    TimeLayout = "January"
	TimeLayoutMonthShort   TimeLayout = "Jan"
)

var (
	// Time zone locations cached at package level
	pacificLocation  *time.Location
	easternLocation  *time.Location
	centralLocation  *time.Location
	mountainLocation *time.Location

	// Initialize locations
	locationError = initLocations()
)

// initLocations loads all time zone locations at startup
// Use lazy initialization to avoid upfront tz loading cost, while preserving
// compatibility via locationError observation.
var tzOnce sync.Once
var tzInitErr error

func initLocations() error {
	tzOnce.Do(func() {
		var err error
		pacificLocation, err = time.LoadLocation("America/Los_Angeles")
		if err != nil {
			tzInitErr = fmt.Errorf("failed to load Pacific timezone: %w", err)
			return
		}

		easternLocation, err = time.LoadLocation("America/New_York")
		if err != nil {
			tzInitErr = fmt.Errorf("failed to load Eastern timezone: %w", err)
			return
		}

		centralLocation, err = time.LoadLocation("America/Chicago")
		if err != nil {
			tzInitErr = fmt.Errorf("failed to load Central timezone: %w", err)
			return
		}

		mountainLocation, err = time.LoadLocation("America/Denver")
		if err != nil {
			tzInitErr = fmt.Errorf("failed to load Mountain timezone: %w", err)
			return
		}
	})
	return tzInitErr
}

// ValidateTimezoneAvailability checks if all timezone locations were properly initialized
// Returns nil if initialization was successful, otherwise returns the initialization error
func ValidateTimezoneAvailability() error {
	if locationError != nil {
		return fmt.Errorf("timezone locations not properly initialized: %w", locationError)
	}
	return nil
}

// Time is an alias for time.Time that defaults to UTC time.
type Time struct {
	time.Time
}

// Now returns the current time in UTC
func Now() Time {
	return Time{time.Now().UTC()}
}

// New returns a new Time from a time.Time
func New(t time.Time) Time {
	return Time{t.UTC()}
}

// ParseRFC3339 parses a time string in RFC3339 format and returns a utc.Time
func ParseRFC3339(s string) (Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Time{}, err
	}
	return Time{t.UTC()}, nil
}

// ParseRFC3339Nano parses a time string in RFC3339Nano format and returns a utc.Time
func ParseRFC3339Nano(s string) (Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return Time{}, err
	}
	return Time{t.UTC()}, nil
}

// Parse parses a time string in the specified format and returns a utc.Time
func Parse(layout string, s string) (Time, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return Time{}, err
	}
	return Time{t.UTC()}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for utc.Time
func (t *Time) UnmarshalJSON(data []byte) error {
	// Handle empty data
	if len(data) == 0 {
		return errors.New("cannot unmarshal empty data into utc.Time")
	}

	// Handle null or empty string
	if string(data) == "null" || string(data) == `""` {
		t.Time = time.Time{}
		return nil
	}

	// Remove quotes
	if len(data) > 2 && (data[0] == '"' && data[len(data)-1] == '"') {
		data = data[1 : len(data)-1]
	}

	// Parse the time (allow a few flexible formats)
	parsedTime, err := parse(string(data))
	if err != nil {
		return err
	}

	// Convert to UTC
	t.Time = parsedTime.UTC()
	return nil
}

// MarshalJSON implements the json.Marshaler interface for utc.Time.
// Returns an error for nil receivers to maintain consistency with standard marshaling behavior.
func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		debugLog("MarshalJSON() called on nil *Time receiver")
		return nil, errors.New("cannot marshal nil utc.Time")
	}
	return []byte(`"` + t.Time.Format(time.RFC3339) + `"`), nil
}

// Ensure Time implements encoding.TextMarshaler/TextUnmarshaler for broader codec support.
var (
	_ encoding.TextMarshaler   = Time{}
	_ encoding.TextUnmarshaler = (*Time)(nil)
)

// MarshalText implements encoding.TextMarshaler.
func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.Time.Format(time.RFC3339)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Time) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		t.Time = time.Time{}
		return nil
	}
	parsed, err := parse(string(text))
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for utc.Time
func (t *Time) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// Handle empty string
	if s == "" {
		t.Time = time.Time{}
		return nil
	}

	// Parse the time string using our flexible parser
	parsed, err := parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse time %q: %w", s, err)
	}

	t.Time = parsed
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface for utc.Time
func (t Time) MarshalYAML() (any, error) {
	if t.Time.IsZero() {
		return nil, nil
	}

	// Use RFC3339 format for YAML output
	return t.Time.Format(time.RFC3339), nil
}

// String implements the Stringer interface for utc.Time. It prints the time in RFC3339 format.
//
// Unlike many Go types that panic on nil receivers, this method returns "<nil>" to match
// stdlib conventions (e.g., bytes.Buffer) and improve production safety. In debug builds
// (compiled with -tags debug), nil receivers are logged to help identify potential bugs.
func (t *Time) String() string {
	if t == nil {
		debugLog("String() called on nil *Time receiver")
		return "<nil>"
	}
	return t.Time.Format(time.RFC3339)
}

// Value implements the driver.Valuer interface for database operations for utc.Time.
// It returns the time.Time value and assumes the time is already in UTC.
//
// Returns an error if called on a nil receiver instead of panicking to allow graceful
// error handling in database operations. In debug builds, nil receivers are logged.
func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		debugLog("Value() called on nil *Time receiver")
		return nil, errors.New("cannot call Value() on nil utc.Time")
	}
	// Preserve previous behavior: zero value still returns a non-nil time
	return t.Time, nil
}

// Scan implements the sql.Scanner interface for database operations for utc.Time
// It does this by scanning the value into a time.Time, converting the time.Time to UTC,
// and then assigning the UTC time to the utc.Time.
func (t *Time) Scan(value any) error {
	if value == nil {
		return errors.New("cannot scan nil into utc.Time")
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v.UTC()
		return nil
	case string:
		parsed, err := parse(v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	case []byte:
		parsed, err := parse(string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return errors.New("cannot scan non-time value into utc.Time")
	}
}

// Before reports whether the time is before u
func (t Time) Before(u Time) bool {
	return t.Time.Before(u.Time)
}

// After reports whether the time is after u
func (t Time) After(u Time) bool {
	return t.Time.After(u.Time)
}

// Equal reports whether t and u represent the same time instant
func (t Time) Equal(u Time) bool {
	return t.Time.Equal(u.Time)
}

// Add returns the time t+d
func (t Time) Add(d time.Duration) Time {
	return Time{t.Time.Add(d)}
}

// Sub returns the duration t-u
func (t Time) Sub(u Time) time.Duration {
	return t.Time.Sub(u.Time)
}

// UTC returns t in UTC
func (t Time) UTC() time.Time {
	return t.Time
}

// PST returns t in PST
func (t Time) PST() time.Time {
	return t.Time.In(time.FixedZone("PST", -8*60*60))
}

// EST returns t in EST
func (t Time) EST() time.Time {
	return t.Time.In(time.FixedZone("EST", -5*60*60))
}

// CST returns t in CST
func (t Time) CST() time.Time {
	return t.Time.In(time.FixedZone("CST", -6*60*60))
}

// MST returns t in MST
func (t Time) MST() time.Time {
	return t.Time.In(time.FixedZone("MST", -7*60*60))
}

// Pacific returns t in Pacific time (handles PST/PDT automatically)
func (t Time) Pacific() time.Time {
	if locationError != nil {
		return t.PST() // Fall back to fixed PST if location isn't available
	}
	return t.Time.In(pacificLocation)
}

// Eastern returns t in Eastern time (handles EST/EDT automatically)
func (t Time) Eastern() time.Time {
	if locationError != nil {
		return t.EST() // Fall back to fixed EST if location isn't available
	}
	return t.Time.In(easternLocation)
}

// Central returns t in Central time (handles CST/CDT automatically)
func (t Time) Central() time.Time {
	if locationError != nil {
		return t.CST() // Fall back to fixed CST if location isn't available
	}
	return t.Time.In(centralLocation)
}

// Mountain returns t in Mountain time (handles MST/MDT automatically)
func (t Time) Mountain() time.Time {
	if locationError != nil {
		return t.MST() // Fall back to fixed MST if location isn't available
	}
	return t.Time.In(mountainLocation)
}

// Add the useful utility methods while maintaining chainability
func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

// Format formats the time using the specified layout
func (t Time) Format(layout string) string {
	return t.Time.Format(layout)
}

// TimeFormat formats the time using the specified layout
func (t Time) TimeFormat(layout TimeLayout) string {
	return t.Time.Format(string(layout))
}

// Standard/ISO formats
// -------------------

// RFC3339 formats time as "2006-01-02T15:04:05Z07:00"
func (t Time) RFC3339() string {
	return t.Time.Format(time.RFC3339)
}

// RFC3339Nano formats time as "2006-01-02T15:04:05.999999999Z07:00"
func (t Time) RFC3339Nano() string {
	return t.Time.Format(time.RFC3339Nano)
}

// ISO8601 formats time as "2006-01-02T15:04:05Z07:00" (same as RFC3339)
func (t Time) ISO8601() string {
	return t.Time.Format(time.RFC3339)
}

// RFC822 formats time as "02 Jan 06 15:04 MST"
func (t Time) RFC822() string {
	return t.Time.Format(time.RFC822)
}

// RFC822Z formats time as "02 Jan 06 15:04 -0700"
func (t Time) RFC822Z() string {
	return t.Time.Format(time.RFC822Z)
}

// RFC850 formats time as "Monday, 02-Jan-06 15:04:05 MST"
func (t Time) RFC850() string {
	return t.Time.Format(time.RFC850)
}

// ANSIC formats time as "Mon Jan _2 15:04:05 2006"
func (t Time) ANSIC() string {
	return t.Time.Format(time.ANSIC)
}

// US Regional formats (MM/DD/YYYY)
// ------------------------------

// USDateShort formats time as "01/02/2006"
func (t Time) USDateShort() string {
	return t.TimeFormat(TimeLayoutUSDateShort)
}

// USDateLong formats time as "January 2, 2006"
func (t Time) USDateLong() string {
	return t.TimeFormat(TimeLayoutUSDateLong)
}

// USDateTime12 formats time as "01/02/2006 03:04:05 PM"
func (t Time) USDateTime12() string {
	return t.TimeFormat(TimeLayoutUSDateTime12)
}

// USDateTime24 formats time as "01/02/2006 15:04:05"
func (t Time) USDateTime24() string {
	return t.TimeFormat(TimeLayoutUSDateTime24)
}

// USTime12 formats time as "3:04 PM"
func (t Time) USTime12() string {
	return t.TimeFormat(TimeLayoutUSTime12)
}

// USTime24 formats time as "15:04"
func (t Time) USTime24() string {
	return t.TimeFormat(TimeLayoutUSTime24)
}

// European formats (DD/MM/YYYY)
// ---------------------------

// EUDateShort formats time as "02/01/2006"
func (t Time) EUDateShort() string {
	return t.TimeFormat(TimeLayoutEUDateShort)
}

// EUDateLong formats time as "2 January 2006"
func (t Time) EUDateLong() string {
	return t.TimeFormat(TimeLayoutEUDateLong)
}

// EUDateTime12 formats time as "02/01/2006 03:04:05 PM"
func (t Time) EUDateTime12() string {
	return t.TimeFormat(TimeLayoutEUDateTime12)
}

// EUDateTime24 formats time as "02/01/2006 15:04:05"
func (t Time) EUDateTime24() string {
	return t.TimeFormat(TimeLayoutEUDateTime24)
}

// EUTime12 formats time as "3:04 PM"
func (t Time) EUTime12() string {
	return t.TimeFormat(TimeLayoutEUTime12)
}

// EUTime24 formats time as "15:04"
func (t Time) EUTime24() string {
	return t.TimeFormat(TimeLayoutEUTime24)
}

// Common Components
// ---------------

// WeekdayLong formats time as "Monday"
func (t Time) WeekdayLong() string {
	return t.TimeFormat(TimeLayoutWeekdayLong)
}

// WeekdayShort formats time as "Mon"
func (t Time) WeekdayShort() string {
	return t.TimeFormat(TimeLayoutWeekdayShort)
}

// MonthLong formats time as "January"
func (t Time) MonthLong() string {
	return t.TimeFormat(TimeLayoutMonthLong)
}

// MonthShort formats time as "Jan"
func (t Time) MonthShort() string {
	return t.TimeFormat(TimeLayoutMonthShort)
}

// DateOnly formats time as "2006-01-02"
func (t Time) DateOnly() string {
	return t.TimeFormat(TimeLayoutDateOnly)
}

// TimeOnly formats time as "15:04:05"
func (t Time) TimeOnly() string {
	return t.TimeFormat(TimeLayoutTimeOnly)
}

// Kitchen formats time as "3:04PM"
func (t Time) Kitchen() string {
	return t.Time.Format(time.Kitchen)
}

// Generic location helpers and utilities

// In converts time to a named location (e.g., "America/Los_Angeles").
func (t Time) In(name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, err
	}
	return t.Time.In(loc), nil
}

// InLocation converts time to a provided *time.Location.
func (t Time) InLocation(loc *time.Location) time.Time {
	return t.Time.In(loc)
}

// Unix helpers
func FromUnix(sec int64) Time     { return Time{time.Unix(sec, 0).UTC()} }
func FromUnixMilli(ms int64) Time { return Time{time.Unix(0, ms*int64(time.Millisecond)).UTC()} }
func (t Time) Unix() int64        { return t.Time.Unix() }
func (t Time) UnixMilli() int64   { return t.Time.UnixMilli() }

// Day helpers - times are always in UTC within this package
func (t Time) StartOfDay() Time {
	y, m, d := t.Time.UTC().Date()
	return Time{time.Date(y, m, d, 0, 0, 0, 0, time.UTC)}
}

func (t Time) EndOfDay() Time {
	y, m, d := t.Time.UTC().Date()
	// One nanosecond before next midnight
	return Time{time.Date(y, m, d+1, 0, 0, 0, -1, time.UTC)}
}

// Internal: parse a variety of common layouts to UTC.
func parse(s string) (time.Time, error) {
	tryLayouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01", // YYYY-MM format
		"2006",    // YYYY format
	}
	var firstErr error
	for _, layout := range tryLayouts {
		if parsed, err := time.Parse(layout, s); err == nil {
			return parsed.UTC(), nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	return time.Time{}, firstErr
}
