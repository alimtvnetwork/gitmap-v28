package payloadconv

import (
	"encoding/json"
	"fmt"
	"strings"

	"coding-guidelines/common/pkg/result"
)

// ToBytes smartly converts any generic payload into a byte array for writing.
// It avoids variable mutation and uses early returns.
func ToBytes(payload any) result.Wrap[[]byte] {
	if payload == nil {
		return result.Success([]byte{})
	}

	switch v := payload.(type) {
	case []byte:
		return result.Success(v)
	case string:
		return result.Success([]byte(v))
	case []string:
		// Convert array of strings into line-by-line format
		return result.Success([]byte(strings.Join(v, "\n") + "\n"))
	case error:
		return result.Success([]byte(v.Error()))
	default:
		// Attempt JSON encoding for structs, maps, etc.
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			// Fallback to Sprint if it cannot be marshaled
			return result.Success([]byte(fmt.Sprint(v)))
		}
		
		// Append trailing newline for better file readability
		b = append(b, '\n')
		return result.Success(b)
	}
}

// ToBytesMust is a convenience wrapper that panics on failure, useful if you know the type is valid.
func ToBytesMust(payload any) []byte {
	res := ToBytes(payload)
	if res.IsFailure() {
		panic(res.Fault().Error())
	}
	return res.Data()
}
