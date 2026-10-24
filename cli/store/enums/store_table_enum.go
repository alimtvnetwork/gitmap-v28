package enums

// StoreTableType defines canonical table names across GitMap SQLite stores.
type StoreTableType string

const (
	TableRepo                  StoreTableType = "Repo"
	TableGroup                 StoreTableType = "Group"
	TablePipelineRun           StoreTableType = "PipelineRun"
	TableSplitDatabaseRegistry StoreTableType = "SplitDatabaseRegistry"
	TableInstallerPins         StoreTableType = "installer_pins"
	TableSiteRegistry          StoreTableType = "SiteRegistry"
	TableSSHHosts              StoreTableType = "ssh_hosts"
	TableParentTask            StoreTableType = "ParentTask"
	TableSubtask               StoreTableType = "Subtask"
	TableCompletedTask         StoreTableType = "CompletedTask"
	TablePendingTask           StoreTableType = "PendingTask"
)

// IsValid reports whether the table type matches a known canonical table.
func (t StoreTableType) IsValid() bool {
	switch t {
	case TableRepo, TableGroup, TablePipelineRun, TableSplitDatabaseRegistry,
		TableInstallerPins, TableSiteRegistry, TableSSHHosts, TableParentTask,
		TableSubtask, TableCompletedTask, TablePendingTask:
		return true
	default:
		return false
	}
}

// String returns the string representation of the table name.
func (t StoreTableType) String() string {
	return string(t)
}

// TableName returns the canonical table identifier.
func (t StoreTableType) TableName() string {
	return string(t)
}

// AllStoreTables returns all registered StoreTableType variants.
func AllStoreTables() []StoreTableType {
	return []StoreTableType{
		TableRepo, TableGroup, TablePipelineRun, TableSplitDatabaseRegistry,
		TableInstallerPins, TableSiteRegistry, TableSSHHosts, TableParentTask,
		TableSubtask, TableCompletedTask, TablePendingTask,
	}
}
