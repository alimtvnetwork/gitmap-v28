// Package main — Cross-Language Coding Guidelines Validator (Go Edition)
//
// Version: 1.5.0  (2026-04-19) — Added P2/P3/P5/P7 boolean-principle checks.
//
// A Go port of validate-guidelines.py that validates Go, PHP, TypeScript,
// and Rust source files against the coding guidelines defined in
// spec/02-coding-guidelines/03-coding-guidelines-spec/.
//
// Usage:
//
//	go run scripts/validate-guidelines.go [--path <dir>] [--json] [--max-lines <n>]
//
// Rules Enforced:
//
//	CODE-RED-001  No nested if statements
//	CODE-RED-002  Boolean naming (is/has/can/should/was prefix)         [P1]
//	CODE-RED-003  No magic strings in comparisons
//	CODE-RED-004  Max 15 lines per function
//	CODE-RED-005  No fmt.Errorf() in Go (use apperror)
//	CODE-RED-006  No (T, error) returns in Go services (use Result[T])
//	CODE-RED-007  No string-based Go enums (use byte + iota)
//	CODE-RED-008  No raw string error codes — use apperrtype enum
//	CODE-RED-011  No magic numbers in logic
//	CODE-RED-012  Immutable by default (const over let/var, no reassignment)
//	CODE-RED-022  Negative words in boolean identifiers (isNot*, hasNo*) [P2]
//	CODE-RED-023  Raw `!` on function/method calls                       [P3]
//	CODE-RED-024  Bare true/false as positional argument                 [P5]
//	CODE-RED-025  Assignment inside if/while condition                   [P7]
//	STYLE-001     Blank line before return
//	STYLE-002     No else after return
//	STYLE-003     Blank line after closing brace
//	STYLE-004     Blank line before if/else if block
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ═══════════════════════════════════════════════════════════════════════
// Types
// ═══════════════════════════════════════════════════════════════════════

type Violation struct {
	File        string `json:"file"`
	Line        int    `json:"line"`
	Rule        string `json:"rule"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	CodeSnippet string `json:"code_snippet,omitempty"`
}

type ValidationReport struct {
	TotalFiles      int                `json:"totalFiles"`
	TotalViolations int                `json:"totalViolations"`
	CodeRedCount    int                `json:"codeRedCount"`
	StyleCount      int                `json:"styleCount"`
	ByRule          map[string]int     `json:"byRule"`
	Violations      []Violation        `json:"violations"`
	byFile          map[string]int
}

func newReport() *ValidationReport {
	return &ValidationReport{
		ByRule:     make(map[string]int),
		byFile:     make(map[string]int),
		Violations: []Violation{},
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Language Detection
// ═══════════════════════════════════════════════════════════════════════

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".php":
		return "php"
	case ".rs":
		return "rust"
	default:
		return ""
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Utility
// ═══════════════════════════════════════════════════════════════════════

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max]
}

func isExemptBoolName(name string) bool {
	exempt := map[string]bool{
		"ok": true, "done": true, "found": true,
		"exists": true, "err": true, "error": true,
		"true": true, "false": true,
	}

	return exempt[strings.ToLower(name)]
}

func hasBoolPrefix(name string) bool {
	prefixes := []string{"is", "has", "can", "should", "was", "will"}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) && len(name) > len(p) {
			next := name[len(p)]
			if next >= 'A' && next <= 'Z' {
				return true
			}
		}
	}

	return false
}

// ═══════════════════════════════════════════════════════════════════════
// Rule Checkers
// ═══════════════════════════════════════════════════════════════════════

func checkNestedIf(lines []string, path string) []Violation {
	var violations []Violation
	type stackEntry struct {
		indent int
		line   int
	}

	var stack []stackEntry

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		if strings.HasPrefix(stripped, "if ") || strings.HasPrefix(stripped, "if(") {
			currentIndent := len(line) - len(strings.TrimLeft(line, " \t"))

			for _, parent := range stack {
				if currentIndent > parent.indent {
					violations = append(violations, Violation{
						File:        path,
						Line:        i + 1,
						Rule:        "CODE-RED-001",
						Severity:    "CODE-RED",
						Message:     "Nested `if` is forbidden. Flatten with early returns or combined conditions.",
						CodeSnippet: truncate(stripped, 120),
					})
					break
				}
			}

			stack = append(stack, stackEntry{currentIndent, i})
		}

		if stripped == "}" && len(stack) > 0 {
			currentIndent := len(line) - len(strings.TrimLeft(line, " \t"))
			var newStack []stackEntry
			for _, e := range stack {
				if e.indent < currentIndent {
					newStack = append(newStack, e)
				}
			}

			stack = newStack
		}
	}

	return violations
}

var (
	goBoolPattern = regexp.MustCompile(`(\w+)\s*:=\s*(true|false)\b`)
	phpBoolPattern = regexp.MustCompile(`\$(\w+)\s*=\s*(?i)(true|false)\b`)
	tsBoolPattern  = regexp.MustCompile(`(?:const|let|var)\s+(\w+)\s*(?::\s*boolean)?\s*=\s*(true|false)\b`)
)

func checkBooleanNaming(lines []string, path string, lang string) []Violation {
	var violations []Violation
	var pattern *regexp.Regexp

	switch lang {
	case "go":
		pattern = goBoolPattern
	case "php":
		pattern = phpBoolPattern
	case "typescript", "javascript":
		pattern = tsBoolPattern
	default:
		return nil
	}

	for i, line := range lines {
		matches := pattern.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			name := m[1]
			if isExemptBoolName(name) || strings.HasPrefix(name, "_") || hasBoolPrefix(name) {
				continue
			}

			violations = append(violations, Violation{
				File:        path,
				Line:        i + 1,
				Rule:        "CODE-RED-002",
				Severity:    "CODE-RED",
				Message:     fmt.Sprintf(`Boolean variable "%s" must start with is/has/can/should/was/will.`, name),
				CodeSnippet: truncate(strings.TrimSpace(line), 120),
			})
		}
	}

	return violations
}

var (
	goMagicPattern  = regexp.MustCompile(`(?:==|!=)\s*"([^"]+)"`)
	phpMagicPattern = regexp.MustCompile(`(?:===|!==|==|!=)\s*['"]([^'"]+)['"]`)
	tsMagicPattern  = regexp.MustCompile(`(?:===|!==)\s*['"]([^'"]+)['"]`)
)

