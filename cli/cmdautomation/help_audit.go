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
	reCobraCmd  = regexp.MustCompile(`(?s)([a-zA-Z0-9_]+)\s*[:=]+\s*&cobra\.Command\s*\{([^}]+)\}`)
	reUseProp   = regexp.MustCompile(`Use:\s*"([^"\s]+)`)
	reShortProp = regexp.MustCompile(`Short:\s*"([^"]+)"`)
)

// RunHelpAudit audits CLI commands for Short descriptions and helptext docs.
func RunHelpAudit(opts HelpAuditOptions) HelpAuditResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	files := findCliSourceFiles(root, opts.Extensions)
	res := auditCliFiles(root, files)
	res.Duration = time.Since(start)
	res.IsClean = len(res.Violations) == 0
	return result.Ok(res)
}

func findCliSourceFiles(root string, extensions []string) []string {
	var files []string
	cliDir := filepath.Join(root, "cli")
	_ = filepath.WalkDir(cliDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if isCliSourceFile(path) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func isCliSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	isGo := ext == ".go"
	isTest := strings.HasSuffix(path, "_test.go")
	return isGo && !isTest
}

func auditCliFiles(root string, files []string) HelpAuditResult {
	res := HelpAuditResult{ScannedFiles: len(files)}
	helpDir := filepath.Join(root, "cli", "helptext")
	for _, f := range files {
		vios := auditFileCommands(f, helpDir)
		if len(vios) > 0 {
			res.Violations = append(res.Violations, vios...)
		}
	}
	res.TotalCommands = countTotalCommands(files)
	return res
}

func countTotalCommands(files []string) int {
	total := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		matches := reCobraCmd.FindAllString(string(data), -1)
		total += len(matches)
	}
	return total
}

func auditFileCommands(path, helpDir string) []HelpViolation {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if !strings.Contains(content, "cobra.Command") {
		return nil
	}
	return scanCobraViolations(path, content, helpDir)
}

func scanCobraViolations(path, content, helpDir string) []HelpViolation {
	var vios []HelpViolation
	matches := reCobraCmd.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		cmdVar := m[1]
		body := m[2]
		cmdVios := checkSingleCobraCommand(path, cmdVar, body, helpDir)
		vios = append(vios, cmdVios...)
	}
	return vios
}

func checkSingleCobraCommand(path, cmdVar, body, helpDir string) []HelpViolation {
	var vios []HelpViolation
	hasShort := reShortProp.MatchString(body)
	if !hasShort {
		vios = append(vios, HelpViolation{
			File:    path,
			Command: cmdVar,
			Issue:   "Missing Short description in cobra.Command",
		})
	}
	useMatch := reUseProp.FindStringSubmatch(body)
	if len(useMatch) > 1 {
		docVio := checkDocParity(path, cmdVar, useMatch[1], helpDir)
		if docVio.Issue != "" {
			vios = append(vios, docVio)
		}
	}
	return vios
}

func checkDocParity(path, cmdVar, cmdName, helpDir string) HelpViolation {
	docFile := filepath.Join(helpDir, cmdName+".md")
	_, err := os.Stat(docFile)
	isDocMissing := os.IsNotExist(err)
	if isDocMissing {
		return HelpViolation{
			File:    path,
			Command: cmdVar,
			Issue:   "Missing helptext doc: cli/helptext/" + cmdName + ".md",
		}
	}
	return HelpViolation{}
}
