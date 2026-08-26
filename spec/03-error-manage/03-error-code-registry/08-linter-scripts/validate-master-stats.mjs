#!/usr/bin/env node

/**
 * Validates that error-codes-master.json Stats.TotalIndexedCodes and
 * Stats.TotalRetryableCodes match the sum of per-module values.
 * Exits non-zero if they are stale.
 */

import { readFileSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __dir = dirname(fileURLToPath(import.meta.url));
const MASTER = resolve(__dir, '../../..', 'spec/07-error-code-registry/error-codes-master.json');

const master = JSON.parse(readFileSync(MASTER, 'utf-8'));
const modules = master.Modules;

const computedCodes = modules.reduce((s, m) => s + (m.TotalCodes || 0), 0);
const computedRetryable = modules.reduce((s, m) => s + (m.RetryableCodes || 0), 0);

let failed = false;

if (master.Stats.TotalIndexedCodes !== computedCodes) {
  console.error(`❌ Stats.TotalIndexedCodes is ${master.Stats.TotalIndexedCodes} but sum of modules is ${computedCodes}`);
  failed = true;
}

if (master.Stats.TotalRetryableCodes !== computedRetryable) {
  console.error(`❌ Stats.TotalRetryableCodes is ${master.Stats.TotalRetryableCodes} but sum of modules is ${computedRetryable}`);
  failed = true;
}

if (failed) {
  console.error('\nRun the following to fix:\n  Update Stats.TotalIndexedCodes and Stats.TotalRetryableCodes in error-codes-master.json');
  process.exit(1);
} else {
  console.log(`✅ Master index stats are consistent (${computedCodes} codes, ${computedRetryable} retryable)`);
}
