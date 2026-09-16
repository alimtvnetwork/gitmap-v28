package cmdupdate

import (
	"encoding/json"
	"io"
	"strings"
)

func decodeReleaseTagName(body io.Reader) string {
	var releaseData struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
	}

	if err := json.NewDecoder(body).Decode(&releaseData); err != nil {
		return ""
	}

	return extractCleanTag(releaseData.TagName, releaseData.Name)
}

func extractCleanTag(tagName, name string) string {
	tag := strings.TrimPrefix(tagName, "v")
	if len(tag) > 0 {
		return tag
	}

	return strings.TrimPrefix(name, "v")
}

func decodeVersionFromMap(body io.Reader) string {
	var rawMap map[string]interface{}
	if err := json.NewDecoder(body).Decode(&rawMap); err != nil {
		return "unknown"
	}

	return extractVersionValue(rawMap)
}

func extractVersionValue(rawMap map[string]interface{}) string {
	for _, key := range []string{"Version", "version"} {
		if v, isString := rawMap[key].(string); isString && len(v) > 0 {
			return v
		}
	}

	return "unknown"
}
