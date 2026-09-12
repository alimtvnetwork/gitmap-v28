package checkpoint

import "os"

// osWriteFile is a thin alias kept in its own file so tests can poke
// the on-disk state via the same primitive checkpoint.go uses, without
// forcing checkpoint_test.go to import os directly.
func osWriteFile(p string, b []byte) error {
	return os.WriteFile(p, b, 0o644)
}
