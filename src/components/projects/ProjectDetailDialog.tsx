import { format } from "date-fns";
import {
  FolderGit2,
  Code2,
  Hash,
  MapPin,
  FileText,
  FileCode,
  Clock,
  FolderOpen,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
import { ProjectTypes } from "@/components/projects/TypeBadge";
import { ROOT_RELATIVE_PATH, ROOT_RELATIVE_LABEL } from "@/constants";
import type { DetectedProject } from "@/components/projects/types";

interface DetailRowProps {
  icon: typeof MapPin;
  label: string;
  value: string;
}

const DetailRow = ({ icon: Icon, label, value }: DetailRowProps) => (
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
      <DialogContent className="sm:max-w-lg max-h-[85vh] overflow-y-auto">
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
          <DetailRow
            icon={MapPin}
            label="Relative Path"
            value={isRootPath ? ROOT_RELATIVE_LABEL : project.relativePath}
          />
          <DetailRow icon={FileText} label="Primary Indicator" value={project.primaryIndicator} />
          <DetailRow
            icon={Clock}
            label="Detected At"
            value={format(new Date(project.detectedAt), "PPpp")}
          />
        </div>

        {project.goMetadata && (
          <GoMetadataSection
            moduleName={project.goMetadata.moduleName}
            goVersion={project.goMetadata.goVersion}
            runnables={project.goMetadata.runnables}
          />
        )}

        {project.csharpMetadata && (
          <CSharpMetadataSection
            slnName={project.csharpMetadata.slnName}
            sdkVersion={project.csharpMetadata.sdkVersion}
            projectFiles={project.csharpMetadata.projectFiles}
          />
        )}
      </DialogContent>
    </Dialog>
  );
};

const GoMetadataSection = ({ moduleName, goVersion, runnables }: {
  moduleName: string;
  goVersion: string;
  runnables: { name: string; relativePath: string }[];
}) => {
  const hasRunnables = runnables.length > 0;

  return (
    <>
      <Separator />
      <div className="space-y-3">
        <h3 className="font-mono text-sm font-semibold text-foreground flex items-center gap-2">
          <Code2 className="h-4 w-4 text-primary" />
          Go Metadata
        </h3>
        <div className="grid grid-cols-2 gap-3">
          <MetadataCell label="Module" value={moduleName} isBreakAll />
          <MetadataCell label="Go Version" value={goVersion} />
        </div>
        {hasRunnables && (
          <div>
            <span className="text-xs text-muted-foreground block mb-2">Runnable Entry Points</span>
            <div className="space-y-1.5">
              {runnables.map((runnable) => (
                <div key={runnable.name} className="flex items-center gap-2 rounded-md bg-muted px-3 py-2">
                  <FileCode className="h-3.5 w-3.5 text-primary shrink-0" />
                  <span className="font-mono text-sm font-medium text-foreground">{runnable.name}</span>
                  <span className="font-mono text-xs text-muted-foreground ml-auto truncate">{runnable.relativePath}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </>
  );
};

const CSharpMetadataSection = ({ slnName, sdkVersion, projectFiles }: {
  slnName: string;
  sdkVersion: string;
  projectFiles: { fileName: string; targetFramework: string; outputType: string }[];
}) => {
  const hasProjectFiles = projectFiles.length > 0;

  return (
    <>
      <Separator />
      <div className="space-y-3">
        <h3 className="font-mono text-sm font-semibold text-foreground flex items-center gap-2">
          <Hash className="h-4 w-4 text-primary" />
          C# Metadata
        </h3>
        <div className="grid grid-cols-2 gap-3">
          <MetadataCell label="Solution" value={slnName} />
          <MetadataCell label="SDK Version" value={sdkVersion} />
        </div>
        {hasProjectFiles && (
          <div>
            <span className="text-xs text-muted-foreground block mb-2">Project Files</span>
            <div className="space-y-1.5">
              {projectFiles.map((file) => (
                <div key={file.fileName} className="flex items-center gap-2 rounded-md bg-muted px-3 py-2">
                  <FileCode className="h-3.5 w-3.5 text-primary shrink-0" />
                  <span className="font-mono text-sm text-foreground">{file.fileName}</span>
                  <Badge variant="outline" className="text-[10px] px-1.5 py-0 ml-auto">{file.targetFramework}</Badge>
                  <span className="text-xs text-muted-foreground">{file.outputType}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </>
  );
};

const MetadataCell = ({ label, value, isBreakAll = false }: { label: string; value: string; isBreakAll?: boolean }) => (
  <div className="rounded-md bg-muted p-3">
    <span className="text-xs text-muted-foreground block">{label}</span>
    <span className={`font-mono text-sm text-foreground ${isBreakAll ? "break-all" : ""}`}>{value}</span>
  </div>
);

export default ProjectDetailDialog;