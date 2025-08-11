package utc

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestUTC_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid UTC time",
			input:   `"2023-01-01T12:00:00Z"`,
			want:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid time with offset",
			input:   `"2023-01-01T12:00:00+02:00"`,
			want:    time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), // Converted to UTC
			wantErr: false,
		},
		{
			name:    "invalid time format",
			input:   `"2023-13-01T12:00:00Z"`,
			wantErr: true,
		},
		{
			name:    "null value",
			input:   `null`,
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "empty string value",
			input:   `""`,
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "empty data",
			input:   ``,
			wantErr: true,
		},
		{
			name:    "malformed JSON",
			input:   `"2023-01-01T12:00:00Z`,
			wantErr: true,
		},
		{
			name:    "non-string JSON",
			input:   `123`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ut Time
			err := ut.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !ut.Time.Equal(tt.want) {
					t.Errorf("Time.UnmarshalJSON() = %v, want %v", ut.Time, tt.want)
				}
				// Verify location is UTC
				if ut.Time.Location() != time.UTC {
					t.Error("Unmarshaled time is not in UTC")
				}
			}
		})
	}
}
func TestUTC_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		time    time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "UTC time",
			time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			want:    `"2023-01-01T12:00:00Z"`,
			wantErr: false,
		},
		{
			name:    "non-UTC time",
			time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			want:    `"2023-01-01T17:00:00Z"`, // Converted to UTC
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := Time{tt.time.UTC()}
			got, err := ut.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("Time.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
func TestUTC_RoundTrip(t *testing.T) {
	type TestStruct struct {
		Timestamp Time `json:"timestamp"`
	}
	original := TestStruct{
		Timestamp: Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
	}
	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}
	// Unmarshal back
	var decoded TestStruct
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	// Compare
	if !decoded.Timestamp.Time.Equal(original.Timestamp.Time) {
		t.Errorf("Round trip failed: got %v, want %v",
			decoded.Timestamp.Time, original.Timestamp.Time)
	}
}
func TestUTC_Now(t *testing.T) {
	now := Now()
	time.Sleep(time.Millisecond)
	later := Now()
	if !now.Before(later) {
		t.Error("Now() did not return increasing times")
	}
	if now.Time.Location() != time.UTC {
		t.Error("Now() time is not in UTC")
	}
}
func TestUTC_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UTC time",
			input:   "2023-01-01T12:00:00Z",
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "2023-01-01",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut, err := ParseRFC3339(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ut.Time.Location() != time.UTC {
				t.Error("Parsed time is not in UTC")
			}
		})
	}
}
func TestUTC_String(t *testing.T) {
	ut := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	want := "2023-01-01T12:00:00Z"
	if got := ut.String(); got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}
