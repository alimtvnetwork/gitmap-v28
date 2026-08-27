"""Tests for the ``--similarity-legend={auto,on,off}`` flag.

Three concerns under test:

* **Resolver unit** — :func:`_should_emit_similarity_legend` decides
  on/off from the mode + stream TTY status. We drive every branch
  directly (without a subprocess) so the boolean truth table is
  pinned independently of the renderer composition above it.

* **Renderer unit** — :func:`_render_changed_files_audit` only emits
  the legend when ``with_similarity=True`` AND the resolver says
  on. JSON mode and the ``--with-similarity``-off case must NOT
  emit prose, regardless of mode (machine consumers + log scrapers
  parsing the legacy schema must keep working byte-for-byte).

* **CLI** — argparse wires the choice into ``args.similarity_legend``,
  the default is ``auto``, and the renderer receives it. We invoke
  the linter as a subprocess (where STDERR is a pipe → not a TTY)
  so ``auto`` resolves to off and the suite stays portable across
  CI runners that don't have a controlling terminal. Explicit
  ``on`` is exercised over the same pipe to prove the override
  works without needing a real PTY.

The legend itself is intentionally a small, stable cheat-sheet
(four short lines, five with labels). Tests assert the
greppable substrings that downstream docs key on (``kind``,
``score``, ``old``, optional ``meaning``) rather than the full
prose so future copy tweaks don't churn the test file.
"""

from __future__ import annotations

import io
import os
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
    """spec/intro.md + spec/copy.md, both clean — same shape as the
    --with-similarity test fixture so the rename payload below has
    a real post-state file to point at."""
    spec = td / "spec"
    spec.mkdir()
    (spec / "intro.md").write_text("# spec\nplain prose.\n")
    (spec / "copy.md").write_text("# spec\nplain prose.\n")
    return spec


# ---------------------------------------------------------------
# Resolver unit — every branch of _should_emit_similarity_legend
# ---------------------------------------------------------------
class ShouldEmitSimilarityLegendUnit(unittest.TestCase):
    def setUp(self) -> None:
        self.mod = load_placeholder_linter()

    def test_mode_on_always_true_even_without_fileno(self) -> None:
        # ``on`` short-circuits the TTY probe entirely so a StringIO
        # (which has no fileno) still gets the legend.
        self.assertTrue(self.mod._should_emit_similarity_legend(
            "on", io.StringIO()))

    def test_mode_off_always_false_even_on_real_tty(self) -> None:
        # The real STDERR fd may or may not be a TTY in CI, but
        # ``off`` must override either way. We assert against the
        # actual sys.stderr because that's the production stream.
        self.assertFalse(self.mod._should_emit_similarity_legend(
            "off", sys.stderr))

    def test_auto_off_when_stream_has_no_fileno(self) -> None:
        # StringIO has no underlying fd → resolver must default to
        # off, mirroring "non-interactive" semantics.
        self.assertFalse(self.mod._should_emit_similarity_legend(
            "auto", io.StringIO()))

    def test_auto_off_when_stream_is_a_pipe(self) -> None:
        # An os.pipe() read end is a real fd but not a TTY. This is
        # the exact shape of subprocess.PIPE in CI, so the test
        # pins the most common production code path.
        r, w = os.pipe()
        try:
            with os.fdopen(r, "r") as rh:
                self.assertFalse(self.mod._should_emit_similarity_legend(
                    "auto", rh))
        finally:
            os.close(w)

    def test_auto_on_when_stream_is_a_pty(self) -> None:
        # Allocate a real pseudo-terminal so the auto-on branch is
        # exercised end-to-end. Skipped on platforms without
        # ``os.openpty`` (Windows) — the test is a portability
        # canary, not a hard requirement.
        if not hasattr(os, "openpty"):
            self.skipTest("os.openpty not available on this platform")
        master, slave = os.openpty()
        try:
            with os.fdopen(slave, "w") as sh:
                self.assertTrue(self.mod._should_emit_similarity_legend(
                    "auto", sh))
        finally:
            os.close(master)

    def test_auto_off_when_fileno_raises(self) -> None:
        # A wrapper that has fileno() but raises (e.g. closed fd)
        # must be treated as non-interactive — never crash the
        # renderer just to decide whether to print prose.
        class Broken:
            def fileno(self) -> int:  # noqa: D401 — test stub
                raise OSError("closed")
        self.assertFalse(self.mod._should_emit_similarity_legend(
            "auto", Broken()))


