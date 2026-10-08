package cmdas

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func WriteShellHandoff(targetPath string) {
	handoffFile := os.Getenv(constants.EnvGitmapHandoffFile)
	if len(handoffFile) == 0 {
		return
	}

	if len(targetPath) == 0 {
		return
	}

	if err := os.WriteFile(handoffFile, []byte(targetPath), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrShellHandoffWriteFmt, handoffFile, err)
	}
}

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}
