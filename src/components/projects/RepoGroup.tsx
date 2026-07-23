import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { FolderGit2, ChevronDown } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import ProjectCard from "./ProjectCard";
import type { DetectedProject } from "./types";

interface Props {
  repoName: string;
  projects: DetectedProject[];
  onSelectProject: (project: DetectedProject) => void;
  defaultOpen?: boolean;
}

const RepoGroup = ({ repoName, projects, onSelectProject, defaultOpen = true }: Props) => {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <div>
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2 mb-3 w-full text-left group hover:opacity-80 transition-opacity"
      >
        <motion.div
          animate={{ rotate: open ? 0 : -90 }}
          transition={{ duration: 0.2 }}
        >
          <ChevronDown className="h-4 w-4 text-muted-foreground" />
        </motion.div>
        <FolderGit2 className="h-4 w-4 text-primary" />
        <h2 className="font-mono font-semibold text-foreground">{repoName}</h2>
        <span className="text-xs text-muted-foreground font-mono hidden sm:inline">{projects[0].repoPath}</span>
        <Badge variant="secondary" className="ml-auto text-xs">
          {projects.length} project{projects.length > 1 ? "s" : ""}
        </Badge>
      </button>
      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25, ease: "easeInOut" }}
            className="overflow-hidden"
          >
            <div className="grid gap-3 grid-cols-1 md:grid-cols-2">
              {projects.map((p, i) => (
                <ProjectCard key={p.id} project={p} onClick={() => onSelectProject(p)} index={i} />
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default RepoGroup;
