"""End-to-end parity test: STDERR audit ↔ CSV export.

Drives the linter via subprocess with a hand-crafted
``--changed-files`` payload that exercises every shape relevant
to the CSV's score-cell semantics:

* Plain matched row (no rename signal → empty kind/score/old_path).
* Scored rename ``R092\\told\\tnew`` (integer score → ``"92"``).
* Observed-dissimilar rename ``R000\\told\\tnew`` (score=0 → the
  literal string ``"0"``, **not** an empty cell).
* Unscored rename ``R\\told\\tnew`` (score=None → an EMPTY cell).
* Ignored-deleted row ``D\\tpath`` (no post-state → empty cells).

It then re-reads the generated CSV with the stdlib ``csv`` module
and asserts:

1. **Row-count parity** — the CSV body has the same number of
   rows as the JSON audit emitted on STDERR (post-dedupe, post-
   filter), so neither stream is silently dropping or fabricating
   rows.
2. **Per-row alignment** — for every (path, status) pair on
   STDERR there's a matching CSV body row.
3. **Empty-vs-0 score semantics** — the CSV's `score` cell is
   the empty string for unscored / plain rows and the literal
   ``"0"`` for the observed-dissimilar case. ``ISBLANK`` and
   ``=0`` are not the same condition; this test pins that
   distinction so a future refactor can't blur it.

Hermetic — no real git, just a tempdir with a ``--changed-files``
manifest. Runs the published CLI exactly as a CI job would.
"""
from __future__ import annotations

import csv
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")


def _write_clean_md(p: Path) -> None:
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text("# Title\n\nNo placeholders.\n", encoding="utf-8")


def _setup_intake(root: Path) -> Path:
    """Build a tree + manifest covering every score-cell semantic.

    Files actually present on disk so they can resolve to ``matched``:

    * ``spec/plain.md``  — plain A row (no rename signal).
    * ``spec/new1.md``   — destination of a scored rename.
    * ``spec/new2.md``   — destination of a score-0 rename.
    * ``spec/new3.md``   — destination of an UNSCORED rename.

    The manifest also includes a ``D\\t...`` row for an absent
    path so we get an ``ignored-deleted`` audit entry whose CSV
    cells must all be empty.
    """
    spec = root / "spec"
    _write_clean_md(spec / "plain.md")
    _write_clean_md(spec / "new1.md")
    _write_clean_md(spec / "new2.md")
    _write_clean_md(spec / "new3.md")
    manifest = root / "changed.txt"
    manifest.write_text(
        "spec/plain.md\n"                       # matched, no rename
        "R092\tspec/old1.md\tspec/new1.md\n"   # scored: score="92"
        "R000\tspec/old2.md\tspec/new2.md\n"   # observed dissim: "0"
        "R\tspec/old3.md\tspec/new3.md\n"      # unscored: empty cell
        "D\tspec/gone.md\n",                    # ignored-deleted
        encoding="utf-8",
    )
    return manifest


