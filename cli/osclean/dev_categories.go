package osclean

import (
	"os"
	"path/filepath"
	"runtime"
)

func cleanGoCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "go-buildcache",
		Label:    "Go build cache + module downloads (~/go/bin SAFE)",
	}
	invokeGoClean(isDryRun, &res)
	for _, dir := range resolveGoPaths() {
		sub := SweepTarget(dir, isDryRun, true)
		mergeCleanStats(&res, sub)
	}

	return res
}

func invokeGoClean(isDryRun bool, res *CategoryCleanStats) {
	if isDryRun || !hasTool("go") {
		return
	}
	_, _ = runToolCommand("go", "clean", "-cache", "-testcache", "-fuzzcache")
	_, _ = runToolCommand("go", "clean", "-modcache")
	res.Notes = append(res.Notes, "Invoked 'go clean -cache -modcache -testcache -fuzzcache'")
}

func resolveGoPaths() []string {
	var paths []string
	if out, ok := runToolCommand("go", "env", "GOCACHE"); ok && len(out) > 0 {
		paths = append(paths, out)
	}
	if out, ok := runToolCommand("go", "env", "GOMODCACHE"); ok && len(out) > 0 {
		paths = append(paths, out)
	}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths = append(paths, filepath.Join(local, "go-build"))
	paths = append(paths, filepath.Join(home, "go", "pkg", "mod"))
	paths = append(paths, filepath.Join(home, ".cache", "go-build"))
	paths = append(paths, resolveDevDrives("go/pkg/mod")...)

	return filterUniquePaths(paths)
}

func cleanPnpmCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "pnpm-store",
		Label:    "pnpm CAS store + download cache (runtime SAFE)",
	}
	if !isDryRun && hasTool("pnpm") {
		_, _ = runToolCommand("pnpm", "store", "prune")
		res.Notes = append(res.Notes, "Invoked 'pnpm store prune'")
	}
	for _, dir := range resolvePnpmPaths() {
		sub := SweepTarget(dir, isDryRun, false)
		mergeCleanStats(&res, sub)
	}

	return res
}

func resolvePnpmPaths() []string {
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(home, ".pnpm-store"),
		filepath.Join(local, "pnpm", "store"),
		filepath.Join(local, "pnpm-cache"),
		filepath.Join(home, ".local", "share", "pnpm", "store"),
	}
	paths = append(paths, resolveDevDrives("pnpm/store")...)

	return filterUniquePaths(paths)
}

func cleanNpmCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "npm-cache",
		Label:    "npm cache + tarball store",
	}
	if !isDryRun && hasTool("npm") {
		_, _ = runToolCommand("npm", "cache", "clean", "--force")
		res.Notes = append(res.Notes, "Invoked 'npm cache clean --force'")
	}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "npm-cache"),
		filepath.Join(home, ".npm"),
		filepath.Join(home, "AppData", "Roaming", "npm-cache"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func cleanChocoCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "choco-cache",
		Label:    "Chocolatey package cache + installer downloads",
	}
	if runtime.GOOS != "windows" {
		return res
	}
	invokeChocoClean(isDryRun, &res)
	for _, dir := range resolveChocoPaths() {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func invokeChocoClean(isDryRun bool, res *CategoryCleanStats) {
	if !isDryRun && hasTool("choco") {
		_, _ = runToolCommand("choco", "cache", "clean", "-y")
		res.Notes = append(res.Notes, "Invoked 'choco cache clean -y'")
	}
}

func resolveChocoPaths() []string {
	local := resolveLocalAppData()
	progData := os.Getenv("ProgramData")
	tempDir := os.Getenv("TEMP")
	return filterUniquePaths([]string{
		filepath.Join(local, "Chocolatey", "Cache"),
		filepath.Join(local, "Temp", "chocolatey"),
		filepath.Join(progData, "chocolatey", "cache"),
		filepath.Join(tempDir, "chocolatey"),
	})
}

func mergeCleanStats(target *CategoryCleanStats, source CategoryCleanStats) {
	target.ItemsRemoved += source.ItemsRemoved
	target.DirsRemoved += source.DirsRemoved
	target.BytesFreed += source.BytesFreed
	target.Errors = append(target.Errors, source.Errors...)
}
