package cmdos

// UpdateToolchain defines package managers discovered on the host system.
type UpdateToolchain struct {
	Name         string
	Binary       string
	UpdateArgs   []string
	UpgradeArgs  []string
	NeedsSudo    bool
	RequiresSudo bool
}

// UpdateResult records execution status for a package manager update.
type UpdateResult struct {
	Name    string
	Success bool
	Output  string
	Error   error
}
