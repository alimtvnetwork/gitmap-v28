import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import CloneNextCommandBuilder from "@/components/docs/CloneNextCommandBuilder";

const flags = [
  { flag: "--force, -f", default: "false", desc: "Force flatten when cwd IS the target folder (chdir to parent first; refuses versioned-folder fallback)" },
  { flag: "--delete", default: "false", desc: "Auto-remove current versioned folder after clone" },
  { flag: "--keep", default: "false", desc: "Keep current folder without prompting" },
  { flag: "--no-desktop", default: "false", desc: "Skip GitHub Desktop registration" },
  { flag: "--ssh-key, -K <name>", default: "(none)", desc: "Use a named SSH key for the clone" },
  { flag: "--verbose", default: "false", desc: "Write detailed debug log" },
  { flag: "--create-remote", default: "false", desc: "Create target GitHub repo if missing (requires GITHUB_TOKEN)" },
  { flag: "--csv <path>", default: "(none)", desc: "Batch mode: read repo paths from CSV (one path per row, header optional)" },
  { flag: "--all", default: "false", desc: "Batch mode: walk cwd and run cn on every git repo one level deep" },
];

const versionArgs = [
  { arg: "v++", desc: "Increment current version by one (e.g. v11 → v12)" },
  { arg: "v+1", desc: "Alias for v++" },
  { arg: "vN", desc: "Jump to a specific version number (e.g. v15)" },
];

