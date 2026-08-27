#!/usr/bin/env python3
"""check-placeholder-comments.py

SPEC-PLACEHOLDER-001 — Lint placeholder HTML comment blocks in spec files.

The cross-link checker (``check-spec-cross-links.py``) ignores links inside
``<spec-placeholder>...</spec-placeholder>`` blocks (and, for backward
compatibility, the older ``<!-- TODO|FIXME ... -->`` HTML-comment form)
so authors can stash pending references without breaking CI. The
trade-off is that *malformed* placeholder blocks (missing markers,
broken bullet rows, stray text) silently slip through — the very rot
the placeholder was supposed to prevent.

This script validates every placeholder block of either supported
format, per the conventions in ``spec/_template.md``
§"Placeholder cross-references":

    <spec-placeholder reason="activate when target is created">
    - [Target Title](../NN-module-name/00-overview.md)
    - [Target Title](../NN-module-name/01-file-name.md#section-anchor)
    </spec-placeholder>

    <!-- legacy form, still supported -->
    <!-- TODO: activate when target is created
    - [Target Title](../NN-module-name/00-overview.md)
    - [Target Title](../NN-module-name/01-file-name.md#section-anchor)
    -->

Rules enforced (lightweight, no AST):

  P-001  Placeholder *intent text* must be a complete imperative
         sentence so reviewers see actionable language. Specifically:

           * Wording follows the marker (``TODO:`` / ``FIXME:`` for
             legacy comments, ``reason="…"`` for ``<spec-placeholder>``).
           * It must be non-empty and start with a recognised
             imperative verb (``activate``, ``add``, ``link``,
             ``replace``, ``wire``, ``update``, ``write``, ``create``,
             ``document``, ``cross-reference``). Extend the allowlist
             via ``--allow-verb <verb>`` (repeatable).
           * It must end with a period — half-sentences like
             ``activate later`` are rejected; ``Activate when target is
             created.`` passes.

         The verb is matched case-insensitively. Articles/auxiliaries
         like ``please`` are stripped before the check so
         ``please add the link.`` passes.
  P-002  Every non-blank body line must be a markdown bullet
         (``- [text](link)``) — no stray prose, no orphan list markers.
  P-003  Bullet links must be relative paths ending in ``.md``
         (optionally with ``#anchor``); ``http(s)://`` and bare anchors
         are rejected because placeholders are meant for *future*
         internal targets, not external references.
  P-004  Block must contain at least one bullet (empty placeholders are
         dead code).
  P-005  Block must not contain blank lines (keeps the snippet
         contiguous per template guidance).
  P-006  Every opening marker must have a matching closer
         (``-->`` or ``</spec-placeholder>``).
  P-007  Two or more placeholder bullets must not point at the same
         target ``.md`` file (anchor ignored — duplicates pointing at
         different sections of the same file still collapse to one
         pending activation). Detected within a file and across files.

Only multi-line comment blocks that start with the ``TODO:``/``FIXME:``
marker on the opening line are linted. Single-line comments and
non-placeholder comments (e.g. licence headers) are left alone. The
``<spec-placeholder>`` form is *always* linted because the tag itself
declares intent.

Exit codes:
  0  = no malformed placeholder blocks
  1  = one or more violations found
  2  = invocation error

Usage:
  python3 linter-scripts/check-placeholder-comments.py [--root spec] [--json]
"""
from __future__ import annotations

import argparse
import csv as _csv
import json
import hashlib
import os
import re
import subprocess
import sys
from dataclasses import dataclass, asdict
from pathlib import Path
from typing import Iterable


# --- HTML-comment placeholder (legacy form) ---------------------------
# Opening marker must be on the same line as ``<!--`` so we can detect
# placeholder intent without scanning the whole comment first.
PLACEHOLDER_OPEN_RE = re.compile(r"<!--\s*(TODO|FIXME)\b[^\n]*$")
# Capture the wording that follows the marker so P-001 can lint it.
PLACEHOLDER_INTENT_RE = re.compile(
    r"<!--\s*(?P<marker>TODO|FIXME)\s*:?\s*(?P<text>[^\n]*?)\s*$"
)
COMMENT_CLOSE = "-->"

# --- Custom-tag placeholder (preferred form) --------------------------
# Single-line opener pattern: ``<spec-placeholder ...>`` (attributes
# optional). The closer is ``</spec-placeholder>`` on its own line or
# at end-of-line. Self-closing ``<spec-placeholder/>`` is rejected as
# P-004 (empty block).
TAG_OPEN_RE = re.compile(r"<spec-placeholder\b[^>]*>")
TAG_SELF_CLOSE_RE = re.compile(r"<spec-placeholder\b[^>]*/\s*>")
TAG_CLOSE = "</spec-placeholder>"
# Capture the ``reason="…"`` attribute on the opening tag (single or
# double quotes). Absent reason → P-001 with a "missing reason" hint.
TAG_REASON_RE = re.compile(
    r"<spec-placeholder\b[^>]*?\breason\s*=\s*(?P<q>[\"'])(?P<text>.*?)(?P=q)",
    re.IGNORECASE,
)

# Curated set of imperative verbs that signal actionable intent. Kept
# deliberately small — authors who want a different verb can extend
# the set via ``--allow-verb`` rather than the linter silently accepting
# any leading word. All entries are lowercase; matching is
# case-insensitive.
DEFAULT_INTENT_VERBS: frozenset[str] = frozenset({
    "activate",
    "add",
    "link",
    "replace",
    "wire",
    "update",
    "write",
    "create",
    "document",
    "cross-reference",
})

# Soft prefixes that authors may stack before the imperative verb
# without the wording becoming non-actionable. Stripped (case-
# insensitively) before the verb-match check.
INTENT_PREFIXES: tuple[str, ...] = ("please ",)

BULLET_LINK_RE = re.compile(r"^-\s+\[[^\]]+\]\(([^)\s]+)\)\s*$")


# --- Suggested-patch fix hints ----------------------------------------
# Per-rule one-line replacement scaffold inserted into the suggested
# `git apply` patch in place of the offending post-state line. The
# linter cannot infer the author's correct fix — these are TODO
# markers labelled with the rule code so a reviewer doing a copy-
# paste apply immediately sees what kind of edit is required.
#
# Rules NOT in this table (e.g. the structural P-006 "missing closer"
# where the violation line is the *opener*, not a wrong line in
# place) fall back to ``_RULE_FIX_FALLBACK``. Adding a more specific
# hint later is purely additive and does not break the schema.
_RULE_FIX_HINTS: dict[str, str] = {
    "P-001": '<!-- TODO(P-001): replace with a complete imperative '
             'reason="…" (e.g. reason="Document RAG eviction policy"). -->',
    "P-002": "- [TODO(P-002): describe target](relative/path/to/spec.md)",
    "P-003": "- [TODO(P-003): use a relative .md path](relative/path/to/spec.md)",
    "P-004": "- [TODO(P-004): add at least one bullet](relative/path/to/spec.md)",
    "P-005": "<!-- TODO(P-005): remove the blank line above; "
             "placeholder bodies must be contiguous. -->",
    "P-006": "<!-- TODO(P-006): add a matching closing marker "
             "(--> or </spec-placeholder>) below this opener. -->",
    "P-007": "- [TODO(P-007): point at a different target — "
             "duplicate of an earlier placeholder](relative/path/to/different-spec.md)",
    "P-008": "<!-- TODO(P-008): see linter rule docs for the exact fix. -->",
}
_RULE_FIX_FALLBACK = (
    "<!-- TODO: see linter-scripts/check-placeholder-comments.py "
    "for the rule's required fix. -->"
)


@dataclass(frozen=True)
class Violation:
    file: str
    line: int
    code: str
    message: str


# --- Diff-mode audit trail --------------------------------------------
# A ``--list-changed-files`` row classifying one path that the diff-
# mode intake considered. Statuses are a closed set so downstream
# consumers (CI dashboards, JSON parsers) can switch on the literal
# value:
#
#   matched              — picked up for linting (under --root,
#                          extension is in the allowlist, file exists)
#   ignored-extension    — under --root but extension not allowed
#                          (e.g. a `.txt` change when the allowlist
#                          is `md`/`mdx`)
#   ignored-out-of-root  — repo path outside --root (e.g. a README
#                          change while linting `spec/`)
#   ignored-missing      — A/M/R/C row whose post-state path no
#                          longer exists on disk (reverted in a later
#                          commit of the same push)
#   ignored-deleted      — git emitted a D-status row, or a rename's
#                          OLD side that the linter intentionally
#                          drops because there's no post-state file
#                          to scan
#
# ``reason`` is a one-line human-readable explanation safe to print
# in a CI log; never None. Rows are emitted in stable input order so
# diffs against a previous run are reviewable.
_AUDIT_STATUSES: tuple[str, ...] = (
    "matched",
    "ignored-extension",
    "ignored-out-of-root",
    "ignored-missing",
    "ignored-deleted",
)


@dataclass(frozen=True)
class _ChangedFileAudit:
    path: str
    status: str
    reason: str
    # Optional rename/copy provenance. Populated only when the diff-
    # mode intake observed an ``R``/``C`` row that resolved to this
    # ``path`` (the post-rename / post-copy "new" side). ``None`` for
    # plain ``A``/``M``/``D`` rows AND for ``R``/``C`` rows that
    # failed earlier filters before reaching the audit constructor.
    # The audit *renderer* (``--with-similarity``) substitutes a
    # dash ("-") wherever this is None or a sub-field is missing — we
    # don't bake the dash into the data so JSON consumers see real
    # ``null``s rather than a sentinel string.
    similarity: "_RenameSimilarity | None" = None
    # Optional intake-provenance tag. Populated only on
    # ``ignored-deleted`` rows; ``None`` everywhere else. The value is
    # one of :data:`_DELETED_REASON`'s keys (today: ``"diff-D"`` for
    # a true ``git diff --name-status`` ``D`` row, or
    # ``"changed-files-D"`` for an authored ``--changed-files``
    # payload row shaped exactly ``D\tpath``). Surfaced verbatim by
    # ``--list-changed-files-verbose`` so a CI reviewer can tell
    # which intake produced a given delete without parsing the
    # ``reason`` text. We deliberately keep the raw tag instead of
    # re-deriving it from the reason string: the reason wording is
    # human-readable and may be re-worded for clarity, but the
    # source vocabulary is part of the machine contract once the
    # verbose flag is on.
    source: "str | None" = None


@dataclass(frozen=True)
class _RenameSimilarity:
    """Rename/copy provenance for one ``_ChangedFileAudit`` row.

    Captures the three pieces of metadata git emits on an ``R`` / ``C``
    name-status row: the kind letter (``"R"`` or ``"C"``), the
    similarity score (0–100, or ``None`` when the row was scoreless —
    e.g. an authored ``--changed-files`` payload that wrote ``R\\told\\tnew``
    without a percentage), and the OLD-side path.

    Frozen + slot-free so it's hashable and round-trips through
    :func:`dataclasses.asdict` without surprises.
    """
    kind: str
    score: int | None
    old_path: str


# Default extension allowlist for spec discovery. Kept as a tuple so
# the value is hashable + cache-segment-friendly. Extending this set
# at runtime is exposed via ``--extension`` (repeatable) and feeds
# both the file iterator and the cache-segment naming below.
DEFAULT_EXTENSIONS: tuple[str, ...] = ("md",)


def iter_markdown_files(
    root: Path,
    *,
    extensions: tuple[str, ...] = DEFAULT_EXTENSIONS,
) -> Iterable[Path]:
    """Yield every file under ``root`` matching one of ``extensions``,
    sorted by path with hidden directories (``.foo/``) excluded.

    ``extensions`` is a tuple of bare extension strings without the
    leading dot (e.g. ``("md", "mdx")``). The function unions the
    per-extension globs into a single sorted, deduplicated stream so
    a future ``--extension md --extension mdx`` run can't yield the
    same path twice (e.g. on a case-insensitive filesystem).

    Extension matching is **case-insensitive** to keep behaviour
    identical across platforms. On Linux CI ``rglob("*.md")`` only
    finds the lowercase form, so a file checked in as ``README.MD``
    (legitimate on Windows + macOS where the FS folds case) would be
    silently invisible to the linter — but the diff-mode audit (which
    lowercases the suffix at classification time) would still mark it
    ``matched``. That asymmetry produced a "passes locally, fails on
    Windows" class of bug. We now lowercase the suffix in the
    iterator too, so the full-tree walk and the diff-mode
    classification agree on every platform.
    """
    # Pre-compute the lowercase, dot-prefixed allowlist once so the
    # per-file check is a single set lookup. Matches the shape of
    # ``allowed_exts`` inside :func:`_resolve_changed_md` so the two
    # codepaths apply identical rules.
    allowed = {("." + e.lstrip(".").lower()) for e in extensions}
    seen: set[Path] = set()
    # ``rglob("*")`` walks every entry once and we filter by suffix
    # ourselves. This is the same number of FS calls as the previous
    # per-extension globs (which all walked the tree internally) and
    # gives us case-insensitive matching for free.
    for p in root.rglob("*"):
        if not p.is_file():
            continue
        if p.suffix.lower() not in allowed:
            continue
        if p in seen:
            continue
        seen.add(p)
    for p in sorted(seen):
        if any(part.startswith(".") for part in p.relative_to(root).parts):
            continue
        yield p


def strip_code_fences(text: str) -> str:
    """Blank out fenced code blocks while preserving line numbers.

    Documentation that *shows* a placeholder snippet inside ```` ```markdown ````
    fences would otherwise be linted as a real placeholder. We replace
    every line inside a fence with an empty string but keep the newline,
    so reported line numbers in surrounding prose stay accurate.
    """
    out: list[str] = []
    in_fence = False
    fence_marker = ""
    for line in text.splitlines():
        stripped = line.lstrip()
        if not in_fence and (stripped.startswith("```") or stripped.startswith("~~~")):
            in_fence = True
            fence_marker = stripped[:3]
            out.append("")
            continue
        if in_fence and stripped.startswith(fence_marker):
            in_fence = False
            out.append("")
            continue
        if in_fence:
            out.append("")
        else:
            out.append(line)
    return "\n".join(out)


# Inline code spans (`...`) on a single line. We blank the contents but
# preserve the line so subsequent line numbers stay accurate. Multi-line
# spans are not standard CommonMark, so per-line handling is enough.
INLINE_CODE_RE = re.compile(r"`+[^`\n]*?`+")


def strip_inline_code(text: str) -> str:
    return "\n".join(INLINE_CODE_RE.sub(lambda m: " " * len(m.group(0)), line)
                     for line in text.splitlines())


def _validate_intent(rel: str, line_no: int, marker: str, text: str,
                     out: list[Violation], verbs: frozenset[str]) -> None:
    """Apply P-001 to a placeholder's intent text.

    ``marker`` is shown verbatim in the message ("TODO:", "FIXME:",
    or "reason"). ``text`` is the wording that follows it (already
    stripped of surrounding whitespace). ``verbs`` is the active
    imperative-verb allowlist for this run.
    """
    if not text:
        out.append(Violation(rel, line_no, "P-001",
            f"Placeholder `{marker}` is empty — describe the pending action "
            "(e.g. `Activate when target is created.`)."))
        return

    # Strip soft prefixes (e.g. "please ") so they don't shadow the verb.
    lowered = text.lower()
    for prefix in INTENT_PREFIXES:
        if lowered.startswith(prefix):
            text = text[len(prefix):]
            lowered = text.lower()
            break

    if not text:
        out.append(Violation(rel, line_no, "P-001",
            f"Placeholder `{marker}` has no actionable wording after `please`."))
        return

    # Extract the leading word (or hyphenated compound like
    # ``cross-reference``). Compare case-insensitively.
    head_match = re.match(r"[A-Za-z][A-Za-z-]*", text)
    head = head_match.group(0).lower() if head_match else ""
    if head not in verbs:
        sample = ", ".join(sorted(list(verbs))[:6])
        out.append(Violation(rel, line_no, "P-001",
            f"Placeholder `{marker}` must start with an imperative verb "
            f"(got `{head or text[:20]}`). Allowed verbs include: {sample}…. "
            "Extend with `--allow-verb <verb>` if needed."))
        return

    if not text.rstrip().endswith("."):
        out.append(Violation(rel, line_no, "P-001",
            f"Placeholder `{marker}` wording must end with a period "
            f"(got `{text.rstrip()[-30:]}`)."))
        return


