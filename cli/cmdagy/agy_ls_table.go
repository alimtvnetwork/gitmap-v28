// Package cmdagy — agy_ls_table.go handles table row structures and dimensions for Antigravity projects.
package cmdagy

type agyTableRow struct {
	ID        string
	Name      string
	Branch    string
	Status    string
	Updated   string
	Path      string
	IsMissing bool
}

type agyTableContext struct {
	Rows       []agyTableRow
	MaxProject int
	MaxID      int
	MaxBranch  int
	MaxStatus  int
	MaxUpdated int
}

func newAgyTableContext() *agyTableContext {
	return &agyTableContext{
		Rows:       make([]agyTableRow, 0),
		MaxProject: 18,
		MaxID:      10,
		MaxBranch:  8,
		MaxStatus:  14,
		MaxUpdated: 9,
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
	c.MaxBranch = max(c.MaxBranch, len(r.Branch))
	c.MaxUpdated = max(c.MaxUpdated, len(r.Updated))
}

func buildAgyTableRow(p AgyProject) agyTableRow {
	path := p.GetPath()
	isMissing := path != "" && !checkDirExists(path)
	status, displayPath := resolveRowStatusAndPath(path)

	return agyTableRow{
		ID:        shortProjectId(p.ID),
		Name:      truncateMiddle(p.Name, 26),
		Branch:    resolveRowBranch(p.GetBranch()),
		Status:    status,
		Updated:   formatRelativeTime(p.UpdatedAt),
		Path:      truncateMiddle(displayPath, 42),
		IsMissing: isMissing,
	}
}

func resolveRowBranch(branch string) string {
	if branch == "" {
		return "—"
	}

	return branch
}

func resolveRowStatusAndPath(path string) (string, string) {
	if path == "" {
		return "global", "—"
	}

	return "active", path
}
