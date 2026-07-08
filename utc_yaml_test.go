package utc

import (
	"errors"
	"testing"
	"time"
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
			input:   "2023-01-01T12:00:00Z",
			want:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date only format",
			input:   "2023-06-15",
			want:    time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Year-month format",
			input:   "2024-06",
			want:    time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Year only format",
			input:   "2024",
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date time format",
			input:   "2023-01-01 15:30:45",
			want:    time.Date(2023, 1, 1, 15, 30, 45, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "not-a-date",
			wantErr: true,
		},
		{
			name:    "Invalid month",
			input:   "2023-13-01",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ut Time
			err := ut.UnmarshalYAML(func(v any) error {
				target, ok := v.(*string)
				if !ok {
					t.Fatalf("UnmarshalYAML requested %T, want *string", v)
				}
				*target = tt.input
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !ut.UTC().Equal(tt.want) {
				t.Errorf("UnmarshalYAML() = %v, want %v", ut.UTC(), tt.want)
			}
		})
	}
}

func TestUTC_UnmarshalYAMLPropagatesDecoderError(t *testing.T) {
	var ut Time
	want := errors.New("decoder failed")
	err := ut.UnmarshalYAML(func(any) error {
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("UnmarshalYAML() error = %v, want %v", err, want)
	}
}

func TestUTC_MarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		time    Time
		want    any
		wantErr bool
	}{
		{
			name:    "Normal time",
			time:    New(time.Date(2023, 6, 15, 12, 30, 45, 0, time.UTC)),
			want:    "2023-06-15T12:30:45Z",
			wantErr: false,
		},
		{
			name:    "Zero time",
			time:    Time{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Time with nanoseconds",
			time:    New(time.Date(2023, 6, 15, 12, 30, 45, 123456789, time.UTC)),
			want:    "2023-06-15T12:30:45.123456789Z",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.time.MarshalYAML()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("MarshalYAML() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestUTC_YAMLRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		time Time
	}{
		{
			name: "Normal time",
			time: New(time.Date(2023, 6, 15, 12, 30, 45, 0, time.UTC)),
		},
		{
			name: "Date only",
			time: New(time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "Nanoseconds",
			time: New(time.Date(2023, 6, 15, 12, 30, 45, 123456789, time.UTC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := tt.time.MarshalYAML()
			if err != nil {
				t.Fatalf("MarshalYAML() error = %v", err)
			}
			text, ok := marshaled.(string)
			if !ok {
				t.Fatalf("MarshalYAML() = %#v, want string", marshaled)
			}

			var result Time
			err = result.UnmarshalYAML(func(v any) error {
				target, ok := v.(*string)
				if !ok {
					t.Fatalf("UnmarshalYAML requested %T, want *string", v)
				}
				*target = text
				return nil
			})
			if err != nil {
				t.Fatalf("UnmarshalYAML() error = %v", err)
			}

			if !result.UTC().Equal(tt.time.UTC()) {
				t.Errorf("Round trip failed: got %v, want %v", result.UTC(), tt.time.UTC())
			}
		})
	}
}
