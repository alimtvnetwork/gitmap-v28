package model

// TempRelease represents a temporary release branch record.
type TempRelease struct {
	ID             int64  `json:"id"`
	Branch         string `json:"branch"`
	VersionPrefix  string `json:"versionPrefix"`
	SequenceNumber int    `json:"sequenceNumber"`
	CommitSha      string `json:"commit"`
	CommitMessage  string `json:"commitMessage"`
	CreatedAt      string `json:"createdAt"`
}
