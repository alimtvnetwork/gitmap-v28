#!/usr/bin/env python3
"""Cross-platform check for changelog version synchronization."""
import os
import re
import sys


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    constants_file = os.path.join(repo_root, "gitmap", "constants", "constants.go")
    changelog_file = os.path.join(repo_root, "changelog.md")

    if not os.path.isfile(constants_file):
        print(f"✗ constants.go not found at {constants_file}", file=sys.stderr)
        sys.exit(1)
    if not os.path.isfile(changelog_file):
        print(f"✗ changelog.md not found at {changelog_file}", file=sys.stderr)
        sys.exit(1)

    version = None
    with open(constants_file, "r", encoding="utf-8", errors="replace") as fh:
        for line in fh:
            m = re.search(r'^(?:var|const)\s+Version\s*=\s*"([0-9]+\.[0-9]+\.[0-9]+)"', line)
            if m:
                version = m.group(1)
                break

    if not version:
        print(f"✗ Could not parse Version from {constants_file}", file=sys.stderr)
        print('  Expected line: var Version = "X.Y.Z"', file=sys.stderr)
        sys.exit(1)

    print(f"→ constants.Version = {version}")

    with open(changelog_file, "r", encoding="utf-8", errors="replace") as fh:
        changelog_text = fh.read()

    escaped_ver = re.escape(version)
    heading_re = re.compile(rf"^##\s+\[?v?{escaped_ver}\]?(?:\s|$)", re.MULTILINE)

    if not heading_re.search(changelog_text):
        print("", file=sys.stderr)
        print(f"✗ CHANGELOG drift: constants.Version is {version} but no matching", file=sys.stderr)
        print(f"  '## v{version}' heading exists in changelog.md.", file=sys.stderr)
        print("", file=sys.stderr)
        print(f"  Fix: add a '## v{version}' section to {changelog_file} describing", file=sys.stderr)
        print("  the release, then re-run this check.", file=sys.stderr)
        sys.exit(1)

    print(f"✓ changelog.md has entry for v{version}")
    sys.exit(0)


if __name__ == "__main__":
    main()
