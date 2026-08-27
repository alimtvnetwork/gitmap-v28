"""Tests for ``check-tunable-constants.py`` Rule T4 — session-TTL invariant.

Rule T4 (defined in ``spec/19-main-worker-service/15-tunable-constants.md`` §6,
authored under FU-16) requires::

    MainWorker.Auth.MainSessionAbsoluteMaxSeconds
        >= MainWorker.Auth.MainSessionTtlSeconds

in BOTH the §2 catalogue defaults and the §4 ``config.seed.json`` defaults.

This file covers the two public T4 entry points:

* ``rule_t4_pair(label, ttl_raw, abs_raw)`` — pure function, no I/O.
  All branches: missing TTL, missing AbsoluteMax, non-numeric inputs,
  passing pair (cap >= ttl), failing pair (cap < ttl), exact-equal edge.

* ``rule_t4()`` — orchestrator that reads the live spec via
  ``collect_catalogue_defaults`` / ``collect_seed_defaults``. We
  monkey-patch those so the tests do not depend on the on-disk spec
  values (otherwise a future spec edit would silently break the suite).

Hardened ``parse_seconds`` (see same module) is exercised transitively
via the unit-bearing fixtures (``"28800s"``, ``"8h"``, ``"1 day"``).
The dedicated ``parse_seconds`` unit-test surface lives in this same
file under ``ParseSecondsTests`` so all T4-adjacent guarantees live
together.
"""
from __future__ import annotations

import sys
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from conftest_shim import load_tunable_constants_linter  # noqa: E402

LINT = load_tunable_constants_linter()

TTL_KEY = LINT.SESSION_TTL_KEY            # MainWorker.Auth.MainSessionTtlSeconds
ABS_KEY = LINT.SESSION_ABS_KEY            # MainWorker.Auth.MainSessionAbsoluteMaxSeconds
TTL_SEED = LINT.SESSION_TTL_SEED          # MainSessionTtlSeconds
ABS_SEED = LINT.SESSION_ABS_SEED          # MainSessionAbsoluteMaxSeconds


def _has_substr(items: list[str], needle: str) -> bool:
    return any(needle in item for item in items)


class RuleT4PairTests(unittest.TestCase):
    """Pure unit tests for ``rule_t4_pair`` — covers every branch."""

    def test_pair_ok_when_cap_greater_than_ttl(self) -> None:
        out = LINT.rule_t4_pair("§2", "28800", "86400")
        self.assertEqual(out, [], "Expected no violations when cap > ttl")

    def test_pair_ok_when_cap_equals_ttl(self) -> None:
        # The invariant is `cap >= ttl`, so equality is a pass.
        out = LINT.rule_t4_pair("§2", "3600", "3600")
        self.assertEqual(out, [])

    def test_pair_fails_when_cap_below_ttl(self) -> None:
        out = LINT.rule_t4_pair("§4", "28800", "3600")
        self.assertEqual(len(out), 1)
        self.assertIn("invariant violated", out[0])
        self.assertIn("AbsoluteMax(3600)", out[0])
        self.assertIn("Ttl(28800)", out[0])
        self.assertIn("§4", out[0], "label must be threaded through")

    def test_pair_missing_ttl_raw(self) -> None:
        out = LINT.rule_t4_pair("§2", None, "86400")
        self.assertEqual(len(out), 1)
        self.assertIn(TTL_KEY, out[0])
        self.assertIn("missing", out[0])

    def test_pair_missing_abs_raw(self) -> None:
        out = LINT.rule_t4_pair("§2", "28800", None)
        self.assertEqual(len(out), 1)
        self.assertIn(ABS_KEY, out[0])
        self.assertIn("missing", out[0])

    def test_pair_both_missing_reports_ttl_first(self) -> None:
        # Implementation contract: TTL absence is reported before AbsMax
        # (early-return order). Locking this in prevents accidental
        # reordering that would change CI error output.
        out = LINT.rule_t4_pair("§2", None, None)
        self.assertEqual(len(out), 1)
        self.assertIn(TTL_KEY, out[0])

    def test_pair_non_numeric_ttl(self) -> None:
        out = LINT.rule_t4_pair("§2", "not a number", "86400")
        self.assertEqual(len(out), 1)
        self.assertIn("non-numeric", out[0])
        self.assertIn("'not a number'", out[0])

    def test_pair_non_numeric_abs(self) -> None:
        out = LINT.rule_t4_pair("§4", "28800", "tbd")
        self.assertEqual(len(out), 1)
        self.assertIn("non-numeric", out[0])
        self.assertIn("'tbd'", out[0])

    def test_pair_empty_string_treated_as_non_numeric(self) -> None:
        # Empty string is truthy-falsy edge: parse_seconds returns None,
        # which the pair surface reports as non-numeric (NOT as "missing"
        # — that path is reserved for the explicit None sentinel).
        out = LINT.rule_t4_pair("§2", "", "86400")
        self.assertEqual(len(out), 1)
        self.assertIn("non-numeric", out[0])

    def test_pair_handles_unit_bearing_inputs(self) -> None:
        # Hardened parse_seconds: `8h` == 28800s, `1 day` == 86400s.
        # Pair must accept and compare them correctly.
        out = LINT.rule_t4_pair("§2", "8h", "1 day")
        self.assertEqual(out, [])

    def test_pair_unit_mix_failing(self) -> None:
        # 1 day (86400) > 1h (3600) — pair where TTL=1day, cap=1h fails.
        out = LINT.rule_t4_pair("§4", "1 day", "1h")
        self.assertEqual(len(out), 1)
        self.assertIn("invariant violated", out[0])
        self.assertIn("AbsoluteMax(3600)", out[0])
        self.assertIn("Ttl(86400)", out[0])

    def test_pair_handles_existing_spec_format(self) -> None:
        # Smoke-test the actual on-disk spec strings to guarantee
        # no-regression against the live `15-tunable-constants.md`.
        out = LINT.rule_t4_pair("§2", "**28800** (8h)", "**86400** (24h)")
        self.assertEqual(out, [])


