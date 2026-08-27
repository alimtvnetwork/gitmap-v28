"""Unit + CLI tests for ``--similarity-csv-format {csv,tsv}``.

Locks in three contracts:

1. ``csv`` dialect (the default) is byte-for-byte identical to the
   legacy output — nothing changes when the new flag is absent.
2. ``tsv`` dialect uses tab separators and avoids quoting cells
   that only contained commas (the whole reason a user would
   reach for TSV in the first place).
3. Cells containing tabs / newlines / quotes still round-trip
   safely under TSV — quoting kicks in via the stdlib
   ``excel-tab`` dialect, so the export is lossless.

The header row, column order, and empty-vs-``0`` score-cell
convention must be IDENTICAL across both dialects — only the
separator changes.
"""
from __future__ import annotations

import csv
import io
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
_write_similarity_csv = _MOD._write_similarity_csv
_FORMATS = _MOD._SIMILARITY_CSV_FORMATS

_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")


def _write(rows, *, dialect: str = "csv", with_labels: bool = False) -> str:
    with tempfile.TemporaryDirectory() as d:
        out = Path(d) / "audit.out"
        _write_similarity_csv(
            rows, str(out),
            with_labels=with_labels, dialect=dialect,
        )
        return out.read_text(encoding="utf-8")


def _sample_rows() -> list:
    return [
        _ChangedFileAudit(
            path="readme.md", status="matched", reason="ok",
        ),
        _ChangedFileAudit(
            path="docs/new.md", status="matched", reason="ok",
            similarity=_RenameSimilarity(
                kind="R", score=92, old_path="docs/old.md",
            ),
        ),
        _ChangedFileAudit(
            path="copy.md", status="matched", reason="zero observed",
            similarity=_RenameSimilarity(
                kind="C", score=0, old_path="src.md",
            ),
        ),
    ]


class TestFormatVocabulary(unittest.TestCase):

    def test_recognised_formats_are_frozen(self) -> None:
        # Pin the contract: exactly two dialects, in this order.
        # Adding a third (e.g. semicolon for German Excel) is a
        # deliberate breaking decision.
        self.assertEqual(_FORMATS, ("csv", "tsv"))

    def test_default_dialect_is_csv(self) -> None:
        # Regression guard: the writer must keep producing legacy
        # CSV when no dialect kwarg is passed.
        out = _write(_sample_rows())
        first_line = out.splitlines()[0]
        self.assertIn(",", first_line)
        self.assertNotIn("\t", first_line)


class TestCsvDialectUnchanged(unittest.TestCase):
    """csv dialect must be byte-for-byte identical with/without the kwarg."""

    def test_explicit_csv_matches_default(self) -> None:
        rows = _sample_rows()
        self.assertEqual(_write(rows), _write(rows, dialect="csv"))

    def test_csv_quotes_cells_containing_commas(self) -> None:
        rows = [_ChangedFileAudit(
            path="a, with comma.md", status="matched", reason="ok",
        )]
        out = _write(rows, dialect="csv")
        self.assertIn('"a, with comma.md"', out)


