"""Cross-surface parity for ``ignored-deleted`` rows.

The audit trail surfaces three rendered representations from the
same in-memory row set:

* a human-readable text table on STDERR (default),
* a ``--json`` array on STDERR (machine schema, dashboards),
* a ``--similarity-csv`` export to disk (spreadsheet review).

An ``ignored-deleted`` row identified by its ``(path, status)``
pair MUST carry the same ``reason`` string in all three. Drift
between surfaces would let a CI dashboard show one explanation
while the operator scrolling logs sees another — exactly the kind
of silent inconsistency the audit is supposed to prevent.

Sibling test ``test_ignored_deleted_reason_exact_match.py``
already checks each surface's reasons land in the canonical
vocabulary allow-list and that the ``(path, reason)`` *sets*
match. This file is the per-pair contract: for every
``(path, status)`` key, the three surfaces must agree on the
exact ``reason`` value — not just on the unordered set.

The fixture deliberately includes:

* multiple deletes from different provenance tags
  (``changed-files-D`` via ``D\\tpath`` payload entries) so the
  per-pair check has more than one row to compare;
* a non-deleted row (``matched``) that the parity check ignores —
  it confirms the surface filters work and that adding non-delete
  rows to the audit doesn't perturb the deleted-row parity;
* a deleted ``path`` whose textual content contains spaces in
  *neighbouring* paths, ensuring the text-table regex doesn't
  misalign columns when adjacent rows have different widths.
"""
from __future__ import annotations

import csv
import json
import re
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

_MOD = load_placeholder_linter()
_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")

# Multi-row fixture: three deletes (so per-pair parity has real
# work to do) + one matched row (so the per-surface deleted filter
# is exercised against non-target rows too). Path widths differ so
# the text-table regex must handle variable column padding.
_PAYLOAD = (
    "spec/present.md\n"             # matched   — ignored by parity check
    "D\tspec/a.md\n"                # deleted   — short path
    "D\tspec/much-longer-name.md\n" # deleted   — wider column
    "D\tspec/b.md\n"                # deleted   — short again
)

# The five known audit statuses; used to build a status-aware
# regex that won't accidentally match header / footer / separator
# lines in the text rendering.
_KNOWN_STATUSES: tuple[str, ...] = (
    "matched", "ignored-extension", "ignored-out-of-root",
    "ignored-missing", "ignored-deleted",
)
_STATUS_ALT = "|".join(re.escape(s) for s in _KNOWN_STATUSES)

# Text-table row shape (matches what `--list-changed-files` prints):
#   <leading-ws> <status> <ws> <path> <ws> <reason...>
# Anchor on the closed status set so we don't capture banner /
# totals / separator lines by accident.
_TEXT_ROW_RE = re.compile(
    rf"^\s*(?P<status>{_STATUS_ALT})\s+(?P<path>\S+)\s+"
    rf"(?P<reason>.+?)\s*$")


class _Sandbox:
    """Temp repo with a `spec/present.md` file (so the matched row
    in the payload resolves to a real on-disk path) and a writable
    `--changed-files` payload. The CLI is invoked from `cwd=root`
    so the relative payload paths resolve underneath the sandbox.
    """

    def __init__(self, payload: str) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.root = Path(self._tmp.name).resolve()
        self.spec = self.root / "spec"
        self.spec.mkdir()
        (self.spec / "present.md").write_text("# present\n",
                                              encoding="utf-8")
        self.changed = self.root / "changed.txt"
        self.changed.write_text(payload, encoding="utf-8")

    def __enter__(self) -> "_Sandbox":
        return self

    def __exit__(self, *exc: object) -> None:
        self._tmp.cleanup()

    def run(self, *extra: str) -> subprocess.CompletedProcess:
        cmd = [
            sys.executable, str(_LINTER),
            "--root", str(self.spec),
            "--changed-files", str(self.changed),
            "--list-changed-files",
            *extra,
        ]
        return subprocess.run(cmd, capture_output=True, text=True,
                              check=False, cwd=str(self.root))


