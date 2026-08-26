#!/usr/bin/env node

/**
 * generate-dashboard-data.cjs
 *
 * Scans the spec/ tree to:
 *   1. Validate all markdown cross-references (broken link detection)
 *   2. Check for required files (00-overview.md, 99-consistency-report.md)
 *   3. Count files per subfolder
 *   4. Output a JSON report to spec/dashboard-data.json
 *
 * Usage:  node linter-scripts/generate-dashboard-data.cjs [--json] [--quiet]
 */

const fs = require("fs");
const path = require("path");

// ── CLI flags ───────────────────────────────────────────────
const args = process.argv.slice(2);
const jsonOnly = args.includes("--json");
const quiet = args.includes("--quiet");

const SPEC_ROOT = path.resolve(__dirname, "..", "spec");
const ARCHIVE_SEGMENTS = ["_archive", "archive"];

// ── Helpers ─────────────────────────────────────────────────

function isArchivePath(filePath) {
  const segments = filePath.split(path.sep);
  return segments.some((s) => ARCHIVE_SEGMENTS.includes(s.toLowerCase()));
}

function walkMarkdown(dir) {
  const results = [];
  if (!fs.existsSync(dir)) return results;

  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      results.push(...walkMarkdown(full));
    } else if (entry.name.endsWith(".md") && !isArchivePath(full)) {
      results.push(full);
    }
  }

  return results;
}

function walkDirs(dir) {
  const results = [];
  if (!fs.existsSync(dir)) return results;

  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory() && !isArchivePath(full)) {
      results.push(full);
      results.push(...walkDirs(full));
    }
  }

  return results;
}

// ── 1. Broken link detection ────────────────────────────────

const LINK_RE = /\[([^\]]*)\]\((\.[^)]+)\)/g;

function extractLinks(filePath, content) {
  const links = [];
  const lines = content.split("\n");
  let inCodeBlock = false;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line.trimStart().startsWith("```")) {
      inCodeBlock = !inCodeBlock;
      continue;
    }
    if (inCodeBlock) continue;

    let match;
    LINK_RE.lastIndex = 0;
    while ((match = LINK_RE.exec(line)) !== null) {
      const linkTarget = match[2];
      const filePart = linkTarget.split("#")[0];
      if (!filePart) continue;

      links.push({
        Line: i + 1,
        Text: match[1],
        Target: linkTarget,
        FilePart: filePart,
      });
    }
  }

  return links;
}

function validateLinks(mdFiles) {
  const broken = [];
  const total = { Checked: 0, Ok: 0, Broken: 0 };

  for (const filePath of mdFiles) {
    const content = fs.readFileSync(filePath, "utf8");
    const links = extractLinks(filePath, content);

    for (const link of links) {
      total.Checked++;
      const resolved = path.resolve(path.dirname(filePath), link.FilePart);

      if (fs.existsSync(resolved)) {
        total.Ok++;
      } else {
        total.Broken++;
        broken.push({
          Source: path.relative(SPEC_ROOT, filePath),
          Line: link.Line,
          Text: link.Text,
          Target: link.Target,
          Resolved: path.relative(SPEC_ROOT, resolved),
        });
      }
    }
  }

  return { Broken: broken, Total: total };
}

// ── 2. Required-file checks ────────────────────────────────

function checkRequiredFiles(dirs) {
  const missingOverview = [];
  const missingConsistency = [];

  for (const dir of dirs) {
    const files = fs.readdirSync(dir).filter((f) => f.endsWith(".md"));
    const rel = path.relative(SPEC_ROOT, dir);

    if (files.length === 0) continue;

    const hasOverview = files.includes("00-overview.md");
    const hasConsistency = files.includes("99-consistency-report.md");

    if (!hasOverview && files.length >= 2) {
      missingOverview.push({ Folder: rel, FileCount: files.length });
    }

    if (!hasConsistency && files.length >= 3) {
      missingConsistency.push({ Folder: rel, FileCount: files.length });
    }
  }

  return { MissingOverview: missingOverview, MissingConsistency: missingConsistency };
}

// ── 3. Folder inventory ─────────────────────────────────────

function buildInventory(dirs) {
  const inventory = [];

  for (const dir of dirs) {
    const allEntries = fs.readdirSync(dir, { withFileTypes: true });
    const mdFiles = allEntries.filter(
      (e) => e.isFile() && e.name.endsWith(".md")
    );
    const subDirs = allEntries.filter((e) => e.isDirectory());
    const rel = path.relative(SPEC_ROOT, dir);

    if (mdFiles.length === 0 && subDirs.length === 0) continue;

    inventory.push({
      Folder: rel,
      MdFiles: mdFiles.length,
      Subfolders: subDirs.length,
      HasOverview: mdFiles.some((f) => f.name === "00-overview.md"),
      HasConsistency: mdFiles.some((f) => f.name === "99-consistency-report.md"),
      HasChangelog: mdFiles.some((f) => f.name === "98-changelog.md"),
      HasAcceptance: mdFiles.some((f) => f.name === "97-acceptance-criteria.md"),
      Files: mdFiles.map((f) => f.name).sort(),
    });
  }

  return inventory;
}

