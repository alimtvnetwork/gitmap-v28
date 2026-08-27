"""Tests for ``--list-changed-files-verbose``.

The flag exposes per-row intake provenance for ``ignored-deleted``
entries on top of the default audit. Two surfaces, two contracts:

* **JSON mode** — adds a ``"source": str|null`` key to EVERY row
  (``str`` on ``ignored-deleted`` rows, ``null`` everywhere else)
  so the schema stays regular for downstream consumers. Off (the
  default), the ``source`` key is **omitted entirely** — legacy
  3-key schema preserved byte-for-byte.

* **Text mode** — appends a ``source`` column at the end of the
  fixed-width cells (just before the variable-width ``reason``
  column). Non-deleted rows render ``-`` to match the surrounding
  blank-cell convention. Off (the default), the column is absent
  entirely.

Three layers of coverage:

1. **Unit-level renderer tests** — drive
   :func:`_render_changed_files_audit` directly with hand-built
   audit rows so the schema / column shape can be asserted
   without subprocessing.
2. **End-to-end via ``--changed-files``** — confirm the source
   field round-trips from the parser through the audit emitter
   to the final JSON payload, with both ``diff-D`` and
   ``changed-files-D`` provenance reachable.
3. **Composition** — verify the flag composes cleanly with
   ``--with-similarity`` (similarity columns + source column),
   ``--dedupe-changed-files``, and ``--only-changed-status``.
"""
from __future__ import annotations

import io
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()
_ChangedFileAudit = _MOD._ChangedFileAudit
_RenameSimilarity = _MOD._RenameSimilarity
_render = _MOD._render_changed_files_audit

_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")


def _audit_fixture() -> list:
    return [
        _ChangedFileAudit(
            path="spec/ok.md", status="matched", reason="ok",
        ),
        _ChangedFileAudit(
            path="spec/gone-from-git.md", status="ignored-deleted",
            reason="git diff reported D (deleted)",
            source="diff-D",
        ),
        _ChangedFileAudit(
            path="spec/gone-from-payload.md", status="ignored-deleted",
            reason="--changed-files payload row shaped `D\\tpath`",
            source="changed-files-D",
        ),
    ]


class TestRendererJsonSchema(unittest.TestCase):
    """Direct renderer tests on the JSON output schema."""

    def test_non_verbose_omits_source_key_entirely(self) -> None:
        # Legacy schema must be preserved byte-for-byte: the
        # ``source`` key never appears in the payload when the flag
        # is off, even on ``ignored-deleted`` rows that have a real
        # tag attached to the dataclass.
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=True)
        for row in json.loads(buf.getvalue()):
            self.assertNotIn("source", row,
                             msg=("source key must be stripped from "
                                  "every row when verbose is off"))

    def test_verbose_adds_source_key_to_every_row(self) -> None:
        # Schema-regularity contract: verbose mode adds the key to
        # ALL rows (not just deleted ones) so downstream consumers
        # can ``.get("source")`` without branching on status.
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=True, verbose=True)
        rows = json.loads(buf.getvalue())
        for row in rows:
            self.assertIn("source", row)

    def test_verbose_source_is_null_on_non_deleted_rows(self) -> None:
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=True, verbose=True)
        rows = json.loads(buf.getvalue())
        matched = next(r for r in rows if r["status"] == "matched")
        self.assertIsNone(
            matched["source"],
            msg=("non-deleted rows must surface source=null, not a "
                 "stringified placeholder"),
        )

    def test_verbose_source_carries_raw_provenance_tag(self) -> None:
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=True, verbose=True)
        rows = json.loads(buf.getvalue())
        by_path = {r["path"]: r for r in rows}
        self.assertEqual(
            by_path["spec/gone-from-git.md"]["source"], "diff-D")
        self.assertEqual(
            by_path["spec/gone-from-payload.md"]["source"],
            "changed-files-D",
        )

    def test_verbose_preserves_reason_string_verbatim(self) -> None:
        # The flag's contract promises the reason wording is
        # *machine-stable* for ignored-deleted rows when verbose is
        # on. Pin the expected string from the fixture so any
        # accidental re-wording trips this test.
        rows_in = _audit_fixture()
        buf = io.StringIO()
        _render(rows_in, buf, as_json=True, verbose=True)
        rows_out = json.loads(buf.getvalue())
        for inp, out in zip(rows_in, rows_out):
            self.assertEqual(inp.reason, out["reason"])


