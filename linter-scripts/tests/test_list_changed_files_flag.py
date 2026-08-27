"""End-to-end tests for the ``--list-changed-files`` audit flag.

The flag turns the diff-mode intake into an explicit, reviewable
table on STDERR. STDOUT must remain untouched (clean human summary
or single-document JSON) regardless of this flag's mode.

Coverage:

* All five status values surface with the expected ``reason``:
  ``matched``, ``ignored-extension``, ``ignored-out-of-root``,
  ``ignored-missing``, ``ignored-deleted``.
* The new ``ignored-deleted`` status is produced from a ``D\\tpath``
  row in ``--changed-files`` input (the file-based intake; the
  git-based intake exercises the same code path in a unit test).
* ``--json`` keeps STDOUT a single parseable document while the
  audit JSON array goes to STDERR.
* Without the flag, neither STDOUT nor STDERR mention the audit
  table — i.e. the flag really is opt-in.
* Outside diff mode the flag is a true no-op (no audit table).

The helpers below avoid invoking ``git`` so the suite runs in
sandboxed CI without a real repo.
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
    """Write a placeholder-free ``.md`` file so the linter passes."""
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text("# Title\n\nNo placeholders here.\n", encoding="utf-8")


class ListChangedFilesAudit(unittest.TestCase):
    def test_all_five_statuses_via_changed_files(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            spec = root / "spec"
            _write_clean_md(spec / "ok.md")
            # `outside.md` lives outside spec/ — out-of-root status.
            _write_clean_md(root / "outside.md")
            # `not-md.txt` exists but isn't in the extension allowlist.
            (root / "spec" / "notes.txt").write_text("x", encoding="utf-8")

            changed = root / "changed.txt"
            changed.write_text(
                "spec/ok.md\n"            # matched
                "outside.md\n"            # ignored-out-of-root
                "spec/notes.txt\n"        # ignored-extension
                "spec/missing.md\n"       # ignored-missing
                "D\tspec/old.md\n",       # ignored-deleted
                encoding="utf-8",
            )

            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                cwd=root,
            )
            self.assertEqual(rc, 0, msg=f"stdout={out!r} stderr={err!r}")
            # All five status labels appear in the stderr table.
            for status in (
                "matched",
                "ignored-extension",
                "ignored-out-of-root",
                "ignored-missing",
                "ignored-deleted",
            ):
                self.assertIn(status, err,
                    msg=f"missing {status!r} in stderr:\n{err}")
            # The deleted path should be reported with its original name.
            self.assertIn("spec/old.md", err)
            # STDOUT carries the human summary, not the audit table.
            self.assertNotIn("changed-file audit", out)

    def test_json_audit_goes_to_stderr_stdout_stays_parseable(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            changed = root / "changed.txt"
            changed.write_text(
                "spec/ok.md\nD\tspec/gone.md\n", encoding="utf-8")

            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                "--list-changed-files",
                "--json",
                cwd=root,
            )
            self.assertEqual(rc, 0, msg=f"stdout={out!r} stderr={err!r}")
            # STDOUT is a single JSON document (the violations list).
            parsed_out = json.loads(out)
            self.assertEqual(parsed_out, [])
            # STDERR is a separate JSON document (the audit array).
            parsed_err = json.loads(err)
            self.assertIsInstance(parsed_err, list)
            statuses = {row["status"] for row in parsed_err}
            self.assertIn("matched", statuses)
            self.assertIn("ignored-deleted", statuses)
            for row in parsed_err:
                self.assertIn("path", row)
                self.assertIn("status", row)
                self.assertIn("reason", row)
                self.assertTrue(row["reason"],
                    msg=f"empty reason: {row!r}")

    def test_flag_omitted_means_no_audit_output(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            changed = root / "changed.txt"
            changed.write_text("spec/ok.md\n", encoding="utf-8")

            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(changed),
                cwd=root,
            )
            self.assertEqual(rc, 0)
            self.assertNotIn("changed-file audit", out)
            self.assertNotIn("changed-file audit", err)

    def test_flag_outside_diff_mode_is_noop(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            _write_clean_md(root / "spec" / "ok.md")
            rc, out, err = _run(
                "--root", "spec",
                "--repo-root", str(root),
                "--list-changed-files",
                cwd=root,
            )
            self.assertEqual(rc, 0)
            self.assertNotIn("changed-file audit", err)
            self.assertNotIn("changed-file audit", out)


class DeletedPathParserUnit(unittest.TestCase):
    """Direct unit tests for the ``D``-row capture in both intake
    parsers, so the audit trail is correct without invoking git.
    """

    def setUp(self) -> None:
        # The linter file uses a hyphen so a normal ``import`` won't
        # work. Load via importlib + register in ``sys.modules`` BEFORE
        # exec so dataclass introspection (which looks the module up
        # by name to resolve forward references) can find it.
        import importlib.util
        spec = importlib.util.spec_from_file_location(
            "_pc_under_test", LINTER)
        assert spec is not None and spec.loader is not None
        mod = importlib.util.module_from_spec(spec)
        sys.modules["_pc_under_test"] = mod
        spec.loader.exec_module(mod)
        self.mod = mod

    def test_parse_name_status_captures_d_rows(self) -> None:
        # ``deleted`` now carries ``(path, source)`` tuples so the
        # audit emitter can look up a per-provenance reason string.
        # ``_parse_name_status`` is the diff-D path → source is
        # always ``"diff-D"`` for every row it captures.
        deleted: list[tuple[str, str]] = []
        out = self.mod._parse_name_status(
            "A\tspec/new.md\n"
            "D\tspec/old.md\n"
            "M\tspec/edit.md\n"
            "D\tspec/another-old.md\n",
            deleted=deleted,
        )
        self.assertEqual(out, ["spec/new.md", "spec/edit.md"])
        self.assertEqual(deleted, [
            ("spec/old.md", "diff-D"),
            ("spec/another-old.md", "diff-D"),
        ])

    def test_parse_name_status_no_deleted_arg_drops_d_rows(self) -> None:
        # Backwards-compat: omitting ``deleted=`` keeps the historical
        # behaviour of silently dropping D rows from the returned list.
        out = self.mod._parse_name_status(
            "A\tspec/new.md\nD\tspec/old.md\n")
        self.assertEqual(out, ["spec/new.md"])

    def test_normalise_changed_lines_diverts_d_rows(self) -> None:
        # ``--changed-files`` payload's ``D\tpath`` rows are tagged
        # ``changed-files-D`` so the audit emitter can distinguish
        # them from real ``git diff`` D rows.
        deleted: list[tuple[str, str]] = []
        out = self.mod._normalise_changed_lines(
            ["spec/keep.md", "D\tspec/gone.md", "spec/also.md"],
            deleted=deleted,
        )
        self.assertEqual(out, ["spec/keep.md", "spec/also.md"])
        self.assertEqual(deleted, [("spec/gone.md", "changed-files-D")])


if __name__ == "__main__":
    unittest.main()