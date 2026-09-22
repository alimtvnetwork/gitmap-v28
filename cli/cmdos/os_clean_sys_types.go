package cmdos

// SystemCleanEngine abstracts OS system-level package and log cleanup.
type SystemCleanEngine interface {
	CleanSystemPackages() (int64, error)
	VacuumJournals() error
}

// defaultSystemCleanEngine is the active engine instance.
var defaultSystemCleanEngine SystemCleanEngine = newPlatformSystemCleanEngine()

// GetSystemCleanEngine returns the active system clean engine.
func GetSystemCleanEngine() SystemCleanEngine {
	return defaultSystemCleanEngine
}

// SetSystemCleanEngine sets the active engine (used in tests).
func SetSystemCleanEngine(engine SystemCleanEngine) {
	defaultSystemCleanEngine = engine
}
