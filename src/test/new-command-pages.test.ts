import { describe, it, expect } from "vitest";
import { readFileSync, existsSync } from "node:fs";
import { resolve } from "node:path";

/**
 * new-command-pages.test.ts — regression guards for the v4.31.0
 * docs additions (commit-in, replace, fix-repo, clone-fix-repo,
 * make-public). Ensures every new command page:
 *
 *   1. Exists on disk under src/pages/
 *   2. Is imported and routed in src/App.tsx
 *   3. Has a corresponding link in the DocsSidebar nav list
 *   4. Renders the expected <h1> command name
 *
 * Pure source-scan — no React render. Fast, deterministic, and
 * catches the failure mode where a page is created but the route
 * or sidebar entry is forgotten (the page is then orphaned and the
 * sidebar link 404s).
 */

interface NewCmd {
  /** Filename under src/pages (no extension). */
  page: string;
  /** Route path registered in App.tsx. */
  route: string;
  /** Visible <h1> text on the page. */
  heading: string;
  /** Sidebar nav title (substring match against DocsSidebar.tsx). */
  sidebarTitle: string;
}

const NEW_COMMANDS: NewCmd[] = [
  { page: "CommitIn",     route: "/commit-in",      heading: "commit-in",      sidebarTitle: "Commit In (cin)" },
  { page: "Replace",      route: "/replace",        heading: "replace",        sidebarTitle: "Replace (rpl)" },
  { page: "FixRepo",      route: "/fix-repo",       heading: "fix-repo",       sidebarTitle: "Fix Repo (fr)" },
  { page: "CloneFixRepo", route: "/clone-fix-repo", heading: "clone-fix-repo", sidebarTitle: "Clone + Fix Repo (cfr)" },
  { page: "MakePublic",   route: "/make-public",    heading: "make-public",    sidebarTitle: "Make Public Repo" },
];

const root = resolve(__dirname, "..");
const appSource = readFileSync(resolve(root, "App.tsx"), "utf8");
const sidebarSource = readFileSync(
  resolve(root, "components/docs/DocsSidebar.tsx"),
  "utf8",
);

describe("new command pages (v4.31.0)", () => {
  for (const cmd of NEW_COMMANDS) {
    describe(cmd.page, () => {
      const pagePath = resolve(root, `pages/${cmd.page}.tsx`);

      it("page file exists", () => {
        expect(existsSync(pagePath), `missing src/pages/${cmd.page}.tsx`).toBe(true);
      });

      it("page renders the expected <h1> heading", () => {
        const src = readFileSync(pagePath, "utf8");
        // Match the literal heading inside any <h1 ...>{heading}</h1>.
        // Tolerates className attrs and surrounding whitespace.
        const re = new RegExp(`<h1[^>]*>\\s*${escapeRegex(cmd.heading)}\\s*</h1>`);
        expect(src, `<h1>${cmd.heading}</h1> not found in ${cmd.page}.tsx`).toMatch(re);
      });

      it("is imported in App.tsx", () => {
        // e.g. import ReplacePage from "./pages/Replace";
        const re = new RegExp(`from\\s+["']\\./pages/${cmd.page}["']`);
        expect(appSource, `App.tsx missing import for ${cmd.page}`).toMatch(re);
      });

      it("has a Route entry in App.tsx", () => {
        const re = new RegExp(`path=["']${escapeRegex(cmd.route)}["']`);
        expect(appSource, `App.tsx missing <Route path="${cmd.route}">`).toMatch(re);
      });

      it("has a sidebar nav entry pointing at the route", () => {
        // Look for an object literal containing both the title and url.
        // The DocsSidebar nav array uses { title: "...", url: "/path", ... }.
        const titleRe = new RegExp(
          `title:\\s*["']${escapeRegex(cmd.sidebarTitle)}["']`,
        );
        const urlRe = new RegExp(`url:\\s*["']${escapeRegex(cmd.route)}["']`);
        expect(sidebarSource, `Sidebar missing title "${cmd.sidebarTitle}"`).toMatch(titleRe);
        expect(sidebarSource, `Sidebar missing url "${cmd.route}"`).toMatch(urlRe);
      });
    });
  }

  it("App.tsx wildcard NotFound route stays last (so new routes are reachable)", () => {
    const wildcardIdx = appSource.lastIndexOf('path="*"');
    expect(wildcardIdx).toBeGreaterThan(-1);
    for (const cmd of NEW_COMMANDS) {
      const routeIdx = appSource.indexOf(`path="${cmd.route}"`);
      expect(routeIdx, `route ${cmd.route} not found`).toBeGreaterThan(-1);
      expect(
        routeIdx,
        `route ${cmd.route} declared after wildcard NotFound — would be unreachable`,
      ).toBeLessThan(wildcardIdx);
    }
  });
});

function escapeRegex(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