func checkMagicStrings(lines []string, path string, lang string) []Violation {
	var violations []Violation

	exempt := map[string]bool{
		"": true, "GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
		"string": true, "number": true, "boolean": true, "undefined": true,
		"object": true, "function": true, "utf-8": true, "utf8": true,
		// UI/layout prop values — not business logic
		"horizontal": true, "vertical": true, "left": true, "right": true,
		"top": true, "bottom": true, "sm": true, "md": true, "lg": true, "xl": true, "2xl": true,
		"none": true, "auto": true, "popper": true, "dot": true, "line": true, "dashed": true,
		"collapsed": true, "expanded": true, "floating": true, "inset": true,
		"sidebar": true, "header": true, "footer": true,
		"default": true, "destructive": true, "outline": true, "secondary": true, "ghost": true, "link": true,
	}

	var pattern *regexp.Regexp
	switch lang {
	case "go":
		pattern = goMagicPattern
	case "php":
		pattern = phpMagicPattern
	case "typescript", "javascript":
		pattern = tsMagicPattern
	default:
		return nil
	}

	for i, line := range lines {
		stripped := strings.TrimSpace(line)
		if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") || strings.HasPrefix(stripped, "*") {
			continue
		}

		matches := pattern.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			value := m[1]
			if exempt[value] || len(value) <= 1 {
				continue
			}

			if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "http") || strings.HasPrefix(value, ".") {
				continue
			}

			violations = append(violations, Violation{
				File:        path,
				Line:        i + 1,
				Rule:        "CODE-RED-003",
				Severity:    "CODE-RED",
				Message:     fmt.Sprintf(`Magic string "%s" in comparison. Use an enum constant.`, value),
				CodeSnippet: truncate(stripped, 120),
			})
		}
	}

	return violations
}

var (
	goFuncPattern  = regexp.MustCompile(`^func\s+(?:\(.*?\)\s+)?(\w+)`)
	phpFuncPattern = regexp.MustCompile(`(?:public|private|protected|static)?\s*function\s+(\w+)`)
	tsFuncPattern  = regexp.MustCompile(`(?:function\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?(?:\([^)]*\)|[^=])\s*=>)`)
)

