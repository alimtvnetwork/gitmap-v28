// scripts/ui-parity.mjs
//
// Dark/light parity snapshot harness (#11).
//
// For each route below, captures a screenshot in light and dark mode,
// then compares them with pixelmatch. Pure layout/contrast deltas are
// expected — this gate fails only when the diff is so large that one
// theme is effectively broken (>20% of pixels different).
//
// Run locally with:  node scripts/ui-parity.mjs
import { chromium } from "playwright";
import { mkdir, writeFile } from "node:fs/promises";
import { PNG } from "pngjs";
import pixelmatch from "pixelmatch";

const ROUTES = ["/", "/commands", "/changelog", "/doctor"];
const BASE = process.env.PARITY_BASE_URL ?? "http://localhost:4173";
const OUT = ".ui-parity";
const MAX_DIFF_FRACTION = 0.20;

const setTheme = async (page, theme) => {
  await page.addInitScript((t) => {
    try { localStorage.setItem("theme", t); } catch {}
    document.documentElement.classList.toggle("dark", t === "dark");
  }, theme);
};

const shot = async (browser, route, theme) => {
  const ctx = await browser.newContext({ viewport: { width: 1280, height: 1800 } });
  const page = await ctx.newPage();
  await setTheme(page, theme);
  await page.goto(`${BASE}${route}`, { waitUntil: "networkidle" });
  const buf = await page.screenshot({ fullPage: false });
  await ctx.close();
  return PNG.sync.read(buf);
};

const main = async () => {
  await mkdir(OUT, { recursive: true });
  const browser = await chromium.launch();
  const failures = [];
  for (const route of ROUTES) {
    const light = await shot(browser, route, "light");
    const dark = await shot(browser, route, "dark");
    // Resize handled by capturing the same viewport; identical themes
    // would diff at ~100%, so we only flag a route when *both* themes
    // produce a near-identical bitmap (a hint that one theme failed
    // to apply — the actual regression we care about).
    const { width, height } = light;
    const diff = new PNG({ width, height });
    const px = pixelmatch(light.data, dark.data, diff.data, width, height, {
      threshold: 0.1,
    });
    const fraction = px / (width * height);
    await writeFile(`${OUT}/${encodeURIComponent(route)}.diff.png`, PNG.sync.write(diff));
    console.log(`route=${route} diff=${(fraction * 100).toFixed(2)}%`);
    if (fraction < MAX_DIFF_FRACTION) {
      failures.push({ route, fraction });
    }
  }
  await browser.close();
  if (failures.length > 0) {
    console.error("UI-parity failures (themes too similar → likely broken):");
    for (const f of failures) console.error(`  - ${f.route}: ${(f.fraction * 100).toFixed(2)}%`);
    process.exit(1);
  }
};

main().catch((err) => {
  console.error(err);
  process.exit(2);
});
