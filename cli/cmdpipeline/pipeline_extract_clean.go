package cmdpipeline

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

func parseLogLine(raw string) (string, string, string, bool) {
	clean := ansiRegex.ReplaceAllString(raw, "")
	parts := strings.Split(clean, "\t")
	isError := strings.Contains(raw, "##[error]") || hasFailureMarker(raw)
	if len(parts) >= 3 {
		return extractThreePartLogLine(parts, isError)
	}

	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), "", cleanLogText(parts[1]), isError
	}

	return "", "", cleanLogText(clean), isError
}

func extractThreePartLogLine(parts []string, isError bool) (string, string, string, bool) {
	job := strings.TrimSpace(parts[0])
	step := strings.TrimSpace(parts[1])
	text := stripTimestamp(strings.Join(parts[2:], "\t"))

	return job, step, cleanLogText(text), isError || hasFailureMarker(text)
}

func stripTimestamp(text string) string {
	if idx := strings.Index(text, "Z "); idx != -1 {
		return text[idx+2:]
	}

	if idx := strings.Index(text, "Z\t"); idx != -1 {
		return text[idx+2:]
	}

	return text
}

func cleanLogText(text string) string {
	t := strings.TrimPrefix(text, "##[error]")
	t = strings.TrimPrefix(t, "##[group]")
	t = strings.TrimPrefix(t, "##[endgroup]")
	t = cleanAnnotationError(t)

	return strings.TrimSpace(t)
}

func cleanAnnotationError(text string) string {
	idx := strings.Index(text, "::error")
	if idx < 0 {
		return text
	}

	return formatExtractedAnnotation(text[idx+len("::error"):])
}

func formatExtractedAnnotation(sub string) string {
	sepIdx := strings.Index(sub, "::")
	if sepIdx < 0 {
		return sub
	}

	meta := strings.TrimSpace(sub[:sepIdx])
	msg := strings.TrimSpace(sub[sepIdx+2:])
	loc := formatAnnotationLocation(meta)
	if len(loc) > 0 {
		return loc + ": " + msg
	}

	return msg
}

func formatAnnotationLocation(meta string) string {
	if len(meta) == 0 {
		return ""
	}

	file, line, col := parseAnnotationParts(meta)
	if len(file) == 0 {
		return ""
	}

	return buildLocationString(file, line, col)
}

func buildLocationString(file, line, col string) string {
	loc := file
	if len(line) > 0 {
		loc += ":" + line
	}
	if len(col) > 0 {
		loc += ":" + col
	}

	return loc
}

func parseAnnotationParts(meta string) (string, string, string) {
	var file, line, col string
	for _, part := range strings.Split(meta, ",") {
		k, v := parseKeyValue(part)
		switch k {
		case "file":
			file = v
		case "line":
			line = v
		case "col":
			col = v
		}
	}

	return file, line, col
}

func parseKeyValue(part string) (string, string) {
	kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
	if len(kv) == 2 {
		return kv[0], kv[1]
	}

	return "", ""
}

func isIgnoredLogLine(text string) bool {
	trimmed := strings.TrimSpace(text)
	if len(trimmed) == 0 {
		return true
	}
	if isRunnerCleanupNoise(trimmed) || isRunnerSetupNoise(trimmed) {
		return true
	}

	return isToolProgressNoise(trimmed)
}

func isRunnerCleanupNoise(trimmed string) bool {
	lower := strings.ToLower(trimmed)
	if strings.Contains(lower, "post job cleanup") || strings.Contains(lower, "safe.directory") {
		return true
	}
	if strings.Contains(lower, "removing ssh command") || strings.Contains(lower, "removing http extra header") {
		return true
	}
	if strings.Contains(lower, "removing includeif") || strings.Contains(lower, "includeif.gitdir:") {
		return true
	}
	if strings.Contains(lower, "git-credentials-") || strings.Contains(lower, "orphan process") {
		return true
	}
	if strings.Contains(lower, "temporarily overriding home") || strings.Contains(lower, "adding repository directory") {
		return true
	}

	return strings.HasPrefix(trimmed, "[command]/usr/bin/git") || strings.HasPrefix(trimmed, "git version ")
}

func isRunnerSetupNoise(trimmed string) bool {
	lower := strings.ToLower(trimmed)
	if strings.Contains(lower, "runner image") || strings.Contains(lower, "hosted compute agent") {
		return true
	}
	if strings.Contains(lower, "azure region:") || strings.Contains(lower, "github_token permissions") {
		return true
	}
	if strings.Contains(lower, "secret source:") || strings.Contains(lower, "prepare workflow directory") {
		return true
	}
	if strings.Contains(lower, "prepare all required actions") || strings.Contains(lower, "action download info") {
		return true
	}
	if strings.Contains(lower, "syncing repository:") || strings.Contains(lower, "disabling automatic garbage collection") {
		return true
	}
	if strings.Contains(lower, "setting up auth") || strings.Contains(lower, "fetching the repository") {
		return true
	}
	if strings.Contains(lower, "checking out the ref") || strings.Contains(lower, "switched to a new branch") {
		return true
	}

	return strings.HasPrefix(trimmed, "hint: ") || strings.Contains(lower, "node.js 20 is deprecated")
}

var ranTestsNoiseRegex = lazyregex.New(`^Ran\s+\d+\s+tests?\s+in\s+`)

func isToolProgressNoise(trimmed string) bool {
	if isRanTestsNoise(trimmed) {
		return true
	}
	if strings.HasPrefix(trimmed, "go: downloading ") || strings.HasPrefix(trimmed, "go install ") {
		return true
	}
	if strings.Contains(trimmed, "Checking boolean & enum compliance:") {
		return true
	}
	if strings.HasPrefix(trimmed, "Successfully set up CPython") || strings.HasPrefix(trimmed, "Complete job name:") {
		return true
	}
	if strings.HasPrefix(trimmed, "Browserslist:") || strings.Contains(trimmed, "caniuse-lite") || strings.Contains(trimmed, "npx update-browserslist-db") {
		return true
	}
	if strings.HasPrefix(trimmed, "transforming...") || strings.HasPrefix(trimmed, "vite v") || strings.Contains(trimmed, "building client environment") {
		return true
	}
	if strings.HasPrefix(trimmed, "Info Looking up installed tauri packages") || strings.HasPrefix(trimmed, "Running beforeBuildCommand") {
		return true
	}

	return strings.HasPrefix(trimmed, "Installed versions")
}

func isRanTestsNoise(trimmed string) bool {
	if !strings.HasPrefix(trimmed, "Ran ") {
		return false
	}

	return ranTestsNoiseRegex.CompileMust().MatchString(trimmed)
}

func hasFailureMarker(line string) bool {
	for _, m := range failureMarkers {
		if strings.Contains(line, m) {
			return true
		}
	}

	return false
}

func getOrCreateJobItem(jobMap map[string]*FailedJobItem, order *[]string, key, job, step string) *FailedJobItem {
	if item, exists := jobMap[key]; exists {
		return item
	}

	if job == "" {
		job = "Workflow Execution"
	}
	if step == "UNKNOWN STEP" || step == "" {
		step = "Job Execution"
	}

	item := &FailedJobItem{
		JobName:  job,
		StepName: step,
	}

	jobMap[key] = item
	*order = append(*order, key)

	return item
}
