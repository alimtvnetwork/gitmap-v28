"""Per-source ``reason`` strings for ``ignored-deleted`` audit rows.

Each provenance tag emitted by the parsers must map to a distinct,
human-readable ``reason`` so an operator scanning the audit table
can see *why* a path was classified as deleted (true diff `D` row
vs. authored `--changed-files` `D\tpath` payload).

Covers the end-to-end pipeline by feeding hand-built parser output
through ``_resolve_changed_md`` (no real git, no temp repo) and
inspecting the resulting ``_ChangedFileAudit`` rows.
"""
from __future__ import annotations

import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()

_resolve_changed_md = _MOD._resolve_changed_md
_DELETED_REASON = _MOD._DELETED_REASON
_DELETED_REASON_FALLBACK = _MOD._DELETED_REASON_FALLBACK


def _audit_for(payload: str) -> list:
    """Run the changed-files intake against an authored payload and
    return the resulting audit rows (deleted entries only).

    Uses a temp dir as the root so the path-existence check classifies
    surviving paths as ``ignored-missing`` rather than ``matched`` —
    that's irrelevant here, the assertions only look at the
    ``ignored-deleted`` rows.
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
    return [r for r in audit if r.status == "ignored-deleted"]


class TestDeletedReasonVocabulary(unittest.TestCase):
    def test_keys_cover_known_provenance_tags(self) -> None:
        # The two provenance tags the parsers emit today must each
        # have a dedicated reason string. Future tags should land in
        # this map alongside their parser changes.
        self.assertIn("diff-D", _DELETED_REASON)
        self.assertIn("changed-files-D", _DELETED_REASON)

    def test_each_reason_is_unique(self) -> None:
        # Distinct provenance must produce distinct reason text —
        # otherwise the per-source split is invisible to operators.
        values = list(_DELETED_REASON.values())
        self.assertEqual(len(values), len(set(values)))

    def test_fallback_is_distinct_from_known_reasons(self) -> None:
        self.assertNotIn(_DELETED_REASON_FALLBACK,
                         _DELETED_REASON.values())


class TestChangedFilesDReasonReachesAudit(unittest.TestCase):
    def test_d_row_carries_changed_files_reason(self) -> None:
        rows = _audit_for("D\tspec/gone.md\n")
        self.assertEqual(len(rows), 1)
        self.assertEqual(rows[0].path, "spec/gone.md")
        self.assertEqual(rows[0].reason,
                         _DELETED_REASON["changed-files-D"])

    def test_reason_mentions_changed_files_provenance(self) -> None:
        # A reviewer scanning the table must see at a glance that
        # this row came from the authored payload, not from a real
        # git diff D row. Asserting on a stable substring rather
        # than the full string keeps the test resilient to copy
        # tweaks while still pinning the *meaning*.
        rows = _audit_for("D\tspec/gone.md\n")
        self.assertIn("--changed-files", rows[0].reason)

    def test_multiple_d_rows_each_get_their_own_row(self) -> None:
        # The pipeline preserves input order; each delete becomes
        # its own audit entry with the same per-source reason.
        rows = _audit_for("D\tspec/a.md\nD\tspec/b.md\n")
        self.assertEqual([r.path for r in rows],
                         ["spec/a.md", "spec/b.md"])
        for r in rows:
            self.assertEqual(r.reason,
                             _DELETED_REASON["changed-files-D"])


if __name__ == "__main__":
    unittest.main()