import { useState, useMemo, useRef, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import DocsLayout from "@/components/docs/DocsLayout";
import { motion } from "framer-motion";
import { Search } from "lucide-react";
import { sections } from "@/components/spec/specData";
import SpecSearchBar from "@/components/spec/SpecSearchBar";
import SpecSectionCard from "@/components/spec/SpecSectionCard";

const container = { hidden: {}, show: { transition: { staggerChildren: 0.08 } } };
const item = { hidden: { opacity: 0, y: 12 }, show: { opacity: 1, y: 0 } };

const SpecIndexPage = () => {
  const location = useLocation();
  const [query, setQuery] = useState("");
  const [collapsed, setCollapsed] = useState<Record<string, boolean>>({});
  const inputRef = useRef<HTMLInputElement>(null);

  const toggleSection = (folder: string) => {
    setCollapsed((prev) => ({ ...prev, [folder]: !prev[folder] }));
  };

  const allCollapsed = sections.every((s) => collapsed[s.folder]);
  const toggleAll = () => {
    const next: Record<string, boolean> = {};
    sections.forEach((s) => (next[s.folder] = !allCollapsed));
    setCollapsed(next);
  };

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "/" && !["INPUT", "TEXTAREA", "SELECT"].includes((e.target as HTMLElement).tagName)) {
        e.preventDefault();
        inputRef.current?.focus();
      }
      if (e.key === "Escape" && document.activeElement === inputRef.current) {
        inputRef.current?.blur();
        setQuery("");
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  // Scroll to hash anchor and auto-expand that section
  useEffect(() => {
    const hash = location.hash.replace("#", "");
    if (hash) {
      setCollapsed((prev) => ({ ...prev, [hash]: false }));
      requestAnimationFrame(() => {
        document.getElementById(hash)?.scrollIntoView({ behavior: "smooth", block: "start" });
      });
    }
  }, [location.hash]);

  const filtered = useMemo(() => {
    const q = query.toLowerCase().trim();
    if (!q) return sections;
    return sections
      .map((section) => ({
        ...section,
        entries: section.entries.filter(
          (e) =>
            e.title.toLowerCase().includes(q) ||
            e.id.toLowerCase().includes(q) ||
            section.folder.toLowerCase().includes(q) ||
            section.title.toLowerCase().includes(q)
        ),
      }))
      .filter((s) => s.entries.length > 0);
  }, [query]);

  const totalResults = filtered.reduce((sum, s) => sum + s.entries.length, 0);

  return (
    <DocsLayout>
      <motion.div initial={{ opacity: 0, y: -10 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.3 }}>
        <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Spec Index</h1>
        <p className="text-muted-foreground mb-2">
          Complete table of contents for all specification documents, issue post-mortems, design guidelines, and the generic CLI blueprint.
        </p>

        <SpecSearchBar ref={inputRef} query={query} onChange={setQuery} />

        {!query && (
          <div className="flex flex-wrap gap-2 mb-4">
            {sections.map((s) => (
              <a
                key={s.folder}
                href={`#${s.folder}`}
                onClick={(e) => {
                  e.preventDefault();
                  setCollapsed((prev) => ({ ...prev, [s.folder]: false }));
                  requestAnimationFrame(() => {
                    document.getElementById(s.folder)?.scrollIntoView({ behavior: "smooth", block: "start" });
                  });
                }}
                className="text-xs font-mono px-2.5 py-1 rounded border border-border bg-muted/30 text-muted-foreground hover:text-primary hover:border-primary/50 transition-colors"
              >
                {s.title}
              </a>
            ))}
          </div>
        )}

        <div className="flex items-center justify-between mb-6">
          <p className="text-xs text-muted-foreground/60 font-mono">
            {query
              ? `${totalResults} result${totalResults !== 1 ? "s" : ""} matching "${query}"`
              : `${totalResults} documents across ${sections.length} sections`}
          </p>
          {!query && (
            <button onClick={toggleAll} className="text-xs font-mono text-muted-foreground hover:text-foreground transition-colors">
              {allCollapsed ? "Expand all" : "Collapse all"}
            </button>
          )}
        </div>
      </motion.div>

      <motion.div variants={container} initial="hidden" animate="show" className="space-y-8" key={query}>
        {filtered.length === 0 && (
          <div className="text-center py-12 text-muted-foreground">
            <Search className="h-8 w-8 mx-auto mb-3 opacity-40" />
            <p className="text-sm font-mono">No specs matching "{query}"</p>
          </div>
        )}
        {filtered.map((section) => (
          <motion.div key={section.folder} variants={item}>
            <SpecSectionCard
              section={section}
              isCollapsed={!query && !!collapsed[section.folder]}
              onToggle={() => toggleSection(section.folder)}
            />
          </motion.div>
        ))}
      </motion.div>

      {/* See Also */}
      <div className="mt-10 pt-6 border-t border-border">
        <h3 className="text-sm font-mono font-semibold text-muted-foreground mb-3">See Also</h3>
        <div className="flex flex-wrap gap-2">
          {[
            { label: "Architecture", to: "/architecture" },
            { label: "Commands", to: "/commands" },
            { label: "Generic CLI", to: "/generic-cli" },
            { label: "Changelog", to: "/changelog" },
          ].map((link) => (
            <Link
              key={link.to}
              to={link.to}
              className="text-xs font-mono px-3 py-1.5 rounded border border-border bg-muted/30 text-muted-foreground hover:text-primary hover:border-primary/50 transition-colors"
            >
              {link.label}
            </Link>
          ))}
        </div>
      </div>
    </DocsLayout>
  );
};

export default SpecIndexPage;
