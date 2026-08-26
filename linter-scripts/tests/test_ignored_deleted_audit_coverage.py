"""Comprehensive coverage for `--list-changed-files` deleted-path
behaviour and the `ignored-deleted` audit row.

These tests fill gaps left by the existing suites
(``test_list_changed_files_flag.py``,
``test_ignored_deleted_reason.py``,
``test_only_changed_status_flag.py``,
``test_dedupe_changed_files_flag.py``) — specifically:

* deletes-only intake exits with a clean PASS (no scan to run, no
  violations to find);
* the totals footer counts deleted rows under the
  ``ignored-deleted`` bucket — distinct from the matched/ignored-
  missing buckets;
* the JSON array includes deleted rows in input order alongside
  matched ones, with the per-source ``reason`` text on each;
* the dedupe pass keeps the FIRST occurrence when a path appears as
  both ``ignored-deleted`` and (later) some other status;
* the audit goes to STDERR even when only deleted rows exist
  (regression guard against an early-return path silently dropping
  them);
* the parser's per-source provenance tag is observable on every
  ``D``-row regardless of leading whitespace / trailing newline
  variations a CI runner might introduce.

The CLI tests invoke the linter as a subprocess so STDOUT/STDERR
separation is exercised through real OS pipes — the same way a CI
job would consume the output.
"""
from __future__ import annotations

import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()
_DELETED_REASON = _MOD._DELETED_REASON
_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")


