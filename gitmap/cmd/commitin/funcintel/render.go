package funcintel

import (
	"sort"
	"strings"
)

// Render builds the §6.3 per-file added-function block. Files where
// nothing was added AND that were not newly-added are omitted.
func Render(changes []FileChange) string {
	sorted := append([]FileChange(nil), changes...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].RelativePath < sorted[j].RelativePath
	})
	var b strings.Builder
	for _, c := range sorted {
		appendFileBlock(&b, c)
	}
	return strings.TrimRight(b.String(), "\n")
}

func appendFileBlock(b *strings.Builder, c FileChange) {
	added := detectAdded(c)
	if len(added) == 0 && !c.NewlyAddedFile {
		return
	}
	b.WriteString("- ")
	b.WriteString(c.RelativePath)
	b.WriteByte('\n')
	if len(added) > 0 {
		b.WriteString("  - added: ")
		b.WriteString(strings.Join(added, ", "))
		b.WriteByte('\n')
	}
}

func detectAdded(c FileChange) []string {
	if c.Language == "" {
		return nil
	}
	d := Get(c.Language)
	if d == nil {
		return nil
	}
	return d.Detect(c.PrevSource, c.NewSource)
}
