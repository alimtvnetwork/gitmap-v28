package cmdchromeprofile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func readSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	lower := strings.ToLower(srcFile)
	if strings.HasSuffix(lower, constants.ExtZIP) {
		return readZipSnapshotMetadata(srcFile)
	}

	if isDirectoryPath(srcFile) {
		return readDirectorySnapshotMetadata(srcFile)
	}

	if strings.HasSuffix(lower, constants.ExtJSON) {
		return readJSONSnapshotMetadata(srcFile)
	}

	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return readYAMLSnapshotMetadata(srcFile)
	}

	return readGenericSnapshotMetadata(srcFile)
}

func readZipSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	candidates, err := discoverZipCandidates(srcFile)
	if err != nil || len(candidates) == 0 {
		return readGenericSnapshotMetadata(srcFile)
	}

	base := filepath.Base(srcFile)
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}

	if len(candidates) == 1 {
		c := candidates[0]
		exp := &chromeExport{Name: c.ProfileDirName, DisplayName: c.DisplayName, Email: c.Email}

		return &snapshotMetadata{
			FilePath:          srcFile,
			FileName:          base,
			FileSize:          size,
			Export:            exp,
			BookmarksCount:    c.BookmarksCount,
			ExtensionsCount:   c.ExtensionsCount,
			HasEmail:          c.Email != "",
			Email:             c.Email,
			DisplayName:       c.DisplayName,
			ProfileName:       c.ProfileDirName,
			TargetDestination: c.TargetDestination,
		}, nil
	}

	totalBM := 0
	totalExt := 0
	for _, c := range candidates {
		totalBM += c.BookmarksCount
		totalExt += c.ExtensionsCount
	}

	profLabel := fmt.Sprintf("%d profiles", len(candidates))
	exp := &chromeExport{Name: profLabel, DisplayName: profLabel}
	target := resolveImportDestination(exp, "", false)

	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          base,
		FileSize:          size,
		Export:            exp,
		BookmarksCount:    totalBM,
		ExtensionsCount:   totalExt,
		DisplayName:       profLabel,
		ProfileName:       profLabel,
		TargetDestination: target,
	}, nil
}

func readDirectorySnapshotMetadata(srcDir string) (*snapshotMetadata, error) {
	base := filepath.Base(srcDir)
	cand, ok := readSubdirProfileCandidate(filepath.Dir(srcDir), base)
	if !ok {
		return readGenericSnapshotMetadata(srcDir)
	}

	info, _ := os.Stat(srcDir)
	var size int64
	if info != nil {
		size = info.Size()
	}

	exp := &chromeExport{Name: cand.ProfileDirName, DisplayName: cand.DisplayName, Email: cand.Email}

	return &snapshotMetadata{
		FilePath:          srcDir,
		FileName:          base,
		FileSize:          size,
		Export:            exp,
		BookmarksCount:    cand.BookmarksCount,
		ExtensionsCount:   cand.ExtensionsCount,
		HasEmail:          cand.Email != "",
		Email:             cand.Email,
		DisplayName:       cand.DisplayName,
		ProfileName:       cand.ProfileDirName,
		TargetDestination: cand.TargetDestination,
	}, nil
}

func readJSONSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	raw, err := os.ReadFile(srcFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", srcFile, err)
	}

	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", srcFile, err)
	}

	if exp.Name == "" {
		base := filepath.Base(srcFile)
		exp.Name = strings.TrimSuffix(base, filepath.Ext(base))
	}

	extractEmailIfMissing(&exp, raw)
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}

	bmsCount := countBookmarks(exp.Bookmarks)
	extsCount := len(exp.ExtensionIDs)
	target := resolveImportDestination(&exp, "", false)

	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          filepath.Base(srcFile),
		FileSize:          size,
		Export:            &exp,
		BookmarksCount:    bmsCount,
		ExtensionsCount:   extsCount,
		HasEmail:          exp.Email != "",
		Email:             exp.Email,
		DisplayName:       exp.DisplayName,
		ProfileName:       exp.Name,
		ExportedAt:        exp.ExportedAt,
		TargetDestination: target,
	}, nil
}
