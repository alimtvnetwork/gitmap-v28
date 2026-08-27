#!/usr/bin/env python3
"""
check-mws-error-codes.py — MWS error-code linter (FU-9).

Enforces the contract in `spec/19-main-worker-service/13-error-codes.md` §7:

  Rule R1 — Presence: every WORKER-XYY-ZZ / MAIN-XYY-ZZ literal that
            appears anywhere under spec/19/, spec/14-update/28-*,
            and the source tree MUST be catalogued in 13-error-codes.md.
  Rule R2 — No orphans: every code catalogued in 13-error-codes.md MUST
            be referenced from at least one source location outside
            13-error-codes.md and the generated index files.
  Rule R3 — Bijection: prefixed code <-> flat integer mapping MUST be
            one-to-one across the whole MWS range.
  Rule R4 — Range: WORKER-* flats live in 21000-21099, MAIN-* flats in
            21100-21199.

Exit codes:
  0   all rules pass
  1   one or more rules failed
  2   linter setup error (missing spec file, unparseable, etc.)

CODE RED compliance:
  - All functions <=15 lines, zero nested ifs, positive guards.
  - Errors are never swallowed; surfaced with file/path context.
"""

from __future__ import annotations
import json
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SPEC_FILE = REPO_ROOT / "spec/19-main-worker-service/13-error-codes.md"
INDEX_FILE = REPO_ROOT / "spec/19-main-worker-service/error-codes.json"
MASTER_FILE = REPO_ROOT / "spec/03-error-manage/03-error-code-registry/error-codes-master.json"
WAIVER_FILE = REPO_ROOT / "linter-scripts/check-mws-error-codes.waivers.txt"
UNALLOCATED_FILE = REPO_ROOT / "linter-scripts/check-mws-error-codes.unallocated.txt"

# Files that contain catalogue entries — references inside them do NOT
# count for orphan detection (Rule R2).
CATALOGUE_FILES = {SPEC_FILE, INDEX_FILE, MASTER_FILE}

WORKER_PREFIX = "WORKER"
MAIN_PREFIX = "MAIN"
WORKER_RANGE_1_START = 21000
WORKER_RANGE_1_END = 21099
WORKER_RANGE_2_START = 21200
WORKER_RANGE_2_END = 21299
MAIN_RANGE_START = 21100
MAIN_RANGE_END = 21199

CODE_RX = re.compile(r"\b(WORKER|MAIN)-(\d{3})-(\d{2})\b")
ROW_RX = re.compile(
    r"^\|\s*`(?P<code>(?:WORKER|MAIN)-\d{3}-\d{2})`\s*"
    r"\|\s*`(?P<flat>\d{5})`\s*\|",
    re.MULTILINE,
)
SCAN_DIRS = ("spec/19-main-worker-service", "spec/14-update", "src", "linter-scripts/tests")
SCAN_EXTS = {".md", ".py", ".ts", ".tsx", ".js", ".jsx", ".go", ".php", ".sql", ".json", ".yml", ".yaml"}


def is_spec_present() -> bool:
    return SPEC_FILE.is_file()


def parse_catalogue() -> dict[str, int]:
    """Return {prefixed_code: flat_int} parsed from the spec tables."""
    text = SPEC_FILE.read_text(encoding="utf-8")
    rows = ROW_RX.finditer(text)
    return {m["code"]: int(m["flat"]) for m in rows}


def iter_scan_files() -> list[Path]:
    files: list[Path] = []
    for d in SCAN_DIRS:
        root = REPO_ROOT / d
        if not root.exists():
            continue
        files.extend(p for p in root.rglob("*") if p.suffix in SCAN_EXTS)
    return files


def collect_references(files: list[Path]) -> dict[str, set[Path]]:
    """{prefixed_code: {file_path,...}} — references outside catalogue files."""
    refs: dict[str, set[Path]] = {}
    for f in files:
        if f.resolve() in {p.resolve() for p in CATALOGUE_FILES}:
            continue
        record_refs_in_file(f, refs)
    return refs