func checkFunctionLength(lines []string, path string, lang string, maxLines int) []Violation {
	var violations []Violation
	var pattern *regexp.Regexp

	switch lang {
	case "go":
		pattern = goFuncPattern
	case "php":
		pattern = phpFuncPattern
	case "typescript", "javascript":
		pattern = tsFuncPattern
	default:
		return nil
	}

	for i := 0; i < len(lines); i++ {
		m := pattern.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}

		funcName := "anonymous"
		for _, g := range m[1:] {
			if g != "" {
				funcName = g
				break
			}
		}

		// Find opening brace
		braceLine := i
		for braceLine < len(lines) && !strings.Contains(lines[braceLine], "{") {
			braceLine++
		}

		if braceLine >= len(lines) {
			continue
		}

		// Count body lines
		depth := 0
		bodyLines := 0
		for j := braceLine; j < len(lines); j++ {
			depth += strings.Count(lines[j], "{") - strings.Count(lines[j], "}")
			if j > braceLine {
				stripped := strings.TrimSpace(lines[j])
				if stripped != "" && !strings.HasPrefix(stripped, "//") && !strings.HasPrefix(stripped, "*") && stripped != "}" {
					bodyLines++
				}
			}

			if depth <= 0 && j > braceLine {
				break
			}
		}

		if bodyLines > maxLines {
			violations = append(violations, Violation{
				File:        path,
				Line:        i + 1,
				Rule:        "CODE-RED-004",
				Severity:    "CODE-RED",
				Message:     fmt.Sprintf(`Function "%s" has %d lines (max %d). Extract into smaller helpers.`, funcName, bodyLines, maxLines),
				CodeSnippet: truncate(strings.TrimSpace(lines[i]), 120),
			})
		}
	}

	return violations
}

var (
	goStringEnumPattern      = regexp.MustCompile(`^type\s+\w+\s+string\s*$`)
	goTupleReturnPattern     = regexp.MustCompile(`func\s.*\)\s*\(\*?\w+,\s*error\)`)
	goRawErrorCodePattern    = regexp.MustCompile(`apperror\.(New|Wrap|Fail\w*)\(\s*"E\d+`)
)

func checkGoSpecific(lines []string, path string) []Violation {
	var violations []Violation

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		// CODE-RED-005: No fmt.Errorf()
		if strings.Contains(stripped, "fmt.Errorf(") && !strings.HasPrefix(stripped, "//") {
			violations = append(violations, Violation{
				File:     path,
				Line:     i + 1,
				Rule:     "CODE-RED-005",
				Severity: "CODE-RED",
				Message:  "Use apperror.New()/apperror.Wrap() instead of fmt.Errorf().",
				CodeSnippet: truncate(stripped, 120),
			})
		}

		// CODE-RED-007: No string-based enums
		if goStringEnumPattern.MatchString(stripped) {
			violations = append(violations, Violation{
				File:     path,
				Line:     i + 1,
				Rule:     "CODE-RED-007",
				Severity: "CODE-RED",
				Message:  "String-based enums forbidden. Use `type Variant byte` with iota.",
				CodeSnippet: truncate(stripped, 120),
		})
		}

		// CODE-RED-008: No raw string error codes — use apperrtype enum
		if goRawErrorCodePattern.MatchString(stripped) && !strings.HasPrefix(stripped, "//") {
			violations = append(violations, Violation{
				File:        path,
				Line:        i + 1,
				Rule:        "CODE-RED-008",
				Severity:    "CODE-RED",
				Message:     "Raw string error code (e.g. \"E2010\") found. Use apperrtype enum: apperror.NewType(apperrtype.SiteNotFound).",
				CodeSnippet: truncate(stripped, 120),
			})
		}
	}

	// CODE-RED-006: No (T, error) returns in service functions
	lowerPath := strings.ToLower(path)
	if !strings.Contains(lowerPath, "test") && !strings.HasSuffix(lowerPath, "_test.go") {
		for i, line := range lines {
			if goTupleReturnPattern.MatchString(line) {
				violations = append(violations, Violation{
					File:     path,
					Line:     i + 1,
					Rule:     "CODE-RED-006",
					Severity: "CODE-RED",
					Message:  "Service functions must return Result[T], not (T, error).",
					CodeSnippet: truncate(strings.TrimSpace(line), 120),
				})
			}
		}
	}

	return violations
}

