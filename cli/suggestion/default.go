package suggestion

var defaultEngineInstance = NewEngine()

// DefaultEngine returns the shared default instance of the suggestion Engine.
func DefaultEngine() Engine {
	return defaultEngineInstance
}
