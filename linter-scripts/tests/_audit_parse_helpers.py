"""Shared parsing helpers for the changed-file audit surfaces.

The placeholder-comments linter renders the same in-memory audit
through three surfaces:

* ``text``  — a banner, a header row, a separator, the data
  rows, and a trailing ``totals: …`` summary line on STDERR.
* ``--json`` — a JSON array on STDERR (currently the *only*
  thing on STDERR, but future work may add a banner or summary
  before/after the array, and the parser must not break then).
* ``--similarity-csv`` — a CSV file on disk; today it's just the
  header + rows, but a trailing ``# totals: …`` comment line
  could be appended as a summary footer in the future.

Several test modules previously rolled their own ad-hoc parsers
(``json.loads(stderr.strip())``, raw ``csv.DictReader`` over the
whole file, ``re.match`` line-by-line over STDERR). Those were
brittle:

* ``json.loads(stderr.strip())`` blows up the moment a banner or
  totals line is appended around the JSON array.
* ``csv.DictReader`` happily yields a row for any trailing
  comment / totals line because it doesn't know what an
  "audit row" looks like.
* the per-line text regex is well-anchored on the closed status
  vocabulary, but several call sites duplicated the regex with
  small drifts.

These helpers centralise the parsing and, critically, **stop at
the first non-audit line** — so any banner, blank line, or
trailing ``totals: …`` summary is treated as an end-of-audit
marker rather than data.

The companion test
``test_audit_parse_helpers_trailing_summary.py`` pins the
contract: each helper must drop trailing summary lines and
return only the real audit rows.
"""
from __future__ import annotations

import csv
import io
import json
import re
from typing import Iterable

# The closed set of audit statuses the linter emits. Keeping the
# vocabulary here lets every helper anchor on the same alternation
# and refuse to treat a header / banner / totals line as a row.
KNOWN_STATUSES: tuple[str, ...] = (
    "matched",
    "ignored-extension",
    "ignored-out-of-root",
    "ignored-missing",
    "ignored-deleted",
)
_STATUS_ALT = "|".join(re.escape(s) for s in KNOWN_STATUSES)

# Text-table row shape (matches what `--list-changed-files`
# prints on STDERR):
#   <leading-ws> <status> <ws> <path> <ws> <reason...>
# Anchored on the closed status set so banner / header /
# separator / totals lines never match.
_TEXT_ROW_RE = re.compile(
    rf"^\s*(?P<status>{_STATUS_ALT})\s+(?P<path>\S+)\s+"
    rf"(?P<reason>.+?)\s*$"
)

# A line is "audit-shaped" if it matches the row regex above.
# Anything else — banner (``── … ──``), header (``status path
# reason``), separator (``----``), blank lines, totals
# (``totals: …``), or arbitrary trailing chatter — is treated as
# end-of-audit by the streaming parser.


def parse_text_audit_rows(stderr: str) -> list[dict[str, str]]:
    """Parse the rendered text audit table on STDERR.

    Walks the lines in order and **stops at the first non-audit
    line that appears after at least one audit row has been
    consumed**. This way:

    * leading banner / header / separator lines are skipped
      (they don't match the row regex, so they're ignored
      until the first real row arrives);
    * trailing summary lines (``totals: …``, blank line,
      arbitrary footer) terminate the scan instead of being
      accidentally parsed as rows.

    Returns a list of ``{"path", "status", "reason"}`` dicts in
    the order they appeared.
    """
    rows: list[dict[str, str]] = []
    seen_first_row = False
    for line in stderr.splitlines():
        m = _TEXT_ROW_RE.match(line)
        if m is not None:
            rows.append({
                "status": m.group("status"),
                "path": m.group("path"),
                "reason": m.group("reason"),
            })
            seen_first_row = True
            continue
        # Non-matching line. Before any data row was seen we
        # treat it as preamble (banner / header / separator) and
        # keep scanning. Once data rows have started, the first
        # non-matching line is the end of the audit block —
        # anything after it is a trailing summary or unrelated
        # output and must NOT be parsed as a row.
        if seen_first_row:
            break
    return rows


def parse_json_audit_rows(stderr: str) -> list[dict]:
    """Parse the ``--json`` audit array off STDERR.

    Today STDERR is exactly the JSON array. Tomorrow it may be
    preceded by a banner line or followed by a ``totals: …``
    summary. This parser locates the first ``[`` and the
    matching closing ``]`` (respecting strings and escapes) and
    feeds *only* that slice to ``json.loads``.

    Raises ``ValueError`` if no JSON array can be located —
    callers should treat that as a hard parse failure rather
    than silently returning an empty list.
    """
    start = stderr.find("[")
    if start == -1:
        raise ValueError(
            "no JSON array found in audit STDERR; "
            f"got: {stderr!r}"
        )
    end = _find_matching_bracket(stderr, start)
    if end == -1:
        raise ValueError(
            "JSON array opened at offset "
            f"{start} but never closed in: {stderr!r}"
        )
    payload = stderr[start:end + 1]
    parsed = json.loads(payload)
    if not isinstance(parsed, list):
        raise ValueError(
            f"expected JSON array, got {type(parsed).__name__}"
        )
    return parsed


def _find_matching_bracket(text: str, start: int) -> int:
    """Return the index of the ``]`` that closes the ``[`` at
    ``text[start]``, or ``-1`` if no match.

    Respects JSON string boundaries and backslash escapes so a
    ``"]"`` inside a string value doesn't close the array
    prematurely. Doesn't validate JSON — that's ``json.loads``'s
    job — only finds the slice boundary.
    """
    assert text[start] == "[", "start must point at '['"
    depth = 0
    in_str = False
    escape = False
    for i in range(start, len(text)):
        ch = text[i]
        if in_str:
            if escape:
                escape = False
            elif ch == "\\":
                escape = True
            elif ch == '"':
                in_str = False
            continue
        if ch == '"':
            in_str = True
        elif ch == "[":
            depth += 1
        elif ch == "]":
            depth -= 1
            if depth == 0:
                return i
    return -1


def parse_csv_audit_rows(csv_text: str) -> list[dict[str, str]]:
    """Parse a ``--similarity-csv`` audit file's contents.

    Stops at the first non-audit row, defined as any row whose
    ``status`` field is not in :data:`KNOWN_STATUSES`. This
    handles trailing summary lines such as ``# totals: …``
    (which CSV would otherwise yield as a row with the literal
    ``"# totals: …"`` in the ``path`` column) or a blank
    trailing line (which ``DictReader`` skips, but a
    ``"totals,…"`` line would be consumed verbatim).

    Accepts the CSV as text (rather than a file path) so callers
    can inject synthetic trailing content in regression tests
    without writing it to disk.
    """
    reader = csv.DictReader(io.StringIO(csv_text))
    rows: list[dict[str, str]] = []
    for row in reader:
        status = (row.get("status") or "").strip()
        if status not in KNOWN_STATUSES:
            # First non-audit row terminates the scan. We do NOT
            # raise: trailing comment / summary rows are an
            # acceptable future extension — the contract is just
            # that the helper must not return them as data.
            break
        rows.append(dict(row))
    return rows


__all__: Iterable[str] = (
    "KNOWN_STATUSES",
    "parse_text_audit_rows",
    "parse_json_audit_rows",
    "parse_csv_audit_rows",
)