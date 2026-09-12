package fileutil

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

// ExportText writes a string to a file, overwriting if it exists.
func ExportText(path string, content string, perm FilePermType) result.Wrap[bool] {
	fRes := CreateFile(path, perm)
	if fRes.HasError() {
		return result.WrapFailure[bool](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	_, err := f.WriteString(content)
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, err, "failed to write text to file: "+path))
	}

	return result.WrapSuccess(true)
}

// ExportLines writes an array of strings to a file, separated by newlines.
func ExportLines(path string, lines []string, perm FilePermType) result.Wrap[bool] {
	content := strings.Join(lines, "\n") + "\n"
	if len(lines) == 0 {
		content = ""
	}

	return ExportText(path, content, perm)
}

// ExportJSON writes a data structure to a file as formatted JSON.
func ExportJSON(path string, data any, perm FilePermType) result.Wrap[bool] {
	fRes := CreateFile(path, perm)
	if fRes.HasError() {
		return result.WrapFailure[bool](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(data)
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.Serialization, err, "failed to encode JSON to: "+path))
	}

	return result.WrapSuccess(true)
}

// ExportYAML writes a data structure to a file as YAML.
func ExportYAML(path string, data any, perm FilePermType) result.Wrap[bool] {
	fRes := CreateFile(path, perm)
	if fRes.HasError() {
		return result.WrapFailure[bool](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()

	encoder.SetIndent(2)
	err := encoder.Encode(data)
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.Serialization, err, "failed to encode YAML to: "+path))
	}

	return result.WrapSuccess(true)
}
