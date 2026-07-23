import { useMemo, useState } from "react";
import CodeBlock from "./CodeBlock";

/**
 * CloneNextCommandBuilder
 *
 * Reproduces the exact `git clone` invocation that `gitmap clone-next` issues
 * for any combination of its public flags. The logic mirrors:
 *   - clonenext.ParseRepoName / TargetRepoName / ReplaceRepoInURL
 *   - the --ssh-key flow that sets GIT_SSH_COMMAND
 *   - the flatten-folder default (base name, no -vN suffix)
 *
 * Pure presentational state — no network, no app state.
 */

type Protocol = "https" | "ssh" | "ssh-alias";
type VersionMode = "v++" | "v+1" | "vN";

const PROTOCOL_LABEL: Record<Protocol, string> = {
  https: "HTTPS",
  ssh: "SSH (git@github.com)",
  "ssh-alias": "SSH alias (git@github.com-work)",
};

interface BuilderState {
  protocol: Protocol;
  owner: string;
  baseName: string;
  currentVersion: number;
  hasVersion: boolean;
  versionMode: VersionMode;
  explicitVersion: number;
  flatten: boolean;
  force: boolean;
  sshKeyName: string;
  branch: string;
}

const DEFAULTS: BuilderState = {
  protocol: "https",
  owner: "alimtvnetwork",
  baseName: "macro-ahk",
  currentVersion: 11,
  hasVersion: true,
  versionMode: "v++",
  explicitVersion: 15,
  flatten: true,
  force: false,
  sshKeyName: "",
  branch: "",
};

function buildCurrentRepoName(s: BuilderState): string {
  return s.hasVersion ? `${s.baseName}-v${s.currentVersion}` : s.baseName;
}

function resolveTargetVersion(s: BuilderState): number {
  if (s.versionMode === "vN") return Math.max(1, s.explicitVersion);
  // v++ and v+1 both increment by 1
  return (s.hasVersion ? s.currentVersion : 1) + 1;
}

function buildTargetRepoName(s: BuilderState): string {
  return `${s.baseName}-v${resolveTargetVersion(s)}`;
}

function buildOriginURL(s: BuilderState): string {
  const repo = buildCurrentRepoName(s);
  switch (s.protocol) {
    case "https":
      return `https://github.com/${s.owner}/${repo}.git`;
    case "ssh":
      return `git@github.com:${s.owner}/${repo}.git`;
    case "ssh-alias":
      return `git@github.com-work:${s.owner}/${repo}.git`;
  }
}

function buildTargetURL(s: BuilderState): string {
  const current = buildCurrentRepoName(s);
  const target = buildTargetRepoName(s);
  // Mirrors clonenext.ReplaceRepoInURL: strings.Replace(url, current, target, 1)
  return buildOriginURL(s).replace(current, target);
}

function buildLocalFolder(s: BuilderState): string {
  return s.flatten ? s.baseName : buildTargetRepoName(s);
}

function buildGitmapCommand(s: BuilderState): string {
  const parts = ["gitmap", "clone-next"];
  parts.push(s.versionMode === "vN" ? `v${Math.max(1, s.explicitVersion)}` : s.versionMode);
  if (s.force) parts.push("-f");
  if (s.sshKeyName.trim().length > 0) parts.push(`--ssh-key ${s.sshKeyName.trim()}`);
  if (s.branch.trim().length > 0) parts.push(`--branch ${s.branch.trim()}`);
  if (!s.flatten) parts.push("--no-flatten");
  return parts.join(" ");
}

function buildGitCloneCommand(s: BuilderState): string {
  const url = buildTargetURL(s);
  const folder = buildLocalFolder(s);
  const branchArg = s.branch.trim().length > 0 ? `--branch ${s.branch.trim()} ` : "";
  const base = `git clone ${branchArg}${url} ${folder}`;
  if (s.protocol !== "https" && s.sshKeyName.trim().length > 0) {
    const keyName = s.sshKeyName.trim();
    return [
      `# Routed through named SSH key "${keyName}"`,
      `GIT_SSH_COMMAND="ssh -i ~/.ssh/id_${keyName} -o IdentitiesOnly=yes" \\`,
      `  ${base}`,
    ].join("\n");
  }
  return base;
}

function buildResolvedSummary(s: BuilderState): string {
  const lines = [
    `cwd        : D:\\repos\\${buildCurrentRepoName(s)}`,
    `origin     : ${buildOriginURL(s)}`,
    `current    : ${buildCurrentRepoName(s)} (v${s.hasVersion ? s.currentVersion : 1}${s.hasVersion ? "" : ", implicit"})`,
    `target     : ${buildTargetRepoName(s)} (v${resolveTargetVersion(s)})`,
    `target url : ${buildTargetURL(s)}`,
    `folder     : ${buildLocalFolder(s)}${s.flatten ? "/   (flattened, base name)" : "/   (versioned, no flatten)"}`,
  ];
  if (s.force) {
    lines.push(`force      : -f set — chdir-to-parent if cwd == folder; aborts on lock instead of fallback`);
  }
  if (s.protocol !== "https" && s.sshKeyName.trim().length > 0) {
    lines.push(`ssh key    : ~/.ssh/id_${s.sshKeyName.trim()}`);
  }
  if (s.branch.trim().length > 0) {
    lines.push(`branch     : ${s.branch.trim()}`);
  }
  return lines.join("\n");
}

