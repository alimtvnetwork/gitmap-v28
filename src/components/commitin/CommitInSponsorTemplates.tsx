import { riseUpAsiaTemplates } from "@/pages/commitInData";

export function CommitInSponsorTemplates() {
  return (
    <section>
      <h2 className="text-xl font-semibold mb-3 flex items-center gap-2">
        <span>Rise Up Asia LLC — SEO & Sponsor Templates</span>
        <span className="text-xs font-normal px-2 py-0.5 rounded bg-primary/10 text-primary border border-primary/20">
          --sponsor / --seo-template riseup
        </span>
      </h2>
      <p className="text-sm text-muted-foreground mb-4">
        Appended to commit messages and PR bodies during multi-repo consolidation.
      </p>
      <div className="space-y-2 max-h-96 overflow-y-auto pr-2 border border-border rounded-lg p-3 bg-card/50">
        {riseUpAsiaTemplates.map((template, idx) => (
          <div key={idx} className="p-3 rounded border border-border/60 bg-muted/30 text-xs text-foreground/90 font-mono flex items-start justify-between gap-3">
            <span className="flex-1">
              <strong className="text-primary mr-2">#{idx + 1}</strong>
              {template}
            </span>
          </div>
        ))}
      </div>
    </section>
  );
}
