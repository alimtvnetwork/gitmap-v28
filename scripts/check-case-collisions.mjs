// Detects case-only filename collisions that break Windows builds.
// Two paths differing only in casing (e.g. DocsTooltip.tsx vs docsTooltip.tsx
// in the same directory) silently work on case-insensitive filesystems
// (Windows, macOS default) but cause hard import failures on case-sensitive
// ones (Linux CI) — and vice versa after a `git mv` rename. Fail fast.
import { readdirSync, statSync } from "node:fs";
import { join, relative } from "node:path";

const ROOT = new URL("../src/", import.meta.url).pathname;
const SKIP = new Set(["node_modules", ".git", "dist", "build"]);

function walk(dir, collisions) {
  const entries = readdirSync(dir);
  const seen = new Map(); // lowercase -> original
  for (const name of entries) {
    if (SKIP.has(name)) continue;
    const lower = name.toLowerCase();
    if (seen.has(lower) && seen.get(lower) !== name) {
      collisions.push({
        dir: relative(ROOT, dir) || ".",
        a: seen.get(lower),
        b: name,
      });
    } else {
      seen.set(lower, name);
    }
    const full = join(dir, name);
    if (statSync(full).isDirectory()) walk(full, collisions);
  }
  return collisions;
}

export function findCaseCollisions(root = ROOT) {
  return walk(root, []);
}

if (import.meta.url === `file://${process.argv[1]}`) {
  const found = findCaseCollisions();
  if (found.length) {
    console.error("\n[case-collision] Case-only filename collisions detected:");
    for (const c of found) {
      console.error(`  ${c.dir}/  ->  ${c.a}  vs  ${c.b}`);
    }
    console.error(
      "\nThese break the build on Windows / case-insensitive filesystems.\n" +
        "Rename one of each pair (use `git mv old new` to preserve history).\n",
    );
    process.exit(1);
  }
  console.log("[case-collision] OK — no case-only filename collisions in src/");
}