def _validate_body(rel: str, open_line: int, body: list[tuple[int, str]],
                   out: list[Violation],
                   bullets: list[tuple[int, str]] | None = None) -> int:
    """Apply P-002/P-003/P-005 to a body and return valid bullet count.

    When ``bullets`` is provided, every valid bullet is appended as
    ``(line, target)`` for later cross-block duplicate analysis (P-007).
    """
    bullet_count = 0
    for ln, content in body:
        if not content.strip():
            out.append(Violation(rel, ln, "P-005",
                "Blank line inside placeholder block; keep it contiguous."))
            continue
        bm = BULLET_LINK_RE.match(content)
        if not bm:
            out.append(Violation(rel, ln, "P-002",
                "Placeholder body line is not a `- [text](link)` bullet."))
            continue
        target = bm.group(1)
        if target.startswith(("http://", "https://", "mailto:", "#")):
            out.append(Violation(rel, ln, "P-003",
                f"Placeholder link `{target}` must be a relative `.md` path, "
                "not external/anchor-only."))
            continue
        path_part = target.split("#", 1)[0]
        if not path_part.endswith(".md"):
            out.append(Violation(rel, ln, "P-003",
                f"Placeholder link `{target}` must point at a `.md` file."))
            continue
        bullet_count += 1
        if bullets is not None:
            bullets.append((ln, target))
    return bullet_count


def lint_file(path: Path, repo_root: Path,
              valid_bullets: list[tuple[str, int, str]] | None = None,
              intent_verbs: frozenset[str] = DEFAULT_INTENT_VERBS,
              ) -> list[Violation]:
    """Lint one markdown file.

    When ``valid_bullets`` is provided, every successfully-validated
    bullet is appended as ``(rel_file, line, target)`` so the caller
    can run cross-file duplicate detection (P-007).

    ``intent_verbs`` controls the imperative-verb allowlist for P-001;
    defaults to ``DEFAULT_INTENT_VERBS`` and can be widened from the
    CLI via ``--allow-verb``.
    """
    rel = str(path.relative_to(repo_root))
    text = path.read_text(encoding="utf-8")
    lines = strip_inline_code(strip_code_fences(text)).splitlines()
    out: list[Violation] = []
    file_bullets: list[tuple[int, str]] = []

    i = 0
    n = len(lines)
    while i < n:
        line = lines[i]

        # ---- Custom-tag placeholder (preferred) ---------------------
        tag_self = TAG_SELF_CLOSE_RE.search(line)
        if tag_self:
            out.append(Violation(rel, i + 1, "P-004",
                "Self-closing `<spec-placeholder/>` has no bullet rows; remove or expand it."))
            i += 1
            continue
        tag_open = TAG_OPEN_RE.search(line)
        if tag_open:
            open_line = i + 1
            # P-001: validate `reason="…"` wording. A missing attribute
            # is itself a P-001 ("placeholder has no documented intent").
            reason_match = TAG_REASON_RE.search(line)
            reason_text = reason_match.group("text").strip() if reason_match else ""
            if not reason_match:
                out.append(Violation(rel, open_line, "P-001",
                    "`<spec-placeholder>` is missing a `reason=\"…\"` attribute "
                    "describing why the link is pending."))
            else:
                _validate_intent(rel, open_line, "reason", reason_text, out, intent_verbs)
            # Same-line open+close — degenerate empty block.
            if TAG_CLOSE in line[tag_open.end():]:
                out.append(Violation(rel, open_line, "P-004",
                    "`<spec-placeholder>` block is empty (no bullet rows)."))
                i += 1
                continue
            body, i, closed = _consume_block(lines, i + 1, TAG_CLOSE)
            if not closed:
                out.append(Violation(rel, open_line, "P-006",
                    "`<spec-placeholder>` opened but never closed "
                    "(missing `</spec-placeholder>`)."))
                continue
            bullet_count = _validate_body(rel, open_line, body, out, file_bullets)
            if bullet_count == 0:
                out.append(Violation(rel, open_line, "P-004",
                    "`<spec-placeholder>` block contains no valid bullet rows."))
            continue

        # ---- HTML-comment placeholder (legacy) ----------------------
        m = PLACEHOLDER_OPEN_RE.search(line)
        if not m:
            i += 1
            continue
        # P-001: lint the wording that follows the marker.
        intent_match = PLACEHOLDER_INTENT_RE.search(line)
        if intent_match:
            marker = intent_match.group("marker") + ":"
            intent_text = intent_match.group("text").strip()
            # Trim a trailing ``-->`` if the open + close share a line —
            # P-004 below will catch the structural problem; here we
            # just need clean intent text.
            if intent_text.endswith(COMMENT_CLOSE):
                intent_text = intent_text[: -len(COMMENT_CLOSE)].rstrip()
            _validate_intent(rel, i + 1, marker, intent_text, out, intent_verbs)
        if COMMENT_CLOSE in line[m.end():] or COMMENT_CLOSE in line[m.start():]:
            out.append(Violation(rel, i + 1, "P-004",
                "Placeholder comment has no bullet rows; remove or expand it."))
            i += 1
            continue
        open_line = i + 1
        body, i, closed = _consume_block(lines, i + 1, COMMENT_CLOSE)
        if not closed:
            out.append(Violation(rel, open_line, "P-006",
                "Placeholder comment opened but never closed (missing `-->`)."))
            continue
        bullet_count = _validate_body(rel, open_line, body, out, file_bullets)
        if bullet_count == 0:
            out.append(Violation(rel, open_line, "P-004",
                "Placeholder block contains no valid bullet rows."))

    # ---- P-007 within-file duplicates ------------------------------
    # Resolve each bullet to a canonical (file, target_path) key. We
    # strip the anchor because two placeholders pointing at different
    # sections of the same target file still collapse to a single
    # activation step, which is what P-007 is designed to surface.
    seen: dict[str, tuple[int, str]] = {}
    for ln, target in file_bullets:
        key = _canonical_target(rel, target, repo_root)
        if key in seen:
            first_ln, first_target = seen[key]
            out.append(Violation(rel, ln, "P-007",
                f"Duplicate placeholder target `{target}` — already "
                f"declared at L{first_ln} as `{first_target}` "
                "(anchor differences are ignored)."))
        else:
            seen[key] = (ln, target)

    if valid_bullets is not None:
        for ln, target in file_bullets:
            valid_bullets.append((rel, ln, target))

    return out


def _consume_block(lines: list[str], start: int, close_marker: str
                   ) -> tuple[list[tuple[int, str]], int, bool]:
    """Walk ``lines[start:]`` collecting body rows until ``close_marker``.

    Returns ``(body, next_index, closed)`` where ``body`` is a list of
    ``(line_number, content)`` tuples (1-indexed line numbers) and
    ``next_index`` is the position to resume scanning from.
    """
    body: list[tuple[int, str]] = []
    n = len(lines)
    i = start
    while i < n:
        cur = lines[i]
        if close_marker in cur:
            pre = cur.split(close_marker, 1)[0]
            if pre.strip():
                body.append((i + 1, pre.rstrip()))
            return body, i + 1, True
        body.append((i + 1, cur.rstrip()))
        i += 1
    return body, i, False


def _canonical_target(source_rel: str, target: str, repo_root: Path) -> str:
    """Resolve a placeholder bullet's link to a canonical repo-relative
    path string. Anchor is stripped so different anchors on the same
    target file collapse to the same key. Path resolution is purely
    syntactic (no I/O) — the target file may not exist yet, which is
    the whole point of placeholders.
    """
    path_part = target.split("#", 1)[0]
    source_dir = (repo_root / source_rel).parent
    try:
        resolved = (source_dir / path_part).resolve()
        return str(resolved.relative_to(repo_root.resolve()))
    except (ValueError, OSError) as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        # Fall back to the literal path if it escapes the repo root —
        # still gives consistent grouping for duplicate detection.
        return path_part


