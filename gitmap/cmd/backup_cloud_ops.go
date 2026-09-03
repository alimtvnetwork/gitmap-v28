package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type SnapshotManifest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Note      string    `json:"note,omitempty"`
	Files     []string  `json:"files"`
}

func copyBackupArtifacts(destDir string) error {
	dataDir := store.BinaryDataDir()
	masterDB := filepath.Join(dataDir, "gitmap.db")
	if err := copyFileIfExists(masterDB, filepath.Join(destDir, "gitmap.db")); err != nil {
		return err
	}

	profFile := store.GitProfilesPath()
	if err := copyFileIfExists(profFile, filepath.Join(destDir, "git_profiles.json")); err != nil {
		return err
	}

	if err := copyFolderIfExists(filepath.Join(dataDir, "pipeline_db"), filepath.Join(destDir, "pipeline_db")); err != nil {
		return err
	}
	return copyFolderIfExists(filepath.Join(dataDir, "repo_search"), filepath.Join(destDir, "repo_search"))
}

func copyFileIfExists(src, dst string) error {
	data, err := os.ReadFile(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return apperror.WrapSimple(err, "read file:")
	}
	if writeErr := os.WriteFile(dst, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write file:")
	}
	return nil
}

func copyFolderIfExists(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return apperror.WrapSimple(err, "read dir:")
	}
	if mkErr := os.MkdirAll(dstDir, 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "create dir:")
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		s := filepath.Join(srcDir, e.Name())
		d := filepath.Join(dstDir, e.Name())
		if copyErr := copyFileIfExists(s, d); copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func writeSnapshotManifest(
	destDir string,
	snapID string,
	note string,
) error {
	m := SnapshotManifest{
		ID:        snapID,
		CreatedAt: time.Now(),
		Note:      note,
		Files:     []string{"gitmap.db", "git_profiles.json", "pipeline_db", "repo_search"},
	}
	data, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		return apperror.WrapSimple(err, "marshal manifest:")
	}

	if writeErr := os.WriteFile(filepath.Join(destDir, "manifest.json"), data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write manifest:")
	}

	return nil
}

func runBackupCloudList(args []string) error {
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	snapsDir := filepath.Join(cloudDir, "snapshots")
	entries, err := os.ReadDir(snapsDir)
	if err != nil || len(entries) == 0 {
		fmt.Println("\n  No cloud backup snapshots found.")
		return nil
	}
	if hasArgFlag(args, "--json") {
		return outputCloudSnapshotsJSON(entries)
	}
	printCloudSnapshotsTable(entries, snapsDir)
	return nil
}

func outputCloudSnapshotsJSON(entries []os.DirEntry) error {
	list := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			list = append(list, e.Name())
		}
	}
	data, _ := json.MarshalIndent(map[string]any{"snapshots": list}, "", "  ")
	fmt.Println(string(data))
	return nil
}

func printCloudSnapshotsTable(entries []os.DirEntry, snapsDir string) {
	fmt.Printf("\n  %s● Cloud Backup Snapshots%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  --------------------------------------------------------------------------------")
	fmt.Println("  #    Snapshot ID                Date                 Note")
	fmt.Println("  --------------------------------------------------------------------------------")
	idx := 1
	for _, e := range entries {
		if e.IsDir() {
			note := readSnapshotNote(filepath.Join(snapsDir, e.Name()))
			fmt.Printf("  [%d]  %-26s %-20s %s\n", idx, e.Name(), formatDirDate(e.Name()), note)
			idx++
		}
	}
	fmt.Println("  --------------------------------------------------------------------------------")
	fmt.Println()
}

func readSnapshotNote(snapPath string) string {
	data, err := os.ReadFile(filepath.Join(snapPath, "manifest.json"))
	if err != nil {
		return "-"
	}
	var m SnapshotManifest
	if json.Unmarshal(data, &m) == nil && m.Note != "" {
		return m.Note
	}
	return "-"
}

func formatDirDate(dirName string) string {
	parts := strings.Split(dirName, "-")
	if len(parts) >= 4 {
		return parts[1] + "-" + parts[2] + "-" + parts[3]
	}
	return "-"
}

