package cmdclone

import (
	"fmt"
)

func cloneAllManifests(manifests []DiscoveredManifest, mode CloneDispatchMode, passthrough []string) error {
	for _, m := range manifests {
		fmt.Printf("\n▸ Cloning manifest [%d]: %s (%d repositories)...\n", m.Index, m.RelativePath, m.TotalRepos)
		if err := executeManifestClone(m, mode, passthrough); err != nil {
			return err
		}
	}
	return nil
}

func executeManifestClone(manifest DiscoveredManifest, mode CloneDispatchMode, passthrough []string) error {
	switch mode {
	case CloneModeCFR:
		f := cloneFixRepoFlags{
			url:     manifest.FullPath,
			autoYes: true,
		}
		return runCFRManifestPipeline(f, false, CfrModifierFlags{})
	case CloneModeCFRP:
		f := cloneFixRepoFlags{
			url:     manifest.FullPath,
			autoYes: true,
		}
		return runCFRManifestPipeline(f, true, CfrModifierFlags{PromotePublic: true})
	default:
		combined := append([]string{manifest.FullPath}, passthrough...)
		cf := parseCloneFlags(combined)
		return executeParsedClone(cf)
	}
}

func executeEntriesClone(entries []RepoCacheEntry, mode CloneDispatchMode, passthrough []string) error {
	for i, e := range entries {
		fmt.Printf("\n[%d/%d] Ingesting %s/%s (%s)...\n", i+1, len(entries), e.Owner, e.RepoName, e.CloneUrl)
		switch mode {
		case CloneModeCFR:
			if err := runCloneFixRepoPipeline([]string{e.CloneUrl, "-y"}, false); err != nil {
				return err
			}
		case CloneModeCFRP:
			if err := runCloneFixRepoPipeline([]string{e.CloneUrl, "-y"}, true); err != nil {
				return err
			}
		default:
			executeDirectClone(DirectCloneParams{
				URL:       e.CloneUrl,
				GHDesktop: true,
			})
		}
	}
	return nil
}
