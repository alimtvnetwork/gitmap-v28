"""Unit tests for ``--similarity-labels`` (per-kind score discriminator).

Covers the three surfaces the flag touches in tandem:

* JSON audit  — adds ``score_kind`` to the nested ``similarity``
  object; absent on plain rows whose ``similarity`` is ``null``;
  legacy schema is preserved byte-for-byte when the flag is off.
* Text table  — appends a ``meaning`` column after ``old`` (only
  when ``--with-similarity`` is also on); plain rows render ``-``.
* CSV export  — appends a 7th ``score_kind`` column without
  reordering the existing six, so positional readers keep working.

The classification matrix is exhaustive: R/scored,
R/unscored, C/scored, C/unscored, plain → label ∈
{rename-similarity, unscored, copy-similarity, unscored, None}.
"""
from __future__ import annotations

import csv
import io
import json
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()

_ChangedFileAudit = _MOD._ChangedFileAudit
_RenameSimilarity = _MOD._RenameSimilarity
_render_changed_files_audit = _MOD._render_changed_files_audit
_write_similarity_csv = _MOD._write_similarity_csv
_score_kind_for = _MOD._score_kind_for
_SIMILARITY_CSV_HEADER = _MOD._SIMILARITY_CSV_HEADER
_SIMILARITY_CSV_HEADER_LABELED = _MOD._SIMILARITY_CSV_HEADER_LABELED


def _sample_rows() -> list:
    """One row per cell of the (kind, scored?) matrix + a plain row."""
    return [
        # R, scored=92 → rename-similarity
        _ChangedFileAudit(
            path="docs/new.md", status="matched", reason="renamed",
            similarity=_RenameSimilarity(kind="R", score=92,
                                         old_path="docs/old.md"),
        ),
        # R, unscored → unscored
        _ChangedFileAudit(
            path="docs/moved.md", status="matched",
            reason="renamed (no score)",
            similarity=_RenameSimilarity(kind="R", score=None,
                                         old_path="docs/src.md"),
        ),
        # C, scored=75 → copy-similarity
        _ChangedFileAudit(
            path="docs/copy.md", status="matched", reason="copied",
            similarity=_RenameSimilarity(kind="C", score=75,
                                         old_path="docs/orig.md"),
        ),
        # C, scored=0 → copy-similarity (zero is meaningful, not unscored)
        _ChangedFileAudit(
            path="docs/zero.md", status="matched", reason="copied dissim",
            similarity=_RenameSimilarity(kind="C", score=0,
                                         old_path="docs/seed.md"),
        ),
        # plain A row → no label (similarity is None)
        _ChangedFileAudit(
            path="docs/added.md", status="matched", reason="added",
        ),
    ]


class TestScoreKindClassifier(unittest.TestCase):
    def test_rename_scored_maps_to_rename_similarity(self) -> None:
        sim = _RenameSimilarity(kind="R", score=87, old_path="a")
        self.assertEqual(_score_kind_for(sim), "rename-similarity")

    def test_copy_scored_maps_to_copy_similarity(self) -> None:
        sim = _RenameSimilarity(kind="C", score=50, old_path="a")
        self.assertEqual(_score_kind_for(sim), "copy-similarity")

    def test_zero_is_scored_not_unscored(self) -> None:
        # The zero-vs-null distinction must survive through the
        # classifier — score=0 is a legitimate similarity rating.
        sim_r = _RenameSimilarity(kind="R", score=0, old_path="a")
        sim_c = _RenameSimilarity(kind="C", score=0, old_path="a")
        self.assertEqual(_score_kind_for(sim_r), "rename-similarity")
        self.assertEqual(_score_kind_for(sim_c), "copy-similarity")

    def test_none_score_is_unscored_regardless_of_kind(self) -> None:
        for kind in ("R", "C"):
            with self.subTest(kind=kind):
                sim = _RenameSimilarity(kind=kind, score=None,
                                        old_path="a")
                self.assertEqual(_score_kind_for(sim), "unscored")

    def test_plain_row_returns_none(self) -> None:
        # Distinct from the 'unscored' label: plain rows have no
        # rename provenance at all, so there's nothing to discriminate.
        self.assertIsNone(_score_kind_for(None))


