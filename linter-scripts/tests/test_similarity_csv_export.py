"""Unit tests for ``--similarity-csv`` export.

Exercises ``_write_similarity_csv`` directly so the cases stay
hermetic (no git, no temp repo) while still covering every shape the
CLI flag will encounter in the wild:

* plain A/M/D rows (no rename → empty similarity cells)
* scored R/C rows (kind + integer score + old_path)
* unscored R/C rows (score cell stays empty, distinct from "0")
* the literal score 0 (rendered as ``"0"``, NOT empty)
* paths with commas / quotes (RFC 4180 quoting round-trips)
* writing to a real file vs writing to STDOUT (``"-"``)

The header row is asserted byte-for-byte so a future re-ordering of
the column contract trips immediately.
"""
from __future__ import annotations

import csv
import io
import sys
import tempfile
import unittest
from contextlib import redirect_stdout
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()

_ChangedFileAudit = _MOD._ChangedFileAudit
_RenameSimilarity = _MOD._RenameSimilarity
_write_similarity_csv = _MOD._write_similarity_csv
_SIMILARITY_CSV_HEADER = _MOD._SIMILARITY_CSV_HEADER


def _read_csv(text: str) -> list[list[str]]:
    return list(csv.reader(io.StringIO(text)))


def _write_to_tmp(rows: list) -> list[list[str]]:
    with tempfile.TemporaryDirectory() as d:
        out = Path(d) / "audit.csv"
        _write_similarity_csv(rows, str(out))
        return _read_csv(out.read_text(encoding="utf-8"))


class TestSimilarityCsvExport(unittest.TestCase):

    def test_header_is_stable_six_columns(self) -> None:
        self.assertEqual(_SIMILARITY_CSV_HEADER, (
            "path", "status", "reason", "kind", "score", "old_path",
        ))

    def test_plain_row_writes_empty_similarity_cells(self) -> None:
        rows = [_ChangedFileAudit(
            path="readme.md", status="matched", reason="ok",
        )]
        grid = _write_to_tmp(rows)
        self.assertEqual(grid[0], list(_SIMILARITY_CSV_HEADER))
        self.assertEqual(grid[1],
                         ["readme.md", "matched", "ok", "", "", ""])

    def test_scored_rename_writes_integer_score(self) -> None:
        rows = [_ChangedFileAudit(
            path="docs/new.md", status="matched", reason="ok",
            similarity=_RenameSimilarity(
                kind="R", score=92, old_path="docs/old.md",
            ),
        )]
        grid = _write_to_tmp(rows)
        self.assertEqual(grid[1], ["docs/new.md", "matched", "ok",
                                   "R", "92", "docs/old.md"])

    def test_unscored_rename_leaves_score_cell_empty(self) -> None:
        """Authored payload without a percentage → ``score=None``."""
        rows = [_ChangedFileAudit(
            path="b.md", status="matched", reason="ok",
            similarity=_RenameSimilarity(
                kind="R", score=None, old_path="a.md",
            ),
        )]
        grid = _write_to_tmp(rows)
        # kind + old_path populated, score cell EMPTY (not "0", not "-").
        self.assertEqual(grid[1],
                         ["b.md", "matched", "ok", "R", "", "a.md"])

    def test_zero_score_is_distinct_from_unscored(self) -> None:
        rows = [
            _ChangedFileAudit(
                path="zero.md", status="matched", reason="ok",
                similarity=_RenameSimilarity(
                    kind="C", score=0, old_path="src.md",
                ),
            ),
            _ChangedFileAudit(
                path="none.md", status="matched", reason="ok",
                similarity=_RenameSimilarity(
                    kind="C", score=None, old_path="src.md",
                ),
            ),
        ]
        grid = _write_to_tmp(rows)
        self.assertEqual(grid[1][4], "0")  # explicit zero
        self.assertEqual(grid[2][4], "")   # unscored
        self.assertNotEqual(grid[1][4], grid[2][4])

    def test_paths_with_commas_and_quotes_roundtrip(self) -> None:
        rows = [_ChangedFileAudit(
            path='weird, name.md', status="matched",
            reason='says "hi", with comma',
            similarity=_RenameSimilarity(
                kind="R", score=50, old_path='old, "name".md',
            ),
        )]
        grid = _write_to_tmp(rows)
        self.assertEqual(grid[1], ["weird, name.md", "matched",
                                   'says "hi", with comma',
                                   "R", "50", 'old, "name".md'])

    def test_dash_target_writes_to_stdout(self) -> None:
        rows = [_ChangedFileAudit(
            path="a.md", status="matched", reason="ok",
        )]
        buf = io.StringIO()
        with redirect_stdout(buf):
            _write_similarity_csv(rows, "-")
        grid = _read_csv(buf.getvalue())
        self.assertEqual(grid[0], list(_SIMILARITY_CSV_HEADER))
        self.assertEqual(grid[1],
                         ["a.md", "matched", "ok", "", "", ""])

    def test_empty_rows_still_writes_header(self) -> None:
        grid = _write_to_tmp([])
        self.assertEqual(grid, [list(_SIMILARITY_CSV_HEADER)])

    def test_mixed_batch_preserves_input_order(self) -> None:
        rows = [
            _ChangedFileAudit(path="1.md", status="matched",
                              reason="ok"),
            _ChangedFileAudit(
                path="2.md", status="matched", reason="ok",
                similarity=_RenameSimilarity(
                    kind="R", score=80, old_path="0.md"),
            ),
            _ChangedFileAudit(
                path="3.md", status="ignored-deleted", reason="gone"),
        ]
        grid = _write_to_tmp(rows)
        self.assertEqual([r[0] for r in grid[1:]],
                         ["1.md", "2.md", "3.md"])


if __name__ == "__main__":
    unittest.main()