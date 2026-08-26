#!/usr/bin/env python3
"""Validate a ``rename_intake`` JSON document against the schema
documented in ``linter-scripts/README-rename-intake.md``.

This validator is **stdlib-only** so CI doesn't need an extra
``pip install jsonschema`` step. It checks the same constraints
that the formal JSON Schema (printed by ``--print-schema``) would
enforce, plus a couple of cross-field rules that are easier to
express in code than in pure schema:

* the array is non-empty (an empty array is technically valid
  JSON but always indicates the intake produced nothing — almost
  certainly a misconfigured invocation; use ``--allow-empty`` to
  permit it intentionally);
* every record has the right key set, governed by the
  ``--with-similarity`` mode (legacy 3-key schema vs enriched
  4-key schema);
* ``status`` is in the closed vocabulary;
* when present, ``similarity`` is either ``null`` or a strict
  ``{kind, score, old_path}`` object (with an optional
  ``score_kind`` when ``--with-labels`` is set);
* ``ignored-deleted`` rows always have ``similarity: null`` (the
  pre-state has no post-state to score).

Exit codes:

* ``0`` — every record is valid against the requested mode.
* ``1`` — at least one schema violation; details on STDERR.
* ``2`` — the input itself isn't valid JSON, or a CLI usage error.

CI usage is documented in ``README-rename-intake.md`` under the
*"Validating `rename_intake` output in CI"* section.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path
from typing import Any

# ---------------------------------------------------------------
# Schema constants (kept in lock-step with the renderer)
# ---------------------------------------------------------------
_STATUSES = frozenset({
    "matched", "ignored-extension", "ignored-out-of-root",
    "ignored-missing", "ignored-deleted",
})
_SIM_KINDS = frozenset({"R", "C"})
_SCORE_KINDS = frozenset({
    "rename-similarity", "copy-similarity", "unscored",
})

_LEGACY_KEYS = frozenset({"path", "status", "reason"})
_ENRICHED_KEYS = frozenset({"path", "status", "reason", "similarity"})
_SIM_OBJECT_KEYS = frozenset({"kind", "score", "old_path"})
_SIM_OBJECT_KEYS_WITH_LABEL = (
    _SIM_OBJECT_KEYS | {"score_kind"})


# ---------------------------------------------------------------
# JSON Schema (Draft 2020-12) — emitted by --print-schema for
# documentation / external tooling. The Python validator below
# re-implements the same rules to stay zero-dep.
# ---------------------------------------------------------------
def _build_json_schema(*, with_similarity: bool,
                       with_labels: bool) -> dict[str, Any]:
    sim_props: dict[str, Any] = {
        "kind": {"type": "string", "enum": sorted(_SIM_KINDS)},
        "score": {
            "oneOf": [
                {"type": "integer", "minimum": 0, "maximum": 100},
                {"type": "null"},
            ],
        },
        "old_path": {"type": "string", "minLength": 1},
    }
    sim_required = ["kind", "score", "old_path"]
    if with_labels:
        sim_props["score_kind"] = {
            "type": "string", "enum": sorted(_SCORE_KINDS),
        }
        sim_required.append("score_kind")

    record_props: dict[str, Any] = {
        "path": {"type": "string", "minLength": 1},
        "status": {"type": "string", "enum": sorted(_STATUSES)},
        "reason": {"type": "string"},
    }
    record_required = ["path", "status", "reason"]
    if with_similarity:
        record_props["similarity"] = {
            "oneOf": [
                {
                    "type": "object",
                    "properties": sim_props,
                    "required": sim_required,
                    "additionalProperties": False,
                },
                {"type": "null"},
            ],
        }
        record_required.append("similarity")

    return {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": ("https://example.invalid/rename-intake/"
                f"{'enriched' if with_similarity else 'legacy'}"
                f"{'-labels' if with_labels else ''}.json"),
        "title": ("rename_intake "
                  f"({'enriched' if with_similarity else 'legacy'}"
                  f"{' + labels' if with_labels else ''})"),
        "type": "array",
        "items": {
            "type": "object",
            "properties": record_props,
            "required": record_required,
            "additionalProperties": False,
        },
    }


# ---------------------------------------------------------------
# Validator (stdlib-only)
# ---------------------------------------------------------------
class _Errors:
    """Accumulates ``(json-pointer, message)`` errors with a cap so
    a totally garbled input doesn't produce thousands of lines."""

    def __init__(self, limit: int = 50) -> None:
        self.items: list[tuple[str, str]] = []
        self.limit = limit
        self.truncated = False

    def add(self, pointer: str, msg: str) -> None:
        if len(self.items) >= self.limit:
            self.truncated = True
            return
        self.items.append((pointer, msg))

    def __bool__(self) -> bool:
        return bool(self.items)