class RuleT4OrchestratorTests(unittest.TestCase):
    """Tests for ``rule_t4`` — exercises the catalogue + seed pair lookup."""

    def setUp(self) -> None:
        # Snapshot the originals so each test restores cleanly even if
        # an assertion raises mid-test.
        self._orig_cat = LINT.collect_catalogue_defaults
        self._orig_seed = LINT.collect_seed_defaults

    def tearDown(self) -> None:
        LINT.collect_catalogue_defaults = self._orig_cat
        LINT.collect_seed_defaults = self._orig_seed

    def _patch(self, cat: dict[str, str], seed: dict[str, str]) -> None:
        LINT.collect_catalogue_defaults = lambda: cat
        LINT.collect_seed_defaults = lambda: seed

    def test_orchestrator_passes_when_both_layers_satisfy(self) -> None:
        self._patch(
            cat={TTL_KEY: "28800", ABS_KEY: "86400"},
            seed={TTL_SEED: "28800", ABS_SEED: "86400"},
        )
        self.assertEqual(LINT.rule_t4(), [])

    def test_orchestrator_passes_with_unit_bearing_layers(self) -> None:
        # Catalogue and seed allowed to use different surface forms
        # (unit-bearing vs bare seconds) as long as the numeric value
        # satisfies the invariant. This guards the parse_seconds /
        # rule_t4 contract: same cap regardless of how authors wrote it.
        self._patch(
            cat={TTL_KEY: "**28800** (8h)", ABS_KEY: "**86400** (24h)"},
            seed={TTL_SEED: "8h", ABS_SEED: "1 day"},
        )
        self.assertEqual(LINT.rule_t4(), [])

    def test_orchestrator_flags_catalogue_failure(self) -> None:
        # §2 broken, §4 fine — exactly one violation labelled `§2`.
        self._patch(
            cat={TTL_KEY: "28800", ABS_KEY: "3600"},
            seed={TTL_SEED: "28800", ABS_SEED: "86400"},
        )
        out = LINT.rule_t4()
        self.assertEqual(len(out), 1)
        self.assertTrue(out[0].startswith("§2"))
        self.assertIn("invariant violated", out[0])

    def test_orchestrator_flags_seed_failure(self) -> None:
        # §4 broken, §2 fine — exactly one violation labelled `§4`.
        self._patch(
            cat={TTL_KEY: "28800", ABS_KEY: "86400"},
            seed={TTL_SEED: "28800", ABS_SEED: "3600"},
        )
        out = LINT.rule_t4()
        self.assertEqual(len(out), 1)
        self.assertTrue(out[0].startswith("§4"))

    def test_orchestrator_flags_both_layers_independently(self) -> None:
        self._patch(
            cat={TTL_KEY: "28800", ABS_KEY: "3600"},
            seed={TTL_SEED: "28800", ABS_SEED: "3600"},
        )
        out = LINT.rule_t4()
        self.assertEqual(len(out), 2)
        self.assertTrue(_has_substr(out, "§2"))
        self.assertTrue(_has_substr(out, "§4"))

    def test_orchestrator_flags_missing_catalogue_keys(self) -> None:
        # Empty catalogue → both keys reported missing for §2; §4 ok.
        self._patch(
            cat={},
            seed={TTL_SEED: "28800", ABS_SEED: "86400"},
        )
        out = LINT.rule_t4()
        # Order is: missing TTL is the early-return short-circuit, so
        # §2 reports exactly one item (the TTL absence). The orchestrator
        # then continues to §4 which is clean.
        self.assertEqual(len(out), 1)
        self.assertIn(TTL_KEY, out[0])
        self.assertTrue(out[0].startswith("§2"))

    def test_orchestrator_flags_missing_seed_keys(self) -> None:
        self._patch(
            cat={TTL_KEY: "28800", ABS_KEY: "86400"},
            seed={TTL_SEED: "28800"},  # ABS_SEED missing
        )
        out = LINT.rule_t4()
        self.assertEqual(len(out), 1)
        self.assertIn(ABS_KEY, out[0])
        self.assertTrue(out[0].startswith("§4"))

    def test_orchestrator_flags_non_numeric_in_seed(self) -> None:
        self._patch(
            cat={TTL_KEY: "28800", ABS_KEY: "86400"},
            seed={TTL_SEED: "tbd", ABS_SEED: "86400"},
        )
        out = LINT.rule_t4()
        self.assertEqual(len(out), 1)
        self.assertIn("non-numeric", out[0])
        self.assertTrue(out[0].startswith("§4"))


