package constants

// Workgroup and workspace table names (PascalCase).
const (
	TableWorkDirectory     = "WorkDirectory"
	TableFolderWorkDir     = "FolderWorkDir"
	TableFolderGroup       = "FolderGroup"
	TableSshKeyAssociation = "SshKeyAssociation"
	ViewWorkDirectories    = "work_directories"
)

// SSH Key association types.
const (
	AssocTypeRepo    = "repo"
	AssocTypeFolder  = "folder"
	AssocTypeGroup   = "group"
	AssocTypeWorkDir = "workdir"
	AssocTypeDefault = "default"
)

// WorkgroupCmdType defines strongly-typed workgroup subcommands.
type WorkgroupCmdType string

const (
	WorkgroupCmdList    WorkgroupCmdType = "list"
	WorkgroupCmdAdd     WorkgroupCmdType = "add"
	WorkgroupCmdBind    WorkgroupCmdType = "bind"
	WorkgroupCmdUnbind  WorkgroupCmdType = "unbind"
	WorkgroupCmdShow    WorkgroupCmdType = "show"
	WorkgroupCmdTestSSH WorkgroupCmdType = "test-ssh"
	WorkgroupCmdRemove  WorkgroupCmdType = "remove"
	WorkgroupCmdLs      WorkgroupCmdType = "ls"
	WorkgroupCmdTest    WorkgroupCmdType = "test"
	WorkgroupCmdRm      WorkgroupCmdType = "rm"
	WorkgroupCmdDelete  WorkgroupCmdType = "delete"
)

// WorkgroupFlagType defines strongly-typed CLI flags for workgroup commands.
type WorkgroupFlagType string

const (
	WorkgroupFlagSSHKey    WorkgroupFlagType = "--ssh-key"
	WorkgroupFlagDesc      WorkgroupFlagType = "--desc"
	WorkgroupFlagHost      WorkgroupFlagType = "--host"
	WorkgroupFlagUserName  WorkgroupFlagType = "--user-name"
	WorkgroupFlagUserEmail WorkgroupFlagType = "--user-email"
	WorkgroupFlagDefault   WorkgroupFlagType = "--default"
	WorkgroupFlagJSON      WorkgroupFlagType = "--json"
	WorkgroupFlagJSONShort WorkgroupFlagType = "-j"
	WorkgroupFlagHelp      WorkgroupFlagType = "--help"
	WorkgroupFlagHelpShort WorkgroupFlagType = "-h"
)
