package cmdchromeprofile

type importDestination struct {
	Dir         string
	Path        string
	DisplayName string
	Email       string
	IsNew       bool
	Action      string
}

type snapshotMetadata struct {
	FilePath          string
	FileName          string
	FileSize          int64
	Export            *chromeExport
	BookmarksCount    int
	ExtensionsCount   int
	HasEmail          bool
	Email             string
	DisplayName       string
	ProfileName       string
	ExportedAt        string
	TargetDestination importDestination
}

type bookmarkRoot struct {
	Roots map[string]bookmarkNode `json:"roots"`
}

type bookmarkNode struct {
	Type     string         `json:"type"`
	Children []bookmarkNode `json:"children"`
}
