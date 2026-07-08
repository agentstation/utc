package yamlintegration_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/agentstation/utc"
	"github.com/goccy/go-yaml"
)

type event struct {
	Timestamp utc.Time `yaml:"timestamp"`
	Optional  utc.Time `yaml:"optional"`
}

func TestGoccyYAMLRoundTripPreservesUTCAndNanoseconds(t *testing.T) {
	original := event{
		Timestamp: utc.New(time.Date(2024, 1, 2, 3, 4, 5, 123456789, time.FixedZone("EST", -5*3600))),
	}

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}
	if !bytes.Contains(data, []byte("2024-01-02T08:04:05.123456789Z")) {
		t.Fatalf("yaml.Marshal() = %q, want UTC RFC3339Nano timestamp", data)
	}
	if !bytes.Contains(data, []byte("optional: null")) {
		t.Fatalf("yaml.Marshal() = %q, want zero utc.Time encoded as null", data)
	}

	var decoded event
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}
	if !decoded.Timestamp.UTC().Equal(original.Timestamp.UTC()) {
		t.Fatalf("decoded timestamp = %v, want %v", decoded.Timestamp.UTC(), original.Timestamp.UTC())
	}
	if !decoded.Optional.IsZero() {
		t.Fatalf("decoded optional = %v, want zero", decoded.Optional)
	}
}

func TestGoccyYAMLUnmarshalFlexibleScalar(t *testing.T) {
	var decoded event
	if err := yaml.Unmarshal([]byte("timestamp: 2024-06\n"), &decoded); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}
	want := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if !decoded.Timestamp.UTC().Equal(want) {
		t.Fatalf("decoded timestamp = %v, want %v", decoded.Timestamp.UTC(), want)
	}
}

func TestGoccyYAMLUnmarshalInvalidScalar(t *testing.T) {
	var decoded event
	if err := yaml.Unmarshal([]byte("timestamp: not-a-date\n"), &decoded); err == nil {
		t.Fatal("yaml.Unmarshal() unexpectedly succeeded")
	}
}
