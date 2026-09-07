package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readManifestProfilesFromDir(dir string) ([]DiscoveredProfileCandidate, bool) {
	mPath := filepath.Join(dir, "manifest.json")
	raw, err := os.ReadFile(mPath)
	if err != nil {
		return nil, false
	}

	var m chromeProfileManifest
	if err := json.Unmarshal(raw, &m); err != nil || len(m.Profiles) == 0 {
		return nil, false
	}

	var candidates []DiscoveredProfileCandidate
	for _, p := range m.Profiles {
		pDir := filepath.Join(dir, p.Name)
		cand, ok := readSubdirProfileCandidate(dir, p.Name)
		if !ok {
			cand = buildFallbackCandidateFromManifest(pDir, p)
		}
		candidates = append(candidates, cand)
	}

	return candidates, true
}

func buildFallbackCandidateFromManifest(pDir string, p chromeManifestProfile) DiscoveredProfileCandidate {
	exp := &chromeExport{Name: p.Name, DisplayName: p.DisplayName, Email: p.Email}
	target := resolveImportDestination(exp, "", false)

	return DiscoveredProfileCandidate{
		SourcePath:        pDir,
		ProfileDirName:    p.Name,
		DisplayName:       p.DisplayName,
		Email:             p.Email,
		ExtensionsCount:   p.ExtensionCount,
		TargetDestination: target,
	}
}
