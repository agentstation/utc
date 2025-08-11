//go:build yaml
// +build yaml

package utc

import (
	"testing"
	"time"

	"github.com/goccy/go-yaml"
)

func TestUTC_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "RFC3339 format",
			input:   `"2023-01-01T12:00:00Z"`,
			want:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date only format",
			input:   `"2023-06-15"`,
			want:    time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Year-month format",
			input:   `"2024-06"`,
			want:    time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Year only format",
			input:   `"2024"`,
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date time format",
			input:   `"2023-01-01 15:30:45"`,
			want:    time.Date(2023, 1, 1, 15, 30, 45, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Empty string",
			input:   `""`,
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   `"not-a-date"`,
			wantErr: true,
		},
		{
			name:    "Invalid month",
			input:   `"2023-13-01"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ut Time
			err := yaml.Unmarshal([]byte(tt.input), &ut)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !ut.Time.Equal(tt.want) {
				t.Errorf("UnmarshalYAML() = %v, want %v", ut.Time, tt.want)
			}
		})
	}
}

func TestUTC_MarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		time    Time
		want    string
		wantErr bool
	}{
		{
			name:    "Normal time",
			time:    Time{time.Date(2023, 6, 15, 12, 30, 45, 0, time.UTC)},
			want:    "\"2023-06-15T12:30:45Z\"\n",
			wantErr: false,
		},
		{
			name:    "Zero time",
			time:    Time{time.Time{}},
			want:    "null\n",
			wantErr: false,
		},
		{
			name:    "Time with nanoseconds",
			time:    Time{time.Date(2023, 6, 15, 12, 30, 45, 123456789, time.UTC)},
			want:    "\"2023-06-15T12:30:45Z\"\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(data) != tt.want {
				t.Errorf("MarshalYAML() = %q, want %q", string(data), tt.want)
			}
		})
	}
}

func TestUTC_YAMLRoundTrip(t *testing.T) {
	// Test that we can marshal and unmarshal back to the same value
	tests := []struct {
		name string
		time Time
	}{
		{
			name: "Normal time",
			time: Time{time.Date(2023, 6, 15, 12, 30, 45, 0, time.UTC)},
		},
		{
			name: "Date only",
			time: Time{time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
		},
		{
			name: "Year month",
			time: Time{time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to YAML
			data, err := yaml.Marshal(tt.time)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			// Unmarshal back
			var result Time
			err = yaml.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			// Compare (ignoring nanoseconds which might be lost in formatting)
			if !result.Time.Truncate(time.Second).Equal(tt.time.Time.Truncate(time.Second)) {
				t.Errorf("Round trip failed: got %v, want %v", result.Time, tt.time.Time)
			}
		})
	}
}
