import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { Play, Trash2, Save, Database } from "lucide-react";
import { BookmarkTerminalPreview } from "@/components/bookmarks/BookmarkTerminalPreview";
import { BookmarkSchemaTable } from "@/components/bookmarks/BookmarkSchemaTable";
import { BookmarkFileLayoutTable } from "@/components/bookmarks/BookmarkFileLayoutTable";

const features = [
  { icon: Save, title: "Save Commands", desc: "Store frequently-used command+flag combinations under a memorable name." },
  { icon: Play, title: "Replay via Dispatch", desc: "Replayed bookmarks go through standard dispatch with audit hooks." },
  { icon: Database, title: "SQLite Persistence", desc: "Bookmarks are stored in a dedicated Bookmarks table with unique constraints." },
  { icon: Trash2, title: "Safe Deletion", desc: "Delete bookmarks by name. Unique name constraint prevents overwriting." },
];

const BookmarksPage = () => (
  <DocsLayout>
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 space-y-8">
      <div>
        <h1 className="text-3xl font-heading font-bold mb-2 docs-h1 tracking-[-0.02em]">Bookmarks</h1>
        <p className="text-muted-foreground">Save and replay frequently-used CLI command+flag combinations by name.</p>
      </div>

      <section>
        <h2 className="text-xl font-heading font-semibold mb-2 tracking-[-0.02em]">Live Preview</h2>
        <BookmarkTerminalPreview />
      </section>

      <section>
        <h2 className="text-xl font-heading font-semibold mb-4 tracking-[-0.02em]">Features</h2>
        <div className="grid md:grid-cols-2 gap-4">
          {features.map((f) => (
            <div key={f.title} className="rounded-lg border border-border bg-card p-4">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-mono font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="space-y-3">
        <h2 className="text-xl font-heading font-semibold tracking-[-0.02em]">Subcommands</h2>
        <CodeBlock code="gitmap bookmark save ssh-scan scan --mode ssh" title="Save a command bookmark" />
        <CodeBlock code="gitmap bookmark list" title="List all saved bookmarks" />
        <CodeBlock code="gitmap bookmark run ssh-scan" title="Replay a saved bookmark" />
        <CodeBlock code="gitmap bookmark delete ssh-scan" title="Remove a bookmark" />
      </section>

      <section>
        <h2 className="text-xl font-heading font-semibold mb-3 tracking-[-0.02em]">Table Schema</h2>
        <BookmarkSchemaTable />
      </section>

      <section>
        <h2 className="text-xl font-heading font-semibold mb-3 tracking-[-0.02em]">File Layout</h2>
        <BookmarkFileLayoutTable />
      </section>
    </div>
  </DocsLayout>
);

export default BookmarksPage;
