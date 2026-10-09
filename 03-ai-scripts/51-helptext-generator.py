#!/usr/bin/env python3
"""
51-helptext-generator.py - DRY help topic generator (spec 243.3).

Renders cli/helpdoc/*.md topics from the HelpDisplay structs (the single
source of truth) so hand-editing those topics is never needed again. The Go
emitter library at cli/tool/helptextemitter owns the topic registry; this
script only drives it via `go run` on an ephemeral runner file, so no stray
package mains live in the repository.

Usage:
  gitmap py 03-ai-scripts/51-helptext-generator.py [--topics space,scan] [--check]
  gitmap py 03-ai-scripts/51-helptext-generator.py --topics space --out /tmp/preview

  --topics   Comma-separated topic subset (default: all registered topics).
  --check    Render to a temp dir and exit 1 when any committed topic drifts.
             For CI / pre-commit enforcement of the generated files.
  --out      Destination dir for <topic>.md (default: <repo>/cli/helpdoc).
"""

import argparse
import filecmp
import os
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

# Repository root discovery (script lives in <repo>/03-ai-scripts/).
REPO_ROOT = Path(__file__).resolve().parent.parent
MODULE_DIR = REPO_ROOT / "cli"
HELPTEXT_DIR = MODULE_DIR / "helpdoc"

RUNNER_SOURCE = """package main

import (
\t"flag"
\t"fmt"
\t"os"

\t"github.com/alimtvnetwork/gitmap-v28/cli/tool/helptextemitter"
)

func main() {
\toutDir := flag.String("out", "", "destination directory for <topic>.md files")
\ttopics := flag.String("topics", "", "comma-separated topic subset (default: all)")
\tflag.Parse()
\tif err := helptextemitter.Generate(*outDir, *topics); err != nil {
\t\tfmt.Fprintln(os.Stderr, "helptext-generator: "+err.Error())
\t\tos.Exit(1)
\t}
}
"""


def find_go():
    """Locates the Go toolchain: PATH first, then ~/go/bin/go."""
    found = shutil.which("go")
    if found:
        return found
    home_go = os.path.expanduser("~/go/bin/go")
    if os.path.isfile(home_go) and os.access(home_go, os.X_OK):
        return home_go
    raise RuntimeError("go toolchain not found on PATH or at ~/go/bin/go")


def run_emitter(out_dir, topics):
    """Renders topics via `go run` on an ephemeral runner; raises on failure."""
    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".go", prefix="helptext-gen-", delete=False
    ) as runner:
        runner.write(RUNNER_SOURCE)
        runner_path = runner.name
    try:
        result = subprocess.run(
            [find_go(), "run", runner_path, "--out", str(out_dir), "--topics", topics],
            cwd=str(MODULE_DIR),
            capture_output=True,
            text=True,
        )
    finally:
        os.unlink(runner_path)
    if result.returncode != 0:
        detail = result.stderr.strip() or "no stderr captured"
        raise RuntimeError("go emitter failed: " + detail)
    return result


def parse_args(argv):
    """Parses CLI flags."""
    parser = argparse.ArgumentParser(
        description="Generate cli/helpdoc/*.md from HelpDisplay structs."
    )
    parser.add_argument(
        "--topics",
        default="",
        help="Comma-separated topic subset (default: all registered topics).",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Fail (exit 1) when any committed topic drifts from the structs.",
    )
    parser.add_argument(
        "--out",
        default=str(HELPTEXT_DIR),
        help="Destination dir for <topic>.md files.",
    )
    return parser.parse_args(argv)


def check_topics(topics):
    """Renders to a temp dir and reports drift vs committed files."""
    with tempfile.TemporaryDirectory(prefix="helptext-check-") as tmp:
        run_emitter(tmp, topics)
        drifted = find_drifted(tmp)
    report_drift(drifted)
    return 1 if drifted else 0


def find_drifted(rendered_dir):
    """Returns topic names whose rendered output differs from the committed file."""
    drifted = []
    for rendered in sorted(Path(rendered_dir).glob("*.md")):
        committed = HELPTEXT_DIR / rendered.name
        if not committed.exists():
            drifted.append(rendered.stem + " (missing)")
        elif not filecmp.cmp(str(rendered), str(committed), shallow=False):
            drifted.append(rendered.stem)
    return drifted


def report_drift(drifted):
    """Prints the drift report."""
    if not drifted:
        print("helptext topics are up to date.")
        return
    print("DRIFTED helptext topics (regenerate with 51-helptext-generator.py):")
    for name in drifted:
        print("  - " + name)


def main(argv):
    """Entry point."""
    args = parse_args(argv)
    if args.check:
        return check_topics(args.topics)
    run_emitter(args.out, args.topics)
    print("Wrote helptext topics to " + args.out)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
