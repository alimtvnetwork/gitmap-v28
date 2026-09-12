package cmdchromeprofile

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func discoverZipCandidates(zipPath string) ([]DiscoveredProfileCandidate, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip %s: %w", zipPath, err)
	}

	defer r.Close()

	m := readZipManifest(r)
	if m != nil && len(m.Profiles) > 0 {
		return buildZipCandidatesFromManifest(zipPath, m, r), nil
	}

	return discoverZipCandidatesByFolder(zipPath, r), nil
}

func buildZipCandidatesFromManifest(zipPath string, m *chromeProfileManifest, r *zip.ReadCloser) []DiscoveredProfileCandidate {
	var candidates []DiscoveredProfileCandidate
	for _, p := range m.Profiles {
		exp := &chromeExport{Name: p.Name, DisplayName: p.DisplayName, Email: p.Email}
		target := resolveImportDestination(exp, "", false)
		bmCount := countZipProfileBookmarks(r, p.Name)
		candidates = append(candidates, DiscoveredProfileCandidate{
			SourcePath:        zipPath,
			ProfileDirName:    p.Name,
			DisplayName:       p.DisplayName,
			Email:             p.Email,
			BookmarksCount:    bmCount,
			ExtensionsCount:   p.ExtensionCount,
			IsFromZip:         true,
			ZipArchive:        filepath.Base(zipPath),
			TargetDestination: target,
		})
	}

	return candidates
}

func countZipProfileBookmarks(r *zip.ReadCloser, profName string) int {
	prefix := profName + "/Bookmarks"
	for _, f := range r.File {
		if f.Name != prefix {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return 0
		}

		defer rc.Close()
		raw, _ := io.ReadAll(rc)

		return countBookmarks(raw)
	}

	return 0
}

func discoverZipCandidatesByFolder(zipPath string, r *zip.ReadCloser) []DiscoveredProfileCandidate {
	seen := make(map[string]bool)
	var candidates []DiscoveredProfileCandidate
	for _, f := range r.File {
		parts := strings.Split(filepath.ToSlash(f.Name), "/")
		if len(parts) >= 2 && !seen[parts[0]] && !strings.HasPrefix(parts[0], "__") {
			seen[parts[0]] = true
			exp := &chromeExport{Name: parts[0], DisplayName: parts[0]}
			target := resolveImportDestination(exp, "", false)
			candidates = append(candidates, DiscoveredProfileCandidate{
				SourcePath:        zipPath,
				ProfileDirName:    parts[0],
				DisplayName:       parts[0],
				IsFromZip:         true,
				ZipArchive:        filepath.Base(zipPath),
				TargetDestination: target,
			})
		}
	}

	return candidates
}

func discoverSingleFileCandidate(filePath string) ([]DiscoveredProfileCandidate, error) {
	meta, err := readSnapshotMetadata(filePath)
	if err != nil {
		return nil, err
	}

	exp := meta.Export
	cand := DiscoveredProfileCandidate{
		SourcePath:        filePath,
		ProfileDirName:    meta.ProfileName,
		DisplayName:       meta.DisplayName,
		Email:             meta.Email,
		BookmarksCount:    meta.BookmarksCount,
		ExtensionsCount:   meta.ExtensionsCount,
		TargetDestination: meta.TargetDestination,
	}

	if exp != nil && cand.ProfileDirName == "" {
		cand.ProfileDirName = exp.Name
	}

	return []DiscoveredProfileCandidate{cand}, nil
}

func importZipSnapshotWithStepLogging(zipPath, explicitTarget string) error {
	fmt.Printf("  \033[1;94m[Step 1/5]\033[0m Inspecting archive: %s\n", zipPath)
	candidates, err := discoverZipCandidates(zipPath)
	if err != nil {
		return err
	}

	fmt.Printf("        → Discovered %d profile(s) in ZIP archive\n", len(candidates))
	fmt.Printf("  \033[1;94m[Step 2/5]\033[0m Unpacking profiles into Chrome User Data...\n")

	if err := applyChromeExportZIPWithOptions(zipPath, explicitTarget, 0); err != nil {
		return err
	}

	for _, c := range candidates {
		dstDir := c.TargetDestination.Dir
		if explicitTarget != "" {
			dstDir = explicitTarget
		}

		_ = registerImportedProfileToLocalState(dstDir, c.DisplayName, c.Email)
	}

	fmt.Printf("  \033[1;94m[Step 5/5]\033[0m Reconciling Chrome Local State...\n")

	return nil
}

func importDirectorySnapshotWithStepLogging(dirPath, explicitTarget string) error {
	base := filepath.Base(dirPath)
	fmt.Printf("  \033[1;94m[Step 1/5]\033[0m Inspecting profile directory: %s\n", dirPath)
	jsonPath := filepath.Join(dirPath, base+".json")
	if chromeProfilePathExists(jsonPath) {
		return importSingleSnapshotWithStepLogging(jsonPath, explicitTarget)
	}

	target := base
	if explicitTarget != "" {
		target = explicitTarget
	}

	destPath := chromeProfilePath(target)
	if err := os.MkdirAll(destPath, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", destPath, err)
	}

	copyProfileDirectoryDiskFiles(dirPath, destPath)
	disp, email := resolveProfileNameAndEmail(target, nil)
	_ = patchImportedChromeProfilePreferences(destPath, disp)
	_ = registerImportedProfileToLocalState(target, disp, email)
	fmt.Printf("  \033[1;92m✓\033[0m %s → imported from directory\n", target)

	return nil
}

func copySingleProfileFile(src, dst string) {
	if !chromeProfilePathExists(src) {
		return
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return
	}

	_ = os.MkdirAll(filepath.Dir(dst), constants.DirPermission)
	_ = os.WriteFile(dst, data, constants.FilePermission)
}

func copyProfileDirectoryDiskFiles(srcDir, destDir string) {
	for _, name := range []string{"Bookmarks", "Preferences"} {
		copySingleProfileFile(filepath.Join(srcDir, name), filepath.Join(destDir, name))
	}

	for _, dbName := range constants.ChromeProfileSQLiteEntries {
		copySingleProfileFile(filepath.Join(srcDir, dbName), filepath.Join(destDir, dbName))
	}

	copySingleProfileFile(filepath.Join(srcDir, "Cookies"), filepath.Join(destDir, "Cookies"))
	copySingleProfileFile(filepath.Join(srcDir, "Network", "Cookies"), filepath.Join(destDir, "Network", "Cookies"))
}
