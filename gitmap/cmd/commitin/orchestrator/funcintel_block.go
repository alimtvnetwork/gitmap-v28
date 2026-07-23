package orchestrator

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/funcintel"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/walk"
)

// renderFunctionIntel returns the §6.3 block for one source commit, or
// "" when the feature is disabled / no detected language matches.
// Pure I/O helper — never returns an error (best-effort per spec §6.3:
// parser failures must NOT abort a commit).
func renderFunctionIntel(repoDir string, c walk.SourceCommit, files []string, fi profile.FunctionIntel) string {
	if !fi.IsEnabled || len(files) == 0 {
		return ""
	}
	enabled := funcintel.EnabledLanguages(fi.Languages)
	if len(enabled) == 0 {
		return ""
	}
	wanted := toSet(enabled)
	changes := collectChanges(repoDir, c.Sha, files, wanted)
	if len(changes) == 0 {
		return ""
	}
	return funcintel.Render(changes)
}

func collectChanges(repoDir, sha string, files []string, wanted map[string]struct{}) []funcintel.FileChange {
	out := make([]funcintel.FileChange, 0, len(files))
	for _, rel := range files {
		lang := funcintel.LanguageForPath(rel)
		if lang == "" {
			continue
		}
		if _, ok := wanted[lang]; !ok {
			continue
		}
		out = append(out, buildChange(repoDir, sha, rel, lang))
	}
	return out
}

func buildChange(repoDir, sha, rel, lang string) funcintel.FileChange {
	prev := readBlob(repoDir, sha+"^:"+rel)
	curr := readBlob(repoDir, sha+":"+rel)
	return funcintel.FileChange{
		RelativePath:   rel,
		Language:       lang,
		PrevSource:     prev,
		NewSource:      curr,
		NewlyAddedFile: prev == "" && curr != "",
	}
}

// readBlob runs `git -C <repoDir> show <ref>` and returns "" on any
// error (missing parent for first commit, deleted file, etc.). The
// caller treats "" as "did not exist before".
func readBlob(repoDir, ref string) string {
	out, err := exec.Command("git", "-C", repoDir, "show", ref).Output()
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(out), "\n")
}

func toSet(in []string) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for _, s := range in {
		out[s] = struct{}{}
	}
	return out
}
