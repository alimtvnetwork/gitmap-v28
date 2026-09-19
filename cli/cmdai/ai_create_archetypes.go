package cmdai

import (
	"strings"
)

func buildLinterBody(opts CreateScriptOptions) string {
	workerPart := resolveWorkerDeclaration(opts.IsParallel)
	lines := []string{
		"def check_file(file_path: Path) -> list[str]:",
		"    issues: list[str] = []",
		"    # Inspect file content using 02-shared-engine helpers",
		"    return issues",
		"",
		"def main() -> int:",
		"    parser = argparse.ArgumentParser(description=__doc__)",
		"    parser.add_argument('--changed-only', action='store_true')",
		"    parser.add_argument('--json', action='store_true')",
		workerPart,
		"    args = parser.parse_args()",
		"    print('Scanning repository...')",
		"    return 0",
		"",
		"if __name__ == '__main__':",
		"    sys.exit(main())",
		"",
	}

	return filterEmptyLines(lines)
}

func buildFixerBody(opts CreateScriptOptions) string {
	lines := []string{
		"def fix_file(file_path: Path, is_apply: bool = False) -> bool:",
		"    # Apply automated fixes cleanly using write_file_lf",
		"    return True",
		"",
		"def main() -> int:",
		"    parser = argparse.ArgumentParser(description=__doc__)",
		"    parser.add_argument('--fix', action='store_true')",
		"    parser.add_argument('--dry-run', action='store_true')",
		"    args = parser.parse_args()",
		"    print(f'Running fixer (apply={args.fix})...')",
		"    return 0",
		"",
		"if __name__ == '__main__':",
		"    sys.exit(main())",
		"",
	}

	return strings.Join(lines, "\n")
}

func buildAuditorBody(opts CreateScriptOptions) string {
	lines := []string{
		"def audit_component() -> int:",
		"    # Perform cross-verification between specs and code",
		"    print('Auditing repository assets...')",
		"    return 0",
		"",
		"def main() -> int:",
		"    parser = argparse.ArgumentParser(description=__doc__)",
		"    parser.add_argument('--strict', action='store_true')",
		"    args = parser.parse_args()",
		"    return audit_component()",
		"",
		"if __name__ == '__main__':",
		"    sys.exit(main())",
		"",
	}

	return strings.Join(lines, "\n")
}

func buildGeneratorBody(opts CreateScriptOptions) string {
	lines := []string{
		"def generate_assets(output_dir: Path) -> int:",
		"    # Generate code, schema, or model assets",
		"    print(f'Generating assets into {output_dir}...')",
		"    return 0",
		"",
		"def main() -> int:",
		"    parser = argparse.ArgumentParser(description=__doc__)",
		"    parser.add_argument('--out', default='.', help='Output directory')",
		"    args = parser.parse_args()",
		"    return generate_assets(Path(args.out))",
		"",
		"if __name__ == '__main__':",
		"    sys.exit(main())",
		"",
	}

	return strings.Join(lines, "\n")
}

func resolveWorkerDeclaration(isParallel bool) string {
	if isParallel {
		return "    parser.add_argument('--workers', type=int, default=4)"
	}

	return ""
}

func filterEmptyLines(lines []string) string {
	var filtered []string
	for _, l := range lines {
		if l != "" || len(filtered) > 0 {
			filtered = append(filtered, l)
		}
	}

	return strings.Join(filtered, "\n")
}