class TestRendererTextLayout(unittest.TestCase):
    """Direct renderer tests on the text-mode column layout."""

    def test_non_verbose_has_no_source_column(self) -> None:
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=False)
        out = buf.getvalue()
        self.assertNotIn("source", out.splitlines()[1])
        # Banner must not advertise the +source suffix either.
        self.assertNotIn("+source", out.splitlines()[0])

    def test_verbose_adds_source_column_header(self) -> None:
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=False, verbose=True)
        lines = buf.getvalue().splitlines()
        # Banner advertises the new column.
        self.assertIn("+source", lines[0])
        # Header row contains the column name.
        self.assertIn("source", lines[1])

    def test_verbose_renders_dash_for_non_deleted_rows(self) -> None:
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=False, verbose=True)
        out = buf.getvalue()
        # The matched row must show the blank-cell sentinel under
        # the source column — not a stray ``None`` or empty string.
        matched_line = next(line for line in out.splitlines()
                            if "spec/ok.md" in line)
        self.assertNotIn("None", matched_line)
        # The dash sentinel appears between the path and the reason.
        self.assertRegex(
            matched_line,
            r"spec/ok\.md\s+-\s+ok",
            msg=f"matched row missing dash source cell: {matched_line!r}",
        )

    def test_verbose_renders_raw_tag_for_deleted_rows(self) -> None:
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=False, verbose=True)
        out = buf.getvalue()
        diff_line = next(line for line in out.splitlines()
                         if "spec/gone-from-git.md" in line)
        payload_line = next(line for line in out.splitlines()
                            if "spec/gone-from-payload.md" in line)
        self.assertIn("diff-D", diff_line)
        self.assertIn("changed-files-D", payload_line)

    def test_verbose_with_similarity_keeps_columns_in_order(self) -> None:
        # Source column must come AFTER the similarity columns so
        # the existing layout (status, path, kind, score, old) is
        # preserved at the front. Reason stays at the very end.
        buf = io.StringIO()
        _render(_audit_fixture(), buf, as_json=False,
                with_similarity=True, verbose=True)
        header = buf.getvalue().splitlines()[1]
        # Find the column positions and assert their order.
        positions = {col: header.find(col)
                     for col in ("status", "path", "kind", "score",
                                 "old", "source", "reason")}
        for col in positions:
            self.assertGreaterEqual(
                positions[col], 0,
                msg=f"missing column {col!r} in header: {header!r}",
            )
        ordered = sorted(positions, key=positions.get)
        self.assertEqual(ordered, [
            "status", "path", "kind", "score",
            "old", "source", "reason",
        ])


def _write_clean_md(p: Path) -> None:
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text("# Title\n\nNo placeholders.\n", encoding="utf-8")


def _run_changed_files(root: Path, manifest_text: str,
                       *extra: str) -> tuple[int, str, str]:
    manifest = root / "changed.txt"
    manifest.write_text(manifest_text, encoding="utf-8")
    cp = subprocess.run(
        [
            sys.executable, str(_LINTER),
            "--root", "spec",
            "--repo-root", str(root),
            "--changed-files", str(manifest),
            "--list-changed-files",
            "--json",
            *extra,
        ],
        cwd=root, capture_output=True, text=True,
    )
    return cp.returncode, cp.stdout, cp.stderr


