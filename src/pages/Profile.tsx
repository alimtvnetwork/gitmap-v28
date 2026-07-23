import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { UserCircle, Database, Shield, ArrowRightLeft, Trash2 } from "lucide-react";

const TerminalPreview = () => (
  <div className="rounded-lg border border-border overflow-hidden my-6">
    <div className="bg-terminal px-4 py-2 flex items-center gap-2 border-b border-border">
      <div className="flex gap-1.5">
        <span className="w-3 h-3 rounded-full bg-red-500/80" />
        <span className="w-3 h-3 rounded-full bg-yellow-500/80" />
        <span className="w-3 h-3 rounded-full bg-green-500/80" />
      </div>
      <span className="text-xs font-mono text-muted-foreground ml-2">gitmap profile list</span>
    </div>
    <div className="bg-terminal p-4 font-mono text-sm leading-relaxed">
      <div className="text-primary font-bold text-xs mb-1">
        {"  "}PROFILE{"              "}REPOS{"  "}STATUS
      </div>
      <div className="text-foreground">{"  "}default{"              "}42{"     "}</div>
      <div className="text-foreground">{"  "}work{"                 "}18{"     "}<span className="text-green-400">✓ active</span></div>
      <div className="text-foreground">{"  "}personal{"             "}7</div>
      <div className="text-muted-foreground mt-2 text-xs">3 profiles</div>
    </div>
  </div>
);

const subcommands = [
  {
    name: "create",
    usage: "gitmap profile create <name>",
    description: "Create a new profile with its own empty database.",
    example: "gitmap profile create work\ngitmap pf create personal",
    output: "✓ Profile 'work' created",
  },
  {
    name: "list",
    usage: "gitmap profile list",
    description: "Show all profiles with repo counts and active marker.",
    example: "gitmap pf list",
    output: null,
  },
  {
    name: "switch",
    usage: "gitmap profile switch <name>",
    description: "Switch to a different profile. All subsequent commands use the new profile's database.",
    example: "gitmap profile switch work\ngitmap pf switch personal",
    output: "✓ Switched to profile 'work'",
  },
  {
    name: "show",
    usage: "gitmap profile show",
    description: "Display the currently active profile name and metadata.",
    example: "gitmap pf show",
    output: "Active profile: work\nRepos: 18\nCreated: 2025-03-01",
  },
  {
    name: "delete",
    usage: "gitmap profile delete <name>",
    description: "Delete a profile and its database file. Cannot delete 'default' or the active profile.",
    example: "gitmap profile delete personal",
    output: "✓ Profile 'personal' deleted",
  },
];

const features = [
  { icon: Database, title: "Isolated Databases", desc: "Each profile uses a separate SQLite file (e.g., gitmap-work.db)." },
  { icon: ArrowRightLeft, title: "Instant Switching", desc: "Switch profiles and all commands immediately use the new database." },
  { icon: Shield, title: "Safe Defaults", desc: "The 'default' profile always exists and cannot be deleted." },
  { icon: Trash2, title: "Clean Deletion", desc: "Deleting a profile removes its database file completely." },
];

