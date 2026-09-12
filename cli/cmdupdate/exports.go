package cmdupdate

// RunUpdate executes the update command logic.
func RunUpdate() error {
	return runUpdate()
}

// RunUpdateCleanup executes the update cleanup command logic.
func RunUpdateCleanup() error {
	return runUpdateCleanup()
}

// RunUpdateRunner executes the update runner command logic.
func RunUpdateRunner() error {
	return runUpdateRunner()
}

// ScheduleDeployedCleanupHandoff hands off cleanup to deployed binary.
func ScheduleDeployedCleanupHandoff() {
	scheduleDeployedCleanupHandoff()
}

// InitRunnerVerbose initializes verbose mode for runner.
func InitRunnerVerbose() {
	initRunnerVerbose()
}

// CreateHandoffCopy creates a handoff copy of the binary.
func CreateHandoffCopy(selfPath string) string {
	return createHandoffCopy(selfPath)
}

// HandleHandoffError handles errors during binary handoff.
func HandleHandoffError(err error) {
	handleHandoffError(err)
}

// WriteScriptToTemp writes script content to a temporary file.
func WriteScriptToTemp(script string) (string, error) {
	return writeScriptToTemp(script)
}

// NormalizeRepoPath normalizes repository path.
func NormalizeRepoPath(path string) string {
	return normalizeRepoPath(path)
}

// SaveRepoPathToDB saves the repo path to SQLite database.
func SaveRepoPathToDB(path string) {
	saveRepoPathToDB(path)
}

// HasFlag checks if flag exists.
func HasFlag(flagName string) bool {
	return hasFlag(flagName)
}

// ExpandTilde expands leading ~ to home directory.
func ExpandTilde(path string) string {
	return expandTilde(path)
}
