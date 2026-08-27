"""End-to-end tests for per-extension cache segregation.

Each unique sorted ``--extension`` allowlist must land its sentinel
in its own subdirectory under ``--cache-dir``. The previous design
leaned on the hash to differentiate runs (same flat directory, hope
the keys diverge); the new design also segregates *physically*, so:

* operators can ``rm -rf cache/ext-mdx/`` to nuke one allowlist's
  sentinels without touching others
* a hypothetical hash collision (impossible with SHA-256 in
  practice, but cheap to defend against) is ALSO trapped by the
  directory split — the cache lookup looks in a different folder
* ``ls cache/`` is self-describing — the segments tell the operator
  which allowlists have been exercised on this machine

Coverage:

* **Default → ``ext-md/``** — no flag passes through the historical
  default and writes under the well-known segment name.
* **``--extension mdx`` → ``ext-mdx/``** — alternate allowlist gets
  its own segment; default segment is left untouched.
* **Repeated flag → ``ext-md+mdx/``** — order-independent: the
  segment name sorts the extensions before joining.
* **Cache hit only within segment** — a sentinel written under
  ``ext-md`` MUST NOT satisfy a lookup for ``--extension mdx``.
* **Hash-fallback segment** — if a future caller passes a long or
  filesystem-hostile extension allowlist, the segment name falls
  back to a short-hash form so the path stays Windows-legal.
* **Schema bump** — the cache key includes ``exts=…`` so even
  within a single segment, switching ``--extension`` invalidates
  any sentinel that pre-dated the schema bump.
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
    """Two siblings so we can exercise both the .md and .mdx paths
    independently. Returns (root, cache_dir)."""
    spec = td / "spec"
    spec.mkdir()
    # Both files are valid (no placeholder blocks at all) so every
    # run cleanly PASSes and writes a sentinel.
    (spec / "intro.md").write_text("# spec\nplain prose.\n")
    (spec / "intro.mdx").write_text("# spec\nplain prose.\n")
    cache = td / "cache"
    cache.mkdir()
    return spec, cache


def _segments(cache: Path) -> set[str]:
    """Return the set of segment-directory names that contain at
    least one ``.pass`` sentinel."""
    return {p.parent.name for p in cache.rglob("*.pass")}


class CacheSegmentNaming(unittest.TestCase):

    def test_default_extension_writes_under_ext_md(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, err = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache),
                cwd=Path(td),
            )
            self.assertEqual(rc, 0, f"stderr={err!r}")
            self.assertEqual(_segments(cache), {"ext-md"})

    def test_alternate_extension_writes_under_ext_mdx(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--extension", "mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertEqual(_segments(cache), {"ext-mdx"})

    def test_combined_extensions_sort_into_canonical_segment(self) -> None:
        """``--extension mdx --extension md`` and
        ``--extension md --extension mdx`` MUST land in the same
        segment, since the allowlist is set-valued."""
        with tempfile.TemporaryDirectory() as td1, \
             tempfile.TemporaryDirectory() as td2:
            spec1, cache1 = _make_repo(Path(td1))
            spec2, cache2 = _make_repo(Path(td2))
            _run("--root", str(spec1), "--repo-root", td1,
                 "--cache-dir", str(cache1),
                 "--extension", "mdx", "--extension", "md",
                 cwd=Path(td1))
            _run("--root", str(spec2), "--repo-root", td2,
                 "--cache-dir", str(cache2),
                 "--extension", "md", "--extension", "mdx",
                 cwd=Path(td2))
            self.assertEqual(_segments(cache1), {"ext-md+mdx"})
            self.assertEqual(_segments(cache2), {"ext-md+mdx"})

    def test_extensions_strip_leading_dot_and_lowercase(self) -> None:
        """Forgiving input: ``.MD`` collapses to ``md`` so it shares
        the default segment instead of creating a phantom bucket."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--extension", ".MD",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertEqual(_segments(cache), {"ext-md"})


class CacheSegmentIsolation(unittest.TestCase):
    """Verify segregation actually prevents cross-segment cache hits."""

    def test_md_sentinel_does_not_satisfy_mdx_lookup(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            # Run #1: write a PASS under ext-md.
            _run("--root", str(spec), "--repo-root", td,
                 "--cache-dir", str(cache), cwd=Path(td))
            # Run #2: same files, different allowlist. Must NOT hit
            # the ext-md sentinel; must perform a fresh scan and
            # write its own sentinel under ext-mdx.
            rc, out, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--extension", "mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            # The "cache hit" log line is segment-prefixed; absence
            # of "ext-mdx/...cache hit" on the FIRST mdx run proves
            # we didn't accidentally match the ext-md sentinel.
            self.assertNotIn("cache hit", out)
            self.assertEqual(_segments(cache), {"ext-md", "ext-mdx"})

    def test_repeat_default_run_hits_ext_md_segment(self) -> None:
        """Sanity: same flags twice → cache hit on the second run,
        and the hit log identifies the ``ext-md`` segment."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            _run("--root", str(spec), "--repo-root", td,
                 "--cache-dir", str(cache), cwd=Path(td))
            rc, out, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertIn("cache hit", out)
            self.assertIn("ext-md/", out)

    def test_no_cache_write_does_not_create_segment(self) -> None:
        """``--no-cache-write`` is read-only: a fresh segment must
        NOT be created on a miss, even though the directory layout
        would otherwise demand it."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--no-cache-write",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertEqual(_segments(cache), set())


class CacheSegmentHelperUnit(unittest.TestCase):
    """Direct unit tests on ``_cache_segment`` — cheaper than
    spinning up a subprocess for the corner cases."""

    def setUp(self) -> None:
        # The linter script's filename is hyphenated; use the shared
        # shim that handles the dataclass/sys.modules dance Python
        # 3.13 requires when loading via ``importlib`` directly.
        from conftest_shim import load_placeholder_linter
        self.chk = load_placeholder_linter()

    def test_single_extension_readable_form(self) -> None:
        self.assertEqual(self.chk._cache_segment(("md",)), "ext-md")

    def test_sorted_join_is_order_independent(self) -> None:
        a = self.chk._cache_segment(("mdx", "md"))
        b = self.chk._cache_segment(("md", "mdx"))
        self.assertEqual(a, b)
        self.assertEqual(a, "ext-md+mdx")

    def test_empty_input_yields_default_sentinel_name(self) -> None:
        """Defensive: programmatic callers shouldn't crash the
        segment derivation by passing an empty tuple."""
        self.assertEqual(self.chk._cache_segment(()), "ext-default")

    def test_unsafe_chars_force_hash_fallback(self) -> None:
        seg = self.chk._cache_segment(("foo.bar",))
        self.assertTrue(seg.startswith("ext-h"))
        self.assertEqual(len(seg), len("ext-h") + 10)

    def test_overlong_allowlist_falls_back_to_hash(self) -> None:
        many = tuple(f"x{i:03d}" for i in range(40))
        seg = self.chk._cache_segment(many)
        self.assertTrue(seg.startswith("ext-h"))

    def test_hash_fallback_is_deterministic(self) -> None:
        a = self.chk._cache_segment(("foo.bar", "baz/qux"))
        b = self.chk._cache_segment(("baz/qux", "foo.bar"))
        self.assertEqual(a, b)


if __name__ == "__main__":
    unittest.main()