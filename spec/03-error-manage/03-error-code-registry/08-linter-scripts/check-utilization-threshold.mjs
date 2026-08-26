#!/usr/bin/env node

/**
 * Flags modules whose error code utilization exceeds a threshold.
 * Emits GitHub Actions warnings but does not fail the build.
 * Usage: node spec/07-error-code-registry/scripts/check-utilization-threshold.mjs
 */

import { readFileSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const THRESHOLD = 0.30; // 30%

const __dir = dirname(fileURLToPath(import.meta.url));
const MASTER = resolve(__dir, '../../..', 'spec/07-error-code-registry/error-codes-master.json');
const master = JSON.parse(readFileSync(MASTER, 'utf-8'));

function rangeCapacity(mod) {
  if (mod.Range) return mod.Range.Max - mod.Range.Min + 1;
  if (mod.Ranges) return mod.Ranges.reduce((s, r) => s + (r.Max - r.Min + 1), 0);
  return 0;
}

const warnings = [];

for (const m of master.Modules) {
  const cap = rangeCapacity(m);
  if (cap === 0 || m.TotalCodes === 0) continue;
  const util = m.TotalCodes / cap;
  if (util >= THRESHOLD) {
    const pct = (util * 100).toFixed(1);
    const remaining = cap - m.TotalCodes;
    warnings.push({ Project: m.Project, Name: m.Name, Pct: pct, Used: m.TotalCodes, Cap: cap, Remaining: remaining });
  }
}

if (warnings.length === 0) {
  console.log(`✅ All modules below ${(THRESHOLD * 100).toFixed(0)}% utilization threshold`);
} else {
  console.log(`⚠️  ${warnings.length} module(s) exceed ${(THRESHOLD * 100).toFixed(0)}% utilization:\n`);
  for (const w of warnings) {
    const msg = `${w.Name} (${w.Project}): ${w.Pct}% — ${w.Used}/${w.Cap} used, ${w.Remaining} remaining`;
    console.log(`  • ${msg}`);
    // Emit GitHub Actions warning annotation
    console.log(`::warning::${msg}`);
  }
}
