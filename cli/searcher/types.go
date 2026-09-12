package searcher

type SearchResult struct {
	MatchedText   string `json:"matched_text" yaml:"matched_text"`
	StartPosition int    `json:"start_position" yaml:"start_position"`
	EndPosition   int    `json:"end_position" yaml:"end_position"`
	FilePath      string `json:"file_path" yaml:"file_path"`
	RelativePath  string `json:"relative_path" yaml:"relative_path"`
}
