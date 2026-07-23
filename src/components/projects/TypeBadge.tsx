import { Code2, Braces, Cpu, Hash } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import type { ProjectType, ProjectTypeConfig } from "./types";

export const ProjectTypes: Record<ProjectType, ProjectTypeConfig> = {
  go: { label: "Go", color: "bg-cyan-500/15 text-cyan-700 dark:text-cyan-400 border-cyan-500/30", icon: Code2 },
  node: { label: "Node.js", color: "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/30", icon: Braces },
  react: { label: "React", color: "bg-sky-500/15 text-sky-700 dark:text-sky-400 border-sky-500/30", icon: Braces },
  cpp: { label: "C++", color: "bg-violet-500/15 text-violet-700 dark:text-violet-400 border-violet-500/30", icon: Cpu },
  csharp: { label: "C#", color: "bg-purple-500/15 text-purple-700 dark:text-purple-400 border-purple-500/30", icon: Hash },
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