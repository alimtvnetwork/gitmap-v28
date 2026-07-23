package cmd

// parseHistoryPaths normalizes the multi-form path list accepted by
// history-purge / history-pin. The spec says paths may come as
// separate args, joined by `,`, or joined by `, ` (comma-space).
// Quoting is irrelevant — we always normalize all three forms. Empty
// tokens are dropped; order is preserved; duplicates are removed.
// splitOnComma (defined in root.go) already trims whitespace and
// drops empty pieces, so this function only needs to dedupe.
func parseHistoryPaths(args []string) []string {
	seen := make(map[string]bool, len(args))
	out := make([]string, 0, len(args))
	for _, raw := range args {
		for _, tok := range splitOnComma(raw) {
			if seen[tok] {
				continue
			}
			seen[tok] = true
			out = append(out, tok)
		}
	}
	return out
}