var (
	magicNumberPattern      = regexp.MustCompile(`(?:==|!=|===|!==|>=|<=|[><]|\*|/|%)\s*(-?\d+\.?\d*)`)
	bracketIndexPattern     = regexp.MustCompile(`\[\s*\d+\s*\]`)
	tailwindNumberPattern   = regexp.MustCompile(`[a-z]-\d|/\d|opacity-|z-\d|gap-|p-|m-|w-|h-|text-\d|rounded-`)
	stringContextPattern    = regexp.MustCompile(`className=|class=|["'` + "`" + `].*[-/]\d`)
	jsxTextNumberPattern    = regexp.MustCompile(`>\s*\d+\s*</`)
)

func checkMagicNumbers(lines []string, path string, lang string) []Violation {
	var violations []Violation

	exempt := map[string]bool{
		"0": true, "1": true, "-1": true, "0.0": true, "1.0": true, "100": true,
		"2": true, "404": true,
	}

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") || strings.HasPrefix(stripped, "*") {
			continue
		}

		if strings.HasPrefix(stripped, "import") || strings.HasPrefix(stripped, "require") {
			continue
		}

		if strings.HasPrefix(stripped, "const ") || strings.HasPrefix(stripped, "const(") {
			continue
		}

		if strings.Contains(stripped, "= iota") {
			continue
		}

		// Skip lines with Tailwind/CSS class strings
		if tailwindNumberPattern.MatchString(stripped) || stringContextPattern.MatchString(stripped) {
			continue
		}

		// Skip JSX text content like <h1>404</h1>
		if jsxTextNumberPattern.MatchString(stripped) {
			continue
		}

		// Remove array indices before checking
		cleaned := bracketIndexPattern.ReplaceAllString(line, "[]")

		matches := magicNumberPattern.FindAllStringSubmatch(cleaned, -1)
		for _, m := range matches {
			value := m[1]
			if exempt[value] || value == "" {
				continue
			}

			lowerStripped := strings.ToLower(stripped)
			if strings.Contains(lowerStripped, "line") || strings.Contains(lowerStripped, "column") {
				continue
			}

			violations = append(violations, Violation{
				File:        path,
				Line:        i + 1,
				Rule:        "CODE-RED-011",
				Severity:    "CODE-RED",
				Message:     fmt.Sprintf("Magic number %s in logic. Use a named constant.", value),
				CodeSnippet: truncate(stripped, 120),
			})
		}
	}

	return violations
}

var (
	tsLetVarPattern    = regexp.MustCompile(`^\s*(?:let|var)\s+(\w+)\s*(?::\s*\w+)?\s*=`)
	goDeclPattern      = regexp.MustCompile(`^\s+(\w+)\s*:=\s+`)
	goReassignPattern  = regexp.MustCompile(`^\s+(\w+)\s*=\s+`)
)

func checkVariableMutation(lines []string, path string, lang string) []Violation {
	var violations []Violation

	if lang == "typescript" || lang == "javascript" {
		for i, line := range lines {
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "/*") || strings.HasPrefix(stripped, "*") {
				continue
			}

			m := tsLetVarPattern.FindStringSubmatch(line)
			if m == nil {
				continue
			}

			name := m[1]
			// Exemptions: loop vars, React hooks
			exemptNames := map[string]bool{"i": true, "j": true, "k": true, "idx": true, "index": true}
			if exemptNames[name] {
				continue
			}

			if strings.Contains(line, "useState") || strings.Contains(line, "useRef") {
				continue
			}

			if strings.HasPrefix(stripped, "for ") || strings.HasPrefix(stripped, "for(") {
				continue
			}

			violations = append(violations, Violation{
				File:        path,
				Line:        i + 1,
				Rule:        "CODE-RED-012",
				Severity:    "CODE-RED",
				Message:     fmt.Sprintf(`Use "const" instead of "let"/"var" for "%s". Assign once.`, name),
				CodeSnippet: truncate(stripped, 120),
			})
		}
	} else if lang == "go" {
		declared := make(map[string]int) // name -> line number

		exemptGoNames := map[string]bool{
			"err": true, "ok": true, "ctx": true, "cancel": true, "mu": true, "wg": true,
		}

		for i, line := range lines {
			stripped := strings.TrimSpace(line)
			if strings.HasPrefix(stripped, "//") {
				continue
			}

			if dm := goDeclPattern.FindStringSubmatch(line); dm != nil {
				declared[dm[1]] = i + 1
			}

			if rm := goReassignPattern.FindStringSubmatch(line); rm != nil {
				name := rm[1]
				declLine, wasDeclared := declared[name]
				if wasDeclared && (i+1)-declLine > 1 {
					if exemptGoNames[name] || strings.HasPrefix(name, "_") {
						continue
					}

					violations = append(violations, Violation{
						File:        path,
						Line:        i + 1,
						Rule:        "CODE-RED-012",
						Severity:    "CODE-RED",
						Message:     fmt.Sprintf(`Variable "%s" reassigned (first declared L%d). Avoid mutation.`, name, declLine),
						CodeSnippet: truncate(stripped, 120),
					})
				}
			}
		}
	}

	return violations
}

