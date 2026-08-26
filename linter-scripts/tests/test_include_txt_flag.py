"""End-to-end tests for the ``--include-txt`` convenience flag.

The flag mirrors ``--include-mdx`` exactly — same union semantics,
same baseline interaction, same cache-segment canonicalisation —
just with ``txt`` substituted for ``mdx``. Tests under exam:

* **Default + flag** — without any explicit ``--extension``, the
  baseline stays ``md`` and the union becomes ``("md", "txt")``.
* **Explicit ``--extension`` + flag** — the flag *adds* ``txt`` to
  whatever the user already passed; it never replaces the baseline.
  In particular, ``--extension rst --include-txt`` drops ``md``
  because the user explicitly took control of the baseline (the
  documented "you opt out of `md` when you pass --extension" rule).
* **Idempotence** — passing both ``--include-txt`` and
  ``--extension txt`` is a no-op (txt is deduped, segment unchanged).
* **Composition with --include-mdx** — both flags coexist and
  produce a 3-element baseline ``("md", "mdx", "txt")``; the cache
  segment sorts independently so flag order doesn't fork the cache.
* **Off by default** — the legacy ``.md``-only baseline is preserved
  for any invocation that omits the flag.

The first three are subprocess tests against the CLI so we exercise
argparse + the resolution block end-to-end. The composition test
leans on the cache-hit log line (the same black-box signal the
include-mdx suite uses) to assert two different invocations
produced byte-identical sentinels.
"""

from __future__ import annotations

import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")


def _run(*args: str, cwd: Path) -> tuple[int, str, str]:
    r = subprocess.run([sys.executable, str(LINTER), *args],
                       cwd=cwd, capture_output=True, text=True)
    return r.returncode, r.stdout, r.stderr


def _make_repo(td: Path) -> tuple[Path, Path]:
    """spec/intro.md + spec/notes.txt, both clean, plus a fresh
    cache dir. Mirrors the layout used by the include-mdx suite so
    both flag tests exercise the same fixture shape."""
    spec = td / "spec"
    spec.mkdir()
    (spec / "intro.md").write_text("# spec\nplain prose.\n")
    (spec / "notes.txt").write_text("# spec\nplain prose.\n")
    cache = td / "cache"
    cache.mkdir()
    return spec, cache


def _segments(cache: Path) -> set[str]:
    """Return the set of cache-segment names that hold a sentinel.

    Sentinels live one directory deep under the cache root; the
    segment name encodes the canonicalised extension allowlist
    (e.g. ``ext-md+txt`` for the union). We collect distinct parent
    names so a test can assert "two invocations wrote into the same
    segment" without depending on the exact sentinel filename."""
    return {p.parent.name for p in cache.rglob("*.pass")}


class IncludeTxtBaseline(unittest.TestCase):
    """Bare ``--include-txt`` keeps the historical ``md`` baseline."""

    def test_default_baseline_unioned_with_txt(self) -> None:
        # No --extension flags → baseline (md,) AND --include-txt →
        # final allowlist (md, txt). The .txt fixture must therefore
        # be discovered + linted (exits clean since it's well-formed).
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec, cache = _make_repo(tdp)
            code, _, _ = _run(
                "--root", str(spec),
                "--repo-root", str(tdp),
                "--include-txt",
                "--cache-dir", str(cache),
                cwd=tdp)
            self.assertEqual(code, 0)
            # The cache segment encodes the allowlist; ``md+txt``
            # proves both extensions reached the iterator (a missing
            # baseline would have produced ``ext-txt``). Asserted
            # INSIDE the ``with`` so the TemporaryDirectory hasn't
            # been cleaned up by the time ``_segments`` walks it.
            self.assertEqual(_segments(cache), {"ext-md+txt"})

    def test_off_by_default_keeps_md_only_baseline(self) -> None:
        # Without the flag, the .txt fixture must NOT be scanned —
        # the segment is the legacy ``ext-md`` and a hypothetical
        # .txt-only violation would slip past silently. This is the
        # contract every CI relying on the historical default expects.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec, cache = _make_repo(tdp)
            code, _, _ = _run(
                "--root", str(spec),
                "--repo-root", str(tdp),
                "--cache-dir", str(cache),
                cwd=tdp)
            self.assertEqual(code, 0)
            self.assertEqual(_segments(cache), {"ext-md"})


