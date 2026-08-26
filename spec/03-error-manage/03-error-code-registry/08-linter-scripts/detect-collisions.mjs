#!/usr/bin/env node
/**
 * Collision Detection Script for Ecosystem Error Codes
 * 
 * Validates that no two modules have overlapping ecosystem integer codes
 * across all error-codes.json index files.
 * 
 * Usage: node spec/07-error-code-registry/scripts/detect-collisions.mjs
 */

import { readFileSync, existsSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT = resolve(__dirname, '../../..');

const MASTER_INDEX = resolve(ROOT, 'spec/07-error-code-registry/error-codes-master.json');

// ── Helpers ──────────────────────────────────────────────────────────

function loadJson(path) {
  if (!existsSync(path)) return null;
  return JSON.parse(readFileSync(path, 'utf-8'));
}

function extractEcosystemCodes(index) {
  const codes = [];
  if (!index?.Categories) return codes;

  for (const cat of index.Categories) {
    if (!cat.Codes) continue;
    for (const entry of cat.Codes) {
      // Integer code = ecosystem code
      if (typeof entry.Code === 'number') {
        codes.push({
          Code: entry.Code,
          Constant: entry.Constant,
          Category: cat.Name,
          LocalCode: entry.LocalCode ?? null,
        });
      }
      // String codes like "E1001" are local/prefixed — skip unless mapped
    }
  }
  return codes;
}

// ── Main ─────────────────────────────────────────────────────────────

function main() {
  console.log('╔══════════════════════════════════════════════════════════╗');
  console.log('║   Error Code Collision Detection                       ║');
  console.log('╚══════════════════════════════════════════════════════════╝\n');

  const master = loadJson(MASTER_INDEX);
  if (!master) {
    console.error('❌ Master index not found:', MASTER_INDEX);
    process.exit(1);
  }

  // Collect all codes from all modules
  const allCodes = new Map(); // code → [{ Project, Constant, Category, file }]
  const moduleResults = [];
  let totalCodes = 0;
  let filesScanned = 0;

  for (const mod of master.Modules) {
    const indexPath = resolve(ROOT, mod.IndexFile);
    const index = loadJson(indexPath);

    if (!index) {
      moduleResults.push({ Project: mod.Project, Name: mod.Name, Status: 'MISSING', Codes: 0 });
      continue;
    }

    filesScanned++;
    const codes = extractEcosystemCodes(index);
    totalCodes += codes.length;

    moduleResults.push({ Project: mod.Project, Name: mod.Name, Status: 'OK', Codes: codes.length });

    for (const c of codes) {
      if (!allCodes.has(c.Code)) {
        allCodes.set(c.Code, []);
      }
      allCodes.get(c.Code).push({
        Project: mod.Project,
        Constant: c.Constant,
        Category: c.Category,
        File: mod.IndexFile,
      });
    }
  }

  // ── Module summary ───────────────────────────────────────────────

  console.log('Modules scanned:');
  console.log('─'.repeat(60));
  for (const m of moduleResults) {
    const icon = m.Status === 'OK' ? '✅' : m.Status === 'MISSING' ? '⚠️ ' : '❌';
    const codesStr = m.Codes > 0 ? `${m.Codes} codes` : m.Status === 'MISSING' ? 'file missing' : '0 codes';
    console.log(`  ${icon} ${m.Project.padEnd(8)} ${m.Name.padEnd(40)} ${codesStr}`);
  }
  console.log('─'.repeat(60));
  console.log(`  Files: ${filesScanned}/${master.Modules.length}  |  Ecosystem codes: ${totalCodes}\n`);

  // ── Collision detection ──────────────────────────────────────────

  const collisions = [];
  for (const [code, owners] of allCodes) {
    if (owners.length > 1) {
      collisions.push({ Code: code, Owners: owners });
    }
  }

  // ── Range overlap detection ──────────────────────────────────────

  const rangeOverlaps = [];
  const moduleRanges = [];
  for (const mod of master.Modules) {
    const ranges = mod.Ranges ?? (mod.Range ? [mod.Range] : []);
    for (const r of ranges) {
      moduleRanges.push({ Project: mod.Project, Min: r.Min, Max: r.Max });
    }
  }

  for (let i = 0; i < moduleRanges.length; i++) {
    for (let j = i + 1; j < moduleRanges.length; j++) {
      const a = moduleRanges[i];
      const b = moduleRanges[j];
      if (a.Min <= b.Max && b.Min <= a.Max) {
        // Check if it's the known intentional PS/AB overlap
        const knownOverlap =
          (a.Project === 'AB' || b.Project === 'AB') &&
          master.SpecialRanges?.some(s => s.Project === 'PS/AB');
        rangeOverlaps.push({
          A: `${a.Project} [${a.Min}-${a.Max}]`,
          B: `${b.Project} [${b.Min}-${b.Max}]`,
          Intentional: knownOverlap,
        });
      }
    }
  }

  // ── Gap analysis ─────────────────────────────────────────────────

  const sortedRanges = [...moduleRanges].sort((a, b) => a.Min - b.Min);
  const gaps = [];
  for (let i = 0; i < sortedRanges.length - 1; i++) {
    const gapStart = sortedRanges[i].Max + 1;
    const gapEnd = sortedRanges[i + 1].Min - 1;
    if (gapEnd >= gapStart && (gapEnd - gapStart) >= 10) {
      gaps.push({ Min: gapStart, Max: gapEnd, Size: gapEnd - gapStart + 1 });
    }
  }

  // ── Duplicate constant detection (within same module) ────────────

  const dupConstants = [];
  for (const mod of master.Modules) {
    const indexPath = resolve(ROOT, mod.IndexFile);
    const index = loadJson(indexPath);
    if (!index?.Categories) continue;

    const seen = new Map();
    for (const cat of index.Categories) {
      if (!cat.Codes) continue;
      for (const entry of cat.Codes) {
        if (seen.has(entry.Constant)) {
          dupConstants.push({
            Project: mod.Project,
            Constant: entry.Constant,
            First: seen.get(entry.Constant),
            Second: typeof entry.Code === 'number' ? entry.Code : entry.LocalCode ?? entry.Code,
          });
        } else {
          seen.set(entry.Constant, typeof entry.Code === 'number' ? entry.Code : entry.LocalCode ?? entry.Code);
        }
      }
    }
  }

  // ── Results ──────────────────────────────────────────────────────

  console.log('══════════════════════════════════════════════════════════');
  console.log('  RESULTS');
  console.log('══════════════════════════════════════════════════════════\n');

  // Code collisions
  if (collisions.length === 0) {
    console.log('✅ Code Collisions: NONE — all ecosystem integer codes are unique\n');
  } else {
    console.log(`❌ Code Collisions: ${collisions.length} FOUND\n`);
    for (const c of collisions) {
      console.log(`   Code ${c.Code}:`);
      for (const o of c.Owners) {
        console.log(`     → ${o.Project} / ${o.Constant} (${o.Category}) in ${o.File}`);
      }
      console.log();
    }
  }

  // Range overlaps
  if (rangeOverlaps.length === 0) {
    console.log('✅ Range Overlaps: NONE\n');
  } else {
    const intentional = rangeOverlaps.filter(r => r.Intentional);
    const unintentional = rangeOverlaps.filter(r => !r.Intentional);
    if (unintentional.length > 0) {
      console.log(`❌ Range Overlaps: ${unintentional.length} UNINTENTIONAL\n`);
      for (const r of unintentional) {
        console.log(`     ${r.A}  ↔  ${r.B}`);
      }
      console.log();
    }
    if (intentional.length > 0) {
      console.log(`ℹ️  Range Overlaps: ${intentional.length} intentional (format-separated)\n`);
      for (const r of intentional) {
        console.log(`     ${r.A}  ↔  ${r.B}  [known PS/AB SEO overlap]`);
      }
      console.log();
    }
  }

  // Duplicate constants
  if (dupConstants.length === 0) {
    console.log('✅ Duplicate Constants: NONE\n');
  } else {
    console.log(`⚠️  Duplicate Constants: ${dupConstants.length} found\n`);
    for (const d of dupConstants) {
      console.log(`     ${d.Project}: "${d.Constant}" at codes ${d.First} and ${d.Second}`);
    }
    console.log();
  }

  // Gaps
  if (gaps.length > 0) {
    console.log(`ℹ️  Unallocated Gaps (≥10 codes):`);
    for (const g of gaps) {
      console.log(`     [${g.Min}-${g.Max}] (${g.Size} codes)`);
    }
    console.log();
  }

  // Final verdict
  const hasFailures = collisions.length > 0 || rangeOverlaps.some(r => !r.Intentional);
  console.log('══════════════════════════════════════════════════════════');
  if (hasFailures) {
    console.log('  ❌ VALIDATION FAILED — collisions detected');
    process.exit(1);
  } else {
    console.log('  ✅ VALIDATION PASSED — no collisions detected');
    process.exit(0);
  }
}

main();
