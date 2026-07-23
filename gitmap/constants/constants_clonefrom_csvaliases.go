package constants

// CSV header aliases for clone-from. Spreadsheets, copy-pasted docs,
// and hand-edited files in the wild use a handful of common
// variations for the same logical column ("URL" vs "httpsURL",
// "dest" vs "relpath", etc.). Rather than fail those files with an
// opaque "url column missing" error, we normalize them to the
// canonical column name before validation.
//
// Rules:
//   - Lookup is case-insensitive — keys are stored lowercase and the
//     header value is lowercased + trimmed before lookup.
//   - The canonical name itself is INTENTIONALLY included so callers
//     can do a single map lookup with no fallback branch.
//   - Aliases must be unambiguous: a single alias may map to exactly
//     one canonical column. Adding a conflicting alias is a bug.
//
// Keep this list short and obvious. Exotic spellings should be
// fixed in the source CSV, not silently accepted here.
var CSVColumnAliases = map[string]string{
	// url
	"url":      CSVColumnURL,
	"urls":     CSVColumnURL,
	"httpsurl": CSVColumnURL,
	"httpurl":  CSVColumnURL,
	"giturl":   CSVColumnURL,
	"repo":     CSVColumnURL,
	"repourl":  CSVColumnURL,
	"clone":    CSVColumnURL,
	"cloneurl": CSVColumnURL,

	// dest
	"dest":     CSVColumnDest,
	"relpath":  CSVColumnDest,
	"path":     CSVColumnDest,
	"folder":   CSVColumnDest,
	"dir":      CSVColumnDest,
	"target":   CSVColumnDest,
	"destpath": CSVColumnDest,

	// branch
	"branch": CSVColumnBranch,
	"ref":    CSVColumnBranch,
	"tag":    CSVColumnBranch,

	// depth
	"depth":      CSVColumnDepth,
	"clonedepth": CSVColumnDepth,

	// checkout
	"checkout":     CSVColumnCheckout,
	"checkoutmode": CSVColumnCheckout,
	"mode":         CSVColumnCheckout,
}

// CanonicalCSVColumn returns the canonical column name for a raw
// header cell, or "" if the header is unknown. Caller is responsible
// for skipping unknown columns silently — extra columns must not
// cause a hard parse failure.
func CanonicalCSVColumn(rawHeader string) string {
	return CSVColumnAliases[normalizeHeader(rawHeader)]
}

// normalizeHeader lowercases + trims a header cell and strips
// inner ASCII whitespace, underscores, and hyphens so "Https URL",
// "https_url", and "https-url" all collapse to "httpsurl".
func normalizeHeader(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '_' || c == '-' {
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		out = append(out, c)
	}

	return string(out)
}