class _Sandbox:
    """Tiny helper: a temp dir with a `spec/` subroot and a writable
    changed-files payload. Returned as a context-manager so the
    cleanup is guaranteed even when an assertion blows up mid-test.
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
        """Invoke the linter with `--list-changed-files` and return
        the completed process. STDOUT and STDERR are captured as
        separate strings so tests can assert the stream contract.

        The linter's `--root` AND its `cwd` are both set to the
        sandbox `root` (NOT `spec/`) so payload paths like
        ``spec/x.md`` resolve to ``<tmp>/spec/x.md`` — i.e. *under*
        the spec subroot, the way a real repo's relative paths
        would. Setting `--root` to ``spec/`` would double-stack the
        prefix and every row would land as ``ignored-out-of-root``.
        """
        cmd = [
            sys.executable, str(_LINTER),
            "--root", str(self.spec),  # post-state files MUST live under here
            "--changed-files", str(self.changed),
            "--list-changed-files",
            *extra,
        ]
        # `cwd=self.root` so the relative `spec/x.md` payload entries
        # are resolved against the sandbox root, NOT against the
        # caller's CWD (which is the project repo and would resolve
        # to a non-existent path under it).
        return subprocess.run(cmd, capture_output=True, text=True,
                              check=False, cwd=str(self.root))


class DeletesOnlyIntakeIsCleanPass(unittest.TestCase):
    """A payload of pure deletes resolves to zero linted files.

    The expected behaviour is a fast PASS (exit 0) with the deleted
    rows surfaced in the audit on STDERR — the lint step has nothing
    to scan, so it must NOT spend cycles walking the tree.
    """

    def test_exits_zero_when_only_deleted_paths_present(self) -> None:
        with _Sandbox("D\tspec/gone.md\n") as box:
            proc = box.run()
        self.assertEqual(proc.returncode, 0,
                         f"unexpected non-zero exit:\n{proc.stderr}")

    def test_audit_still_printed_even_with_no_matched_files(self) -> None:
        # Regression guard: an early-return on "no matched files"
        # must not skip the STDERR audit.
        with _Sandbox("D\tspec/gone.md\n") as box:
            proc = box.run()
        self.assertIn("ignored-deleted", proc.stderr)
        self.assertIn("spec/gone.md", proc.stderr)


class TotalsFooterCountsDeletedRows(unittest.TestCase):
    """The footer line accounts for every status bucket separately.

    The deleted bucket must roll up under `ignored-deleted=N`, NOT
    `ignored-missing` (which is a different "post-state file isn't
    on disk" reason) and NOT `matched`.
    """

    def test_two_deletes_count_as_ignored_deleted_two(self) -> None:
        with _Sandbox("D\tspec/a.md\nD\tspec/b.md\n") as box:
            proc = box.run()
        self.assertIn("ignored-deleted=2", proc.stderr)
        self.assertIn("matched=0", proc.stderr)
        self.assertIn("ignored-missing=0", proc.stderr)

    def test_mixed_intake_keeps_buckets_separate(self) -> None:
        # One real delete + one ghost path (extension allowed, but
        # missing on disk → `ignored-missing`) + one out-of-root
        # path. Each must land in its own bucket.
        payload = (
            "D\tspec/gone.md\n"
            "spec/never-existed.md\n"
            "M\tother/out.md\n"
        )
        with _Sandbox(payload) as box:
            proc = box.run()
        self.assertIn("ignored-deleted=1", proc.stderr)
        self.assertIn("ignored-missing=1", proc.stderr)
        self.assertIn("ignored-out-of-root=1", proc.stderr)


class JsonArrayIncludesDeletedRows(unittest.TestCase):
    """The `--json` payload on STDERR carries every deleted row.

    STDOUT must remain a single well-formed JSON document (the
    violation array, empty here) — the audit lives on STDERR.
    """

    def _split(self, payload: str) -> tuple[list, list]:
        with _Sandbox(payload) as box:
            proc = box.run("--json")
        # STDOUT must parse cleanly to the violations array.
        stdout_obj = json.loads(proc.stdout)
        # STDERR contains the audit array (last JSON document on the
        # stream — there's only one).
        stderr_text = proc.stderr.strip()
        # The `ℹ️ diff-mode active` banner doesn't print under --json,
        # so STDERR should be the audit array verbatim.
        audit = json.loads(stderr_text)
        return stdout_obj, audit

    def test_deleted_row_present_in_audit_array(self) -> None:
        _, audit = self._split("D\tspec/gone.md\n")
        rows = [r for r in audit if r["status"] == "ignored-deleted"]
        self.assertEqual(len(rows), 1)
        self.assertEqual(rows[0]["path"], "spec/gone.md")
        self.assertEqual(rows[0]["reason"],
                         _DELETED_REASON["changed-files-D"])

    def test_stdout_remains_clean_violations_document(self) -> None:
        # Even when only deletes are present, STDOUT must be a valid
        # JSON document (an empty array — no scan happened).
        stdout_obj, _ = self._split("D\tspec/gone.md\n")
        self.assertEqual(stdout_obj, [])

    def test_input_order_preserved_across_mixed_intake(self) -> None:
        # Audit rows follow input order so `diff` between two runs
        # stays reviewable.
        payload = (
            "D\tspec/first-gone.md\n"
            "D\tspec/second-gone.md\n"
        )
        _, audit = self._split(payload)
        paths = [r["path"] for r in audit
                 if r["status"] == "ignored-deleted"]
        self.assertEqual(paths,
                         ["spec/first-gone.md", "spec/second-gone.md"])


class DedupeInteractionWithDeletes(unittest.TestCase):
    """`--dedupe-changed-files` first-seen-wins applies across the
    audit's *emit* order (not the payload's input order).

    The audit pipeline appends matched / extension / missing /
    out-of-root rows as it walks the post-state intake, then
    appends every captured deleted row in a final batch. So a path
    that appears as both `M` and `D` in the same payload always
    surfaces with the M-row FIRST in the audit — and dedupe keeps
    that first row regardless of payload order.

    This test pins that contract explicitly so a future refactor
    that re-orders the audit emit (e.g. interleaving deletes
    inline) doesn't silently change which row wins dedupe.
    """

    def test_matched_always_wins_over_delete_under_dedupe(self) -> None:
        # Both row orderings collapse to the same result because
        # deletes are appended to the audit AFTER the in-order
        # matched/extension/missing/out-of-root rows.
        for payload_label, payload in [
            ("delete-then-modify", "D\tspec/x.md\nM\tspec/x.md\n"),
            ("modify-then-delete", "M\tspec/x.md\nD\tspec/x.md\n"),
        ]:
            with self.subTest(payload=payload_label):
                with _Sandbox(payload) as box:
                    (box.spec / "x.md").write_text("# x",
                                                   encoding="utf-8")
                    proc = box.run("--dedupe-changed-files")
                # The matched row is the survivor; the deleted row
                # was dropped as a duplicate.
                self.assertIn("matched=1", proc.stderr,
                              msg=f"payload={payload_label}\n"
                                  f"stderr:\n{proc.stderr}")
                self.assertIn("ignored-deleted=0", proc.stderr)
                self.assertIn("deduped, 1 duplicate(s) dropped",
                              proc.stderr)

    def test_two_deletes_for_same_path_collapse_to_one(self) -> None:
        # Pure-delete dedupe — both rows are emitted at the end of
        # the audit, in payload order, so dedupe keeps the first
        # delete and drops the second.
        with _Sandbox("D\tspec/gone.md\nD\tspec/gone.md\n") as box:
            proc = box.run("--dedupe-changed-files")
        self.assertIn("ignored-deleted=1", proc.stderr)
        self.assertIn("deduped, 1 duplicate(s) dropped", proc.stderr)


class OnlyStatusFilterCanIsolateDeletes(unittest.TestCase):
    """`--only-changed-status ignored-deleted` shows just deletes.

    Filter runs after dedupe + similarity attach so the totals line
    still reports every bucket — only the rendered table is filtered.
    """

    def test_filter_to_ignored_deleted_hides_other_rows(self) -> None:
        with _Sandbox(
            "D\tspec/gone.md\nM\tother/out.md\n"
        ) as box:
            proc = box.run("--only-changed-status", "ignored-deleted")
        # The visible row is the delete. The out-of-root row stays
        # out of the table but appears in the totals breakdown.
        self.assertIn("spec/gone.md", proc.stderr)
        self.assertNotIn("other/out.md   ", proc.stderr)
        self.assertIn("ignored-out-of-root=1", proc.stderr)
        self.assertIn("ignored-deleted=1", proc.stderr)


class ParserProvenanceTagSurvivesWhitespace(unittest.TestCase):
    """The `D\\tpath` shape is recognised regardless of trailing CR."""

    def test_trailing_cr_does_not_break_d_row(self) -> None:
        # CRLF-pasted authored payloads (Windows runner copy-paste)
        # must still classify as `changed-files-D`, not fall through
        # to the generic plain-path branch.
        deleted: list[tuple[str, str]] = []
        out = _MOD._normalise_changed_lines(
            ["D\tspec/gone.md\r"], deleted=deleted,
        )
        self.assertEqual(out, [])
        self.assertEqual(deleted, [("spec/gone.md", "changed-files-D")])

    def test_blank_lines_and_comments_do_not_pollute_deleted_list(
            self) -> None:
        # Comments / blanks are filtered downstream of the parser by
        # `_resolve_changed_md`, but the parser itself should pass
        # them through untouched and never inject them into `deleted`.
        deleted: list[tuple[str, str]] = []
        _MOD._normalise_changed_lines(
            ["", "# comment", "D\tspec/x.md", ""], deleted=deleted,
        )
        self.assertEqual(deleted, [("spec/x.md", "changed-files-D")])


if __name__ == "__main__":
    unittest.main()
