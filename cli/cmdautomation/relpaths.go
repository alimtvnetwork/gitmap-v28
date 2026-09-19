package cmdautomation

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	reWinDrive = regexp.MustCompile(`(?i)[a-z]:[\\/][^ \r\n\t"'` + "`" + `<>]+`)
	reFileUri  = regexp.MustCompile(`(?i)file:///[^ \r\n\t"'` + "`" + `<>]+`)
)

// RunRelPathAudit scans target files for forbidden absolute filesystem paths.
func RunRelPathAudit(opts RelPathOptions) RelPathResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	files, err := collectTextFiles(root, opts.Extensions)
	if err != nil {
		return result.Fail[RelPathResult](err)
	}
	res := auditFilesForPaths(files, opts)
	res.Duration = time.Since(start)
	res.IsClean = len(res.Violations) == 0
	return result.Ok(res)
}

func resolveAuditDir(dir string) string {
	if len(dir) > 0 {
		return dir
	}
	return findRepoRoot()
}

func auditFilesForPaths(files []string, opts RelPathOptions) RelPathResult {
	res := RelPathResult{ScannedFiles: len(files)}
	for _, f := range files {
		vios, fixed := auditSingleFilePaths(f, opts.IsFixMode)
		if len(vios) > 0 {
			res.Violations = append(res.Violations, vios...)
		}
		res.FixedCount += fixed
	}
	return res
}

func auditSingleFilePaths(path string, isFix bool) ([]PathViolation, int) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0
	}
	content := string(data)
	if !hasPathMarkers(content) {
		return nil, 0
	}
	return inspectAndFixContent(path, content, isFix)
}

func hasPathMarkers(content string) bool {
	return strings.Contains(content, "file:") ||
		strings.Contains(content, `:\`) ||
		strings.Contains(content, ":/")
}

func inspectAndFixContent(path, content string, isFix bool) ([]PathViolation, int) {
	vios := findPathViolationsInText(path, content)
	if len(vios) == 0 {
		return nil, 0
	}
	fixed := applyPathFixIfEnabled(path, content, isFix)
	return vios, fixed
}
