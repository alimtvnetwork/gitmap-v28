package fileutil

import (
	"encoding/json"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

// StreamJSON sequentially decodes a massive JSON array from a file, passing each element to the handler.
// This prevents excessive RAM usage when dealing with huge datasets.
func StreamJSON[T any](path string, handler func(T) *appfault.AppError) result.Wrap[bool] {
	fRes := OpenFile(path, FileOpenReadOnly, FilePermStandard)
	if fRes.HasError() {
		return result.WrapFailure[bool](fRes.Fault())
	}

	f := fRes.Data()
	defer f.Close()

	decoder := json.NewDecoder(f)

	// Read the opening bracket
	t, err := decoder.Token()
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.Serialization, err, "failed to read JSON array start in: "+path))
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return result.WrapFailureWithId[bool](errtype.Serialization, "StreamJSON requires the root element to be a JSON array")
	}

	for decoder.More() {
		var item T
		if err := decoder.Decode(&item); err != nil {
			return result.WrapFailure[bool](appfault.Wrap(errtype.Serialization, err, "failed to decode array element in: "+path))
		}

		if err := handler(item); err != nil {
			return result.WrapFailure[bool](err)
		}
	}

	// Read the closing bracket
	_, err = decoder.Token()
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.Serialization, err, "failed to read JSON array end in: "+path))
	}

	return result.WrapSuccess(true)
}
