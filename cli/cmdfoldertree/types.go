package cmdfoldertree

// FolderTreeNode represents a directory or file node in the hierarchy.
type FolderTreeNode struct {
	Sequence  int               `json:"sequence" yaml:"sequence"`
	Name      string            `json:"name" yaml:"name"`
	Path      string            `json:"path" yaml:"path"`
	RelPath   string            `json:"relPath" yaml:"relPath"`
	IsDir     bool              `json:"isDir" yaml:"isDir"`
	IsGit     bool              `json:"isGit" yaml:"isGit"`
	GitRepo   string            `json:"gitRepo,omitempty" yaml:"gitRepo,omitempty"`
	GitBranch string            `json:"gitBranch,omitempty" yaml:"gitBranch,omitempty"`
	Children  []*FolderTreeNode `json:"children,omitempty" yaml:"children,omitempty"`
}

// FolderTreeExportDoc represents the exported top-level document structure.
type FolderTreeExportDoc struct {
	RootPath   string          `json:"rootPath" yaml:"rootPath"`
	RootName   string          `json:"rootName" yaml:"rootName"`
	TotalNodes int             `json:"totalNodes" yaml:"totalNodes"`
	TotalDirs  int             `json:"totalDirs" yaml:"totalDirs"`
	TotalFiles int             `json:"totalFiles" yaml:"totalFiles"`
	Tree       *FolderTreeNode `json:"tree" yaml:"tree"`
}

// FolderTreeOptions holds CLI options for folder-tree commands.
type FolderTreeOptions struct {
	TargetDir     string
	Subcommand    string
	Format        string
	MaxDepth      int
	DirsOnly      bool
	ShowFiles     bool
	ShowNumbers   bool
	GitOnly       bool
	IncludeHidden bool
	OutputFile    string
	InputFile     string
	DryRun        bool
	IsPreview     bool
	IsTree        bool
}

// ImportSummary records the result of an import operation.
type ImportSummary struct {
	TargetDir    string   `json:"targetDir"`
	DirsCreated  int      `json:"dirsCreated"`
	FilesCreated int      `json:"filesCreated"`
	CreatedPaths []string `json:"createdPaths,omitempty"`
}
