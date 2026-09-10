package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runPurge(args []string) error {
	isRestore := false
	autoConfirm := false
	var pattern string

	for _, arg := range args {
		if arg == "--restore" {
			isRestore = true
		} else if arg == "-y" || arg == "--confirm" {
			autoConfirm = true
		} else if !strings.HasPrefix(arg, "-") && pattern == "" {
			pattern = arg
		}
	}

	repoPath, err := os.Getwd()
	if err != nil {
		return apperror.Wrap(err, "failed to get current directory", nil)
	}
	db, err := store.OpenDefault()
	if err != nil {
		return apperror.Wrap(err, "failed to open database", nil)
	}

	if isRestore {
		return doRestore(db, repoPath)
	}

	if pattern == "" {
		return apperror.NewSimple("EXECUTION", "pattern argument is required for purging (e.g., '*.db')")
	}

	return doPurge(db, repoPath, pattern, autoConfirm)
}

type shFileOpStruct struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

func sendToRecycleBin(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shFileOperation := shell32.NewProc("SHFileOperationW")

	pFrom, err := syscall.UTF16PtrFromString(absPath + "\x00")
	if err != nil {
		return err
	}

	op := shFileOpStruct{
		wFunc:  3, // FO_DELETE
		pFrom:  pFrom,
		fFlags: 0x40 | 0x10 | 0x0400 | 0x0004, // FOF_ALLOWUNDO | FOF_NOCONFIRMATION | FOF_NOERRORUI | FOF_SILENT
	}

	ret, _, _ := shFileOperation.Call(uintptr(unsafe.Pointer(&op)))
	if ret != 0 {
		return fmt.Errorf("SHFileOperation failed with code %d", ret)
	}
	return nil
}

func runPurgeCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("command %s failed: %s\n%s", name, err, string(out))
	}
	return string(out), nil
}

