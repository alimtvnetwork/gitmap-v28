"""Tests for the ``--with-similarity`` audit-table flag.

Two layers under test:

* **Unit** — :func:`_fmt_similarity` and the ``_RenameSimilarity`` /
  ``_ChangedFileAudit`` integration. Renaming logic that captures
  metadata into the audit lives in three different parsers
  (``_parse_name_status`` for git output, ``_normalise_changed_lines``
  for both tab- and arrow-form ``--changed-files`` payloads); we
  exercise each so a future refactor can't quietly drop a code path.

* **CLI** — argparse wiring + renderer composition with
  ``--list-changed-files``. We invoke the linter as a subprocess
  against a hand-rolled ``--changed-files`` payload (which exercises
  the same audit code path as ``--diff-base`` without requiring a
  scratch git repo) and inspect both the text-mode columns and the
  JSON ``similarity`` sub-object.

Why the dash sentinel matters: a CI dashboard scraping "kind=R" out
of the table needs the absence of a rename to render as a stable,
greppable token rather than a blank cell that vanishes between
columns. ``-`` is single-char, never legal as a kind / score / path,
and survives ``str.split()`` cleanly. Tests assert the dash on every
blank/unknown axis (kind, score, old) independently so a partial
record stays maximally informative.
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

LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")


def _run(*args: str, cwd: Path) -> tuple[int, str, str]:
    r = subprocess.run([sys.executable, str(LINTER), *args],
                       cwd=cwd, capture_output=True, text=True)
    return r.returncode, r.stdout, r.stderr


def _make_repo(td: Path) -> Path:
    """spec/intro.md + spec/copy.md, both clean. The second file
    exists so a rename/copy payload can name it as the post-state
    side without tripping the ``ignored-missing`` audit branch."""
    spec = td / "spec"
    spec.mkdir()
    (spec / "intro.md").write_text("# spec\nplain prose.\n")
    (spec / "copy.md").write_text("# spec\nplain prose.\n")
    return spec


class FmtSimilarityUnit(unittest.TestCase):
    """:func:`_fmt_similarity` cell-by-cell contract."""

    def setUp(self) -> None:
        self.mod = load_placeholder_linter()

    def test_none_yields_three_dashes(self) -> None:
        # Plain A/M rows have no rename provenance. Every cell must
        # render as the documented ``-`` sentinel — no empty cells,
        # no ``None`` leaking into the text table.
        self.assertEqual(self.mod._fmt_similarity(None), ("-", "-", "-"))

    def test_full_record_renders_each_cell(self) -> None:
        sim = self.mod._RenameSimilarity(
            kind="R", score=87, old_path="spec/old.md")
        self.assertEqual(self.mod._fmt_similarity(sim),
                         ("R", "87", "spec/old.md"))

    def test_scoreless_record_dashes_only_score_cell(self) -> None:
        # Authored arrow-form payloads (`OLD => NEW`) carry a kind
        # and an old-path but no percentage. The kind + old must
        # survive; only the score cell falls back to the dash.
        sim = self.mod._RenameSimilarity(
            kind="R", score=None, old_path="spec/old.md")
        self.assertEqual(self.mod._fmt_similarity(sim),
                         ("R", "-", "spec/old.md"))

    def test_empty_old_path_falls_back_independently(self) -> None:
        # A pathological R/C row with an empty old slot still keeps
        # the kind and score visible — only the missing cell dashes.
        sim = self.mod._RenameSimilarity(kind="C", score=50, old_path="")
        self.assertEqual(self.mod._fmt_similarity(sim), ("C", "50", "-"))


class ParseNameStatusSimilarity(unittest.TestCase):
    """``_parse_name_status`` populates the optional out-map."""

    def setUp(self) -> None:
        self.mod = load_placeholder_linter()

    def test_rename_row_records_kind_score_old(self) -> None:
        sims: dict = {}
        out = self.mod._parse_name_status(
            "R087\tspec/old.md\tspec/new.md\n",
            similarities=sims)
        self.assertEqual(out, ["spec/new.md"])
        self.assertIn("spec/new.md", sims)
        rec = sims["spec/new.md"]
        self.assertEqual((rec.kind, rec.score, rec.old_path),
                         ("R", 87, "spec/old.md"))

    def test_copy_row_records_C_kind(self) -> None:
        sims: dict = {}
        self.mod._parse_name_status(
            "C100\tspec/src.md\tspec/dup.md\n", similarities=sims)
        self.assertEqual(sims["spec/dup.md"].kind, "C")
        self.assertEqual(sims["spec/dup.md"].score, 100)

    def test_plain_modify_row_leaves_map_empty(self) -> None:
        # M rows must NOT be recorded — the renderer treats absence
        # as "no rename signal" and prints the dash. Recording them
        # with kind="M" would muddy the cell vocabulary.
        sims: dict = {}
        self.mod._parse_name_status("M\tspec/intro.md\n", similarities=sims)
        self.assertEqual(sims, {})


class NormaliseChangedLinesSimilarity(unittest.TestCase):
    """Both tab- and arrow-form rename intakes populate the map."""

    def setUp(self) -> None:
        self.mod = load_placeholder_linter()

    def test_tab_form_score_captured(self) -> None:
        sims: dict = {}
        out = self.mod._normalise_changed_lines(
            ["R092\tspec/old.md\tspec/new.md"], similarities=sims)
        self.assertEqual(out, ["spec/new.md"])
        rec = sims["spec/new.md"]
        self.assertEqual((rec.kind, rec.score, rec.old_path),
                         ("R", 92, "spec/old.md"))

    def test_arrow_form_records_no_score_but_keeps_old(self) -> None:
        # `git status -s` short form has no percentage — we know it's
        # a rename but score is genuinely unknown. The dataclass
        # records ``None``; the renderer will dash that cell.
        sims: dict = {}
        out = self.mod._normalise_changed_lines(
            ["spec/old.md => spec/new.md"], similarities=sims)
        self.assertEqual(out, ["spec/new.md"])
        rec = sims["spec/new.md"]
        self.assertEqual((rec.kind, rec.score, rec.old_path),
                         ("R", None, "spec/old.md"))

    def test_plain_path_leaves_map_empty(self) -> None:
        sims: dict = {}
        self.mod._normalise_changed_lines(["spec/intro.md"],
                                          similarities=sims)
        self.assertEqual(sims, {})


class WithSimilarityCli(unittest.TestCase):
    """End-to-end --list-changed-files + --with-similarity composition."""

    def test_text_table_shows_extra_columns_and_dashes(self) -> None:
        # Two intake rows: a rename (R087) of spec/intro.md → spec/copy.md,
        # and a plain modify of spec/intro.md. The renamed row must
        # show kind=R / score=87 / old=spec/intro.md; the modify row
        # must show three dashes.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text(
                "R087\tspec/intro.md\tspec/copy.md\n"
                "M\tspec/intro.md\n"
            )
            code, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--with-similarity",
                cwd=tdp)
        self.assertIn(code, (0, 1))  # exit code is about lint, not audit
        # Header row carries the new columns + the suffix annotation.
        self.assertIn("kind", err)
        self.assertIn("score", err)
        self.assertIn("old", err)
        self.assertIn("+similarity columns", err)
        # Rename row: every metadata cell is populated.
        self.assertRegex(err, r"matched\s+spec/copy\.md\s+R\s+87\s+spec/intro\.md")
        # Plain modify row: triple dash sentinel between path and reason.
        self.assertRegex(err, r"matched\s+spec/intro\.md\s+-\s+-\s+-")

    def test_off_by_default_keeps_legacy_three_columns(self) -> None:
        # No --with-similarity → no kind/score/old in the header.
        # Protects the legacy log shape that downstream parsers and
        # human reviewers already depend on.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("R087\tspec/intro.md\tspec/copy.md\n")
            _, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                cwd=tdp)
        self.assertIn("status", err)
        self.assertIn("path", err)
        self.assertIn("reason", err)
        # Header line specifically — we look at the column header to
        # avoid matching the path "spec/copy.md" which contains "old"
        # transitively in some test names. The header line is the one
        # right after the "── …audit…" banner.
        header_line = next(
            (l for l in err.splitlines()
             if "status" in l and "path" in l and "reason" in l),
            "")
        self.assertNotIn("kind", header_line)
        self.assertNotIn("score", header_line)

    def test_json_includes_similarity_object_when_flag_set(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text(
                "R087\tspec/intro.md\tspec/copy.md\n"
                "M\tspec/intro.md\n"
            )
            _, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--with-similarity",
                "--json",
                cwd=tdp)
        # The audit JSON is on STDERR (STDOUT carries the lint
        # report). It's the first JSON document we can decode there.
        decoded = None
        for chunk_start in range(len(err)):
            if err[chunk_start] != "[":
                continue
            try:
                decoded = json.loads(err[chunk_start:].split("\n\n", 1)[0])
                break
            except json.JSONDecodeError as exc:
                import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
                continue
        self.assertIsNotNone(decoded, f"no JSON array on STDERR: {err!r}")
        by_path = {row["path"]: row for row in decoded}
        self.assertIn("similarity", by_path["spec/copy.md"])
        self.assertEqual(by_path["spec/copy.md"]["similarity"],
                         {"kind": "R", "score": 87,
                          "old_path": "spec/intro.md"})
        # Plain modify row: similarity present in schema, value null.
        self.assertIsNone(by_path["spec/intro.md"]["similarity"])

    def test_json_omits_similarity_key_without_flag(self) -> None:
        # Legacy schema preservation: dashboards parsing the
        # historical 3-key shape with a strict validator must keep
        # working. We assert the absence of the key, not its null-ness.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("R087\tspec/intro.md\tspec/copy.md\n")
            _, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--json",
                cwd=tdp)
        decoded = None
        for i, ch in enumerate(err):
            if ch == "[":
                try:
                    decoded = json.loads(err[i:].split("\n\n", 1)[0])
                    break
                except json.JSONDecodeError as exc:
                    import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
                    continue
        self.assertIsNotNone(decoded)
        for row in decoded:
            self.assertNotIn("similarity", row,
                             f"legacy schema must not carry similarity: {row}")


if __name__ == "__main__":  # pragma: no cover - manual runner
    unittest.main()