class IncludeTxtUnionSemantics(unittest.TestCase):
    """Flag composes with ``--extension`` and with ``--include-mdx``."""

    def test_explicit_extension_drops_md_baseline(self) -> None:
        # Documented gotcha (now also in --help): passing
        # ``--extension rst`` REPLACES the (md,) baseline with (rst,);
        # ``--include-txt`` then unions ``txt`` onto THAT, yielding
        # (rst, txt) — NOT (md, rst, txt). The cache segment proves
        # which set the iterator actually used.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec, cache = _make_repo(tdp)
            code, _, _ = _run(
                "--root", str(spec),
                "--repo-root", str(tdp),
                "--extension", "rst",
                "--include-txt",
                "--cache-dir", str(cache),
                cwd=tdp)
            self.assertEqual(code, 0)
            # Sorted canonicalisation in segment: ``rst+txt`` (no md).
            self.assertEqual(_segments(cache), {"ext-rst+txt"})

    def test_idempotent_with_explicit_txt_extension(self) -> None:
        # ``--include-txt`` + ``--extension md --extension txt`` must
        # collapse to the same canonical segment as either form on
        # its own. The dedupe in the resolution block + the sort in
        # ``_cache_segment`` together guarantee this; we assert the
        # observable cache-segment side of that contract.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec, cache = _make_repo(tdp)
            code, _, _ = _run(
                "--root", str(spec),
                "--repo-root", str(tdp),
                "--extension", "md",
                "--extension", "txt",
                "--include-txt",
                "--cache-dir", str(cache),
                cwd=tdp)
            self.assertEqual(code, 0)
            self.assertEqual(_segments(cache), {"ext-md+txt"})

    def test_composes_with_include_mdx(self) -> None:
        # Both flags + default baseline → (md, mdx, txt). The segment
        # is sorted so the canonical form is ``ext-md+mdx+txt``
        # regardless of which ``--include-*`` flag came first on the
        # CLI. We run the flags in BOTH orders against a SHARED
        # cache and assert the second run is a cache hit — proof
        # that the segment + key are byte-identical.
        with tempfile.TemporaryDirectory() as td:
            tdp = Path(td)
            spec, cache = _make_repo(tdp)
            (spec / "extra.mdx").write_text("# spec\nplain prose.\n")
            code1, out1, _ = _run(
                "--root", str(spec),
                "--repo-root", str(tdp),
                "--include-mdx", "--include-txt",
                "--cache-dir", str(cache),
                cwd=tdp)
            code2, out2, _ = _run(
                "--root", str(spec),
                "--repo-root", str(tdp),
                "--include-txt", "--include-mdx",
                "--cache-dir", str(cache),
                cwd=tdp)
            self.assertEqual(code1, 0)
            self.assertEqual(code2, 0)
            # Single canonical segment after both runs — proves the
            # flag order didn't fork the cache.
            self.assertEqual(_segments(cache), {"ext-md+mdx+txt"})
            # Second run hit the cache (sentinel from run 1). The
            # cache-hit log line goes to STDOUT (not STDERR) — it
            # replaces the normal scan summary, so STDOUT is the
            # right stream to grep. The first run wrote the sentinel
            # AFTER its scan, so its STDOUT is the success summary;
            # only run 2's STDOUT carries "cache hit".
            self.assertIn("cache hit", out2.lower())
            self.assertNotIn("cache hit", out1.lower())


if __name__ == "__main__":  # pragma: no cover - manual runner
    unittest.main()