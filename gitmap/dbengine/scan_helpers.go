package dbengine

import (
	"fmt"
	"math"
)

// ScanString safely converts an arbitrary scanned interface value into a string.
func ScanString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprint(val)
	}
}

// ScanInt safely converts an arbitrary scanned interface value into an int.
func ScanInt(v any) int {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case int:
		return val
	case int64:
		return clampInt64ToInt(val)
	case int32:
		return int(val)
	case uint64:
		return clampUint64ToInt(val)
	default:
		return 0
	}
}

func clampInt64ToInt(val int64) int {
	if val > math.MaxInt || val < math.MinInt {
		return 0
	}

	return int(val)
}

func clampUint64ToInt(val uint64) int {
	if val > math.MaxInt {
		return 0
	}

	return int(val)
}

// ScanInt64 safely converts an arbitrary scanned interface value into an int64.
func ScanInt64(v any) int64 {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case uint64:
		return clampUint64ToInt64(val)
	default:
		return 0
	}
}

func clampUint64ToInt64(val uint64) int64 {
	if val > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(val)
}

// ScanUint64 safely converts an arbitrary scanned interface value into a uint64.
func ScanUint64(v any) uint64 {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case uint64:
		return val
	case int64:
		return clampInt64ToUint64(val)
	case int:
		return clampIntToUint64(val)
	case uint:
		return uint64(val)
	default:
		return 0
	}
}

func clampInt64ToUint64(val int64) uint64 {
	if val < 0 {
		return 0
	}

	return uint64(val)
}

func clampIntToUint64(val int) uint64 {
	if val < 0 {
		return 0
	}

	return uint64(val)
}

// ScanUint safely converts an arbitrary scanned interface value into a uint.
func ScanUint(v any) uint {
	return uint(ScanUint64(v))
}

// ScanBool safely converts an arbitrary scanned interface value into a bool.
func ScanBool(v any) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case int64:
		return val != 0
	case int:
		return val != 0
	case uint64:
		return val != 0
	default:
		return false
	}
}

// ScanFloat64 safely converts an arbitrary scanned interface value into a float64.
func ScanFloat64(v any) float64 {
	if v == nil {
		return 0.0
	}

	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int64:
		return float64(val)
	case int:
		return float64(val)
	default:
		return 0.0
	}
}
