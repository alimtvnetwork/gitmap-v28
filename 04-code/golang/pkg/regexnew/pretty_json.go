package regexnew

import (
	"bytes"
	"encoding/json"
)

// prettyJson formats an object as indented JSON, swallowing errors safely.
func prettyJson(anyItem any) string {
	if anyItem == nil {
		return ""
	}

	allBytes, err := json.Marshal(anyItem)
	if err != nil || len(allBytes) == 0 {
		return ""
	}

	var prettyJSON bytes.Buffer
	_ = json.Indent(
		&prettyJSON,
		allBytes,
		EmptyString,
		Tab)

	return prettyJSON.String()
}