def _text_deleted_pairs(stderr: str) -> dict[tuple[str, str], str]:
    """Extract the ``{(path, status): reason}`` map of
    ``ignored-deleted`` rows from the rendered text table.

    Returns a mapping (not a set) so the per-pair comparison can
    diff *values* per key and report which key disagrees, rather
    than just "the sets differ".
    """
    out: dict[tuple[str, str], str] = {}
    for line in stderr.splitlines():
        m = _TEXT_ROW_RE.match(line)
        if m and m.group("status") == "ignored-deleted":
            out[(m.group("path"), m.group("status"))] = m.group("reason")
    return out


def _json_deleted_pairs(audit: list[dict]) -> dict[tuple[str, str], str]:
    """Same shape as :func:`_text_deleted_pairs`, sourced from the
    parsed ``--json`` audit array.
    """
    return {
        (r["path"], r["status"]): r["reason"]
        for r in audit if r["status"] == "ignored-deleted"
    }


def _csv_deleted_pairs(rows: list[dict]) -> dict[tuple[str, str], str]:
    """Same shape as the helpers above, sourced from the parsed
    ``--similarity-csv`` rows.
    """
    return {
        (r["path"], r["status"]): r["reason"]
        for r in rows if r["status"] == "ignored-deleted"
    }


class IgnoredDeletedReasonsAgreeAcrossSurfaces(unittest.TestCase):
    """For every ``(path, status)`` key in the ``ignored-deleted``
    bucket, the ``reason`` string must be byte-identical across
    text, JSON, and CSV renderings of the same payload.
    """

    @classmethod
    def setUpClass(cls) -> None:
        # Run all three surfaces once and parse the deleted-row
        # maps. Done in setUpClass so the subprocess + temp-dir
        # cost is paid once across the assertions below.
        with _Sandbox(_PAYLOAD) as box:
            text_proc = box.run()
            json_proc = box.run("--json")
            csv_path = box.root / "audit.csv"
            csv_proc = box.run("--similarity-csv", str(csv_path))
            with csv_path.open(encoding="utf-8", newline="") as fh:
                csv_rows = list(csv.DictReader(fh))

        # Stash exit context so failure messages below are
        # self-contained (the temp dir is gone by the time tests
        # actually run their assertions).
        cls.text_stderr = text_proc.stderr
        cls.text_rc = text_proc.returncode
        cls.json_stderr = json_proc.stderr
        cls.json_rc = json_proc.returncode
        cls.csv_rc = csv_proc.returncode
        cls.csv_stderr = csv_proc.stderr

        cls.text_pairs = _text_deleted_pairs(text_proc.stderr)
        cls.json_pairs = _json_deleted_pairs(
            json.loads(json_proc.stderr.strip()))
        cls.csv_pairs = _csv_deleted_pairs(csv_rows)

    # --- Sanity guards (so a degraded fixture fails loud) --------

    def test_all_three_surfaces_exited_clean(self) -> None:
        # If any surface non-zero'd, downstream parity assertions
        # would either KeyError or compare against a partial map —
        # surface the exit failure first.
        self.assertEqual(self.text_rc, 0,
                         msg=f"text run failed:\n{self.text_stderr}")
        self.assertEqual(self.json_rc, 0,
                         msg=f"json run failed:\n{self.json_stderr}")
        self.assertEqual(self.csv_rc, 0,
                         msg=f"csv run failed:\n{self.csv_stderr}")

    def test_each_surface_reports_three_deleted_rows(self) -> None:
        # The fixture contains exactly three `D\tpath` rows. If a
        # surface drops one (rendering bug, dedupe regression),
        # the per-pair diff below would still pass on the
        # remaining intersection — pin the count first.
        self.assertEqual(len(self.text_pairs), 3,
                         msg=f"text deleted rows: {self.text_pairs!r}")
        self.assertEqual(len(self.json_pairs), 3,
                         msg=f"json deleted rows: {self.json_pairs!r}")
        self.assertEqual(len(self.csv_pairs), 3,
                         msg=f"csv deleted rows: {self.csv_pairs!r}")

    # --- Key-set agreement (path, status) ------------------------

    def test_path_status_keys_match_text_vs_json(self) -> None:
        self.assertEqual(
            set(self.text_pairs), set(self.json_pairs),
            msg=f"text-only keys: {set(self.text_pairs) - set(self.json_pairs)}\n"
                f"json-only keys: {set(self.json_pairs) - set(self.text_pairs)}")

    def test_path_status_keys_match_json_vs_csv(self) -> None:
        self.assertEqual(
            set(self.json_pairs), set(self.csv_pairs),
            msg=f"json-only keys: {set(self.json_pairs) - set(self.csv_pairs)}\n"
                f"csv-only keys: {set(self.csv_pairs) - set(self.json_pairs)}")

    # --- Per-pair value parity (the actual contract) -------------

    def test_reason_matches_text_vs_json_per_pair(self) -> None:
        # Iterate the union so a missing key in either surface
        # surfaces with a clear "key X is in <surface> but not
        # <other>" failure rather than a silent set-diff.
        for key in set(self.text_pairs) | set(self.json_pairs):
            with self.subTest(pair=key):
                self.assertIn(key, self.text_pairs,
                              msg=f"key {key!r} missing from text")
                self.assertIn(key, self.json_pairs,
                              msg=f"key {key!r} missing from JSON")
                self.assertEqual(
                    self.text_pairs[key], self.json_pairs[key],
                    msg=f"reason drift for {key!r}:\n"
                        f"  text: {self.text_pairs[key]!r}\n"
                        f"  json: {self.json_pairs[key]!r}")

    def test_reason_matches_json_vs_csv_per_pair(self) -> None:
        for key in set(self.json_pairs) | set(self.csv_pairs):
            with self.subTest(pair=key):
                self.assertIn(key, self.json_pairs,
                              msg=f"key {key!r} missing from JSON")
                self.assertIn(key, self.csv_pairs,
                              msg=f"key {key!r} missing from CSV")
                self.assertEqual(
                    self.json_pairs[key], self.csv_pairs[key],
                    msg=f"reason drift for {key!r}:\n"
                        f"  json: {self.json_pairs[key]!r}\n"
                        f"  csv:  {self.csv_pairs[key]!r}")

    def test_reason_matches_text_vs_csv_per_pair(self) -> None:
        # Transitive equality is implied by the two tests above,
        # but pin the direct text↔CSV link explicitly so a future
        # refactor that replaces text-vs-JSON with a separate
        # canonicaliser doesn't silently let text and CSV drift.
        for key in set(self.text_pairs) | set(self.csv_pairs):
            with self.subTest(pair=key):
                self.assertIn(key, self.text_pairs,
                              msg=f"key {key!r} missing from text")
                self.assertIn(key, self.csv_pairs,
                              msg=f"key {key!r} missing from CSV")
                self.assertEqual(
                    self.text_pairs[key], self.csv_pairs[key],
                    msg=f"reason drift for {key!r}:\n"
                        f"  text: {self.text_pairs[key]!r}\n"
                        f"  csv:  {self.csv_pairs[key]!r}")

    # --- Whole-map equality (one assertion, full diff) -----------

    def test_full_pair_to_reason_map_is_identical_across_surfaces(
            self) -> None:
        # Final belt-and-braces: equality between dicts gives
        # unittest's diff a chance to render the full map mismatch
        # in one go, complementing the per-pair subTests above.
        self.assertEqual(self.text_pairs, self.json_pairs,
                         msg="text vs JSON full-map drift")
        self.assertEqual(self.json_pairs, self.csv_pairs,
                         msg="json vs CSV full-map drift")


if __name__ == "__main__":
    unittest.main()