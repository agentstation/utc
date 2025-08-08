//go:build !debug
// +build !debug

package utc

// debugLog is a no-op in non-debug builds
func debugLog(format string, v ...any) {
	// No-op in production builds
}
