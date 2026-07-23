import { COLORS } from "./theme";
import type { TerminalLine } from "./Terminal";

const C = COLORS;

// Compact helpers
const cmd = (cwd: string, ...tokens: { t: string; c?: string }[]): TerminalLine => ({
  kind: "prompt", cwd, tokens,
});
const out = (...tokens: { t: string; c?: string }[]): TerminalLine => ({
  kind: "out", tokens,
});
const blank = (): TerminalLine => ({ kind: "blank" });

const T = (t: string, c?: string) => ({ t, c });

export type CommandDemo = { caption: string; cwd: string; lines: TerminalLine[] };

export const COMMANDS: CommandDemo[] = [
  // 1) clone
  {
    caption: "Clone any repo, anywhere",
    cwd: "~/repos",
    lines: [
      cmd("~/repos", T("gitmap ", C.command), T("clone ", C.brandGold), T("https://github.com/alimtvnetwork/gitmap-v27", C.path)),
      out(T("→ Resolving remote… ", C.muted), T("OK", C.success)),
      out(T("→ Cloning into ", C.muted), T("./gitmap-v27", C.path), T(" (default branch: ", C.muted), T("main", C.flag), T(")", C.muted)),
      out(T("→ Registering with GitHub Desktop ✓", C.success)),
      out(T("✔ done in 2.4s", C.success)),
    ],
  },
  // 2) scan
  {
    caption: "Scan a folder tree for repos",
    cwd: "~/repos",
    lines: [
      cmd("~/repos", T("gitmap ", C.command), T("scan ", C.brandGold), T(". ", C.path), T("--json --csv", C.flag)),
      out(T("scanning ", C.muted), T("~/repos", C.path), T(" …", C.muted)),
      out(T("  ✓ ", C.success), T("18 repos", C.text), T(" indexed", C.muted)),
      out(T("  ✓ wrote ", C.success), T(".gitmap/output/gitmap.json", C.path)),
      out(T("  ✓ wrote ", C.success), T(".gitmap/output/gitmap.csv", C.path)),
      out(T("  ✓ wrote ", C.success), T(".gitmap/output/clone.ps1", C.path)),
    ],
  },
  // 3) history
  {
    caption: "history — recent activity",
    cwd: "~/repos/gitmap-v27",
    lines: [
      cmd("~/repos/gitmap-v27", T("gitmap ", C.command), T("history ", C.brandGold), T("--limit 5", C.flag)),
      out(T("a3f1c2e ", C.warn), T("feat: fix-repo --strict runs go test on touched packages", C.text)),
      out(T("9b22f01 ", C.warn), T("ci: jq → python3 → awk fallback for json parsing", C.text)),
      out(T("4d5e8aa ", C.warn), T("fix(fixrepo): width-crossing v9→v12 desync guard", C.text)),
      out(T("77b09f3 ", C.warn), T("feat(merge): merge-both / -left / -right with prompt", C.text)),
      out(T("12c0d11 ", C.warn), T("feat(release-alias): as / ra / rap with auto-stash", C.text)),
    ],
  },
  // 4) release
  {
    caption: "release — tag, build, publish",
    cwd: "~/repos/gitmap-v27",
    lines: [
      cmd("~/repos/gitmap-v27", T("gitmap ", C.command), T("release ", C.brandGold), T("--bump minor", C.flag)),
      out(T("→ current: ", C.muted), T("v4.14.0", C.flag), T("  →  next: ", C.muted), T("v4.15.0", C.success)),
      out(T("→ cross-compiling linux/amd64, darwin/arm64, windows/amd64 …", C.muted)),
      out(T("✓ checksums written", C.success), T("  ✓ assets uploaded", C.success)),
      out(T("✔ released ", C.success), T("v4.15.0", C.brandGold)),
    ],
  },
  // 5) ssh
  {
    caption: "ssh — manage keys per repo",
    cwd: "~/repos",
    lines: [
      cmd("~/repos", T("gitmap ", C.command), T("ssh ", C.brandGold), T("--add ", C.flag), T("~/.ssh/id_ed25519_work", C.path)),
      out(T("→ key fingerprint: ", C.muted), T("SHA256:Yk2…", C.flag)),
      out(T("→ added to ssh-agent ✓", C.success)),
      out(T("→ scoped to repos under ", C.muted), T("~/repos/work/*", C.path)),
    ],
  },
  // 6) make-public / make-private
  {
    caption: "Toggle visibility on GitHub",
    cwd: "~/repos/lovable-cloud",
    lines: [
      cmd("~/repos/lovable-cloud", T("gitmap ", C.command), T("make-private", C.brandGold)),
      out(T("→ ", C.muted), T("alimtvnetwork/lovable-cloud", C.path), T(" is now ", C.muted), T("private", C.warn), T(" ✓", C.success)),
      blank(),
      cmd("~/repos/lovable-cloud", T("gitmap ", C.command), T("make-public", C.brandGold)),
      out(T("→ ", C.muted), T("alimtvnetwork/lovable-cloud", C.path), T(" is now ", C.muted), T("public", C.success), T(" ✓", C.success)),
    ],
  },
  // 7) merge family
  {
    caption: "merge-left · merge-right · merge-both",
    cwd: "~/repos/design-system",
    lines: [
      cmd("~/repos/design-system", T("gitmap ", C.command), T("merge-left", C.brandGold)),
      out(T("← prefer LEFT side (", C.muted), T("local", C.flag), T(") on conflicts", C.muted)),
      out(T("✓ 3 files merged · 0 conflicts", C.success)),
      blank(),
      cmd("~/repos/design-system", T("gitmap ", C.command), T("merge-right", C.brandGold)),
      out(T("→ prefer RIGHT side (", C.muted), T("remote", C.flag), T(") on conflicts", C.muted)),
      out(T("✓ 2 files merged · 0 conflicts", C.success)),
      blank(),
      cmd("~/repos/design-system", T("gitmap ", C.command), T("merge-both", C.brandGold)),
      out(T("⇄ interactive: ", C.muted), T("[L]eft  [R]ight  [B]oth  [S]kip", C.flag)),
      out(T("✓ resolved 4/4 hunks ✓ committed ✓ pushed", C.success)),
    ],
  },
  // 8) interactive
  {
    caption: "interactive — full TUI dashboard",
    cwd: "~/repos",
    lines: [
      cmd("~/repos", T("gitmap ", C.command), T("interactive", C.brandGold)),
      out(T("launching TUI…", C.muted)),
      out(T("9 views: Repos · Actions · Groups · Status · Releases · Logs …", C.flag)),
      out(T("press ", C.muted), T("?", C.brandGold), T(" for help, ", C.muted), T("q", C.brandGold), T(" to quit", C.muted)),
    ],
  },
];