"""JSON schema guard: ``reason`` is present on every audit row in
both the legacy and ``--list-changed-files-verbose`` shapes.

The ``--list-changed-files --json`` audit array is the
machine-stable surface dashboards consume. The shape comes in two
flavours today:

* **Legacy** (no extra flag) — three required keys per row:
  ``path``, ``status``, ``reason``. The ``source`` key is
  *stripped entirely* (key absent, not ``null``) so historical
  validators that close on the 3-key shape keep working.
* **Verbose** (``--list-changed-files-verbose``) — adds a
  top-level ``"source"`` key to **every** row (``str`` on
  ``ignored-deleted`` rows, ``null`` everywhere else) so
  consumers can ``.get("source")`` without branching on status.

Across both shapes, the ``reason`` field is the universal
contract: every row, regardless of status, must carry it as a
non-empty string. Sibling tests cover the field's *value*
vocabulary and cross-surface parity; this file pins the
**JSON schema** specifically — both shapes, both legacy + verbose
modes, every status the audit can emit.
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
_LINTER = (Path(__file__).resolve().parent.parent
           / "check-placeholder-comments.py")

# Closed status set the audit emits today. Used to assert the
# fixture exercises every branch so the per-row schema check
# can't pass on a degraded subset.
_KNOWN_STATUSES: frozenset[str] = frozenset({
    "matched", "ignored-extension", "ignored-out-of-root",
    "ignored-missing", "ignored-deleted",
})

# Legacy mode required keys. Frozen here as the *contract* — if
# the linter ever stops emitting one, this constant must change
# in the same patch as the schema change.
_LEGACY_REQUIRED_KEYS: frozenset[str] = frozenset({
    "path", "status", "reason",
})

# Verbose mode adds exactly one top-level key on top of the
# legacy three. The audit also carries an optional ``similarity``
# sub-object, but only when ``--with-similarity`` is also passed —
# the schema tested here is verbose-without-similarity, mirroring
# what dashboards that opt in *only* to provenance receive.
_VERBOSE_REQUIRED_KEYS: frozenset[str] = (
    _LEGACY_REQUIRED_KEYS | {"source"})


class _MultiStatusSandbox:
    """Same multi-status fixture used by
    ``test_audit_reason_field_present.py``: one payload that
    hits every status branch in a single CLI run, so schema
    assertions cover all five row classes simultaneously.
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
        (self.spec / "present.md").write_text("# present\n",
                                              encoding="utf-8")
        self.changed = self.root / "changed.txt"
        self.changed.write_text(self.PAYLOAD, encoding="utf-8")

    def __enter__(self) -> "_MultiStatusSandbox":
        return self

    def __exit__(self, *exc: object) -> None:
        self._tmp.cleanup()

    def run_json(self, *extra: str) -> list[dict]:
        cmd = [
            sys.executable, str(_LINTER),
            "--root", str(self.spec),
            "--changed-files", str(self.changed),
            "--list-changed-files", "--json",
            *extra,
        ]
        proc = subprocess.run(cmd, capture_output=True, text=True,
                              check=False, cwd=str(self.root))
        if proc.returncode != 0:
            raise AssertionError(
                f"linter exited {proc.returncode}\n"
                f"stdout:\n{proc.stdout}\n"
                f"stderr:\n{proc.stderr}"
            )
        # The audit array lives on STDERR in JSON mode — STDOUT
        # carries the (empty) violations array. See
        # ``test_ignored_deleted_audit_coverage.py``.
        return json.loads(proc.stderr.strip())


# --- 1. Legacy shape ---------------------------------------------

class LegacyJsonShapeIncludesReason(unittest.TestCase):
    """Without ``--list-changed-files-verbose``, every row in the
    JSON audit array must have exactly the three legacy keys
    (``path``, ``status``, ``reason``) — and ``source`` must be
    *absent* (key omitted, not ``null``) to preserve the historical
    schema byte-for-byte.
    """

    @classmethod
    def setUpClass(cls) -> None:
        with _MultiStatusSandbox() as box:
            cls.audit = box.run_json()

    def test_audit_is_a_list(self) -> None:
        # Schema root is an array, not an object.
        self.assertIsInstance(self.audit, list)

    def test_audit_covers_every_known_status(self) -> None:
        # Coverage guard: per-row schema checks below would
        # vacuously pass if the fixture skipped a branch.
        observed = {r["status"] for r in self.audit}
        self.assertEqual(observed, _KNOWN_STATUSES,
                         msg=f"missed statuses "
                             f"{_KNOWN_STATUSES - observed}")

    def test_every_row_has_reason_key(self) -> None:
        for row in self.audit:
            self.assertIn("reason", row,
                          msg=f"row missing 'reason' key: {row!r}")

    def test_every_row_reason_is_non_empty_string(self) -> None:
        # ``reason`` must be a populated str — null / empty / blank
        # all violate the contract.
        for row in self.audit:
            value = row.get("reason")
            self.assertIsInstance(
                value, str,
                msg=f"row {row!r} 'reason' is not str: {value!r}")
            self.assertNotEqual(
                value.strip(), "",
                msg=f"row {row!r} 'reason' is blank: {value!r}")

    def test_every_row_has_exactly_legacy_keys(self) -> None:
        # Exact key set, not just superset — a stray key
        # (e.g. accidental ``source: null`` leak) would break
        # validators that close on the 3-key shape.
        for row in self.audit:
            self.assertEqual(
                set(row), _LEGACY_REQUIRED_KEYS,
                msg=f"legacy row has unexpected keys: {row!r}")

    def test_legacy_mode_strips_source_key_entirely(self) -> None:
        # Absence ≠ null. The legacy schema preservation strips
        # the key; verbose mode adds it. Both modes must be
        # distinguishable by ``"source" in row``.
        for row in self.audit:
            self.assertNotIn("source", row,
                             msg=f"legacy row leaked 'source': "
                                 f"{row!r}")


