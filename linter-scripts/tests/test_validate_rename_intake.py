"""Tests for ``linter-scripts/validate-rename-intake.py``.

The validator is the CI gate for `rename_intake` JSON, so the
test surface is large on purpose: every documented mode, every
exit code, every cross-field rule. Three concerns:

* **Positive paths** — the shipped example artifact (and a real
  linter run captured to a temp file) validate cleanly in legacy
  / enriched / labels modes;
* **Negative paths** — every documented constraint produces a
  schema-violation exit (code 1) with a useful pointer + message;
* **Schema export** — ``--print-schema`` emits a Draft 2020-12
  schema whose constraints match what the in-process validator
  enforces (so an external `check-jsonschema` run agrees with
  this script row-for-row).
"""

from __future__ import annotations

import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

VALIDATOR = (Path(__file__).resolve().parent.parent
             / "validate-rename-intake.py")
LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")
EXAMPLE = (Path(__file__).resolve().parent.parent
           / "examples" / "rename-intake-audit.json")


def _run(*args: str, stdin: str | None = None
         ) -> tuple[int, str, str]:
    r = subprocess.run(
        [sys.executable, str(VALIDATOR), *args],
        input=stdin, capture_output=True, text=True)
    return r.returncode, r.stdout, r.stderr


def _write(td: Path, name: str, payload) -> Path:
    p = td / name
    p.write_text(json.dumps(payload), encoding="utf-8")
    return p


