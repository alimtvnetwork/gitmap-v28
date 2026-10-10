package cmdclone

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CloneDispatchMode defines which cloning pipeline handles discovered repositories.
type CloneDispatchMode string

const (
	CloneModeStandard CloneDispatchMode = "standard" // standard git clone
	CloneModeCFR      CloneDispatchMode = "cfr"      // private clone-fix-repo
	CloneModeCFRP     CloneDispatchMode = "cfrp"     // public clone-fix-repo-pub
)

// ResolveRepoCacheRootFn is a test hook to override the repo-cache root directory.
var ResolveRepoCacheRootFn func() string

// ResolveRepoCacheRoot resolves the directory path for repo-cache.
func ResolveRepoCacheRoot() string {
	if ResolveRepoCacheRootFn != nil {
		return ResolveRepoCacheRootFn()
	}
	if info, err := os.Stat("repo-cache"); err == nil && info.IsDir() {
		return "repo-cache"
	}

	db, err := store.OpenSpecialReposSplitDB()
	if err == nil {
		defer db.Close()
		rec, err := db.GetSpecialRepo("repo-cache")
		if err == nil && rec != nil {
			if len(rec.LocalPath) > 0 {
				if info, err := os.Stat(rec.LocalPath); err == nil && info.IsDir() {
					return rec.LocalPath
				}
			}
			workBase := resolveWorkBaseDir()
			confName := rec.ConfiguredName
			if len(confName) == 0 {
				confName = "repo-cache"
			}
			target := filepath.Join(workBase, confName)
			_ = os.MkdirAll(target, 0755)
			return target
		}
	}

	workBase := resolveWorkBaseDir()
	target := filepath.Join(workBase, "repo-cache")
	_ = os.MkdirAll(target, 0755)
	return target
}

func resolveWorkBaseDir() string {
	mainDB, err := store.OpenDefault()
	if err == nil {
		defer mainDB.Close()
		if wd, err := mainDB.GetDefaultWorkDir(); err == nil && wd != nil && len(wd.AbsolutePath) > 0 {
			if info, err := os.Stat(wd.AbsolutePath); err == nil && info.IsDir() {
				return wd.AbsolutePath
			}
		}
	}
	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return `D:\work`
	}
	cwd, err := os.Getwd()
	if err == nil {
		return filepath.Dir(cwd)
	}
	return "."
}

func hasRCFlagOrToken(args []string) bool {
	for _, a := range args {
		normalized := strings.ToLower(strings.TrimSpace(a))
		if normalized == "rc" || normalized == "repo-cache" ||
			normalized == "--rc" || normalized == "-rc" ||
			normalized == "--repo-cache" {
			return true
		}
	}
	return false
}

type rcParsedFlags struct {
	isAssumeYes      bool
	manifestIndex    int
	passthroughFlags []string
}

func parseRCArgs(args []string) rcParsedFlags {
	var pf rcParsedFlags
	for _, a := range args {
		low := strings.ToLower(strings.TrimSpace(a))
		if low == "rc" || low == "repo-cache" || low == "--rc" || low == "-rc" || low == "--repo-cache" {
			continue
		}
		if low == "-y" || low == "--yes" || low == "-yes" || low == "--assume-yes" {
			pf.isAssumeYes = true
			continue
		}
		if num, err := strconv.Atoi(a); err == nil && num > 0 {
			pf.manifestIndex = num
			continue
		}
		pf.passthroughFlags = append(pf.passthroughFlags, a)
	}
	return pf
}

// RunCloneRC handles high-speed repository ingestion and cloning from repo-cache manifests.
func RunCloneRC(args []string, mode CloneDispatchMode) error {
	repoCacheRoot := ResolveRepoCacheRoot()
	manifests, err := DiscoverRCManifests(repoCacheRoot)
	if err != nil {
		return fmt.Errorf("failed to discover repo-cache manifests in %s: %w", repoCacheRoot, err)
	}

	if len(manifests) == 0 {
		fmt.Printf("ℹ No repository manifests found in repo-cache (%s).\n", repoCacheRoot)
		return nil
	}

	flags := parseRCArgs(args)

	isInteractive := stdinIsTerminal() && !flags.isAssumeYes && flags.manifestIndex == 0
	if !isInteractive {
		targetIndex := flags.manifestIndex
		if targetIndex <= 0 || targetIndex > len(manifests) {
			targetIndex = 1
		}
		targetManifest := manifests[targetIndex-1]
		fmt.Printf("▸ Ingesting repo-cache manifest [%d]: %s (%d repositories)\n",
			targetManifest.Index, targetManifest.RelativePath, targetManifest.TotalRepos)
		return executeManifestClone(targetManifest, mode, flags.passthroughFlags)
	}

	return runInteractiveRCMenu(manifests, mode, flags.passthroughFlags)
}