def main(argv: list[str] | None = None) -> int:
    import sys
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--root", default="spec",
        help="Directory to scan recursively for markdown files (default: spec).")
    ap.add_argument("--repo-root", default=".",
        help="Repository root for relative path reporting (default: cwd).")
    ap.add_argument("--json", action="store_true",
        help="Emit findings as JSON instead of human text.")
    ap.add_argument("--allow-verb", action="append", default=[],
        metavar="VERB",
        help="Add an extra imperative verb to the P-001 allowlist "
             "(repeatable). Use lowercase, hyphens allowed.")
    ap.add_argument("--extension", action="append", default=None,
        metavar="EXT",
        help="Restrict spec discovery to files with this extension "
             "(repeatable, no leading dot, case-sensitive). Default "
             "is `md`. Each unique sorted allowlist gets its own "
             "physical cache subdirectory under --cache-dir (e.g. "
             "`<cache-dir>/ext-md/`, `<cache-dir>/ext-md+mdx/`) so "
             "switching the allowlist never reads a sentinel written "
             "for a different file set — even if a future bug made "
             "the content hash collide. Side benefit: `rm -rf "
             "<cache-dir>/ext-mdx/` nukes one allowlist's sentinels "
             "without touching the others.")
    ap.add_argument("--include-mdx", action="store_true",
        help="Convenience shortcut for `--extension mdx`: scan `.mdx` "
             "files in addition to whatever the active extension "
             "allowlist already includes (which still defaults to "
             "`.md`). Composes with explicit `--extension` flags — "
             "the resulting allowlist is the union, deduped and "
             "sorted, so `--include-mdx` and `--extension mdx` land "
             "in the same cache segment (`ext-md+mdx/`). Use this "
             "when your spec tree mixes Markdown and MDX (e.g. a "
             "Docusaurus site) and you want the linter to cover "
             "both with a single short flag instead of repeating "
             "`--extension md --extension mdx` in every CI invocation. "
             "UNION vs REPLACEMENT: the flag is ALWAYS a union — it "
             "appends `mdx` to the active baseline, never replaces "
             "it. The baseline is `(md,)` when no `--extension` is "
             "given, OR exactly the set you passed via `--extension` "
             "when one or more of those flags are present. Examples: "
             "(1) bare `--include-mdx` → `(md, mdx)` because the "
             "default baseline is preserved; (2) `--extension rst "
             "--include-mdx` → `(rst, mdx)` because YOU dropped `md` "
             "from the baseline by passing an explicit `--extension`; "
             "(3) `--extension md --extension mdx` ≡ `--include-mdx` "
             "(same canonical cache segment + key). The rule of "
             "thumb: if you want `md` in the result, either omit "
             "`--extension` or include `--extension md` explicitly — "
             "`--include-mdx` will not re-add it for you.")
    ap.add_argument("--include-txt", action="store_true",
        help="Convenience shortcut for `--extension txt`: scan `.txt` "
             "files in addition to whatever the active extension "
             "allowlist already includes (which still defaults to "
             "`.md`). Mirrors `--include-mdx` exactly — same union "
             "semantics, same cache-segment canonicalisation. The "
             "resulting allowlist is the union, deduped and sorted, "
             "so `--include-txt` and `--extension txt` land in the "
             "same cache segment (`ext-md+txt/`). Useful when your "
             "spec tree mixes Markdown with plain-text artifacts "
             "(e.g. legacy `README.txt`, hand-authored release notes) "
             "and you want the linter to cover both with one short "
             "flag instead of repeating `--extension md --extension "
             "txt` in every CI invocation. UNION vs REPLACEMENT: "
             "always a union — appends `txt` to the active baseline, "
             "never replaces it. The baseline is `(md,)` when no "
             "`--extension` is given, OR exactly the set you passed "
             "via `--extension` when one or more of those flags are "
             "present. Examples: (1) bare `--include-txt` → `(md, "
             "txt)`; (2) `--extension rst --include-txt` → `(rst, "
             "txt)` because YOU dropped `md` from the baseline; (3) "
             "`--extension md --extension txt` ≡ `--include-txt` "
             "(same canonical cache segment + key). Composes with "
             "`--include-mdx`: passing both yields `(md, mdx, txt)` "
             "(or `(md, txt, mdx)` — order is CLI-order-stable for "
             "diagnostics, but the cache segment + key sort "
             "independently so the two orderings collapse).")
    ap.add_argument("--cache-dir", default=None, metavar="DIR",
        help="Enable a content-addressed PASS cache. On a hit (the linter "
             "script + every scanned `.md` hash to the same key as a "
             "previously-cached PASS) the scan is skipped and exit 0 is "
             "returned immediately. Misses run normally and write a fresh "
             "sentinel only on success. Stale or poisoned sentinels are "
             "ignored because the key is recomputed from the working tree "
             "every run. Sentinels are stored under an extension-derived "
             "subdirectory (see --extension) so different allowlists "
             "never share a sentinel pool.")
    ap.add_argument("--no-cache-write", action="store_true",
        help="With --cache-dir, read the sentinel but never write it. "
             "Useful for read-only / forked-repo CI runs.")
    ap.add_argument("--diff-base", default=None, metavar="REF",
        help="Diff-mode: only report per-file violations (P-001…P-006, "
             "P-008, within-file P-007) for `.md` files changed vs. "
             "REF (e.g. `origin/main`, `HEAD~1`, or a SHA). Resolved "
             "with `git diff --name-only --diff-filter=AM REF...HEAD` "
             "so renames + deletions are excluded. Cross-file P-007 "
             "still scans the full tree so new duplicates introduced "
             "by a changed file always surface, even if the colliding "
             "first declaration lives in an unchanged file. Mutually "
             "exclusive with --changed-files. Numeric shorthand: a "
             "bare positive integer N expands to `HEAD~N` (e.g. "
             "`--diff-base 1` ≡ `--diff-base HEAD~1`); leading `~N` "
             "or `^N` are also expanded against `HEAD` so "
             "`--diff-base ~2` ≡ `HEAD~2`. Any other value is passed "
             "to git verbatim.")
    ap.add_argument("--diff-prev", nargs="?", const="1", default=None,
        metavar="N",
        help="Diff-mode shorthand: compare against `HEAD~N` (default "
             "N=1, i.e. the previous commit). Equivalent to "
             "`--diff-base HEAD~N`. Mutually exclusive with "
             "--diff-base and --changed-files.")
    ap.add_argument("--changed-files", default=None, metavar="PATH",
        help="Diff-mode: read the changed-file list from PATH (one "
             "repo-relative path per line, blanks/`#` comments ignored) "
             "instead of invoking git. Use `-` to read from stdin. "
             "Same semantics as --diff-base for cross-file P-007. "
             "Mutually exclusive with --diff-base.")
    ap.add_argument("--diff-empty-passes", action="store_true",
        help="With --diff-base/--changed-files, when the resolved "
             "changed-file set has no `.md` under --root, exit 0 "
             "without scanning. Default behaviour is the same; this "
             "flag is accepted for explicitness in CI configs.")
    ap.add_argument("--list-changed-files", action="store_true",
        help="Diff-mode audit: after resolving the changed-file set "
             "from --diff-base or --changed-files, print one row per "
             "considered path showing whether the linter picked it "
             "up and, if not, why. Status values are `matched`, "
             "`ignored-extension`, `ignored-out-of-root`, "
             "`ignored-missing`, and `ignored-deleted`. The table is "
             "always written to STDERR so STDOUT remains a clean, "
             "single-document JSON payload (when --json is set) or "
             "the usual human summary. With --json the audit is "
             "serialised as a JSON array (still on STDERR) using the "
             "schema `{\"path\":str, \"status\":str, \"reason\":str}` "
             "so dashboards can ingest it without scraping the text "
             "table. No-op outside diff mode.")
    ap.add_argument("--list-changed-files-verbose", action="store_true",
        help="With --list-changed-files, expose the per-row intake "
             "provenance for every `ignored-deleted` entry. Surfaces "
             "two pieces of information that the default audit hides: "
             "(1) the exact `reason` string is preserved verbatim "
             "(it already is in the default output, but verbose mode "
             "promises the wording is *machine-stable* — no future "
             "re-wording — for `ignored-deleted` rows specifically), "
             "and (2) a new `source` field carrying the raw "
             "provenance tag — `diff-D` (a true `D`-status row from "
             "`git diff --name-status`) or `changed-files-D` (an "
             "authored `--changed-files` payload row shaped exactly "
             "`D\\tpath`). Useful when a CI pipeline needs to attribute "
             "deletes back to their intake source without scraping "
             "the reason text. Surfaces:\n"
             "  • Text mode — appends a `source` column at the end of "
             "the row; non-`ignored-deleted` rows render `-` to keep "
             "the column aligned.\n"
             "  • JSON mode — adds a `\"source\": str|null` key to "
             "EVERY row (null for non-deleted rows) so the schema is "
             "regular for downstream consumers; the legacy 3-key "
             "schema is preserved byte-for-byte when the flag is "
             "off.\n"
             "Composes with --with-similarity (similarity columns "
             "render first, then `source` last), --dedupe-changed-"
             "files (first-seen `source` wins), and --only-changed-"
             "status (filter runs after source attachment). No-op "
             "without --list-changed-files.")
    ap.add_argument("--dedupe-changed-files", action="store_true",
        help="With --list-changed-files, collapse repeated `path` "
             "values in the audit so each repo-relative path appears "
             "at most once. The FIRST occurrence wins — its `status` "
             "and `reason` are preserved verbatim and later rows for "
             "the same path are dropped (no merging, no \"...and 2 "
             "more\" annotations). Useful when the upstream intake "
             "feeds the same path from multiple sources (e.g. a "
             "`--changed-files` file concatenated from several "
             "`paths-filter` jobs, or a rebase that touched the same "
             "file in two commits and the diff surfaced it twice). "
             "The dedupe footer reports how many duplicates were "
             "collapsed so the count is auditable. No-op without "
             "--list-changed-files; the underlying linted-file set "
             "is unaffected (it's already a `set[Path]`, so "
             "duplicates can never reach the scan).")
    ap.add_argument("--only-changed-status", action="append",
        default=None, choices=list(_AUDIT_STATUSES), metavar="STATUS",
        help="With --list-changed-files, restrict the printed audit "
             "to rows whose `status` is in this set (repeatable). "
             "Valid values match the closed audit vocabulary: "
             "`matched`, `ignored-extension`, `ignored-out-of-root`, "
             "`ignored-missing`, `ignored-deleted`. Filtering happens "
             "AFTER `--dedupe-changed-files` collapses duplicates, so "
             "first-seen semantics are evaluated against the full "
             "intake (a `matched` row that lost the dedupe race to an "
             "earlier `ignored-extension` for the same path stays "
             "hidden, same as without the filter). The footer's "
             "totals line still counts EVERY status in the canonical "
             "order so you can see what was filtered out. Common "
             "uses: `--only-changed-status matched` to feed a downstream "
             "tool the exact list the linter scanned, or "
             "`--only-changed-status ignored-extension` to debug why a "
             "PR's docs aren't being checked. No-op without "
             "--list-changed-files.")
    ap.add_argument("--with-similarity", action="store_true",
        help="With --list-changed-files, include the rename/copy "
             "similarity metadata in the printed audit. Three extra "
             "columns are appended: `kind` (`R` for rename, `C` for "
             "copy, `-` for plain A/M/D rows), `score` (git's 0–100 "
             "similarity percentage, `-` when absent — e.g. plain "
             "rows or arrow-form `--changed-files` payloads that "
             "don't carry a percentage, `?` is reserved for future "
             "use), and `old` (the OLD-side path on R/C rows, `-` "
             "otherwise). With --json the metadata is emitted as a "
             "nested object `{\"kind\":str, \"score\":int|null, "
             "\"old_path\":str}` under the `similarity` key (or "
             "`null` when no rename/copy was observed) so dashboards "
             "can reason about provenance without parsing the text "
             "table. No-op without --list-changed-files; safe to "
             "combine with --dedupe-changed-files (first-seen "
             "semantics also apply to the similarity record) and "
             "--only-changed-status (filtering runs after the "
             "metadata is attached).")
    ap.add_argument("--similarity-csv", default=None, metavar="PATH",
        help="With --list-changed-files, ALSO export the audit rows "
             "as CSV to PATH for spreadsheet review (Excel, Numbers, "
             "LibreOffice, `csvkit`, etc.). Use `-` to write the CSV "
             "to STDOUT instead — only safe when STDOUT isn't already "
             "carrying the violation summary or `--json` payload, so "
             "the recommended pattern is a real file path. The header "
             "row is always `path,status,reason,kind,score,old_path` "
             "(stable column order regardless of `--with-similarity`); "
             "the four similarity columns are populated when the "
             "underlying row carries a `_RenameSimilarity` and left "
             "EMPTY otherwise — empty `score` cells distinguish "
             "*unscored* rename/copy rows (authored `--changed-files` "
             "payloads without a percentage) from `score=0` (git rated "
             "the pair entirely dissimilar). Plain A/M/D rows have "
             "all four similarity cells empty. RFC 4180 quoting via "
             "the stdlib `csv` module so paths with commas, quotes, "
             "or newlines round-trip safely. Dedupe + filter run "
             "BEFORE the export, so the CSV contains exactly the rows "
             "you'd see in the text/JSON audit. No-op without "
             "--list-changed-files.")
    ap.add_argument("--similarity-csv-format", default="csv",
        choices=("csv", "tsv"),
        help="Field-separator dialect for the --similarity-csv export. "
             "`csv` (default) writes RFC 4180 with comma separators "
             "and double-quote quoting (the legacy behaviour, "
             "byte-for-byte unchanged when the flag is absent). `tsv` "
             "writes tab-separated values using the stdlib "
             "`csv.excel_tab` dialect — handy when your paths or "
             "reasons contain commas (avoids needing quoted cells in "
             "the spreadsheet) or when piping into tools that prefer "
             "tabs (`cut -f`, `awk -F'\\t'`, `column -t -s$'\\t'`, "
             "`q -t`). The header row, column order, and the empty-"
             "vs-`0` score-cell convention are identical across both "
             "dialects — only the separator changes. The output file "
             "extension is NOT auto-rewritten: pass an explicit "
             "`.tsv` path if you want one. No-op without "
             "--similarity-csv.")
    ap.add_argument("--similarity-labels", action="store_true",
        help="With --with-similarity, attach a per-kind discriminator "
             "to every rename/copy row so the score's *meaning* is "
             "explicit instead of implied by the kind letter. Three "
             "canonical labels: `rename-similarity` (R rows — score "
             "is how alike the two paths are, 100 = byte-identical "
             "move), `copy-similarity` (C rows — score is how much of "
             "the source survived in the copy, 100 = verbatim "
             "duplicate), and `unscored` (R/C row whose percentage is "
             "absent — kind is still meaningful, magnitude isn't). "
             "Plain A/M/D rows carry no label. In the text table the "
             "label appears as a new `meaning` column appended after "
             "`old`; in --json mode it's added as a `score_kind` "
             "field on the nested `similarity` object (omitted on "
             "plain rows where `similarity` itself is null); in "
             "--similarity-csv exports a 7th `score_kind` column is "
             "APPENDED (never inserted) so positional readers that "
             "hard-code indices 0–5 keep working unchanged. Opt-in "
             "to preserve the legacy schema byte-for-byte for "
             "downstream consumers; no-op without --with-similarity.")
    ap.add_argument("--similarity-legend",
        choices=("auto", "on", "off"), default="auto",
        help="Control whether the audit table is followed by a short "
             "human-readable legend explaining the similarity columns "
             "(`kind`, `score`, `old`, and — when --similarity-labels "
             "is on — `meaning`). Three modes: `auto` (default) emits "
             "the legend only when the audit stream (STDERR) is "
             "attached to an interactive terminal, so a human reading "
             "the run live gets the cheat-sheet but a CI log capture / "
             "pipe / file redirect stays byte-for-byte identical to "
             "the legacy output (no surprise prose for log scrapers); "
             "`on` forces the legend regardless of TTY (useful when "
             "rendering to a paged-pretty wrapper that strips TTY "
             "detection but a human is still reading); `off` "
             "suppresses the legend even on a live terminal (useful "
             "when streaming straight into clipboard / paste-into-"
             "ticket flows where the prose is noise). The legend is "
             "always printed AFTER the totals footer so consumers "
             "parsing column-aligned rows see the same byte sequence "
             "they always have up to the totals line. No-op without "
             "--with-similarity (the columns being explained simply "
             "aren't there) and no-op in --json mode (machine "
             "consumers don't need prose).")
    ap.add_argument("--github", dest="github", action="store_true",
        default=None,
        help="Emit one GitHub Actions `::error file=…,line=…,title=…::` "
             "workflow command per violation in addition to the human-"
             "readable summary, so each finding lights up inline on "
             "the PR diff with its rule code (P-001 … P-008). Auto-"
             "enabled when the `GITHUB_ACTIONS` env var is `true` "
             "(i.e. inside any GitHub Actions runner). Use "
             "`--no-github` to force-disable.")
    ap.add_argument("--no-github", dest="github", action="store_false",
        help="Disable GitHub Actions annotations even when the "
             "`GITHUB_ACTIONS` env var would auto-enable them.")
    ap.add_argument("--diff-context", type=int, default=3, metavar="N",
        help="When --diff-base is set, fetch N lines of unified-diff "
             "context around each violation (via `git diff -UN <base> "
             "-- <file>`) and print the post-state excerpt under the "
             "human-readable summary so authors can patch without "
             "switching to git. Default 3; 0 disables. Ignored in "
             "--changed-files mode (no diff-base to query). In --json "
             "mode, excerpts are suppressed by default to keep the "
             "schema byte-identical for legacy consumers; pass "
             "`--json-excerpts` to emit them as a structured array.")
    ap.add_argument("--json-excerpts", action="store_true",
        help="Only meaningful with --json + --diff-base. Adds an "
             "`excerpt` array to each violation row containing the "
             "post-state diff window around the violation line. Each "
             "element is `{\"line\": int, \"kind\": \"+\"|\" \", "
             "\"text\": str, \"focus\": bool}` — `focus` marks the "
             "exact violation line so consumers don't need to re-do "
             "the centering math. The schema is additive: violations "
             "with no available excerpt simply omit the key, so "
             "parsers that don't know about it are unaffected. "
             "Window size is governed by --diff-context.")
    ap.add_argument("--suggest-patch", action="store_true",
        help="Under each human-readable diff excerpt, append a "
             "`git apply`-style unified-diff scaffold that removes "
             "the offending line and inserts a rule-specific TODO "
             "marker in its place. Designed for copy-paste: pipe the "
             "block (between the `--- BEGIN SUGGESTED PATCH ---` and "
             "`--- END SUGGESTED PATCH ---` fences) into "
             "`git apply -p0` and the file is staged with a clearly "
             "marked spot to fix. The replacement text is a TODO "
             "hint, not a real fix — the linter cannot infer the "
             "author's intent. No-op without --diff-base (no post-"
             "state line numbers available outside diff mode).")
    ap.add_argument("--json-suggest-patch", action="store_true",
        help="Only meaningful with --json + --diff-base. Adds a "
             "`suggested_patch` string field to each violation row "
             "containing the same `git apply`-ready unified diff as "
             "--suggest-patch. The schema is strictly additive: "
             "violations the linter can't generate a patch for "
             "(e.g. line not in any captured hunk) simply omit the "
             "key, so legacy parsers keying off `file`/`line`/`code`/"
             "`message` keep working unchanged.")
    args = ap.parse_args(argv)

    root = Path(args.root).resolve()
    repo_root = Path(args.repo_root).resolve()
    if not root.is_dir():
        print(f"error: --root {args.root!r} is not a directory", file=sys.stderr)
        return 2

    if args.diff_base and args.changed_files:
        print("error: --diff-base and --changed-files are mutually exclusive",
              file=sys.stderr)
        return 2

    # ``--diff-prev`` is a shorthand for ``--diff-base HEAD~N`` (default
    # N=1). It's mutually exclusive with both --diff-base and
    # --changed-files: stacking it with --diff-base would silently lose
    # one of the two refs, and stacking it with --changed-files would
    # mix two distinct intake sources. We resolve --diff-prev into
    # ``args.diff_base`` so every downstream code path keeps using a
    # single attribute.
    if args.diff_prev is not None:
        if args.diff_base:
            print("error: --diff-prev and --diff-base are mutually exclusive",
                  file=sys.stderr)
            return 2
        if args.changed_files:
            print("error: --diff-prev and --changed-files are mutually exclusive",
                  file=sys.stderr)
            return 2
        prev_raw = str(args.diff_prev).strip()
        if not prev_raw.isdigit() or int(prev_raw) < 0:
            print(f"error: --diff-prev requires a non-negative integer "
                  f"(got {args.diff_prev!r})", file=sys.stderr)
            return 2
        args.diff_base = f"HEAD~{int(prev_raw)}"
    elif args.diff_base:
        # Apply numeric / ~N / ^N shorthand expansion to whatever the
        # user typed. Non-shorthand refs pass through unchanged, so a
        # plain ``--diff-base origin/main`` still hits git verbatim.
        args.diff_base = _normalize_diff_base(args.diff_base)

    if args.diff_context < 0:
        print(f"error: --diff-context must be >= 0 (got {args.diff_context})",
              file=sys.stderr)
        return 2

    # Tri-state: --github → True, --no-github → False, neither → auto.
    if args.github is None:
        github_annotations = os.environ.get("GITHUB_ACTIONS", "").lower() == "true"
    else:
        github_annotations = args.github

    intent_verbs = DEFAULT_INTENT_VERBS | {v.lower() for v in args.allow_verb}

    # ---- Resolve --extension allowlist ---------------------------
    # ``--extension`` is repeatable; ``None`` (no flag passed) keeps
    # the historical ``("md",)`` behaviour. We normalise to lowercase,
    # strip any leading dot the user typed by accident, and dedupe via
    # ``dict.fromkeys`` so the FIRST occurrence wins (preserves the
    # CLI order in error messages without affecting the cache segment,
    # which sorts independently). The result is a tuple so it can
    # flow through ``iter_markdown_files`` and the cache key as a
    # hashable, append-safe value.
    if args.extension is None:
        # No explicit --extension flags → start from the historical
        # default. ``--include-mdx`` may augment this below; without
        # it we behave exactly like the legacy ``.md``-only baseline.
        cleaned: list[str] = list(DEFAULT_EXTENSIONS)
    else:
        cleaned = [e.lstrip(".").lower() for e in args.extension if e.strip()]
        if not cleaned:
            print("error: --extension requires at least one non-empty value",
                  file=sys.stderr)
            return 2
    # ``--include-mdx`` is a convenience union, NOT a replacement: it
    # adds ``mdx`` to whatever allowlist the previous block built so
    # the user keeps their baseline (default ``md`` or any explicit
    # ``--extension`` set) AND picks up ``.mdx`` in one short flag.
    # Implemented as a list-append + dedupe so the resulting tuple
    # is order-stable for diagnostics, while ``_cache_segment`` and
    # ``_compute_cache_key`` both sort independently — that means
    # ``--include-mdx`` and ``--extension mdx`` (with default ``md``
    # baseline) collapse to the same canonical segment + cache key.
    if args.include_mdx and "mdx" not in cleaned:
        cleaned.append("mdx")
    # ``--include-txt`` mirrors ``--include-mdx`` exactly: it's a
    # convenience union (NOT a replacement) appending ``txt`` to
    # whatever baseline the previous block built. Order in
    # ``cleaned`` reflects CLI order (``--include-mdx`` before
    # ``--include-txt`` → ``(..., mdx, txt)``) which only affects
    # diagnostic output; ``_cache_segment`` and ``_compute_cache_key``
    # sort independently so ``--include-mdx --include-txt`` and
    # ``--include-txt --include-mdx`` produce byte-identical
    # sentinels and the segment collapses to ``ext-md+mdx+txt/``.
    if args.include_txt and "txt" not in cleaned:
        cleaned.append("txt")
    extensions = tuple(dict.fromkeys(cleaned))

    # ---- Resolve diff-mode changed-file allowlist (if any) -------
    # ``changed_md`` is None ⇒ full-tree mode (legacy behaviour).
    # ``changed_md`` is a set of resolved Paths ⇒ diff mode: only
    # those files emit per-file violations. Cross-file P-007 still
    # walks every `.md` so a changed file colliding with an
    # unchanged target is reported.
    changed_md: set[Path] | None = None
    changed_audit: list[_ChangedFileAudit] | None = (
        [] if args.list_changed_files else None
    )
    if args.diff_base or args.changed_files:
        try:
            changed_md = _resolve_changed_md(
                repo_root, root,
                diff_base=args.diff_base,
                changed_files=args.changed_files,
                extensions=extensions,
                audit=changed_audit,
            )
        except RuntimeError as e:
            print(f"error: {e}", file=sys.stderr)
            return 2
        # Emit the audit trail BEFORE the fast-empty-PASS branch so
        # the user always sees why their diff resolved to zero linted
        # files (e.g. the only changes were deletes, or all hits
        # were filtered as out-of-root). Goes to STDERR so STDOUT is
        # untouched in --json mode.
        if changed_audit is not None:
            _render_changed_files_audit(
                changed_audit, sys.stderr,
                as_json=args.json,
                dedupe=args.dedupe_changed_files,
                only_statuses=(frozenset(args.only_changed_status)
                               if args.only_changed_status else None),
                with_similarity=args.with_similarity,
                with_labels=args.similarity_labels,
                legend_mode=args.similarity_legend,
                verbose=args.list_changed_files_verbose,
            )
            # CSV export mirrors the same dedupe + filter pipeline so
            # the spreadsheet always matches what the operator just
            # saw on STDERR — no surprise extra rows. Independently
            # computed (rather than threaded through the renderer) to
            # keep the renderer's signature stable and to allow CSV
            # without --with-similarity (the export carries all four
            # similarity columns regardless; they just stay empty).
            if args.similarity_csv:
                csv_rows = changed_audit
                if args.dedupe_changed_files:
                    csv_rows, _ = _dedupe_audit_rows(csv_rows)
                if args.only_changed_status:
                    only = frozenset(args.only_changed_status)
                    csv_rows = [r for r in csv_rows if r.status in only]
                _write_similarity_csv(
                    csv_rows, args.similarity_csv,
                    with_labels=args.similarity_labels,
                    dialect=args.similarity_csv_format,
                )
        if not args.json:
            print(f"ℹ️  placeholder-comments: diff-mode active — "
                  f"{len(changed_md)} changed `.md` file(s) under {args.root}/")
        if not changed_md:
            # Nothing under --root changed → fast PASS. Cross-file P-007
            # has nothing new to report by definition (no new bullets).
            if args.json:
                print("[]")
            else:
                print(f"✅ placeholder-comments: no spec changes vs. diff base, "
                      "skipping scan.")
            return 0

    # ---- Cache fast-path ------------------------------------------
    # The cache key fingerprints every input that can change the
    # linter's verdict: the linter script itself, the resolved root,
    # the imperative-verb allowlist (P-001 widening), and every `.md`
    # under the root. A hit short-circuits the scan; a miss falls
    # through to the full lint and writes a sentinel only if the
    # scan ends clean (exit 0).
    #
    # Diff mode bypasses the PASS-cache: the cache is keyed on the
    # full-tree fingerprint, but a diff-mode run only inspects a
    # subset of files, so its PASS verdict is *narrower* and must
    # never satisfy a future full-tree query. Skipping cache I/O
    # entirely keeps the invariant trivial: only full-tree PASSes
    # are ever cached.
    cache_key: str | None = None
    sentinel: Path | None = None
    if args.cache_dir and changed_md is None:
        cache_key = _compute_cache_key(root, intent_verbs,
                                       extensions=extensions)
        # Sentinels live under a per-extension-allowlist subdirectory
        # so different ``--extension`` runs are physically segregated
        # on disk. The segment name is deterministic from the sorted
        # extensions: ``ext-md``, ``ext-md+mdx``, etc. Long or
        # otherwise filesystem-hostile allowlists fall back to a
        # short-hash form (see ``_cache_segment``) to keep the path
        # legal on Windows + tar-friendly. ``mkdir(parents=True)``
        # below creates the segment directory on first PASS.
        segment = _cache_segment(extensions)
        sentinel = Path(args.cache_dir) / segment / f"{cache_key}.pass"
        if sentinel.is_file():
            if not args.json:
                print(f"✅ placeholder-comments: cache hit "
                      f"({segment}/{cache_key[:12]}…), "
                      f"skipping scan of {args.root}/")
            else:
                print("[]")
            return 0

    violations: list[Violation] = []
    cross_file_bullets: list[tuple[str, int, str]] = []
    for md in iter_markdown_files(root, extensions=extensions):
        if changed_md is not None and md.resolve() not in changed_md:
            # Unchanged file: still collect its bullets so cross-file
            # P-007 can detect a new collision introduced by a
            # changed file, but suppress its per-file violations.
            _collect_bullets_only(md, repo_root, cross_file_bullets)
            continue
        violations.extend(lint_file(md, repo_root, cross_file_bullets, intent_verbs))

    # ---- P-007 cross-file duplicates -------------------------------
    # Group every valid bullet across the scan by canonical target.
    # Within-file duplicates are already reported above, so we only
    # surface groups whose entries span ≥2 distinct files. The
    # *second* and later occurrences are flagged, pointing back at
    # the first declaration site for fast triage.
    #
    # In diff mode, only collisions whose *later* side lives in a
    # changed file are reported — an unchanged file colliding with
    # another unchanged file is pre-existing and out of scope for
    # the push under review.
    by_target: dict[str, list[tuple[str, int, str]]] = {}
    for rel, ln, target in cross_file_bullets:
        by_target.setdefault(_canonical_target(rel, target, repo_root), []).append(
            (rel, ln, target)
        )
    changed_rels: set[str] | None = None
    if changed_md is not None:
        changed_rels = {
            str(p.relative_to(repo_root)) for p in changed_md
            if p.is_relative_to(repo_root)
        }
    for key, entries in by_target.items():
        files_seen = {e[0] for e in entries}
        if len(files_seen) < 2:
            continue
        first_rel, first_ln, first_target = entries[0]
        for rel, ln, target in entries[1:]:
            if rel == first_rel:
                continue  # already reported by the within-file pass
            if changed_rels is not None and rel not in changed_rels:
                continue  # pre-existing collision in unchanged file
            violations.append(Violation(rel, ln, "P-007",
                f"Duplicate placeholder target `{target}` — also declared at "
                f"`{first_rel}:L{first_ln}` as `{first_target}` "
                "(anchor differences are ignored)."))

    # ---- Pre-fetch diff excerpts once for both human + JSON modes
    # (only when the user actually wants them). Excerpts are bounded
    # by changed-file count, not violation count — a single bad
    # block tripping P-001 + P-002 + P-004 still costs one git
    # invocation per file, not three.
    diff_excerpts: dict[str, _DiffExcerpts] = {}
    want_excerpts = (
        violations
        and args.diff_base
        and args.diff_context > 0
        and changed_md is not None
        and (not args.json or args.json_excerpts)
    )
    # Suggested patches reuse the same ``_DiffExcerpts`` data (post-
    # state line index + hunk windows) so we widen the fetch trigger
    # to also cover --suggest-patch / --json-suggest-patch when
    # excerpts themselves are off (e.g. --diff-context=0 but the
    # author still wants a copy-paste fix scaffold).
    want_patches = (
        violations
        and args.diff_base
        and changed_md is not None
        and ((not args.json and args.suggest_patch)
             or (args.json and args.json_suggest_patch))
    )
    if want_excerpts or want_patches:
        affected = sorted({v.file for v in violations
                           if (repo_root / v.file).resolve() in changed_md})
        for rel in affected:
            # When only patches are requested, fetch a minimal
            # 1-line context so suggest_patch() still has the
            # above/below anchor rows it uses for hunk math. The
            # excerpt renderer would emit a tiny window in that
            # case, but we already gate human/JSON excerpt output
            # on ``want_excerpts`` separately so nothing leaks.
            ctx = args.diff_context if want_excerpts else max(
                1, args.diff_context,
            )
            excerpt = _fetch_diff_excerpts(
                repo_root, args.diff_base, rel, ctx,
            )
            if excerpt is not None:
                diff_excerpts[rel] = excerpt

    if args.json:
        # Backward-compatible: when --json-excerpts is OFF the
        # payload is byte-identical to the legacy schema (only the
        # four Violation fields). When ON, an ``excerpt`` array is
        # appended to violations that have one — never present as
        # ``null`` or ``[]``, so legacy parsers that key only off
        # ``file``/``line``/``code``/``message`` keep working
        # without seeing a new always-present field.
        if not args.json_excerpts and not args.json_suggest_patch:
            print(json.dumps([asdict(v) for v in violations], indent=2))
        else:
            payload: list[dict[str, object]] = []
            for v in violations:
                row = asdict(v)
                excerpt = diff_excerpts.get(v.file)
                if excerpt is not None and args.json_excerpts:
                    snippet = excerpt.render_structured(
                        v.line, args.diff_context,
                    )
                    if snippet:
                        row["excerpt"] = snippet
                if excerpt is not None and args.json_suggest_patch:
                    patch_text = excerpt.suggest_patch(v.file, v.line, v.code)
                    if patch_text:
                        row["suggested_patch"] = patch_text
                payload.append(row)
            # ``ensure_ascii=False`` so non-ASCII spec content
            # (e.g. quoted UTF-8 paths from the rename hardening)
            # round-trips as readable characters instead of
            # ``\uXXXX`` escapes. Still valid JSON; downstream
            # parsers don't care.
            print(json.dumps(payload, indent=2, ensure_ascii=False))
    else:
        if not violations:
            print(f"✅ placeholder-comments: no malformed blocks under {args.root}/")
        else:
            print(f"❌ placeholder-comments: {len(violations)} violation(s):\n")
            for v in violations:
                print(f"  {v.file}:{v.line}  [{v.code}] {v.message}")
                excerpt = diff_excerpts.get(v.file)
                if excerpt is not None and args.diff_context > 0:
                    snippet = excerpt.render(v.line, args.diff_context)
                    if snippet:
                        # Two-space indent under the violation line
                        # so the excerpt is visually attached to it
                        # in the log without breaking grep on the
                        # leading `<file>:<line>` shape.
                        for sline in snippet:
                            print(f"    {sline}")
                if excerpt is not None and args.suggest_patch:
                    patch_text = excerpt.suggest_patch(v.file, v.line, v.code)
                    if patch_text:
                        # Wrap the patch in clear fences so authors
                        # can mouse-select the body and pipe it
                        # straight into ``git apply -p0 --recount``.
                        # Indented to nest visually under the
                        # violation, but the fences themselves stay
                        # at column 4 so awk/sed extraction is
                        # trivial (look for the literal sentinel).
                        print("    --- BEGIN SUGGESTED PATCH "
                              "(git apply -p0 --recount) ---")
                        for pline in patch_text.rstrip("\n").split("\n"):
                            print(f"    {pline}")
                        print("    --- END SUGGESTED PATCH ---")
            print("\n  See linter-scripts/check-placeholder-comments.py for rule docs.")

    # ---- GitHub Actions annotations (always after the human summary)
    # so a reviewer scrolling the log sees the digest first, then the
    # auto-attached inline annotations on the PR diff. Workflow
    # commands go to STDOUT regardless of --json so JSON consumers
    # still get clean output on a separate channel (the annotations
    # stream is interleaved but parseable by the runner, not by us).
    if github_annotations and violations:
        for line in _format_github_annotations(violations):
            print(line)

    # ---- Persist sentinel on clean runs only ----------------------
    # Failed runs MUST NOT poison the cache: a future "fix" might
    # re-introduce the same hash via revert, and we'd then skip the
    # scan and miss the regression. Only PASS gets cached.
    if sentinel is not None and not violations and not args.no_cache_write:
        try:
            sentinel.parent.mkdir(parents=True, exist_ok=True)
            sentinel.write_text(
                f"placeholder-comments PASS\nkey={cache_key}\n",
                encoding="utf-8",
            )
        except OSError as e:
            # Cache write failures are advisory — never fail the run.
            print(f"::warning::placeholder-comments: cache write failed: {e}",
                  file=sys.stderr)

    return 1 if violations else 0


