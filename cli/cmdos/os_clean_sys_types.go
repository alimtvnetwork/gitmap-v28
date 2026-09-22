package cmdos

// SystemCleanOperator abstracts OS system-level package and log cleanup.
type SystemCleanOperator interface {
	CleanSystemPackages() (int64, error)
	VacuumJournals() error
}

// defaultSystemCleanEngine is the active engine instance.
var defaultSystemCleanEngine SystemCleanOperator = newPlatformSystemCleanEngine()

// GetSystemCleanEngine returns the active system clean engine.
func GetSystemCleanEngine() SystemCleanOperator {
	return defaultSystemCleanEngine
}

// SetSystemCleanEngine sets the active engine (used in tests).
func SetSystemCleanEngine(engine SystemCleanOperator) {
	defaultSystemCleanEngine = engine
}
