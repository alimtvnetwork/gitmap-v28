package cmd

// StorageBreakdown encapsulates storage utilization metrics of a workspace.
type StorageBreakdown struct {
	RepoRoot         string `json:"repoRoot"`
	WorkingTreeBytes int64  `json:"workingTreeBytes"`
	WorkingTreeFiles int    `json:"workingTreeFiles"`
	GitDirBytes      int64  `json:"gitDirBytes"`
	GitObjectsBytes  int64  `json:"gitObjectsBytes"`
	DatabaseBytes    int64  `json:"databaseBytes"`
	DatabaseCount    int    `json:"databaseCount"`
	TempAndLogsBytes int64  `json:"tempAndLogsBytes"`
	TempFilesCount   int    `json:"tempFilesCount"`
	TotalRepoBytes   int64  `json:"totalRepoBytes"`
}

// StorageCleanOptions holds user-configured flags for storage clean.
type StorageCleanOptions struct {
	IsDryRun   bool `json:"isDryRun"`
	IsVerbose  bool `json:"isVerbose"`
	IsForce    bool `json:"isForce"`
	IsVacuumDB bool `json:"isVacuumDb"`
}

// StorageCleanStats captures reclaimed space metrics.
type StorageCleanStats struct {
	DeletedLogsCount int   `json:"deletedLogsCount"`
	DeletedTempCount int   `json:"deletedTempCount"`
	ReclaimedBytes   int64 `json:"reclaimedBytes"`
	VacuumFreedBytes int64 `json:"vacuumFreedBytes"`
}
