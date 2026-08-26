"""Exhaustive cross-surface check that every ``ignored-deleted``
row's ``reason`` exactly matches one of the allowed resolved
values from the shared vocabulary (``_DELETED_REASON``).

Sibling tests already cover specific provenance tags
(``test_ignored_deleted_reason.py``,
``test_ignored_deleted_audit_coverage.py``). This module is the
contract guard: regardless of *how* a deleted row reaches the
audit (in-process via :func:`_resolve_changed_md`, parser-direct
via :func:`_parse_name_status`, or out via the
CLI's text / JSON / CSV renderings), its ``reason`` must be a
member of the canonical set ``set(_DELETED_REASON.values())``.

The fallback string (``_DELETED_REASON_FALLBACK``) is *intentionally
excluded* from the allow-list: it exists for forward-compat with
future provenance tags the parsers may add before the vocabulary
catches up. A production row landing on the fallback today means
the parser emitted a tag the vocabulary doesn't know about — that
is the bug these tests are designed to catch.

Identity guard
--------------
The dataclass-level test compares with ``is`` (object identity),
not just ``==``. The current implementation looks reasons up via
``_DELETED_REASON.get(src, fallback)`` so the vocabulary string
object is reused verbatim. A future refactor that returns a
freshly-formatted string (``template.format(...)``) would still
compare equal but break identity — which would also break the
"closed vocabulary" contract this test is pinning. Equality is
used for the rendered surfaces (text / JSON / CSV) because those
round-trip through serialisation and lose Python identity by
construction.
"""
from __future__ import annotations

import csv
import io
import json
import re
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()
_DELETED_REASON = _MOD._DELETED_REASON
_DELETED_REASON_FALLBACK = _MOD._DELETED_REASON_FALLBACK
_resolve_changed_md = _MOD._resolve_changed_md
_parse_name_status = _MOD._parse_name_status
_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")

# Canonical allow-list: the set of strings that may legitimately
# appear as the ``reason`` on an ``ignored-deleted`` row today.
# Frozen at module load — a vocabulary change must reflect here
# explicitly, not silently widen the allow-list at test time.
ALLOWED_DELETED_REASONS: frozenset[str] = frozenset(
    _DELETED_REASON.values())


def _audit_via_changed_files(payload: str) -> list:
    """Drive the live ``--changed-files`` intake and return the
    full audit list. The test harnesses can then filter it down to
    ``ignored-deleted`` rows and assert on ``reason``.
    """
    audit: list = []
    with tempfile.TemporaryDirectory() as d:
        root = Path(d).resolve()
        cf = root / "changed.txt"
        cf.write_text(payload, encoding="utf-8")
        _resolve_changed_md(
            repo_root=root, root=root,
            diff_base=None, changed_files=str(cf),
            extensions=("md",), audit=audit,
        )
    return audit


def _audit_via_diff_parser(diff_text: str) -> list:
    """Run the ``git diff --name-status`` parser directly so the
    ``diff-D`` provenance tag is exercised without a real git
    repo. The parser fills ``deleted`` in-place; the assertions
    look up each tag in ``_DELETED_REASON`` to mirror what the
    audit emitter does at line 2154.
    """
    deleted: list[tuple[str, str]] = []
    _parse_name_status(diff_text, deleted=deleted)
    return deleted


# --- 1. Vocabulary-level guarantees ------------------------------

