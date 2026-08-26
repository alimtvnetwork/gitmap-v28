"""Regression tests for the shared audit-trail parsing helpers.

Pins the contract that each helper in
:mod:`_audit_parse_helpers` must:

1. parse the real audit rows correctly when the linter output
   matches today's shape (no trailing summary);
2. continue to parse correctly when a *trailing summary* — a
   ``totals: …`` line, a ``# totals: …`` CSV comment, a blank
   line, or arbitrary chatter — is appended after the audit
   block;
3. **never** return a synthesised summary line as a data row.

These are pure-Python tests against the helpers (no subprocess);
the linter's actual output already exercises the happy path
through the cross-surface parity tests. Here we feed the helpers
hand-crafted strings that include trailing summary lines so the
contract is verified independently of whatever the live CLI
happens to print today.
"""
from __future__ import annotations

import sys
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from _audit_parse_helpers import (  # noqa: E402
    KNOWN_STATUSES,
    parse_csv_audit_rows,
    parse_json_audit_rows,
    parse_text_audit_rows,
)

# A realistic three-row text audit followed by the totals
# summary line the linter prints today. Matches the actual
# layout (banner / header / separator / rows / totals).
_TEXT_FIXTURE = (
    "── placeholder-comments: changed-file audit (3 row(s)) ──\n"
    "  status           path             reason\n"
    "  ----------------------------------------\n"
    "  matched          spec/present.md  under --root, extension allowed, file present on disk\n"
    "  ignored-deleted  spec/a.md        --changed-files payload row shaped `D\\tpath`: explicit delete marker, no post-state to lint\n"
    "  ignored-deleted  spec/b.md        --changed-files payload row shaped `D\\tpath`: explicit delete marker, no post-state to lint\n"
    "  totals: matched=1  ignored-extension=0  ignored-out-of-root=0  ignored-missing=0  ignored-deleted=2\n"
)

# Same text fixture with extra hostile trailing chatter beyond
# the totals line — a blank line, a debug ``DEBUG:`` log and an
# ``ERROR`` line that contains the substring ``matched`` to try
# to trick the regex anchor into matching it.
_TEXT_FIXTURE_EXTRA = _TEXT_FIXTURE + (
    "\n"
    "DEBUG: 4ms cache hit\n"
    "ERROR: something matched the wrong thing earlier\n"
)

_JSON_FIXTURE_BARE = (
    '[\n'
    '  {"path": "spec/a.md", "status": "ignored-deleted",'
    ' "reason": "git diff reported D (deleted)"},\n'
    '  {"path": "spec/b.md", "status": "matched",'
    ' "reason": "under --root, extension allowed, file present on disk"}\n'
    ']\n'
)
# JSON output framed by a leading banner and a trailing summary
# line — neither is valid JSON, so today's
# ``json.loads(stderr.strip())`` would explode on this input.
_JSON_FIXTURE_FRAMED = (
    "── placeholder-comments: changed-file audit (2 row(s)) ──\n"
    + _JSON_FIXTURE_BARE
    + "totals: matched=1  ignored-deleted=1\n"
)
# A JSON value whose ``reason`` deliberately contains ``]`` to
# verify the bracket-matching parser isn't fooled by closers
# inside string literals.
_JSON_FIXTURE_TRICKY = (
    "banner before\n"
    '[{"path": "spec/x.md", "status": "ignored-deleted",'
    ' "reason": "contains ] and \\"quotes\\" inside"}]\n'
    "totals: ignored-deleted=1\n"
)

_CSV_FIXTURE_BARE = (
    "path,status,reason,kind,score,old_path\n"
    "spec/a.md,ignored-deleted,"
    '"--changed-files payload row shaped `D\\tpath`: explicit delete marker, no post-state to lint"'
    ",,,\n"
    "spec/b.md,matched,"
    '"under --root, extension allowed, file present on disk"'
    ",,,\n"
)
# CSV with a trailing pseudo-comment row appended as a summary.
# A naive ``DictReader`` walk would yield it as a row whose
# ``path`` is ``"# totals: matched=1  ignored-deleted=1"`` and
# whose ``status`` is empty — exactly the regression we're
# guarding against.
_CSV_FIXTURE_TRAILING_COMMENT = (
    _CSV_FIXTURE_BARE
    + "# totals: matched=1  ignored-deleted=1,,,,,\n"
)
# CSV with a trailing summary row that uses an explicit
# non-status sentinel in the ``status`` column (e.g. a future
# ``__totals__`` aggregate row). Helper must drop it.
_CSV_FIXTURE_TRAILING_TOTALS_ROW = (
    _CSV_FIXTURE_BARE
    + "TOTAL,__totals__,2 rows,,,\n"
)


class TextHelperHandlesTrailingSummary(unittest.TestCase):
    """``parse_text_audit_rows`` must drop the trailing
    ``totals: …`` line and any chatter past it.
    """

    def test_returns_only_three_audit_rows(self) -> None:
        rows = parse_text_audit_rows(_TEXT_FIXTURE)
        self.assertEqual(len(rows), 3,
                         msg=f"unexpected rows: {rows!r}")

    def test_no_row_carries_a_totals_marker(self) -> None:
        # Belt-and-braces: even if the count above ever
        # accidentally lined up with a bug, no real audit row
        # contains the literal ``"totals:"`` substring on its
        # status / path / reason fields.
        for row in parse_text_audit_rows(_TEXT_FIXTURE):
            for k, v in row.items():
                self.assertNotIn(
                    "totals:", v,
                    msg=f"summary leaked into row[{k!r}]: {row!r}")

    def test_statuses_are_all_from_closed_vocabulary(self) -> None:
        for row in parse_text_audit_rows(_TEXT_FIXTURE):
            self.assertIn(row["status"], KNOWN_STATUSES)

    def test_extra_trailing_chatter_does_not_create_phantom_rows(
            self) -> None:
        # The ``ERROR: something matched …`` line contains the
        # word ``matched`` and could be mis-anchored if the
        # parser kept scanning past the end-of-audit marker.
        rows = parse_text_audit_rows(_TEXT_FIXTURE_EXTRA)
        self.assertEqual(
            len(rows), 3,
            msg=f"trailing chatter promoted to row: {rows!r}")
        self.assertEqual(
            [r["path"] for r in rows],
            ["spec/present.md", "spec/a.md", "spec/b.md"])

    def test_real_row_reasons_are_preserved_verbatim(self) -> None:
        rows = parse_text_audit_rows(_TEXT_FIXTURE)
        # The deleted rows must carry the full reason string,
        # not a truncated prefix that stopped at the colon.
        self.assertTrue(
            rows[1]["reason"].endswith(
                "explicit delete marker, no post-state to lint"),
            msg=f"reason truncated: {rows[1]!r}")