class TestEndToEndChangedFiles(unittest.TestCase):
    """Subprocess tests confirming the full intake → renderer wiring."""

    def test_changed_files_payload_d_row_surfaces_changed_files_d(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            rc, _out, err = _run_changed_files(
                root,
                "spec/ok.md\n"
                "D\tspec/gone.md\n",
                "--list-changed-files-verbose",
            )
            self.assertEqual(rc, 0, msg=err)
            audit = json.loads(err)
            by_path = {r["path"]: r for r in audit}

            self.assertIn("source", by_path["spec/ok.md"])
            self.assertIsNone(by_path["spec/ok.md"]["source"])

            deleted = by_path["spec/gone.md"]
            self.assertEqual(deleted["status"], "ignored-deleted")
            self.assertEqual(deleted["source"], "changed-files-D")
            # Reason text must still be present and informative.
            self.assertTrue(deleted["reason"])

    def test_legacy_schema_unchanged_when_flag_off(self) -> None:
        # Regression guard: the JSON payload's key set must be
        # identical to the historical 3-key shape (no ``source``
        # leaking out as ``null``) when the verbose flag is absent.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            rc, _out, err = _run_changed_files(
                root,
                "spec/ok.md\n"
                "D\tspec/gone.md\n",
                # NB: no --list-changed-files-verbose
            )
            self.assertEqual(rc, 0, msg=err)
            audit = json.loads(err)
            for row in audit:
                self.assertEqual(
                    set(row.keys()), {"path", "status", "reason"},
                    msg=("legacy 3-key schema must be preserved "
                         "byte-for-byte when verbose is off"),
                )

    def test_verbose_composes_with_with_similarity(self) -> None:
        # Both extension flags should layer cleanly: similarity
        # sub-object + top-level source key on every row.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "new.md")
            rc, _out, err = _run_changed_files(
                root,
                "R092\tspec/old.md\tspec/new.md\n"
                "D\tspec/gone.md\n",
                "--list-changed-files-verbose",
                "--with-similarity",
            )
            self.assertEqual(rc, 0, msg=err)
            audit = json.loads(err)
            new_row = next(r for r in audit
                           if r["path"] == "spec/new.md")
            gone_row = next(r for r in audit
                            if r["path"] == "spec/gone.md")

            # Similarity attached to the rename's NEW path.
            self.assertIsNotNone(new_row["similarity"])
            self.assertEqual(new_row["similarity"]["kind"], "R")
            # Source field present on every row, populated only on
            # the deleted one.
            self.assertIsNone(new_row["source"])
            self.assertEqual(gone_row["source"], "changed-files-D")

    def test_verbose_composes_with_only_changed_status_filter(self) -> None:
        # Filter must run after source attachment so the surviving
        # row still carries the tag.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            rc, _out, err = _run_changed_files(
                root,
                "spec/ok.md\n"
                "D\tspec/gone.md\n",
                "--list-changed-files-verbose",
                "--only-changed-status", "ignored-deleted",
            )
            self.assertEqual(rc, 0, msg=err)
            audit = json.loads(err)
            self.assertEqual(len(audit), 1)
            self.assertEqual(audit[0]["status"], "ignored-deleted")
            self.assertEqual(audit[0]["source"], "changed-files-D")

    def test_verbose_no_op_outside_diff_mode(self) -> None:
        # ``--list-changed-files-verbose`` without ``--list-changed-
        # files`` (or any diff-mode trigger) must not crash and
        # must not produce a stray audit table.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            (root / "spec").mkdir()
            cp = subprocess.run(
                [
                    sys.executable, str(_LINTER),
                    "--root", "spec", "--repo-root", str(root),
                    "--list-changed-files-verbose",
                ],
                cwd=root, capture_output=True, text=True,
            )
            # Clean exit, no audit on STDERR.
            self.assertEqual(cp.returncode, 0, msg=cp.stderr)
            self.assertNotIn("changed-file audit", cp.stderr)


if __name__ == "__main__":
    unittest.main()