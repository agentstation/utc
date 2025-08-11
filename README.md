# utc

```sh
                                _            _    _                    
                         _   _ | |_  ___    | |_ (_)_ __ ___   ___    
                        | | | || __|/ __|   | __|| | '_ ` _ \ / _ \   
                        | |_| || |_| (__    | |_ | | | | | | |  __/   
                         \___/  \__|\___|    \__||_|_| |_| |_|\___|   
```

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/agentstation/utc)
[![Go Report Card](https://goreportcard.com/badge/github.com/agentstation/utc?style=flat-square)](https://goreportcard.com/report/github.com/agentstation/utc)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/agentstation/utc/ci.yaml?style=flat-square)](https://github.com/agentstation/utc/actions)
[![codecov](https://codecov.io/gh/agentstation/utc/graph/badge.svg?token=EOAZUVVH7H)](https://codecov.io/gh/agentstation/utc)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/agentstation/utc/master/LICENSE)

The `utc` package provides an enhanced, **zero-dependency** wrapper around Go's `time.Time` that ensures your times are consistently in UTC while adding powerful convenience methods for real-world applications.

## Key Features 🌟

### **🛡️ Safety & Reliability**
- **Nil-safe operations** - No more panic on nil receivers
- **Guaranteed UTC storage** - Eliminates timezone confusion
- **Race condition tested** - Safe for concurrent applications
- **Comprehensive error handling** - Graceful failures instead of crashes

### **🎯 Developer Productivity**  
- **Rich formatting options** - US/EU dates, RFC standards, custom layouts
- **Automatic timezone handling** - PST/PDT, EST/EDT transitions
- **Flexible parsing** - Handles multiple input formats automatically
- **Serialization ready** - JSON, YAML, Database, Text encoding

### **⚡ Performance & Compatibility**
- **Zero dependencies** - No external packages required
- **Lightweight footprint** - Minimal impact on your binary size
- **Go 1.18+ compatible** - Works with modern and legacy Go versions
- **Drop-in replacement** - Compatible with standard `time.Time` methods

## Installation

To install the `utc` package, use the following command:

```sh
go get github.com/agentstation/utc
```

**Requirements**: Go 1.18 or later

**YAML Usage**: The package includes full YAML support. To use YAML functionality, install a YAML library:
```bash
go get github.com/goccy/go-yaml  # Recommended YAML library
```

## Why Choose UTC? 🚀

### ✅ **Zero Dependencies**
- **No external dependencies** for core functionality
- **Lightweight** - adds minimal footprint to your project
- **Fast installation** - `go get` with no dependency resolution delays

### ✅ **Maximum Compatibility**  
- **Go 1.18+** support (broader than most time libraries)
- **Cross-platform** - works on all Go-supported platforms
- **Future-proof** - extensively tested across Go 1.18-1.24

### ✅ **Production-Ready Safety**
- **Nil-safe operations** - methods return errors instead of panicking
- **Race condition tested** - safe for concurrent use
- **Comprehensive test coverage** - battle-tested with 100+ test cases

### ✅ **Developer Experience**
- **Intuitive API** - familiar time.Time wrapper with enhanced functionality  
- **Rich formatting options** - US/EU date formats, RFC standards, custom layouts
- **Automatic timezone handling** - PST/PDT, EST/EDT transitions handled correctly
- **Serialization ready** - JSON, YAML, Database, Text encoding

### ✅ **Optional Advanced Features**
- **Debug mode** - development-time nil pointer detection  
- **Flexible parsing** - handles multiple time formats automatically

## Quick Start

Get up and running in seconds:

```go
package main

import (
    "fmt"
    "github.com/agentstation/utc"
)

func main() {
    // Get current time in UTC
    now := utc.Now()
    fmt.Println("Current UTC time:", now.RFC3339())
    
    // Parse and convert to different formats
    t, _ := utc.ParseRFC3339("2024-01-15T10:30:00Z")
    fmt.Println("US format:", t.USDateShort())     // "01/15/2024"
    fmt.Println("EU format:", t.EUDateShort())     // "15/01/2024" 
    fmt.Println("Pacific time:", t.Pacific())      // Auto PST/PDT
}
```

## UTC vs Standard Library ⚡

See the difference between `utc` and Go's standard `time` package:

| Feature | Standard `time.Time` | `utc.Time` |
|---------|---------------------|------------|
| **Timezone Safety** | ❌ Manual timezone handling | ✅ **Always UTC, automatic conversion** |
| **Nil Safety** | ❌ Panics on nil receiver | ✅ **Returns errors gracefully** |
| **Dependencies** | ✅ Zero deps | ✅ **Zero deps (core)** |
| **Rich Formatting** | ❌ Manual layout strings | ✅ **Built-in US/EU/ISO formats** |
| **Timezone Conversion** | ❌ Manual location loading | ✅ **Auto PST/PDT, EST/EDT handling** |
| **JSON Support** | ✅ Basic marshal/unmarshal | ✅ **Enhanced parsing & formatting** |
| **YAML Support** | ❌ No built-in support | ✅ **Full YAML marshal/unmarshal** |
| **Text Encoding** | ❌ Limited support | ✅ **Full MarshalText/UnmarshalText** |
| **Database Ready** | ✅ Basic support | ✅ **Enhanced Scan/Value methods** |
| **Unix Timestamps** | ✅ Basic Unix() method | ✅ **Unix + UnixMilli helpers** |
| **Day Boundaries** | ❌ Manual calculation | ✅ **StartOfDay/EndOfDay methods** |
| **Production Safety** | ❌ Can panic unexpectedly | ✅ **Error-first design** |
| **Debug Support** | ❌ No debugging aids | ✅ **Optional debug mode** |

**Before** (standard library):
```go
loc, _ := time.LoadLocation("America/New_York")
t := time.Now().In(loc)  // Hope timezone exists!
if t != nil {            // Manual nil checking
    fmt.Println(t.Format("01/02/2006"))  // Remember layout
}
```

**After** (with UTC):
```go
t := utc.Now()
fmt.Println(t.Eastern().Format(time.Kitchen))  // Auto EST/EDT
fmt.Println(t.USDateShort())                   // "01/15/2024"
// No panics, no manual timezone loading, no layout memorization!
```

## Detailed Usage

1. Import the package:

```go
import "github.com/agentstation/utc"
```

2. Create a new UTC time:

```go
// Get current time in UTC
now := utc.Now()

// Convert existing time.Time to UTC
myTime := utc.New(someTime)

// Parse a time string
t, err := utc.ParseRFC3339("2023-01-01T12:00:00Z")
```

3. Format times using various layouts:

```go
t := utc.Now()

// US formats
fmt.Println(t.USDateShort())     // "01/02/2024"
fmt.Println(t.USDateTime12())    // "01/02/2024 03:04:05 PM"

// EU formats
fmt.Println(t.EUDateShort())     // "02/01/2024"
fmt.Println(t.EUDateTime24())    // "02/01/2024 15:04:05"

// ISO/RFC formats
fmt.Println(t.RFC3339())         // "2024-01-02T15:04:05Z"
fmt.Println(t.ISO8601())         // "2024-01-02T15:04:05Z"

// Components
fmt.Println(t.WeekdayLong())     // "Tuesday"
fmt.Println(t.MonthShort())      // "Jan"
```

4. Convert between timezones:

```go
t := utc.Now()

// Get time in different US timezones
pacific := t.Pacific()   // Handles PST/PDT automatically
eastern := t.Eastern()   // Handles EST/EDT automatically
central := t.Central()   // Handles CST/CDT automatically
mountain := t.Mountain() // Handles MST/MDT automatically
```

5. Serialization and Database operations:

```go
// JSON marshaling
type Event struct {
    StartTime utc.Time `json:"start_time"`
    EndTime   utc.Time `json:"end_time"`
}

// YAML marshaling (requires a YAML library like go-yaml)
type Config struct {
    StartTime utc.Time `yaml:"start_time"`
    EndTime   utc.Time `yaml:"end_time"`
}

// Database operations
type Record struct {
    CreatedAt utc.Time `db:"created_at"`
    UpdatedAt utc.Time `db:"updated_at"`
}
```

## YAML Support

The package includes full YAML marshaling/unmarshaling support through `MarshalYAML`/`UnmarshalYAML` methods that implement the standard YAML interfaces. Works with any Go YAML library that follows these interfaces.

**Requirements for YAML testing**:
- Go 1.21.0+ (required by go-yaml dependency)
- `github.com/goccy/go-yaml` package

```sh
# Install the YAML testing dependency
go get github.com/goccy/go-yaml@v1.18.0

# Run all tests including YAML functionality
go test -tags yaml ./...

# Run only YAML-specific tests
go test -tags yaml -run YAML ./...
```

**Note**: The YAML marshal/unmarshal methods are available in the main package, but the actual YAML processing requires the go-yaml dependency. Most users won't need YAML functionality for production use.

## Testing

The project includes comprehensive Makefile targets for testing:

```sh
# Run core tests (Go 1.18+, no dependencies)
make test

# Run tests with YAML support (requires Go 1.21.0+)
make test-yaml

# Run all tests (core + YAML)
make test-all

# Generate coverage reports
make coverage         # Core tests only
make coverage-yaml    # Include YAML tests
make coverage-all     # Both coverage reports
```

## Debug Mode

The package includes a debug mode that helps identify potential bugs during development:

```sh
# Build with debug mode enabled
go build -tags debug

# Run tests with debug mode
go test -tags debug ./...
```

When debug mode is enabled, the package logs warnings when methods are called on nil receivers:

```
[UTC DEBUG] 2024/01/02 15:04:05 debug.go:26: String() called on nil *Time receiver
[UTC DEBUG] 2024/01/02 15:04:05 debug.go:26: Value() called on nil *Time receiver
```

## Additional Utilities

The package includes several convenience methods:

```go
// Unix timestamp conversions
t1 := utc.FromUnix(1704199445)           // From Unix seconds
t2 := utc.FromUnixMilli(1704199445000)   // From Unix milliseconds
seconds := t.Unix()                       // To Unix seconds
millis := t.UnixMilli()                   // To Unix milliseconds

// Day boundaries
start := t.StartOfDay()  // 2024-01-02 00:00:00.000000000 UTC
end := t.EndOfDay()      // 2024-01-02 23:59:59.999999999 UTC

// Generic timezone conversion
eastern, err := t.In("America/New_York")
tokyo, err := t.In("Asia/Tokyo")
```

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# utc

```go
import "github.com/agentstation/utc"
```

Package utc provides a time.Time wrapper that ensures all times are in UTC.

The package offers enhanced safety by gracefully handling nil receivers instead of panicking, making it more suitable for production environments. When compiled with the debug build tag \(\-tags debug\), it provides additional logging for nil receiver method calls to help identify potential bugs during development.

Key features:

- All times are automatically converted to and stored in UTC
- JSON marshaling/unmarshaling with flexible parsing
- Full YAML marshaling/unmarshaling support
- SQL database compatibility with enhanced type support
- Timezone conversion helpers with automatic DST handling
- Extensive formatting options for US and EU date formats
- Nil\-safe operations that return errors instead of panicking

Debug mode:

```
To enable debug logging, compile with: go build -tags debug
This will log warnings when methods are called on nil receivers.
```

## Index

- [func ValidateTimezoneAvailability\(\) error](<#ValidateTimezoneAvailability>)
- [type Time](<#Time>)
  - [func FromUnix\(sec int64\) Time](<#FromUnix>)
  - [func FromUnixMilli\(ms int64\) Time](<#FromUnixMilli>)
  - [func New\(t time.Time\) Time](<#New>)
  - [func Now\(\) Time](<#Now>)
  - [func Parse\(layout string, s string\) \(Time, error\)](<#Parse>)
  - [func ParseRFC3339\(s string\) \(Time, error\)](<#ParseRFC3339>)
  - [func ParseRFC3339Nano\(s string\) \(Time, error\)](<#ParseRFC3339Nano>)
  - [func \(t Time\) ANSIC\(\) string](<#Time.ANSIC>)
  - [func \(t Time\) Add\(d time.Duration\) Time](<#Time.Add>)
  - [func \(t Time\) After\(u Time\) bool](<#Time.After>)
  - [func \(t Time\) Before\(u Time\) bool](<#Time.Before>)
  - [func \(t Time\) CST\(\) time.Time](<#Time.CST>)
  - [func \(t Time\) Central\(\) time.Time](<#Time.Central>)
  - [func \(t Time\) DateOnly\(\) string](<#Time.DateOnly>)
  - [func \(t Time\) EST\(\) time.Time](<#Time.EST>)
  - [func \(t Time\) EUDateLong\(\) string](<#Time.EUDateLong>)
  - [func \(t Time\) EUDateShort\(\) string](<#Time.EUDateShort>)
  - [func \(t Time\) EUDateTime12\(\) string](<#Time.EUDateTime12>)
  - [func \(t Time\) EUDateTime24\(\) string](<#Time.EUDateTime24>)
  - [func \(t Time\) EUTime12\(\) string](<#Time.EUTime12>)
  - [func \(t Time\) EUTime24\(\) string](<#Time.EUTime24>)
  - [func \(t Time\) Eastern\(\) time.Time](<#Time.Eastern>)
  - [func \(t Time\) EndOfDay\(\) Time](<#Time.EndOfDay>)
  - [func \(t Time\) Equal\(u Time\) bool](<#Time.Equal>)
  - [func \(t Time\) Format\(layout string\) string](<#Time.Format>)
  - [func \(t Time\) ISO8601\(\) string](<#Time.ISO8601>)
  - [func \(t Time\) In\(name string\) \(time.Time, error\)](<#Time.In>)
  - [func \(t Time\) InLocation\(loc \*time.Location\) time.Time](<#Time.InLocation>)
  - [func \(t Time\) IsZero\(\) bool](<#Time.IsZero>)
  - [func \(t Time\) Kitchen\(\) string](<#Time.Kitchen>)
  - [func \(t Time\) MST\(\) time.Time](<#Time.MST>)
  - [func \(t \*Time\) MarshalJSON\(\) \(\[\]byte, error\)](<#Time.MarshalJSON>)
  - [func \(t Time\) MarshalText\(\) \(\[\]byte, error\)](<#Time.MarshalText>)
  - [func \(t Time\) MarshalYAML\(\) \(any, error\)](<#Time.MarshalYAML>)
  - [func \(t Time\) MonthLong\(\) string](<#Time.MonthLong>)
  - [func \(t Time\) MonthShort\(\) string](<#Time.MonthShort>)
  - [func \(t Time\) Mountain\(\) time.Time](<#Time.Mountain>)
  - [func \(t Time\) PST\(\) time.Time](<#Time.PST>)
  - [func \(t Time\) Pacific\(\) time.Time](<#Time.Pacific>)
  - [func \(t Time\) RFC3339\(\) string](<#Time.RFC3339>)
  - [func \(t Time\) RFC3339Nano\(\) string](<#Time.RFC3339Nano>)
  - [func \(t Time\) RFC822\(\) string](<#Time.RFC822>)
  - [func \(t Time\) RFC822Z\(\) string](<#Time.RFC822Z>)
  - [func \(t Time\) RFC850\(\) string](<#Time.RFC850>)
  - [func \(t \*Time\) Scan\(value any\) error](<#Time.Scan>)
  - [func \(t Time\) StartOfDay\(\) Time](<#Time.StartOfDay>)
  - [func \(t \*Time\) String\(\) string](<#Time.String>)
  - [func \(t Time\) Sub\(u Time\) time.Duration](<#Time.Sub>)
  - [func \(t Time\) TimeFormat\(layout TimeLayout\) string](<#Time.TimeFormat>)
  - [func \(t Time\) TimeOnly\(\) string](<#Time.TimeOnly>)
  - [func \(t Time\) USDateLong\(\) string](<#Time.USDateLong>)
  - [func \(t Time\) USDateShort\(\) string](<#Time.USDateShort>)
  - [func \(t Time\) USDateTime12\(\) string](<#Time.USDateTime12>)
  - [func \(t Time\) USDateTime24\(\) string](<#Time.USDateTime24>)
  - [func \(t Time\) USTime12\(\) string](<#Time.USTime12>)
  - [func \(t Time\) USTime24\(\) string](<#Time.USTime24>)
  - [func \(t Time\) UTC\(\) time.Time](<#Time.UTC>)
  - [func \(t Time\) Unix\(\) int64](<#Time.Unix>)
  - [func \(t Time\) UnixMilli\(\) int64](<#Time.UnixMilli>)
  - [func \(t \*Time\) UnmarshalJSON\(data \[\]byte\) error](<#Time.UnmarshalJSON>)
  - [func \(t \*Time\) UnmarshalText\(text \[\]byte\) error](<#Time.UnmarshalText>)
  - [func \(t \*Time\) UnmarshalYAML\(unmarshal func\(any\) error\) error](<#Time.UnmarshalYAML>)
  - [func \(t \*Time\) Value\(\) \(driver.Value, error\)](<#Time.Value>)
  - [func \(t Time\) WeekdayLong\(\) string](<#Time.WeekdayLong>)
  - [func \(t Time\) WeekdayShort\(\) string](<#Time.WeekdayShort>)
- [type TimeLayout](<#TimeLayout>)


<a name="ValidateTimezoneAvailability"></a>
## func [ValidateTimezoneAvailability](<https://github.com/agentstation/utc/blob/master/utc.go#L106>)

```go
func ValidateTimezoneAvailability() error
```

ValidateTimezoneAvailability checks if all timezone locations were properly initialized Returns nil if initialization was successful, otherwise returns the initialization error

<a name="Time"></a>
## type [Time](<https://github.com/agentstation/utc/blob/master/utc.go#L114-L116>)

Time is an alias for time.Time that defaults to UTC time.

```go
type Time struct {
    time.Time
}
```

<a name="FromUnix"></a>
### func [FromUnix](<https://github.com/agentstation/utc/blob/master/utc.go#L566>)

```go
func FromUnix(sec int64) Time
```

Unix helpers

<a name="FromUnixMilli"></a>
### func [FromUnixMilli](<https://github.com/agentstation/utc/blob/master/utc.go#L567>)

```go
func FromUnixMilli(ms int64) Time
```



<a name="New"></a>
### func [New](<https://github.com/agentstation/utc/blob/master/utc.go#L124>)

```go
func New(t time.Time) Time
```

New returns a new Time from a time.Time

<a name="Now"></a>
### func [Now](<https://github.com/agentstation/utc/blob/master/utc.go#L119>)

```go
func Now() Time
```

Now returns the current time in UTC

<a name="Parse"></a>
### func [Parse](<https://github.com/agentstation/utc/blob/master/utc.go#L147>)

```go
func Parse(layout string, s string) (Time, error)
```

Parse parses a time string in the specified format and returns a utc.Time

<a name="ParseRFC3339"></a>
### func [ParseRFC3339](<https://github.com/agentstation/utc/blob/master/utc.go#L129>)

```go
func ParseRFC3339(s string) (Time, error)
```

ParseRFC3339 parses a time string in RFC3339 format and returns a utc.Time

<a name="ParseRFC3339Nano"></a>
### func [ParseRFC3339Nano](<https://github.com/agentstation/utc/blob/master/utc.go#L138>)

```go
func ParseRFC3339Nano(s string) (Time, error)
```

ParseRFC3339Nano parses a time string in RFC3339Nano format and returns a utc.Time

<a name="Time.ANSIC"></a>
### func \(Time\) [ANSIC](<https://github.com/agentstation/utc/blob/master/utc.go#L441>)

```go
func (t Time) ANSIC() string
```

ANSIC formats time as "Mon Jan \_2 15:04:05 2006"

<a name="Time.Add"></a>
### func \(Time\) [Add](<https://github.com/agentstation/utc/blob/master/utc.go#L326>)

```go
func (t Time) Add(d time.Duration) Time
```

Add returns the time t\+d

<a name="Time.After"></a>
### func \(Time\) [After](<https://github.com/agentstation/utc/blob/master/utc.go#L316>)

```go
func (t Time) After(u Time) bool
```

After reports whether the time is after u

<a name="Time.Before"></a>
### func \(Time\) [Before](<https://github.com/agentstation/utc/blob/master/utc.go#L311>)

```go
func (t Time) Before(u Time) bool
```

Before reports whether the time is before u

<a name="Time.CST"></a>
### func \(Time\) [CST](<https://github.com/agentstation/utc/blob/master/utc.go#L351>)

```go
func (t Time) CST() time.Time
```

CST returns t in CST

<a name="Time.Central"></a>
### func \(Time\) [Central](<https://github.com/agentstation/utc/blob/master/utc.go#L377>)

```go
func (t Time) Central() time.Time
```

Central returns t in Central time \(handles CST/CDT automatically\)

<a name="Time.DateOnly"></a>
### func \(Time\) [DateOnly](<https://github.com/agentstation/utc/blob/master/utc.go#L535>)

```go
func (t Time) DateOnly() string
```

DateOnly formats time as "2006\-01\-02"

<a name="Time.EST"></a>
### func \(Time\) [EST](<https://github.com/agentstation/utc/blob/master/utc.go#L346>)

```go
func (t Time) EST() time.Time
```

EST returns t in EST

<a name="Time.EUDateLong"></a>
### func \(Time\) [EUDateLong](<https://github.com/agentstation/utc/blob/master/utc.go#L487>)

```go
func (t Time) EUDateLong() string
```

EUDateLong formats time as "2 January 2006"

<a name="Time.EUDateShort"></a>
### func \(Time\) [EUDateShort](<https://github.com/agentstation/utc/blob/master/utc.go#L482>)

```go
func (t Time) EUDateShort() string
```

EUDateShort formats time as "02/01/2006"

<a name="Time.EUDateTime12"></a>
### func \(Time\) [EUDateTime12](<https://github.com/agentstation/utc/blob/master/utc.go#L492>)

```go
func (t Time) EUDateTime12() string
```

EUDateTime12 formats time as "02/01/2006 03:04:05 PM"

<a name="Time.EUDateTime24"></a>
### func \(Time\) [EUDateTime24](<https://github.com/agentstation/utc/blob/master/utc.go#L497>)

```go
func (t Time) EUDateTime24() string
```

EUDateTime24 formats time as "02/01/2006 15:04:05"

<a name="Time.EUTime12"></a>
### func \(Time\) [EUTime12](<https://github.com/agentstation/utc/blob/master/utc.go#L502>)

```go
func (t Time) EUTime12() string
```

EUTime12 formats time as "3:04 PM"

<a name="Time.EUTime24"></a>
### func \(Time\) [EUTime24](<https://github.com/agentstation/utc/blob/master/utc.go#L507>)

```go
func (t Time) EUTime24() string
```

EUTime24 formats time as "15:04"

<a name="Time.Eastern"></a>
### func \(Time\) [Eastern](<https://github.com/agentstation/utc/blob/master/utc.go#L369>)

```go
func (t Time) Eastern() time.Time
```

Eastern returns t in Eastern time \(handles EST/EDT automatically\)

<a name="Time.EndOfDay"></a>
### func \(Time\) [EndOfDay](<https://github.com/agentstation/utc/blob/master/utc.go#L577>)

```go
func (t Time) EndOfDay() Time
```



<a name="Time.Equal"></a>
### func \(Time\) [Equal](<https://github.com/agentstation/utc/blob/master/utc.go#L321>)

```go
func (t Time) Equal(u Time) bool
```

Equal reports whether t and u represent the same time instant

<a name="Time.Format"></a>
### func \(Time\) [Format](<https://github.com/agentstation/utc/blob/master/utc.go#L398>)

```go
func (t Time) Format(layout string) string
```

Format formats the time using the specified layout

<a name="Time.ISO8601"></a>
### func \(Time\) [ISO8601](<https://github.com/agentstation/utc/blob/master/utc.go#L421>)

```go
func (t Time) ISO8601() string
```

ISO8601 formats time as "2006\-01\-02T15:04:05Z07:00" \(same as RFC3339\)

<a name="Time.In"></a>
### func \(Time\) [In](<https://github.com/agentstation/utc/blob/master/utc.go#L552>)

```go
func (t Time) In(name string) (time.Time, error)
```

In converts time to a named location \(e.g., "America/Los\_Angeles"\).

<a name="Time.InLocation"></a>
### func \(Time\) [InLocation](<https://github.com/agentstation/utc/blob/master/utc.go#L561>)

```go
func (t Time) InLocation(loc *time.Location) time.Time
```

InLocation converts time to a provided \*time.Location.

<a name="Time.IsZero"></a>
### func \(Time\) [IsZero](<https://github.com/agentstation/utc/blob/master/utc.go#L393>)

```go
func (t Time) IsZero() bool
```

Add the useful utility methods while maintaining chainability

<a name="Time.Kitchen"></a>
### func \(Time\) [Kitchen](<https://github.com/agentstation/utc/blob/master/utc.go#L545>)

```go
func (t Time) Kitchen() string
```

Kitchen formats time as "3:04PM"

<a name="Time.MST"></a>
### func \(Time\) [MST](<https://github.com/agentstation/utc/blob/master/utc.go#L356>)

```go
func (t Time) MST() time.Time
```

MST returns t in MST

<a name="Time.MarshalJSON"></a>
### func \(\*Time\) [MarshalJSON](<https://github.com/agentstation/utc/blob/master/utc.go#L186>)

```go
func (t *Time) MarshalJSON() ([]byte, error)
```

MarshalJSON implements the json.Marshaler interface for utc.Time. Returns an error for nil receivers to maintain consistency with standard marshaling behavior.

<a name="Time.MarshalText"></a>
### func \(Time\) [MarshalText](<https://github.com/agentstation/utc/blob/master/utc.go#L201>)

```go
func (t Time) MarshalText() ([]byte, error)
```

MarshalText implements encoding.TextMarshaler.

<a name="Time.MarshalYAML"></a>
### func \(Time\) [MarshalYAML](<https://github.com/agentstation/utc/blob/master/utc.go#L243>)

```go
func (t Time) MarshalYAML() (any, error)
```

MarshalYAML implements the yaml.Marshaler interface for utc.Time

<a name="Time.MonthLong"></a>
### func \(Time\) [MonthLong](<https://github.com/agentstation/utc/blob/master/utc.go#L525>)

```go
func (t Time) MonthLong() string
```

MonthLong formats time as "January"

<a name="Time.MonthShort"></a>
### func \(Time\) [MonthShort](<https://github.com/agentstation/utc/blob/master/utc.go#L530>)

```go
func (t Time) MonthShort() string
```

MonthShort formats time as "Jan"

<a name="Time.Mountain"></a>
### func \(Time\) [Mountain](<https://github.com/agentstation/utc/blob/master/utc.go#L385>)

```go
func (t Time) Mountain() time.Time
```

Mountain returns t in Mountain time \(handles MST/MDT automatically\)

<a name="Time.PST"></a>
### func \(Time\) [PST](<https://github.com/agentstation/utc/blob/master/utc.go#L341>)

```go
func (t Time) PST() time.Time
```

PST returns t in PST

<a name="Time.Pacific"></a>
### func \(Time\) [Pacific](<https://github.com/agentstation/utc/blob/master/utc.go#L361>)

```go
func (t Time) Pacific() time.Time
```

Pacific returns t in Pacific time \(handles PST/PDT automatically\)

<a name="Time.RFC3339"></a>
### func \(Time\) [RFC3339](<https://github.com/agentstation/utc/blob/master/utc.go#L411>)

```go
func (t Time) RFC3339() string
```

RFC3339 formats time as "2006\-01\-02T15:04:05Z07:00"

<a name="Time.RFC3339Nano"></a>
### func \(Time\) [RFC3339Nano](<https://github.com/agentstation/utc/blob/master/utc.go#L416>)

```go
func (t Time) RFC3339Nano() string
```

RFC3339Nano formats time as "2006\-01\-02T15:04:05.999999999Z07:00"

<a name="Time.RFC822"></a>
### func \(Time\) [RFC822](<https://github.com/agentstation/utc/blob/master/utc.go#L426>)

```go
func (t Time) RFC822() string
```

RFC822 formats time as "02 Jan 06 15:04 MST"

<a name="Time.RFC822Z"></a>
### func \(Time\) [RFC822Z](<https://github.com/agentstation/utc/blob/master/utc.go#L431>)

```go
func (t Time) RFC822Z() string
```

RFC822Z formats time as "02 Jan 06 15:04 \-0700"

<a name="Time.RFC850"></a>
### func \(Time\) [RFC850](<https://github.com/agentstation/utc/blob/master/utc.go#L436>)

```go
func (t Time) RFC850() string
```

RFC850 formats time as "Monday, 02\-Jan\-06 15:04:05 MST"

<a name="Time.Scan"></a>
### func \(\*Time\) [Scan](<https://github.com/agentstation/utc/blob/master/utc.go#L282>)

```go
func (t *Time) Scan(value any) error
```

Scan implements the sql.Scanner interface for database operations for utc.Time It does this by scanning the value into a time.Time, converting the time.Time to UTC, and then assigning the UTC time to the utc.Time.

<a name="Time.StartOfDay"></a>
### func \(Time\) [StartOfDay](<https://github.com/agentstation/utc/blob/master/utc.go#L572>)

```go
func (t Time) StartOfDay() Time
```

Day helpers \- times are always in UTC within this package

<a name="Time.String"></a>
### func \(\*Time\) [String](<https://github.com/agentstation/utc/blob/master/utc.go#L257>)

```go
func (t *Time) String() string
```

String implements the Stringer interface for utc.Time. It prints the time in RFC3339 format.

Unlike many Go types that panic on nil receivers, this method returns "\<nil\>" to match stdlib conventions \(e.g., bytes.Buffer\) and improve production safety. In debug builds \(compiled with \-tags debug\), nil receivers are logged to help identify potential bugs.

<a name="Time.Sub"></a>
### func \(Time\) [Sub](<https://github.com/agentstation/utc/blob/master/utc.go#L331>)

```go
func (t Time) Sub(u Time) time.Duration
```

Sub returns the duration t\-u

<a name="Time.TimeFormat"></a>
### func \(Time\) [TimeFormat](<https://github.com/agentstation/utc/blob/master/utc.go#L403>)

```go
func (t Time) TimeFormat(layout TimeLayout) string
```

TimeFormat formats the time using the specified layout

<a name="Time.TimeOnly"></a>
### func \(Time\) [TimeOnly](<https://github.com/agentstation/utc/blob/master/utc.go#L540>)

```go
func (t Time) TimeOnly() string
```

TimeOnly formats time as "15:04:05"

<a name="Time.USDateLong"></a>
### func \(Time\) [USDateLong](<https://github.com/agentstation/utc/blob/master/utc.go#L454>)

```go
func (t Time) USDateLong() string
```

USDateLong formats time as "January 2, 2006"

<a name="Time.USDateShort"></a>
### func \(Time\) [USDateShort](<https://github.com/agentstation/utc/blob/master/utc.go#L449>)

```go
func (t Time) USDateShort() string
```

USDateShort formats time as "01/02/2006"

<a name="Time.USDateTime12"></a>
### func \(Time\) [USDateTime12](<https://github.com/agentstation/utc/blob/master/utc.go#L459>)

```go
func (t Time) USDateTime12() string
```

USDateTime12 formats time as "01/02/2006 03:04:05 PM"

<a name="Time.USDateTime24"></a>
### func \(Time\) [USDateTime24](<https://github.com/agentstation/utc/blob/master/utc.go#L464>)

```go
func (t Time) USDateTime24() string
```

USDateTime24 formats time as "01/02/2006 15:04:05"

<a name="Time.USTime12"></a>
### func \(Time\) [USTime12](<https://github.com/agentstation/utc/blob/master/utc.go#L469>)

```go
func (t Time) USTime12() string
```

USTime12 formats time as "3:04 PM"

<a name="Time.USTime24"></a>
### func \(Time\) [USTime24](<https://github.com/agentstation/utc/blob/master/utc.go#L474>)

```go
func (t Time) USTime24() string
```

USTime24 formats time as "15:04"

<a name="Time.UTC"></a>
### func \(Time\) [UTC](<https://github.com/agentstation/utc/blob/master/utc.go#L336>)

```go
func (t Time) UTC() time.Time
```

UTC returns t in UTC

<a name="Time.Unix"></a>
### func \(Time\) [Unix](<https://github.com/agentstation/utc/blob/master/utc.go#L568>)

```go
func (t Time) Unix() int64
```



<a name="Time.UnixMilli"></a>
### func \(Time\) [UnixMilli](<https://github.com/agentstation/utc/blob/master/utc.go#L569>)

```go
func (t Time) UnixMilli() int64
```



<a name="Time.UnmarshalJSON"></a>
### func \(\*Time\) [UnmarshalJSON](<https://github.com/agentstation/utc/blob/master/utc.go#L156>)

```go
func (t *Time) UnmarshalJSON(data []byte) error
```

UnmarshalJSON implements the json.Unmarshaler interface for utc.Time

<a name="Time.UnmarshalText"></a>
### func \(\*Time\) [UnmarshalText](<https://github.com/agentstation/utc/blob/master/utc.go#L206>)

```go
func (t *Time) UnmarshalText(text []byte) error
```

UnmarshalText implements encoding.TextUnmarshaler.

<a name="Time.UnmarshalYAML"></a>
### func \(\*Time\) [UnmarshalYAML](<https://github.com/agentstation/utc/blob/master/utc.go#L220>)

```go
func (t *Time) UnmarshalYAML(unmarshal func(any) error) error
```

UnmarshalYAML implements the yaml.Unmarshaler interface for utc.Time

<a name="Time.Value"></a>
### func \(\*Time\) [Value](<https://github.com/agentstation/utc/blob/master/utc.go#L270>)

```go
func (t *Time) Value() (driver.Value, error)
```

Value implements the driver.Valuer interface for database operations for utc.Time. It returns the time.Time value and assumes the time is already in UTC.

Returns an error if called on a nil receiver instead of panicking to allow graceful error handling in database operations. In debug builds, nil receivers are logged.

<a name="Time.WeekdayLong"></a>
### func \(Time\) [WeekdayLong](<https://github.com/agentstation/utc/blob/master/utc.go#L515>)

```go
func (t Time) WeekdayLong() string
```

WeekdayLong formats time as "Monday"

<a name="Time.WeekdayShort"></a>
### func \(Time\) [WeekdayShort](<https://github.com/agentstation/utc/blob/master/utc.go#L520>)

```go
func (t Time) WeekdayShort() string
```

WeekdayShort formats time as "Mon"

<a name="TimeLayout"></a>
## type [TimeLayout](<https://github.com/agentstation/utc/blob/master/utc.go#L31>)



```go
type TimeLayout string
```

<a name="TimeLayoutUSDateShort"></a>Add layout constants at package level

```go
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
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
