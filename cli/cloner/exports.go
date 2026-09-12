package cloner

import "github.com/alimtvnetwork/gitmap-v28/cli/model"

// CloneAll exposes the internal cloneAll runner for integration and heavy tests.
func CloneAll(records []model.ScanRecord, targetDir string, opts CloneOptions) model.CloneSummary {
	return cloneAll(records, targetDir, opts)
}
