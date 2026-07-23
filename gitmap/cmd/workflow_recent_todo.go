// Package cmd — workflow_recent_todo.go: `gitmap recent` jumps back to
// recent repos via the navigation helper history; `gitmap todo` greps
// TODO/FIXME/XXX with blame.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runRecent(args []string) {
	checkHelp("recent", args)

	path := recentLogPath()
	f, err := os.Open(path) //nolint:gosec
	if err != nil {
		fmt.Fprintf(os.Stderr, "recent: no history yet (%s)\n", path)
		os.Exit(0)
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if l != "" {
			lines = append(lines, l)
		}
	}
	uniq := dedupeReverse(lines, 10)
	if len(args) > 0 && args[0] == "--print" {
		for _, p := range uniq {
			fmt.Println(p)
		}
		return
	}
	fmt.Println("\033[1;94mRecent repos\033[0m")
	for i, p := range uniq {
		fmt.Printf("  \033[2;37m%2d\033[0m %s\n", i+1, p)
	}
	fmt.Println("\n\033[2;37mhint:\033[0m pipe to fzf — \033[1;96mgitmap recent --print | fzf\033[0m")
}

func recentLogPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gitmap", "recent.log")
}

func dedupeReverse(in []string, max int) []string {
	seen := map[string]bool{}
	out := []string{}
	for i := len(in) - 1; i >= 0 && len(out) < max; i-- {
		if !seen[in[i]] {
			seen[in[i]] = true
			out = append(out, in[i])
		}
	}
	return out
}

func runTodo(args []string) {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	out, err := exec.Command("git", "-C", root, "grep", "-nE", `TODO|FIXME|XXX`).Output()
	if err != nil && len(out) == 0 {
		fmt.Println("todo: no matches")
		return
	}
	fmt.Printf("\033[1;94mTODO / FIXME / XXX\033[0m in %s\n", root)
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		line := sc.Text()
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		file, lno := parts[0], parts[1]
		blame, _ := exec.Command("git", "-C", root, "blame", "-L", lno+","+lno, "--porcelain", file).Output()
		author := "?"
		for _, bl := range strings.Split(string(blame), "\n") {
			if strings.HasPrefix(bl, "author ") {
				author = strings.TrimPrefix(bl, "author ")
				break
			}
		}
		fmt.Printf("  \033[2;37m%s:%s\033[0m  \033[1;93m%s\033[0m  %s\n",
			file, lno, author, strings.TrimSpace(parts[2]))
	}
}
