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

func stageAndSequenceFiles(destDir string, sources []SourceRepoConfig) ([]FileCollisionRecord, int, error) {
	collisionMap := make(map[string]*FileCollisionRecord)
	placedOrigin := make(map[string]FileCollisionVariant)
	uniquePlaced := 0

	for seqIdx, src := range sources {
		seqPrefix := fmt.Sprintf("%02d", seqIdx+1)
		sourceSlug := filepath.Base(src.LocalStaging)
		headSha, _ := gitExec(src.LocalStaging, "rev-parse", "HEAD")
		if len(headSha) > 8 {
			headSha = headSha[:8]
		}

		for _, relPath := range src.FolderTree {
			destPath := filepath.Join(destDir, relPath)
			dirName := filepath.Dir(relPath)
			baseName := filepath.Base(relPath)

			// Case 1: Already known collision on this relPath
			if col, isKnownCollision := collisionMap[relPath]; isKnownCollision {
				newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
				newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), newSeqDest); errCopy != nil {
					return nil, 0, errCopy
				}

				col.Variants = append(col.Variants, FileCollisionVariant{
					Sequence:   seqPrefix,
					File:       newSeqRel,
					SourceRepo: sourceSlug,
					Commit:     headSha,
				})
				continue
			}

			// Case 2: File does not exist yet in destination
			if !fileExists(destPath) {
				if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), destPath); errCopy != nil {
					return nil, 0, errCopy
				}
				placedOrigin[relPath] = FileCollisionVariant{
					Sequence:   seqPrefix,
					File:       relPath,
					SourceRepo: sourceSlug,
					Commit:     headSha,
				}
				uniquePlaced++
				continue
			}

			// Case 3: First collision detected on relPath:
			orig, hasOrig := placedOrigin[relPath]
			if !hasOrig {
				baseHeadSha, _ := gitExec(sources[0].LocalStaging, "rev-parse", "HEAD")
				if len(baseHeadSha) > 8 {
					baseHeadSha = baseHeadSha[:8]
				}
				orig = FileCollisionVariant{
					Sequence:   "01",
					File:       relPath,
					SourceRepo: filepath.Base(sources[0].LocalStaging),
					Commit:     baseHeadSha,
				}
			}

			origSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", orig.Sequence, baseName)))
			origSeqPath := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", orig.Sequence, baseName))
			_ = os.Rename(destPath, origSeqPath)
			if uniquePlaced > 0 {
				uniquePlaced--
			}
			delete(placedOrigin, relPath)

			col := &FileCollisionRecord{
				CanonicalPath: relPath,
				Status:        "pending_ai_consolidation",
				Variants: []FileCollisionVariant{
					{
						Sequence:   orig.Sequence,
						File:       origSeqRel,
						SourceRepo: orig.SourceRepo,
						Commit:     orig.Commit,
					},
				},
			}
			collisionMap[relPath] = col

			// Now place incoming colliding version as seqPrefix_<filename>
			newSeqRel := filepath.ToSlash(filepath.Join(dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName)))
			newSeqDest := filepath.Join(destDir, dirName, fmt.Sprintf("%s_%s", seqPrefix, baseName))
			if errCopy := copyFile(filepath.Join(src.LocalStaging, relPath), newSeqDest); errCopy != nil {
				return nil, 0, errCopy
			}

			col.Variants = append(col.Variants, FileCollisionVariant{
				Sequence:   seqPrefix,
				File:       newSeqRel,
				SourceRepo: sourceSlug,
				Commit:     headSha,
			})
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
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
		entry.ReleaseInfo.LastTag = s.ReleaseTag
		entry.FolderTree = s.FolderTree
		m.Data.RepoSequence = append(m.Data.RepoSequence, entry)
	}

	m.Data.SkippedSources = skipped
	m.Data.FileCollisions = collisions
	m.Data.UniqueFilesDirectlyPlaced = uniqueCount
	return m
}
