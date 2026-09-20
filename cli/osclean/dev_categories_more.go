package osclean

import (
	"path/filepath"
)

func cleanYarnCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "yarn-cache",
		Label:    "Yarn package cache",
	}
	if !isDryRun && hasTool("yarn") {
		_, _ = runToolCommand("yarn", "cache", "clean")
		res.Notes = append(res.Notes, "Invoked 'yarn cache clean'")
	}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "Yarn", "Cache"),
		filepath.Join(home, ".yarn", "cache"),
		filepath.Join(home, ".cache", "yarn"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func cleanBunCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "bun-cache",
		Label:    "Bun package cache",
	}
	if !isDryRun && hasTool("bun") {
		_, _ = runToolCommand("bun", "pm", "cache", "rm")
		res.Notes = append(res.Notes, "Invoked 'bun pm cache rm'")
	}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "bun", "install", "cache"),
		filepath.Join(home, ".bun", "install", "cache"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func cleanPipCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "pip-cache",
		Label:    "Python pip download cache",
	}
	if !isDryRun && (hasTool("pip") || hasTool("python")) {
		invokePipClean(&res)
	}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "pip", "cache"),
		filepath.Join(home, ".cache", "pip"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
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
	home := resolveHomeDir()
	cargo := filepath.Join(home, ".cargo")
	paths := []string{
		filepath.Join(cargo, "registry", "cache"),
		filepath.Join(cargo, "registry", "src"),
		filepath.Join(cargo, "git", "db"),
		filepath.Join(cargo, "git", "checkouts"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func cleanNugetCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "nuget-cache",
		Label:    ".NET / NuGet package and HTTP cache",
	}
	if !isDryRun && hasTool("dotnet") {
		_, _ = runToolCommand("dotnet", "nuget", "locals", "all", "--clear")
		res.Notes = append(res.Notes, "Invoked 'dotnet nuget locals all --clear'")
	}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "NuGet", "v3-cache"),
		filepath.Join(local, "NuGet", "plugins-cache"),
		filepath.Join(home, ".nuget", "packages"),
		filepath.Join(home, ".local", "share", "NuGet", "v3-cache"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}

func cleanGradleMavenCache(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{
		Category: "gradle-maven-cache",
		Label:    "Gradle and Maven build caches",
	}
	home := resolveHomeDir()
	paths := []string{
		filepath.Join(home, ".gradle", "caches"),
		filepath.Join(home, ".m2", "repository"),
	}
	for _, dir := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(dir, isDryRun, false))
	}

	return res
}
