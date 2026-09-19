package cmdautomation

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	reExplicitTrue   = regexp.MustCompile(`===?\s*(?:true|True)\b|!==?\s*(?:true|True)\b`)
	reNonAffirmParam = regexp.MustCompile(`\b(?:trigger|showRole|forceLifecycle|autoConfirm)\s+bool\b`)
	reNonAffirmField = regexp.MustCompile(`^\s+(?:Exists|Empty)\s+bool\b`)
)

// RunNamingAudit scans repository source files for boolean naming and comparison anti-patterns.
func RunNamingAudit(opts NamingOptions) NamingResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	files, err := collectTextFiles(root, opts.Extensions)
	if err != nil {
		return result.Fail[NamingResult](err)
	}
	res := auditFilesForNaming(files)
	res.Duration = time.Since(start)
	res.IsClean = len(res.Violations) == 0
	return result.Ok(res)
}

func auditFilesForNaming(files []string) NamingResult {
	res := NamingResult{ScannedFiles: len(files)}
	for _, f := range files {
		vios := auditSingleFileNaming(f)
		if len(vios) > 0 {
			res.Violations = append(res.Violations, vios...)
		}
	}
	return res
}

func auditSingleFileNaming(path string) []NamingViolation {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if !hasNamingMarkers(content) {
		return nil
	}
	return inspectNamingContent(path, content)
}

func hasNamingMarkers(content string) bool {
	return strings.Contains(content, "true") ||
		strings.Contains(content, "True") ||
		strings.Contains(content, "bool")
}

func inspectNamingContent(path, content string) []NamingViolation {
	var vios []NamingViolation
	lines := strings.Split(content, "\n")
	for idx, line := range lines {
		stripped := strings.TrimSpace(line)
		if isIgnoredCommentLine(stripped) {
			continue
		}
		vios = append(vios, checkLineNamingViolations(path, idx+1, stripped)...)
	}
	return vios
}

func checkLineNamingViolations(file string, lineNum int, line string) []NamingViolation {
	var vios []NamingViolation
	if reExplicitTrue.MatchString(line) {
		vios = append(vios, NamingViolation{
			File:        file,
			LineNumber:  lineNum,
			Kind:        "EXPLICIT_TRUE",
			LineContent: line,
		})
	}
	if reNonAffirmParam.MatchString(line) || reNonAffirmField.MatchString(line) {
		vios = append(vios, NamingViolation{
			File:        file,
			LineNumber:  lineNum,
			Kind:        "NON_AFFIRMATIVE_BOOLEAN",
			LineContent: line,
		})
	}
	return vios
}

func isIgnoredCommentLine(line string) bool {
	return strings.HasPrefix(line, "//") ||
		strings.HasPrefix(line, "/*") ||
		strings.HasPrefix(line, "*") ||
		strings.HasPrefix(line, "#")
}