const CloneNextCommandBuilder = () => {
  const [s, setS] = useState<BuilderState>(DEFAULTS);

  const update = <K extends keyof BuilderState>(key: K, value: BuilderState[K]) => {
    setS((prev) => ({ ...prev, [key]: value }));
  };

  const gitmapCmd = useMemo(() => buildGitmapCommand(s), [s]);
  const gitCmd = useMemo(() => buildGitCloneCommand(s), [s]);
  const resolved = useMemo(() => buildResolvedSummary(s), [s]);

  return (
    <div className="rounded-lg border border-border bg-card/40 p-4 space-y-4">
      <div>
        <h3 className="text-base font-heading font-semibold mb-1 docs-h3">Copy commands</h3>
        <p className="text-xs text-muted-foreground">
          Pick a flag combo. The resolved <code className="docs-inline-code">git clone</code> command updates
          live, mirroring the same origin-rewrite logic gitmap runs internally
          (<code className="docs-inline-code">strings.Replace(url, currentRepo, targetRepo, 1)</code>).
        </p>
      </div>

      <div className="grid gap-3 md:grid-cols-2">
        <Field label="Origin protocol">
          <select
            className="w-full rounded-md bg-background border border-border px-2 py-1 text-sm font-mono"
            value={s.protocol}
            onChange={(e) => update("protocol", e.target.value as Protocol)}
          >
            {(Object.keys(PROTOCOL_LABEL) as Protocol[]).map((p) => (
              <option key={p} value={p}>{PROTOCOL_LABEL[p]}</option>
            ))}
          </select>
        </Field>

        <Field label="Owner / org">
          <input
            type="text"
            className="w-full rounded-md bg-background border border-border px-2 py-1 text-sm font-mono"
            value={s.owner}
            onChange={(e) => update("owner", e.target.value)}
          />
        </Field>

        <Field label="Base repo name">
          <input
            type="text"
            className="w-full rounded-md bg-background border border-border px-2 py-1 text-sm font-mono"
            value={s.baseName}
            onChange={(e) => update("baseName", e.target.value)}
          />
        </Field>

        <Field label="Current version">
          <div className="flex items-center gap-2">
            <label className="flex items-center gap-1 text-xs text-muted-foreground">
              <input
                type="checkbox"
                checked={s.hasVersion}
                onChange={(e) => update("hasVersion", e.target.checked)}
              />
              has -vN suffix
            </label>
            <input
              type="number"
              min={1}
              disabled={!s.hasVersion}
              className="w-20 rounded-md bg-background border border-border px-2 py-1 text-sm font-mono disabled:cursor-not-allowed disabled:opacity-60 disabled:bg-muted disabled:text-muted-foreground disabled:border-border"
              value={s.currentVersion}
              onChange={(e) => update("currentVersion", Math.max(1, Number(e.target.value) || 1))}
            />
          </div>
        </Field>

        <Field label="Version argument">
          <div className="flex items-center gap-2">
            <select
              className="rounded-md bg-background border border-border px-2 py-1 text-sm font-mono"
              value={s.versionMode}
              onChange={(e) => update("versionMode", e.target.value as VersionMode)}
            >
              <option value="v++">v++</option>
              <option value="v+1">v+1</option>
              <option value="vN">vN (specific)</option>
            </select>
            {s.versionMode === "vN" && (
              <input
                type="number"
                min={1}
                className="w-20 rounded-md bg-background border border-border px-2 py-1 text-sm font-mono"
                value={s.explicitVersion}
                onChange={(e) => update("explicitVersion", Math.max(1, Number(e.target.value) || 1))}
              />
            )}
          </div>
        </Field>

        <Field label="Branch (optional)">
          <input
            type="text"
            placeholder="main"
            className="w-full rounded-md bg-background border border-border px-2 py-1 text-sm font-mono"
            value={s.branch}
            onChange={(e) => update("branch", e.target.value)}
          />
        </Field>

        <Field label="Named SSH key (optional)">
          <input
            type="text"
            placeholder="work"
            disabled={s.protocol === "https"}
            className="w-full rounded-md bg-background border border-border px-2 py-1 text-sm font-mono disabled:cursor-not-allowed disabled:opacity-60 disabled:bg-muted disabled:text-muted-foreground disabled:border-border"
            value={s.sshKeyName}
            onChange={(e) => update("sshKeyName", e.target.value)}
          />
        </Field>

        <Field label="Behavior flags">
          <div className="flex flex-wrap gap-3 text-xs">
            <label className="flex items-center gap-1">
              <input
                type="checkbox"
                checked={s.flatten}
                onChange={(e) => update("flatten", e.target.checked)}
              />
              flatten to base name
            </label>
            <label className="flex items-center gap-1">
              <input
                type="checkbox"
                checked={s.force}
                onChange={(e) => update("force", e.target.checked)}
              />
              -f / --force
            </label>
          </div>
        </Field>
      </div>

      <div className="space-y-2">
        <p className="text-xs font-mono text-muted-foreground">gitmap invocation</p>
        <CodeBlock language="bash" title="Terminal" code={gitmapCmd} />

        <p className="text-xs font-mono text-muted-foreground mt-3">Resolved git command</p>
        <CodeBlock language="bash" title="Underlying git" code={gitCmd} />

        <p className="text-xs font-mono text-muted-foreground mt-3">Rewrite trace</p>
        <CodeBlock language="bash" title="What gitmap computes" code={resolved} />
      </div>
    </div>
  );
};

interface FieldProps {
  label: string;
  children: React.ReactNode;
}

const Field = ({ label, children }: FieldProps) => (
  <div className="space-y-1">
    <label className="block text-xs font-mono text-muted-foreground">{label}</label>
    {children}
  </div>
);

export default CloneNextCommandBuilder;