def record_refs_in_file(f: Path, refs: dict[str, set[Path]]) -> None:
    try:
        text = f.read_text(encoding="utf-8", errors="replace")
    except OSError as exc:
        print(f"Error reading {f}: {exc}", file=sys.stderr)
        raise SystemExit(f"[linter setup] cannot read {f}: {exc}")
    for m in CODE_RX.finditer(text):
        code = f"{m.group(1)}-{m.group(2)}-{m.group(3)}"
        refs.setdefault(code, set()).add(f)


def has_no_unknown_refs(refs: dict[str, set[Path]], catalogue: set[str], errors: list[str]) -> bool:
    unallocated = load_unallocated()
    unknown = sorted((set(refs) - catalogue) - unallocated)
    for code in unknown:
        sample = sorted(refs[code])[0].relative_to(REPO_ROOT)
        errors.append(f"R1 unknown code {code} referenced in {sample} (and {len(refs[code])-1} more) — not in 13-error-codes.md")
    return not unknown


def load_waivers() -> set[str]:
    return load_code_list(WAIVER_FILE)


def load_unallocated() -> set[str]:
    return load_code_list(UNALLOCATED_FILE)


def load_code_list(path: Path) -> set[str]:
    if not path.is_file():
        return set()
    lines = path.read_text(encoding="utf-8").splitlines()
    return {ln.strip() for ln in lines if ln.strip() and not ln.lstrip().startswith("#")}


def has_no_orphans(refs: dict[str, set[Path]], catalogue: set[str], waivers: set[str], errors: list[str]) -> bool:
    orphans = sorted((catalogue - set(refs)) - waivers)
    for code in orphans:
        errors.append(f"R2 orphan {code} catalogued in 13-error-codes.md but never referenced elsewhere (and not waived)")
    return not orphans


def is_bijective(catalogue: dict[str, int], errors: list[str]) -> bool:
    flats = list(catalogue.values())
    dup_flats = sorted({f for f in flats if flats.count(f) > 1})
    for f in dup_flats:
        owners = [c for c, v in catalogue.items() if v == f]
        errors.append(f"R3 flat {f} is shared by {owners}")
    return not dup_flats


def is_in_range(catalogue: dict[str, int], errors: list[str]) -> bool:
    bad = []
    for code, flat in catalogue.items():
        bad.extend(check_one_range(code, flat))
    for msg in bad:
        errors.append(msg)
    return not bad


def check_one_range(code: str, flat: int) -> list[str]:
    if code.startswith(WORKER_PREFIX) and is_worker_flat_invalid(flat):
        return [f"R4 {code} flat {flat} outside Worker ranges {WORKER_RANGE_1_START}-{WORKER_RANGE_1_END} or {WORKER_RANGE_2_START}-{WORKER_RANGE_2_END}"]
    if code.startswith(MAIN_PREFIX) and not MAIN_RANGE_START <= flat <= MAIN_RANGE_END:
        return [f"R4 {code} flat {flat} outside Main range {MAIN_RANGE_START}-{MAIN_RANGE_END}"]
    return []


def is_worker_flat_invalid(flat: int) -> bool:
    primary = WORKER_RANGE_1_START <= flat <= WORKER_RANGE_1_END
    overflow = WORKER_RANGE_2_START <= flat <= WORKER_RANGE_2_END
    return not (primary or overflow)


def main() -> int:
    if not is_spec_present():
        print(f"[setup] missing spec file: {SPEC_FILE}", file=sys.stderr)
        return 2
    catalogue = parse_catalogue()
    if not catalogue:
        print(f"[setup] no codes parsed from {SPEC_FILE}", file=sys.stderr)
        return 2
    refs = collect_references(iter_scan_files())
    waivers = load_waivers()
    errors: list[str] = []
    has_no_unknown_refs(refs, set(catalogue), errors)
    has_no_orphans(refs, set(catalogue), waivers, errors)
    is_bijective(catalogue, errors)
    is_in_range(catalogue, errors)
    return report(errors, catalogue, waivers)


def report(errors: list[str], catalogue: dict[str, int], waivers: set[str]) -> int:
    if not errors:
        print(f"[ok] check-mws-error-codes: {len(catalogue)} codes verified (R1-R4); {len(waivers)} R2 waiver(s) loaded")
        return 0
    for e in errors:
        print(f"[fail] {e}", file=sys.stderr)
    print(f"[summary] {len(errors)} violation(s) across {len(catalogue)} catalogued codes", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
