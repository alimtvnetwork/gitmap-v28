// Package cmd — macro_add_file_ops.go: in-builder file manipulation helpers (cat, touch, mkfile, rmfile, cpfile).
package cmdmacro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func runCatCmd(args []string) error {
	opts, err := ParseCatOptions(args)
	if err != nil {
		return err
	}

	return ExecuteCatWithOptions(opts)
}

func printCatUsage() {
	fmt.Println("Usage: gitmap cat <filepath> [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -L, --lines <start:end>   Slice lines from start to end (e.g. -L 140:220)")
	fmt.Println("  -n, --head <count>        Show first N lines (e.g. --head 30, -n 25)")
	fmt.Println("  --tail <count>            Show last N lines (e.g. --tail 30)")
	fmt.Println("  -num, --number            Display 1-based line numbers in gutter")
	fmt.Println("  -q, --query <text>        Filter lines matching query string")
	fmt.Println("  -l, --limit <count>       Cap maximum matching lines (default: 50)")
	fmt.Println("  --raw                     Output raw bytes without borders or banners")
	fmt.Println()
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

	return removeTargetFile(target)
}

func removeTargetFile(target string) error {
	if err := os.Remove(target); err != nil {
		if os.IsNotExist(err) {
			return apperror.NewNotFound("rmfile", "E4004", fmt.Sprintf("file %q not found", target))
		}

		return apperror.WrapSimple(err, fmt.Sprintf("remove file %q", target))
	}

	fmt.Printf("  %s✓ Removed file: %s%s\n\n", constants.ColorGreen, target, constants.ColorReset)

	return nil
}
