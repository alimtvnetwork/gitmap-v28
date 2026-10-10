package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"strings"
)

func isFlagToken(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func loadStatusRecords(path string) ([]model.ScanRecord, error) {
	return model.LoadStatusRecords(path)
}