class TestTsvDialect(unittest.TestCase):

    def test_tsv_uses_tab_separators(self) -> None:
        out = _write(_sample_rows(), dialect="tsv")
        for line in out.splitlines():
            self.assertIn("\t", line)
            # The default header has no embedded commas, so a clean
            # TSV must contain none either — confirms we're not
            # accidentally writing CSV-with-tabs-too.
            self.assertNotIn(",", line.split("\t")[0])

    def test_tsv_header_row_is_tab_separated_canonical_columns(self) -> None:
        out = _write(_sample_rows(), dialect="tsv")
        header = out.splitlines()[0].split("\t")
        self.assertEqual(header, [
            "path", "status", "reason", "kind", "score", "old_path",
        ])

    def test_tsv_avoids_quoting_cells_that_only_contained_commas(self) -> None:
        # The poster-child reason to switch dialects: paths with
        # commas no longer need quoting under TSV, so spreadsheets
        # show the raw path instead of `"a, b.md"`.
        rows = [_ChangedFileAudit(
            path="a, with comma.md", status="matched", reason="ok",
        )]
        out = _write(rows, dialect="tsv")
        self.assertIn("a, with comma.md", out)
        self.assertNotIn('"a, with comma.md"', out)

    def test_tsv_still_quotes_cells_containing_tabs(self) -> None:
        # Lossless round-trip: a literal tab inside a cell must
        # NOT be confused with a column separator.
        rows = [_ChangedFileAudit(
            path="weird\tpath.md", status="matched", reason="ok",
        )]
        out = _write(rows, dialect="tsv")
        # Round-trip via csv.reader to prove the cell came back
        # whole (the exact escape form is the dialect's business).
        grid = list(csv.reader(io.StringIO(out), dialect="excel-tab"))
        self.assertEqual(grid[1][0], "weird\tpath.md")

    def test_tsv_still_quotes_cells_containing_quotes_or_newlines(self) -> None:
        rows = [_ChangedFileAudit(
            path='quote"path.md', status="matched",
            reason="multi\nline reason",
        )]
        out = _write(rows, dialect="tsv")
        grid = list(csv.reader(io.StringIO(out), dialect="excel-tab"))
        self.assertEqual(grid[1][0], 'quote"path.md')
        self.assertEqual(grid[1][2], "multi\nline reason")

    def test_header_and_columns_match_csv_dialect(self) -> None:
        # Only the separator should change between dialects — the
        # column order, the header row, and the cell contents per
        # row must be identical.
        rows = _sample_rows()
        csv_grid = list(csv.reader(
            io.StringIO(_write(rows, dialect="csv"))))
        tsv_grid = list(csv.reader(
            io.StringIO(_write(rows, dialect="tsv")),
            dialect="excel-tab"))
        self.assertEqual(csv_grid, tsv_grid)

    def test_with_labels_appends_score_kind_under_tsv_too(self) -> None:
        # The labels feature is dialect-agnostic — same 7-column
        # layout, just with tabs.
        out = _write(_sample_rows(), dialect="tsv", with_labels=True)
        header = out.splitlines()[0].split("\t")
        self.assertEqual(header[-1], "score_kind")
        grid = list(csv.reader(io.StringIO(out), dialect="excel-tab"))
        # Plain row: empty score_kind. Scored R row: rename-similarity.
        # Score-0 C row: copy-similarity (NOT unscored — the row was
        # observed, just rated dissimilar).
        self.assertEqual(grid[1][-1], "")
        self.assertEqual(grid[2][-1], "rename-similarity")
        self.assertEqual(grid[3][-1], "copy-similarity")


class TestCliGate(unittest.TestCase):
    """Reject anything outside the {csv,tsv} vocabulary at parse time."""

    def _run(self, *args: str) -> subprocess.CompletedProcess:
        with tempfile.TemporaryDirectory() as d:
            return subprocess.run(
                [sys.executable, str(_LINTER), "--root", d, *args],
                capture_output=True, text=True,
            )

    def test_unknown_format_is_rejected_by_argparse(self) -> None:
        cp = self._run("--similarity-csv-format", "xls")
        # argparse's "invalid choice" message exits with code 2 and
        # writes to STDERR — same failure-stream contract as every
        # other CLI gate documented in README-rename-intake.md.
        self.assertEqual(cp.returncode, 2)
        self.assertIn("--similarity-csv-format", cp.stderr)
        self.assertIn("invalid choice", cp.stderr.lower())
        self.assertEqual(cp.stdout, "")

    def test_csv_choice_accepted(self) -> None:
        cp = self._run("--similarity-csv-format", "csv")
        self.assertEqual(cp.returncode, 0, msg=cp.stderr)

    def test_tsv_choice_accepted(self) -> None:
        cp = self._run("--similarity-csv-format", "tsv")
        self.assertEqual(cp.returncode, 0, msg=cp.stderr)


if __name__ == "__main__":
    unittest.main()