def _dedupe_audit_rows(
    rows: list[_ChangedFileAudit],
) -> tuple[list[_ChangedFileAudit], int]:
    """Collapse repeated ``path`` values, keeping the FIRST row.

    Returns ``(deduped_rows, dropped_count)``. Order of the surviving
    rows matches their first-occurrence order in ``rows`` so the
    rendered table stays stable + reviewable. The dropped rows'
    ``status``/``reason`` are *not* merged into the survivor — first-
    seen wins, full stop. This matches the documented contract on
    ``--dedupe-changed-files`` and keeps the function trivially
    idempotent (running it twice returns the same list with
    ``dropped_count == 0`` on the second pass).
    """
    seen: set[str] = set()
    out: list[_ChangedFileAudit] = []
    dropped = 0
    for r in rows:
        if r.path in seen:
            dropped += 1
            continue
        seen.add(r.path)
        out.append(r)
    return out, dropped


# Sentinel substituted for blank/unknown similarity cells in the
# text-mode renderer. Single dash so the column stays narrow and the
# eye picks out "no rename here" at a glance. JSON consumers see real
# ``null`` instead.
_SIMILARITY_BLANK = "-"


# Per-source ``reason`` text for ``ignored-deleted`` rows. The KEY is
# the provenance tag emitted by :func:`_parse_name_status` /
# :func:`_normalise_changed_lines` alongside each captured delete; the
# VALUE is the human-readable explanation that lands on the audit
# row. Centralised so the four callers (text table, JSON payload,
# CSV export, dedupe footer) all surface the same wording for a
# given source.
#
# Vocabulary is intentionally tiny + closed:
#   * ``diff-D``           — true ``D``-status row from
#                            ``git diff --name-status``.
#   * ``changed-files-D``  — authored ``--changed-files`` payload
#                            shaped exactly ``D\tpath`` (the verbatim
#                            git wire format some CI runners
#                            forward).
#
# ``_DELETED_REASON_FALLBACK`` covers any future provenance the
# parsers add before this map catches up — keeps the audit
# self-explanatory rather than crashing on a missing key.
_DELETED_REASON: dict[str, str] = {
    "diff-D": ("git diff reported D (deleted): file removed in the "
               "diff range, no post-state to lint"),
    "changed-files-D": ("--changed-files payload row shaped `D\\tpath`: "
                        "explicit delete marker, no post-state to lint"),
}
_DELETED_REASON_FALLBACK = ("path captured as a delete by the diff "
                            "intake but provenance is unknown — "
                            "treated as `ignored-deleted` for safety")


# Stable header for the ``--similarity-csv`` export. Frozen at module
# scope so the column order is part of the contract — downstream
# spreadsheets / pandas readers can hard-code positions if they want.
_SIMILARITY_CSV_HEADER: tuple[str, ...] = (
    "path", "status", "reason", "kind", "score", "old_path",
)

# Extended header used when ``--similarity-labels`` is on. The extra
# trailing column is appended (rather than inserted) so consumers that
# read positionally can keep using indices 0–5 unchanged and only need
# to opt into index 6 when they care about the per-kind label.
_SIMILARITY_CSV_HEADER_LABELED: tuple[str, ...] = (
    *_SIMILARITY_CSV_HEADER, "score_kind",
)

# Recognised values for ``--similarity-csv-format``. Centralised so
# the argparse ``choices=`` list, the writer's dispatch, and the
# tests all agree on the exact spelling. ``csv`` is the legacy
# default (RFC 4180 with comma separators); ``tsv`` switches to
# the stdlib ``csv.excel_tab`` dialect (tab separators) for
# spreadsheets / pipelines where commas are inconvenient.
_SIMILARITY_CSV_FORMAT_CSV = "csv"
_SIMILARITY_CSV_FORMAT_TSV = "tsv"
_SIMILARITY_CSV_FORMATS: tuple[str, ...] = (
    _SIMILARITY_CSV_FORMAT_CSV,
    _SIMILARITY_CSV_FORMAT_TSV,
)

# Canonical labels for the ``score_kind`` discriminator. Centralised so
# the JSON serializer, the text-table renderer, the CSV exporter, and
# the tests all agree on the exact spelling. The vocabulary is
# deliberately tiny and hyphenated so downstream grep / jq pipelines
# can pattern-match without ambiguity.
#
# Semantics:
#   * ``rename-similarity`` — kind ``R``, score = how similar the two
#     paths are (100 = byte-identical move, lower = more edits during
#     the rename).
#   * ``copy-similarity``   — kind ``C``, score = how much of the
#     source survived in the copy (100 = verbatim duplicate).
#   * ``unscored``          — kind ``R`` / ``C`` row whose score is
#     ``None`` (authored ``--changed-files`` payload that omitted the
#     percentage, or arrow-form rename). The kind is still meaningful;
#     the magnitude isn't.
_SCORE_KIND_RENAME = "rename-similarity"
_SCORE_KIND_COPY = "copy-similarity"
_SCORE_KIND_UNSCORED = "unscored"


def _score_kind_for(sim: "_RenameSimilarity | None") -> str | None:
    """Map a ``_RenameSimilarity`` to its canonical ``score_kind`` label.

    Returns ``None`` for plain A/M/D rows (no rename provenance at
    all) so callers can distinguish "no label applies" from "label is
    ``unscored``" — the latter still carries a kind letter and an
    old-side path, only the magnitude is missing.
    """
    if sim is None:
        return None
    if sim.score is None:
        return _SCORE_KIND_UNSCORED
    if sim.kind == "R":
        return _SCORE_KIND_RENAME
    if sim.kind == "C":
        return _SCORE_KIND_COPY
    # Defensive: an unknown kind letter shouldn't reach here (the
    # parsers only emit R/C), but if a future git format adds one we
    # surface it as unscored rather than crashing.
    return _SCORE_KIND_UNSCORED