class CanonicalAllowedSetIsWellFormed(unittest.TestCase):
    """Sanity checks on the allow-list itself, before exercising
    any pipeline. These would catch a vocabulary edit that breaks
    invariants the downstream tests rely on.
    """

    def test_allow_list_is_non_empty(self) -> None:
        # An empty vocabulary would make every downstream assertion
        # vacuously fail — fail loud here instead.
        self.assertGreater(len(ALLOWED_DELETED_REASONS), 0)

    def test_every_allowed_value_is_a_non_empty_string(self) -> None:
        for value in ALLOWED_DELETED_REASONS:
            self.assertIsInstance(value, str)
            self.assertNotEqual(value.strip(), "")

    def test_fallback_is_not_in_allow_list(self) -> None:
        # The fallback exists for unknown future tags. A production
        # row landing on it must NOT pass this contract.
        self.assertNotIn(_DELETED_REASON_FALLBACK,
                         ALLOWED_DELETED_REASONS)

    def test_allow_list_size_matches_vocabulary_size(self) -> None:
        # If two tags collapse to the same reason string, the
        # per-source split is invisible to operators. This guard
        # catches an accidental dedupe by typo.
        self.assertEqual(len(ALLOWED_DELETED_REASONS),
                         len(_DELETED_REASON))


# --- 2. Dataclass / in-process pipeline --------------------------

class DataclassRowsCarryAllowedReason(unittest.TestCase):
    """Rows produced by :func:`_resolve_changed_md` for every
    delete-shape payload land with a ``reason`` that is exactly a
    member of the allow-list — and the SAME object the vocabulary
    holds (identity, not just equality).
    """

    def test_changed_files_d_row_reason_is_in_allow_list(self) -> None:
        rows = [r for r in _audit_via_changed_files("D\tspec/x.md\n")
                if r.status == "ignored-deleted"]
        self.assertEqual(len(rows), 1)
        self.assertIn(rows[0].reason, ALLOWED_DELETED_REASONS)

    def test_changed_files_d_row_reason_is_same_object_as_vocab(
            self) -> None:
        # Identity guard: the lookup must return the vocabulary's
        # own string object, not a freshly-formatted copy. See the
        # module docstring for the rationale.
        rows = [r for r in _audit_via_changed_files("D\tspec/x.md\n")
                if r.status == "ignored-deleted"]
        self.assertIs(rows[0].reason,
                      _DELETED_REASON["changed-files-D"])

    def test_no_deleted_row_uses_the_fallback_string(self) -> None:
        # Mixed payload across every delete-shape we know how to
        # author + a non-delete row. Every captured deleted row's
        # reason must be in the allow-list — fallback is forbidden.
        payload = (
            "D\tspec/a.md\n"
            "D\tspec/b.md\n"
            "M\tspec/keep.md\n"
        )
        rows = [r for r in _audit_via_changed_files(payload)
                if r.status == "ignored-deleted"]
        self.assertGreater(len(rows), 0)
        for r in rows:
            self.assertIn(r.reason, ALLOWED_DELETED_REASONS,
                          msg=f"path {r.path!r} fell back to "
                              f"{r.reason!r}")
            self.assertNotEqual(r.reason, _DELETED_REASON_FALLBACK)


# --- 3. Parser-direct (diff-D) coverage --------------------------

class DiffParserDRowReasonResolution(unittest.TestCase):
    """The ``diff-D`` provenance tag is only reachable via the
    real ``git diff --name-status`` parser, which the in-process
    ``--changed-files`` harness above doesn't touch. Drive the
    parser directly so this contract is independently exercised.
    """

    def test_diff_d_tag_resolves_to_allowed_reason(self) -> None:
        deleted = _audit_via_diff_parser("D\tspec/gone.md\n")
        self.assertEqual(len(deleted), 1)
        path, src = deleted[0]
        self.assertEqual(src, "diff-D")
        # Mirror the audit emitter's lookup (line 2154) to confirm
        # the resolved reason is in the allow-list.
        resolved = _DELETED_REASON.get(src, _DELETED_REASON_FALLBACK)
        self.assertIn(resolved, ALLOWED_DELETED_REASONS)
        self.assertIs(resolved, _DELETED_REASON["diff-D"])

    def test_every_known_provenance_tag_resolves_into_allow_list(
            self) -> None:
        # Iterate the vocabulary itself so adding a new tag
        # automatically extends coverage — no test edit required.
        for src in _DELETED_REASON:
            with self.subTest(source=src):
                resolved = _DELETED_REASON.get(
                    src, _DELETED_REASON_FALLBACK)
                self.assertIn(resolved, ALLOWED_DELETED_REASONS)