const CloneNextCommandPage = () => {
  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">gitmap clone-next</h1>
      <p className="text-muted-foreground mb-2">
        Clone the next (or a specific) versioned iteration of the current repository into the parent
        directory, using the base name (no version suffix) as the local folder. Reads the remote from
        the cwd's git origin and walks <code className="docs-inline-code">-vN</code> suffixes.
      </p>
      <p className="text-sm text-muted-foreground mb-8">
        Alias: <code className="docs-inline-code">cn</code>
      </p>

      <section className="space-y-8">
        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Usage</h2>
          <CodeBlock code={`gitmap clone-next <v++|v+1|vN> [flags]`} title="Syntax" />
          <p className="text-sm text-muted-foreground mt-2">
            Must be run inside a Git repository with a remote origin configured. The remote URL is
            rewritten by swapping the <code className="docs-inline-code">-vN</code> suffix.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Version argument</h2>
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-sm">
              <thead className="bg-muted/40">
                <tr>
                  <th className="text-left px-4 py-2 font-mono">Argument</th>
                  <th className="text-left px-4 py-2 font-mono">Behavior</th>
                </tr>
              </thead>
              <tbody>
                {versionArgs.map((v) => (
                  <tr key={v.arg} className="border-t border-border">
                    <td className="px-4 py-2"><code className="docs-inline-code">{v.arg}</code></td>
                    <td className="px-4 py-2 text-muted-foreground">{v.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Flags</h2>
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-sm">
              <thead className="bg-muted/40">
                <tr>
                  <th className="text-left px-4 py-2 font-mono">Flag</th>
                  <th className="text-left px-4 py-2 font-mono">Default</th>
                  <th className="text-left px-4 py-2 font-mono">Description</th>
                </tr>
              </thead>
              <tbody>
                {flags.map((f) => (
                  <tr key={f.flag} className="border-t border-border">
                    <td className="px-4 py-2"><code className="docs-inline-code">{f.flag}</code></td>
                    <td className="px-4 py-2 text-muted-foreground">{f.default}</td>
                    <td className="px-4 py-2 text-muted-foreground">{f.desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="text-sm text-muted-foreground mt-3">
            Flags placed after the positional version (e.g. <code className="docs-inline-code">cn v+1 -f</code>)
            are reordered before parsing — both prefix and suffix forms work.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Flatten behavior</h2>
          <p className="text-sm text-muted-foreground mb-3">
            By default, clone-next clones into the <strong>base name</strong> folder (without version suffix).
            Running <code className="docs-inline-code">gitmap cn v++</code> inside <code className="docs-inline-code">macro-ahk-v11</code> will:
          </p>
          <ol className="list-decimal list-inside text-sm text-muted-foreground space-y-1 mb-3">
            <li>Clone <code className="docs-inline-code">macro-ahk-v12</code> into <code className="docs-inline-code">macro-ahk/</code></li>
            <li>If <code className="docs-inline-code">macro-ahk/</code> already exists, remove it first</li>
            <li>Keep the remote URL pointing at <code className="docs-inline-code">macro-ahk-v12</code> on GitHub</li>
            <li>Record the version transition (v11 → v12) in the database</li>
          </ol>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Edge cases</h2>

          <h3 className="text-base font-heading font-semibold mb-2 mt-2 docs-h3">Target version folder already exists</h3>
          <p className="text-sm text-muted-foreground mb-2">
            When the flattened target folder (e.g. <code className="docs-inline-code">macro-ahk/</code>)
            already exists, clone-next removes it <strong>without prompting</strong> before cloning. There is
            no <code className="docs-inline-code">--yes</code> confirmation — flatten is destructive by design.
          </p>
          <ul className="list-disc list-inside text-sm text-muted-foreground space-y-1 mb-3">
            <li>
              <strong>cwd is not the target folder:</strong>{" "}
              <code className="docs-inline-code">os.RemoveAll(target)</code> runs immediately, then the clone
              proceeds into the freshly emptied path.
            </li>
            <li>
              <strong>cwd IS the target folder (Windows lock):</strong> removal fails because the shell holds a
              file handle on the cwd. Without <code className="docs-inline-code">-f</code>, gitmap falls back
              to a versioned folder name (e.g. <code className="docs-inline-code">macro-ahk-v22/</code>) and
              prints <code className="docs-inline-code">MsgFlattenFallback</code> with a hint to retry with{" "}
              <code className="docs-inline-code">-f</code>.
            </li>
            <li>
              <strong>cwd IS the target folder + <code className="docs-inline-code">-f</code>:</strong> gitmap-v27{" "}
              <code className="docs-inline-code">chdir</code>s to the parent first to release the lock, removes
              the folder, clones, then <code className="docs-inline-code">chdir</code>s back into the new
              flattened folder. If removal still fails (another process holds the lock), gitmap aborts with{" "}
              <code className="docs-inline-code">ErrCloneNextForceFailed</code> and exits 1 — it
              <strong> never silently degrades</strong> to a versioned folder when{" "}
              <code className="docs-inline-code">-f</code> is set.
            </li>
            <li>
              <strong>Target folder is locked by another process:</strong> see the lock-detection example below
              — gitmap lists the offending PIDs (e.g. <code className="docs-inline-code">Code.exe</code>,{" "}
              <code className="docs-inline-code">explorer.exe</code>) and offers to terminate them.
            </li>
          </ul>

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Origin has no <code className="docs-inline-code">-vN</code> suffix</h3>
          <p className="text-sm text-muted-foreground mb-2">
            clone-next derives the next remote URL by parsing a trailing{" "}
            <code className="docs-inline-code">-vN</code> suffix on the current repo name. The parser
            (<code className="docs-inline-code">clonenext.ParseRepoName</code>) handles unsuffixed repos as
            <strong> implicit v1</strong>:
          </p>
          <ul className="list-disc list-inside text-sm text-muted-foreground space-y-1 mb-3">
            <li>
              <code className="docs-inline-code">macro-ahk</code> → base{" "}
              <code className="docs-inline-code">macro-ahk</code>, current version{" "}
              <code className="docs-inline-code">1</code>,{" "}
              <code className="docs-inline-code">HasVersion=false</code>.
            </li>
            <li>
              <code className="docs-inline-code">cn v++</code> targets{" "}
              <code className="docs-inline-code">macro-ahk-v2</code> (current + 1 = 2).
            </li>
            <li>
              <code className="docs-inline-code">cn v15</code> targets{" "}
              <code className="docs-inline-code">macro-ahk-v15</code> directly.
            </li>
            <li>
              Remote URL rewrite uses{" "}
              <code className="docs-inline-code">strings.Replace(url, currentRepo, targetRepo, 1)</code>.
              When the origin has no suffix, gitmap looks for the bare repo name in the URL and produces a
              suffixed URL — so the target repo <strong>must already exist on GitHub</strong> (or pass{" "}
              <code className="docs-inline-code">--create-remote</code> with a{" "}
              <code className="docs-inline-code">GITHUB_TOKEN</code> to provision it).
            </li>
            <li>
              Local folder is still flattened to the base name, so cloning{" "}
              <code className="docs-inline-code">macro-ahk-v2</code> from inside{" "}
              <code className="docs-inline-code">macro-ahk/</code> places it back into{" "}
              <code className="docs-inline-code">macro-ahk/</code> (the existing-folder rules above apply).
            </li>
          </ul>
          <p className="text-sm text-muted-foreground">
            If the version argument cannot be parsed (e.g.{" "}
            <code className="docs-inline-code">cn foo</code>,{" "}
            <code className="docs-inline-code">cn v0</code>,{" "}
            <code className="docs-inline-code">cn v-1</code>), gitmap exits with{" "}
            <code className="docs-inline-code">invalid version argument: ... (expected v++, v+1, or vN)</code>{" "}
            and makes no network or filesystem changes.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">URL rewrite — HTTPS and SSH</h2>
          <p className="text-sm text-muted-foreground mb-3">
            clone-next reads <code className="docs-inline-code">git remote get-url origin</code> from the
            current repo and rewrites it with{" "}
            <code className="docs-inline-code">strings.Replace(url, currentRepo, targetRepo, 1)</code>. The
            same rule applies whether the origin is HTTPS or SSH — only the substring matching the current
            repo name is swapped, so credentials, hosts, and paths are preserved untouched.
          </p>

          <h3 className="text-base font-heading font-semibold mb-2 mt-4 docs-h3">HTTPS origin</h3>
          <p className="text-sm text-muted-foreground mb-2">
            Inside <code className="docs-inline-code">D:\repos\macro-ahk-v11\</code> with origin{" "}
            <code className="docs-inline-code">https://github.com/alimtvnetwork/macro-ahk-v11.git</code>:
          </p>
          <CodeBlock code={`gitmap cn v++`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Resolved rewrite"
            code={`current : https://github.com/alimtvnetwork/macro-ahk-v11.git
target  : https://github.com/alimtvnetwork/macro-ahk-v12.git
folder  : macro-ahk/   (flattened, base name)
git cmd : git clone https://github.com/alimtvnetwork/macro-ahk-v12.git macro-ahk`}
          />
          <p className="text-sm text-muted-foreground mt-2">
            Jumping to a specific version uses the same substitution:
          </p>
          <CodeBlock code={`gitmap cn v15`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Resolved rewrite"
            code={`current : https://github.com/alimtvnetwork/macro-ahk-v11.git
target  : https://github.com/alimtvnetwork/macro-ahk-v15.git
folder  : macro-ahk/`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">SSH origin</h3>
          <p className="text-sm text-muted-foreground mb-2">
            Inside <code className="docs-inline-code">D:\repos\macro-ahk-v11\</code> with origin{" "}
            <code className="docs-inline-code">git@github.com:alimtvnetwork/macro-ahk-v11.git</code>:
          </p>
          <CodeBlock code={`gitmap cn v++`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Resolved rewrite"
            code={`current : git@github.com:alimtvnetwork/macro-ahk-v11.git
target  : git@github.com:alimtvnetwork/macro-ahk-v12.git
folder  : macro-ahk/
git cmd : git clone git@github.com:alimtvnetwork/macro-ahk-v12.git macro-ahk`}
          />
          <p className="text-sm text-muted-foreground mt-2">
            Use <code className="docs-inline-code">--ssh-key, -K &lt;name&gt;</code> to route the clone
            through a named key from <code className="docs-inline-code">gitmap ssh list</code> — gitmap sets{" "}
            <code className="docs-inline-code">GIT_SSH_COMMAND="ssh -i &lt;path&gt; -o IdentitiesOnly=yes"</code>{" "}
            for the duration of the clone:
          </p>
          <CodeBlock code={`gitmap cn v++ --ssh-key work`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Resolved rewrite"
            code={`current : git@github.com:alimtvnetwork/macro-ahk-v11.git
target  : git@github.com:alimtvnetwork/macro-ahk-v12.git
folder  : macro-ahk/
ssh key : ~/.ssh/id_work     (resolved from named key "work")
git cmd : GIT_SSH_COMMAND="ssh -i ~/.ssh/id_work -o IdentitiesOnly=yes" \\
          git clone git@github.com:alimtvnetwork/macro-ahk-v12.git macro-ahk`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">SSH alias host (per-key Host entries)</h3>
          <p className="text-sm text-muted-foreground mb-2">
            If <code className="docs-inline-code">~/.ssh/config</code> defines a per-key alias such as{" "}
            <code className="docs-inline-code">Host github.com-work</code>, the rewrite still operates only on
            the repo-name segment — the alias host is preserved:
          </p>
          <CodeBlock
            language="bash"
            title="Resolved rewrite"
            code={`current : git@github.com-work:alimtvnetwork/macro-ahk-v11.git
target  : git@github.com-work:alimtvnetwork/macro-ahk-v12.git
folder  : macro-ahk/`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Unsuffixed origin (implicit v1)</h3>
          <p className="text-sm text-muted-foreground mb-2">
            When the origin has no <code className="docs-inline-code">-vN</code> suffix, gitmap treats the repo
            as v1 and rewrites by appending <code className="docs-inline-code">-vN</code> to the bare name —
            same rule for HTTPS and SSH:
          </p>
          <CodeBlock
            language="bash"
            title="HTTPS"
            code={`current : https://github.com/alimtvnetwork/macro-ahk.git
target  : https://github.com/alimtvnetwork/macro-ahk-v2.git
folder  : macro-ahk/`}
          />
          <CodeBlock
            language="bash"
            title="SSH"
            code={`current : git@github.com:alimtvnetwork/macro-ahk.git
target  : git@github.com:alimtvnetwork/macro-ahk-v2.git
folder  : macro-ahk/`}
          />
          <p className="text-sm text-muted-foreground mt-2">
            The target repo must already exist on GitHub. To auto-provision it, pass{" "}
            <code className="docs-inline-code">--create-remote</code> with{" "}
            <code className="docs-inline-code">GITHUB_TOKEN</code> set in the environment.
          </p>
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Copy commands</h2>
          <p className="text-sm text-muted-foreground mb-3">
            Configure any flag combo and copy the exact <code className="docs-inline-code">gitmap</code>{" "}
            invocation plus the underlying <code className="docs-inline-code">git clone</code> it expands to.
            Useful for previewing what gitmap will run before committing to a destructive flatten, or for
            scripting the same action without gitmap installed.
          </p>
          <CloneNextCommandBuilder />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">Examples</h2>

          <h3 className="text-base font-heading font-semibold mb-2 mt-4 docs-h3">Increment version by one</h3>
          <CodeBlock code={`gitmap cn v++`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`Removing existing macro-ahk for fresh clone...
Cloning macro-ahk-v12 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v12 into macro-ahk
✓ Recorded version transition v11 -> v12
✓ Registered macro-ahk-v12 with GitHub Desktop`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Jump to a specific version with auto-delete</h3>
          <CodeBlock code={`gitmap cn v15 --delete`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`Cloning macro-ahk-v15 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v15 into macro-ahk
✓ Recorded version transition v12 -> v15
✓ Registered macro-ahk-v15 with GitHub Desktop
✓ Removed macro-ahk-v12`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Force-flatten from inside an already-flat folder</h3>
          <p className="text-sm text-muted-foreground mb-2">
            You're working in <code className="docs-inline-code">D:\repos\macro-ahk\</code> (flattened from v21)
            and want to bump to v22 without ending up in <code className="docs-inline-code">macro-ahk-v22/</code>.
          </p>
          <CodeBlock code={`gitmap cn v++ -f`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`  → Stepping out of D:\\repos\\macro-ahk to release the file lock

  ── 1/3  Preparing flatten (macro-ahk → macro-ahk) ──
Removing existing macro-ahk for fresh clone...

  ── 2/3  Cloning macro-ahk-v22 ──
Cloning macro-ahk-v22 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v22 into macro-ahk

  ── 3/3  Finalizing ──
✓ Recorded version transition v21 -> v22
✓ Registered macro-ahk-v22 with GitHub Desktop

  ── ✓ Done — now in macro-ahk ──`}
          />
          <p className="text-sm text-muted-foreground mt-2">
            Without <code className="docs-inline-code">-f</code>, gitmap falls back to creating
            <code className="docs-inline-code">macro-ahk-v22/</code> (the shell holds a file lock on cwd) and
            prints a hint. With <code className="docs-inline-code">-f</code>, the versioned-folder fallback is
            disabled — if removal still fails, gitmap aborts with a clear error.
          </p>

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Use a named SSH key</h3>
          <CodeBlock code={`gitmap cn v++ --ssh-key work`} title="Terminal" />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Batch mode — every repo one level deep</h3>
          <CodeBlock code={`gitmap cn v++ --all`} title="Terminal" />
          <CodeBlock
            language="bash"
            title="Output"
            code={`→ Batch cn over 4 repo(s)
  • macro-ahk: v11 -> v12
  • wp-onboarding: v13 -> v14
  • wp-alim: v7 (no update needed)
  • dashboard-kit: v3 -> v4
✓ Batch complete: 3 ok, 0 failed, 1 skipped
  Report: .gitmap/output/cn-batch-2025-04-21T14-30.csv`}
          />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Batch mode — explicit CSV input</h3>
          <CodeBlock code={`gitmap cn v++ --csv repos.csv`} title="Terminal" />

          <h3 className="text-base font-heading font-semibold mb-2 mt-6 docs-h3">Create the remote repo if it doesn't exist</h3>
          <CodeBlock code={`GITHUB_TOKEN=ghp_... gitmap cn v++ --create-remote`} title="Terminal" />
        </div>

        <hr className="docs-hr" />

        <div>
          <h2 className="text-xl font-heading font-semibold mb-3 docs-h2">See also</h2>
          <ul className="list-disc list-inside text-muted-foreground space-y-1">
            <li><a href="/clone-command" className="text-primary hover:underline">clone</a> — Clone repos from output files or a direct URL</li>
            <li><a href="/scan-command" className="text-primary hover:underline">scan</a> — Generate the output files clone consumes</li>
            <li><a href="/clone-next" className="text-primary hover:underline">clone-next (spec)</a> — Full design + flatten spec</li>
            <li><a href="/scan-clone-flags" className="text-primary hover:underline">scan/clone flags</a> — Cross-command flag reference</li>
          </ul>
        </div>
      </section>
    </DocsLayout>
  );
};

export default CloneNextCommandPage;
