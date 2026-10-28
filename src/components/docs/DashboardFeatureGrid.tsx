import { LucideIcon } from "lucide-react";

export interface DashboardFeature {
  icon: LucideIcon;
  title: string;
  desc: string;
}

interface DashboardFeatureGridProps {
  features: DashboardFeature[];
}

export const DashboardFeatureGrid = ({ features }: DashboardFeatureGridProps) => {
  return (
    <div className="grid md:grid-cols-2 gap-4 mb-8">
      {features.map((f) => (
        <div
          key={f.title}
          className="rounded-lg border border-border bg-card p-4 shadow-sm hover:border-primary/40 transform-gpu hover:-translate-y-0.5 transition-all"
        >
          <f.icon className="h-5 w-5 text-primary mb-2" />
          <h3 className="font-heading font-semibold tracking-[-0.01em] text-sm mb-1 text-foreground">
            {f.title}
          </h3>
          <p className="text-xs text-muted-foreground leading-relaxed">{f.desc}</p>
        </div>
      ))}
    </div>
  );
};

export default DashboardFeatureGrid;
