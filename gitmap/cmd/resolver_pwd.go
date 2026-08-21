package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func resolveByPWD(all []model.ScanRecord) (*model.ScanRecord, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get current directory: %w", err)
	}
	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, pwd) {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("current directory (%s) is not a tracked repository", pwd)
}
