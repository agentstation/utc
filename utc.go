// Package utc provides a small time.Time wrapper that stores instants in UTC.
//
// The package keeps the underlying time.Time value unexported, normalizes values
// through constructors, parsers, scanners, and serializers, and exposes UTC
// time.Time values through Time and UTC. When compiled with the debug build tag
// (-tags debug), pointer-based methods log nil receiver calls before returning
// errors where Go permits that behavior.
//
// Key features:
//   - Constructors and parsers normalize values to UTC
//   - JSON marshaling/unmarshaling uses strict string/null inputs and preserves
//     sub-second precision
//   - Text and YAML marshal/unmarshal support
//   - SQL database compatibility
//   - Timezone conversion helpers with automatic DST handling
//   - Extensive formatting options for US and EU date formats
//
// Debug mode:
//
//	To enable debug logging, compile with: go build -tags debug
//	This logs nil receiver calls for pointer-based methods that can return errors.
package utc

import (
	"bytes"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
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

// Time zone locations cached at package level.
var (
	pacificLocation  *time.Location
	easternLocation  *time.Location
	centralLocation  *time.Location
	mountainLocation *time.Location
	locationError    = initLocations()
)

func initLocations() error {
	var err error
	pacificLocation, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return fmt.Errorf("failed to load Pacific timezone: %w", err)
	}

	easternLocation, err = time.LoadLocation("America/New_York")
	if err != nil {
		return fmt.Errorf("failed to load Eastern timezone: %w", err)
	}

	centralLocation, err = time.LoadLocation("America/Chicago")
	if err != nil {
		return fmt.Errorf("failed to load Central timezone: %w", err)
	}

	mountainLocation, err = time.LoadLocation("America/Denver")
	if err != nil {
		return fmt.Errorf("failed to load Mountain timezone: %w", err)
	}

	return nil
}

// ValidateTimezoneAvailability checks if all timezone locations were properly initialized
// Returns nil if initialization was successful, otherwise returns the initialization error
func ValidateTimezoneAvailability() error {
	if locationError != nil {
		return fmt.Errorf("timezone locations not properly initialized: %w", locationError)
	}
	return nil
}

// Time stores a time instant normalized to UTC.
type Time struct {
	t time.Time
}

// Now returns the current time in UTC
func Now() Time {
	return New(time.Now())
}

// New returns a new Time from a time.Time
func New(t time.Time) Time {
	return Time{t: t.UTC()}
}

func (t Time) utc() time.Time {
	return t.t.UTC()
}

// Time returns the underlying time.Time normalized to UTC.
func (t Time) Time() time.Time {
	return t.utc()
}

// ParseRFC3339 parses a time string in RFC3339 format and returns a utc.Time
func ParseRFC3339(s string) (Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Time{}, err
	}
	return New(t), nil
}

// ParseRFC3339Nano parses a time string in RFC3339Nano format and returns a utc.Time
func ParseRFC3339Nano(s string) (Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return Time{}, err
	}
	return New(t), nil
}

// Parse parses a time string in the specified format and returns a utc.Time
func Parse(layout string, s string) (Time, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return Time{}, err
	}
	return New(t), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for utc.Time
func (t *Time) UnmarshalJSON(data []byte) error {
	if t == nil {
		debugLog("UnmarshalJSON() called on nil *Time receiver")
		return errors.New("cannot unmarshal into nil utc.Time")
	}
	data = bytes.TrimSpace(data)

	// Handle empty data
	if len(data) == 0 {
		return errors.New("cannot unmarshal empty data into utc.Time")
	}

	// Handle null
	if string(data) == "null" {
		t.t = time.Time{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("utc.Time must be a JSON string or null: %w", err)
	}

	if s == "" {
		t.t = time.Time{}
		return nil
	}

	// Parse the time (allow a few flexible formats)
	parsedTime, err := parse(s)
	if err != nil {
		return err
	}

	// Convert to UTC
	t.t = parsedTime.UTC()
	return nil
}

// MarshalJSON implements the json.Marshaler interface for utc.Time.
// Returns an error for nil receivers to maintain consistency with standard marshaling behavior.
func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		debugLog("MarshalJSON() called on nil *Time receiver")
		return nil, errors.New("cannot marshal nil utc.Time")
	}
	return t.utc().MarshalJSON()
}