# ---------------------------------------------------------------
# Renderer unit — composition with --with-similarity / --json
# ---------------------------------------------------------------
class RendererLegendComposition(unittest.TestCase):
    def setUp(self) -> None:
        self.mod = load_placeholder_linter()
        sim = self.mod._RenameSimilarity(
            kind="R", score=87, old_path="spec/old.md")
        self.rows = [
            self.mod._ChangedFileAudit(
                path="spec/new.md", status="matched",
                reason="renamed", similarity=sim),
        ]

    def _render(self, **kwargs):  # type: ignore[no-untyped-def]
        buf = io.StringIO()
        defaults = dict(
            as_json=False, with_similarity=True, with_labels=False,
            legend_mode="on",
        )
        defaults.update(kwargs)
        self.mod._render_changed_files_audit(
            self.rows, buf, **defaults)
        return buf.getvalue()

    def test_legend_on_emits_kind_score_old_lines(self) -> None:
        out = self._render(legend_mode="on")
        self.assertIn("legend:", out)
        self.assertIn("kind", out)
        self.assertIn("score", out)
        self.assertIn("old", out)
        # Default (labels off) must NOT mention the meaning column.
        self.assertNotIn("meaning", out)

    def test_legend_with_labels_adds_meaning_line(self) -> None:
        out = self._render(legend_mode="on", with_labels=True)
        self.assertIn("meaning", out)
        self.assertIn("rename-similarity", out)
        self.assertIn("copy-similarity", out)

    def test_legend_off_suppresses_legend(self) -> None:
        out = self._render(legend_mode="off")
        self.assertNotIn("legend:", out)
        # Sanity: the table itself is still there — we suppressed
        # only the trailing prose, not the data rows.
        self.assertIn("spec/new.md", out)
        self.assertIn("totals:", out)

    def test_legend_skipped_without_with_similarity(self) -> None:
        # Even with mode=on, the legend documents columns that
        # aren't in the table — emitting it would be misleading.
        out = self._render(legend_mode="on", with_similarity=False)
        self.assertNotIn("legend:", out)

    def test_legend_skipped_in_json_mode(self) -> None:
        # JSON mode short-circuits before the totals/legend block.
        # Pin that here so a future refactor can't regress and
        # poison the JSON document with a trailing prose line.
        out = self._render(legend_mode="on", as_json=True)
        self.assertNotIn("legend:", out)

    def test_legend_appears_after_totals(self) -> None:
        # Strict ordering: byte-stream consumers parsing rows + the
        # totals footer must see the same prefix; the legend is
        # purely additive at the end.
        out = self._render(legend_mode="on")
        i_totals = out.find("totals:")
        i_legend = out.find("legend:")
        self.assertGreater(i_totals, -1)
        self.assertGreater(i_legend, i_totals,
            msg="legend must be emitted AFTER the totals footer")

    def test_auto_default_quiet_on_pipe_stream(self) -> None:
        # StringIO is the test analogue of a pipe (no TTY). Default
        # ``auto`` must stay quiet — this is the property that
        # protects CI log scrapers from gaining surprise prose.
        buf = io.StringIO()
        self.mod._render_changed_files_audit(
            self.rows, buf,
            as_json=False, with_similarity=True,
            legend_mode="auto",
        )
        self.assertNotIn("legend:", buf.getvalue())


# ---------------------------------------------------------------
# CLI smoke — argparse default + explicit on over a non-TTY pipe
# ---------------------------------------------------------------
class CliSimilarityLegendSmoke(unittest.TestCase):
    def test_default_auto_quiet_under_subprocess_pipe(self) -> None:
        # subprocess.PIPE is never a TTY → ``auto`` must emit
        # nothing, preserving the legacy STDERR byte sequence.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("R087\tspec/intro.md\tspec/copy.md\n")
            code, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--with-similarity",
                cwd=tdp,
            )
            self.assertEqual(code, 0,
                msg=f"unexpected exit; stderr={err!r}")
            self.assertNotIn("legend:", err,
                msg="auto mode under a pipe must not emit the "
                    "legend (would change the legacy byte stream)")
            # Sanity: the audit table itself IS there.
            self.assertIn("spec/copy.md", err)
            self.assertIn("totals:", err)

    def test_explicit_on_emits_legend_under_pipe(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("R087\tspec/intro.md\tspec/copy.md\n")
            code, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--with-similarity",
                "--similarity-legend", "on",
                cwd=tdp,
            )
            self.assertEqual(code, 0,
                msg=f"unexpected exit; stderr={err!r}")
            self.assertIn("legend:", err)
            self.assertIn("kind", err)
            self.assertIn("score", err)

    def test_explicit_off_suppresses_legend(self) -> None:
        # We can't easily attach a real TTY in pytest, but we can
        # at least pin that ``off`` over a pipe stays off (and
        # combined with the resolver unit tests above this nails
        # the off-on-real-TTY case too).
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("R087\tspec/intro.md\tspec/copy.md\n")
            code, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--with-similarity",
                "--similarity-legend", "off",
                cwd=tdp,
            )
            self.assertEqual(code, 0,
                msg=f"unexpected exit; stderr={err!r}")
            self.assertNotIn("legend:", err)

    def test_legend_noop_without_with_similarity(self) -> None:
        # Even ``--similarity-legend on`` must do nothing when the
        # explained columns aren't in the table.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("M\tspec/intro.md\n")
            code, _, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--similarity-legend", "on",
                cwd=tdp,
            )
            self.assertEqual(code, 0,
                msg=f"unexpected exit; stderr={err!r}")
            self.assertNotIn("legend:", err)

    def test_legend_noop_in_json_mode(self) -> None:
        # JSON mode → STDOUT carries the violations doc; STDERR
        # carries a JSON audit array (no totals footer, no legend).
        # Pin that ``--similarity-legend on`` doesn't poison either.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("R087\tspec/intro.md\tspec/copy.md\n")
            code, out, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--with-similarity",
                "--similarity-legend", "on",
                "--json",
                cwd=tdp,
            )
            self.assertEqual(code, 0,
                msg=f"unexpected exit; stdout={out!r} stderr={err!r}")
            self.assertNotIn("legend:", err)
            self.assertNotIn("legend:", out)

    def test_invalid_choice_rejected_by_argparse(self) -> None:
        # argparse must enforce the {auto,on,off} vocabulary so a
        # typo fails loud at the boundary instead of silently
        # falling through to a default branch.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec = _make_repo(tdp)
            payload = tdp / "changed.txt"
            payload.write_text("M\tspec/intro.md\n")
            code, _out, err = _run(
                "--root", str(spec),
                "--changed-files", str(payload),
                "--list-changed-files",
                "--similarity-legend", "yes",  # invalid
                cwd=tdp,
            )
            self.assertEqual(code, 2,
                msg=f"argparse should exit 2; stderr={err!r}")
            self.assertIn("similarity-legend", err)


if __name__ == "__main__":
    unittest.main()
