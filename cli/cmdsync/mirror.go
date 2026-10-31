package cmdsync

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var specDirNumRe = lazyregex.New(`^(\d+)-`)

// MirrorStats records added, updated, removed, preserved, conflict, and in-sync counts.
type MirrorStats struct {
	Added     int
	Updated   int
	Removed   int
	Preserved int
	Conflicts int
	InSync    int
}

type syncDirMapping struct {
	relSrc  string
	relDst  string
	addOnly bool
}

// MirrorAssets synchronizes canonical assets using 3-Way Hash state evaluation.
func MirrorAssets(sourceRoot, targetRoot string) (MirrorStats, error) {
	return MirrorAssetsWithScope(sourceRoot, targetRoot, enums.SyncScopeAll, false)
}

// MirrorAssetsWithScope synchronizes canonical assets with designated scope and 3-way hash protection.
func MirrorAssetsWithScope(sourceRoot, targetRoot string, scope enums.SyncScopeType, isDryRun bool) (MirrorStats, error) {
	var stats MirrorStats
	manifest, _ := LoadTargetManifest(targetRoot)
	if manifest == nil {
		manifest = &SyncManifest{Assets: make(map[string]string)}
	}

	dirs := buildSyncDirMappingsWithScope(sourceRoot, scope)
	for _, d := range dirs {
		if err := mirrorSingleDirectoryMapping(d, sourceRoot, targetRoot, manifest, isDryRun, &stats); err != nil {
			return stats, err
		}
	}

	if !isDryRun {
		manifest.SectionHashes = ComputeSectionHashes(manifest.Assets)
		_ = SaveTargetManifest(targetRoot, manifest)
	}

	return stats, nil
}

func buildSyncDirMappingsWithScope(sourceRoot string, scope enums.SyncScopeType) []syncDirMapping {
	if scope.IsGitmapSkillOnly() {
		return []syncDirMapping{
			{filepath.FromSlash(".agents/skills/gitmap"), filepath.FromSlash(".agents/skills/gitmap"), false},
			{filepath.FromSlash(".cursor/skills/gitmap"), filepath.FromSlash(".cursor/skills/gitmap"), false},
		}
	}

	var dirs []syncDirMapping
	if scope.IsPrompts() {
		dirs = append(dirs, syncDirMapping{constants.DirPrompts, constants.DirPrompts, false})
	}
	if scope.IsSkills() {
		dirs = append(dirs,
			syncDirMapping{filepath.FromSlash(constants.DirAgentsSkills), filepath.FromSlash(constants.DirAgentsSkills), false},
			syncDirMapping{filepath.FromSlash(constants.DirCursorSkills), filepath.FromSlash(constants.DirCursorSkills), false},
		)
	}
	if scope.IsScripts() {
		dirs = append(dirs,
			syncDirMapping{constants.DirAiScripts, constants.DirAiScripts, true},
			syncDirMapping{filepath.FromSlash(constants.DirAgentsScripts), filepath.FromSlash(constants.DirAgentsScripts), true},
		)
	}
	if scope.IsSpecs() {
		specRoot := filepath.Join(sourceRoot, constants.DirSpec)
		dirs = appendSpecDirs(specRoot, dirs)
	}

	return dirs
}

func mirrorSingleDirectoryMapping(d syncDirMapping, sourceRoot, targetRoot string, manifest *SyncManifest, isDryRun bool, stats *MirrorStats) error {
	srcDir := filepath.Join(sourceRoot, d.relSrc)
	dstDir := filepath.Join(targetRoot, d.relDst)

	if _, err := os.Stat(srcDir); err != nil {
		return nil
	}

	s, err := syncDirectory(srcDir, dstDir, sourceRoot, targetRoot, manifest, d.addOnly, isDryRun)
	if err != nil {
		return err
	}

	stats.Added += s.Added
	stats.Updated += s.Updated
	stats.Removed += s.Removed
	stats.Preserved += s.Preserved
	stats.Conflicts += s.Conflicts
	stats.InSync += s.InSync

	return nil
}

func syncDirectory(srcDir, dstDir, sourceRoot, targetRoot string, manifest *SyncManifest, addOnly bool, isDryRun bool) (MirrorStats, error) {
	var stats MirrorStats
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		return processWalkItem(path, info, srcDir, dstDir, sourceRoot, manifest, addOnly, isDryRun, &stats)
	})

	return stats, err
}