# --- 2. Verbose shape --------------------------------------------

class VerboseJsonShapeIncludesReasonAndSource(unittest.TestCase):
    """With ``--list-changed-files-verbose``, every row carries
    the legacy three keys *plus* a top-level ``source``.
    ``reason`` remains required and non-empty on every status;
    ``source`` is a string on ``ignored-deleted`` rows and
    ``null`` everywhere else.
    """

    @classmethod
    def setUpClass(cls) -> None:
        with _MultiStatusSandbox() as box:
            cls.audit = box.run_json("--list-changed-files-verbose")

    def test_audit_covers_every_known_status(self) -> None:
        observed = {r["status"] for r in self.audit}
        self.assertEqual(observed, _KNOWN_STATUSES,
                         msg=f"verbose audit missed statuses "
                             f"{_KNOWN_STATUSES - observed}")

    def test_every_verbose_row_has_reason_key(self) -> None:
        for row in self.audit:
            self.assertIn("reason", row,
                          msg=f"verbose row missing 'reason': "
                              f"{row!r}")

    def test_every_verbose_row_reason_is_non_empty_string(
            self) -> None:
        for row in self.audit:
            value = row.get("reason")
            self.assertIsInstance(
                value, str,
                msg=f"verbose row {row!r} 'reason' not str: "
                    f"{value!r}")
            self.assertNotEqual(
                value.strip(), "",
                msg=f"verbose row {row!r} 'reason' blank: "
                    f"{value!r}")

    def test_every_verbose_row_has_exactly_verbose_keys(self) -> None:
        # Exact key set: legacy three + ``source``. Anything more
        # (e.g. ``similarity`` sub-object leaking in without the
        # ``--with-similarity`` flag) is a schema regression.
        for row in self.audit:
            self.assertEqual(
                set(row), _VERBOSE_REQUIRED_KEYS,
                msg=f"verbose row has unexpected keys: {row!r}")

    def test_verbose_source_is_string_on_deleted_rows(self) -> None:
        # Provenance tag must be the actual string, not omitted
        # or null, on every ignored-deleted row.
        deleted = [r for r in self.audit
                   if r["status"] == "ignored-deleted"]
        self.assertGreater(len(deleted), 0,
                           msg="fixture produced no deleted rows")
        for row in deleted:
            self.assertIsInstance(
                row.get("source"), str,
                msg=f"deleted row 'source' not a str: {row!r}")
            self.assertNotEqual(
                row["source"].strip(), "",
                msg=f"deleted row 'source' is blank: {row!r}")

    def test_verbose_source_is_null_on_non_deleted_rows(self) -> None:
        # Schema regularity: every non-deleted row carries
        # ``"source": null`` so consumers don't branch on status
        # before reading the field.
        for row in self.audit:
            if row["status"] != "ignored-deleted":
                self.assertIsNone(
                    row.get("source", "<missing>"),
                    msg=f"non-deleted row 'source' not null: "
                        f"{row!r}")


# --- 3. Legacy ↔ verbose alignment -------------------------------

class LegacyAndVerboseAgreeOnPathStatusReason(unittest.TestCase):
    """The verbose shape must be a strict *superset* of the legacy
    shape. For every row in the audit, the ``(path, status,
    reason)`` triple in legacy mode must equal the corresponding
    triple in verbose mode — adding ``--list-changed-files-verbose``
    only attaches a ``source`` key, it never changes the existing
    fields.
    """

    @classmethod
    def setUpClass(cls) -> None:
        with _MultiStatusSandbox() as box:
            cls.legacy = box.run_json()
            cls.verbose = box.run_json("--list-changed-files-verbose")

    @staticmethod
    def _triples(audit: list[dict]) -> list[tuple[str, str, str]]:
        return [(r["path"], r["status"], r["reason"]) for r in audit]

    def test_legacy_and_verbose_have_same_row_count(self) -> None:
        self.assertEqual(len(self.legacy), len(self.verbose),
                         msg="row count drift between legacy and "
                             "verbose")

    def test_path_status_reason_triples_match_in_order(self) -> None:
        # Order matters: the audit contract is "input order
        # preserved" so verbose mustn't reshuffle.
        self.assertEqual(self._triples(self.legacy),
                         self._triples(self.verbose),
                         msg="(path, status, reason) drift between "
                             "legacy and verbose JSON shapes")


if __name__ == "__main__":
    unittest.main()