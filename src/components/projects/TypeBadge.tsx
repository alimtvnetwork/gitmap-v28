import { Code2, Braces, Cpu, Hash } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { ProjectType, type ProjectTypeConfig } from "./types";

export const ProjectTypes: Record<ProjectType, ProjectTypeConfig> = {
  [ProjectType.Go]: { label: "Go", color: "bg-muted/60 text-foreground/80 border-border/60 hover:bg-muted/80", icon: Code2 },
  [ProjectType.Node]: { label: "Node.js", color: "bg-muted/60 text-foreground/80 border-border/60 hover:bg-muted/80", icon: Braces },
  [ProjectType.React]: { label: "React", color: "bg-muted/60 text-foreground/80 border-border/60 hover:bg-muted/80", icon: Braces },
  [ProjectType.Cpp]: { label: "C++", color: "bg-muted/60 text-foreground/80 border-border/60 hover:bg-muted/80", icon: Cpu },
  [ProjectType.Csharp]: { label: "C#", color: "bg-muted/60 text-foreground/80 border-border/60 hover:bg-muted/80", icon: Hash },
};

export const TypeBadge = ({ type }: { type: ProjectType }) => {
  const config = ProjectTypes[type];
  const Icon = config.icon;

  return (
    <Badge variant="outline" className={`${config.color} font-mono text-xs gap-1 border`}>
      <Icon className="h-3 w-3" />
      {config.label}
    </Badge>
  );
};
