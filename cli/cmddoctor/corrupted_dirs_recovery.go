package cmddoctor

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

func resolveCanonicalInstallDir() string {
	if runtime.GOOS != "windows" && os.Geteuid() == 0 {
		return "/usr/local/bin"
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "/tmp"
	}

	return filepath.Join(home, ".local", "bin")
}

func copyFileWithPerm(src, dst string, perm fs.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, in)

	return err
}

func isRecoverableAsset(name string) bool {
	base := filepath.Base(name)
	if base == "gitmap" || base == "gitmap.exe" || base == "powershell.json" {
		return true
	}

	if base == ".gitmap-last-profile" || base == ".gitmap-last-install-dir" {
		return true
	}

	return false
}

func assetFilePerm(name string) fs.FileMode {
	if name == "gitmap" || name == "gitmap.exe" {
		return 0755
	}

	return 0644
}

func recoverSingleAsset(path, targetDir, name string) (string, error) {
	dst := filepath.Join(targetDir, name)
	if _, statErr := os.Stat(dst); !os.IsNotExist(statErr) {
		return "", nil
	}

	perm := assetFilePerm(name)
	if err := copyFileWithPerm(path, dst, perm); err != nil {
		return "", err
	}

	return dst, nil
}

// RecoverCorruptedDirAssets copies trapped Gitmap binaries and configs to canonical dir.
func RecoverCorruptedDirAssets(info CorruptedDirInfo, targetDir string) ([]string, error) {
	if !info.HasFiles {
		return nil, nil
	}

	var recovered []string
	err := filepath.WalkDir(info.Path, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || !isRecoverableAsset(d.Name()) {
			return walkErr
		}

		if dst, err := recoverSingleAsset(path, targetDir, d.Name()); err == nil && dst != "" {
			recovered = append(recovered, dst)
		}

		return nil
	})

	return recovered, err
}