def _write_similarity_csv(rows: list[_ChangedFileAudit],
                          target: str,
                          *,
                          with_labels: bool = False,
                          dialect: str = _SIMILARITY_CSV_FORMAT_CSV) -> None:
    """Export the audit rows as RFC 4180 CSV for spreadsheet review.

    ``target`` is either a filesystem path or the literal ``"-"`` to
    write to STDOUT. The header is always
    ``path,status,reason,kind,score,old_path`` regardless of whether
    ``--with-similarity`` was passed — the four similarity columns
    just stay empty when no ``_RenameSimilarity`` is attached.

    When ``with_labels=True`` (driven by ``--similarity-labels``) a
    seventh ``score_kind`` column is APPENDED — never inserted — so
    positional readers that already hard-code indices 0–5 keep
    working unchanged and only need to opt into index 6 when they
    care about the per-kind discriminator. The cell vocabulary is
    ``rename-similarity`` / ``copy-similarity`` / ``unscored`` (R or
    C rows) and empty for plain A/M/D rows, mirroring the JSON
    ``score_kind`` field exactly.

    Empty `score` cells are *intentional* and meaningful: they mark
    *unscored* rename/copy rows (authored ``--changed-files`` payloads
    that omitted the percentage) and distinguish them from
    ``score=0`` (git observed the pair and rated them entirely
    dissimilar). Keep that distinction when filtering in Excel —
    ``ISBLANK`` vs ``=0`` are not the same condition.

    Quoting is delegated to the stdlib ``csv`` writer with the default
    dialect, so paths containing commas, quotes, or newlines are
    escaped per RFC 4180 and round-trip through every mainstream CSV
    reader. ``newline=""`` on the file handle is mandatory per the
    ``csv`` module docs to avoid stray blank lines on Windows.

    ``dialect`` selects the field separator and quoting rules:

    * ``"csv"`` (default) — stdlib default dialect: comma separator,
      double-quote quoting, RFC 4180 escapes. Byte-for-byte
      unchanged from the legacy behaviour when the caller doesn't
      pass the kwarg.
    * ``"tsv"`` — stdlib ``csv.excel_tab`` dialect: tab separator,
      same double-quote quoting rules. Use when commas inside
      paths/reasons would force quoted cells, or when piping into
      ``cut -f`` / ``awk -F'\\t'`` style tools. Tabs and newlines
      inside cell values still get quoted via the underlying
      ``csv`` writer, so the round-trip is lossless.

    Header row, column order, and the empty-vs-``0`` score
    convention are identical across both dialects — only the
    separator changes.
    """
    def _emit(handle) -> None:  # type: ignore[no-untyped-def]
        # ``excel_tab`` is the stdlib's canonical TSV dialect: tab
        # delimiter, same quoting rules as the default ``excel``
        # dialect, so cells containing tabs/newlines/quotes still
        # round-trip safely.
        if dialect == _SIMILARITY_CSV_FORMAT_TSV:
            writer = _csv.writer(handle, dialect="excel-tab")
        else:
            writer = _csv.writer(handle)
        header = (_SIMILARITY_CSV_HEADER_LABELED if with_labels
                  else _SIMILARITY_CSV_HEADER)
        writer.writerow(header)
        for r in rows:
            sim = r.similarity
            if sim is None:
                kind = score = old_path = ""
            else:
                kind = sim.kind
                # Empty cell for unscored rows; ``str(0)`` for the
                # legitimate zero-similarity case so the two stay
                # distinguishable in the spreadsheet.
                score = "" if sim.score is None else str(sim.score)
                old_path = sim.old_path
            row = [r.path, r.status, r.reason, kind, score, old_path]
            if with_labels:
                # Empty cell for plain A/M/D rows — same convention as
                # the other similarity columns — so the "no rename
                # provenance" case is uniform across all four/five
                # similarity-related fields.
                row.append(_score_kind_for(sim) or "")
            writer.writerow(row)

    if target == "-":
        _emit(sys.stdout)
    else:
        # ``newline=""`` per the csv module's documented contract; the
        # writer inserts the platform-correct line terminator itself.
        with open(target, "w", encoding="utf-8", newline="") as fh:
            _emit(fh)


# Canonical legend modes for ``--similarity-legend``. Centralised so
# the argparse ``choices=`` list, the resolver, and the tests all
# spell the vocabulary the same way. ``auto`` is the default and
# means "emit only on an interactive TTY".
_SIMILARITY_LEGEND_AUTO = "auto"
_SIMILARITY_LEGEND_ON = "on"
_SIMILARITY_LEGEND_OFF = "off"
_SIMILARITY_LEGEND_MODES: tuple[str, ...] = (
    _SIMILARITY_LEGEND_AUTO,
    _SIMILARITY_LEGEND_ON,
    _SIMILARITY_LEGEND_OFF,
)


def _should_emit_similarity_legend(
    mode: str,
    stream,  # type: ignore[no-untyped-def]
) -> bool:
    """Decide whether to print the similarity legend for ``stream``.

    * ``mode == "on"``  → always True.
    * ``mode == "off"`` → always False.
    * ``mode == "auto"`` (default) → True iff ``stream`` is attached
      to an interactive terminal. We probe via :func:`os.isatty` on
      the stream's underlying file descriptor; any failure (a
      ``StringIO`` test double, a wrapper without ``fileno``, an
      ``OSError`` because the fd was closed) is treated as
      "not a TTY" so non-interactive contexts default to silent.

    Centralised so the renderer call site stays a one-liner and the
    tests can drive every branch directly without a subprocess.
    """
    if mode == _SIMILARITY_LEGEND_ON:
        return True
    if mode == _SIMILARITY_LEGEND_OFF:
        return False
    # ``auto`` from here on. Anything we can't positively confirm as
    # a TTY is treated as non-interactive (the safe default — a CI
    # log file capture must not gain surprise prose).
    fileno = getattr(stream, "fileno", None)
    if not callable(fileno):
        return False
    try:
        return os.isatty(fileno())
    except (OSError, ValueError) as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return False


def _render_similarity_legend(stream,  # type: ignore[no-untyped-def]
                              *,
                              with_labels: bool) -> None:
    """Print a compact legend for the similarity columns to ``stream``.

    Always printed AFTER the totals footer so consumers parsing the
    column-aligned audit rows + totals see the same byte sequence
    they always have up to that point. The legend itself is
    indentation-prefixed (``  ``) so it visually nests under the
    table header banner — same indent convention as the existing
    ``totals: …`` footer.

    When ``with_labels`` is True an extra line documents the
    ``meaning`` column added by ``--similarity-labels``.
    """
    print("  legend:", file=stream)
    print("    kind   R = rename, C = copy, - = plain A/M/D row",
          file=stream)
    print("    score  git's 0–100 similarity %, - if absent "
          "(authored payload without %, or plain row)",
          file=stream)
    print("    old    OLD-side path on R/C rows, - otherwise",
          file=stream)
    if with_labels:
        print("    meaning  rename-similarity / copy-similarity / "
              "unscored, - on plain rows", file=stream)


def _fmt_similarity(
    sim: "_RenameSimilarity | None",
) -> tuple[str, str, str]:
    """Stringify a similarity record for the text audit table.

    Returns ``(kind, score, old_path)`` as already-padded-friendly
    strings (no internal whitespace, only column-width-friendly tokens
    so :meth:`str.ljust` stays predictable). Each cell falls back to
    ``_SIMILARITY_BLANK`` independently so a partial record (R/C row
    with no numeric score, or an arrow-form rename whose old path
    survived but score didn't) is still maximally informative.

    The cell vocabulary is intentionally tiny so downstream grep
    pipelines can match on it:
      * ``kind`` ∈ {``R``, ``C``, ``-``}
      * ``score`` ∈ {``0`` … ``100``, ``-``}
      * ``old`` is a path or ``-``
    """
    if sim is None:
        return (_SIMILARITY_BLANK, _SIMILARITY_BLANK, _SIMILARITY_BLANK)
    kind = sim.kind or _SIMILARITY_BLANK
    score = str(sim.score) if sim.score is not None else _SIMILARITY_BLANK
    old = sim.old_path or _SIMILARITY_BLANK
    return (kind, score, old)


def _render_changed_files_audit(rows: list[_ChangedFileAudit],
                                stream,  # type: ignore[no-untyped-def]
                                *,
                                as_json: bool,
                                dedupe: bool = False,
                                only_statuses: frozenset[str] | None = None,
                                with_similarity: bool = False,
                                with_labels: bool = False,
                                legend_mode: str = _SIMILARITY_LEGEND_AUTO,
                                verbose: bool = False,
                                ) -> None:
    """Print the diff-mode changed-file audit table to ``stream``.

    Always writes to STDERR (the caller passes ``sys.stderr``) so
    STDOUT remains a clean single-document JSON payload (or the
    usual human summary) regardless of this flag.

    Two output modes:

    * ``as_json=False`` — aligned text table with a header row, a
      counts-by-status footer, and stable input ordering. Empty
      input prints a single ``(no changed files considered)`` line
      so the operator never wonders whether the flag silently
      no-op'd.
    * ``as_json=True``  — JSON array of ``{"path", "status",
      "reason"}`` objects, one per row, in stable input order.
      ``ensure_ascii=False`` so non-ASCII paths round-trip readably.

    When ``dedupe=True`` the rows are passed through
    :func:`_dedupe_audit_rows` first: each repo-relative ``path``
    appears at most once, with the FIRST seen ``status`` + ``reason``
    preserved verbatim. The text-mode header is annotated with the
    drop count so the collapse is auditable; JSON consumers can
    diff the array length against the upstream intake size.

    When ``only_statuses`` is a non-None frozenset, only rows whose
    ``status`` is in that set are printed (or serialised). Filtering
    runs AFTER dedupe so first-seen semantics are evaluated against
    the full intake — a ``matched`` row that lost the dedupe race to
    an earlier ``ignored-extension`` for the same path stays hidden,
    same as without the filter. The text-mode header reports both
    the visible row count and the underlying total so a filter that
    hides everything is obvious (``0 of 12 row(s) shown``); the
    totals line still counts every status in the canonical order so
    the operator can see exactly what was filtered out.

    When ``with_similarity`` is True the rendered table grows three
    extra columns — ``kind``, ``score``, ``old`` — populated from each
    row's ``similarity`` field. ``None`` similarities (plain A/M rows
    + every ``ignored-deleted`` row, which carries no rename
    provenance) render as ``-`` in all three columns. R/C rows whose
    score is ``None`` (scoreless authored payloads) render the score
    cell as ``-`` while still showing the kind and old-path. JSON
    mode emits a nested ``similarity`` object (or ``null``) using
    :func:`dataclasses.asdict` so the schema is regular for downstream
    consumers. Off by default to keep the legacy 3-column shape and
    avoid widening logs that don't care about provenance.

    When ``with_labels`` is True (driven by ``--similarity-labels`` and
    only meaningful in combination with ``with_similarity``) one more
    column — ``meaning`` — is appended to the text table, carrying the
    canonical ``rename-similarity`` / ``copy-similarity`` / ``unscored``
    discriminator (or ``-`` for plain A/M/D rows). In JSON mode the
    same value is added to the nested ``similarity`` object as a
    ``score_kind`` field; for plain rows where ``similarity`` itself
    is ``null`` the discriminator is naturally absent. The legacy
    schema is preserved byte-for-byte when the flag is off.

    When ``verbose`` is True (driven by ``--list-changed-files-verbose``)
    the audit exposes the per-row intake provenance for every
    ``ignored-deleted`` row:

    * **Text mode** appends a trailing ``source`` column whose
      cell is the raw provenance tag (``diff-D`` /
      ``changed-files-D``) on ``ignored-deleted`` rows and the
      blank-cell sentinel (``-``) on every other row.
    * **JSON mode** adds a top-level ``"source"`` key to EVERY
      row in the array — ``str`` on ``ignored-deleted`` rows and
      ``null`` everywhere else — so the schema stays regular and
      downstream JSON consumers can ``.get("source")`` without
      branching on status.

    Off by default (``verbose=False``) the ``source`` field is
    stripped entirely from the JSON payload (legacy 3-key schema
    preserved byte-for-byte) and the text table renders without
    the trailing column.
    """
    dropped = 0
    if dedupe:
        rows, dropped = _dedupe_audit_rows(rows)
    full_rows = rows
    if only_statuses is not None:
        rows = [r for r in rows if r.status in only_statuses]
    if as_json:
        # ``asdict`` recurses into nested dataclasses so the
        # ``similarity`` field becomes a sub-object automatically. When
        # the operator didn't ask for similarity, drop the field
        # entirely from the payload so the legacy schema is preserved
        # byte-for-byte. (Stripping a key is intentional rather than
        # emitting ``"similarity": null``: dashboards parsing the
        # historical schema with a strict object validator must keep
        # working unchanged.)
        payload = []
        for r in rows:
            obj = asdict(r)
            if not with_similarity:
                obj.pop("similarity", None)
            elif with_labels:
                # Inject the discriminator alongside the existing
                # ``kind``/``score``/``old_path`` triple. Skip the
                # injection on plain rows whose ``similarity`` is
                # ``null`` — there's no sub-object to extend, and
                # absence already means "no provenance".
                sim_obj = obj.get("similarity")
                if isinstance(sim_obj, dict):
                    sim_obj["score_kind"] = _score_kind_for(r.similarity)
            if not verbose:
                # Legacy schema preservation: strip the new ``source``
                # key entirely when the operator didn't opt in. Same
                # rationale as the ``similarity`` strip above —
                # downstream validators that close on the historical
                # 3-key shape (``path``/``status``/``reason``) keep
                # working unchanged. Absent ≠ null: callers MUST
                # distinguish "field omitted" (legacy mode) from
                # "field present but null" (verbose mode + non-
                # deleted row).
                obj.pop("source", None)
            payload.append(obj)
        print(json.dumps(payload, indent=2, ensure_ascii=False),
              file=stream)
        return

    suffix = (f"; deduped, {dropped} duplicate(s) dropped"
              if dedupe else "")
    if only_statuses is not None:
        # Surface the filter in the header so a hidden-everything
        # filter doesn't look like a bug. The totals line below still
        # reports the full breakdown.
        suffix += (f"; filtered, {len(rows)} of {len(full_rows)} "
                   f"row(s) shown ({'+'.join(sorted(only_statuses))})")
    if with_similarity:
        # Surface the extra columns in the header so a reviewer
        # scanning the log knows the wider table isn't a layout bug.
        suffix += "; +similarity columns"
        if with_labels:
            suffix += " +meaning"
    if verbose:
        # Mention the verbose `source` column in the header banner
        # so a reviewer scanning the log notices the schema change
        # without comparing against a non-verbose run.
        suffix += "; +source"
    print("── placeholder-comments: changed-file audit "
          f"({len(full_rows)} row(s){suffix}) ──", file=stream)
    if not rows:
        if not full_rows:
            print("  (no changed files considered)", file=stream)
        else:
            # The intake had rows but none matched the filter — make
            # that explicit so the operator can adjust their query.
            print("  (no rows matched --only-changed-status)",
                  file=stream)
        # Still print the totals footer so the filtered-out counts
        # are visible. Falls through to the counting block below.
        path_w = status_w = 0
    else:
        path_w = max(len("path"), max(len(r.path) for r in rows))
        status_w = max(len("status"), max(len(r.status) for r in rows))
    if rows:
        # Pre-compute the verbose ``source`` column when needed so
        # its width participates in the same ``ljust`` math as
        # every other fixed-width cell. ``-`` is the blank-cell
        # sentinel — same convention as the similarity columns —
        # so non-``ignored-deleted`` rows keep the column aligned
        # without leaking a stray ``None`` literal into the table.
        if verbose:
            sources = [(r.source if r.source else "-") for r in rows]
            source_w = max(len("source"), max(len(s) for s in sources))
        else:
            sources = []
            source_w = 0
        if with_similarity:
            # Pre-compute every cell so the column widths bake in the
            # widest dash-substituted value. ``_fmt_similarity`` returns
            # the (kind, score, old) triple as already-stringified
            # cells (with `-` substituted for None / non-rename rows).
            cells = [_fmt_similarity(r.similarity) for r in rows]
            kind_w = max(len("kind"), max(len(c[0]) for c in cells))
            score_w = max(len("score"), max(len(c[1]) for c in cells))
            old_w = max(len("old"), max(len(c[2]) for c in cells))
            if with_labels:
                # Compute the meaning column up front so its width
                # participates in the same ``ljust`` math as the rest.
                # Plain rows render as ``-`` to match the surrounding
                # blank-cell convention.
                meanings = [
                    (_score_kind_for(r.similarity) or _SIMILARITY_BLANK)
                    for r in rows
                ]
                meaning_w = max(len("meaning"),
                                max(len(m) for m in meanings))
                header = (
                    f"  {'status'.ljust(status_w)}  "
                    f"{'path'.ljust(path_w)}  "
                    f"{'kind'.ljust(kind_w)}  "
                    f"{'score'.ljust(score_w)}  "
                    f"{'old'.ljust(old_w)}  "
                    f"{'meaning'.ljust(meaning_w)}  "
                    + (f"{'source'.ljust(source_w)}  " if verbose else "")
                    + "reason"
                )
                print(header, file=stream)
                rule_w = (status_w + path_w + kind_w + score_w
                          + old_w + meaning_w + len("reason") + 6 * 2)
                if verbose:
                    rule_w += source_w + 2
                print("  " + "-" * rule_w, file=stream)
                for i, (r, (k, sc, op), meaning) in enumerate(
                        zip(rows, cells, meanings)):
                    src_cell = (f"{sources[i].ljust(source_w)}  "
                                if verbose else "")
                    print(f"  {r.status.ljust(status_w)}  "
                          f"{r.path.ljust(path_w)}  "
                          f"{k.ljust(kind_w)}  "
                          f"{sc.ljust(score_w)}  "
                          f"{op.ljust(old_w)}  "
                          f"{meaning.ljust(meaning_w)}  "
                          f"{src_cell}{r.reason}",
                          file=stream)
            else:
                header = (
                    f"  {'status'.ljust(status_w)}  "
                    f"{'path'.ljust(path_w)}  "
                    f"{'kind'.ljust(kind_w)}  "
                    f"{'score'.ljust(score_w)}  "
                    f"{'old'.ljust(old_w)}  "
                    + (f"{'source'.ljust(source_w)}  " if verbose else "")
                    + "reason"
                )
                print(header, file=stream)
                rule_w = (status_w + path_w + kind_w + score_w
                          + old_w + len("reason") + 5 * 2)
                if verbose:
                    rule_w += source_w + 2
                print("  " + "-" * rule_w, file=stream)
                for i, (r, (k, sc, op)) in enumerate(zip(rows, cells)):
                    src_cell = (f"{sources[i].ljust(source_w)}  "
                                if verbose else "")
                    print(f"  {r.status.ljust(status_w)}  "
                          f"{r.path.ljust(path_w)}  "
                          f"{k.ljust(kind_w)}  "
                          f"{sc.ljust(score_w)}  "
                          f"{op.ljust(old_w)}  "
                          f"{src_cell}{r.reason}", file=stream)
        else:
            header = (
                f"  {'status'.ljust(status_w)}  "
                f"{'path'.ljust(path_w)}  "
                + (f"{'source'.ljust(source_w)}  " if verbose else "")
                + "reason"
            )
            print(header, file=stream)
            rule_w = status_w + path_w + len("reason") + 4
            if verbose:
                rule_w += source_w + 2
            print("  " + "-" * rule_w, file=stream)
            for i, r in enumerate(rows):
                src_cell = (f"{sources[i].ljust(source_w)}  "
                            if verbose else "")
                print(f"  {r.status.ljust(status_w)}  "
                      f"{r.path.ljust(path_w)}  "
                      f"{src_cell}{r.reason}", file=stream)
    # Counts-by-status footer in the canonical status order so the
    # eye lands on the same column positions run-to-run. Counts
    # against ``full_rows`` (post-dedupe, pre-filter) so the totals
    # line truthfully describes the underlying intake even when a
    # filter is hiding most of it.
    counts = {s: 0 for s in _AUDIT_STATUSES}
    for r in full_rows:
        counts[r.status] = counts.get(r.status, 0) + 1
    summary = "  ".join(f"{s}={counts[s]}" for s in _AUDIT_STATUSES)
    print(f"  totals: {summary}", file=stream)
    # Optional legend — only meaningful when the similarity columns
    # are actually in the table. Resolver decides on/off given the
    # operator's ``--similarity-legend`` choice and the stream's TTY
    # status; suppressed in JSON mode (handled by the early return
    # above, which never reaches this footer block).
    if with_similarity and _should_emit_similarity_legend(
            legend_mode, stream):
        _render_similarity_legend(stream, with_labels=with_labels)