// ── 4. Health score ─────────────────────────────────────────

function computeHealth(linkResult, requiredFiles, inventory) {
  let score = 100;
  const deductions = [];

  if (linkResult.Total.Broken > 0) {
    const d = Math.min(linkResult.Total.Broken * 2, 20);
    score -= d;
    deductions.push(
      `${linkResult.Total.Broken} broken links (-${d})`
    );
  }

  if (requiredFiles.MissingConsistency.length > 0) {
    const d = Math.min(requiredFiles.MissingConsistency.length, 15);
    score -= d;
    deductions.push(
      `${requiredFiles.MissingConsistency.length} missing consistency reports (-${d})`
    );
  }

  if (requiredFiles.MissingOverview.length > 0) {
    const d = Math.min(requiredFiles.MissingOverview.length * 3, 15);
    score -= d;
    deductions.push(
      `${requiredFiles.MissingOverview.length} missing overviews (-${d})`
    );
  }

  const grade =
    score >= 95 ? "A+" :
    score >= 90 ? "A" :
    score >= 85 ? "B+" :
    score >= 80 ? "B" :
    score >= 70 ? "C" :
    score >= 60 ? "D" : "F";

  return { Score: Math.max(score, 0), Grade: grade, Deductions: deductions };
}

// ── Main ────────────────────────────────────────────────────

function main() {
  const mdFiles = walkMarkdown(SPEC_ROOT);
  const allDirs = [SPEC_ROOT, ...walkDirs(SPEC_ROOT)];

  const linkResult = validateLinks(mdFiles);
  const requiredFiles = checkRequiredFiles(allDirs);
  const inventory = buildInventory(allDirs);
  const health = computeHealth(linkResult, requiredFiles, inventory);

  const report = {
    Generated: new Date().toISOString().slice(0, 10),
    Health: health,
    Links: {
      TotalChecked: linkResult.Total.Checked,
      Ok: linkResult.Total.Ok,
      Broken: linkResult.Total.Broken,
      BrokenDetails: linkResult.Broken,
    },
    RequiredFiles: {
      MissingOverview: requiredFiles.MissingOverview,
      MissingConsistency: requiredFiles.MissingConsistency,
    },
    Inventory: {
      TotalFolders: inventory.length,
      TotalMdFiles: mdFiles.length,
      Folders: inventory,
    },
  };

  const outPath = path.join(SPEC_ROOT, "dashboard-data.json");
  fs.writeFileSync(outPath, JSON.stringify(report, null, 2) + "\n");

  if (jsonOnly) {
    process.stdout.write(JSON.stringify(report, null, 2) + "\n");
    return;
  }

  if (!quiet) {
    console.log("╔══════════════════════════════════════════════════╗");
    console.log("║        SPEC HEALTH DASHBOARD GENERATOR          ║");
    console.log("╚══════════════════════════════════════════════════╝\n");

    console.log(`  Health Score:  ${health.Score}/100 (${health.Grade})`);
    if (health.Deductions.length > 0) {
      health.Deductions.forEach((d) => console.log(`    └─ ${d}`));
    }

    console.log(`\n  Files scanned: ${mdFiles.length}`);
    console.log(`  Folders:       ${inventory.length}`);
    console.log(
      `  Links checked: ${linkResult.Total.Checked} (${linkResult.Total.Ok} ok, ${linkResult.Total.Broken} broken)`
    );

    if (linkResult.Broken.length > 0) {
      console.log("\n  ── Broken Links ──────────────────────────────");
      for (const b of linkResult.Broken) {
        console.log(`    ${b.Source}:${b.Line}`);
        console.log(`      → ${b.Target}`);
      }
    }

    if (requiredFiles.MissingConsistency.length > 0) {
      console.log("\n  ── Missing 99-consistency-report.md ──────────");
      for (const m of requiredFiles.MissingConsistency) {
        console.log(`    ${m.Folder}/ (${m.FileCount} files)`);
      }
    }

    if (requiredFiles.MissingOverview.length > 0) {
      console.log("\n  ── Missing 00-overview.md ────────────────────");
      for (const m of requiredFiles.MissingOverview) {
        console.log(`    ${m.Folder}/ (${m.FileCount} files)`);
      }
    }

    console.log(`\n  Output: ${path.relative(process.cwd(), outPath)}`);
  }

  process.exit(linkResult.Total.Broken > 0 ? 1 : 0);
}

main();