const ProfilePage = () => (
  <DocsLayout>
    <div className="max-w-4xl">
      <div className="flex items-center gap-3 mb-2">
        <UserCircle className="h-8 w-8 text-primary" />
        <div>
          <h1 className="text-3xl font-heading font-bold text-foreground docs-h1">Profile</h1>
          <p className="text-muted-foreground font-mono text-sm">gitmap profile (pf)</p>
        </div>
      </div>
      <p className="text-muted-foreground mb-8 text-lg">
        Manage isolated database environments for different contexts — work, personal, client projects.
      </p>

      <TerminalPreview />

      {/* Features */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Features</h2>
        <div className="grid sm:grid-cols-2 gap-4">
          {features.map((f) => (
            <div key={f.title} className="border border-border rounded-lg p-4 bg-card">
              <div className="flex items-center gap-2 mb-2">
                <f.icon className="h-5 w-5 text-primary" />
                <span className="font-mono font-semibold text-foreground">{f.title}</span>
              </div>
              <p className="text-sm text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Switch Behavior */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Switch Behavior</h2>
        <div className="border border-border rounded-lg p-4 bg-card">
          <div className="flex items-center gap-3 mb-3 flex-wrap">
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">1. Update profiles.json</span>
            <span className="text-muted-foreground">→</span>
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">2. Resolve DB path</span>
            <span className="text-muted-foreground">→</span>
            <span className="px-2 py-1 bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary rounded text-xs font-mono transition-colors duration-300 hover:border-primary/40 hover:shadow-sm hover:shadow-primary/10">3. Open new DB</span>
          </div>
          <p className="text-sm text-muted-foreground mb-3">
            Switching sets <code className="text-primary">active</code> in <code className="text-primary">profiles.json</code>,
            then resolves the database file name. The <code className="text-primary">default</code> profile maps to{" "}
            <code className="text-primary">gitmap.db</code> for backward compatibility; all others map to{" "}
            <code className="text-primary">gitmap-&lt;name&gt;.db</code>.
          </p>
          <CodeBlock code={`# profiles.json after switch\n{\n  "active": "work",\n  "profiles": ["default", "work", "personal"]\n}`} />
        </div>
      </section>

      {/* Subcommands */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Subcommands</h2>
        <div className="space-y-6">
          {subcommands.map((sub) => (
            <div key={sub.name} className="border border-border rounded-lg p-4 bg-card">
              <h3 className="font-mono font-bold text-foreground mb-1">{sub.name}</h3>
              <p className="text-sm text-muted-foreground mb-3">{sub.description}</p>
              <CodeBlock code={sub.usage} />
              <div className="mt-3">
                <p className="text-xs font-mono text-muted-foreground uppercase tracking-wider mb-1">Example</p>
                <CodeBlock code={sub.example} />
              </div>
              {sub.output && (
                <div className="mt-3">
                  <p className="text-xs font-mono text-muted-foreground uppercase tracking-wider mb-1">Output</p>
                  <div className="bg-terminal rounded-md p-3 font-mono text-sm text-green-400 whitespace-pre-line">
                    {sub.output}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>

      {/* Database Mapping */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Database Mapping</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm font-mono">
            <thead>
              <tr className="border-b border-border text-muted-foreground">
                <th className="text-left py-2 pr-4">Profile</th>
                <th className="text-left py-2">Database File</th>
              </tr>
            </thead>
            <tbody>
              {[
                ["default", "gitmap.db"],
                ["work", "gitmap-work.db"],
                ["personal", "gitmap-personal.db"],
              ].map(([profile, file]) => (
                <tr key={profile} className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">{profile}</td>
                  <td className="py-2 text-foreground">{file}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* Constraints */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">Constraints</h2>
        <ul className="space-y-2 text-sm text-muted-foreground list-disc list-inside">
          <li>Default profile always exists and cannot be deleted.</li>
          <li>Cannot delete the currently active profile.</li>
          <li>Profile names must be unique.</li>
          <li>Backward compatible: no profiles.json = default profile.</li>
        </ul>
      </section>

      {/* File Layout */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">File Layout</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm font-mono">
            <thead>
              <tr className="border-b border-border text-muted-foreground">
                <th className="text-left py-2 pr-4">File</th>
                <th className="text-left py-2">Purpose</th>
              </tr>
            </thead>
            <tbody>
              {[
                ["constants/constants_profile.go", "Command names, messages"],
                ["model/profile.go", "ProfileConfig struct"],
                ["store/profile.go", "Profile config read/write, DB resolution"],
                ["cmd/profile.go", "Profile command routing"],
                ["cmd/profileops.go", "Create, list, switch, delete, show handlers"],
                ["cmd/profileutil.go", "Shared profile helper functions"],
              ].map(([file, purpose]) => (
                <tr key={file} className="border-b border-border">
                  <td className="py-2 pr-4 text-primary">{file}</td>
                  <td className="py-2 text-foreground">{purpose}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* See Also */}
      <section className="mb-10">
        <h2 className="text-xl font-heading font-bold text-foreground mb-4 docs-h2">See Also</h2>
        <ul className="space-y-1 text-sm font-mono">
          <li><a href="/commands" className="text-primary hover:underline">diff-profiles</a> — Compare repos across profiles</li>
          <li><a href="/export" className="text-primary hover:underline">export</a> — Export current profile data</li>
          <li><a href="/import" className="text-primary hover:underline">import</a> — Import data into a profile</li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default ProfilePage;