def _normalize_diff_base(ref: str) -> str:
    """Expand bare-numeric / ``~N`` / ``^N`` shorthands to ``HEAD~N`` / ``HEAD^N``.

    Rules (applied to the trimmed input):
      * ``"N"`` where N is a positive int  → ``"HEAD~N"``
      * ``"~N"`` (N a positive int)        → ``"HEAD~N"``
      * ``"^N"`` (N a positive int)        → ``"HEAD^N"``
      * Anything else (including refs that already start with a name
        like ``HEAD``, ``main``, ``origin/main``, a SHA, etc.) is
        returned verbatim so we never second-guess git's own ref
        grammar. ``"0"`` / ``"~0"`` resolve to ``HEAD`` itself which
        git accepts and which yields an empty diff — useful for
        smoke-testing the diff plumbing without changing scope.
    """
    s = ref.strip()
    if not s:
        return ref
    if s.isdigit():
        return f"HEAD~{s}"
    if len(s) >= 2 and s[0] in ("~", "^") and s[1:].isdigit():
        return f"HEAD{s}"
    return ref


def _resolve_changed_md(repo_root: Path, root: Path, *,
                        diff_base: str | None,
                        changed_files: str | None,
                        extensions: tuple[str, ...] = DEFAULT_EXTENSIONS,
                        audit: list[_ChangedFileAudit] | None = None,
                        ) -> set[Path]:
    """Resolve the set of `.md` files under ``root`` that are changed.

    Two input modes:
      * ``diff_base`` → invoke
        ``git diff --name-status -M -C --diff-filter=AMRCD
        <base>...HEAD`` from ``repo_root``. The triple-dot syntax
        compares HEAD against the merge-base with ``<base>``, which
        matches GitHub's PR diff and survives force-pushes / rebases
        on the base branch. The ``AMRCD`` filter keeps Added,
        Modified, Renamed, Copied, *and* Deleted paths. Deletes are
        never linted (a deleted file can't carry a new violation in
        the post-state) but they ARE recorded in the audit trail
        with status ``ignored-deleted`` so downstream consumers can
        see the full intake. For ``R``/``C`` rows we take the *new* path so
        the linter re-checks the file under its post-rename location
        — that's where the placeholder block lives now, and a rename
        commit often touches it (e.g. updated back-pointer hints).
        ``-M`` enables rename detection (default 50 % similarity)
        and ``-C`` enables copy detection so the new sibling of a
        copy is also linted.
      * ``changed_files`` → read newline-delimited paths from the
        given file (``-`` = stdin). Blank lines and ``#`` comments are
        ignored. Useful for CI runners that compute the diff
        themselves (e.g. ``dorny/paths-filter``) or for local testing
        without a git invocation. Renames may be expressed on a
        single line as either ``OLD\\tNEW`` (tab-separated, matches
        ``git diff --name-status`` output verbatim) or
        ``OLD => NEW`` (matches ``git status -s`` rename arrows). In
        both forms the OLD path is discarded and the NEW path is
        linted as a normal change. A ``D\\tpath`` row in this input
        is recorded as ``ignored-deleted`` in the audit trail.

    ``audit`` (optional out-parameter): when provided, every path the
    intake considered is appended as a :class:`_ChangedFileAudit` row
    classified by extension/root/existence/delete-status. Pass
    ``None`` (the default) to skip the bookkeeping entirely — the
    legacy hot path stays allocation-free.

    Returned paths are absolute + resolved and filtered to:
      * extension ``.md``
      * residing under ``root`` (so a README change doesn't trigger a
        spec scan)
      * actually present on disk (a Modified path that was reverted
        in a later commit of the same push won't exist)
    """
    # Each entry is the post-state repo-relative path. Rename/copy
    # rows contribute only their NEW side. Delete rows are tagged
    # with a leading ``\x00D\x00`` sentinel so the audit pass can
    # mark them ``ignored-deleted`` without re-parsing the diff.
    raw: list[str] = []
    # Each entry is ``(path, source)`` where ``source`` is one of
    # ``_DELETED_REASON``'s keys. The audit emitter below maps
    # ``source`` → human-readable ``reason`` so each
    # ``ignored-deleted`` row explains *why* it was classified.
    deleted_paths: list[tuple[str, str]] = []
    # When the caller asked for an audit trail, also collect rename/
    # copy provenance per new-path so the audit constructor below can
    # attach a ``_RenameSimilarity`` to every row whose path came from
    # an R/C diff entry. The map is keyed on the *unquoted* new-path
    # string exactly as it appears in ``raw`` so the lookup is a
    # constant-time dict hit per row.
    similarities: dict[str, _RenameSimilarity] | None = (
        {} if audit is not None else None
    )
    if diff_base:
        try:
            proc = subprocess.run(
                ["git", "diff", "--name-status", "-M", "-C",
                 "--diff-filter=AMRCD",
                 f"{diff_base}...HEAD"],
                cwd=repo_root, check=True, capture_output=True, text=True,
            )
        except FileNotFoundError as e:
            import sys as _sys; print(f"Error: {e}", file=_sys.stderr)
            raise RuntimeError(f"git not found on PATH: {e}") from e
        except subprocess.CalledProcessError as e:
            import sys as _sys; print(f"Error: {e}", file=_sys.stderr)
            raise RuntimeError(
                f"git diff vs. {diff_base!r} failed (exit {e.returncode}): "
                f"{e.stderr.strip() or '(no stderr)'}"
            ) from e
        raw = _parse_name_status(proc.stdout, deleted=deleted_paths,
                                 similarities=similarities)
    else:
        assert changed_files is not None
        if changed_files == "-":
            lines = sys.stdin.read().splitlines()
        else:
            try:
                lines = Path(changed_files).read_text(encoding="utf-8").splitlines()
            except OSError as e:
                import sys as _sys; print(f"Error: {e}", file=_sys.stderr)
                raise RuntimeError(
                    f"--changed-files {changed_files!r} unreadable: {e}"
                ) from e
        raw = _normalise_changed_lines(lines, deleted=deleted_paths,
                                       similarities=similarities)
    allowed_exts = {("." + e.lstrip(".").lower()) for e in extensions}
    out: set[Path] = set()
    for line in raw:
        s = line.strip()
        if not s or s.startswith("#"):
            continue
        # Pull the rename/copy provenance (if any) for this path. We
        # look up against the unstripped ``line`` first (which is what
        # ``similarities`` was keyed on) and fall back to the stripped
        # form for symmetry with how ``out.add`` resolves the path.
        # ``None`` means "no rename signal observed" — the renderer
        # will substitute dashes when --with-similarity is on.
        sim = None
        if similarities is not None:
            sim = similarities.get(line) or similarities.get(s)
        ext = Path(s).suffix.lower()
        if ext not in allowed_exts:
            if audit is not None:
                audit.append(_ChangedFileAudit(
                    path=s, status="ignored-extension",
                    reason=(f"extension {ext or '(none)'!r} not in "
                            f"allowlist {sorted(allowed_exts)}"),
                    similarity=sim,
                ))
            continue
        p = (repo_root / s).resolve()
        try:
            p.relative_to(root)
        except ValueError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            if audit is not None:
                audit.append(_ChangedFileAudit(
                    path=s, status="ignored-out-of-root",
                    reason=f"path is outside --root {root}",
                    similarity=sim,
                ))
            continue
        if not p.is_file():
            if audit is not None:
                audit.append(_ChangedFileAudit(
                    path=s, status="ignored-missing",
                    reason="post-state path is not on disk "
                           "(reverted later in the push, or "
                           "filtered by .gitignore on checkout)",
                    similarity=sim,
                ))
            continue
        out.add(p)
        if audit is not None:
            audit.append(_ChangedFileAudit(
                path=s, status="matched",
                reason="under --root, extension allowed, "
                       "file present on disk",
                similarity=sim,
            ))
    if audit is not None:
        for d, src in deleted_paths:
            audit.append(_ChangedFileAudit(
                path=d, status="ignored-deleted",
                reason=_DELETED_REASON.get(
                    src, _DELETED_REASON_FALLBACK),
                # Carry the raw provenance tag onto the audit row so
                # ``--list-changed-files-verbose`` can surface it
                # alongside the human-readable reason. Non-deleted
                # rows leave ``source`` at its ``None`` default.
                source=src,
            ))
    return out


# `git diff --name-status -M` emits one of:
#   A\tpath          (added)
#   M\tpath          (modified)
#   D\tpath          (deleted)              — pre-filtered, never seen here
#   R<score>\told\tnew  (renamed, score 0–100)
#   C<score>\told\tnew  (copied,  score 0–100)
#   T\tpath          (type change)          — pre-filtered, never seen here
# Tabs are the field separator; we split on tab and key off the first
# character of column 0 to decide which column carries the new path.
_NAME_STATUS_RE = re.compile(r"^([AMDRCTUX])(\d{0,3})$")


# Git emits paths in C-quoted form (``"path\twith\ttab"``) when
# ``core.quotePath`` is true (the default) and the path contains a
# byte that isn't safe for the terminal — tabs, newlines, control
# chars, non-ASCII bytes when ``core.quotePath=true``. The quoting
# is a strict subset of C string escapes: surrounding double quotes,
# backslash escapes for ``\a \b \t \n \v \f \r " \\``, and
# ``\NNN`` octal triplets for arbitrary bytes. A path that doesn't
# need escaping is output bare (no quotes). We MUST decode quoted
# paths before splitting on tab — otherwise an embedded ``\t`` in a
# valid filename would be mistaken for a column separator.
#
# Reference: ``git help config`` → ``core.quotePath``;
# ``quote.c::quote_c_style_counted`` in git's source.
# Match a *run* of consecutive ``\NNN`` octal escapes so we can
# decode them as a single UTF-8 byte sequence. Decoding triplet-by-
# triplet would split a multi-byte character (e.g. ``é`` =
# ``\303\251``) across two ``bytes.decode`` calls and each half
# would emit a U+FFFD replacement char. A single ``re.sub`` over
# the whole run hands the bytes to the codec atomically, which is
# the only way to get correct round-trips for non-ASCII paths.
_C_OCT_RUN_RE = re.compile(r"(?:\\[0-7]{1,3})+")
_C_ESC_TBL = {
    "a": "\a", "b": "\b", "t": "\t", "n": "\n",
    "v": "\v", "f": "\f", "r": "\r",
    '"': '"', "\\": "\\",
}


def _unquote_git_path(field: str) -> str:
    """Reverse git's C-style path quoting if ``field`` is wrapped in
    double quotes; otherwise return ``field`` unchanged.

    Tolerant of malformed input: a stray ``\\x`` (where ``x`` is not
    a recognised escape) is passed through verbatim rather than
    raising — this matches how a human would copy-paste the row out
    of ``git status`` and into ``--changed-files``. Also tolerant of
    a trailing CR (Windows line endings) which can survive
    ``splitlines()`` when the file is opened in binary or has lone
    ``\\r`` separators upstream.
    """
    s = field
    # Strip a single trailing CR — harmless on POSIX paths (NUL is
    # the only forbidden byte besides ``/``) and silently fixes
    # Windows-runner inputs.
    if s.endswith("\r"):
        s = s[:-1]
    if not (len(s) >= 2 and s.startswith('"') and s.endswith('"')):
        return s
    inner = s[1:-1]
    # First expand ``\NNN`` octal byte escapes. We decode each
    # *run* of escapes as one UTF-8 byte string so a multi-byte
    # character split across triplets (``\303\251`` = ``é``)
    # round-trips correctly. ``errors="replace"`` keeps malformed
    # input visible (U+FFFD) rather than raising — same posture as
    # the rest of the linter, which never crashes on weird git
    # output, only logs the violation site.
    def _oct_run_sub(m: "re.Match[str]") -> str:
        run = m.group(0)
        try:
            buf = bytes(int(t, 8) for t in run.split("\\")[1:])
        except ValueError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            return run
        return buf.decode("utf-8", "replace")
    inner = _C_OCT_RUN_RE.sub(_oct_run_sub, inner)
    # Then expand single-char escapes.
    out: list[str] = []
    i = 0
    while i < len(inner):
        ch = inner[i]
        if ch == "\\" and i + 1 < len(inner):
            esc = inner[i + 1]
            out.append(_C_ESC_TBL.get(esc, "\\" + esc))
            i += 2
            continue
        out.append(ch)
        i += 1
    return "".join(out)