// Ensure Time implements encoding.TextMarshaler/TextUnmarshaler for broader codec support.
var (
	_ encoding.TextMarshaler   = Time{}
	_ encoding.TextUnmarshaler = (*Time)(nil)
	_ driver.Valuer            = Time{}
)

// MarshalText implements encoding.TextMarshaler.
func (t Time) MarshalText() ([]byte, error) {
	return t.utc().MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Time) UnmarshalText(text []byte) error {
	if t == nil {
		debugLog("UnmarshalText() called on nil *Time receiver")
		return errors.New("cannot unmarshal text into nil utc.Time")
	}
	if len(text) == 0 {
		t.t = time.Time{}
		return nil
	}
	parsed, err := parse(string(text))
	if err != nil {
		return err
	}
	t.t = parsed.UTC()
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for utc.Time
func (t *Time) UnmarshalYAML(unmarshal func(any) error) error {
	if t == nil {
		debugLog("UnmarshalYAML() called on nil *Time receiver")
		return errors.New("cannot unmarshal YAML into nil utc.Time")
	}

	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// Handle empty string
	if s == "" {
		t.t = time.Time{}
		return nil
	}

	// Parse the time string using our flexible parser
	parsed, err := parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse time %q: %w", s, err)
	}

	t.t = parsed.UTC()
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface for utc.Time
func (t Time) MarshalYAML() (any, error) {
	if t.utc().IsZero() {
		return nil, nil
	}

	// Use RFC3339Nano format so YAML preserves sub-second precision.
	return t.utc().Format(time.RFC3339Nano), nil
}

// String implements fmt.Stringer. It prints the time in RFC3339Nano format.
func (t Time) String() string {
	return t.utc().Format(time.RFC3339Nano)
}

// Value implements the driver.Valuer interface for database operations for utc.Time.
// It returns the UTC time.Time value.
func (t Time) Value() (driver.Value, error) {
	// Preserve previous behavior: zero value still returns a non-nil time
	return t.utc(), nil
}

// Scan implements the sql.Scanner interface for database operations for utc.Time
// It does this by scanning the value into a time.Time, converting the time.Time to UTC,
// and then assigning the UTC time to the utc.Time.
func (t *Time) Scan(value any) error {
	if t == nil {
		debugLog("Scan() called on nil *Time receiver")
		return errors.New("cannot scan into nil utc.Time")
	}

	if value == nil {
		return errors.New("cannot scan nil into utc.Time")
	}

	switch v := value.(type) {
	case time.Time:
		t.t = v.UTC()
		return nil
	case string:
		parsed, err := parse(v)
		if err != nil {
			return err
		}
		t.t = parsed.UTC()
		return nil
	case []byte:
		parsed, err := parse(string(v))
		if err != nil {
			return err
		}
		t.t = parsed.UTC()
		return nil
	default:
		return errors.New("cannot scan non-time value into utc.Time")
	}
}

// Before reports whether the time is before u
func (t Time) Before(u Time) bool {
	return t.utc().Before(u.utc())
}

// After reports whether the time is after u
func (t Time) After(u Time) bool {
	return t.utc().After(u.utc())
}

// Equal reports whether t and u represent the same time instant
func (t Time) Equal(u Time) bool {
	return t.utc().Equal(u.utc())
}

// Add returns the time t+d
func (t Time) Add(d time.Duration) Time {
	return New(t.utc().Add(d))
}

// Sub returns the duration t-u
func (t Time) Sub(u Time) time.Duration {
	return t.utc().Sub(u.utc())
}

// UTC returns t in UTC
func (t Time) UTC() time.Time {
	return t.utc()
}

// PST returns t in PST
func (t Time) PST() time.Time {
	return t.utc().In(time.FixedZone("PST", -8*60*60))
}

// EST returns t in EST
func (t Time) EST() time.Time {
	return t.utc().In(time.FixedZone("EST", -5*60*60))
}

// CST returns t in CST
func (t Time) CST() time.Time {
	return t.utc().In(time.FixedZone("CST", -6*60*60))
}

// MST returns t in MST
func (t Time) MST() time.Time {
	return t.utc().In(time.FixedZone("MST", -7*60*60))
}

// Pacific returns t in Pacific time (handles PST/PDT automatically)
func (t Time) Pacific() time.Time {
	if locationError != nil {
		return t.PST() // Fall back to fixed PST if location isn't available
	}
	return t.utc().In(pacificLocation)
}

// Eastern returns t in Eastern time (handles EST/EDT automatically)
func (t Time) Eastern() time.Time {
	if locationError != nil {
		return t.EST() // Fall back to fixed EST if location isn't available
	}
	return t.utc().In(easternLocation)
}

// Central returns t in Central time (handles CST/CDT automatically)
func (t Time) Central() time.Time {
	if locationError != nil {
		return t.CST() // Fall back to fixed CST if location isn't available
	}
	return t.utc().In(centralLocation)
}

// Mountain returns t in Mountain time (handles MST/MDT automatically)
func (t Time) Mountain() time.Time {
	if locationError != nil {
		return t.MST() // Fall back to fixed MST if location isn't available
	}
	return t.utc().In(mountainLocation)
}

// Add the useful utility methods while maintaining chainability
func (t Time) IsZero() bool {
	return t.utc().IsZero()
}

// Format formats the time using the specified layout
func (t Time) Format(layout string) string {
	return t.utc().Format(layout)
}

// TimeFormat formats the time using the specified layout
func (t Time) TimeFormat(layout TimeLayout) string {
	return t.utc().Format(string(layout))
}

// Standard/ISO formats
// -------------------

// RFC3339 formats time as "2006-01-02T15:04:05Z07:00"
func (t Time) RFC3339() string {
	return t.utc().Format(time.RFC3339)
}

// RFC3339Nano formats time as "2006-01-02T15:04:05.999999999Z07:00"
func (t Time) RFC3339Nano() string {
	return t.utc().Format(time.RFC3339Nano)
}

// ISO8601 formats time as "2006-01-02T15:04:05Z07:00" (same as RFC3339)
func (t Time) ISO8601() string {
	return t.utc().Format(time.RFC3339)
}

// RFC822 formats time as "02 Jan 06 15:04 MST"
func (t Time) RFC822() string {
	return t.utc().Format(time.RFC822)
}

// RFC822Z formats time as "02 Jan 06 15:04 -0700"
func (t Time) RFC822Z() string {
	return t.utc().Format(time.RFC822Z)
}

// RFC850 formats time as "Monday, 02-Jan-06 15:04:05 MST"
func (t Time) RFC850() string {
	return t.utc().Format(time.RFC850)
}

// ANSIC formats time as "Mon Jan _2 15:04:05 2006"
func (t Time) ANSIC() string {
	return t.utc().Format(time.ANSIC)
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
	return t.utc().Format(time.Kitchen)
}

// Generic location helpers and utilities

// In converts time to a named location (e.g., "America/Los_Angeles").
func (t Time) In(name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, err
	}
	return t.utc().In(loc), nil
}

// InLocation converts time to a provided *time.Location.
func (t Time) InLocation(loc *time.Location) time.Time {
	if loc == nil {
		return t.utc()
	}
	return t.utc().In(loc)
}

// Unix helpers
func FromUnix(sec int64) Time     { return New(time.Unix(sec, 0)) }
func FromUnixMilli(ms int64) Time { return New(time.UnixMilli(ms)) }
func (t Time) Unix() int64        { return t.utc().Unix() }
func (t Time) UnixMilli() int64   { return t.utc().UnixMilli() }

// Day helpers - times are always in UTC within this package
func (t Time) StartOfDay() Time {
	y, m, d := t.utc().Date()
	return New(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
}

func (t Time) EndOfDay() Time {
	y, m, d := t.utc().Date()
	// One nanosecond before next midnight
	return New(time.Date(y, m, d+1, 0, 0, 0, -1, time.UTC))
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
