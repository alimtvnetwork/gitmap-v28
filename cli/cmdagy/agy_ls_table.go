// Package cmdagy — agy_ls_table.go handles table row structures and dimensions for Antigravity projects.
package cmdagy

type agyTableRow struct {
	Seq       int
	ID        string
	Name      string
	ConvName  string
	ConvID    string
	Branch    string
	Status    string
	Updated   string
	Path      string
	IsMissing bool
}

type agyTableContext struct {
	Rows        []agyTableRow
	MaxSeq      int
	MaxProject  int
	MaxID       int
	MaxConvName int
	MaxConvID   int
	MaxUpdated  int
}

func newAgyTableContext() *agyTableContext {
	return &agyTableContext{
		Rows:        make([]agyTableRow, 0),
		MaxSeq:      3,
		MaxProject:  16,
		MaxID:       10,
		MaxConvName: 16,
		MaxConvID:   10,
		MaxUpdated:  9,
	}
}

func truncateMiddle(s string, maxLen int) string {
	if len(s) <= maxLen || maxLen <= 5 {
		return s
	}

	prefixLen := (maxLen - 3) / 2
	suffixLen := maxLen - 3 - prefixLen

	return s[:prefixLen] + "..." + s[len(s)-suffixLen:]
}

func (c *agyTableContext) addRow(r agyTableRow) {
	c.Rows = append(c.Rows, r)
	c.updateMaxColumnWidths(r)
}

func (c *agyTableContext) updateMaxColumnWidths(r agyTableRow) {
	c.MaxProject = max(c.MaxProject, len(r.Name))
	c.MaxID = max(c.MaxID, len(r.ID))
	c.MaxConvName = max(c.MaxConvName, len(r.ConvName))
	c.MaxConvID = max(c.MaxConvID, len(r.ConvID))
	c.MaxUpdated = max(c.MaxUpdated, len(r.Updated))
}

func buildAgyTableRow(p AgyProject, convMap map[string]AgyLatestConv, seq int) agyTableRow {
	path := p.GetPath()
	isMissing := path != "" && !checkDirExists(path)
	pinnedMap := getPinnedProjectsSet()
	isPinned := pinnedMap[p.ID] || pinnedMap[p.Name]
	status, displayPath := resolveRowStatusAndPath(path, isPinned)
	convTitle, convID := resolveProjectConvDetails(p.ID, convMap)

	return agyTableRow{
		Seq:       seq,
		ID:        shortProjectId(p.ID),
		Name:      truncateMiddle(p.Name, 24),
		ConvName:  truncateMiddle(convTitle, 24),
		ConvID:    shortProjectId(convID),
		Branch:    resolveRowBranch(p.GetBranch()),
		Status:    status,
		Updated:   formatRelativeTime(p.UpdatedAt),
		Path:      truncateMiddle(displayPath, 42),
		IsMissing: isMissing,
	}
}
