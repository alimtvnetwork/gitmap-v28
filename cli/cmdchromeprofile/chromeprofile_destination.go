package cmdchromeprofile

import (
	"fmt"
	"path/filepath"
	"strings"
)

func resolveImportDestination(exp *chromeExport, explicitTarget string, logSteps bool) importDestination {
	if explicitTarget != "" {
		return resolveExplicitDestination(explicitTarget, exp, logSteps)
	}

	if exp.Email != "" {
		return resolveDestinationByEmail(exp, logSteps)
	}

	return resolveDestinationDefault(exp, logSteps)
}

func resolveExplicitDestination(explicitTarget string, exp *chromeExport, logSteps bool) importDestination {
	path := chromeProfilePath(explicitTarget)
	isNew := !chromeProfilePathExists(path)
	logExplicitTargetStep(explicitTarget, logSteps)

	return importDestination{
		Dir:         explicitTarget,
		Path:        path,
		DisplayName: exp.DisplayName,
		Email:       exp.Email,
		IsNew:       isNew,
		Action:      fmt.Sprintf("Explicit target %q", explicitTarget),
	}
}

func logExplicitTargetStep(target string, logSteps bool) {
	if !logSteps {
		return
	}

	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Using explicit profile target: %q\n", target)
}

func resolveDestinationByEmail(exp *chromeExport, logSteps bool) importDestination {
	existingDir, found := findChromeProfileByEmail(exp.Email)
	if found {
		logExistingProfileMatched(existingDir, exp.Email, logSteps)

		return importDestination{
			Dir:         existingDir,
			Path:        chromeProfilePath(existingDir),
			DisplayName: exp.DisplayName,
			Email:       exp.Email,
			IsNew:       false,
			Action:      fmt.Sprintf("Matched existing profile %q (<%s>)", existingDir, exp.Email),
		}
	}

	nextDir := findNextAvailableProfileDir()
	dispName := resolveProfileDisplayName(exp)
	logNewProfileCreated(exp.Email, nextDir, logSteps)

	return importDestination{
		Dir:         nextDir,
		Path:        chromeProfilePath(nextDir),
		DisplayName: dispName,
		Email:       exp.Email,
		IsNew:       true,
		Action:      fmt.Sprintf("New profile %q (<%s>)", nextDir, exp.Email),
	}
}

func logExistingProfileMatched(existingDir, email string, logSteps bool) {
	if !logSteps {
		return
	}

	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Matched existing Chrome profile %q by email <%s>\n", existingDir, email)
}

func logNewProfileCreated(email, nextDir string, logSteps bool) {
	if !logSteps {
		return
	}

	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Email <%s> not found in Chrome. Creating new profile %q (protecting existing profiles)\n", email, nextDir)
}

func resolveProfileDisplayName(exp *chromeExport) string {
	if exp.DisplayName != "" {
		return exp.DisplayName
	}

	if exp.Name != "" {
		return exp.Name
	}

	if exp.Email != "" {
		return strings.Split(exp.Email, "@")[0]
	}

	return "Profile"
}

func findChromeProfileByEmail(email string) (string, bool) {
	for _, dir := range availableChromeProfileNames() {
		_, existingEmail := resolveProfileNameAndEmail(dir, nil)
		if existingEmail != "" && strings.EqualFold(existingEmail, email) {
			return dir, true
		}
	}

	return "", false
}

func resolveDestinationDefault(exp *chromeExport, logSteps bool) importDestination {
	candidate := exp.Name
	if candidate == "" {
		candidate = "Default"
	}

	targetPath := chromeProfilePath(candidate)
	if !chromeProfilePathExists(targetPath) {
		logProfileAvailable(candidate, logSteps)

		return importDestination{
			Dir:         candidate,
			Path:        targetPath,
			DisplayName: exp.DisplayName,
			Email:       "",
			IsNew:       true,
			Action:      fmt.Sprintf("Available profile %q", candidate),
		}
	}

	_, existingEmail := resolveProfileNameAndEmail(candidate, nil)
	isOccupied := existingEmail != "" || hasProfileBookmarks(targetPath)
	if isOccupied {
		nextDir := findNextAvailableProfileDir()
		logProfileOccupied(candidate, nextDir, logSteps)

		return importDestination{
			Dir:         nextDir,
			Path:        chromeProfilePath(nextDir),
			DisplayName: exp.DisplayName,
			Email:       "",
			IsNew:       true,
			Action:      fmt.Sprintf("Allocating new profile %q (protecting occupied %q)", nextDir, candidate),
		}
	}

	logProfileUpdating(candidate, logSteps)

	return importDestination{
		Dir:         candidate,
		Path:        targetPath,
		DisplayName: exp.DisplayName,
		Email:       "",
		IsNew:       false,
		Action:      fmt.Sprintf("Updating existing profile %q", candidate),
	}
}

func logProfileAvailable(name string, logSteps bool) {
	if !logSteps {
		return
	}

	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Target profile directory %q is available\n", name)
}

func logProfileOccupied(candidate, nextDir string, logSteps bool) {
	if !logSteps {
		return
	}

	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Profile %q is occupied; allocating new profile %q to protect existing data\n", candidate, nextDir)
}

func logProfileUpdating(candidate string, logSteps bool) {
	if !logSteps {
		return
	}

	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Updating profile %q\n", candidate)
}

func isExcepted(exceptRules []string, fileName, profileName, displayName, email string) (bool, string) {
	for _, rule := range exceptRules {
		r := strings.TrimSpace(rule)
		if r == "" {
			continue
		}

		if matchExceptRule(r, fileName) {
			return true, r
		}

		baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		if matchExceptRule(r, baseName) {
			return true, r
		}

		if matchExceptRule(r, profileName) {
			return true, r
		}

		if matchExceptRule(r, displayName) {
			return true, r
		}

		if matchExceptRule(r, email) {
			return true, r
		}
	}

	return false, ""
}

func matchExceptRule(rule, val string) bool {
	if val == "" {
		return false
	}

	lowRule := strings.ToLower(rule)
	lowVal := strings.ToLower(val)
	if strings.EqualFold(lowRule, lowVal) {
		return true
	}

	if strings.HasSuffix(lowRule, "*") {
		prefix := strings.TrimSuffix(lowRule, "*")

		return strings.HasPrefix(lowVal, prefix)
	}

	return strings.HasPrefix(lowVal, lowRule)
}
