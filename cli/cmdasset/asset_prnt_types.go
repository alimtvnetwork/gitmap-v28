package cmdasset

// PrntOptions holds parsed command line arguments for the prnt tool.
type PrntOptions struct {
	URLs        []string
	CustomDest  string
	OutputDir   string
	ExtractOnly bool
	JsonFormat  bool
	Concurrency int
}

// PrntResultItem represents the outcome of resolving/downloading a single URL.
type PrntResultItem struct {
	SourceURL     string `json:"sourceUrl"`
	DirectURL     string `json:"directUrl,omitempty"`
	Slug          string `json:"slug,omitempty"`
	SavedPath     string `json:"savedPath,omitempty"`
	BrainSyncPath string `json:"brainSyncPath,omitempty"`
	FileSizeBytes int    `json:"fileSizeBytes,omitempty"`
	IsSuccess     bool   `json:"isSuccess"`
	Error         string `json:"error,omitempty"`
}

// PrntStatus models the status envelope header.
type PrntStatus struct {
	IsSuccess bool   `json:"isSuccess"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// PrntAttributes models the summary metadata of the operation.
type PrntAttributes struct {
	TotalRequested int  `json:"totalRequested"`
	TotalSucceeded int  `json:"totalSucceeded"`
	TotalFailed    int  `json:"totalFailed"`
	ExtractOnly    bool `json:"extractOnly"`
}

// PrntResponseEnvelope is the standard Universal JSON envelope returned when --json is requested.
type PrntResponseEnvelope struct {
	Status     PrntStatus       `json:"status"`
	Attributes PrntAttributes   `json:"attributes"`
	Results    []PrntResultItem `json:"results"`
	Errors     []string         `json:"errors,omitempty"`
}