def _validate_similarity(
        sim: Any, ptr: str, errs: _Errors,
        *, with_labels: bool, status: str) -> None:
    if sim is None:
        return
    if not isinstance(sim, dict):
        errs.add(ptr,
            f"similarity must be an object or null, got "
            f"{type(sim).__name__}")
        return
    keys = set(sim)
    expected = (_SIM_OBJECT_KEYS_WITH_LABEL if with_labels
                else _SIM_OBJECT_KEYS)
    # When labels are enabled, score_kind is REQUIRED on every
    # non-null similarity object the renderer emits, so we treat
    # its absence as an error in --with-labels mode. When labels
    # are off, score_kind must NOT be present.
    missing = expected - keys
    extra = keys - expected
    if missing:
        errs.add(ptr,
            f"similarity object missing key(s): "
            f"{sorted(missing)}")
    if extra:
        errs.add(ptr,
            f"similarity object has unexpected key(s): "
            f"{sorted(extra)}")
    kind = sim.get("kind")
    if kind not in _SIM_KINDS:
        errs.add(ptr + "/kind",
            f"kind must be one of {sorted(_SIM_KINDS)}; got "
            f"{kind!r}")
    score = sim.get("score", "<missing>")
    if score is not None and score != "<missing>":
        if not isinstance(score, int) or isinstance(score, bool):
            errs.add(ptr + "/score",
                f"score must be int|null; got "
                f"{type(score).__name__}")
        elif not (0 <= score <= 100):
            errs.add(ptr + "/score",
                f"score must be in [0, 100]; got {score}")
    old = sim.get("old_path")
    if not isinstance(old, str) or not old:
        errs.add(ptr + "/old_path",
            f"old_path must be a non-empty string; got {old!r}")
    if with_labels:
        sk = sim.get("score_kind")
        if sk not in _SCORE_KINDS:
            errs.add(ptr + "/score_kind",
                f"score_kind must be one of "
                f"{sorted(_SCORE_KINDS)}; got {sk!r}")
    # ignored-deleted rows must never carry a non-null similarity.
    if status == "ignored-deleted":
        errs.add(ptr,
            "ignored-deleted rows must have similarity:null "
            "(no post-state to score)")


def _validate(doc: Any, *, with_similarity: bool,
              with_labels: bool, allow_empty: bool,
              ) -> _Errors:
    errs = _Errors()
    if not isinstance(doc, list):
        errs.add("",
            f"top-level value must be a JSON array; got "
            f"{type(doc).__name__}")
        return errs
    if not doc and not allow_empty:
        errs.add("",
            "array is empty; pass --allow-empty if this is "
            "intentional (e.g. a filter that hid every row)")
    expected_keys = (_ENRICHED_KEYS if with_similarity
                     else _LEGACY_KEYS)
    for i, rec in enumerate(doc):
        ptr = f"/{i}"
        if not isinstance(rec, dict):
            errs.add(ptr,
                f"record must be an object; got "
                f"{type(rec).__name__}")
            continue
        keys = set(rec)
        missing = expected_keys - keys
        extra = keys - expected_keys
        if missing:
            errs.add(ptr,
                f"missing required key(s): {sorted(missing)}")
        if extra:
            errs.add(ptr,
                f"unexpected key(s) for "
                f"{'enriched' if with_similarity else 'legacy'} "
                f"schema: {sorted(extra)}")
        path = rec.get("path")
        if not isinstance(path, str) or not path:
            errs.add(ptr + "/path",
                f"path must be a non-empty string; got {path!r}")
        status = rec.get("status")
        if status not in _STATUSES:
            errs.add(ptr + "/status",
                f"status must be one of {sorted(_STATUSES)}; "
                f"got {status!r}")
        if not isinstance(rec.get("reason"), str):
            errs.add(ptr + "/reason",
                f"reason must be a string; got "
                f"{type(rec.get('reason')).__name__}")
        if with_similarity:
            if "similarity" not in rec:
                # already reported as missing key above; skip the
                # sub-validation to avoid a confusing duplicate.
                continue
            _validate_similarity(
                rec["similarity"], ptr + "/similarity", errs,
                with_labels=with_labels,
                status=status if isinstance(status, str) else "")
    return errs


