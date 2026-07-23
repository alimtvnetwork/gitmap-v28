package release

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// discoverMetadataPending finds .gitmap/release/v*.json files where neither
// the Git branch nor the tag exists. Skips versions already in pendingBranches.
func discoverMetadataPending(pendingBranches []string) []ReleaseMeta {
	metaFiles, err := ListReleaseMetaFiles()
	if err != nil {
		return nil
	}

	branchSet := make(map[string]bool, len(pendingBranches))
	for _, b := range pendingBranches {
		branchSet[b] = true
	}

	var pending []ReleaseMeta

	for _, meta := range metaFiles {
		if isMetaPending(meta, branchSet) {
			pending = append(pending, meta)
		}
	}

	return pending
}

// isMetaPending returns true when the metadata version has no branch/tag
// and is not already covered by a pending branch.
func isMetaPending(meta ReleaseMeta, branchSet map[string]bool) bool {
	if len(meta.Commit) == 0 {
		return false
	}

	branchName := constants.ReleaseBranchPrefix + meta.Tag
	if branchSet[branchName] {
		return false
	}
	if BranchExists(branchName) {
		return false
	}
	if TagExistsLocally(meta.Tag) || TagExistsRemote(meta.Tag) {
		return false
	}

	return true
}

// releasePendingFromMetadata creates branch+tag from stored commit SHA.
func releasePendingFromMetadata(pending []ReleaseMeta, assetsPath, notes string, isDraft, dryRun bool) error {
	for _, meta := range pending {
		err := releaseFromMetadata(meta, assetsPath, notes, isDraft, dryRun)
		if err != nil {
			fmt.Printf(constants.MsgReleasePendingFailed, meta.Tag, err)
			continue
		}
	}

	return nil
}

// releaseFromMetadata creates a release branch+tag from a metadata file's commit SHA.
func releaseFromMetadata(meta ReleaseMeta, assetsPath, notes string, isDraft, dryRun bool) error {
	v, err := Parse(meta.Tag)
	if err != nil {
		return fmt.Errorf("invalid version in metadata: %s", meta.Tag)
	}

	if !CommitExists(meta.Commit) {
		fmt.Printf(constants.WarnPendingMetaNoCommit, meta.Tag, meta.Commit)

		return nil
	}

	shortSHA := meta.Commit
	if len(shortSHA) > 7 {
		shortSHA = shortSHA[:7]
	}

	fmt.Printf(constants.MsgPendingMetaRelease, meta.Tag, shortSHA)

	if dryRun {
		branchName := constants.ReleaseBranchPrefix + v.String()
		fmt.Printf(constants.MsgReleaseDryRun, "Create branch "+branchName+" from commit "+shortSHA)
		fmt.Printf(constants.MsgReleaseDryRun, "Create tag "+v.String())
		fmt.Printf(constants.MsgReleaseDryRun, "Push branch and tag to origin")

		return nil
	}

	branchName := constants.ReleaseBranchPrefix + v.String()

	err = CreateBranch(branchName, meta.Commit)
	if err != nil {
		return fmt.Errorf("create branch from metadata: %w", err)
	}
	fmt.Printf(constants.MsgReleaseBranch, branchName)

	tag := v.String()
	err = CreateTag(tag, constants.ReleaseTagPrefix+tag)
	if err != nil {
		return fmt.Errorf("create tag from metadata: %w", err)
	}
	fmt.Printf(constants.MsgReleaseTag, tag)

	opts := Options{Assets: assetsPath, Notes: notes, IsDraft: isDraft, SkipMeta: true}

	return pushAndFinalize(v, branchName, tag, "metadata:"+meta.Commit, opts)
}

// isPendingBranch returns true when the branch has no released tag.
func isPendingBranch(branch string) bool {
	v, err := extractVersionFromBranch(branch)
	if err != nil {
		return false
	}

	tag := v.String()

	return tagIsMissing(tag)
}

// tagIsMissing returns true when a tag does not exist locally or remotely.
func tagIsMissing(tag string) bool {
	if TagExistsLocally(tag) {
		return false
	}
	if TagExistsRemote(tag) {
		return false
	}

	return true
}