def _parse_name_status(stdout: str,
                       *,
                       deleted: "list[tuple[str, str]] | None" = None,
                       similarities: "dict[str, _RenameSimilarity] | None" = None,
                       ) -> list[str]:
    """Extract the post-state path from each ``git diff --name-status``
    row, mapping renames + copies to their NEW side.

    Unknown / malformed rows are skipped silently — the linter's job
    is to lint placeholders, not to police git plumbing output.

    When ``deleted`` is provided, every ``D``-status row contributes
    a ``(path, source)`` tuple — in input order, after
    :func:`_unquote_git_path`. The ``source`` is always ``"diff-D"``
    for this parser (rows come straight from
    ``git diff --name-status``); the audit-trail emitter uses the
    tag to look up the per-provenance ``reason`` string in
    :data:`_DELETED_REASON` so each ``ignored-deleted`` row explains
    *why* it was classified that way.

    When ``similarities`` is provided, every ``R``/``C`` row contributes
    one ``new_path → _RenameSimilarity`` entry. The mapping key is the
    *unquoted* new path (so it matches what ``raw`` carries downstream)
    and the value records ``kind`` (``R``/``C``), ``score`` (0–100, or
    ``None`` when git emitted no digits — pathological but cheap to
    tolerate), and the unquoted ``old_path``. Plain A/M/D/T rows are
    not recorded — the renderer treats their absence as "no similarity
    metadata" and prints the dash sentinel.

    Hardened against git's path-quoting and whitespace edge cases:

    * Tabs are the field separator. C-quoted paths
      (``"with\\ttab.md"``) are decoded *before* the tab split would
      see them, so an embedded literal tab inside a filename can't
      masquerade as a column separator.
    * Trailing CR on the row (CRLF input from a Windows-piped diff)
      is stripped per-field by :func:`_unquote_git_path`.
    * The R/C arm requires a non-empty *new* path (``cols[2]``) but
      tolerates an empty *old* path slot — git never emits one, but
      hand-rolled diff payloads occasionally do, and there's no
      reason to discard the row when its NEW side is well-formed.
    * Whitespace-only paths (``"   "``) are kept as-is — POSIX
      permits them, and ``Path.is_file()`` downstream will resolve
      whether the file actually exists.
    """
    out: list[str] = []
    for line in stdout.splitlines():
        if not line:
            continue
        # Drop a trailing CR on the *row* before column splitting so
        # the last field doesn't get a stray ``\r`` glued onto it.
        # Per-field stripping handles the in-quote case; this handles
        # the bare (unquoted) case for the row's last path.
        if line.endswith("\r"):
            line = line[:-1]
        cols = line.split("\t")
        if len(cols) < 2:
            continue
        m = _NAME_STATUS_RE.match(cols[0])
        if not m:
            continue
        kind = m.group(1)
        if kind in ("R", "C"):
            # Rename / copy: cols = [R<score>, old, new]. Take new.
            # ``cols[2]`` is required; ``cols[1]`` (old) may be empty
            # in pathological inputs — we don't need it for linting.
            if len(cols) >= 3 and cols[2] != "":
                new_path = _unquote_git_path(cols[2])
                out.append(new_path)
                if similarities is not None:
                    score_raw = m.group(2)
                    score = int(score_raw) if score_raw else None
                    similarities[new_path] = _RenameSimilarity(
                        kind=kind,
                        score=score,
                        old_path=_unquote_git_path(cols[1]) if cols[1] else "",
                    )
        elif kind in ("A", "M"):
            # Add / modify: cols = [A|M, path]. Take path.
            if cols[1] != "":
                out.append(_unquote_git_path(cols[1]))
        elif kind == "D" and deleted is not None:
            # Delete: cols = [D, path]. Path is captured for the
            # audit trail only — never returned for linting because
            # there is no post-state file to scan.
            if cols[1] != "":
                deleted.append((_unquote_git_path(cols[1]), "diff-D"))
        # T / U / X intentionally dropped — see docstring.
    return out


# Two textual rename conventions accepted in `--changed-files` input:
#   1. Tab-separated, matches `git diff --name-status` verbatim:
#        R087\tspec/old.md\tspec/new.md
#      Any leading status token is tolerated (we just take the last
#      tab-separated column as the new path).
#   2. Arrow-separated, matches `git status -s` short output:
#        spec/old.md => spec/new.md
#      Whitespace around the arrow is ignored.
# The arrow form is intentionally permissive on the surrounding
# whitespace because it's authored by humans (or by `git status -s`
# which left-pads the row with two status columns + a space).
# ``\S`` at each end was too strict — it rejected paths that
# legitimately start or end with a space (rare but POSIX-legal). We
# now anchor on ``=>`` and let the path bodies be any non-empty
# trimmed run; trimming is done *after* the split so embedded
# whitespace inside the path is preserved.
_RENAME_ARROW_RE = re.compile(r"^\s*(?P<old>.+?)\s*=>\s*(?P<new>.+?)\s*$")


def _normalise_changed_lines(lines: list[str],
                             *,
                             deleted: "list[tuple[str, str]] | None" = None,
                             similarities: "dict[str, _RenameSimilarity] | None" = None,
                             ) -> list[str]:
    """Collapse rename-bearing rows in a ``--changed-files`` payload
    down to their post-rename path.

    When ``deleted`` is provided and a row is recognisable as a
    delete (``D\\tpath`` — the exact shape ``git diff --name-status``
    emits), the path is captured into ``deleted`` as a
    ``(path, "changed-files-D")`` tuple and the row is NOT forwarded
    to the caller. The provenance tag distinguishes it from true
    diff ``D`` rows so the audit emitter can surface a different
    ``reason`` string per source. Without ``deleted`` (the default),
    such a row falls through to the generic tab-form branch and the
    bare path travels downstream to be filtered by extension/root
    checks — same end result as before this audit-trail addition.

    Plain paths (no tab, no ``=>``) pass through unchanged. Comments
    and blanks are *not* stripped here — the caller does that on the
    normalised output so we don't lose alignment with the source line
    numbers in error messages.

    Hardened against the same whitespace + quoting edge cases as
    :func:`_parse_name_status`:

    * Tab rows: instead of dropping every empty column (which
      silently re-indexes ``R\\t\\told\\tnew`` to ``[R, old, new]``
      and then ``cols[-1]`` is correct, but ``R<score>\\told\\t\\t``
      would re-index to ``[R<score>, old]`` and steal ``old`` as the
      "new" path), we keep the column count intact and pick the
      last non-empty field. Quoted fields are unquoted; trailing
      CR is stripped.
    * Arrow rows: the regex no longer requires ``\\S`` at the path
      boundaries, so a path with a leading/trailing space round-
      trips correctly. The ``new`` group is unquoted to match what
      a user pasted from ``git status``.
    * A line that contains *only* whitespace (or a CR-only line on
      Windows input) is passed through verbatim so the caller's
      blank/comment filter can still discard it on the same line
      number.

    When ``similarities`` is provided, R/C-shaped tab rows AND arrow
    rows contribute one ``new_path → _RenameSimilarity`` entry. Tab
    rows whose first column matches ``R<digits>?`` / ``C<digits>?``
    record both the kind letter and (when present) the numeric score;
    a kind letter without digits records ``score=None``. Arrow rows
    (``OLD => NEW``) are recorded as ``kind="R", score=None`` because
    git status doesn't emit a similarity percentage in short form —
    we know it's a rename, but not how close. Plain paths (no
    rename signal at all) are not recorded.
    """
    out: list[str] = []
    for line in lines:
        # Strip a trailing CR for the whole row before any other
        # parsing. We don't ``rstrip()`` — that would eat legitimate
        # trailing spaces in a path. Only ``\r`` is dropped.
        if line.endswith("\r"):
            line = line[:-1]
        # Tab form: take the last tab-separated column. Works for
        # both `R<score>\told\tnew` (3 cols) and unscored `R\told\tnew`
        # (rare, e.g. when authors hand-edit the file).
        if "\t" in line:
            # Preserve column positions: split without filtering, so
            # padding tabs from copy-pasted output (e.g. an extra
            # ``\t`` after the ``R<score>`` token in some tooling)
            # don't shift our column index. Then pick the last
            # *non-empty* field as the post-rename path.
            cols = line.split("\t")
            # Recognise the exact ``D\tpath`` delete shape so we can
            # divert it into the audit trail instead of letting the
            # path pretend it was modified. Only fires when the
            # caller asked for delete capture; otherwise behaviour is
            # byte-identical to the historical implementation.
            if (deleted is not None
                    and len(cols) == 2
                    and cols[0] == "D"
                    and cols[1] != ""):
                deleted.append((_unquote_git_path(cols[1]),
                                "changed-files-D"))
                continue
            new_col = ""
            for c in reversed(cols):
                if c != "":
                    new_col = c
                    break
            if new_col:
                new_path = _unquote_git_path(new_col)
                out.append(new_path)
                # Recover similarity metadata when the leading column
                # is an R/C marker. We tolerate the same scoreless and
                # scored shapes ``_NAME_STATUS_RE`` accepts; anything
                # else (a bare tab-padded path, e.g. ``\tspec/x.md``)
                # is left without a similarity record so the renderer
                # falls back to dashes.
                if similarities is not None and cols and cols[0]:
                    head = _NAME_STATUS_RE.match(cols[0])
                    if head and head.group(1) in ("R", "C"):
                        # Find the OLD-side path: it's the
                        # *first* non-empty column after cols[0],
                        # excluding the new path we just consumed.
                        old_path = ""
                        for c in cols[1:]:
                            if c != "" and _unquote_git_path(c) != new_path:
                                old_path = _unquote_git_path(c)
                                break
                        score_raw = head.group(2)
                        similarities[new_path] = _RenameSimilarity(
                            kind=head.group(1),
                            score=int(score_raw) if score_raw else None,
                            old_path=old_path,
                        )
            continue
        # Arrow form: `OLD => NEW`.
        m = _RENAME_ARROW_RE.match(line)
        if m:
            new_path_raw = m.group("new")
            # Trim at boundaries (regex already did greedy-min) but
            # don't touch interior whitespace. Then unquote in case
            # the user pasted a C-quoted form from ``git status``.
            unquoted_new = _unquote_git_path(new_path_raw.strip())
            out.append(unquoted_new)
            if similarities is not None:
                old_path_raw = m.group("old")
                similarities[unquoted_new] = _RenameSimilarity(
                    kind="R",
                    score=None,
                    old_path=_unquote_git_path(old_path_raw.strip()),
                )
            continue
        out.append(line)
    return out


# ---- Diff-excerpt rendering (used by the human summary in diff mode)
#
# We keep the parser tiny and tolerant: only the post-state side of a
# unified diff matters (that's where line numbers in violations come
# from), and any malformed hunk gracefully degrades to "no excerpt"
# rather than failing the run. The shape captured per file is:
#
#     {post_line_no: (kind, text)}
#
# where ``kind`` is one of:
#     "+"   line added in the post-state (highlighted in the snippet)
#     " "   context line carried over from both sides
# Removed lines (``-``) are dropped because they don't have a
# post-state line number a violation could reference. The renderer
# slices a ±context window around the violation line, falling back
# to "no diff context available" if the line isn't part of the
# fetched hunks (e.g. a P-007 collision pointing at a hunk that
# wasn't included in the unified diff).

_HUNK_HEADER_RE = re.compile(
    r"^@@\s*-\d+(?:,\d+)?\s+\+(?P<start>\d+)(?:,(?P<count>\d+))?\s*@@"
)


@dataclass(frozen=True)
class _Hunk:
    """One post-state hunk window inside a parsed unified diff.

    ``start`` / ``end`` are inclusive 1-indexed post-state line
    numbers covering every ``+`` and `` `` row this hunk emitted.
    A single-line hunk has ``start == end``. A diff with no
    post-state coverage (rare: pure-deletion file, an empty added
    file) yields zero hunks; the parent ``_DiffExcerpts.hunks``
    list is empty in that case and the renderer returns ``[]``.
    """
    start: int
    end: int

    def contains(self, line: int) -> bool:
        return self.start <= line <= self.end

    def distance_to(self, line: int) -> int:
        """Minimum line-count distance from ``line`` to this hunk
        window. Returns 0 when the line is inside the hunk."""
        if line < self.start:
            return self.start - line
        if line > self.end:
            return line - self.end
        return 0


@dataclass(frozen=True)
class _DiffExcerpts:
    """Post-state line index for one file's `git diff -UN` output.

    ``lines[N]`` is ``(kind, text)`` for post-state line ``N`` (1-
    indexed). ``min_line`` / ``max_line`` describe the union of all
    hunk windows so the renderer can detect "violation outside any
    hunk" cleanly.

    ``hunks`` lists each hunk as a discrete ``_Hunk`` range so the
    renderer can pick the best match for a violation that lives
    between hunks (e.g. a P-007 cross-file collision that points at
    a placeholder block far away from the changed lines). Without
    this list the old renderer would silently emit an empty excerpt
    when the violation fell in the *gap* between two hunks but
    inside the global ``[min_line, max_line]`` bounds — visible to
    the bounds check but not to the line-by-line lookup.
    """
    lines: dict[int, tuple[str, str]]
    min_line: int
    max_line: int
    hunks: tuple[_Hunk, ...] = ()

    def _select_hunk(self, line: int) -> _Hunk | None:
        """Pick the best hunk to use as context for ``line``.

        Selection rules, in order:

        1. If a hunk *contains* ``line``, return it (zero distance).
        2. Otherwise return the hunk with the smallest distance to
           ``line``. Ties are broken in *post-state line order* —
           the earlier hunk wins, which keeps output deterministic
           and matches how a human reads a file top-to-bottom.
        3. No hunks → ``None`` (caller emits "no excerpt").
        """
        if not self.hunks:
            return None
        # Single pass: find the minimum distance and earliest hunk
        # at that distance. ``min(..., key=...)`` over an empty
        # iterable would raise; the early ``not self.hunks`` guard
        # above prevents that.
        return min(self.hunks,
                   key=lambda h: (h.distance_to(line), h.start))

    def render(self, line: int, context: int) -> list[str]:
        """Return a list of human-readable excerpt lines centered on
        ``line`` with up to ``context`` lines on each side, or [] if
        no relevant excerpt is available.

        When the violation line falls between hunks, the renderer
        selects the *nearest* hunk (see :meth:`_select_hunk`) and
        prepends a one-line breadcrumb so the reader knows the
        excerpt is not centered on the violation line itself but on
        the closest changed region. This is far more useful than the
        old behaviour of returning an empty list, which made the
        violation line look like it had no diff context at all.
        """
        if not self.lines:
            return []
        hunk = self._select_hunk(line)
        if hunk is None:
            return ["(line not in current diff hunks — view file directly)"]
        # Window selection:
        #   * Violation INSIDE the hunk → ±context around the
        #     violation line, clamped to hunk bounds. Mirrors the
        #     classic single-hunk behaviour.
        #   * Violation OUTSIDE every hunk → render the *entire*
        #     selected hunk (capped at 2*context+1 lines so we
        #     don't blow up output for a huge nearby hunk). This
        #     is what "best matching hunk context" means: the user
        #     gets to actually see the changed region near their
        #     violation, not just a breadcrumb.
        if hunk.contains(line):
            lo = max(hunk.start, line - context)
            hi = min(hunk.end, line + context)
        else:
            cap = 2 * context + 1
            if hunk.end - hunk.start + 1 <= cap:
                lo, hi = hunk.start, hunk.end
            elif line < hunk.start:
                # Hunk is below the violation: take its leading edge.
                lo, hi = hunk.start, hunk.start + cap - 1
            else:
                # Hunk is above the violation: take its trailing edge.
                lo, hi = hunk.end - cap + 1, hunk.end
        out: list[str] = []
        if not hunk.contains(line):
            # Violation is between hunks (or outside the whole diff
            # but still within the global ±context tolerance). Tell
            # the reader we're showing the *nearest* changed region
            # so they don't assume the excerpt is centered on the
            # violation line itself.
            delta = hunk.distance_to(line)
            out.append(
                f"  ℹ︎ violation at L{line}; nearest changed hunk "
                f"@ L{hunk.start}-{hunk.end} (Δ {delta} line"
                f"{'s' if delta != 1 else ''})"
            )
        for ln in range(lo, hi + 1):
            entry = self.lines.get(ln)
            if entry is None:
                continue
            kind, text = entry
            marker = "►" if ln == line else " "
            sigil = "+" if kind == "+" else " "
            # ``ln:5`` keeps gutter widths aligned for files up to
            # 99,999 lines — well past anything we'll encounter in
            # spec/.
            out.append(f"{marker} {ln:5d} {sigil} {text}")
        return out

    def render_structured(self, line: int,
                          context: int) -> list[dict[str, object]]:
        """Return a JSON-friendly window around ``line``.

        Each element is ``{"line": <int>, "kind": "+"|" ", "text":
        <str>, "focus": <bool>, "nearest": <bool>}`` for one post-
        state line in the ±``context`` window. Pure data — no
        Unicode markers, no gutter padding, no truncation. The
        text payload is the raw post-state line content (no
        leading ``+``/`` `` sigil); the sigil is moved to the
        typed ``kind`` field so a JSON consumer doesn't have to
        strip it.

        ``focus`` is ``True`` only on the row whose ``line`` equals
        the violation line, regardless of whether that row is in
        the same hunk as the excerpt (it usually is; for a between-
        hunks violation the focus row simply doesn't appear in the
        excerpt). ``nearest`` is ``True`` on every row of the
        excerpt when the violation is *not* in this hunk — flagging
        the whole window as a fallback to the nearest changed
        region rather than the violation site itself.

        Returns ``[]`` (not a sentinel string like the human
        renderer) when:

        * no hunks were captured for this file, OR
        * no hunk is reachable as nearest (only happens when
          ``hunks`` is empty — see above).

        Returning ``[]`` rather than a "no data" object means the
        caller can simply omit the ``excerpt`` key when the list is
        empty — keeping the JSON schema strictly additive (legacy
        consumers see no new keys on violations the linter has no
        excerpt for).
        """
        if not self.lines:
            return []
        hunk = self._select_hunk(line)
        if hunk is None:
            return []
        # Same window logic as the human renderer — see render()
        # for the rationale. The ``cap`` keeps the JSON payload
        # bounded for consumers when a giant nearby hunk would
        # otherwise dump hundreds of rows per violation.
        if hunk.contains(line):
            lo = max(hunk.start, line - context)
            hi = min(hunk.end, line + context)
        else:
            cap = 2 * context + 1
            if hunk.end - hunk.start + 1 <= cap:
                lo, hi = hunk.start, hunk.end
            elif line < hunk.start:
                lo, hi = hunk.start, hunk.start + cap - 1
            else:
                lo, hi = hunk.end - cap + 1, hunk.end
        is_nearest = not hunk.contains(line)
        out: list[dict[str, object]] = []
        for ln in range(lo, hi + 1):
            entry = self.lines.get(ln)
            if entry is None:
                continue
            kind, text = entry
            out.append({
                "line": ln,
                "kind": kind,             # "+" (added) or " " (context)
                "text": text,
                "focus": ln == line,      # exact violation line
                "nearest": is_nearest,    # True ⇒ this row is from
                                          # the *nearest* hunk, not
                                          # the violation's own
            })
        return out


    def suggest_patch(self, file: str, line: int,
                      rule_code: str) -> str:
        """Return a ``git apply``-ready unified-diff scaffold that
        replaces the violation line with a TODO marker keyed off
        ``rule_code``.

        Returns ``""`` (empty string — *not* ``None``) when no
        post-state line for ``line`` was captured (violation outside
        every hunk, or pure-removal file). Callers treat ``""`` as
        "no patch available" and simply omit the suggestion.

        The emitted diff is intentionally minimal: at most one line
        of pre-existing context above and below the violation line
        (whatever the captured hunk has — the patch is shorter at
        hunk boundaries). Headers use ``a/<path>`` / ``b/<path>``
        and post-state line numbers on both sides; downstream users
        run ``git apply -p0 --recount`` (documented in the human
        renderer's fence) so any pre/post line drift is reconciled
        automatically.

        For a context line (``kind == " "``) the patch *swaps* the
        bad line for the TODO. For an added line (``kind == "+"``)
        the patch keeps the bad addition and inserts the TODO
        immediately after — the linter cannot safely delete an
        added line because the surrounding pre-state coordinates
        would shift; the author then manually removes the bad
        addition once they've written the real fix.
        """
        entry = self.lines.get(line)
        if entry is None:
            return ""
        kind, text = entry
        # Pre-existing context (one line above + below if available
        # in the captured hunk; the parser only stores ``+`` and
        # `` `` rows so anything we get back is safe to render).
        above = self.lines.get(line - 1)
        below = self.lines.get(line + 1)
        replacement = _RULE_FIX_HINTS.get(rule_code, _RULE_FIX_FALLBACK)

        body: list[str] = []
        # Anchor lines must use a leading single space — git's
        # tolerance for missing-space context lines is undefined
        # across versions; we always emit the canonical form.
        if above is not None:
            body.append(f" {above[1]}")
        if kind == "+":
            # Pure insertion: keep the bad ``+`` line, then add the
            # TODO immediately after it. The author removes the bad
            # line manually once they write the real fix.
            body.append(f" {text}")
            body.append(f"+{replacement}")
        else:
            # Context line in the post-state ⇒ it exists in the
            # pre-state too, so a swap (`-bad` / `+todo`) lands
            # cleanly with ``git apply --recount``.
            body.append(f"-{text}")
            body.append(f"+{replacement}")
        if below is not None:
            body.append(f" {below[1]}")

        # Hunk math: pre-state and post-state both span ``hunk_len``
        # lines starting at ``hunk_start`` (the row above the
        # violation, or the violation itself when it's the first
        # line of the file). For a swap, pre and post counts match.
        # For an insertion (`+` violation), post is one larger than
        # pre — git accepts mismatched counts and ``--recount``
        # would fix them anyway, but we emit the precise numbers so
        # vanilla ``git apply -p0`` works without the flag.
        hunk_start = line if above is None else line - 1
        ctx_above = 0 if above is None else 1
        ctx_below = 0 if below is None else 1
        if kind == "+":
            pre_len = ctx_above + 1 + ctx_below          # bad + ctx
            post_len = ctx_above + 2 + ctx_below         # bad + TODO + ctx
        else:
            pre_len = ctx_above + 1 + ctx_below          # bad + ctx
            post_len = ctx_above + 1 + ctx_below         # TODO + ctx

        header = (f"@@ -{hunk_start},{pre_len} "
                  f"+{hunk_start},{post_len} @@")
        return (
            f"--- a/{file}\n"
            f"+++ b/{file}\n"
            f"{header}\n"
            + "\n".join(body)
            + "\n"
        )