func processWalkItem(path string, info os.FileInfo, srcDir, dstDir, sourceRoot string, manifest *SyncManifest, addOnly bool, isDryRun bool, stats *MirrorStats) error {
	relToSrc, _ := filepath.Rel(srcDir, path)
	if relToSrc == "." || relToSrc == "" {
		return nil
	}

	relFromRoot, _ := filepath.Rel(sourceRoot, path)
	normRel := filepath.ToSlash(relFromRoot)

	if isExcludedSyncItem(normRel, info.Name()) {
		return handleExcludedItem(info)
	}

	dstPath := filepath.Join(dstDir, relToSrc)
	// WHY flat: the nested isDryRun check inside IsDir tripped the
	// nested-if linter on the v6.643.0 tip; both guards are depth-1 now.
	if info.IsDir() && !isDryRun {
		_ = os.MkdirAll(dstPath, 0755)
	}
	if info.IsDir() {
		return nil
	}

	return syncSingleFileThreeWay(relFromRoot, path, dstPath, manifest, addOnly, isDryRun, stats)
}

func isExcludedSyncItem(normRel, name string) bool {
	return IsSpec21OrHigher(normRel) || IsBumpScript(name) ||
		IsMemoryOrPlans(normRel) || IsArchive(normRel) || IsSecret(name)
}

func handleExcludedItem(info os.FileInfo) error {
	if info.IsDir() {
		return filepath.SkipDir
	}

	return nil
}

func syncSingleFileThreeWay(relFromRoot, srcPath, dstPath string, manifest *SyncManifest, addOnly bool, isDryRun bool, stats *MirrorStats) error {
	eval := EvaluateThreeWay(relFromRoot, srcPath, dstPath, manifest.Assets)
	if eval.State == "InSync" {
		stats.InSync++
		return nil
	}
	if eval.State == "NewUpstream" {
		return applyNewUpstream(eval, relFromRoot, srcPath, dstPath, manifest, isDryRun, stats)
	}
	if addOnly {
		return nil
	}
	if eval.State == "CleanUpdate" {
		return applyCleanUpdate(eval, relFromRoot, srcPath, dstPath, manifest, isDryRun, stats)
	}
	if eval.State == "LocalEdited" {
		stats.Preserved++
		return nil
	}
	if eval.State == "Collision" || eval.State == "UntrackedConflict" {
		return applyConflictSidecar(eval, srcPath, dstPath, isDryRun, stats)
	}
	if eval.State == "Converged" {
		stats.InSync++
		manifest.Assets[filepath.ToSlash(relFromRoot)] = eval.SourceHash
		return nil
	}

	return nil
}

func applyNewUpstream(eval ThreeWayResult, relFromRoot, src, dst string, m *SyncManifest, dry bool, stats *MirrorStats) error {
	if err := ApplyThreeWayDecision(eval, src, dst, dry); err != nil {
		return err
	}
	stats.Added++
	m.Assets[filepath.ToSlash(relFromRoot)] = eval.SourceHash

	return nil
}

func applyCleanUpdate(eval ThreeWayResult, relFromRoot, src, dst string, m *SyncManifest, dry bool, stats *MirrorStats) error {
	if err := ApplyThreeWayDecision(eval, src, dst, dry); err != nil {
		return err
	}
	stats.Updated++
	m.Assets[filepath.ToSlash(relFromRoot)] = eval.SourceHash

	return nil
}

func applyConflictSidecar(eval ThreeWayResult, src, dst string, dry bool, stats *MirrorStats) error {
	if err := ApplyThreeWayDecision(eval, src, dst, dry); err != nil {
		return err
	}
	stats.Conflicts++

	return nil
}

func appendSpecDirs(specRoot string, dirs []syncDirMapping) []syncDirMapping {
	entries, err := os.ReadDir(specRoot)
	if err != nil {
		return dirs
	}

	for _, e := range entries {
		dirs = appendSingleSpecDir(e, dirs)
	}

	return dirs
}

func appendSingleSpecDir(e os.DirEntry, dirs []syncDirMapping) []syncDirMapping {
	if !e.IsDir() {
		return dirs
	}

	m := specDirNumRe.FindStringSubmatch(e.Name())
	if len(m) < 2 {
		return dirs
	}

	n, _ := strconv.Atoi(m[1])
	if n < 1 || n > 20 {
		return dirs
	}

	rel := filepath.Join(constants.DirSpec, e.Name())

	return append(dirs, syncDirMapping{rel, rel, false})
}
