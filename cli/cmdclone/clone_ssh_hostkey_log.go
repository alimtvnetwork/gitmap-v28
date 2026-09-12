package cmdclone

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func fmtCloneEnvError(err error) {
	fmt.Fprintf(os.Stderr, constants.ErrCloneSSHEnvFmt, err)
}