def _parse_unified_diff_post(stdout: str) -> _DiffExcerpts:
    """Parse `git diff -UN` output and index post-state lines only.

    File-header / index / mode lines are ignored — we already know
    which file we asked about. Hunk headers (``@@ -a,b +c,d @@``)
    reset the post-state line counter; ``+`` and `` `` rows advance
    it; ``-`` rows are skipped (no post-state coordinate).

    In addition to the flat ``lines`` index, we record each hunk
    as a discrete ``_Hunk(start, end)`` range so a multi-hunk file
    can route a violation to the *nearest* hunk's window rather
    than blindly slicing the global min/max range. A hunk that
    covers no post-state lines (e.g. pure removals — ``+`` count
    of 0 in the header) is dropped from the list because it has no
    coordinate a violation could land on.
    """
    lines: dict[int, tuple[str, str]] = {}
    cur_post = 0
    in_hunk = False
    min_line = 10**9
    max_line = 0
    hunks: list[_Hunk] = []
    cur_hunk_start: int | None = None
    cur_hunk_last: int | None = None

    def _flush_hunk() -> None:
        # Capture the in-progress hunk's [start, end] range. We
        # use the ``last`` post-line we actually wrote to (rather
        # than ``cur_post``, which has already been incremented
        # past it) so single-line hunks get start == end.
        nonlocal cur_hunk_start, cur_hunk_last
        if cur_hunk_start is not None and cur_hunk_last is not None:
            hunks.append(_Hunk(start=cur_hunk_start, end=cur_hunk_last))
        cur_hunk_start = None
        cur_hunk_last = None

    for raw in stdout.splitlines():
        if raw.startswith("@@"):
            _flush_hunk()
            m = _HUNK_HEADER_RE.match(raw)
            if not m:
                in_hunk = False
                continue
            cur_post = int(m.group("start"))
            in_hunk = True
            continue
        if not in_hunk:
            continue
        if not raw:
            # Blank line inside a hunk = a context line whose payload
            # is empty. Treat as " " (context) to keep the line
            # counter honest.
            lines[cur_post] = (" ", "")
            min_line = min(min_line, cur_post)
            max_line = max(max_line, cur_post)
            if cur_hunk_start is None:
                cur_hunk_start = cur_post
            cur_hunk_last = cur_post
            cur_post += 1
            continue
        kind = raw[0]
        body = raw[1:]
        if kind == "+":
            lines[cur_post] = ("+", body)
            min_line = min(min_line, cur_post)
            max_line = max(max_line, cur_post)
            if cur_hunk_start is None:
                cur_hunk_start = cur_post
            cur_hunk_last = cur_post
            cur_post += 1
        elif kind == " ":
            lines[cur_post] = (" ", body)
            min_line = min(min_line, cur_post)
            max_line = max(max_line, cur_post)
            if cur_hunk_start is None:
                cur_hunk_start = cur_post
            cur_hunk_last = cur_post
            cur_post += 1
        elif kind == "-":
            # Removed line — no post-state coordinate. Skip silently.
            pass
        elif kind == "\\":
            # `\ No newline at end of file` marker — ignore.
            pass
        else:
            # Unknown row inside a hunk (extra header from combined
            # diff, etc.). Be defensive: bail out of this hunk so we
            # don't desync the line counter.
            _flush_hunk()
            in_hunk = False
    # Capture the final hunk (the last ``@@`` block has no
    # successor header to trigger a flush).
    _flush_hunk()

    if max_line == 0:
        return _DiffExcerpts(lines={}, min_line=0, max_line=0, hunks=())
    return _DiffExcerpts(lines=lines, min_line=min_line,
                         max_line=max_line, hunks=tuple(hunks))


def _fetch_diff_excerpts(repo_root: Path, diff_base: str, rel_path: str,
                         context: int) -> _DiffExcerpts | None:
    """Run `git diff -U<context> <base>...HEAD -- <rel_path>` and
    return the parsed post-state excerpt, or ``None`` if git fails
    (missing binary, unreachable base, file not in diff, etc.).

    Failures are silent on purpose: the violation summary is still
    printed without an excerpt — we never want a missing snippet to
    fail the lint run, only to degrade gracefully.
    """
    try:
        proc = subprocess.run(
            ["git", "diff", f"-U{context}", f"{diff_base}...HEAD",
             "--", rel_path],
            cwd=repo_root, check=True, capture_output=True, text=True,
        )
    except (FileNotFoundError, subprocess.CalledProcessError) as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return None
    if not proc.stdout.strip():
        return None
    return _parse_unified_diff_post(proc.stdout)


def _collect_bullets_only(path: Path, repo_root: Path,
                          bullets_out: list[tuple[str, int, str]]) -> None:
    """Cross-file P-007 helper for diff mode.

    Re-uses ``lint_file`` to extract every valid bullet from an
    unchanged file but discards the per-file violations. The bullets
    are needed so a *changed* file's new bullet can collide with a
    pre-existing target in an unchanged file and still trip P-007.
    """
    # ``lint_file`` already appends to ``bullets_out`` via its
    # ``valid_bullets`` parameter; we drop the violations list.
    lint_file(path, repo_root, bullets_out, DEFAULT_INTENT_VERBS)


# Per-rule one-liner shown in the annotation title so reviewers see
# *why* the line is flagged without opening the linter docs. Keep
# each entry ≤ ~50 chars — GitHub truncates long titles in the diff
# gutter tooltip.
RULE_TITLES: dict[str, str] = {
    "P-001": "Placeholder intent must be an imperative sentence",
    "P-002": "Placeholder body must be `- [text](link)` bullets",
    "P-003": "Placeholder link must be a relative `.md` path",
    "P-004": "Placeholder block must contain ≥1 valid bullet",
    "P-005": "Placeholder block must not contain blank lines",
    "P-006": "Placeholder opener has no matching closer",
    "P-007": "Duplicate placeholder target",
    "P-008": "Placeholder opener missing `@path:line` back-pointer",
}

# GitHub Actions workflow commands use `,` and `:` as field separators
# and `\n` / `\r` as line terminators. Any of these in the message
# corrupt the annotation, so we URL-style escape them per the
# documented contract:
# https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions
_ANNOTATION_ESCAPES: tuple[tuple[str, str], ...] = (
    ("%", "%25"),  # MUST be first — every other replacement uses %.
    ("\r", "%0D"),
    ("\n", "%0A"),
    (":", "%3A"),
    (",", "%2C"),
)


def _escape_annotation(value: str) -> str:
    out = value
    for src, dst in _ANNOTATION_ESCAPES:
        out = out.replace(src, dst)
    return out


def _format_github_annotations(violations: list[Violation]) -> Iterable[str]:
    """Yield one ``::error file=…,line=…,col=1,title=…::message`` per
    violation, preserving input order.

    * ``file`` is the repo-relative path stored on ``Violation`` —
      matches GitHub's checkout layout so the gutter pin lands on the
      right file in the PR diff.
    * ``line`` is the 1-indexed source line of the offending opener.
    * ``col=1`` is included so the annotation pins to the gutter
      rather than column 0 (some renderers hide col=0 entirely).
    * ``title`` is ``"<P-NNN> <one-liner>"`` so the rule code is the
      first thing reviewers see in the diff tooltip; unknown codes
      degrade to just the bare code (forward-compatible with future
      P-009+ rules added before this map is updated).
    * The message body is the full ``Violation.message`` so the
      remediation hint stays visible when the user clicks through.
    """
    for v in violations:
        title = RULE_TITLES.get(v.code)
        head = f"{v.code} {title}" if title else v.code
        yield (
            f"::error file={_escape_annotation(v.file)},"
            f"line={v.line},col=1,"
            f"title={_escape_annotation(head)}::"
            f"{_escape_annotation(v.message)}"
        )


def _compute_cache_key(
    root: Path,
    intent_verbs: frozenset[str] | set[str],
    *,
    extensions: tuple[str, ...] = DEFAULT_EXTENSIONS,
) -> str:
    """Build a SHA-256 fingerprint of every input that affects the verdict.

    Inputs (in deterministic order):
      1. The absolute, resolved scan root.
      2. The sorted, canonicalised imperative-verb allowlist.
      3. The sorted, canonicalised extension allowlist (so a future
         ``--extension mdx`` run can never collide with the default
         ``--extension md`` set, even before the cache-segment
         directory split would catch it on disk).
      4. The SHA-256 of the linter script itself (so a logic change
         invalidates every cached PASS automatically).
      5. For every file matching ``extensions`` under the root
         (sorted by path, dotfiles excluded — same filter as
         ``iter_markdown_files``):
         ``<repo-relative-path>\\0<sha256-of-bytes>\\n``

    Anything outside this set (mtimes, permissions, sibling files,
    environment variables) is intentionally excluded so the key is
    reproducible across machines and CI shards.
    """
    h = hashlib.sha256()
    # Schema tag bumped to v2 when the extension allowlist became
    # part of the key. A v1 sentinel would have been written without
    # the ``exts=`` line, so its hash domain is disjoint from v2 —
    # old sentinels are inert (never collide with new lookups) rather
    # than dangerous, but the explicit version tag documents intent.
    h.update(b"placeholder-comments-cache-v2\n")
    h.update(f"root={root}\n".encode("utf-8"))
    h.update(("verbs=" + ",".join(sorted(intent_verbs)) + "\n").encode("utf-8"))
    h.update(("exts=" + ",".join(sorted(extensions)) + "\n").encode("utf-8"))
    try:
        script_bytes = Path(__file__).resolve().read_bytes()
        h.update(b"script=" + hashlib.sha256(script_bytes).hexdigest().encode() + b"\n")
    except OSError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        # __file__ unreadable (zipapp / frozen). Fall back to a stable
        # tag so the cache still works, just with coarser invalidation.
        h.update(b"script=unknown\n")
    for md in iter_markdown_files(root, extensions=extensions):
        try:
            data = md.read_bytes()
        except OSError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            continue
        rel = str(md.relative_to(root)).encode("utf-8")
        h.update(rel + b"\0" + hashlib.sha256(data).hexdigest().encode() + b"\n")
    return h.hexdigest()


# Filesystem-safe extension chars: lowercase ASCII letters + digits.
# Anything outside this set (dots in compound extensions like
# ``tar.gz``, unicode, slashes) forces the segment name into the
# hash-suffix form so we never produce a path that would explode on
# Windows, NTFS, or a tarball extracted on a case-insensitive FS.
_SAFE_EXT_RE = re.compile(r"^[a-z0-9]+$")

# Cap the readable form before we fall back to a hash. NTFS' 255-char
# filename limit is the binding constraint, but we want headroom for
# the surrounding ``ext-`` prefix, ``+`` joiners, AND the eventual
# ``<key>.pass`` filename inside the segment directory. 64 chars
# leaves the segment well under any practical limit while still
# accommodating ~10 typical extensions joined by ``+``.
_MAX_SEGMENT_BODY_LEN = 64


def _cache_segment(extensions: tuple[str, ...]) -> str:
    """Derive a filesystem-safe, deterministic cache subdirectory
    name from the active extension allowlist.

    Format::

        ext-<sorted-extensions-joined-by-plus>           (readable form)
        ext-h<10-char-sha256>                            (hash fallback)

    The readable form is preferred because it makes cache contents
    self-describing (``ls cache/`` shows ``ext-md/``, ``ext-md+mdx/``
    at a glance). We fall back to the hash form when ANY extension
    contains a character outside ``[a-z0-9]`` (so we never emit a
    Windows-illegal path), or when the joined name would exceed
    :data:`_MAX_SEGMENT_BODY_LEN`. The fallback is keyed on the same
    sorted-and-joined string the readable form would have used, so
    two runs with identical allowlists always land in the same
    bucket regardless of which branch they took.

    The function is pure: same input → same output, no I/O. That
    matters because the segment is consulted both on read (cache
    fast-path) and write (sentinel persistence), and any drift would
    silently bypass the cache.
    """
    # Sort to canonicalise: ``("mdx", "md")`` and ``("md", "mdx")``
    # MUST share a segment. Empty input is defensive — the CLI
    # validator rejects it earlier, but a programmatic caller might
    # not. Treat it as "default" so we still produce a stable name.
    if not extensions:
        return "ext-default"
    body = "+".join(sorted(extensions))
    if (len(body) <= _MAX_SEGMENT_BODY_LEN
            and all(_SAFE_EXT_RE.match(e) for e in extensions)):
        return f"ext-{body}"
    digest = hashlib.sha256(body.encode("utf-8")).hexdigest()[:10]
    return f"ext-h{digest}"


if __name__ == "__main__":
    sys.exit(main())