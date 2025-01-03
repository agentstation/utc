package utctime

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestUTCTime_UnmarshalJSON(t *testing.T) {
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

func TestUTCTime_MarshalJSON(t *testing.T) {
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

func TestUTCTime_RoundTrip(t *testing.T) {
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

func TestUTCTime_UTCTimeNow(t *testing.T) {
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

func TestUTCTime_ParseUTCTime(t *testing.T) {
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

func TestUTCTime_String(t *testing.T) {
	ut := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	want := "2023-01-01T12:00:00Z"
	if got := ut.String(); got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}

func TestUTCTime_DatabaseOperations(t *testing.T) {
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

func TestUTCTime_Comparisons(t *testing.T) {
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

func TestUTCTime_Arithmetic(t *testing.T) {
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

func TestUTCTime_IsZero(t *testing.T) {
	var zero Time
	if !zero.IsZero() {
		t.Error("IsZero() should be true for zero value")
	}

	nonZero := Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)}
	if nonZero.IsZero() {
		t.Error("IsZero() should be false for non-zero value")
	}
}

func TestUTCTime_UnixTimestamps(t *testing.T) {
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

func TestNewUTCTime(t *testing.T) {
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
				t.Errorf("NewUTCTime() = %v, want %v", got.UTC(), tt.want)
			}
			if got.UTC().Location() != time.UTC {
				t.Error("NewUTCTime() time is not in UTC")
			}
		})
	}
}

func TestUTCTime_TimeZoneConversions(t *testing.T) {
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

func TestUTCTime_Formatting(t *testing.T) {
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

func TestUTCTime_TimezoneError(t *testing.T) {
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
