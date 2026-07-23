#!/usr/bin/env python3
"""Audit gitmap/constants/constants_*.go for identifier collisions.

Fails (exit 1) when any of the following is true:

  [1] CROSS-FILE EXACT-NAME COLLISIONS
      Same identifier declared in two files of the `constants` package.
      Would break `go build` with "redeclared in this block".

  [2] CROSS-FILE CASE-INSENSITIVE COLLISIONS
      Different exact names that lowercase to the same string and live in
      different files (e.g. `HelpFoo` vs `helpFoo`). Latent confusion risk
      and a sign of inconsistent naming, even though Go accepts it.

  [3] INTRA-FILE DUPLICATE DECLARATIONS
      Same identifier declared twice in one file. `go build` catches this,
      but the script reports it for completeness so reviewers see the
      offending lines without waiting for a Go compile.

The parser is string-literal aware: it tracks raw-string (`...`) and
"..." quoted regions so that SQL keywords like FROM, WHERE, VALUES that
appear INSIDE multi-line raw-string SQL constants are NEVER mistaken for
top-level identifiers.

Run locally:
    python3 .github/scripts/check-constants-collisions.py

CI invocation lives in .github/workflows/ci.yml.
"""

from __future__ import annotations

import glob
import os
import re
import sys
from collections import defaultdict

CONSTANTS_GLOB = "gitmap/constants/constants*.go"
IDENT_TOP_RE = re.compile(r"^([A-Z][A-Za-z0-9_]*(?:\s*,\s*[A-Z][A-Za-z0-9_]*)*)\b")


def parse_file(path: str):
    """Yield (kind, name, lineno) for each top-level const/var declaration.

    Tracks raw-string and quoted-string state to avoid matching tokens
    inside string literals.
    """
    out: list[tuple[str, str, int]] = []
    in_raw = False
    in_block: str | None = None
    with open(path, encoding="utf-8") as fh:
        for lineno, raw in enumerate(fh, 1):
            line = raw.rstrip("\n")
            if in_raw:
                idx = line.find("`")
                if idx == -1:
                    continue
                in_raw = False
                line = line[idx + 1 :]
            decl_part = line
            for q in ("`", '"'):
                idx = decl_part.find(q)
                if idx != -1:
                    decl_part = decl_part[:idx]
            if line.count("`") % 2 == 1:
                in_raw = True
            stripped = decl_part.strip()
            if in_block is None:
                if re.match(r"^\s*const\s*\(\s*(?://.*)?$", decl_part):
                    in_block = "const"
                    continue
                if re.match(r"^\s*var\s*\(\s*(?://.*)?$", decl_part):
                    in_block = "var"
                    continue
                m = re.match(r"^\s*const\s+([A-Z][A-Za-z0-9_]*)\b", decl_part)
                if m:
                    out.append(("const", m.group(1), lineno))
                    continue
                m = re.match(r"^\s*var\s+([A-Z][A-Za-z0-9_]*)\b", decl_part)
                if m:
                    out.append(("var", m.group(1), lineno))
                continue
            if re.match(r"^\s*\)\s*(?://.*)?$", decl_part):
                in_block = None
                continue
            m = IDENT_TOP_RE.match(stripped)
            if not m:
                continue
            for name in (n.strip() for n in m.group(1).split(",")):
                out.append((in_block, name, lineno))
    return out


def main() -> int:
    repo_root = os.getcwd()
    files = sorted(glob.glob(os.path.join(repo_root, CONSTANTS_GLOB)))
    files = [f for f in files if not f.endswith("_test.go")]
    if not files:
        print(f"::error::No files matched {CONSTANTS_GLOB}", file=sys.stderr)
        return 1

    exact: dict[str, list[tuple[str, int, str]]] = defaultdict(list)
    ci: dict[str, list[tuple[str, str, int, str]]] = defaultdict(list)
    for path in files:
        base = os.path.basename(path)
        for kind, name, ln in parse_file(path):
            exact[name].append((base, ln, kind))
            ci[name.lower()].append((name, base, ln, kind))

    cross_file = {
        n: locs for n, locs in exact.items() if len({f for f, _, _ in locs}) > 1
    }
    ci_dupes = {}
    for low, entries in ci.items():
        names = {n for n, _, _, _ in entries}
        files_set = {f for _, f, _, _ in entries}
        if len(names) > 1 and len(files_set) > 1:
            ci_dupes[low] = entries
    intra: dict[str, list[tuple[str, int, str]]] = {}
    for name, locs in exact.items():
        by_file: dict[str, list[tuple[int, str]]] = defaultdict(list)
        for f, ln, k in locs:
            by_file[f].append((ln, k))
        for f, occs in by_file.items():
            if len(occs) > 1:
                intra.setdefault(name, []).extend((f, ln, k) for ln, k in occs)

    print(f"Scanned {len(files)} constants_*.go files")
    print(f"Total unique top-level identifiers: {len(exact)}")

    failed = False
    if cross_file:
        failed = True
        print(f"\n::error::[1] {len(cross_file)} cross-file exact-name collision(s):")
        for name in sorted(cross_file):
            print(f"  {name}")
            for f, ln, k in cross_file[name]:
                print(f"      [{k}] {f}:{ln}")
    if ci_dupes:
        failed = True
        print(f"\n::error::[2] {len(ci_dupes)} cross-file case-insensitive collision(s):")
        for low in sorted(ci_dupes):
            print(f"  '{low}':")
            for name, f, ln, k in ci_dupes[low]:
                print(f"      [{k}] {name} @ {f}:{ln}")
    if intra:
        failed = True
        print(f"\n::error::[3] {len(intra)} intra-file duplicate declaration(s):")
        for name in sorted(intra):
            print(f"  {name}")
            for f, ln, k in intra[name]:
                print(f"      [{k}] {f}:{ln}")

    if failed:
        return 1
    print("\nOK No collisions detected.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
