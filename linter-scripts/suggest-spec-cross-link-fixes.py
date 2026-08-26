#!/usr/bin/env python3
"""suggest-spec-cross-link-fixes.py

Companion to ``check-spec-cross-links.py``. For every broken internal
markdown link in ``spec/``, propose the closest matching fix by fuzzy
matching against:

  * existing markdown files (for ``missing-file`` failures), and
  * existing headings inside the resolved target file (for
    ``missing-section`` failures).

Modes
-----
  --report   (default) Print suggestions, never touch files. Exit 0 always
             (purely advisory). Designed for CI annotations.
  --apply    Rewrite files in place when a suggestion has confidence
             >= --min-confidence (default 0.82). Anything below threshold
             is reported but NOT auto-applied.

Exit codes
----------
  0  Success (report mode always; apply mode when no IO errors).
  1  Apply mode: at least one broken link could not be auto-fixed
     (so CI can still gate on residual breakage if desired).
  2  Invocation error.

The matcher is deterministic: it relies on ``difflib.SequenceMatcher``
against the slug form of headings and the POSIX path form of files. No
network, no LLM.
"""
from __future__ import annotations

import argparse
import difflib
import json
import re
import sys
from pathlib import Path
from typing import Iterable

# Reuse the same regexes / helpers as the checker. We re-implement here
# (rather than importing) so this script stays a single-file dependency
# of CI and can be run standalone in any checkout.
MD_LINK_RE = re.compile(r"\[([^\]]+)\]\(([^)\s]+?)(?:\s+\"[^\"]*\")?\)")
HEADING_RE = re.compile(r"^(#{1,6})\s+(.+?)\s*$", re.MULTILINE)
EXTERNAL_PREFIXES = ("http://", "https://", "mailto:", "tel:", "ftp://", "#")
SKIP_SCHEMES = ("mem://", "user-uploads://", "knowledge://")

DEFAULT_MIN_CONFIDENCE = 0.82


def slugify(text: str) -> str:
    text = text.strip().lower()
    text = re.sub(r"[^\w\s-]", "", text, flags=re.UNICODE)
    text = re.sub(r"\s+", "-", text)
    return text.strip("-")


def is_external(target: str) -> bool:
    return target.startswith(EXTERNAL_PREFIXES) or target.startswith(SKIP_SCHEMES)


def strip_code_fences(text: str) -> str:
    out: list[str] = []
    in_fence = False
    fence = ""
    for line in text.splitlines():
        s = line.lstrip()
        opening = (s.startswith("```") or s.startswith("~~~")) and not in_fence
        closing = in_fence and s.startswith(fence)
        if opening:
            in_fence = True
            fence = "```" if s.startswith("```") else "~~~"
            out.append("")
            continue
        if closing:
            in_fence = False
            out.append("")
            continue
        out.append("" if in_fence else line)
        
    text = "\n".join(out)
    
    def _blank(match: re.Match[str]) -> str:
        body = match.group(0)
        return "".join(ch if ch == "\n" else " " for ch in body)
        
    text = re.sub(r"``.*?``", _blank, text, flags=re.DOTALL)
    text = re.sub(r"`[^`\n]+`", _blank, text)
    return text


def collect_headings(path: Path) -> list[tuple[str, str]]:
    """Return list of (raw_heading_text, slug) tuples."""
    try:
        content = path.read_text(encoding="utf-8", errors="ignore")
    except OSError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return []
    return [(m.group(2).strip(), slugify(m.group(2))) for m in HEADING_RE.finditer(content)]


def iter_markdown_files(root: Path) -> Iterable[Path]:
    for path in sorted(root.rglob("*.md")):
        yield path


def load_allowlist(repo_root: Path) -> set[str]:
    """Reuse the same waiver file used by check-spec-cross-links.py so the
    suggester operates on the same set of genuine failures."""
    path = repo_root / "linter-scripts" / "spec-cross-links.allowlist"
    if not path.exists():
        return set()
    out: set[str] = set()
    for raw in path.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if line and not line.startswith("#"):
            out.add(line)
    return out


def resolve_target_path(source: Path, target: str, repo_root: Path) -> Path:
    raw = target.split("#", 1)[0]
    if not raw:
        return source
    if raw.startswith("/"):
        return (repo_root / raw.lstrip("/")).resolve()
    return (source.parent / raw).resolve()