// ═══════════════════════════════════════════════════════════════════════
// Style Checks
// ═══════════════════════════════════════════════════════════════════════

func checkStyleRules(lines []string, path string) []Violation {
	var violations []Violation

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		// STYLE-001: Blank line before return
		if strings.HasPrefix(stripped, "return ") || stripped == "return" {
			if i > 0 {
				prevStripped := strings.TrimSpace(lines[i-1])
				if prevStripped != "" && prevStripped != "{" && !strings.HasPrefix(prevStripped, "//") {
					violations = append(violations, Violation{
						File:     path,
						Line:     i + 1,
						Rule:     "STYLE-001",
						Severity: "STYLE",
						Message:  "Add a blank line before `return` statement.",
						CodeSnippet: truncate(stripped, 120),
					})
				}
			}
		}

		// STYLE-002: No else after return
		if stripped == "} else {" || stripped == "} else if" || strings.HasPrefix(stripped, "} else ") {
			if i > 0 {
				// Look back for a return in the preceding block
				for j := i - 1; j >= 0 && j >= i-5; j-- {
					ps := strings.TrimSpace(lines[j])
					if strings.HasPrefix(ps, "return ") || ps == "return" {
						violations = append(violations, Violation{
							File:     path,
							Line:     i + 1,
							Rule:     "STYLE-002",
							Severity: "STYLE",
							Message:  "No `else` after `return`. Use early return pattern.",
							CodeSnippet: truncate(stripped, 120),
						})
						break
					}

					if ps == "}" || ps == "{" {
						break
					}
				}
			}
		}

		// STYLE-003: Blank line after closing brace (when followed by code)
		if stripped == "}" && i+1 < len(lines) {
			nextStripped := strings.TrimSpace(lines[i+1])
			if nextStripped != "" && nextStripped != "}" && !strings.HasPrefix(nextStripped, "//") &&
				nextStripped != ")" && !strings.HasPrefix(nextStripped, "} else") &&
				!strings.HasPrefix(nextStripped, "case ") && nextStripped != "default:" {
				violations = append(violations, Violation{
					File:     path,
					Line:     i + 1,
					Rule:     "STYLE-003",
					Severity: "STYLE",
					Message:  "Add a blank line after closing `}`.",
					CodeSnippet: truncate(stripped, 120),
				})
			}
		}

		// STYLE-004: Blank line before if block (when preceded by a statement)
		if strings.HasPrefix(stripped, "if ") || strings.HasPrefix(stripped, "if(") {
			if i > 0 {
				prev := strings.TrimSpace(lines[i-1])
				if prev != "" && prev != "{" && prev != "}" &&
					!strings.HasPrefix(prev, "//") && !strings.HasPrefix(prev, "#") && !strings.HasPrefix(prev, "*") {
					violations = append(violations, Violation{
						File:        path,
						Line:        i + 1,
						Rule:        "STYLE-004",
						Severity:    "STYLE",
						Message:     "Add a blank line before `if` block.",
						CodeSnippet: truncate(stripped, 120),
					})
				}
			}
		}
	}

	return violations
}

// ═══════════════════════════════════════════════════════════════════════
// File Validation
// ═══════════════════════════════════════════════════════════════════════

func shouldSkip(path string) bool {
	normalized := strings.ReplaceAll(path, "\\", "/")
	skipPatterns := []string{
		"vendor/", "node_modules/", "dist/", ".min.", "_test.go",
		"components/ui/", // Auto-generated shadcn/ui — not business logic
	}

	for _, p := range skipPatterns {
		if strings.Contains(normalized, p) {
			return true
		}
	}

	base := filepath.Base(path)

	return strings.HasSuffix(base, ".test.ts") || strings.HasSuffix(base, ".spec.ts")
}

// ═══════════════════════════════════════════════════════════════════════
// Boolean Principles — P2/P3/P5/P7 (added v1.5.0)
// ═══════════════════════════════════════════════════════════════════════

