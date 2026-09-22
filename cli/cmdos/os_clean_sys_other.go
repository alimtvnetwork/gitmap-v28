//go:build !linux

package cmdos

type fallbackSystemCleanEngine struct{}

func newPlatformSystemCleanEngine() SystemCleanOperator {
	return &fallbackSystemCleanEngine{}
}

func (f *fallbackSystemCleanEngine) CleanSystemPackages() (int64, error) {
	return 0, nil
}

func (f *fallbackSystemCleanEngine) VacuumJournals() error {
	return nil
}