def best_file_match(missing_path: Path, candidates: list[Path]) -> tuple[Path, float] | None:
    """Find the candidate whose POSIX path is most similar to ``missing_path``.
    Bias toward files with the same basename for stronger signal.
    """
    if not candidates:
        return None
    target_str = missing_path.as_posix()
    target_name = missing_path.name
    same_name = [c for c in candidates if c.name == target_name]
    pool = same_name if same_name else candidates
    scored: list[tuple[Path, float]] = []
    for cand in pool:
        ratio = difflib.SequenceMatcher(None, target_str, cand.as_posix()).ratio()
        if same_name and cand in same_name:
            ratio = min(1.0, ratio + 0.05)  # small bonus for exact basename
        scored.append((cand, ratio))
    scored.sort(key=lambda x: x[1], reverse=True)
    return scored[0]


def best_heading_match(anchor: str, headings: list[tuple[str, str]]) -> tuple[str, str, float] | None:
    """Return (raw_heading, slug, score) of the heading whose slug best
    matches the requested anchor slug.
    """
    if not headings:
        return None
    target_slug = slugify(anchor)
    scored = [
        (raw, slug, difflib.SequenceMatcher(None, target_slug, slug).ratio())
        for raw, slug in headings
    ]
    scored.sort(key=lambda x: x[2], reverse=True)
    return scored[0]


def relativize(target_file: Path, source_file: Path) -> str:
    """Express ``target_file`` as a path relative to ``source_file``'s parent,
    using POSIX separators and a leading ``./`` for same-folder targets so
    the result matches the project's convention.
    """
    rel = Path(*target_file.relative_to(source_file.parent.anchor).parts) if False else None
    try:
        rel_path = Path(target_file).resolve().relative_to(source_file.parent.resolve())
        rel_str = "./" + rel_path.as_posix()
    except ValueError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        # Need ../ traversal — compute manually via os.path.relpath semantics.
        import os
        rel_str = os.path.relpath(target_file.resolve(), source_file.parent.resolve())
        rel_str = Path(rel_str).as_posix()
    return rel_str


def find_link_failures(root: Path, repo_root: Path) -> list[dict]:
    """Re-scan spec/ and return broken links the same way check-spec-cross-links does."""
    failures: list[dict] = []
    allowlist = load_allowlist(repo_root)
    for md in iter_markdown_files(root):
        try:
            text = md.read_text(encoding="utf-8", errors="ignore")
        except OSError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            continue
        scan_text = strip_code_fences(text)
        for match in MD_LINK_RE.finditer(scan_text):
            target = match.group(2).strip()
            if is_external(target):
                continue
            if "#" in target:
                path_part, anchor = target.split("#", 1)
            else:
                path_part, anchor = target, ""
            resolved = resolve_target_path(md, path_part or ".", repo_root)
            line_num = scan_text.count("\n", 0, match.start()) + 1
            rel_file = str(md.relative_to(repo_root))
            waiver_key = f"{rel_file}:{line_num}:{target}"
            if waiver_key in allowlist:
                continue
            entry = {
                "file": md,
                "line": line_num,
                "link_text": match.group(1),
                "target": target,
                "path_part": path_part,
                "anchor": anchor,
            }
            if not resolved.exists():
                entry["kind"] = "missing-file"
                failures.append(entry)
                continue
            if anchor:
                if resolved.is_dir():
                    entry["kind"] = "anchor-on-directory"
                    failures.append(entry)
                    continue
                slugs = {s for _, s in collect_headings(resolved)}
                if slugify(anchor) not in slugs:
                    entry["kind"] = "missing-section"
                    entry["resolved"] = resolved
                    failures.append(entry)
    return failures


def build_suggestions(failures: list[dict], root: Path, repo_root: Path) -> list[dict]:
    all_md = list(iter_markdown_files(root))
    suggestions: list[dict] = []
    for f in failures:
        if f["kind"] == "missing-file":
            missing = resolve_target_path(f["file"], f["path_part"] or ".", repo_root)
            best = best_file_match(missing, all_md)
            if not best:
                continue
            cand_file, score = best
            new_path = relativize(cand_file, f["file"])
            new_target = new_path + (f"#{f['anchor']}" if f["anchor"] else "")
            suggestions.append({**f, "suggestion": new_target, "confidence": round(score, 3)})
        elif f["kind"] == "missing-section":
            headings = collect_headings(f["resolved"])
            best = best_heading_match(f["anchor"], headings)
            if not best:
                continue
            raw, slug, score = best
            new_target = (f["path_part"] or "") + f"#{slug}"
            suggestions.append({
                **f,
                "suggestion": new_target,
                "matched_heading": raw,
                "confidence": round(score, 3),
            })
    return suggestions