func doPurge(db *store.DB, repoPath, pattern string, autoConfirm bool) error {
	out, err := runPurgeCmd("git", "status", "--porcelain")
	if err != nil {
		return apperror.Wrap(err, "git status failed", nil)
	}
	if strings.TrimSpace(out) != "" {
		return apperror.NewSimple("EXECUTION", "Working tree is not clean. Please commit or stash changes before running.")
	}

	out, err = runPurgeCmd("git", "branch", "--show-current")
	if err != nil || strings.TrimSpace(out) == "" {
		return apperror.NewSimple("EXECUTION", "Could not determine current branch. Are you in a detached HEAD?")
	}

	out, err = runPurgeCmd("git", "ls-files", pattern)
	if err != nil {
		return apperror.Wrap(err, "git ls-files failed", nil)
	}
	var matchingFiles []string
	for _, f := range strings.Split(out, "\n") {
		f = strings.TrimSpace(f)
		if f != "" {
			matchingFiles = append(matchingFiles, f)
		}
	}

	remoteURL := ""
	out, err = runPurgeCmd("git", "remote", "get-url", "origin")
	if err == nil {
		remoteURL = strings.TrimSpace(out)
	}

	timestamp := time.Now().Unix()
	backupBranch := fmt.Sprintf("backup-purge-%d", timestamp)
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("gitmap_purge_%d", timestamp))

	fmt.Printf("\n--- GIT HISTORY PURGE WARNING ---\n")
	fmt.Printf("Target pattern: %s\n", pattern)
	fmt.Printf("Matching files in working tree: %d\n", len(matchingFiles))
	for i, f := range matchingFiles {
		if i < 5 {
			fmt.Printf("  - %s\n", f)
		}
	}
	if len(matchingFiles) > 5 {
		fmt.Printf("  ... and %d more.\n", len(matchingFiles)-5)
	}

	fmt.Println("\nThis operation will:")
	fmt.Println("1. Copy the current matching files to your Windows Temp directory.")
	fmt.Println("2. Send the original matching files in your repo to the Recycle Bin.")
	fmt.Println("3. Create a backup branch of your current git history.")
	fmt.Println("4. Rewrite your ENTIRE Git history using git filter-repo to remove the files.")
	fmt.Println("5. Add the pattern to .gitignore and commit it.")
	fmt.Printf("\nBackup Temp Directory will be: %s\n", tempDir)
	fmt.Printf("Backup Git Branch will be: %s\n", backupBranch)
	fmt.Println("---------------------------------")

	if !autoConfirm {
		fmt.Print("\nPlease type 'I confirm' to proceed: ")
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		if strings.TrimSpace(ans) != "I confirm" {
			fmt.Println("Operation aborted.")
			os.Exit(0)
		}
	}

	fmt.Println("\nStarting purge process...")

	_, err = runPurgeCmd("git", "branch", backupBranch)
	if err != nil {
		return apperror.Wrap(err, "failed to create backup branch", nil)
	}
	fmt.Printf("Created backup branch: %s\n", backupBranch)

	os.MkdirAll(tempDir, 0755)
	var backedUpFiles []string

	for _, f := range matchingFiles {
		if _, err := os.Stat(f); err == nil {
			dest := filepath.Join(tempDir, f)
			os.MkdirAll(filepath.Dir(dest), 0755)
			copyPurgeFile(f, dest)
			backedUpFiles = append(backedUpFiles, f)
			sendToRecycleBin(f)
		}
	}

	if len(backedUpFiles) > 0 {
		fmt.Printf("Backed up and recycled %d files.\n", len(backedUpFiles))
	}

	fmt.Println("Rewriting history with git filter-repo (this may take a moment)...")
	out, err = runPurgeCmd("git", "filter-repo", "--path-glob", pattern, "--invert-paths", "--force")
	if err != nil {
		fmt.Printf("git filter-repo output: %s\n", out)
		return apperror.Wrap(err, "git filter-repo failed", nil)
	}

	if remoteURL != "" {
		runPurgeCmd("git", "remote", "add", "origin", remoteURL)
	}

	f, _ := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(fmt.Sprintf("\n%s\n", pattern))
	f.Close()

	runPurgeCmd("git", "add", ".gitignore")
	runPurgeCmd("git", "commit", "-m", fmt.Sprintf("chore: add %s to .gitignore", pattern))

	filesJSON, _ := json.Marshal(backedUpFiles)
	log := store.PurgeHistoryLog{
		RepoPath:     repoPath,
		Pattern:      pattern,
		BackupBranch: backupBranch,
		TempDir:      tempDir,
		Files:        string(filesJSON),
		Timestamp:    time.Now().Unix(),
		Restored:     false,
	}
	if err := db.InsertPurgeHistoryLog(&log); err != nil {
		return apperror.Wrap(err, "failed to save purge state to database", nil)
	}

	fmt.Println("\n✅ Purge completed successfully!")
	fmt.Println("\nIf everything looks good, force push your changes:")
	fmt.Println("  git push origin --force --all")
	fmt.Println("  git push origin --force --tags")
	fmt.Println("\nIf you need to revert this operation, run:")
	fmt.Println("  gitmap purge --restore")

	return nil
}

func doRestore(db *store.DB, repoPath string) error {
	log, err := db.GetLastPurgeHistoryLog(repoPath)
	if err != nil || log == nil {
		return apperror.NewSimple("EXECUTION", "No active purge state found to restore.")
	}

	fmt.Printf("Restoring history from backup branch: %s...\n", log.BackupBranch)
	_, err = runPurgeCmd("git", "reset", "--hard", log.BackupBranch)
	if err != nil {
		return apperror.Wrap(err, "failed to reset to backup branch", nil)
	}

	if log.TempDir != "" {
		if _, err := os.Stat(log.TempDir); err == nil {
			fmt.Printf("Restoring files from temp directory: %s...\n", log.TempDir)
			filepath.WalkDir(log.TempDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				rel, _ := filepath.Rel(log.TempDir, path)
				os.MkdirAll(filepath.Dir(rel), 0755)
				copyPurgeFile(path, rel)
				return nil
			})
		}
	}

	db.MarkPurgeHistoryRestored(log.ID)

	fmt.Println("\n✅ Restore completed successfully!")
	fmt.Printf("You can now safely delete the backup branch if desired: git branch -D %s\n", log.BackupBranch)

	return nil
}

func copyPurgeFile(src, dst string) error {
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