func runInteractiveRCMenu(manifests []DiscoveredManifest, mode CloneDispatchMode, passthrough []string) error {
	fmt.Print(RenderShortTreeView(manifests))
	fmt.Println("\nChoose an option:")
	fmt.Println("  [A] Clone all repositories across all manifests")
	for _, m := range manifests {
		repoWord := "repos"
		if m.TotalRepos == 1 {
			repoWord = "repo"
		}
		fmt.Printf("  [%d] Clone manifest [%d] (%s - %d %s)\n", m.Index, m.Index, m.RelativePath, m.TotalRepos, repoWord)
	}
	fmt.Println("  [S] Select specific repositories from a manifest")
	fmt.Println("  [R] Rotate and inspect next manifest")
	fmt.Println("  [Q] Quit")
	fmt.Print("\nSelect [A/1/2/S/R/Q] (or pass -y to clone all): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := strings.ToUpper(strings.TrimSpace(input))

	switch choice {
	case "A":
		return cloneAllManifests(manifests, mode, passthrough)
	case "S":
		return runInteractiveSelectRepos(manifests, mode, passthrough, reader)
	case "R":
		fmt.Println("\nRotating manifests overview:")
		fmt.Print(RenderShortTreeView(manifests))
		return runInteractiveRCMenu(manifests, mode, passthrough)
	case "Q", "":
		fmt.Println("Operation cancelled.")
		return nil
	default:
		if num, err := strconv.Atoi(choice); err == nil && num >= 1 && num <= len(manifests) {
			target := manifests[num-1]
			fmt.Printf("\n▸ Cloning manifest [%d]: %s (%d repositories)...\n", target.Index, target.RelativePath, target.TotalRepos)
			return executeManifestClone(target, mode, passthrough)
		}
		fmt.Println("Invalid choice. Operation cancelled.")
		return nil
	}
}

func cloneAllManifests(manifests []DiscoveredManifest, mode CloneDispatchMode, passthrough []string) error {
	for _, m := range manifests {
		fmt.Printf("\n▸ Cloning manifest [%d]: %s (%d repositories)...\n", m.Index, m.RelativePath, m.TotalRepos)
		if err := executeManifestClone(m, mode, passthrough); err != nil {
			return err
		}
	}
	return nil
}

func runInteractiveSelectRepos(manifests []DiscoveredManifest, mode CloneDispatchMode, passthrough []string, reader *bufio.Reader) error {
	fmt.Printf("Select manifest number [1-%d]: ", len(manifests))
	numStr, _ := reader.ReadString('\n')
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil || num < 1 || num > len(manifests) {
		fmt.Println("Invalid manifest index.")
		return nil
	}

	target := manifests[num-1]
	fmt.Printf("\nRepositories in [%d] %s:\n", target.Index, target.RelativePath)
	for i, e := range target.Entries {
		proto := ClassifyProtocol(e.CloneUrl)
		fmt.Printf("  [%d] %s/%s %s\n", i+1, e.Owner, e.RepoName, proto)
	}

	fmt.Print("\nEnter repository numbers to clone (e.g. 1, 3) or 'all': ")
	selStr, _ := reader.ReadString('\n')
	selStr = strings.TrimSpace(selStr)
	if strings.EqualFold(selStr, "all") {
		return executeManifestClone(target, mode, passthrough)
	}

	var selectedEntries []RepoCacheEntry
	parts := strings.Split(selStr, ",")
	for _, p := range parts {
		idx, pErr := strconv.Atoi(strings.TrimSpace(p))
		if pErr == nil && idx >= 1 && idx <= len(target.Entries) {
			selectedEntries = append(selectedEntries, target.Entries[idx-1])
		}
	}

	if len(selectedEntries) == 0 {
		fmt.Println("No valid repositories selected.")
		return nil
	}

	return executeEntriesClone(selectedEntries, mode, passthrough)
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
