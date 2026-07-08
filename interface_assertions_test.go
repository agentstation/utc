package utc

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"time"
)

// Keep assertion-only imports in tests when they are not part of production
// method signatures. Value returns driver.Value, so database/sql/driver remains
// a production import; Scanner is structural and can be proven here.
var (
	_ UTC                      = time.Time{}
	_ UTC                      = Time{}
	_ json.Marshaler           = (*Time)(nil)
	_ json.Unmarshaler         = (*Time)(nil)
	_ encoding.TextMarshaler   = Time{}
	_ encoding.TextUnmarshaler = (*Time)(nil)
	_ driver.Valuer            = Time{}
	_ sql.Scanner              = (*Time)(nil)
)
