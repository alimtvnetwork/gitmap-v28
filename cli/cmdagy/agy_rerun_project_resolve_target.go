package cmdagy

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func isCwdMatchingProject(p AgyProject, cleanCwd, cleanRealCwd string) bool {
	pWs := cleanProjectWorkspace(p.GetPath())
	if isWorkspaceMatchOrSubdir(pWs, cleanCwd, cleanRealCwd) {
		return true
	}
	realWs, err := filepath.EvalSymlinks(p.GetPath())
	if err != nil || realWs == "" {
		return false
	}
	cleanRealWs := cleanProjectWorkspace(realWs)

	return isWorkspaceMatchOrSubdir(cleanRealWs, cleanCwd, cleanRealCwd)
}

func isWorkspaceMatchOrSubdir(ws, cleanCwd, cleanRealCwd string) bool {
	if ws == "" {
		return false
	}
	if ws == cleanCwd || strings.HasPrefix(cleanCwd, ws+"/") {
		return true
	}
	if cleanRealCwd != "" && (ws == cleanRealCwd || strings.HasPrefix(cleanRealCwd, ws+"/")) {
		return true
	}

	return false
}

func resolveCleanRealCwd(cwd string) string {
	realCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil || realCwd == "" {
		return ""
	}
	return cleanProjectWorkspace(realCwd)
}

func findProjectByFlexibleTarget(projects []AgyProject, target string) (AgyProject, error) {
	if isHelpKeyword(target) {
		return AgyProject{}, apperror.NewNotFound("project_target", "E9039", "help is not a valid project target")
	}

	if proj, isFound := matchDirectTarget(projects, target); isFound {
		return proj, nil
	}

	return matchFuzzyTarget(projects, target)
}

func matchDirectTarget(projects []AgyProject, target string) (AgyProject, bool) {
	matched, err := ResolveAgyProjectTargets([]string{target}, projects)
	if err == nil && len(matched) > 0 {
		return matched[0], true
	}

	return AgyProject{}, false
}

func matchFuzzyTarget(projects []AgyProject, target string) (AgyProject, error) {
	cleanTarget := normalizeSlugForComparison(target)
	for _, p := range projects {
		if isProjectTargetMatch(p, cleanTarget) {
			return p, nil
		}
	}

	return AgyProject{}, apperror.NewNotFound("project_target", "E9039", "no project matching target: "+target)
}

func isHelpKeyword(target string) bool {
	return isRerunHelpToken(target)
}

func isProjectTargetMatch(p AgyProject, cleanTarget string) bool {
	if isHelpKeyword(cleanTarget) {
		return false
	}

	nameNorm := normalizeSlugForComparison(p.Name)
	wsNorm := normalizeSlugForComparison(filepath.Base(p.GetPath()))

	if strings.Contains(nameNorm, cleanTarget) || strings.Contains(wsNorm, cleanTarget) {
		return true
	}

	if cleanTarget != "" && (strings.Contains(cleanTarget, nameNorm) || strings.Contains(cleanTarget, wsNorm)) {
		return true
	}

	return isPhoneticOrAliasMatch(cleanTarget, nameNorm) || isPhoneticOrAliasMatch(cleanTarget, wsNorm)
}

func isPhoneticOrAliasMatch(target, candidate string) bool {
	if strings.Contains(target, "xampp") && strings.Contains(candidate, "exam") {
		return true
	}
	if strings.Contains(target, "exam") && strings.Contains(candidate, "xampp") {
		return true
	}
	return false
}

func normalizeSlugForComparison(raw string) string {
	s := strings.ToLower(raw)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}
