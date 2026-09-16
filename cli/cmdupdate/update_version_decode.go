package cmdupdate

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func decodeReleaseTagName(body io.Reader) string {
	var releaseData struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
	}

	err := json.NewDecoder(body).Decode(&releaseData)
	if err != nil {
		return ""
	}

	return extractCleanTag(releaseData.TagName, releaseData.Name)
}

func extractCleanTag(tagName, name string) string {
	tag := strings.TrimPrefix(tagName, constants.VersionPrefixV)
	hasTag := len(tag) > 0

	if hasTag {
		return tag
	}

	return strings.TrimPrefix(name, constants.VersionPrefixV)
}

func decodeVersionFromMap(body io.Reader) string {
	var rawMap map[string]interface{}
	err := json.NewDecoder(body).Decode(&rawMap)
	if err != nil {
		return constants.VersionUnknown
	}

	return extractVersionValue(rawMap)
}

func extractVersionValue(rawMap map[string]interface{}) string {
	for _, key := range constants.VersionKeys {
		v, isString := rawMap[key].(string)
		hasContent := isString && len(v) > 0

		if hasContent {
			return v
		}
	}

	return constants.VersionUnknown
}