var (
	negBoolPattern   = regexp.MustCompile(`\b(is|has|can|should|was|will)(Not|No)[A-Z_]\w*`)
	bangCallPattern  = regexp.MustCompile(`(?:^|[\s(=,&|!])!\s*[A-Za-z_$][\w$.]*\s*\(`)
	bareBoolArgPat   = regexp.MustCompile(`[A-Za-z_$][\w$.]*\s*\([^)\n]*?(?:^|[^.\w])(true|false)(?:[^.\w]|$)[^)\n]*\)`)
	assignInCondPat  = regexp.MustCompile(`\b(?:if|while)\s*\(?[^()=!<>]*?[^=!<>]=[^=][^)]*\)?\s*\{`)
	bangAllowedAfter = regexp.MustCompile(`!\s*[A-Za-z_$][\w$.]*\s*\(.*\)\s*!=`)
)

func isCommentOrEmpty(stripped string) bool {
	if stripped == "" {
		return true
	}

	return strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") || strings.HasPrefix(stripped, "*")
}

func checkNegativeWords(lines []string, path string) []Violation {
	var violations []Violation

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		if isCommentOrEmpty(stripped) {
			continue
		}

		match := negBoolPattern.FindString(line)

		if match == "" {
			continue
		}

		violations = append(violations, Violation{
			File:        path,
			Line:        i + 1,
			Rule:        "CODE-RED-022",
			Severity:    "CODE-RED",
			Message:     fmt.Sprintf(`Negative-word boolean identifier "%s". Use a positive synonym (e.g. isPending, isInvalid, lacksAccess).`, match),
			CodeSnippet: truncate(stripped, 120),
		})
	}

	return violations
}

func checkBangOnCall(lines []string, path string) []Violation {
	var violations []Violation

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		if isCommentOrEmpty(stripped) {
			continue
		}

		if !bangCallPattern.MatchString(line) {
			continue
		}

		if bangAllowedAfter.MatchString(line) {
			continue
		}

		violations = append(violations, Violation{
			File:        path,
			Line:        i + 1,
			Rule:        "CODE-RED-023",
			Severity:    "CODE-RED",
			Message:     "Raw `!` on a function/method call is forbidden. Use a positive guard function or semantic inverse method.",
			CodeSnippet: truncate(stripped, 120),
		})
	}

	return violations
}

func isExemptBoolArgCall(stripped string) bool {
	exempt := []string{"expect", "assert", "should", "describe", "it(", "test(", "fmt.Print", "log.Print", "console.log"}

	for _, c := range exempt {
		if strings.Contains(stripped, c) {
			return true
		}
	}

	return false
}

func checkBareBoolArgs(lines []string, path string) []Violation {
	var violations []Violation

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		if isCommentOrEmpty(stripped) {
			continue
		}

		if isExemptBoolArgCall(stripped) {
			continue
		}

		match := bareBoolArgPat.FindStringSubmatch(stripped)

		if match == nil {
			continue
		}

		violations = append(violations, Violation{
			File:        path,
			Line:        i + 1,
			Rule:        "CODE-RED-024",
			Severity:    "CODE-RED",
			Message:     fmt.Sprintf("Bare `%s` as positional argument. Use a named flag, options object, or dedicated method.", match[1]),
			CodeSnippet: truncate(stripped, 120),
		})
	}

	return violations
}

func checkAssignInCondition(lines []string, path string, lang string) []Violation {
	var violations []Violation

	for i, line := range lines {
		stripped := strings.TrimSpace(line)

		if isCommentOrEmpty(stripped) {
			continue
		}

		if lang == "go" && strings.Contains(stripped, ":=") {
			continue
		}

		if !assignInCondPat.MatchString(stripped) {
			continue
		}

		violations = append(violations, Violation{
			File:        path,
			Line:        i + 1,
			Rule:        "CODE-RED-025",
			Severity:    "CODE-RED",
			Message:     "Assignment inside if/while condition is forbidden. Hoist the assignment to a prior line.",
			CodeSnippet: truncate(stripped, 120),
		})
	}

	return violations
}