def serialize(s: dict) -> dict:
    out = {k: v for k, v in s.items() if k not in {"file", "resolved"}}
    out["file"] = str(s["file"])
    if "resolved" in s:
        out["resolved"] = str(s["resolved"])
    return out


def emit_human(sugs: list[dict], min_conf: float) -> None:
    if not sugs:
        print("OK No broken cross-references — nothing to suggest.")
        return
    high = [s for s in sugs if s["confidence"] >= min_conf]
    low = [s for s in sugs if s["confidence"] < min_conf]
    print(f"INFO {len(sugs)} broken link(s); {len(high)} auto-fixable (>= {min_conf}), {len(low)} need manual review.\n")
    for s in sugs:
        marker = "AUTO" if s["confidence"] >= min_conf else "MANUAL"
        rel = Path(s["file"]).name
        print(f"  [{marker} {s['confidence']:.2f}] {rel}:{s['line']}  ({s['kind']})")
        print(f"    current:    {s['target']}")
        print(f"    suggestion: {s['suggestion']}")
        if "matched_heading" in s:
            print(f"    heading:    {s['matched_heading']}")
        print()


def emit_github(sugs: list[dict], min_conf: float) -> None:
    for s in sugs:
        level = "warning" if s["confidence"] >= min_conf else "notice"
        msg = f"Broken link → suggested fix: {s['suggestion']} (confidence {s['confidence']:.2f})"
        print(f"::{level} file={s['file']},line={s['line']}::{msg}")


def apply_fixes(sugs: list[dict], min_conf: float) -> tuple[int, int]:
    """Rewrite files in place. Returns (applied_count, skipped_count)."""
    by_file: dict[Path, list[dict]] = {}
    for s in sugs:
        if s["confidence"] < min_conf:
            continue
        by_file.setdefault(Path(s["file"]), []).append(s)
    applied = 0
    skipped = sum(1 for s in sugs if s["confidence"] < min_conf)
    for path, items in by_file.items():
        text = path.read_text(encoding="utf-8")
        # Sort longest target first so partial overlaps don't clobber each other.
        items.sort(key=lambda x: len(x["target"]), reverse=True)
        for it in items:
            old = f"]({it['target']})"
            new = f"]({it['suggestion']})"
            if old in text:
                text = text.replace(old, new, 1)
                applied += 1
        path.write_text(text, encoding="utf-8")
    return applied, skipped


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Suggest or auto-fix broken spec cross-references.")
    p.add_argument("--root", default="spec")
    p.add_argument("--repo-root", default=".")
    p.add_argument("--apply", action="store_true", help="Rewrite files in place when confidence >= threshold")
    p.add_argument("--min-confidence", type=float, default=DEFAULT_MIN_CONFIDENCE)
    p.add_argument("--json", action="store_true")
    p.add_argument("--github", action="store_true", help="Emit GitHub Actions annotations")
    return p.parse_args()


def main() -> int:
    args = parse_args()
    root = Path(args.root).resolve()
    repo_root = Path(args.repo_root).resolve()
    if not root.exists():
        print(f"::error::spec root not found: {root}", file=sys.stderr)
        return 2
    failures = find_link_failures(root, repo_root)
    sugs = build_suggestions(failures, root, repo_root)
    payload = {
        "broken_count": len(failures),
        "suggestion_count": len(sugs),
        "min_confidence": args.min_confidence,
        "suggestions": [serialize(s) for s in sugs],
    }
    if args.json:
        print(json.dumps(payload, indent=2))
    else:
        emit_human(sugs, args.min_confidence)
    if args.github:
        emit_github(sugs, args.min_confidence)
    if args.apply:
        applied, skipped = apply_fixes(sugs, args.min_confidence)
        print(f"\nAUTO-FIX applied={applied} skipped_low_confidence={skipped}")
        return 1 if skipped else 0
    return 0


if __name__ == "__main__":
    sys.exit(main())