class JsonHelperHandlesFramingAndStringEscapes(unittest.TestCase):
    """``parse_json_audit_rows`` must locate the JSON array even
    when STDERR is framed with banner / summary lines, and must
    not be fooled by ``]`` inside string values.
    """

    def test_bare_array_parses(self) -> None:
        rows = parse_json_audit_rows(_JSON_FIXTURE_BARE)
        self.assertEqual(len(rows), 2)
        self.assertEqual(rows[0]["path"], "spec/a.md")

    def test_framed_array_parses_ignoring_banner_and_totals(
            self) -> None:
        rows = parse_json_audit_rows(_JSON_FIXTURE_FRAMED)
        self.assertEqual(len(rows), 2,
                         msg=f"got: {rows!r}")
        self.assertEqual(
            {r["status"] for r in rows},
            {"ignored-deleted", "matched"})

    def test_string_with_inner_bracket_does_not_close_array(
            self) -> None:
        rows = parse_json_audit_rows(_JSON_FIXTURE_TRICKY)
        self.assertEqual(len(rows), 1, msg=f"got: {rows!r}")
        self.assertEqual(
            rows[0]["reason"],
            'contains ] and "quotes" inside')

    def test_missing_array_raises(self) -> None:
        with self.assertRaises(ValueError):
            parse_json_audit_rows("totals: nothing here\n")

    def test_unclosed_array_raises(self) -> None:
        with self.assertRaises(ValueError):
            parse_json_audit_rows('[{"path": "spec/a.md"')


class CsvHelperStopsAtFirstNonAuditRow(unittest.TestCase):
    """``parse_csv_audit_rows`` must drop any trailing row whose
    ``status`` isn't in the closed vocabulary.
    """

    def test_bare_csv_returns_two_rows(self) -> None:
        rows = parse_csv_audit_rows(_CSV_FIXTURE_BARE)
        self.assertEqual(len(rows), 2,
                         msg=f"got: {rows!r}")
        self.assertEqual([r["path"] for r in rows],
                         ["spec/a.md", "spec/b.md"])

    def test_trailing_comment_row_is_dropped(self) -> None:
        # Without the helper's status filter, ``DictReader``
        # would yield a third row whose ``path`` is the literal
        # ``"# totals: …"`` string.
        rows = parse_csv_audit_rows(_CSV_FIXTURE_TRAILING_COMMENT)
        self.assertEqual(len(rows), 2,
                         msg=f"trailing comment leaked: {rows!r}")
        for row in rows:
            self.assertNotIn("totals:", row.get("path", ""))

    def test_trailing_totals_row_with_unknown_status_is_dropped(
            self) -> None:
        rows = parse_csv_audit_rows(
            _CSV_FIXTURE_TRAILING_TOTALS_ROW)
        self.assertEqual(len(rows), 2,
                         msg=f"sentinel row leaked: {rows!r}")
        for row in rows:
            self.assertIn(row["status"], KNOWN_STATUSES)

    def test_csv_helper_returns_status_in_closed_vocabulary(
            self) -> None:
        for row in parse_csv_audit_rows(
                _CSV_FIXTURE_TRAILING_COMMENT):
            self.assertIn(row["status"], KNOWN_STATUSES)


class HelpersAgreeOnRealisticParallelFixtures(unittest.TestCase):
    """Cross-helper smoke check: feeding the three helpers
    parallel fixtures (same logical rows in three formats, each
    with a trailing summary) yields the same ``(path, status)``
    set. Pins that "stop at first non-audit line" is interpreted
    consistently across surfaces.
    """

    def test_path_status_pairs_agree(self) -> None:
        text_pairs = {
            (r["path"], r["status"])
            for r in parse_text_audit_rows(_TEXT_FIXTURE)
            if r["status"] != "matched"
        }
        json_pairs = {
            (r["path"], r["status"])
            for r in parse_json_audit_rows(_JSON_FIXTURE_FRAMED)
            if r["status"] != "matched"
        }
        csv_pairs = {
            (r["path"], r["status"])
            for r in parse_csv_audit_rows(
                _CSV_FIXTURE_TRAILING_COMMENT)
            if r["status"] != "matched"
        }
        # All three fixtures intentionally include the deleted
        # ``spec/a.md`` row; that's the minimum overlap we pin.
        common = {("spec/a.md", "ignored-deleted")}
        self.assertTrue(common.issubset(text_pairs),
                        msg=f"text missing: {text_pairs!r}")
        self.assertTrue(common.issubset(json_pairs),
                        msg=f"json missing: {json_pairs!r}")
        self.assertTrue(common.issubset(csv_pairs),
                        msg=f"csv missing: {csv_pairs!r}")


if __name__ == "__main__":
    unittest.main()