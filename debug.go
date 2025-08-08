//go:build debug
// +build debug

package utc

import (
	"log"
	"os"
	"sync"
)

// debugLogger is only available in debug builds
var (
	debugLogger *log.Logger
	debugOnce   sync.Once
)

func initDebugLogger() {
	debugOnce.Do(func() {
		debugLogger = log.New(os.Stderr, "[UTC DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	})
}

func debugLog(format string, v ...any) {
	initDebugLogger()
	debugLogger.Printf(format, v...)
}