func validateFile(path string, maxLines int) []Violation {
	lang := detectLanguage(path)
	if lang == "" || shouldSkip(path) {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var violations []Violation

	// Universal rules
	violations = append(violations, checkNestedIf(lines, path)...)
	violations = append(violations, checkBooleanNaming(lines, path, lang)...)
	violations = append(violations, checkMagicStrings(lines, path, lang)...)
	violations = append(violations, checkMagicNumbers(lines, path, lang)...)
	violations = append(violations, checkFunctionLength(lines, path, lang, maxLines)...)
	violations = append(violations, checkVariableMutation(lines, path, lang)...)
	violations = append(violations, checkStyleRules(lines, path)...)
	violations = append(violations, checkNegativeWords(lines, path)...)
	violations = append(violations, checkBangOnCall(lines, path)...)
	violations = append(violations, checkBareBoolArgs(lines, path)...)
	violations = append(violations, checkAssignInCondition(lines, path, lang)...)

	// Language-specific
	if lang == "go" {
		violations = append(violations, checkGoSpecific(lines, path)...)
	}

	return violations
}

// ═══════════════════════════════════════════════════════════════════════
// Report Output
// ═══════════════════════════════════════════════════════════════════════

func printReport(report *ValidationReport) {
	sep := strings.Repeat("=", 72)
	fmt.Println(sep)
	fmt.Println("  CODING GUIDELINES VALIDATION REPORT")
	fmt.Println(sep)
	fmt.Printf("  Files scanned:      %d\n", report.TotalFiles)
	fmt.Printf("  Total violations:   %d\n", report.TotalViolations)
	fmt.Printf("  🔴 CODE RED:        %d\n", report.CodeRedCount)
	fmt.Printf("  ⚠️  Style:           %d\n", report.StyleCount)
	fmt.Println(sep)

	if len(report.Violations) == 0 {
		fmt.Println("\n  ✅ ALL CLEAR — No violations found.\n")

		return
	}

	// Group by file
	byFile := make(map[string][]Violation)
	for _, v := range report.Violations {
		byFile[v.File] = append(byFile[v.File], v)
	}

	// Sort file keys
	var files []string
	for f := range byFile {
		files = append(files, f)
	}

	sort.Strings(files)

	for _, f := range files {
		vs := byFile[f]
		sort.Slice(vs, func(i, j int) bool { return vs[i].Line < vs[j].Line })
		fmt.Printf("\n📄 %s (%d violations)\n", f, len(vs))
		for _, v := range vs {
			icon := "🔴"
			if v.Severity != "CODE-RED" {
				icon = "⚠️ "
			}

			fmt.Printf("  %s L%-5d [%s] %s\n", icon, v.Line, v.Rule, v.Message)
			if v.CodeSnippet != "" {
				fmt.Printf("           │ %d: %s\n", v.Line, v.CodeSnippet)
			}
		}
	}

	fmt.Printf("\n%s\n", sep)
	fmt.Println("  BY RULE:")

	var rules []string
	for r := range report.ByRule {
		rules = append(rules, r)
	}

	sort.Strings(rules)

	for _, r := range rules {
		fmt.Printf("    %s: %d\n", r, report.ByRule[r])
	}

	fmt.Printf("%s\n\n", sep)

	if report.CodeRedCount > 0 {
		fmt.Printf("  ❌ FAILED — %d CODE RED violation(s) must be fixed before merge.\n\n", report.CodeRedCount)
	} else {
		fmt.Printf("  ⚠️  PASSED with %d style warning(s).\n\n", report.StyleCount)
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Main
// ═══════════════════════════════════════════════════════════════════════

func main() {
	scanPath := flag.String("path", "src", "Directory to scan (default: src)")
	outputJSON := flag.Bool("json", false, "Output as JSON")
	maxLines := flag.Int("max-lines", 15, "Max function lines (default: 15)")
	flag.Parse()

	report := newReport()

	extensions := map[string]bool{
		".go": true, ".ts": true, ".tsx": true,
		".js": true, ".jsx": true, ".php": true, ".rs": true,
	}

	err := filepath.Walk(*scanPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !extensions[ext] {
			return nil
		}

		report.TotalFiles++
		violations := validateFile(path, *maxLines)
		for _, v := range violations {
			report.Violations = append(report.Violations, v)
			report.TotalViolations++
			report.ByRule[v.Rule]++
			report.byFile[v.File]++
			if v.Severity == "CODE-RED" {
				report.CodeRedCount++
			} else {
				report.StyleCount++
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking path: %v\n", err)
		os.Exit(2)
	}

	if *outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(report)
	} else {
		printReport(report)
	}

	if report.CodeRedCount > 0 {
		os.Exit(1)
	}
}