# --- 4. CLI surface coverage (JSON / text / CSV) -----------------

class _Sandbox:
    """Tiny temp-dir harness around a `--changed-files` payload —
    same pattern as ``test_ignored_deleted_audit_coverage.py``'s
    sandbox so the CLI invocation here matches what real CI runs.
    """

    def __init__(self, payload: str) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.root = Path(self._tmp.name).resolve()
        self.spec = self.root / "spec"
        self.spec.mkdir()
        self.changed = self.root / "changed.txt"
        self.changed.write_text(payload, encoding="utf-8")

    def __enter__(self) -> "_Sandbox":
        return self

    def __exit__(self, *exc: object) -> None:
        self._tmp.cleanup()

    def run(self, *extra: str) -> subprocess.CompletedProcess:
        cmd = [
            sys.executable, str(_LINTER),
            "--root", str(self.spec),
            "--changed-files", str(self.changed),
            "--list-changed-files",
            *extra,
        ]
        return subprocess.run(cmd, capture_output=True, text=True,
                              check=False, cwd=str(self.root))


# Multi-delete payload reused across all surfaces so cross-surface
# parity can be asserted on a non-trivial row count.
_MULTI_DELETE_PAYLOAD = (
    "D\tspec/alpha.md\n"
    "D\tspec/beta.md\n"
    "D\tspec/gamma.md\n"
)


class JsonOutputReasonsAreInAllowList(unittest.TestCase):
    """The ``--json`` audit array on STDERR is the schema-stable
    surface dashboards consume. Every ``ignored-deleted`` row must
    carry a ``reason`` from the allow-list — string equality, since
    the value has round-tripped through JSON serialisation.
    """

    def test_every_deleted_row_in_json_has_allowed_reason(
            self) -> None:
        with _Sandbox(_MULTI_DELETE_PAYLOAD) as box:
            proc = box.run("--json")
        self.assertEqual(proc.returncode, 0, msg=proc.stderr)
        audit = json.loads(proc.stderr.strip())
        deleted = [r for r in audit if r["status"] == "ignored-deleted"]
        self.assertEqual(len(deleted), 3)
        for row in deleted:
            self.assertIn(row["reason"], ALLOWED_DELETED_REASONS,
                          msg=f"path {row.get('path')!r} reason "
                              f"{row.get('reason')!r} not in allow-list")

    def test_json_deleted_row_reason_never_fallback(self) -> None:
        with _Sandbox(_MULTI_DELETE_PAYLOAD) as box:
            proc = box.run("--json")
        audit = json.loads(proc.stderr.strip())
        for row in audit:
            if row["status"] == "ignored-deleted":
                self.assertNotEqual(row["reason"],
                                    _DELETED_REASON_FALLBACK)


class TextOutputReasonsAreInAllowList(unittest.TestCase):
    """The default human text table is what an operator scrolling
    CI logs sees first. Pull each ``ignored-deleted`` row's reason
    cell out of the rendered table and check the allow-list.
    """

    # The text table's row shape is:
    #   <leading-ws> <status> <whitespace> <path> <whitespace> <reason...>
    # Column boundaries are whitespace-padded; the reason cell
    # extends to end-of-line. Anchor on ``ignored-deleted`` to
    # avoid pulling other status rows by accident.
    _ROW_RE = re.compile(
        r"^\s*ignored-deleted\s+(?P<path>\S+)\s+(?P<reason>.+?)\s*$")

    def _extract_reasons(self, stderr: str) -> list[str]:
        out: list[str] = []
        for line in stderr.splitlines():
            m = self._ROW_RE.match(line)
            if m:
                out.append(m.group("reason"))
        return out

    def test_text_table_reasons_in_allow_list(self) -> None:
        with _Sandbox(_MULTI_DELETE_PAYLOAD) as box:
            proc = box.run()
        self.assertEqual(proc.returncode, 0, msg=proc.stderr)
        reasons = self._extract_reasons(proc.stderr)
        self.assertEqual(len(reasons), 3,
                         msg=f"expected 3 deleted rows, got "
                             f"{len(reasons)}\nstderr:\n{proc.stderr}")
        for reason in reasons:
            self.assertIn(reason, ALLOWED_DELETED_REASONS,
                          msg=f"text reason {reason!r} not in allow-list")