class TestCsvStderrParity(unittest.TestCase):

    def _run(self, root: Path, csv_path: Path,
             *extra: str) -> tuple[int, str, str]:
        manifest = _setup_intake(root)
        cp = subprocess.run(
            [
                sys.executable, str(_LINTER),
                "--root", "spec",
                "--repo-root", str(root),
                "--changed-files", str(manifest),
                "--list-changed-files",
                "--with-similarity",
                "--similarity-csv", str(csv_path),
                "--json",
                *extra,
            ],
            cwd=root, capture_output=True, text=True,
        )
        return cp.returncode, cp.stdout, cp.stderr

    # -------------------------------------------------------------
    # 1. Row-count parity + per-row (path, status) alignment.
    # -------------------------------------------------------------

    def test_csv_row_count_matches_stderr_audit(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            csv_path = root / "audit.csv"
            rc, _out, err = self._run(root, csv_path)
            self.assertEqual(rc, 0, msg=f"stderr={err!r}")

            audit = json.loads(err)
            # 4 matched + 1 ignored-deleted = 5 audit rows.
            self.assertEqual(len(audit), 5)

            with csv_path.open(encoding="utf-8", newline="") as fh:
                grid = list(csv.reader(fh))
            header, body = grid[0], grid[1:]
            self.assertEqual(header, [
                "path", "status", "reason",
                "kind", "score", "old_path",
            ])
            self.assertEqual(
                len(body), len(audit),
                msg=(
                    f"CSV body has {len(body)} rows but STDERR audit "
                    f"has {len(audit)} — the two streams must agree "
                    f"on row count post-dedupe / post-filter."
                ),
            )

    def test_csv_path_status_pairs_match_stderr(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            csv_path = root / "audit.csv"
            _rc, _out, err = self._run(root, csv_path)

            audit = json.loads(err)
            stderr_pairs = sorted(
                (r["path"], r["status"]) for r in audit
            )

            with csv_path.open(encoding="utf-8", newline="") as fh:
                grid = list(csv.reader(fh))
            csv_pairs = sorted(
                (row[0], row[1]) for row in grid[1:]
            )
            self.assertEqual(csv_pairs, stderr_pairs)

    # -------------------------------------------------------------
    # 2. Score cell semantics: empty vs "0" vs integer.
    # -------------------------------------------------------------

    def test_score_empty_vs_zero_semantics_match_stderr(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            csv_path = root / "audit.csv"
            _rc, _out, err = self._run(root, csv_path)

            audit = json.loads(err)
            with csv_path.open(encoding="utf-8", newline="") as fh:
                grid = list(csv.reader(fh))
            csv_by_path = {row[0]: row for row in grid[1:]}
            audit_by_path = {r["path"]: r for r in audit}

            self.assertEqual(set(csv_by_path), set(audit_by_path))

            # Walk every row and assert the CSV `score` cell encodes
            # the same three-way distinction as the JSON
            # `similarity.score` field:
            #
            #   JSON null  → CSV "" (empty cell)
            #   JSON 0     → CSV "0" (literal zero, NOT empty)
            #   JSON N>0   → CSV str(N)
            #
            # And rows without a `similarity` object at all (plain
            # A/M/D) must surface an empty CSV `score` cell too, but
            # with the kind / old_path columns ALSO empty so the
            # "no rename provenance" case is uniform.
            for path, csv_row in csv_by_path.items():
                _, _, _, csv_kind, csv_score, csv_old = csv_row
                with self.subTest(path=path):
                    sim = audit_by_path[path].get("similarity")
                    if sim is None:
                        self.assertEqual(csv_kind, "")
                        self.assertEqual(csv_score, "")
                        self.assertEqual(csv_old, "")
                        continue
                    # Rename row — kind + old_path must be populated.
                    self.assertEqual(csv_kind, sim["kind"])
                    self.assertEqual(csv_old, sim["old_path"])
                    json_score = sim["score"]
                    if json_score is None:
                        self.assertEqual(
                            csv_score, "",
                            msg=("unscored rename must surface as an "
                                 "EMPTY CSV cell — distinguishes it "
                                 "from observed-dissimilar (`0`)"),
                        )
                    elif json_score == 0:
                        self.assertEqual(
                            csv_score, "0",
                            msg=("observed-dissimilar must surface "
                                 "as the literal string `0`, NOT "
                                 "empty — `=0` and `ISBLANK` are "
                                 "not the same condition"),
                        )
                    else:
                        self.assertEqual(csv_score, str(json_score))

    def test_csv_contains_both_empty_and_zero_score_cells(self) -> None:
        # Belt-and-braces: even if the per-row loop above were to
        # silently skip a case, this test fails loudly when the
        # CSV ever loses the empty-vs-"0" distinction across the
        # five-row fixture.
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            csv_path = root / "audit.csv"
            self._run(root, csv_path)
            with csv_path.open(encoding="utf-8", newline="") as fh:
                grid = list(csv.reader(fh))
            score_cells = [row[4] for row in grid[1:]]
            self.assertIn("", score_cells,
                          msg="fixture lost its unscored / plain rows")
            self.assertIn("0", score_cells,
                          msg="fixture lost its observed-dissimilar row")
            self.assertIn("92", score_cells,
                          msg="fixture lost its scored rename row")

    # -------------------------------------------------------------
    # 3. Parity holds under dedupe + filter (the audit pipeline
    #    runs BEFORE the export, so both streams must shrink in
    #    lock-step).
    # -------------------------------------------------------------

    def test_parity_holds_under_only_changed_status_filter(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            csv_path = root / "audit.csv"
            _rc, _out, err = self._run(
                root, csv_path,
                "--only-changed-status", "matched",
            )
            audit = json.loads(err)
            # Filter drops the ignored-deleted row.
            self.assertEqual({r["status"] for r in audit}, {"matched"})
            self.assertEqual(len(audit), 4)

            with csv_path.open(encoding="utf-8", newline="") as fh:
                grid = list(csv.reader(fh))
            body = grid[1:]
            self.assertEqual(len(body), len(audit))
            self.assertTrue(all(row[1] == "matched" for row in body),
                            msg=("filter must apply to the CSV too — "
                                 "the export mirrors STDERR"))


if __name__ == "__main__":
    unittest.main()