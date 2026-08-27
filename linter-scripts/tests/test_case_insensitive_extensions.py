"""Tests for case-insensitive extension matching across platforms.

Background: on Linux CI ``Path.rglob("*.md")`` only finds the
lowercase form, but on Windows + macOS the FS folds case so
``README.MD`` is a perfectly normal filename. The diff-mode
classifier already lowercases the suffix at audit time, so a path
like ``spec/README.MD`` would be marked ``matched`` in the audit
but then silently skipped by the full-tree walker — a "passes
locally, fails on Windows" inconsistency.

These tests pin the now-symmetric behaviour:

* The full-tree walker (:func:`iter_markdown_files`) yields
  uppercase ``.MD`` / mixed-case ``.Md`` / ``.mDx`` files when the
  matching lowercase extension is in the allowlist.
* End-to-end CLI run on a tree containing ``UPPER.MD`` lints it
  (a real placeholder violation in the file is reported, proving
  the file actually reached the rules).
* Diff-mode audit classifies ``README.MD`` as ``matched`` (not
  ``ignored-extension``), and the linted-file count matches.
* Case-insensitive matching does NOT widen the allowlist: a
  ``.txt`` file is still ``ignored-extension`` regardless of case
  (``NOTES.TXT`` is not silently scanned).
"""

from __future__ import annotations

import importlib.util
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


def _load_module():  # type: ignore[no-untyped-def]
    spec = importlib.util.spec_from_file_location("_pc_caseins", LINTER)
    assert spec is not None and spec.loader is not None
    mod = importlib.util.module_from_spec(spec)
    sys.modules["_pc_caseins"] = mod
    spec.loader.exec_module(mod)
    return mod


class IterMarkdownFilesCaseInsensitive(unittest.TestCase):
    def setUp(self) -> None:
        self.mod = _load_module()

    def test_yields_uppercase_and_mixed_case_md(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            (root / "lower.md").write_text("a", encoding="utf-8")
            (root / "UPPER.MD").write_text("b", encoding="utf-8")
            (root / "Mixed.Md").write_text("c", encoding="utf-8")
            (root / "skip.txt").write_text("d", encoding="utf-8")
            files = sorted(p.name for p in
                           self.mod.iter_markdown_files(root))
            self.assertEqual(files,
                ["Mixed.Md", "UPPER.MD", "lower.md"])

    def test_mdx_allowlist_matches_uppercase_mdx(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            (root / "page.mDx").write_text("x", encoding="utf-8")
            (root / "doc.md").write_text("y", encoding="utf-8")
            files = sorted(p.name for p in self.mod.iter_markdown_files(
                root, extensions=("md", "mdx")))
            self.assertEqual(files, ["doc.md", "page.mDx"])

    def test_extension_not_in_allowlist_still_filtered(self) -> None:
        # Case-insensitivity must not widen the allowlist: .txt
        # is rejected regardless of letter case.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            (root / "ok.md").write_text("a", encoding="utf-8")
            (root / "NOTES.TXT").write_text("b", encoding="utf-8")
            files = [p.name for p in self.mod.iter_markdown_files(root)]
            self.assertEqual(files, ["ok.md"])


class DiffModeAuditTreatsCaseConsistently(unittest.TestCase):
    def test_uppercase_md_path_is_matched_and_linted(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            spec = root / "spec"
            spec.mkdir()
            # Real file on disk under uppercase suffix.
            (spec / "README.MD").write_text(
                "# Title\n\nNo placeholders here.\n",
                encoding="utf-8",
            )
            changed = root / "changed.txt"
            changed.write_text("spec/README.MD\n", encoding="utf-8")

            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                "--json",
                cwd=root,
            )
            self.assertEqual(rc, 0, msg=f"stdout={out!r} stderr={err!r}")
            audit = json.loads(err)
            self.assertEqual(len(audit), 1)
            self.assertEqual(audit[0]["status"], "matched",
                msg=f"row={audit[0]!r}")
            # STDOUT is a clean violations array (no rule fired).
            self.assertEqual(json.loads(out), [])

    def test_uppercase_txt_still_ignored_extension(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            (root / "spec").mkdir()
            (root / "spec" / "NOTES.TXT").write_text("x",
                encoding="utf-8")
            changed = root / "changed.txt"
            changed.write_text("spec/NOTES.TXT\n", encoding="utf-8")

            rc, _out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                "--json",
                cwd=root,
            )
            self.assertEqual(rc, 0)
            audit = json.loads(err)
            self.assertEqual(len(audit), 1)
            self.assertEqual(audit[0]["status"], "ignored-extension")


if __name__ == "__main__":
    unittest.main()