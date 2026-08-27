"""Pins the ready-to-copy JSON example payload at
``linter-scripts/examples/rename-intake-audit.json``.

The example is referenced from ``README-rename-intake.md`` as the
canonical demonstration of every ``status`` × ``similarity``-shape
combination the rename/copy intake can emit. If a future schema
change drops a status, renames a key, or alters the scored /
unscored / null trichotomy, this test fails so the docs and the
shipped artifact stay in lock-step with the renderer.

Coverage we lock in:

* every value of the closed ``status`` vocabulary appears at least
  once;
* ``similarity`` only ever takes one of three shapes —
  ``{kind, score: int, old_path}`` (scored),
  ``{kind, score: null, old_path}`` (unscored), or ``null`` (no
  rename/copy provenance);
* every non-``ignored-deleted`` status is exercised in *both* the
  scored-object and the ``null`` shape so consumers can copy a
  representative row for any case;
* ``ignored-deleted`` rows always carry ``similarity: null`` (the
  pre-state has nothing to score against).
"""

from __future__ import annotations

import json
import unittest
from pathlib import Path

EXAMPLE = (Path(__file__).resolve().parent.parent
           / "examples" / "rename-intake-audit.json")

_STATUSES = {
    "matched", "ignored-extension", "ignored-out-of-root",
    "ignored-missing", "ignored-deleted",
}


def _shape(sim) -> str:
    if sim is None:
        return "null"
    if sim.get("score") is None:
        return "unscored"
    return "scored"


class ExamplePayloadCoverage(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.data = json.loads(EXAMPLE.read_text(encoding="utf-8"))

    def test_payload_is_a_nonempty_list(self) -> None:
        self.assertIsInstance(self.data, list)
        self.assertGreater(len(self.data), 0)

    def test_every_record_has_the_documented_keys(self) -> None:
        # Same key set as the JSON audit emits today; if a future
        # version adds a new top-level key the example must be
        # updated alongside the renderer.
        for r in self.data:
            self.assertEqual(set(r),
                {"path", "status", "reason", "similarity"},
                msg=f"unexpected key set in record {r!r}")

    def test_status_values_are_in_closed_vocabulary(self) -> None:
        for r in self.data:
            self.assertIn(r["status"], _STATUSES,
                msg=f"status {r['status']!r} not in closed "
                    f"vocabulary {sorted(_STATUSES)}")

    def test_every_status_appears_at_least_once(self) -> None:
        seen = {r["status"] for r in self.data}
        self.assertEqual(seen, _STATUSES,
            msg=f"missing statuses: {_STATUSES - seen}")

    def test_similarity_object_shape_is_canonical(self) -> None:
        for r in self.data:
            sim = r["similarity"]
            if sim is None:
                continue
            self.assertEqual(set(sim), {"kind", "score", "old_path"},
                msg=f"similarity object has wrong keys: {sim!r}")
            self.assertIn(sim["kind"], ("R", "C"),
                msg=f"kind must be R or C, got {sim['kind']!r}")
            score = sim["score"]
            if score is not None:
                self.assertIsInstance(score, int)
                self.assertGreaterEqual(score, 0)
                self.assertLessEqual(score, 100)
            self.assertIsInstance(sim["old_path"], str)
            self.assertTrue(sim["old_path"],
                msg="old_path must be non-empty when similarity "
                    "is an object")

    def test_each_status_covers_scored_and_null_when_meaningful(
            self) -> None:
        # Every status that CAN carry a similarity object is
        # exercised in BOTH a scored-object and a null shape so
        # operators copying the example don't have to guess.
        # ``ignored-deleted`` is excluded because the renderer
        # never attaches a similarity to it.
        per_status: dict[str, set[str]] = {}
        for r in self.data:
            per_status.setdefault(r["status"], set()).add(
                _shape(r["similarity"]))
        for status in _STATUSES - {"ignored-deleted"}:
            shapes = per_status.get(status, set())
            self.assertIn("scored", shapes,
                msg=f"{status!r} needs at least one scored row "
                    f"in the example, found shapes: {shapes}")
            self.assertIn("null", shapes,
                msg=f"{status!r} needs at least one similarity:"
                    f"null row in the example, found: {shapes}")

    def test_unscored_shape_is_demonstrated(self) -> None:
        # The unscored case (object with ``score: null``) is the
        # subtle one — easy to confuse with ``similarity: null``
        # or ``score: 0``. Pin that the example demonstrates it
        # at least once so the docs stay self-evident.
        unscored = [r for r in self.data
                    if _shape(r["similarity"]) == "unscored"]
        self.assertGreaterEqual(len(unscored), 1,
            msg="example must demonstrate at least one unscored "
                "row (`similarity.score: null`)")

    def test_ignored_deleted_rows_have_null_similarity(self) -> None:
        deleted = [r for r in self.data
                   if r["status"] == "ignored-deleted"]
        self.assertGreaterEqual(len(deleted), 1)
        for r in deleted:
            self.assertIsNone(r["similarity"],
                msg=f"ignored-deleted row must have "
                    f"similarity=null, got {r!r}")


if __name__ == "__main__":
    unittest.main()