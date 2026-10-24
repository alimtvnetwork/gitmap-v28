package constants

// gitmap:cmd top-level
const (
	CmdSCWrite               = "servers-clients write"            // gitmap:cmd top-level
	CmdSCSetDefaultPath      = "servers-clients set-default-path" // gitmap:cmd top-level
	CmdSCSetPathAlias        = "servers-clients set-path-alias"   // gitmap:cmd top-level
	CmdSCUpdate              = "servers-clients update"           // gitmap:cmd top-level
	CmdSCUpdateAll           = "servers-clients update-all"       // gitmap:cmd top-level
	CmdClusterHistory        = "cluster history"                  // gitmap:cmd top-level
	CmdClusterExport         = "cluster export"                   // gitmap:cmd top-level
	CmdClusterImport         = "cluster import"                   // gitmap:cmd top-level
	CmdClusterSetPassword    = "cluster set-password"             // gitmap:cmd top-level
	CmdService               = "service"                          // gitmap:cmd top-level
	CmdServiceAlias          = "srv"
	CmdGroupCreate           = "create" // gitmap:cmd skip
	CmdGroupAdd              = "add"    // gitmap:cmd skip
	CmdGroupRemove           = "remove" // gitmap:cmd skip
	CmdGroupList             = "list"   // gitmap:cmd skip
	CmdGroupShow             = "show"   // gitmap:cmd skip
	CmdGroupDelete           = "delete" // gitmap:cmd skip
	CmdDBReset               = "db-reset"
	CmdDbReset               = CmdDBReset // gitmap:cmd skip
	CmdDB                    = "db"
	CmdDb                    = CmdDB // gitmap:cmd skip
	CmdStartFresh            = "start-fresh"
	CmdFindDuplicates        = "find-duplicates"
	CmdFindDuplicatesAlias   = "dups"
	CmdFindDuplicatesAlias2  = "find-dups"
	CmdReset                 = "reset"
	CmdListVersions          = "list-versions"
	CmdListVersionsAlias     = "lv"
	CmdRevert                = "revert"
	CmdRevertRunner          = "revert-runner" // gitmap:cmd skip
	CmdRm                    = "rm"
	CmdRmAlias               = "remove"
	CmdRmAlias2              = "del"
	CmdListReleases          = "list-releases"
	CmdListReleasesAlias     = "lr"
	CmdReleases              = "releases" // v3.20.0: alias of list-releases, intended for --all-repos batch view
	CmdNodesRunningOnly      = "nodes-running-only"
	CmdNodesRunningOnlyAlias = "nro"
)
