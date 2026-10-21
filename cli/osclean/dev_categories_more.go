package osclean

import (
	"os"
	"path/filepath"
)

func cleanYarnCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "yarn-cache", Label: "Yarn package cache"}
	home, local := resolveHomeDir(), resolveLocalAppData()
	paths := []string{filepath.Join(local, "Yarn", "Cache"), filepath.Join(home, ".yarn", "cache"), filepath.Join(home, ".cache", "yarn")}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}
	if !isDryRun && hasTool("yarn") {
		_, _ = runToolCommand("yarn", "cache", "clean")
		res.Notes = append(res.Notes, "Invoked 'yarn cache clean'")
	}
	return res
}

func cleanBunCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "bun-cache", Label: "Bun package cache"}
	home, local := resolveHomeDir(), resolveLocalAppData()
	paths := []string{filepath.Join(local, "bun", "install", "cache"), filepath.Join(home, ".bun", "install", "cache")}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}
	if !isDryRun && hasTool("bun") {
		_, _ = runToolCommand("bun", "pm", "cache", "rm")
		res.Notes = append(res.Notes, "Invoked 'bun pm cache rm'")
	}
	return res
}

func cleanPipCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "pip-cache", Label: "Python pip download cache"}
	home, local := resolveHomeDir(), resolveLocalAppData()
	paths := []string{filepath.Join(local, "pip", "cache"), filepath.Join(home, ".cache", "pip")}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}
	if !isDryRun && (hasTool("pip") || hasTool("python")) {
		invokePipClean(&res)
	}
	return res
}

func invokePipClean(res *CategoryCleanStats) {
	if _, ok := runToolCommand("pip", "cache", "purge"); ok {
		res.Notes = append(res.Notes, "Invoked 'pip cache purge'")
		return
	}
	if _, ok := runToolCommand("python", "-m", "pip", "cache", "purge"); ok {
		res.Notes = append(res.Notes, "Invoked 'python -m pip cache purge'")
	}
}

func cleanCargoCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "cargo-cache",
		Label:    "Cargo/Rust registry cache and git checkouts",
	}
	for _, dir := range resolveCargoPaths() {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, true))
	}
	invokeCargoClean(isDryRun, &res)

	return res
}

func invokeCargoClean(isDryRun bool, res *CategoryCleanStats) {
	if isDryRun {
		return
	}
	invokeCargoTools(res)
	invokeCargoTargetClean(res)
}

func invokeCargoTools(res *CategoryCleanStats) {
	if hasTool("cargo-cache") {
		_, _ = runToolCommand("cargo", "cache", "--autoclean")
		res.Notes = append(res.Notes, "Invoked 'cargo cache --autoclean'")
	}
	if hasTool("sccache") {
		_, _ = runToolCommand("sccache", "--clear-cache")
		res.Notes = append(res.Notes, "Invoked 'sccache --clear-cache'")
	}
}

func invokeCargoTargetClean(res *CategoryCleanStats) {
	tgt := os.Getenv("CARGO_TARGET_DIR")
	if len(tgt) == 0 || !hasTool("cargo") {
		return
	}
	_, _ = runToolCommand("cargo", "clean", "--target-dir", tgt)
	res.Notes = append(res.Notes, "Invoked 'cargo clean --target-dir "+tgt+"'")
}

func resolveCargoPaths() []string {
	var paths []string
	home, local := resolveHomeDir(), resolveLocalAppData()
	paths = append(paths, resolveCargoRegistryPaths(home)...)
	paths = append(paths, resolveRustToolingPaths(home, local)...)
	paths = append(paths, resolveCargoDevDrivePaths()...)

	return filterUniquePaths(paths)
}

func resolveCargoDevDrivePaths() []string {
	var paths []string
	subs := []string{
		"cargo/registry/cache", "cargo/registry/src", "cargo/git/db",
		"cargo/git/checkouts", "cargo/cache", "cargo/target", "cargo/build",
		"rust/cache", "rust/target", "rust/build",
	}
	for _, sub := range subs {
		paths = append(paths, resolveDevDrives(sub)...)
	}

	return paths
}

