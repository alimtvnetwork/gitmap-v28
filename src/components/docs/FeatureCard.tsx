import { type LucideIcon } from "lucide-react";

interface FeatureCardProps {
  icon: LucideIcon;
  title: string;
  description: string;
}

const FeatureCard = ({ icon: Icon, title, description }: FeatureCardProps) => {
  return (
    <div className="group rounded-md border border-border bg-card p-6 shadow-sm transition-all duration-300 hover:border-primary/40 hover:bg-card">
      <div className="mb-4 flex h-10 w-10 items-center justify-center rounded-sm border border-border bg-secondary transition-colors duration-300 group-hover:border-primary/40 group-hover:bg-accent">
        <Icon className="h-5 w-5 text-primary transition-colors duration-300 group-hover:text-primary" />
      </div>
      <h3 className="font-heading font-semibold text-foreground mb-2">{title}</h3>
      <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
    </div>
  );
};

export default FeatureCard;
