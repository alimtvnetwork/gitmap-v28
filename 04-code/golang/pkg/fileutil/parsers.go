package fileutil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

// ReadText reads the entire file content as a string.
func ReadText(path string) result.Wrap[string] {
	fRes := OpenFile(path, FileOpenReadOnly, FilePermStandard)
	if fRes.HasError() {
		return result.WrapFailure[string](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return result.WrapFailure[string](appfault.Wrap(errtype.IO, err, "failed to read file content: "+path))
	}

	return result.WrapSuccess(string(content))
}

// ReadLines reads the file and splits it into a string array by lines.
func ReadLines(path string) result.Wrap[[]string] {
	fRes := OpenFile(path, FileOpenReadOnly, FilePermStandard)
	if fRes.HasError() {
		return result.WrapFailure[[]string](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return result.WrapFailure[[]string](appfault.Wrap(errtype.IO, err, "error scanning file lines: "+path))
	}

	return result.WrapSuccess(lines)
}

// ReadJSON parses a JSON file into the specified type T.
func ReadJSON[T any](path string) result.Wrap[T] {
	var val T
	fRes := OpenFile(path, FileOpenReadOnly, FilePermStandard)
	if fRes.HasError() {
		return result.WrapFailure[T](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	err := json.NewDecoder(f).Decode(&val)
	if err != nil {
		return result.WrapFailure[T](appfault.Wrap(errtype.Serialization, err, "failed to decode JSON from: "+path))
	}

	return result.WrapSuccess(val)
}

// ReadYAML parses a YAML file into the specified type T.
func ReadYAML[T any](path string) result.Wrap[T] {
	var val T
	fRes := OpenFile(path, FileOpenReadOnly, FilePermStandard)
	if fRes.HasError() {
		return result.WrapFailure[T](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return result.WrapFailure[T](appfault.Wrap(errtype.IO, err, "failed to read file for YAML decoding: "+path))
	}

	// Remove BOM if present, which yaml.v3 doesn't handle natively sometimes
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf"))

	err = yaml.Unmarshal(content, &val)
	if err != nil {
		return result.WrapFailure[T](appfault.Wrap(errtype.Serialization, err, "failed to decode YAML from: "+path))
	}

	return result.WrapSuccess(val)
}
