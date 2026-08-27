"""Universal contract: every audit-trail row carries a non-empty
``reason``.

The ``reason`` field is the single human-readable cell a CI
reviewer scans first when triaging an unfamiliar audit row — it
explains *why* the path was bucketed into its ``status`` (matched,
ignored-extension, ignored-out-of-root, ignored-missing,
ignored-deleted). A row with a missing or blank ``reason`` is a
silent regression: the row still surfaces and the totals still
tally, but the reviewer is left without an explanation, breaking
the audit's whole purpose.

Sibling test files cover the ``ignored-deleted`` vocabulary
specifically (``test_ignored_deleted_reason_exact_match.py``,
``test_ignored_deleted_reason.py``). This module is broader: it
pins the universal "reason field is always populated" contract
across **every status** the audit can emit, and across every
output surface (in-process dataclass, CLI text table, ``--json``
array, ``--similarity-csv`` export).

The single hand-built payload exercises all five status branches
at once so a future branch that forgets the ``reason=`` keyword
(or sets it to ``""``) is caught regardless of which status it
regressed.
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
_resolve_changed_md = _MOD._resolve_changed_md
_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")

# Canonical set of statuses the audit can emit today. Frozen so a
# future status addition has to extend this set explicitly rather
# than silently slip past the coverage check below.
_KNOWN_STATUSES: frozenset[str] = frozenset({
    "matched", "ignored-extension", "ignored-out-of-root",
    "ignored-missing", "ignored-deleted",
})


def _is_blank(value: object) -> bool:
    """Return True for ``None``, non-strings, and strings that are
    empty or whitespace-only after stripping. Centralised so every
    surface uses the same definition of "blank reason".
    """
    return not (isinstance(value, str) and value.strip() != "")


class _MultiStatusSandbox:
    """Build a temp repo whose `--changed-files` payload exercises
    all five known statuses in a single run. The CLI tests reuse
    the same payload so cross-surface assertions stay aligned.

    Layout::

        <tmp>/
        ├── changed.txt              ← payload below
        ├── spec/
        │   └── present.md           ← MATCHED row
        └── outside/                 ← lives OUTSIDE --root spec/

    Payload entries (one per status):

    * ``spec/present.md`` — file exists, ext OK, under root → matched
    * ``spec/notes.txt`` — wrong extension → ignored-extension
    * ``outside/elsewhere.md`` — outside --root → ignored-out-of-root
    * ``spec/never-existed.md`` — under root, ext OK, missing on disk
      → ignored-missing
    * ``D\\tspec/gone.md`` — explicit delete marker → ignored-deleted
    """

    PAYLOAD = (
        "spec/present.md\n"
        "spec/notes.txt\n"
        "outside/elsewhere.md\n"
        "spec/never-existed.md\n"
        "D\tspec/gone.md\n"
    )

    def __init__(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.root = Path(self._tmp.name).resolve()
        self.spec = self.root / "spec"
        self.spec.mkdir()
        # The MATCHED row needs an actual on-disk file.
        (self.spec / "present.md").write_text("# present\n",
                                              encoding="utf-8")
        self.changed = self.root / "changed.txt"
        self.changed.write_text(self.PAYLOAD, encoding="utf-8")

    def __enter__(self) -> "_MultiStatusSandbox":
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


def _audit_in_process(payload: str, *, root: Path,
                       repo_root: Path) -> list:
    """Drive ``_resolve_changed_md`` directly so the dataclass-level
    ``reason`` contract is exercised without a subprocess hop.
    """
    audit: list = []
    cf = repo_root / "changed.txt"
    cf.write_text(payload, encoding="utf-8")
    _resolve_changed_md(
        repo_root=repo_root, root=root,
        diff_base=None, changed_files=str(cf),
        extensions=("md",), audit=audit,
    )
    return audit


# --- 1. In-process dataclass contract ----------------------------

class DataclassRowsAlwaysCarryReason(unittest.TestCase):
    """Every ``_ChangedFileAudit`` produced by the live intake has
    a non-empty ``reason`` regardless of which status branch it
    came from. The single multi-status payload exercises all five
    branches so a regression on any one is caught here.
    """

    def _audit(self) -> list:
        # Build the same on-disk fixture the CLI sandbox uses so
        # the in-process and CLI tests share their payload.
        with _MultiStatusSandbox() as box:
            return _audit_in_process(
                _MultiStatusSandbox.PAYLOAD,
                root=box.spec, repo_root=box.root,
            )

    def test_audit_is_non_empty(self) -> None:
        # Sanity: if the payload silently produced zero rows, the
        # per-row assertions below would vacuously pass and hide a
        # much worse regression.
        rows = self._audit()
        self.assertGreater(len(rows), 0,
                           msg="multi-status payload produced no "
                               "audit rows — fixture likely broken")

    def test_every_row_has_non_empty_reason(self) -> None:
        rows = self._audit()
        for row in rows:
            self.assertIsInstance(
                row.reason, str,
                msg=f"row path={row.path!r} status={row.status!r} "
                    f"reason is not a str: {row.reason!r}")
            self.assertNotEqual(
                row.reason.strip(), "",
                msg=f"row path={row.path!r} status={row.status!r} "
                    f"has blank reason: {row.reason!r}")

    def test_payload_covers_every_known_status(self) -> None:
        # Coverage guard: if a future refactor breaks one of the
        # five branches (e.g. matched rows stop being emitted) the
        # per-row reason check above could pass on a degraded set.
        # Fail loud here so a missing branch is obvious.
        rows = self._audit()
        observed = {r.status for r in rows}
        self.assertEqual(observed, _KNOWN_STATUSES,
                         msg=f"payload missed statuses "
                             f"{_KNOWN_STATUSES - observed}; got "
                             f"{sorted(observed)}")


# --- 2. JSON surface ---------------------------------------------

class JsonAuditRowsAlwaysCarryReason(unittest.TestCase):
    """The ``--json`` audit array on STDERR is the schema-stable
    surface dashboards consume. Every row must include a ``reason``
    key whose value is a non-empty string — missing key, ``null``,
    and ``""`` all fail.
    """

    def _audit_array(self) -> list[dict]:
        with _MultiStatusSandbox() as box:
            proc = box.run("--json")
        self.assertEqual(proc.returncode, 0, msg=proc.stderr)
        return json.loads(proc.stderr.strip())

    def test_every_json_row_has_reason_key(self) -> None:
        for row in self._audit_array():
            self.assertIn("reason", row,
                          msg=f"row missing 'reason' key: {row!r}")

    def test_every_json_reason_is_non_empty_string(self) -> None:
        for row in self._audit_array():
            value = row.get("reason")
            self.assertFalse(
                _is_blank(value),
                msg=f"row path={row.get('path')!r} "
                    f"status={row.get('status')!r} has blank reason: "
                    f"{value!r}")

    def test_json_audit_covers_every_known_status(self) -> None:
        observed = {r["status"] for r in self._audit_array()}
        self.assertEqual(observed, _KNOWN_STATUSES,
                         msg=f"JSON audit missed statuses "
                             f"{_KNOWN_STATUSES - observed}")


# --- 3. Text-table surface ---------------------------------------

class TextTableRowsAlwaysCarryReason(unittest.TestCase):
    """The default human text table is what an operator sees first
    in CI logs. Each row must carry a visible reason cell — a row
    that ends right after its path (no third column) would mean a
    blank reason silently rendered.
    """

    # Row shape:
    #   <leading-ws> <status> <whitespace> <path> <whitespace>
    #   <reason...>
    # Anchor on the known status set so we don't accidentally
    # capture header / footer / separator lines.
    _STATUS_ALT = "|".join(re.escape(s) for s in _KNOWN_STATUSES)
    _ROW_RE = re.compile(
        rf"^\s*(?P<status>{_STATUS_ALT})\s+(?P<path>\S+)"
        rf"(?:\s+(?P<reason>.+?))?\s*$")

    def _rows(self, stderr: str) -> list[dict]:
        out: list[dict] = []
        for line in stderr.splitlines():
            m = self._ROW_RE.match(line)
            if m:
                out.append({
                    "status": m.group("status"),
                    "path": m.group("path"),
                    # `None` when the reason column is absent — the
                    # _is_blank check below treats that as failure.
                    "reason": m.group("reason"),
                })
        return out

    def test_every_text_row_has_visible_reason_cell(self) -> None:
        with _MultiStatusSandbox() as box:
            proc = box.run()
        self.assertEqual(proc.returncode, 0, msg=proc.stderr)
        rows = self._rows(proc.stderr)
        self.assertGreater(len(rows), 0,
                           msg=f"no rows parsed; stderr:\n{proc.stderr}")
        for row in rows:
            self.assertFalse(
                _is_blank(row["reason"]),
                msg=f"text row status={row['status']!r} "
                    f"path={row['path']!r} has blank reason cell "
                    f"{row['reason']!r}\nfull stderr:\n{proc.stderr}")

    def test_text_audit_covers_every_known_status(self) -> None:
        with _MultiStatusSandbox() as box:
            proc = box.run()
        rows = self._rows(proc.stderr)
        observed = {r["status"] for r in rows}
        self.assertEqual(observed, _KNOWN_STATUSES,
                         msg=f"text audit missed statuses "
                             f"{_KNOWN_STATUSES - observed}")


# --- 4. CSV surface ----------------------------------------------

class CsvAuditRowsAlwaysCarryReason(unittest.TestCase):
    """The ``--similarity-csv`` export is the spreadsheet surface.
    The header is fixed (``path,status,reason,kind,score,old_path``)
    so absence of the column would be a parser failure — but a row
    that leaves the cell empty (``,,``) would still be valid CSV
    and silently pass for any consumer doing column-positional
    reads. Catch that here.
    """

    def test_every_csv_row_has_non_empty_reason_cell(self) -> None:
        with _MultiStatusSandbox() as box:
            csv_path = box.root / "audit.csv"
            proc = box.run("--similarity-csv", str(csv_path))
            self.assertEqual(proc.returncode, 0, msg=proc.stderr)
            self.assertTrue(csv_path.is_file(),
                            msg=f"CSV not produced; stderr:\n"
                                f"{proc.stderr}")
            with csv_path.open(encoding="utf-8", newline="") as fh:
                reader = csv.DictReader(fh)
                self.assertIn("reason", reader.fieldnames or [],
                              msg=f"CSV header missing 'reason': "
                                  f"{reader.fieldnames!r}")
                rows = list(reader)

        self.assertGreater(len(rows), 0,
                           msg="CSV produced no data rows")
        for row in rows:
            self.assertFalse(
                _is_blank(row.get("reason")),
                msg=f"CSV row status={row.get('status')!r} "
                    f"path={row.get('path')!r} has blank reason: "
                    f"{row.get('reason')!r}")

    def test_csv_audit_covers_every_known_status(self) -> None:
        with _MultiStatusSandbox() as box:
            csv_path = box.root / "audit.csv"
            box.run("--similarity-csv", str(csv_path))
            with csv_path.open(encoding="utf-8", newline="") as fh:
                rows = list(csv.DictReader(fh))
        observed = {r["status"] for r in rows}
        self.assertEqual(observed, _KNOWN_STATUSES,
                         msg=f"CSV audit missed statuses "
                             f"{_KNOWN_STATUSES - observed}")


if __name__ == "__main__":
    unittest.main()