import { Database } from "lucide-react";

export const DatabaseHeader = () => {
  return (
    <div className="flex flex-wrap items-center justify-between gap-4 border-b border-border/40 pb-5">
      <div className="flex items-center gap-3">
        <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
          <Database className="h-5 w-5 text-primary" />
        </div>
        <div>
          <h1 className="text-2xl font-heading font-bold tracking-[-0.02em] text-foreground">
            Database Visualizer
          </h1>
          <p className="text-sm text-muted-foreground">
            Interactive table browser, schema inspector, and binary/JSON payload decoder.
          </p>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
          Connected
        </span>
      </div>
    </div>
  );
};

export default DatabaseHeader;
