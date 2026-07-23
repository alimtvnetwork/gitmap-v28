import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";
import { componentTagger } from "lovable-tagger";
// @ts-expect-error - plain .mjs helper, no types
import { findCaseCollisions } from "./scripts/check-case-collisions.mjs";

type Collision = { dir: string; a: string; b: string };

// Fails fast when two files in the same directory differ only by case
// (silent on Windows/macOS, fatal on Linux CI — see scripts/check-case-collisions.mjs).
function caseCollisionGuard() {
  return {
    name: "case-collision-guard",
    buildStart() {
      const found = findCaseCollisions() as Collision[];
      if (!found.length) return;
      const lines = found.map(
        (c) => `  ${c.dir}/  ->  ${c.a}  vs  ${c.b}`,
      );
      throw new Error(
        "Case-only filename collisions detected (break Windows builds):\n" +
          lines.join("\n") +
          "\nRename one of each pair with `git mv` to preserve history.",
      );
    },
  };
}

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  server: {
    host: "::",
    port: 8080,
    hmr: {
      overlay: false,
    },
  },
  plugins: [
    caseCollisionGuard(),
    react(),
    mode === "development" && componentTagger(),
  ].filter(Boolean),
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
}));