# ---------------------------------------------------------------
# Positive paths
# ---------------------------------------------------------------
class PositivePaths(unittest.TestCase):
    def test_shipped_example_validates_in_enriched_mode(self) -> None:
        code, out, err = _run(str(EXAMPLE), "--with-similarity")
        self.assertEqual(code, 0,
            msg=f"shipped example must validate; err={err!r}")
        self.assertIn("OK", out)
        self.assertIn("14 record(s)", out)

    def test_legacy_minimal_record_validates(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            p = _write(Path(td), "a.json", [
                {"path": "spec/x.md", "status": "matched",
                 "reason": "ok"},
            ])
            code, _out, err = _run(str(p))
            self.assertEqual(code, 0, msg=f"err={err!r}")

    def test_enriched_with_object_and_null_similarity(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            p = _write(Path(td), "a.json", [
                {"path": "a.md", "status": "matched", "reason": "ok",
                 "similarity": {"kind": "R", "score": 92,
                                "old_path": "b.md"}},
                {"path": "c.md", "status": "matched", "reason": "ok",
                 "similarity": None},
                {"path": "d.md", "status": "matched", "reason": "ok",
                 "similarity": {"kind": "C", "score": None,
                                "old_path": "e.md"}},
            ])
            code, _out, err = _run(str(p), "--with-similarity")
            self.assertEqual(code, 0, msg=f"err={err!r}")

    def test_labels_mode_with_score_kind(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            p = _write(Path(td), "a.json", [
                {"path": "a.md", "status": "matched", "reason": "ok",
                 "similarity": {"kind": "R", "score": 92,
                                "old_path": "b.md",
                                "score_kind": "rename-similarity"}},
                {"path": "c.md", "status": "matched", "reason": "ok",
                 "similarity": None},
            ])
            code, _out, err = _run(str(p), "--with-labels")
            self.assertEqual(code, 0, msg=f"err={err!r}")

    def test_quiet_suppresses_success_line(self) -> None:
        code, out, _ = _run(str(EXAMPLE),
                            "--with-similarity", "--quiet")
        self.assertEqual(code, 0)
        self.assertEqual(out.strip(), "")

    def test_stdin_input_works(self) -> None:
        body = json.dumps([{"path": "a.md", "status": "matched",
                            "reason": "ok"}])
        code, _out, err = _run("-", stdin=body)
        self.assertEqual(code, 0, msg=f"err={err!r}")


# ---------------------------------------------------------------
# Negative paths — every documented constraint
# ---------------------------------------------------------------
class NegativePaths(unittest.TestCase):
    def _expect_fail(self, args, payload, *, fragment: str) -> None:
        with tempfile.TemporaryDirectory() as td:
            p = _write(Path(td), "a.json", payload)
            code, _out, err = _run(str(p), *args)
            self.assertEqual(code, 1,
                msg=f"expected exit 1; got {code}; err={err!r}")
            self.assertIn(fragment, err,
                msg=f"missing {fragment!r} in stderr: {err!r}")

    def test_invalid_json_returns_exit_2(self) -> None:
        code, _out, err = _run("-", stdin="{not json")
        self.assertEqual(code, 2)
        self.assertIn("invalid JSON", err)

    def test_top_level_not_an_array(self) -> None:
        code, _out, err = _run("-", stdin='{"not": "an array"}')
        self.assertEqual(code, 1)
        self.assertIn("must be a JSON array", err)

    def test_empty_array_fails_by_default(self) -> None:
        code, _out, err = _run("-", stdin="[]")
        self.assertEqual(code, 1)
        self.assertIn("array is empty", err)

    def test_empty_array_ok_with_allow_empty(self) -> None:
        code, out, _err = _run("-", "--allow-empty", stdin="[]")
        self.assertEqual(code, 0)
        self.assertIn("0 record(s)", out)

    def test_unknown_status_rejected(self) -> None:
        self._expect_fail([], [{
            "path": "a.md", "status": "bogus", "reason": "ok"}],
            fragment="status must be one of")

    def test_missing_required_key(self) -> None:
        self._expect_fail([], [{"path": "a.md", "status": "matched"}],
            fragment="missing required key(s)")

    def test_legacy_rejects_extra_similarity_key(self) -> None:
        self._expect_fail([], [{
            "path": "a.md", "status": "matched", "reason": "ok",
            "similarity": None}],
            fragment="unexpected key(s) for legacy schema")

    def test_enriched_requires_similarity_key(self) -> None:
        self._expect_fail(["--with-similarity"], [{
            "path": "a.md", "status": "matched", "reason": "ok"}],
            fragment="missing required key(s)")

    def test_similarity_object_bad_kind(self) -> None:
        self._expect_fail(["--with-similarity"], [{
            "path": "a.md", "status": "matched", "reason": "ok",
            "similarity": {"kind": "X", "score": 50,
                           "old_path": "b.md"}}],
            fragment="kind must be one of")

    def test_similarity_score_out_of_range(self) -> None:
        self._expect_fail(["--with-similarity"], [{
            "path": "a.md", "status": "matched", "reason": "ok",
            "similarity": {"kind": "R", "score": 150,
                           "old_path": "b.md"}}],
            fragment="score must be in [0, 100]")

    def test_similarity_extra_key(self) -> None:
        self._expect_fail(["--with-similarity"], [{
            "path": "a.md", "status": "matched", "reason": "ok",
            "similarity": {"kind": "R", "score": 50,
                           "old_path": "b.md", "junk": 1}}],
            fragment="unexpected key(s)")

    def test_labels_required_when_flag_set(self) -> None:
        self._expect_fail(["--with-labels"], [{
            "path": "a.md", "status": "matched", "reason": "ok",
            "similarity": {"kind": "R", "score": 50,
                           "old_path": "b.md"}}],
            fragment="missing key(s): ['score_kind']")

    def test_score_kind_invalid_value(self) -> None:
        self._expect_fail(["--with-labels"], [{
            "path": "a.md", "status": "matched", "reason": "ok",
            "similarity": {"kind": "R", "score": 50,
                           "old_path": "b.md",
                           "score_kind": "bogus"}}],
            fragment="score_kind must be one of")

    def test_ignored_deleted_with_non_null_similarity(self) -> None:
        self._expect_fail(["--with-similarity"], [{
            "path": "a.md", "status": "ignored-deleted",
            "reason": "D",
            "similarity": {"kind": "R", "score": 50,
                           "old_path": "b.md"}}],
            fragment="ignored-deleted rows must have similarity:null")


# ---------------------------------------------------------------
# Schema export
# ---------------------------------------------------------------
class SchemaExport(unittest.TestCase):
    def test_print_schema_legacy_is_valid_draft_2020_12(self) -> None:
        code, out, _err = _run("--print-schema")
        self.assertEqual(code, 0)
        schema = json.loads(out)
        self.assertEqual(schema["$schema"],
            "https://json-schema.org/draft/2020-12/schema")
        self.assertEqual(schema["type"], "array")
        item = schema["items"]
        self.assertEqual(set(item["required"]),
            {"path", "status", "reason"})
        self.assertNotIn("similarity", item["properties"],
            msg="legacy schema must not declare a similarity prop")
        self.assertFalse(item["additionalProperties"])

    def test_print_schema_enriched_includes_similarity(self) -> None:
        code, out, _err = _run("--print-schema", "--with-similarity")
        self.assertEqual(code, 0)
        schema = json.loads(out)
        item = schema["items"]
        self.assertIn("similarity", item["properties"])
        self.assertIn("similarity", item["required"])

    def test_print_schema_labels_requires_score_kind(self) -> None:
        code, out, _err = _run("--print-schema", "--with-labels")
        self.assertEqual(code, 0)
        schema = json.loads(out)
        sim_object = (schema["items"]["properties"]["similarity"]
                      ["oneOf"][0])
        self.assertIn("score_kind", sim_object["required"])
        self.assertIn("score_kind", sim_object["properties"])
        self.assertEqual(
            set(sim_object["properties"]["score_kind"]["enum"]),
            {"rename-similarity", "copy-similarity", "unscored"})

    def test_with_labels_implies_with_similarity(self) -> None:
        # --with-labels alone (no --with-similarity) must still
        # produce the enriched schema; the validator auto-promotes.
        code, out, _err = _run("--print-schema", "--with-labels")
        self.assertEqual(code, 0)
        schema = json.loads(out)
        self.assertIn("similarity", schema["items"]["required"])


# ---------------------------------------------------------------
# End-to-end: real linter run → validator passes
# ---------------------------------------------------------------
class EndToEndAgainstRealLinter(unittest.TestCase):
    def _setup_repo(self, td: Path) -> tuple[Path, Path]:
        spec = td / "spec"; spec.mkdir()
        (spec / "intro.md").write_text("# spec\nplain prose.\n")
        (spec / "copy.md").write_text("# spec\nplain prose.\n")
        payload = td / "changed.txt"
        payload.write_text(
            "R087\tspec/intro.md\tspec/copy.md\n"
            "M\tspec/intro.md\n")
        return spec, payload

    def _run_linter(self, td: Path, *extra: str) -> Path:
        spec, payload = self._setup_repo(td)
        out_path = td / "audit.json"
        proc = subprocess.run(
            [sys.executable, str(LINTER),
             "--root", str(spec),
             "--changed-files", str(payload),
             "--list-changed-files", "--json", *extra],
            capture_output=True, text=True, cwd=td,
        )
        self.assertEqual(proc.returncode, 0,
            msg=f"linter failed; err={proc.stderr!r}")
        out_path.write_text(proc.stderr, encoding="utf-8")
        return out_path

    def test_validator_accepts_real_legacy_output(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            audit = self._run_linter(Path(td))
            code, _out, err = _run(str(audit))
            self.assertEqual(code, 0,
                msg=f"validator rejected real legacy output; "
                    f"err={err!r}")

    def test_validator_accepts_real_enriched_output(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            audit = self._run_linter(Path(td), "--with-similarity")
            code, _out, err = _run(str(audit), "--with-similarity")
            self.assertEqual(code, 0,
                msg=f"validator rejected real enriched output; "
                    f"err={err!r}")

    def test_validator_accepts_real_labels_output(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            audit = self._run_linter(Path(td),
                "--with-similarity", "--similarity-labels")
            code, _out, err = _run(str(audit), "--with-labels")
            self.assertEqual(code, 0,
                msg=f"validator rejected real labels output; "
                    f"err={err!r}")


if __name__ == "__main__":
    unittest.main()