func resolveCargoRegistryPaths(home string) []string {
	cargoHome := os.Getenv("CARGO_HOME")
	if len(cargoHome) == 0 {
		cargoHome = filepath.Join(home, ".cargo")
	}

	return []string{
		filepath.Join(cargoHome, "registry", "cache"),
		filepath.Join(cargoHome, "registry", "src"),
		filepath.Join(cargoHome, "git", "db"),
		filepath.Join(cargoHome, "git", "checkouts"),
		filepath.Join(cargoHome, ".package-cache"),
	}
}

func resolveRustToolingPaths(home, local string) []string {
	paths := resolveRustupDirs(home)
	paths = append(paths, resolveSccacheDirs(home, local)...)
	if tgt := os.Getenv("CARGO_TARGET_DIR"); len(tgt) > 0 {
		paths = append(paths, tgt)
	}

	return paths
}

func resolveRustupDirs(home string) []string {
	rustupHome := os.Getenv("RUSTUP_HOME")
	if len(rustupHome) == 0 {
		rustupHome = filepath.Join(home, ".rustup")
	}

	return []string{filepath.Join(rustupHome, "downloads"), filepath.Join(rustupHome, "tmp")}
}

func resolveSccacheDirs(home, local string) []string {
	var paths []string
	if scc := os.Getenv("SCCACHE_DIR"); len(scc) > 0 {
		paths = append(paths, scc)
	}

	return append(paths,
		filepath.Join(local, "Mozilla", "sccache"),
		filepath.Join(home, ".cache", "sccache"),
		filepath.Join(home, "Library", "Caches", "Mozilla.sccache"),
	)
}

func cleanCcacheCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "ccache-cache",
		Label:    "C/C++ compiler cache (ccache)",
	}
	for _, dir := range resolveCcacheDirs() {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}
	invokeCcacheClean(isDryRun, &res)

	return res
}

func resolveCcacheDirs() []string {
	home, local := resolveHomeDir(), resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "ccache"),
		filepath.Join(home, ".cache", "ccache"),
		filepath.Join(home, ".ccache"),
	}
	if cdir := os.Getenv("CCACHE_DIR"); len(cdir) > 0 {
		paths = append(paths, cdir)
	}

	return filterUniquePaths(paths)
}

func invokeCcacheClean(isDryRun bool, res *CategoryCleanStats) {
	if !isDryRun && hasTool("ccache") {
		_, _ = runToolCommand("ccache", "-C")
		res.Notes = append(res.Notes, "Invoked 'ccache -C'")
	}
}

func cleanNugetCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "nuget-cache",
		Label:    ".NET / NuGet package and HTTP cache",
	}
	invokeNugetClean(isDryRun, &res)
	for _, dir := range resolveNugetPaths() {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func invokeNugetClean(isDryRun bool, res *CategoryCleanStats) {
	if !isDryRun && hasTool("dotnet") {
		_, _ = runToolCommand("dotnet", "nuget", "locals", "all", "--clear")
		res.Notes = append(res.Notes, "Invoked 'dotnet nuget locals all --clear'")
	}
}

func resolveNugetPaths() []string {
	home, local := resolveHomeDir(), resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "NuGet", "v3-cache"),
		filepath.Join(local, "NuGet", "plugins-cache"),
		filepath.Join(home, ".nuget", "packages"),
		filepath.Join(home, ".local", "share", "NuGet", "v3-cache"),
	}

	return filterUniquePaths(paths)
}

func cleanGradleMavenCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "gradle-maven-cache",
		Label:    "Gradle and Maven build caches",
	}
	for _, dir := range resolveGradleMavenPaths() {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func resolveGradleMavenPaths() []string {
	home := resolveHomeDir()
	paths := []string{
		filepath.Join(home, ".gradle", "caches"),
		filepath.Join(home, ".m2", "repository"),
	}

	return filterUniquePaths(paths)
}
