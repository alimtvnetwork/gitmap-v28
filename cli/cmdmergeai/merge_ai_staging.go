// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdresolver"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

func cleanDestStagingDirectory(destDir string) error {
	entries, err := os.ReadDir(destDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		_ = os.RemoveAll(filepath.Join(destDir, entry.Name()))
	}
	return nil
}

func truncateCommitSha(sha string) string {
	if len(sha) > 8 {
		return sha[:8]
	}
	return sha
}

func stageAndSequenceFiles(destDir string, sources []SourceRepoConfig) ([]FileCollisionRecord, int, error) {
	collisionMap := make(map[string]*FileCollisionRecord)
	placedOrigin := make(map[string]FileCollisionVariant)
	uniquePlaced := 0

	for seqIdx, src := range sources {
		seqPrefix := fmt.Sprintf("%02d", seqIdx+1)
		sourceSlug := filepath.Base(src.LocalStaging)
		headSha := src.CommitMetadata.ShortSHA
		if headSha == "" {
			fullSha, _ := gitExec(src.LocalStaging, "rev-parse", "HEAD")
			headSha = truncateCommitSha(fullSha)
		}

		for _, relPath := range src.FolderTree {
			errStage := stageSingleFile(destDir, sources, src, relPath, seqPrefix, sourceSlug, headSha, collisionMap, placedOrigin, &uniquePlaced)
			if errStage != nil {
				return nil, 0, errStage
			}
		}
	}

	var collisionList []FileCollisionRecord
	for _, c := range collisionMap {
		collisionList = append(collisionList, *c)
	}

	sort.Slice(collisionList, func(i, j int) bool {
		return collisionList[i].CanonicalPath < collisionList[j].CanonicalPath
	})

	return collisionList, uniquePlaced, nil
}

func stageSingleFile(
	destDir string,
	sources []SourceRepoConfig,
	src SourceRepoConfig,
	relPath, seqPrefix, sourceSlug, headSha string,
	collisionMap map[string]*FileCollisionRecord,
	placedOrigin map[string]FileCollisionVariant,
	uniquePlaced *int,
) error {
	destPath := filepath.Join(destDir, relPath)

	if col, isKnownCollision := collisionMap[relPath]; isKnownCollision {
		return handleKnownCollision(col, destDir, src.LocalStaging, relPath, seqPrefix, sourceSlug, headSha)
	}

	if !fsutil.FileExists(destPath) {
		return handleUniqueFile(destPath, src.LocalStaging, relPath, seqPrefix, sourceSlug, headSha, placedOrigin, uniquePlaced)
	}

	return handleFirstCollision(destDir, sources, src, relPath, destPath, seqPrefix, sourceSlug, headSha, collisionMap, placedOrigin, uniquePlaced)
}

func handleKnownCollision(
	col *FileCollisionRecord,
	destDir, stagingDir, relPath, seqPrefix, sourceSlug, headSha string,
) error {
	dirName := filepath.Dir(relPath)
	baseName := filepath.Base(relPath)
	newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
	newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
	errCopy := copyFile(filepath.Join(stagingDir, relPath), newSeqDest)
	if errCopy != nil {
		return errCopy
	}

	col.Variants = append(col.Variants, FileCollisionVariant{
		Sequence:   seqPrefix,
		File:       newSeqRel,
		SourceRepo: sourceSlug,
		Commit:     headSha,
	})
	return nil
}

func handleUniqueFile(
	destPath, stagingDir, relPath, seqPrefix, sourceSlug, headSha string,
	placedOrigin map[string]FileCollisionVariant,
	uniquePlaced *int,
) error {
	errCopy := copyFile(filepath.Join(stagingDir, relPath), destPath)
	if errCopy != nil {
		return errCopy
	}
	placedOrigin[relPath] = FileCollisionVariant{
		Sequence:   seqPrefix,
		File:       relPath,
		SourceRepo: sourceSlug,
		Commit:     headSha,
	}
	*uniquePlaced++
	return nil
}

func resolveOriginalVariant(sources []SourceRepoConfig, relPath string, placedOrigin map[string]FileCollisionVariant) FileCollisionVariant {
	orig, hasOrig := placedOrigin[relPath]
	if hasOrig {
		return orig
	}
	baseHeadSha, _ := gitExec(sources[0].LocalStaging, "rev-parse", "HEAD")
	return FileCollisionVariant{
		Sequence:   "01",
		File:       relPath,
		SourceRepo: filepath.Base(sources[0].LocalStaging),
		Commit:     truncateCommitSha(baseHeadSha),
	}
}

func handleFirstCollision(
	destDir string,
	sources []SourceRepoConfig,
	src SourceRepoConfig,
	relPath, destPath, seqPrefix, sourceSlug, headSha string,
	collisionMap map[string]*FileCollisionRecord,
	placedOrigin map[string]FileCollisionVariant,
	uniquePlaced *int,
) error {
	orig := resolveOriginalVariant(sources, relPath, placedOrigin)
	dirName := filepath.Dir(relPath)
	baseName := filepath.Base(relPath)

	origSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", orig.Sequence, baseName)))
	origSeqPath := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", orig.Sequence, baseName))
	_ = os.Rename(destPath, origSeqPath)
	if *uniquePlaced > 0 {
		*uniquePlaced--
	}
	delete(placedOrigin, relPath)

	newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
	newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
	errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), newSeqDest)
	if errCopy != nil {
		return errCopy
	}

	collisionMap[relPath] = &FileCollisionRecord{
		CanonicalPath: relPath,
		Status:        "pending_ai_consolidation",
		Variants: []FileCollisionVariant{
			{
				Sequence:   orig.Sequence,
				File:       origSeqRel,
				SourceRepo: orig.SourceRepo,
				Commit:     orig.Commit,
			},
			{
				Sequence:   seqPrefix,
				File:       newSeqRel,
				SourceRepo: sourceSlug,
				Commit:     headSha,
			},
		},
	}
	return nil
}

func copyFile(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func assembleManifest(dest *cmdresolver.ResolvedRepoFeature, sources []SourceRepoConfig, collisions []FileCollisionRecord, uniqueCount int, skipped []SkippedSourceItem) MergeAIManifest {
	var m MergeAIManifest
	m.Attributes.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	m.Attributes.GitMapVersion = constants.Version
	m.Attributes.Tool = "gitmap-merge-ai"
	m.Attributes.CollisionStrategy = "prefixed-sequence"

	m.Data.Destination.TargetType = string(dest.TargetType)
	m.Data.Destination.Slug = dest.CanonicalSlug
	m.Data.Destination.ResolvedPath = dest.LocalPath
	m.Data.Destination.OriginURL = dest.RemoteURL

	for i, s := range sources {
		var entry SourceRepoSequenceEntry
		entry.SequenceOrder = i + 1
		entry.Title = s.Title
		entry.RepoURL = s.RepoURL
		entry.Branch = s.Branch
		entry.CommitRange = s.CommitRange
		entry.LastCommitDate = s.LastCommitDate.Format(time.RFC3339)
		entry.Description = s.Description
		entry.CommitMetadata = s.CommitMetadata
		entry.ReleaseInfo.LastTag = s.ReleaseTag
		entry.ReleaseInfo.TagHash = s.ReleaseTagHash
		entry.FolderTree = s.FolderTree
		m.Data.RepoSequence = append(m.Data.RepoSequence, entry)
	}

	m.Data.SkippedSources = skipped
	m.Data.FileCollisions = collisions
	m.Data.UniqueFilesDirectlyPlaced = uniqueCount
	return m
}
