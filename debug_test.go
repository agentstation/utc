//go:build debug
// +build debug

package utc

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestDebugLogging verifies that debug logging works when enabled
func TestDebugLogging(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Test nil receiver calls
	var nilTime *Time

	// These should log to stderr in debug mode
	_, _ = nilTime.MarshalJSON()
	_ = nilTime.UnmarshalJSON([]byte(`"2024-01-02T15:04:05Z"`))
	_ = nilTime.UnmarshalText([]byte("2024-01-02T15:04:05Z"))
	_ = nilTime.Scan("2024-01-02T15:04:05Z")

	// Restore stderr and read output
	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify debug logs were written
	expectedLogs := []string{
		"MarshalJSON() called on nil *Time receiver",
		"UnmarshalJSON() called on nil *Time receiver",
		"UnmarshalText() called on nil *Time receiver",
		"Scan() called on nil *Time receiver",
	}

	for _, expected := range expectedLogs {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected debug log containing %q, but got: %s", expected, output)
		}
	}

	// Verify the format includes timestamp and file info
	if !strings.Contains(output, "[UTC DEBUG]") {
		t.Error("Debug logs should include [UTC DEBUG] prefix")
	}
}

// TestDebugLogDirect tests the debugLog function directly
func TestDebugLogDirect(t *testing.T) {
	// This test just verifies debugLog doesn't panic and initializes correctly
	// The actual output verification is done in TestDebugLogging which captures
	// output in a more controlled way by testing real usage patterns

	// This should not panic
	debugLog("test message: %s", "hello")
	debugLog("another test")

	// Verify debugLogger is initialized (it should be after first call)
	if debugLogger == nil {
		t.Error("debugLogger should be initialized after debugLog calls")
	}
}
