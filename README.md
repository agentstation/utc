# utc

```sh
                                _            _    _                    
                         _   _ | |_  ___    | |_ (_)_ __ ___   ___    
                        | | | || __|/ __|   | __|| | '_ ` _ \ / _ \   
                        | |_| || |_| (__    | |_ | | | | | | |  __/   
                         \___/  \__|\___|    \__||_|_| |_| |_|\___|   
```

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/agentstation/utc)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/agentstation/utc/ci.yaml?style=flat-square)](https://github.com/agentstation/utc/actions)
[![codecov](https://codecov.io/gh/agentstation/utc/graph/badge.svg?token=EOAZUVVH7H)](https://codecov.io/gh/agentstation/utc)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/agentstation/utc/main/LICENSE)

The `utc` package provides a **zero-dependency** wrapper around Go's `time.Time` for applications that want timestamp fields normalized to UTC at package boundaries.

## What It Solves

Go's `time.Time` is flexible, but UTC storage is usually a convention that every caller has to remember. `utc.Time` makes that convention a package boundary: values are normalized when they enter, stored privately, and emitted as UTC through JSON, text, YAML, SQL, and standard `time.Time` accessors.

Use it for API models, database records, event payloads, and config structs where timestamps should always represent instants in UTC.

## Key Features

### UTC Storage
- The underlying `time.Time` is unexported, so callers cannot construct non-UTC values with struct literals.
- Constructors, parsers, scanners, and unmarshaling methods normalize values to UTC.
- Use `t.UTC()` or `t.Time()` when you need a standard `time.Time`.

### Serialization
- JSON accepts strings or `null`; non-string JSON values are rejected.
- JSON, text, and YAML output preserve sub-second precision with RFC3339Nano formatting.
- SQL `Value` and `Scan` support UTC-normalized database boundaries.

### Convenience
- Common US/EU date and time formatting helpers.
- Named US timezone helpers with DST-aware `Pacific`, `Eastern`, `Central`, and `Mountain` conversions.
- Unix timestamp and UTC day-boundary helpers.

### Compatibility
- Zero external dependencies.
- Go 1.18 or later.
- A convenient replacement for timestamp fields and serialization/database boundaries; use `t.Time()` or `t.UTC()` when an API requires a concrete `time.Time`.

## Installation

To install the `utc` package, use the following command:

```sh
go get github.com/agentstation/utc
```

**Requirements**: Go 1.18 or later

**YAML Usage**: The package implements YAML marshal/unmarshal interface methods. To encode or decode YAML in an application, install the YAML library you already use:
```bash
go get github.com/goccy/go-yaml
```

## Why Choose UTC?

### ✅ **Zero Dependencies**
- **No external dependencies** for core functionality
- **Lightweight** - adds minimal footprint to your project
- **Fast installation** - `go get` with no dependency resolution delays

### ✅ **Maximum Compatibility**
- **Go 1.18+** support (broader than most time libraries)
- **Cross-platform** - works on all Go-supported platforms

### ✅ **Boundary Safety**
- **Enforced UTC storage** - wrapped values are normalized before storage and output
- **Strict JSON input** - accidental numbers, booleans, objects, and arrays are rejected
- **Precision preserving** - nanoseconds survive JSON/text/YAML round trips
- **Race condition tested** - safe for concurrent use
- **Comprehensive test coverage** - focused coverage for parsing, formatting, serialization, database, and timezone helpers

### ✅ **Developer Experience**
- **Small API** - explicit constructors and explicit conversion back to `time.Time`
- **Rich formatting options** - US/EU date formats, RFC standards, custom layouts
- **Automatic timezone handling** - PST/PDT, EST/EDT transitions handled correctly
- **Serialization ready** - JSON, YAML, Database, Text encoding

### ✅ **Optional Advanced Features**
- **Debug mode** - development-time logging for nil receiver calls on pointer-based methods
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

## UTC vs Standard Library

See the difference between `utc` and Go's standard `time` package:

| Feature | Standard `time.Time` | `utc.Time` |
|---------|---------------------|------------|
| **UTC Storage** | Manual convention | Enforced by unexported storage and constructors |
| **Dependencies** | Zero deps | Zero deps |
| **Rich Formatting** | Manual layout strings | Built-in US/EU/RFC helpers |
| **Timezone Conversion** | Manual location loading | Helpers for common US zones plus generic `In` |
| **JSON Support** | RFC3339 string support | Strict string/null input, UTC normalization, RFC3339Nano output |
| **YAML Support** | No built-in YAML methods | `MarshalYAML`/`UnmarshalYAML` interface methods |
| **Text Encoding** | `time.Time` supports text encoding | UTC-normalized text encoding |
| **Database Ready** | `time.Time` is database-friendly | UTC-normalized `Scan`/`Value` methods |
| **Unix Timestamps** | `Unix`/`UnixMilli` methods | `Unix`/`UnixMilli` plus constructors |
| **Day Boundaries** | Manual calculation | UTC `StartOfDay`/`EndOfDay` methods |
| **Debug Support** | N/A | Optional debug logging for nil pointer calls |

## Migration

Start at boundaries: change API, database, event, and config timestamp fields from `time.Time` to `utc.Time`. Keep existing code that requires `time.Time` by calling `t.Time()` or `t.UTC()`.

```go
createdAt := utc.New(row.CreatedAt) // from time.Time
legacy(createdAt.Time())            // back to time.Time when required
```

For helper functions that should accept both `time.Time` and `utc.Time`, use the `utc.UTC` interface with `utc.From`. `utc.UTC` is an interface for values that expose `UTC() time.Time`, not a timezone value:

```go
func normalize(t utc.UTC) utc.Time {
    return utc.From(t)
}
```

## Library Integration

- JSON, text encoding, YAML libraries, and `database/sql` can use `utc.Time` fields directly through standard marshal/unmarshal, scanner, and valuer interfaces.
- Code generators and ORMs such as sqlc, GORM-style models, PostgreSQL/MySQL layers, MongoDB, or DynamoDB integrations vary: use `utc.Time` when they support custom field types or scanner/valuer interfaces; use `t.Time()`/`utc.New(...)` at the boundary when they require concrete `time.Time`.
- The root module has no external dependencies. Compile-time interface assertions that only prove compatibility live in tests; `database/sql/driver` remains a production import because `Value() (driver.Value, error)` is the standard SQL value interface.

## Notes for Code Generators and AI Assistants

Use `github.com/agentstation/utc` when a Go struct field should store and emit UTC-normalized instants. The main type is `utc.Time`; it wraps an unexported `time.Time` and is not a type alias. Convert into the package with `utc.New(time.Time)`, `utc.From(utc.UTC)`, or the parse/unmarshal/scan methods. Convert out with `t.Time()` or `t.UTC()` when another API requires a concrete `time.Time`.

Do not construct `utc.Time` with struct literals outside this package, do not assume it has every `time.Time` method, and do not add external dependencies to the root module for optional codec support.

**Before** (standard library):
```go
createdAt := time.Now().UTC()
fmt.Println(createdAt.Format("01/02/2006"))

loc, err := time.LoadLocation("America/New_York")
if err != nil {
    return err
}
fmt.Println(createdAt.In(loc).Format(time.Kitchen))
```

**After** (with UTC):
```go
t := utc.Now()
fmt.Println(t.Eastern().Format(time.Kitchen))  // Auto EST/EDT
fmt.Println(t.USDateShort())                   // "01/15/2024"
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

The package includes YAML marshaling/unmarshaling support through `MarshalYAML`/`UnmarshalYAML` methods that follow the common Go YAML interfaces. The package itself does not import a YAML library.

```sh
# Run root tests including YAML interface coverage
go test -tags yaml ./...

# Run root YAML tests plus real YAML codec integration tests
make test-yaml
```

**Note**: Actual YAML file parsing/emission requires a YAML package in your application. `utc` only provides the methods those packages call. Choose a YAML library version compatible with your application's Go version. The real codec integration test lives in `integration/yaml` and currently uses `github.com/goccy/go-yaml` with Go 1.21+, so the root module remains dependency-free and Go 1.18-compatible.

## Testing

The project includes comprehensive Makefile targets for testing:

```sh
# Run tests (Go 1.18+, no dependencies)
make test

# Run tests with YAML interface and real codec coverage
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

When debug mode is enabled, the package logs warnings when pointer-based methods that can return errors are called on nil receivers:

```
[UTC DEBUG] 2024/01/02 15:04:05 debug.go:26: MarshalJSON() called on nil *Time receiver
[UTC DEBUG] 2024/01/02 15:04:05 debug.go:26: Scan() called on nil *Time receiver
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

Package utc provides a small time.Time wrapper that stores instants in UTC.

The package keeps the underlying time.Time value unexported, normalizes values through constructors, parsers, scanners, and serializers, and exposes standard time.Time values through the Time and UTC methods and the UTC interface. When compiled with the debug build tag \(\-tags debug\), pointer\-based methods log nil receiver calls before returning errors where Go permits that behavior.

Key features:

- Constructors and parsers normalize values to UTC
- JSON marshaling/unmarshaling uses strict string/null inputs and preserves sub\-second precision
- Text and YAML marshal/unmarshal support
- SQL database compatibility
- Timezone conversion helpers with automatic DST handling
- Extensive formatting options for US and EU date formats

Debug mode:

```
To enable debug logging, compile with: go build -tags debug
This logs nil receiver calls for pointer-based methods that can return errors.
```

## Index

- [func ValidateTimezoneAvailability\(\) error](<#ValidateTimezoneAvailability>)
- [type Time](<#Time>)
  - [func From\(t UTC\) Time](<#From>)
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
  - [func \(t Time\) String\(\) string](<#Time.String>)
  - [func \(t Time\) Sub\(u Time\) time.Duration](<#Time.Sub>)
  - [func \(t Time\) Time\(\) time.Time](<#Time.Time>)
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
  - [func \(t Time\) Value\(\) \(driver.Value, error\)](<#Time.Value>)
  - [func \(t Time\) WeekdayLong\(\) string](<#Time.WeekdayLong>)
  - [func \(t Time\) WeekdayShort\(\) string](<#Time.WeekdayShort>)
- [type TimeLayout](<#TimeLayout>)
- [type UTC](<#UTC>)


<a name="ValidateTimezoneAvailability"></a>
## func [ValidateTimezoneAvailability](<https://github.com/agentstation/utc/blob/main/utc.go#L97>)

```go
func ValidateTimezoneAvailability() error
```

ValidateTimezoneAvailability checks whether package timezone locations were initialized. It returns nil if initialization succeeded.

<a name="Time"></a>
## type [Time](<https://github.com/agentstation/utc/blob/main/utc.go#L105-L107>)

Time stores a time instant normalized to UTC.

```go
type Time struct {
    // contains filtered or unexported fields
}
```

<a name="From"></a>
### func [From](<https://github.com/agentstation/utc/blob/main/utc.go#L126>)

```go
func From(t UTC) Time
```

From returns a new Time from any value that exposes a UTC time.Time.

<a name="FromUnix"></a>
### func [FromUnix](<https://github.com/agentstation/utc/blob/main/utc.go#L583>)

```go
func FromUnix(sec int64) Time
```

Unix helpers

<a name="FromUnixMilli"></a>
### func [FromUnixMilli](<https://github.com/agentstation/utc/blob/main/utc.go#L584>)

```go
func FromUnixMilli(ms int64) Time
```



<a name="New"></a>
### func [New](<https://github.com/agentstation/utc/blob/main/utc.go#L121>)

```go
func New(t time.Time) Time
```

New returns a new Time from a time.Time.

<a name="Now"></a>
### func [Now](<https://github.com/agentstation/utc/blob/main/utc.go#L116>)

```go
func Now() Time
```

Now returns the current time in UTC.

<a name="Parse"></a>
### func [Parse](<https://github.com/agentstation/utc/blob/main/utc.go#L158>)

```go
func Parse(layout string, s string) (Time, error)
```

Parse parses a time string in the specified format and returns a Time.

<a name="ParseRFC3339"></a>
### func [ParseRFC3339](<https://github.com/agentstation/utc/blob/main/utc.go#L140>)

```go
func ParseRFC3339(s string) (Time, error)
```

ParseRFC3339 parses a time string in RFC3339 format and returns a Time.

<a name="ParseRFC3339Nano"></a>
### func [ParseRFC3339Nano](<https://github.com/agentstation/utc/blob/main/utc.go#L149>)

```go
func ParseRFC3339Nano(s string) (Time, error)
```

ParseRFC3339Nano parses a time string in RFC3339Nano format and returns a Time.

<a name="Time.ANSIC"></a>
### func \(Time\) [ANSIC](<https://github.com/agentstation/utc/blob/main/utc.go#L455>)

```go
func (t Time) ANSIC() string
```

ANSIC formats time as "Mon Jan \_2 15:04:05 2006"

<a name="Time.Add"></a>
### func \(Time\) [Add](<https://github.com/agentstation/utc/blob/main/utc.go#L340>)

```go
func (t Time) Add(d time.Duration) Time
```

Add returns the time t\+d

<a name="Time.After"></a>
### func \(Time\) [After](<https://github.com/agentstation/utc/blob/main/utc.go#L330>)

```go
func (t Time) After(u Time) bool
```

After reports whether the time is after u

<a name="Time.Before"></a>
### func \(Time\) [Before](<https://github.com/agentstation/utc/blob/main/utc.go#L325>)

```go
func (t Time) Before(u Time) bool
```

Before reports whether the time is before u

<a name="Time.CST"></a>
### func \(Time\) [CST](<https://github.com/agentstation/utc/blob/main/utc.go#L365>)

```go
func (t Time) CST() time.Time
```

CST returns t in CST

<a name="Time.Central"></a>
### func \(Time\) [Central](<https://github.com/agentstation/utc/blob/main/utc.go#L391>)

```go
func (t Time) Central() time.Time
```

Central returns t in Central time \(handles CST/CDT automatically\)

<a name="Time.DateOnly"></a>
### func \(Time\) [DateOnly](<https://github.com/agentstation/utc/blob/main/utc.go#L549>)

```go
func (t Time) DateOnly() string
```

DateOnly formats time as "2006\-01\-02"

<a name="Time.EST"></a>
### func \(Time\) [EST](<https://github.com/agentstation/utc/blob/main/utc.go#L360>)

```go
func (t Time) EST() time.Time
```

EST returns t in EST

<a name="Time.EUDateLong"></a>
### func \(Time\) [EUDateLong](<https://github.com/agentstation/utc/blob/main/utc.go#L501>)

```go
func (t Time) EUDateLong() string
```

EUDateLong formats time as "2 January 2006"

<a name="Time.EUDateShort"></a>
### func \(Time\) [EUDateShort](<https://github.com/agentstation/utc/blob/main/utc.go#L496>)

```go
func (t Time) EUDateShort() string
```

EUDateShort formats time as "02/01/2006"

<a name="Time.EUDateTime12"></a>
### func \(Time\) [EUDateTime12](<https://github.com/agentstation/utc/blob/main/utc.go#L506>)

```go
func (t Time) EUDateTime12() string
```

EUDateTime12 formats time as "02/01/2006 03:04:05 PM"

<a name="Time.EUDateTime24"></a>
### func \(Time\) [EUDateTime24](<https://github.com/agentstation/utc/blob/main/utc.go#L511>)

```go
func (t Time) EUDateTime24() string
```

EUDateTime24 formats time as "02/01/2006 15:04:05"

<a name="Time.EUTime12"></a>
### func \(Time\) [EUTime12](<https://github.com/agentstation/utc/blob/main/utc.go#L516>)

```go
func (t Time) EUTime12() string
```

EUTime12 formats time as "3:04 PM"

<a name="Time.EUTime24"></a>
### func \(Time\) [EUTime24](<https://github.com/agentstation/utc/blob/main/utc.go#L521>)

```go
func (t Time) EUTime24() string
```

EUTime24 formats time as "15:04"

<a name="Time.Eastern"></a>
### func \(Time\) [Eastern](<https://github.com/agentstation/utc/blob/main/utc.go#L383>)

```go
func (t Time) Eastern() time.Time
```

Eastern returns t in Eastern time \(handles EST/EDT automatically\)

<a name="Time.EndOfDay"></a>
### func \(Time\) [EndOfDay](<https://github.com/agentstation/utc/blob/main/utc.go#L594>)

```go
func (t Time) EndOfDay() Time
```



<a name="Time.Equal"></a>
### func \(Time\) [Equal](<https://github.com/agentstation/utc/blob/main/utc.go#L335>)

```go
func (t Time) Equal(u Time) bool
```

Equal reports whether t and u represent the same time instant

<a name="Time.Format"></a>
### func \(Time\) [Format](<https://github.com/agentstation/utc/blob/main/utc.go#L412>)

```go
func (t Time) Format(layout string) string
```

Format formats the time using the specified layout

<a name="Time.ISO8601"></a>
### func \(Time\) [ISO8601](<https://github.com/agentstation/utc/blob/main/utc.go#L435>)

```go
func (t Time) ISO8601() string
```

ISO8601 formats time as "2006\-01\-02T15:04:05Z07:00" \(same as RFC3339\)

<a name="Time.In"></a>
### func \(Time\) [In](<https://github.com/agentstation/utc/blob/main/utc.go#L566>)

```go
func (t Time) In(name string) (time.Time, error)
```

In converts time to a named location \(e.g., "America/Los\_Angeles"\).

<a name="Time.InLocation"></a>
### func \(Time\) [InLocation](<https://github.com/agentstation/utc/blob/main/utc.go#L575>)

```go
func (t Time) InLocation(loc *time.Location) time.Time
```

InLocation converts time to a provided \*time.Location.

<a name="Time.IsZero"></a>
### func \(Time\) [IsZero](<https://github.com/agentstation/utc/blob/main/utc.go#L407>)

```go
func (t Time) IsZero() bool
```

Add the useful utility methods while maintaining chainability

<a name="Time.Kitchen"></a>
### func \(Time\) [Kitchen](<https://github.com/agentstation/utc/blob/main/utc.go#L559>)

```go
func (t Time) Kitchen() string
```

Kitchen formats time as "3:04PM"

<a name="Time.MST"></a>
### func \(Time\) [MST](<https://github.com/agentstation/utc/blob/main/utc.go#L370>)

```go
func (t Time) MST() time.Time
```

MST returns t in MST

<a name="Time.MarshalJSON"></a>
### func \(\*Time\) [MarshalJSON](<https://github.com/agentstation/utc/blob/main/utc.go#L208>)

```go
func (t *Time) MarshalJSON() ([]byte, error)
```

MarshalJSON implements the json.Marshaler interface for Time. Returns an error for nil receivers to maintain consistency with standard marshaling behavior.

<a name="Time.MarshalText"></a>
### func \(Time\) [MarshalText](<https://github.com/agentstation/utc/blob/main/utc.go#L217>)

```go
func (t Time) MarshalText() ([]byte, error)
```

MarshalText implements encoding.TextMarshaler.

<a name="Time.MarshalYAML"></a>
### func \(Time\) [MarshalYAML](<https://github.com/agentstation/utc/blob/main/utc.go#L268>)

```go
func (t Time) MarshalYAML() (any, error)
```

MarshalYAML implements the yaml.Marshaler interface for Time.

<a name="Time.MonthLong"></a>
### func \(Time\) [MonthLong](<https://github.com/agentstation/utc/blob/main/utc.go#L539>)

```go
func (t Time) MonthLong() string
```

MonthLong formats time as "January"

<a name="Time.MonthShort"></a>
### func \(Time\) [MonthShort](<https://github.com/agentstation/utc/blob/main/utc.go#L544>)

```go
func (t Time) MonthShort() string
```

MonthShort formats time as "Jan"

<a name="Time.Mountain"></a>
### func \(Time\) [Mountain](<https://github.com/agentstation/utc/blob/main/utc.go#L399>)

```go
func (t Time) Mountain() time.Time
```

Mountain returns t in Mountain time \(handles MST/MDT automatically\)

<a name="Time.PST"></a>
### func \(Time\) [PST](<https://github.com/agentstation/utc/blob/main/utc.go#L355>)

```go
func (t Time) PST() time.Time
```

PST returns t in PST

<a name="Time.Pacific"></a>
### func \(Time\) [Pacific](<https://github.com/agentstation/utc/blob/main/utc.go#L375>)

```go
func (t Time) Pacific() time.Time
```

Pacific returns t in Pacific time \(handles PST/PDT automatically\)

<a name="Time.RFC3339"></a>
### func \(Time\) [RFC3339](<https://github.com/agentstation/utc/blob/main/utc.go#L425>)

```go
func (t Time) RFC3339() string
```

RFC3339 formats time as "2006\-01\-02T15:04:05Z07:00"

<a name="Time.RFC3339Nano"></a>
### func \(Time\) [RFC3339Nano](<https://github.com/agentstation/utc/blob/main/utc.go#L430>)

```go
func (t Time) RFC3339Nano() string
```

RFC3339Nano formats time as "2006\-01\-02T15:04:05.999999999Z07:00"

<a name="Time.RFC822"></a>
### func \(Time\) [RFC822](<https://github.com/agentstation/utc/blob/main/utc.go#L440>)

```go
func (t Time) RFC822() string
```

RFC822 formats time as "02 Jan 06 15:04 MST"

<a name="Time.RFC822Z"></a>
### func \(Time\) [RFC822Z](<https://github.com/agentstation/utc/blob/main/utc.go#L445>)

```go
func (t Time) RFC822Z() string
```

RFC822Z formats time as "02 Jan 06 15:04 \-0700"

<a name="Time.RFC850"></a>
### func \(Time\) [RFC850](<https://github.com/agentstation/utc/blob/main/utc.go#L450>)

```go
func (t Time) RFC850() string
```

RFC850 formats time as "Monday, 02\-Jan\-06 15:04:05 MST"

<a name="Time.Scan"></a>
### func \(\*Time\) [Scan](<https://github.com/agentstation/utc/blob/main/utc.go#L291>)

```go
func (t *Time) Scan(value any) error
```

Scan implements sql.Scanner for database operations. It accepts time.Time, string, and \[\]byte values and stores them in UTC.

<a name="Time.StartOfDay"></a>
### func \(Time\) [StartOfDay](<https://github.com/agentstation/utc/blob/main/utc.go#L589>)

```go
func (t Time) StartOfDay() Time
```

Day helpers \- times are always in UTC within this package

<a name="Time.String"></a>
### func \(Time\) [String](<https://github.com/agentstation/utc/blob/main/utc.go#L278>)

```go
func (t Time) String() string
```

String implements fmt.Stringer. It prints the time in RFC3339Nano format.

<a name="Time.Sub"></a>
### func \(Time\) [Sub](<https://github.com/agentstation/utc/blob/main/utc.go#L345>)

```go
func (t Time) Sub(u Time) time.Duration
```

Sub returns the duration t\-u

<a name="Time.Time"></a>
### func \(Time\) [Time](<https://github.com/agentstation/utc/blob/main/utc.go#L135>)

```go
func (t Time) Time() time.Time
```

Time returns the underlying time.Time normalized to UTC.

<a name="Time.TimeFormat"></a>
### func \(Time\) [TimeFormat](<https://github.com/agentstation/utc/blob/main/utc.go#L417>)

```go
func (t Time) TimeFormat(layout TimeLayout) string
```

TimeFormat formats the time using the specified layout

<a name="Time.TimeOnly"></a>
### func \(Time\) [TimeOnly](<https://github.com/agentstation/utc/blob/main/utc.go#L554>)

```go
func (t Time) TimeOnly() string
```

TimeOnly formats time as "15:04:05"

<a name="Time.USDateLong"></a>
### func \(Time\) [USDateLong](<https://github.com/agentstation/utc/blob/main/utc.go#L468>)

```go
func (t Time) USDateLong() string
```

USDateLong formats time as "January 2, 2006"

<a name="Time.USDateShort"></a>
### func \(Time\) [USDateShort](<https://github.com/agentstation/utc/blob/main/utc.go#L463>)

```go
func (t Time) USDateShort() string
```

USDateShort formats time as "01/02/2006"

<a name="Time.USDateTime12"></a>
### func \(Time\) [USDateTime12](<https://github.com/agentstation/utc/blob/main/utc.go#L473>)

```go
func (t Time) USDateTime12() string
```

USDateTime12 formats time as "01/02/2006 03:04:05 PM"

<a name="Time.USDateTime24"></a>
### func \(Time\) [USDateTime24](<https://github.com/agentstation/utc/blob/main/utc.go#L478>)

```go
func (t Time) USDateTime24() string
```

USDateTime24 formats time as "01/02/2006 15:04:05"

<a name="Time.USTime12"></a>
### func \(Time\) [USTime12](<https://github.com/agentstation/utc/blob/main/utc.go#L483>)

```go
func (t Time) USTime12() string
```

USTime12 formats time as "3:04 PM"

<a name="Time.USTime24"></a>
### func \(Time\) [USTime24](<https://github.com/agentstation/utc/blob/main/utc.go#L488>)

```go
func (t Time) USTime24() string
```

USTime24 formats time as "15:04"

<a name="Time.UTC"></a>
### func \(Time\) [UTC](<https://github.com/agentstation/utc/blob/main/utc.go#L350>)

```go
func (t Time) UTC() time.Time
```

UTC returns t in UTC

<a name="Time.Unix"></a>
### func \(Time\) [Unix](<https://github.com/agentstation/utc/blob/main/utc.go#L585>)

```go
func (t Time) Unix() int64
```



<a name="Time.UnixMilli"></a>
### func \(Time\) [UnixMilli](<https://github.com/agentstation/utc/blob/main/utc.go#L586>)

```go
func (t Time) UnixMilli() int64
```



<a name="Time.UnmarshalJSON"></a>
### func \(\*Time\) [UnmarshalJSON](<https://github.com/agentstation/utc/blob/main/utc.go#L167>)

```go
func (t *Time) UnmarshalJSON(data []byte) error
```

UnmarshalJSON implements the json.Unmarshaler interface for Time.

<a name="Time.UnmarshalText"></a>
### func \(\*Time\) [UnmarshalText](<https://github.com/agentstation/utc/blob/main/utc.go#L222>)

```go
func (t *Time) UnmarshalText(text []byte) error
```

UnmarshalText implements encoding.TextUnmarshaler.

<a name="Time.UnmarshalYAML"></a>
### func \(\*Time\) [UnmarshalYAML](<https://github.com/agentstation/utc/blob/main/utc.go#L240>)

```go
func (t *Time) UnmarshalYAML(unmarshal func(any) error) error
```

UnmarshalYAML implements the yaml.Unmarshaler interface for Time.

<a name="Time.Value"></a>
### func \(Time\) [Value](<https://github.com/agentstation/utc/blob/main/utc.go#L284>)

```go
func (t Time) Value() (driver.Value, error)
```

Value implements driver.Valuer for database operations. It returns the UTC time.Time value as a driver.Value.

<a name="Time.WeekdayLong"></a>
### func \(Time\) [WeekdayLong](<https://github.com/agentstation/utc/blob/main/utc.go#L529>)

```go
func (t Time) WeekdayLong() string
```

WeekdayLong formats time as "Monday"

<a name="Time.WeekdayShort"></a>
### func \(Time\) [WeekdayShort](<https://github.com/agentstation/utc/blob/main/utc.go#L534>)

```go
func (t Time) WeekdayShort() string
```

WeekdayShort formats time as "Mon"

<a name="TimeLayout"></a>
## type [TimeLayout](<https://github.com/agentstation/utc/blob/main/utc.go#L35>)

TimeLayout names one of the package's built\-in formatting layouts.

```go
type TimeLayout string
```

<a name="TimeLayoutUSDateShort"></a>Built\-in formatting layouts.

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

<a name="UTC"></a>
## type [UTC](<https://github.com/agentstation/utc/blob/main/utc.go#L111-L113>)

UTC is implemented by values that can expose themselves as a UTC time.Time. It is an interface, not a timezone value. Both time.Time and utc.Time satisfy it.

```go
type UTC interface {
    UTC() time.Time
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