func TestUTC_DatabaseOperations(t *testing.T) {
	original := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	ut := Time{original}
	// Test Value
	value, err := ut.Value()
	if err != nil {
		t.Fatalf("Value() error = %v", err)
	}
	// Test Scan
	var scanned Time
	err = scanned.Scan(value)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if !scanned.Equal(ut) {
		t.Errorf("Scan() = %v, want %v", scanned, ut)
	}
	// Test scanning nil
	err = scanned.Scan(nil)
	if err == nil {
		t.Error("Scan(nil) should return error")
	}
	// Test scanning invalid type
	err = scanned.Scan(42)
	if err == nil {
		t.Error("Scan(int) should return error")
	}
}
func TestUTC_Comparisons(t *testing.T) {
	t1 := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	t2 := Time{time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)}
	if !t1.Before(t2) {
		t.Error("Before() failed")
	}
	if !t2.After(t1) {
		t.Error("After() failed")
	}
	if t1.Equal(t2) {
		t.Error("Equal() failed")
	}
}
func TestUTC_Arithmetic(t *testing.T) {
	t1 := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	duration := time.Hour
	// Test Add
	t2 := t1.Add(duration)
	if !t2.Equal(Time{time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)}) {
		t.Error("Add() failed")
	}
	// Test Sub
	if t2.Sub(t1) != duration {
		t.Error("Sub() failed")
	}
}
func TestUTC_IsZero(t *testing.T) {
	var zero Time
	if !zero.IsZero() {
		t.Error("IsZero() should be true for zero value")
	}
	nonZero := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	if nonZero.IsZero() {
		t.Error("IsZero() should be false for non-zero value")
	}
}
func TestUTC_UnixTimestamps(t *testing.T) {
	ut := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	// Test Unix
	if ut.Unix() != ut.UTC().Unix() {
		t.Error("Unix() returned incorrect value")
	}
	// Test UnixMilli
	if ut.UnixMilli() != ut.UTC().UnixMilli() {
		t.Error("UnixMilli() returned incorrect value")
	}
}
func TestUTC_New(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want time.Time
	}{
		{
			name: "UTC time",
			time: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			want: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "non-UTC time",
			time: time.Date(2023, 1, 1, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			want: time.Date(2023, 1, 1, 17, 0, 0, 0, time.UTC), // Converted to UTC
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.time)
			if !got.UTC().Equal(tt.want) {
				t.Errorf("New() = %v, want %v", got.UTC(), tt.want)
			}
			if got.UTC().Location() != time.UTC {
				t.Error("New() time is not in UTC")
			}
		})
	}
}
func TestUTC_TimeZoneConversions(t *testing.T) {
	// Use a fixed time that won't be affected by DST
	ut := Time{time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)} // Noon UTC
	tests := []struct {
		name     string
		convert  func(Time) time.Time
		wantHour int // Expected hour in target timezone
	}{
		{
			name:     "UTC to Pacific",
			convert:  Time.Pacific,
			wantHour: 4, // UTC-8
		},
		{
			name:     "UTC to Eastern",
			convert:  Time.Eastern,
			wantHour: 7, // UTC-5
		},
		{
			name:     "UTC to Central",
			convert:  Time.Central,
			wantHour: 6, // UTC-6
		},
		{
			name:     "UTC to Mountain",
			convert:  Time.Mountain,
			wantHour: 5, // UTC-7
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converted := tt.convert(ut)
			if converted.Hour() != tt.wantHour {
				t.Errorf("Expected hour to be %d, got %d", tt.wantHour, converted.Hour())
			}
		})
	}
}
func TestUTC_Formatting(t *testing.T) {
	ut := Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)}
	tests := []struct {
		name     string
		format   func(Time) string
		expected string
	}{
		{
			name:     "USDateShort",
			format:   Time.USDateShort,
			expected: "01/02/2024",
		},
		{
			name:     "USDateLong",
			format:   Time.USDateLong,
			expected: "January 2, 2024",
		},
		{
			name:     "EUDateShort",
			format:   Time.EUDateShort,
			expected: "02/01/2024",
		},
		{
			name:     "EUDateLong",
			format:   Time.EUDateLong,
			expected: "2 January 2024",
		},
		{
			name:     "WeekdayLong",
			format:   Time.WeekdayLong,
			expected: "Tuesday",
		},
		{
			name:     "MonthShort",
			format:   Time.MonthShort,
			expected: "Jan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format(ut)
			if result != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}
func TestUTC_TimezoneError(t *testing.T) {
	// Save original locationErr and restore it after test
	originalErr := locationError
	defer func() { locationError = originalErr }()
	// Simulate timezone initialization failure
	locationError = errors.New("simulated timezone initialization error")
	testTime := Time{time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
	tests := []struct {
		name     string
		got      time.Time
		expected time.Time
	}{
		{
			name:     "Pacific",
			got:      testTime.Pacific(),
			expected: testTime.PST(),
		},
		{
			name:     "Eastern",
			got:      testTime.Eastern(),
			expected: testTime.EST(),
		},
		{
			name:     "Central",
			got:      testTime.Central(),
			expected: testTime.CST(),
		},
		{
			name:     "Mountain",
			got:      testTime.Mountain(),
			expected: testTime.MST(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.expected) {
				t.Errorf("Expected %v to fall back to %v, got %v",
					tt.name, tt.expected, tt.got)
			}
		})
	}
}
func TestValidateTimezoneAvailability(t *testing.T) {
	err := ValidateTimezoneAvailability()
	if !errors.Is(err, locationError) {
		t.Errorf("Expected error to wrap %v, got %v", locationError, err)
	}
}
func TestUTC_AdditionalFormats(t *testing.T) {
	ut := Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)}
	tests := []struct {
		name     string
		format   func(Time) string
		expected string
	}{
		{
			name:     "RFC3339Nano",
			format:   Time.RFC3339Nano,
			expected: "2024-01-02T15:04:05Z",
		},
		{
			name:     "ISO8601",
			format:   Time.ISO8601,
			expected: "2024-01-02T15:04:05Z",
		},
		{
			name:     "RFC822",
			format:   Time.RFC822,
			expected: "02 Jan 24 15:04 UTC",
		},
		{
			name:     "RFC822Z",
			format:   Time.RFC822Z,
			expected: "02 Jan 24 15:04 +0000",
		},
		{
			name:     "RFC850",
			format:   Time.RFC850,
			expected: "Tuesday, 02-Jan-24 15:04:05 UTC",
		},
		{
			name:     "ANSIC",
			format:   Time.ANSIC,
			expected: "Tue Jan  2 15:04:05 2024",
		},
		{
			name:     "Kitchen",
			format:   Time.Kitchen,
			expected: "3:04PM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format(ut)
			if result != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}
func TestUTC_AdditionalTimeFormats(t *testing.T) {
	ut := Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)}
	tests := []struct {
		name     string
		format   func(Time) string
		expected string
	}{
		{
			name:     "USDateTime12",
			format:   Time.USDateTime12,
			expected: "01/02/2024 03:04:05 PM",
		},
		{
			name:     "USDateTime24",
			format:   Time.USDateTime24,
			expected: "01/02/2024 15:04:05",
		},
		{
			name:     "USTime12",
			format:   Time.USTime12,
			expected: "3:04 PM",
		},
		{
			name:     "USTime24",
			format:   Time.USTime24,
			expected: "15:04",
		},
		{
			name:     "EUDateTime12",
			format:   Time.EUDateTime12,
			expected: "02/01/2024 03:04:05 PM",
		},
		{
			name:     "EUDateTime24",
			format:   Time.EUDateTime24,
			expected: "02/01/2024 15:04:05",
		},
		{
			name:     "EUTime12",
			format:   Time.EUTime12,
			expected: "3:04 PM",
		},
		{
			name:     "EUTime24",
			format:   Time.EUTime24,
			expected: "15:04",
		},
		{
			name:     "DateOnly",
			format:   Time.DateOnly,
			expected: "2024-01-02",
		},
		{
			name:     "TimeOnly",
			format:   Time.TimeOnly,
			expected: "15:04:05",
		},
		{
			name:     "WeekdayShort",
			format:   Time.WeekdayShort,
			expected: "Tue",
		},
		{
			name:     "MonthLong",
			format:   Time.MonthLong,
			expected: "January",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format(ut)
			if result != tt.expected {
				t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}
func TestUTC_ParseRFC3339Nano(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid nano time",
			input:   "2023-01-01T12:00:00.123456789Z",
			want:    time.Date(2023, 1, 1, 12, 0, 0, 123456789, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "2023-01-01 12:00:00",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRFC3339Nano(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRFC3339Nano() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(Time{tt.want}) {
				t.Errorf("ParseRFC3339Nano() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUTC_ScanEdgeCases(t *testing.T) {
	var ut Time
	// Test scanning various time formats
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:    "string RFC3339",
			input:   "2023-01-01T12:00:00Z",
			wantErr: false,
		},
		{
			name:    "time.Time",
			input:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid string format",
			input:   "invalid-time",
			wantErr: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   123,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ut.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestUTC_TimezoneFallbacks(t *testing.T) {
	testTime := Time{time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
	tests := []struct {
		name     string
		got      time.Time
		expected time.Time
	}{
		{
			name:     "PST",
			got:      testTime.PST(),
			expected: testTime.Time.In(time.FixedZone("PST", -8*3600)),
		},
		{
			name:     "MST",
			got:      testTime.MST(),
			expected: testTime.Time.In(time.FixedZone("MST", -7*3600)),
		},
		{
			name:     "CST",
			got:      testTime.CST(),
			expected: testTime.Time.In(time.FixedZone("CST", -6*3600)),
		},
		{
			name:     "EST",
			got:      testTime.EST(),
			expected: testTime.Time.In(time.FixedZone("EST", -5*3600)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.got.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, tt.got)
			}
		})
	}
}
func TestUTC_ValueNilReceiver(t *testing.T) {
	var ut Time // Use zero value instead of nil pointer
	value, err := ut.Value()
	if err != nil {
		t.Errorf("Value() on zero value should not return error, got %v", err)
	}
	if value == nil {
		t.Error("Value() on zero value should not return nil")
	}
}
func TestUTC_ParseMultipleFormats(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "RFC3339 with milliseconds",
			input:   "2023-01-01T12:00:00.123Z",
			wantErr: false,
		},
		{
			name:    "RFC3339 with microseconds",
			input:   "2023-01-01T12:00:00.123456Z",
			wantErr: false,
		},
		{
			name:    "RFC3339 with positive offset",
			input:   "2023-01-01T12:00:00+01:00",
			wantErr: false,
		},
		{
			name:    "RFC3339 with negative offset",
			input:   "2023-01-01T12:00:00-01:00",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRFC3339(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRFC3339() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestUTC_NilHandling(t *testing.T) {
	// Test with nil pointer
	var ut *Time = nil
	// Test MarshalJSON with nil
	data, err := ut.MarshalJSON()
	if err == nil {
		t.Error("MarshalJSON() on nil receiver should return error")
	}
	if data != nil {
		t.Errorf("MarshalJSON() on nil receiver should return nil data, got %v", data)
	}
	// Test String() on nil - should not panic and return "<nil>"
	if s := ut.String(); s != "<nil>" {
		t.Errorf("String() on nil receiver should return \"<nil>\", got %q", s)
	}
	// Test Value() on nil - should return error and nil value
	if v, err := ut.Value(); err == nil || v != nil {
		t.Errorf("Value() on nil receiver should return (nil, error), got (%v, %v)", v, err)
	}
}
func TestUTC_ZeroValueHandling(t *testing.T) {
	// Test with zero value (empty struct) instead of nil pointer
	var ut Time
	// Test MarshalJSON with zero value
	data, err := ut.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() on zero value should not return error, got %v", err)
	}
	if string(data) != `"0001-01-01T00:00:00Z"` {
		t.Errorf("MarshalJSON() on zero value = %s, want %s", string(data), `"0001-01-01T00:00:00Z"`)
	}
	// Test String() with zero value
	str := ut.String()
	if str != "0001-01-01T00:00:00Z" {
		t.Errorf("String() on zero value = %q, want %q", str, "0001-01-01T00:00:00Z")
	}
	// Test Value() with zero value
	val, err := ut.Value()
	if err != nil {
		t.Errorf("Value() on zero value should not return error, got %v", err)
	}
	if val == nil {
		t.Error("Value() on zero value should not return nil")
	}
	// Test IsZero() with zero value
	if !ut.IsZero() {
		t.Error("IsZero() should return true for zero value")
	}
	// Test with non-zero value
	nonZero := Time{time.Now()}
	if nonZero.IsZero() {
		t.Error("IsZero() should return false for non-zero value")
	}
}
func TestUTC_TimeFormatErrors(t *testing.T) {
	ut := Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)}
	// Test with invalid layout
	result := ut.Format("invalid")
	expected := "invalid"
	if result != expected {
		t.Errorf("Format() with invalid layout = %q, want %q", result, expected)
	}
	// Test TimeFormat with empty TimeLayout
	result = ut.TimeFormat("")
	if result != "" {
		t.Errorf("TimeFormat() with empty layout should return empty string, got %s", result)
	}
}
func TestUTC_ParseWithCustomLayout(t *testing.T) {
	tests := []struct {
		name    string
		layout  string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "custom layout success",
			layout:  "2006-01-02",
			input:   "2024-01-02",
			want:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "custom layout failure",
			layout:  "2006-01-02",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty layout",
			layout:  "",
			input:   "2024-01-02",
			wantErr: true,
		},
		{
			name:    "empty input",
			layout:  "2006-01-02",
			input:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.layout, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(Time{tt.want}) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUTC_LocationInitialization(t *testing.T) {
	err := ValidateTimezoneAvailability()
	if err != nil {
		// Verify the error message contains useful information
		if !strings.Contains(err.Error(), "timezone locations") {
			t.Errorf("Expected error message to contain 'timezone locations', got %v", err)
		}
	}
}
func TestUTC_ComparisonEdgeCases(t *testing.T) {
	zero := Time{}
	now := Now()
	// Test comparison with zero value
	if zero.After(now) {
		t.Error("Zero time should not be after now")
	}
	if now.Before(zero) {
		t.Error("Now should not be before zero time")
	}
	if zero.Equal(now) {
		t.Error("Zero time should not equal now")
	}
	// Test comparison with same time
	same := Now()
	if !same.Equal(same) {
		t.Error("Same time should equal itself")
	}
}
func TestUTC_ArithmeticEdgeCases(t *testing.T) {
	now := Now()
	// Test adding/subtracting zero duration
	if !now.Add(0).Equal(now) {
		t.Error("Adding zero duration should return same time")
	}
	// Test adding negative duration
	negDur := -time.Hour
	if !now.Add(negDur).Add(time.Hour).Equal(now) {
		t.Error("Adding negative then positive duration should return to original time")
	}
	// Test subtracting same time
	if now.Sub(now) != 0 {
		t.Error("Subtracting same time should return zero duration")
	}
}
func TestUTC_RFC3339(t *testing.T) {
	ut := Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)}
	expected := "2024-01-02T15:04:05Z"
	result := ut.RFC3339()
	if result != expected {
		t.Errorf("RFC3339() = %q, want %q", result, expected)
	}
}
func TestUTC_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "RFC3339 format",
			input:   "2024-01-02T15:04:05Z",
			want:    time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "RFC3339Nano format",
			input:   "2024-01-02T15:04:05.123456789Z",
			want:    time.Date(2024, 1, 2, 15, 4, 5, 123456789, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date only format",
			input:   "2024-01-02",
			want:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "DateTime format",
			input:   "2024-01-02 15:04:05",
			want:    time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Empty input",
			input:   "",
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "invalid-time",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ut Time
			err := ut.UnmarshalText([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !ut.Time.Equal(tt.want) {
				t.Errorf("UnmarshalText() = %v, want %v", ut.Time, tt.want)
			}
		})
	}
}
func TestUTC_MarshalText(t *testing.T) {
	tests := []struct {
		name     string
		time     Time
		expected string
	}{
		{
			name:     "Normal time",
			time:     Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)},
			expected: "2024-01-02T15:04:05Z",
		},
		{
			name:     "Zero time",
			time:     Time{},
			expected: "0001-01-01T00:00:00Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.time.MarshalText()
			if err != nil {
				t.Errorf("MarshalText() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("MarshalText() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}
func TestUTC_UnixHelpers(t *testing.T) {
	// Test FromUnix
	unixSec := int64(1704215045) // 2024-01-02 17:04:05 UTC
	t1 := FromUnix(unixSec)
	expected1 := time.Unix(unixSec, 0).UTC()
	if !t1.Time.Equal(expected1) {
		t.Errorf("FromUnix() = %v, want %v", t1.Time, expected1)
	}
	// Test FromUnixMilli
	unixMilli := int64(1704215045123) // 2024-01-02 17:04:05.123 UTC
	t2 := FromUnixMilli(unixMilli)
	expected2 := time.Unix(0, unixMilli*int64(time.Millisecond)).UTC()
	if !t2.Time.Equal(expected2) {
		t.Errorf("FromUnixMilli() = %v, want %v", t2.Time, expected2)
	}
	// Test Unix() method
	ut := Time{time.Date(2024, 1, 2, 17, 4, 5, 0, time.UTC)}
	if ut.Unix() != 1704215045 {
		t.Errorf("Unix() = %d, want %d", ut.Unix(), 1704215045)
	}
	// Test UnixMilli() method
	utMilli := Time{time.Date(2024, 1, 2, 17, 4, 5, 123000000, time.UTC)}
	if utMilli.UnixMilli() != 1704215045123 {
		t.Errorf("UnixMilli() = %d, want %d", utMilli.UnixMilli(), 1704215045123)
	}
}
func TestUTC_DayBoundaries(t *testing.T) {
	// Test with a time in the middle of the day
	ut := Time{time.Date(2024, 1, 2, 15, 30, 45, 123456789, time.UTC)}
	// Test StartOfDay
	start := ut.StartOfDay()
	expectedStart := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	if !start.Time.Equal(expectedStart) {
		t.Errorf("StartOfDay() = %v, want %v", start.Time, expectedStart)
	}
	// Test EndOfDay
	end := ut.EndOfDay()
	expectedEnd := time.Date(2024, 1, 2+1, 0, 0, 0, -1, time.UTC) // One nanosecond before next midnight
	if !end.Time.Equal(expectedEnd) {
		t.Errorf("EndOfDay() = %v, want %v", end.Time, expectedEnd)
	}
	// Verify that start is before original time and original time is before end
	if !start.Before(ut) {
		t.Error("StartOfDay should be before original time")
	}
	if !ut.Before(end) {
		t.Error("Original time should be before EndOfDay")
	}
}
func TestUTC_LocationHelpers(t *testing.T) {
	ut := Time{time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)}
	// Test In() with valid timezone
	eastern, err := ut.In("America/New_York")
	if err != nil {
		t.Errorf("In() error = %v", err)
	}
	// Should be 5 hours behind UTC in winter
	if eastern.Hour() != 10 {
		t.Errorf("In('America/New_York') hour = %d, want 10", eastern.Hour())
	}
	// Test In() with invalid timezone
	_, err = ut.In("Invalid/Timezone")
	if err == nil {
		t.Error("In() with invalid timezone should return error")
	}
	// Test InLocation()
	loc, _ := time.LoadLocation("America/Los_Angeles")
	pacific := ut.InLocation(loc)
	// Should be 8 hours behind UTC in winter
	if pacific.Hour() != 7 {
		t.Errorf("InLocation() hour = %d, want 7", pacific.Hour())
	}
}
func TestUTC_ScanEnhanced(t *testing.T) {
	var ut Time
	// Test scanning []byte
	err := ut.Scan([]byte("2024-01-02T15:04:05Z"))
	if err != nil {
		t.Errorf("Scan([]byte) error = %v", err)
	}
	expected := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	if !ut.Time.Equal(expected) {
		t.Errorf("Scan([]byte) = %v, want %v", ut.Time, expected)
	}
	// Test scanning []byte with flexible format
	err = ut.Scan([]byte("2024-01-02"))
	if err != nil {
		t.Errorf("Scan([]byte date only) error = %v", err)
	}
	expectedDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	if !ut.Time.Equal(expectedDate) {
		t.Errorf("Scan([]byte date only) = %v, want %v", ut.Time, expectedDate)
	}
	// Test scanning invalid []byte
	err = ut.Scan([]byte("invalid-time"))
	if err == nil {
		t.Error("Scan(invalid []byte) should return error")
	}
	// Test scanning float64 (unsupported type)
	err = ut.Scan(float64(123.45))
	if err == nil {
		t.Error("Scan(float64) should return error")
	}
}
func TestUTC_InternalParse(t *testing.T) {
	// Test our internal parse function through public APIs
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "RFC3339 with timezone",
			input:   "2024-01-02T15:04:05+02:00",
			want:    time.Date(2024, 1, 2, 13, 4, 5, 0, time.UTC), // Converted to UTC
			wantErr: false,
		},
		{
			name:    "Date-time without timezone (treated as UTC)",
			input:   "2024-01-02 15:04:05",
			want:    time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date only (treated as UTC)",
			input:   "2024-01-02",
			want:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through UnmarshalJSON which uses parse
			var ut Time
			jsonData := `"` + tt.input + `"`
			err := ut.UnmarshalJSON([]byte(jsonData))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !ut.Time.Equal(tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", ut.Time, tt.want)
			}
		})
	}
}
func TestUTC_TimezoneInitErrors(t *testing.T) {
	// Test ValidateTimezoneAvailability when there's no error
	if locationError == nil {
		err := ValidateTimezoneAvailability()
		if err != nil {
			t.Errorf("ValidateTimezoneAvailability() with no location error should return nil, got %v", err)
		}
	}
}