class CsvOutputReasonsAreInAllowList(unittest.TestCase):
    """The ``--similarity-csv`` export is the spreadsheet surface.
    Parse it with the stdlib ``csv`` module so the test exercises
    the same quoting rules a real consumer would.
    """

    def test_csv_deleted_rows_have_allowed_reason(self) -> None:
        with _Sandbox(_MULTI_DELETE_PAYLOAD) as box:
            csv_path = box.root / "audit.csv"
            proc = box.run("--similarity-csv", str(csv_path))
            self.assertEqual(proc.returncode, 0, msg=proc.stderr)
            self.assertTrue(csv_path.is_file(),
                            msg=f"CSV not produced; stderr:\n"
                                f"{proc.stderr}")
            with csv_path.open(encoding="utf-8", newline="") as fh:
                reader = csv.DictReader(fh)
                rows = list(reader)
        deleted = [r for r in rows if r["status"] == "ignored-deleted"]
        self.assertEqual(len(deleted), 3)
        for row in deleted:
            self.assertIn(row["reason"], ALLOWED_DELETED_REASONS,
                          msg=f"CSV reason {row['reason']!r} not in "
                              f"allow-list")


# --- 5. Cross-surface parity -------------------------------------

class ReasonsAreIdenticalAcrossSurfaces(unittest.TestCase):
    """The same payload must produce the same set of
    ``(path, reason)`` pairs whether rendered as text, JSON, or
    CSV. Drift between surfaces would make the allow-list test
    above pass on one channel and silently fail on another.
    """

    _ROW_RE = TextOutputReasonsAreInAllowList._ROW_RE

    def _text_pairs(self, stderr: str) -> set[tuple[str, str]]:
        out: set[tuple[str, str]] = set()
        for line in stderr.splitlines():
            m = self._ROW_RE.match(line)
            if m:
                out.add((m.group("path"), m.group("reason")))
        return out

    def test_text_json_csv_emit_identical_path_reason_pairs(
            self) -> None:
        # Run once per surface against the same payload, collect
        # ``(path, reason)`` for every ignored-deleted row, and
        # assert the three sets are equal.
        with _Sandbox(_MULTI_DELETE_PAYLOAD) as box:
            text_proc = box.run()
            json_proc = box.run("--json")
            csv_path = box.root / "audit.csv"
            csv_proc = box.run("--similarity-csv", str(csv_path))
            with csv_path.open(encoding="utf-8", newline="") as fh:
                csv_rows = list(csv.DictReader(fh))
        # Reference ``csv_proc`` so a future failure includes its
        # exit context — silences linters that flag the unused name.
        self.assertEqual(csv_proc.returncode, 0, msg=csv_proc.stderr)

        text_pairs = self._text_pairs(text_proc.stderr)

        json_audit = json.loads(json_proc.stderr.strip())
        json_pairs = {
            (r["path"], r["reason"]) for r in json_audit
            if r["status"] == "ignored-deleted"
        }

        csv_pairs = {
            (r["path"], r["reason"]) for r in csv_rows
            if r["status"] == "ignored-deleted"
        }

        self.assertEqual(text_pairs, json_pairs,
                         msg="text vs JSON (path, reason) drift")
        self.assertEqual(json_pairs, csv_pairs,
                         msg="JSON vs CSV (path, reason) drift")
        # And every reason in the unified set is in the allow-list.
        for _path, reason in json_pairs:
            self.assertIn(reason, ALLOWED_DELETED_REASONS)


if __name__ == "__main__":
    unittest.main()