# ---------------------------------------------------------------
# CLI
# ---------------------------------------------------------------
def _read_input(path: str) -> Any:
    if path == "-":
        text = sys.stdin.read()
        source = "<stdin>"
    else:
        text = Path(path).read_text(encoding="utf-8")
        source = path
    try:
        return json.loads(text)
    except json.JSONDecodeError as exc:
        print(f"validate-rename-intake: {source}: invalid JSON: "
              f"{exc}", file=sys.stderr)
        sys.exit(2)


def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser(
        prog="validate-rename-intake.py",
        description=("Validate a `rename_intake` JSON document "
                     "against the documented schema. Stdlib-only, "
                     "no external dependencies."),
        epilog=("Examples:\n"
                "  # validate the legacy 3-key schema from a file\n"
                "  python3 linter-scripts/validate-rename-intake.py "
                "audit.json\n\n"
                "  # validate the enriched 4-key schema (with "
                "--with-similarity)\n"
                "  python3 linter-scripts/validate-rename-intake.py "
                "audit.json --with-similarity\n\n"
                "  # validate from stdin (typical CI pattern, "
                "captures STDERR audit)\n"
                "  python3 linter-scripts/check-placeholder-comments"
                ".py --diff-base origin/main "
                "--list-changed-files --with-similarity --json "
                "2> audit.json >/dev/null && \\\n"
                "  python3 linter-scripts/validate-rename-intake.py "
                "audit.json --with-similarity\n\n"
                "  # print the formal JSON Schema (Draft 2020-12) "
                "for external tools\n"
                "  python3 linter-scripts/validate-rename-intake.py "
                "--print-schema --with-similarity\n"),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    p.add_argument("input", nargs="?", default="-",
        help="Path to a JSON file, or `-` for stdin (default).")
    p.add_argument("--with-similarity", action="store_true",
        help="Expect the enriched 4-key schema "
             "(`path`, `status`, `reason`, `similarity`). Use this "
             "when the upstream linter was invoked with "
             "`--with-similarity`. Default: legacy 3-key schema.")
    p.add_argument("--with-labels", action="store_true",
        help="Expect a `score_kind` field on every non-null "
             "`similarity` object (the upstream linter was invoked "
             "with `--similarity-labels`). Implies "
             "`--with-similarity`.")
    p.add_argument("--allow-empty", action="store_true",
        help="Treat an empty array as valid. By default an empty "
             "array fails because it almost always indicates a "
             "misconfigured invocation.")
    p.add_argument("--print-schema", action="store_true",
        help="Print the formal JSON Schema (Draft 2020-12) for the "
             "selected mode and exit. Useful for piping into a "
             "richer external validator (e.g. `check-jsonschema`, "
             "`ajv`) or for archiving the contract alongside CI "
             "logs.")
    p.add_argument("--quiet", action="store_true",
        help="On success, print nothing. On failure, still print "
             "the violation list.")
    args = p.parse_args(argv)

    if args.with_labels and not args.with_similarity:
        # --with-labels implies --with-similarity; auto-promote so
        # the operator doesn't have to spell both out.
        args.with_similarity = True

    if args.print_schema:
        print(json.dumps(
            _build_json_schema(
                with_similarity=args.with_similarity,
                with_labels=args.with_labels),
            indent=2, ensure_ascii=False))
        return 0

    doc = _read_input(args.input)
    errs = _validate(doc,
        with_similarity=args.with_similarity,
        with_labels=args.with_labels,
        allow_empty=args.allow_empty)

    if errs:
        n = len(doc) if isinstance(doc, list) else 0
        print(f"validate-rename-intake: {len(errs.items)} "
              f"violation(s) in {n} record(s):", file=sys.stderr)
        for ptr, msg in errs.items:
            print(f"  {ptr or '<root>'}: {msg}", file=sys.stderr)
        if errs.truncated:
            print(f"  ... (truncated at {errs.limit} errors)",
                  file=sys.stderr)
        return 1

    if not args.quiet:
        n = len(doc) if isinstance(doc, list) else 0
        mode = ("enriched" if args.with_similarity else "legacy")
        if args.with_labels:
            mode += " + labels"
        print(f"validate-rename-intake: OK — {n} record(s), "
              f"{mode} schema")
    return 0


if __name__ == "__main__":
    sys.exit(main())