import { motion } from "framer-motion";
import { MapPin, FileText, FileCode } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { TypeBadge } from "./TypeBadge";
import type { DetectedProject } from "./types";

interface Props {
  project: DetectedProject;
  onClick: () => void;
  index?: number;
}

const ProjectCard = ({ project, onClick, index = 0 }: Props) => (
  <motion.div
    initial={{ opacity: 0, y: 12 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ duration: 0.25, delay: index * 0.05 }}
    className="border border-border rounded-lg p-4 hover:border-primary/40 hover:shadow-md transition-all bg-card cursor-pointer group"
    onClick={onClick}
  >
    <div className="flex items-center gap-2 mb-1">
      <TypeBadge type={project.projectType} />
      <span className="font-mono text-sm font-semibold text-foreground truncate group-hover:text-primary transition-colors">
        {project.projectName}
      </span>
    </div>
    <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-2">
      <MapPin className="h-3 w-3 shrink-0" />
      <span className="font-mono truncate">{project.relativePath === "." ? "(root)" : project.relativePath}</span>
    </div>
    <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-1">
      <FileText className="h-3 w-3 shrink-0" />
      <span className="font-mono">{project.primaryIndicator}</span>
    </div>
    {project.goMetadata && project.goMetadata.runnables.length > 0 && (
      <div className="flex flex-wrap gap-1.5 mt-2">
        {project.goMetadata.runnables.map((r) => (
          <span key={r.name} className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-muted text-xs font-mono">
            <FileCode className="h-3 w-3 text-primary" />
            {r.name}
          </span>
        ))}
      </div>
    )}
    {project.csharpMetadata && project.csharpMetadata.projectFiles.length > 0 && (
      <div className="flex flex-wrap gap-1.5 mt-2">
        {project.csharpMetadata.projectFiles.map((f) => (
          <Badge key={f.fileName} variant="outline" className="text-[10px] px-1.5 py-0">{f.fileName}</Badge>
        ))}
      </div>
    )}
  </motion.div>
);

export default ProjectCard;
