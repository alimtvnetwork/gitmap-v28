package cmdlogviewer

// CommitItem represents a single commit row in the upper pane.
type CommitItem struct {
	Hash      string `json:"hash"`
	ShortHash string `json:"shortHash"`
	Author    string `json:"author"`
	Date      string `json:"date"`
	Subject   string `json:"subject"`
}

// FileChange represents a modified file row in the lower pane.
type FileChange struct {
	Action  string `json:"action"`
	Path    string `json:"path"`
	OldPath string `json:"oldPath,omitempty"`
	Stats   string `json:"stats,omitempty"`
}

// CommitDetail represents full commit metadata and file changes for the middle/lower pane.
type CommitDetail struct {
	Hash        string       `json:"hash"`
	ShortHash   string       `json:"shortHash"`
	Author      string       `json:"author"`
	AuthorEmail string       `json:"authorEmail"`
	Date        string       `json:"date"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	Parents     []string     `json:"parents"`
	Files       []FileChange `json:"files"`
}

// DiffRequest defines the payload sent to trigger external compare tools.
type DiffRequest struct {
	Commit string `json:"commit"`
	Path   string `json:"path"`
	Tool   string `json:"tool,omitempty"`
}

// DiffResponse carries the result of a diff invocation.
type DiffResponse struct {
	IsSuccess  bool   `json:"isSuccess"`
	Tool       string `json:"tool,omitempty"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
	StackTrace string `json:"stackTrace,omitempty"`
}

// BranchCompareRequest defines parameters to compare two git branches.
type BranchCompareRequest struct {
	BranchA string `json:"branchA"`
	BranchB string `json:"branchB"`
	Path    string `json:"path,omitempty"`
}

// BranchCompareResponse carries the result of a branch comparison invocation.
type BranchCompareResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

// BranchesResponse lists local branches and the active checkout.
type BranchesResponse struct {
	Branches     []string `json:"branches"`
	ActiveBranch string   `json:"activeBranch"`
}
