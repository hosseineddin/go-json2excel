package converter

import (
	"fmt"
	"time"
)

// ConvertDynamic normalizes dynamic JSON types to standard Go types.
// Highly optimized: avoids reflection where possible.
func ConvertDynamic(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		// Fast path for checking RFC3339 dates without regex
		if len(v) >= 10 && (v[4] == '-' || v[4] == '/') {
			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
				return parsed
			}
		}
		return v
	case float64, bool:
		return v
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}
