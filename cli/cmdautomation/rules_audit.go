package cmdautomation

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	reMapTuple = regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?\w+\s*\([^)]*\)\s*\(\s*map\[[^\]]+\][^,]+,\s*(?:error|\*apperror\.AppError)\s*\)`)
	reSliceTuple = regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?\w+\s*\([^)]*\)\s*\(\s*\[\][a-zA-Z0-9_*.]+\s*,\s*(?:error|\*apperror\.AppError)\s*\)`)
	reRawRuneCast = regexp.MustCompile(`\brune\s*\(\s*\d+\s*\)`)
	reGoEnumNoType = regexp.MustCompile(`type\s+([A-Za-z]\w*?)\s+(?:string|int|int32|int64|byte)\b`)
)

// RunRuleAudit executes a targeted guideline check across repository files.
func RunRuleAudit(opts RuleAuditOptions) RuleAuditResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	files, err := collectTextFiles(root, resolveRuleExtensions(opts.RuleName))
	if err != nil {
		return result.Fail[RuleAuditResult](err)
	}
	res := auditFilesForRule(files, opts.RuleName)
	res.Duration = time.Since(start)
	res.IsClean = len(res.Violations) == 0
	return result.Ok(res)
}

func resolveRuleExtensions(rule string) []string {
	switch rule {
	case "result-wrapper", "params":
		return []string{".go"}
	case "enums":
		return []string{".go", ".ts", ".py"}
	default:
		return []string{".go"}
	}
}

func auditFilesForRule(files []string, rule string) RuleAuditResult {
	res := RuleAuditResult{ScannedFiles: len(files)}
	for _, f := range files {
		vios := auditSingleFileRule(f, rule)
		if len(vios) > 0 {
			res.Violations = append(res.Violations, vios...)
		}
	}
	return res
}

func auditSingleFileRule(path, rule string) []RuleViolation {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	switch rule {
	case "result-wrapper":
		return auditResultWrapperContent(path, content)
	case "params":
		return auditParamsContent(path, content)
	case "enums":
		return auditEnumsContent(path, content)
	default:
		return nil
	}
}

func auditResultWrapperContent(path, content string) []RuleViolation {
	var vios []RuleViolation
	if reMapTuple.MatchString(content) {
		vios = append(vios, RuleViolation{
			File:        path,
			Rule:        "RESULT_MAP_TUPLE",
			Description: "Multi-value (map, error) return found; use ResultMap[K, V] instead",
		})
	}
	if hasNonPrimitiveSliceTuple(content) {
		vios = append(vios, RuleViolation{
			File:        path,
			Rule:        "RESULT_SLICE_TUPLE",
			Description: "Multi-value ([]T, error) return found; use ResultSlice[T] instead",
		})
	}
	return vios
}

func hasNonPrimitiveSliceTuple(content string) bool {
	matches := reSliceTuple.FindAllString(content, -1)
	for _, m := range matches {
		if !strings.Contains(m, "[]byte") && !strings.Contains(m, "[]rune") {
			return true
		}
	}
	return false
}

func auditParamsContent(path, content string) []RuleViolation {
	var vios []RuleViolation
	lines := strings.Split(content, "\n")
	for idx, line := range lines {
		stripped := strings.TrimSpace(line)
		if strings.HasPrefix(stripped, "func ") && countParamsInSignature(stripped) > 4 {
			vios = append(vios, RuleViolation{
				File:        path,
				LineNumber:  idx + 1,
				Rule:        "HIGH_ARITY_PARAMS",
				Description: "Function exceeds 4 parameters; extract to *Params struct in types.go",
			})
		}
	}
	return vios
}

func countParamsInSignature(sig string) int {
	start := strings.Index(sig, "(")
	end := strings.Index(sig, ")")
	if start == -1 || end == -1 || end <= start {
		return 0
	}
	paramsPart := sig[start+1 : end]
	if len(strings.TrimSpace(paramsPart)) == 0 {
		return 0
	}
	return strings.Count(paramsPart, ",") + 1
}

func auditEnumsContent(path, content string) []RuleViolation {
	var vios []RuleViolation
	lines := strings.Split(content, "\n")
	for idx, line := range lines {
		stripped := strings.TrimSpace(line)
		if reRawRuneCast.MatchString(stripped) {
			vios = append(vios, RuleViolation{
				File:        path,
				LineNumber:  idx + 1,
				Rule:        "RAW_RUNE_CAST",
				Description: "Raw numeric rune cast (e.g. rune(10)); use defined constant or rune literal",
			})
		}
		if isGoEnumMissingTypeSuffix(path, stripped) {
			vios = append(vios, RuleViolation{
				File:        path,
				LineNumber:  idx + 1,
				Rule:        "ENUM_MISSING_TYPE_SUFFIX",
				Description: "Enum missing mandatory 'Type' suffix in name",
			})
		}
	}
	return vios
}

func isGoEnumMissingTypeSuffix(path, line string) bool {
	if !strings.HasSuffix(path, ".go") || !strings.HasPrefix(line, "type ") {
		return false
	}
	match := reGoEnumNoType.FindStringSubmatch(line)
	if len(match) > 1 {
		name := match[1]
		return !strings.HasSuffix(name, "Type") && !strings.HasSuffix(name, "Monad")
	}
	return false
}
