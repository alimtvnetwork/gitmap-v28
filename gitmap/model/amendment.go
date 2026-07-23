// Package model — amendment.go defines the amendment audit record.
package model

// AmendmentRecord represents a single author-rewrite operation.
type AmendmentRecord struct {
	ID             int64         `json:"id"`
	Timestamp      string        `json:"timestamp"`
	Branch         string        `json:"branch"`
	FromCommit     string        `json:"fromCommit"`
	ToCommit       string        `json:"toCommit"`
	TotalCommits   int           `json:"totalCommits"`
	PreviousAuthor AmendAuthor   `json:"previousAuthor"`
	NewAuthor      AmendAuthor   `json:"newAuthor"`
	Mode           string        `json:"mode"`
	ForcePushed    bool          `json:"forcePushed"`
	Commits        []CommitEntry `json:"commits"`
}

// AmendAuthor holds a name/email pair for amendment records.
type AmendAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CommitEntry holds a single commit's SHA and message.
type CommitEntry struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
}