class TestJsonAuditScoreKind(unittest.TestCase):
    def _render(self, with_labels: bool) -> list:
        buf = io.StringIO()
        _render_changed_files_audit(
            _sample_rows(), buf,
            as_json=True, with_similarity=True, with_labels=with_labels,
        )
        return json.loads(buf.getvalue())

    def test_legacy_schema_omits_score_kind_when_flag_off(self) -> None:
        payload = self._render(with_labels=False)
        for obj in payload:
            sim = obj.get("similarity")
            if sim is not None:
                self.assertNotIn("score_kind", sim)

    def test_label_added_to_each_kind_when_on(self) -> None:
        payload = self._render(with_labels=True)
        # Index aligns with _sample_rows().
        self.assertEqual(payload[0]["similarity"]["score_kind"],
                         "rename-similarity")
        self.assertEqual(payload[1]["similarity"]["score_kind"],
                         "unscored")
        self.assertEqual(payload[2]["similarity"]["score_kind"],
                         "copy-similarity")
        self.assertEqual(payload[3]["similarity"]["score_kind"],
                         "copy-similarity")
        # Plain row carries similarity=None; the discriminator is
        # naturally absent (no sub-object to extend).
        self.assertIsNone(payload[4]["similarity"])

    def test_existing_similarity_fields_unchanged(self) -> None:
        # Adding the discriminator must not perturb the established
        # kind/score/old_path triple.
        payload = self._render(with_labels=True)
        self.assertEqual(payload[0]["similarity"]["kind"], "R")
        self.assertEqual(payload[0]["similarity"]["score"], 92)
        self.assertEqual(payload[0]["similarity"]["old_path"],
                         "docs/old.md")


class TestTextTableMeaningColumn(unittest.TestCase):
    def _render(self, *, with_labels: bool) -> str:
        buf = io.StringIO()
        _render_changed_files_audit(
            _sample_rows(), buf,
            as_json=False, with_similarity=True, with_labels=with_labels,
        )
        return buf.getvalue()

    def test_meaning_column_absent_by_default(self) -> None:
        out = self._render(with_labels=False)
        # The column header isn't printed and no labels leak into rows.
        self.assertNotIn("meaning", out)
        self.assertNotIn("rename-similarity", out)
        self.assertNotIn("copy-similarity", out)

    def test_meaning_column_renders_with_per_kind_labels(self) -> None:
        out = self._render(with_labels=True)
        self.assertIn("meaning", out)
        self.assertIn("rename-similarity", out)
        self.assertIn("copy-similarity", out)
        self.assertIn("unscored", out)
        # Header annotation calls out the new column so reviewers
        # don't think the wider table is a layout regression.
        self.assertIn("+similarity columns +meaning", out)

    def test_plain_row_meaning_renders_as_dash(self) -> None:
        out = self._render(with_labels=True)
        # The plain row's path shouldn't carry any of the labels.
        plain_line = next(
            line for line in out.splitlines() if "docs/added.md" in line
        )
        self.assertNotIn("rename-similarity", plain_line)
        self.assertNotIn("copy-similarity", plain_line)
        self.assertNotIn("unscored", plain_line)
        # The blank cell uses the same single-dash sentinel as the
        # other similarity columns for visual consistency.
        self.assertIn(" - ", plain_line)


class TestCsvLabeledExport(unittest.TestCase):
    def _write(self, *, with_labels: bool) -> list[list[str]]:
        with tempfile.TemporaryDirectory() as d:
            out = Path(d) / "audit.csv"
            _write_similarity_csv(_sample_rows(), str(out),
                                  with_labels=with_labels)
            return list(csv.reader(io.StringIO(
                out.read_text(encoding="utf-8"))))

    def test_header_extends_not_replaces(self) -> None:
        # Position invariant: the seven-column header must equal the
        # six-column header plus one suffix — never a reorder.
        self.assertEqual(
            _SIMILARITY_CSV_HEADER_LABELED[:6], _SIMILARITY_CSV_HEADER,
        )
        self.assertEqual(_SIMILARITY_CSV_HEADER_LABELED[6], "score_kind")

    def test_six_column_csv_when_labels_off(self) -> None:
        grid = self._write(with_labels=False)
        self.assertEqual(grid[0], list(_SIMILARITY_CSV_HEADER))
        for row in grid[1:]:
            self.assertEqual(len(row), 6)

    def test_seventh_column_carries_per_kind_label(self) -> None:
        grid = self._write(with_labels=True)
        self.assertEqual(grid[0], list(_SIMILARITY_CSV_HEADER_LABELED))
        # Map path → label so the assertions are order-independent.
        labels = {row[0]: row[6] for row in grid[1:]}
        self.assertEqual(labels["docs/new.md"], "rename-similarity")
        self.assertEqual(labels["docs/moved.md"], "unscored")
        self.assertEqual(labels["docs/copy.md"], "copy-similarity")
        self.assertEqual(labels["docs/zero.md"], "copy-similarity")
        # Plain row gets an empty cell — same convention as the other
        # similarity columns when no provenance is attached.
        self.assertEqual(labels["docs/added.md"], "")

    def test_existing_columns_unchanged_by_label_flag(self) -> None:
        # Re-exporting with labels on must not perturb the prior six
        # columns — only append. Compare them position-by-position.
        without = self._write(with_labels=False)
        with_l = self._write(with_labels=True)
        for plain_row, labeled_row in zip(without[1:], with_l[1:]):
            self.assertEqual(plain_row, labeled_row[:6])


if __name__ == "__main__":
    unittest.main()