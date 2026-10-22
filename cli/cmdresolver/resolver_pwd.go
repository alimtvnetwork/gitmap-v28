package cmdresolver

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func resolveByPWD(all []model.ScanRecord) (*model.ScanRecord, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, apperror.WrapSimple(err, "get current directory")
	}

	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, pwd) {
			return &r, nil
		}
	}

	return nil, apperror.NewNotFoundError(fmt.Sprintf("current directory (%s) is not a tracked repository", pwd))
}
