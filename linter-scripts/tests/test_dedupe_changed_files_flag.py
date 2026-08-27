"""Tests for ``--dedupe-changed-files``.

Behaviour under test:

* The flag collapses repeated ``path`` values in the
  ``--list-changed-files`` audit, keeping the FIRST-seen row's
  ``status`` + ``reason`` verbatim.
* Order of surviving rows matches first-occurrence order so the
  rendered table stays stable + reviewable.
* The text-mode header reports the drop count
  (``deduped, N duplicate(s) dropped``); JSON mode reflects the
  collapse via array length.
* The flag is a true no-op when ``--list-changed-files`` is absent
  (no audit table appears, exit code unchanged).
* Idempotence: passing the flag against an already-unique stream
  reports ``0 duplicate(s) dropped`` and emits every row.
* Direct unit test on ``_dedupe_audit_rows`` proves first-seen
  semantics across mixed status values for the same path (e.g. a
  ``matched`` followed by an ``ignored-deleted`` keeps ``matched``).
"""

from __future__ import annotations

import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")


def _run(*args: str, cwd: Path) -> tuple[int, str, str]:
    r = subprocess.run([sys.executable, str(LINTER), *args],
                       cwd=cwd, capture_output=True, text=True)
    return r.returncode, r.stdout, r.stderr


def _write_clean_md(p: Path) -> None:
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text("# Title\n\nNo placeholders.\n", encoding="utf-8")


def _load_module():  # type: ignore[no-untyped-def]
    """Import the linter module by file path (hyphenated name)."""
    import importlib.util
    spec = importlib.util.spec_from_file_location("_pc_dedupe", LINTER)
    assert spec is not None and spec.loader is not None
    mod = importlib.util.module_from_spec(spec)
    sys.modules["_pc_dedupe"] = mod
    spec.loader.exec_module(mod)
    return mod


class DedupeChangedFilesCli(unittest.TestCase):
    def test_text_mode_collapses_duplicates_and_reports_count(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            changed = root / "changed.txt"
            # Same path listed three times → 2 dropped.
            changed.write_text(
                "spec/ok.md\nspec/ok.md\nspec/ok.md\n",
                encoding="utf-8",
            )
            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                "--dedupe-changed-files",
                cwd=root,
            )
            self.assertEqual(rc, 0, msg=f"stderr={err!r}")
            self.assertIn("2 duplicate(s) dropped", err)
            # Only one data row for spec/ok.md (plus header/totals).
            data_rows = [ln for ln in err.splitlines()
                         if "spec/ok.md" in ln]
            self.assertEqual(len(data_rows), 1,
                msg=f"expected one row, got: {data_rows!r}")

    def test_json_mode_reflects_collapse_in_array_length(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            changed = root / "changed.txt"
            changed.write_text(
                "spec/ok.md\nspec/ok.md\nD\tspec/old.md\n"
                "D\tspec/old.md\n",
                encoding="utf-8",
            )
            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                "--dedupe-changed-files",
                "--json",
                cwd=root,
            )
            self.assertEqual(rc, 0, msg=f"stderr={err!r}")
            audit = json.loads(err)
            paths = [r["path"] for r in audit]
            # Each path appears exactly once, in first-seen order.
            self.assertEqual(paths, ["spec/ok.md", "spec/old.md"])
            # STDOUT remains a clean violations array.
            self.assertEqual(json.loads(out), [])

    def test_idempotent_when_no_duplicates(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            changed = root / "changed.txt"
            changed.write_text(
                "spec/ok.md\nD\tspec/old.md\n", encoding="utf-8")
            rc, _out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                "--dedupe-changed-files",
                cwd=root,
            )
            self.assertEqual(rc, 0)
            self.assertIn("0 duplicate(s) dropped", err)

    def test_flag_without_list_changed_files_is_noop(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            changed = root / "changed.txt"
            changed.write_text("spec/ok.md\n", encoding="utf-8")
            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--dedupe-changed-files",
                cwd=root,
            )
            self.assertEqual(rc, 0)
            # No audit table on either stream.
            self.assertNotIn("changed-file audit", err)
            self.assertNotIn("changed-file audit", out)


class DedupeAuditRowsUnit(unittest.TestCase):
    def setUp(self) -> None:
        self.mod = _load_module()

    def test_first_seen_status_wins_across_mixed_statuses(self) -> None:
        Row = self.mod._ChangedFileAudit
        rows = [
            Row(path="a.md", status="matched", reason="r1"),
            # Later row for the same path with a different status —
            # must NOT overwrite the matched verdict.
            Row(path="a.md", status="ignored-deleted", reason="r2"),
            Row(path="b.md", status="ignored-extension", reason="r3"),
            Row(path="a.md", status="ignored-missing", reason="r4"),
        ]
        out, dropped = self.mod._dedupe_audit_rows(rows)
        self.assertEqual(dropped, 2)
        self.assertEqual([r.path for r in out], ["a.md", "b.md"])
        # First-seen status + reason preserved verbatim.
        self.assertEqual(out[0].status, "matched")
        self.assertEqual(out[0].reason, "r1")

    def test_empty_input_returns_empty(self) -> None:
        out, dropped = self.mod._dedupe_audit_rows([])
        self.assertEqual(out, [])
        self.assertEqual(dropped, 0)

    def test_idempotent_on_already_unique_input(self) -> None:
        Row = self.mod._ChangedFileAudit
        rows = [
            Row(path="a.md", status="matched", reason="r1"),
            Row(path="b.md", status="matched", reason="r2"),
        ]
        first, dropped1 = self.mod._dedupe_audit_rows(rows)
        second, dropped2 = self.mod._dedupe_audit_rows(first)
        self.assertEqual(dropped1, 0)
        self.assertEqual(dropped2, 0)
        self.assertEqual(first, second)


if __name__ == "__main__":
    unittest.main()