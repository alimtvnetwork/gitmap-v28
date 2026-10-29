import { format } from "date-fns";
import { FolderGit2, MapPin, FileText, Clock, FolderOpen } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ProjectTypes } from "@/components/projects/TypeBadge";
import { ROOT_RELATIVE_PATH, ROOT_RELATIVE_LABEL } from "@/constants";
import type { DetectedProject } from "@/components/projects/types";
import { GoMetadataSection, CSharpMetadataSection } from "./ProjectDetailSections";

const DetailRow = ({ icon: Icon, label, value }: { icon: typeof MapPin; label: string; value: string }) => (
  <div className="flex items-start gap-3 py-2">
    <Icon className="h-4 w-4 text-muted-foreground mt-0.5 shrink-0" />
    <div className="min-w-0">
      <span className="text-xs text-muted-foreground block">{label}</span>
      <span className="font-mono text-sm text-foreground break-all">{value}</span>
    </div>
  </div>
);

interface ProjectDetailDialogProps {
  project: DetectedProject | null;
  open: boolean;
  onOpenChange: (isOpen: boolean) => void;
}

const ProjectDetailDialog = ({ project, open, onOpenChange }: ProjectDetailDialogProps) => {
  if (!project) return null;

  const typeConfig = ProjectTypes[project.projectType];
  const TypeIcon = typeConfig.icon;
  const isRootPath = project.relativePath === ROOT_RELATIVE_PATH;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg max-h-[85vh] overflow-y-auto transform-gpu will-change-[transform,opacity]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 font-mono text-lg">
            <Badge variant="outline" className={`${typeConfig.color} font-mono text-xs gap-1 border`}>
              <TypeIcon className="h-3 w-3" />
              {typeConfig.label}
            </Badge>
            <span className="truncate">{project.projectName}</span>
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-1">
          <DetailRow icon={FolderGit2} label="Repository" value={project.repoName} />
          <DetailRow icon={FolderOpen} label="Absolute Path" value={project.absolutePath} />
          <DetailRow icon={MapPin} label="Repo Path" value={project.repoPath} />
          <DetailRow icon={MapPin} label="Relative Path" value={isRootPath ? ROOT_RELATIVE_LABEL : project.relativePath} />
          <DetailRow icon={FileText} label="Primary Indicator" value={project.primaryIndicator} />
          <DetailRow icon={Clock} label="Detected At" value={format(new Date(project.detectedAt), "PPpp")} />
        </div>

        {project.goMetadata && <GoMetadataSection {...project.goMetadata} />}
        {project.csharpMetadata && <CSharpMetadataSection {...project.csharpMetadata} />}
      </DialogContent>
    </Dialog>
  );
};

export default ProjectDetailDialog;
