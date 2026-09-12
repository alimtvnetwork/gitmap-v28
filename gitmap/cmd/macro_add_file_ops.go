// Package cmd — macro_add_file_ops.go: in-builder file manipulation helpers (cat, touch, mkfile, rmfile, cpfile).
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func runCatCmd(args []string) error {
	if len(args) == 0 {
		printCatUsage()

		return apperror.NewValidationError("missing filepath for cat")
	}

	target := macro.ExpandPathAndEnv(args[0])
	data, err := os.ReadFile(target)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read file %q", target))
	}

	printCatContent(target, data)

	return nil
}

func printCatUsage() {
	fmt.Println("Usage: gitmap cat <filepath>")
	fmt.Println()
	fmt.Println("Displays file contents cleanly in the terminal across all operating systems.")
	fmt.Println()
}

func printCatContent(target string, data []byte) {
	fmt.Printf("\n  %s--- %s (%d bytes) ---%s\n",
		constants.ColorCyan, target, len(data), constants.ColorReset)
	fmt.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		fmt.Println()
	}

	fmt.Printf("  %s--------------------------------%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func runTouchCmd(args []string) error {
	if len(args) == 0 {
		printTouchUsage()

		return apperror.NewValidationError("missing filepath for touch")
	}

	target := macro.ExpandPathAndEnv(args[0])
	if err := touchTargetFile(target); err != nil {
		return err
	}

	fmt.Printf("  %s✓ Touched file: %s%s\n\n", constants.ColorGreen, target, constants.ColorReset)

	return nil
}

func printTouchUsage() {
	fmt.Println("Usage: gitmap touch <filepath>")
	fmt.Println()
	fmt.Println("Creates an empty file or updates modtime, automatically creating parent directories.")
	fmt.Println()
}

func touchTargetFile(target string) error {
	ensureParentDir(target)
	now := time.Now()
	if err := os.Chtimes(target, now, now); err == nil {
		return nil
	}

	return createEmptyTouchFile(target)
}

func ensureParentDir(target string) {
	dir := filepath.Dir(target)
	if dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
}

func createEmptyTouchFile(target string) error {
	file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("touch file %q", target))
	}

	defer file.Close()

	return nil
}

func runMkfileCmd(args []string) error {
	if len(args) == 0 {
		printMkfileUsage()

		return apperror.NewValidationError("missing filepath for mkfile")
	}

	target := macro.ExpandPathAndEnv(args[0])
	content := resolveMkfileContent(args)

	return writeAndAnnounceNewFile(target, content)
}

func resolveMkfileContent(args []string) string {
	if len(args) > 1 {
		return strings.Join(args[1:], " ")
	}

	return ""
}

func writeAndAnnounceNewFile(target, content string) error {
	if err := writeNewFile(target, content); err != nil {
		return err
	}

	fmt.Printf("  %s✓ Created file: %s (%d bytes)%s\n\n",
		constants.ColorGreen, target, len(content), constants.ColorReset)

	return nil
}

func printMkfileUsage() {
	fmt.Println("Usage: gitmap mkfile <filepath> [content...]")
	fmt.Println()
	fmt.Println("Creates a file with optional initial content, automatically creating parent directories.")
	fmt.Println()
}

func writeNewFile(target, content string) error {
	ensureParentDir(target)
	if err := os.WriteFile(target, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("create file %q", target))
	}

	return nil
}

func runRmFileCmd(args []string) error {
	if len(args) == 0 {
		fmt.Printf("  %s▲ Usage: rmfile <filepath>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return apperror.NewValidationError("missing filepath for rmfile")
	}

	target := macro.ExpandPathAndEnv(args[0])
	if err := os.Remove(target); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("remove file %q", target))
	}

	fmt.Printf("  %s✓ Removed file: %s%s\n\n", constants.ColorGreen, target, constants.ColorReset)

	return nil
}

func runCpFileCmd(args []string) error {
	if len(args) < 2 {
		fmt.Printf("  %s▲ Usage: cpfile <src> <dst>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return apperror.NewValidationError("missing src or dst for cpfile")
	}

	src := macro.ExpandPathAndEnv(args[0])
	dst := macro.ExpandPathAndEnv(args[1])

	return copyFileDirect(src, dst)
}

func copyFileDirect(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read source file %q", src))
	}

	ensureParentDir(dst)
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("write destination file %q", dst))
	}

	fmt.Printf("  %s✓ Copied %s → %s (%d bytes)%s\n\n",
		constants.ColorGreen, src, dst, len(data), constants.ColorReset)

	return nil
}
