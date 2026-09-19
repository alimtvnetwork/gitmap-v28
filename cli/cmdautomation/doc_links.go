package cmdautomation

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	reMdLink = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
)

type docFileAudit struct {
	Violations []DocLinkViolation
	TotalLinks int
	FixedLinks int
}

// RunDocLinks audits and fixes markdown relative link integrity.
func RunDocLinks(opts DocLinksOptions) DocLinksResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	files := collectDocMarkdownFiles(root, opts.Dir)
	res := auditMarkdownFiles(root, files, opts.IsFix)
	res.Duration = time.Since(start)
	res.IsClean = len(res.Violations) == 0
	return result.Ok(res)
}

func collectDocMarkdownFiles(root, customDir string) []string {
	hasCustom := len(customDir) > 0
	if hasCustom {
		return walkMarkdownFiles(customDir)
	}
	var files []string
	targets := []string{"02-spec", ".ai-memory", "README.md"}
	for _, t := range targets {
		p := filepath.Join(root, t)
		files = append(files, resolveTargetFiles(p)...)
	}
	return files
}

func resolveTargetFiles(path string) []string {
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if fi.IsDir() {
		return walkMarkdownFiles(path)
	}
	isMd := strings.HasSuffix(strings.ToLower(path), ".md")
	if isMd {
		return []string{path}
	}
	return nil
}

func walkMarkdownFiles(dir string) []string {
	var files []string
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && isExcludedDir(d.Name()) {
			return filepath.SkipDir
		}
		isMd := !d.IsDir() && strings.HasSuffix(strings.ToLower(p), ".md")
		if isMd {
			files = append(files, p)
		}
		return nil
	})
	return files
}

func auditMarkdownFiles(root string, files []string, isFix bool) DocLinksResult {
	res := DocLinksResult{ScannedFiles: len(files)}
	for _, f := range files {
		audit := auditSingleDocFile(root, f, isFix)
		res.TotalLinks += audit.TotalLinks
		res.FixedLinks += audit.FixedLinks
		res.BrokenLinks += len(audit.Violations)
		if len(audit.Violations) > 0 {
			res.Violations = append(res.Violations, audit.Violations...)
		}
	}
	return res
}

func auditSingleDocFile(root, path string, isFix bool) docFileAudit {
	data, err := os.ReadFile(path)
	if err != nil {
		return docFileAudit{}
	}
	content := string(data)
	audit := scanContentLinks(root, path, content)
	if isFix {
		audit.FixedLinks = applyDocLinkFixes(path, content)
	}
	return audit
}

func scanContentLinks(root, path, content string) docFileAudit {
	var audit docFileAudit
	lines := strings.Split(content, "\n")
	fileDir := filepath.Dir(path)
	for idx, line := range lines {
		vios := checkLineDocLinks(root, path, fileDir, idx+1, line)
		matches := reMdLink.FindAllString(line, -1)
		audit.TotalLinks += len(matches)
		audit.Violations = append(audit.Violations, vios...)
	}
	return audit
}

func checkLineDocLinks(root, file, fileDir string, lineNum int, line string) []DocLinkViolation {
	var vios []DocLinkViolation
	matches := reMdLink.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		target := m[2]
		if isIgnoredLink(target) {
			continue
		}
		vio := validateLinkTarget(root, file, fileDir, lineNum, target)
		if vio.Issue != "" {
			vios = append(vios, vio)
		}
	}
	return vios
}

func isIgnoredLink(target string) bool {
	isHttp := strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://")
	isAnchor := strings.HasPrefix(target, "#")
	isMailto := strings.HasPrefix(target, "mailto:")
	return isHttp || isAnchor || isMailto
}

func validateLinkTarget(root, file, fileDir string, lineNum int, target string) DocLinkViolation {
	cleanTarget := stripLinkAnchor(target)
	resolved := resolveTargetOnDisk(root, fileDir, cleanTarget)
	_, err := os.Stat(resolved)
	isMissing := os.IsNotExist(err)
	if isMissing {
		return DocLinkViolation{
			File:           file,
			LineNumber:     lineNum,
			RawLink:        target,
			ResolvedTarget: resolved,
			Issue:          "Target file not found on disk",
		}
	}
	return DocLinkViolation{}
}

func stripLinkAnchor(target string) string {
	idx := strings.Index(target, "#")
	if idx >= 0 {
		return target[:idx]
	}
	return target
}

func resolveTargetOnDisk(root, fileDir, target string) string {
	relTarget := filepath.Join(fileDir, target)
	_, err := os.Stat(relTarget)
	if err == nil {
		return relTarget
	}
	isRepoRelative := strings.HasPrefix(target, "02-spec") ||
		strings.HasPrefix(target, ".ai-memory") ||
		strings.HasPrefix(target, "01-prompts") ||
		strings.HasPrefix(target, "03-ai-scripts")
	if isRepoRelative {
		return filepath.Join(root, target)
	}
	return relTarget
}

func applyDocLinkFixes(path, content string) int {
	fixed := content
	count := 0
	replacements := []string{
		".ai-memory/coding-guidelines/coding-guidelines.md", ".ai-memory/coding-guidelines.md",
		"lovable/coding-guidelines/coding-guidelines.md", ".ai-memory/coding-guidelines.md",
	}
	for i := 0; i < len(replacements); i += 2 {
		oldStr := replacements[i]
		newStr := replacements[i+1]
		if strings.Contains(fixed, oldStr) {
			fixed = strings.ReplaceAll(fixed, oldStr, newStr)
			count++
		}
	}
	writeDocFileIfChanged(path, content, fixed)
	return count
}

func writeDocFileIfChanged(path, original, modified string) {
	isChanged := original != modified
	if isChanged {
		_ = os.WriteFile(path, []byte(modified), 0o644)
	}
}