func runBackupCloudRestore(args []string) error {
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	snapsDir := filepath.Join(cloudDir, "snapshots")
	target, resolveErr := resolveSnapshotTarget(snapsDir, args)
	if resolveErr != nil {
		return resolveErr
	}
	if !confirmOrSkip("Restore snapshot '"+target+"'? This replaces local databases.", args) {
		fmt.Println("  Restore aborted.")
		return nil
	}
	src := filepath.Join(snapsDir, target)
	dst := store.BinaryDataDir()
	if restoreErr := copyBackupArtifactsToDir(src, dst); restoreErr != nil {
		return restoreErr
	}
	fmt.Printf("  %s✓ Restored databases and profiles from: %s%s\n",
		constants.ColorGreen, target, constants.ColorReset)
	return nil
}

func resolveSnapshotTarget(snapsDir string, args []string) (string, error) {
	entries, err := os.ReadDir(snapsDir)
	if err != nil || len(entries) == 0 {
		return "", apperror.NewSimple("no snapshots available to restore", "E1080")
	}
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return pickSnapshotByNameOrIndex(entries, args[0])
	}
	if !isInteractiveStdin() {
		return entries[len(entries)-1].Name(), nil // default to latest
	}
	return promptSnapshotSelection(entries)
}

func pickSnapshotByNameOrIndex(entries []os.DirEntry, val string) (string, error) {
	num, err := strconv.Atoi(val)
	if err == nil && num >= 1 && num <= len(entries) {
		return entries[num-1].Name(), nil
	}
	for _, e := range entries {
		if e.Name() == val {
			return e.Name(), nil
		}
	}
	return "", apperror.NewSimple("snapshot not found: "+val, "E1081")
}

func promptSnapshotSelection(entries []os.DirEntry) (string, error) {
	fmt.Printf("\nSelect snapshot to restore (1-%d):\n", len(entries))
	for i, e := range entries {
		fmt.Printf("  [%d] %s\n", i+1, e.Name())
	}
	fmt.Print("Choice: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", apperror.WrapSimple(err, "read choice:")
	}
	return pickSnapshotByNameOrIndex(entries, strings.TrimSpace(line))
}

func copyBackupArtifactsToDir(srcDir, dstDir string) error {
	masterDB := filepath.Join(srcDir, "gitmap.db")
	if err := copyFileIfExists(masterDB, filepath.Join(dstDir, "gitmap.db")); err != nil {
		return err
	}

	profFile := filepath.Join(srcDir, "git_profiles.json")
	if err := copyFileIfExists(profFile, filepath.Join(dstDir, "git_profiles.json")); err != nil {
		return err
	}

	if err := copyFolderIfExists(filepath.Join(srcDir, "pipeline_db"), filepath.Join(dstDir, "pipeline_db")); err != nil {
		return err
	}
	return copyFolderIfExists(filepath.Join(srcDir, "repo_search"), filepath.Join(dstDir, "repo_search"))
}

func runBackupCloudRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("usage: gitmap backup rm <snapshot-id|1-N>", "E1082")
	}
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	snapsDir := filepath.Join(cloudDir, "snapshots")
	target, resolveErr := resolveSnapshotTarget(snapsDir, args)
	if resolveErr != nil {
		return resolveErr
	}
	if !confirmOrSkip("Delete cloud backup snapshot '"+target+"'?", args) {
		fmt.Println("  Aborted.")
		return nil
	}
	if rmErr := os.RemoveAll(filepath.Join(snapsDir, target)); rmErr != nil {
		return apperror.WrapSimple(rmErr, "remove snapshot:")
	}
	cmdCommit := exec.Command("git", "commit", "-am", "backup: remove "+target)
	cmdCommit.Dir = cloudDir
	_ = cmdCommit.Run()
	cmdPush := exec.Command("git", "push", "origin", "main")
	cmdPush.Dir = cloudDir
	_ = cmdPush.Run()
	fmt.Printf("  %s✓ Removed snapshot: %s%s\n", constants.ColorGreen, target, constants.ColorReset)
	return nil
}

func runBackupCloudStatus(args []string) error {
	prof, _ := resolveDefaultCloudProfile(args)
	slug := resolveBackupRepoSlug(prof, args)
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	snapsDir := filepath.Join(cloudDir, "snapshots")
	entries, _ := os.ReadDir(snapsDir)

	fmt.Printf("\n  %s● Cloud Backup Status%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  ● Remote Repository: https://github.com/%s\n", slug)
	fmt.Printf("  ● Active Profile:    %s (%s)\n", prof.Name, prof.Provider)
	fmt.Printf("  ● Total Snapshots:   %d\n", len(entries))
	fmt.Printf("  ● Local Cache Path:  %s\n\n", cloudDir)
	return nil
}
