package clonefrom

import (
	"os"
)

// writeFile is a tiny os.WriteFile wrapper with a fixed permission.
func writeFile(path, body string) error {
	return os.WriteFile(path, []byte(body), 0o644)
}