class ParseSecondsTests(unittest.TestCase):
    """Companion tests: lock in the parser contract that T4 depends on."""

    def test_bare_int_returns_seconds(self) -> None:
        self.assertEqual(LINT.parse_seconds("28800"), 28800)

    def test_existing_spec_bolded_form(self) -> None:
        self.assertEqual(LINT.parse_seconds("**28800** (8h)"), 28800)
        self.assertEqual(LINT.parse_seconds("**86400** (24h)"), 86400)

    def test_unit_suffix_seconds(self) -> None:
        self.assertEqual(LINT.parse_seconds("28800s"), 28800)
        self.assertEqual(LINT.parse_seconds("90s"), 90)

    def test_unit_hours(self) -> None:
        self.assertEqual(LINT.parse_seconds("8h"), 28800)
        self.assertEqual(LINT.parse_seconds("8 h"), 28800)
        self.assertEqual(LINT.parse_seconds("12 hours"), 43200)

    def test_unit_minutes_and_days(self) -> None:
        self.assertEqual(LINT.parse_seconds("60 min"), 3600)
        self.assertEqual(LINT.parse_seconds("1 day"), 86400)

    def test_milliseconds_floor(self) -> None:
        # 500ms → 0 whole seconds (floored), 1500ms → 1.
        self.assertEqual(LINT.parse_seconds("500ms"), 0)
        self.assertEqual(LINT.parse_seconds("1500ms"), 1)

    def test_first_token_wins_on_slashed_form(self) -> None:
        # `28800s/8h` — first numeric token is the canonical value.
        self.assertEqual(LINT.parse_seconds("28800s/8h"), 28800)

    def test_code_fenced_value(self) -> None:
        self.assertEqual(LINT.parse_seconds("`28800`"), 28800)

    def test_returns_none_on_empty(self) -> None:
        self.assertIsNone(LINT.parse_seconds(""))
        self.assertIsNone(LINT.parse_seconds("   "))

    def test_returns_none_on_no_digits(self) -> None:
        self.assertIsNone(LINT.parse_seconds("not a number"))


if __name__ == "__main__":  # pragma: no cover
    unittest.main(verbosity=2)
