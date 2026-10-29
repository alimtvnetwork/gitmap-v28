import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Code2, Hash, FileCode } from "lucide-react";

export const MetadataCell = ({ label, value, isBreakAll = false }: { label: string; value: string; isBreakAll?: boolean }) => (
  <div className="rounded-md bg-muted p-3">
    <span className="text-xs text-muted-foreground block">{label}</span>
    <span className={`font-mono text-sm text-foreground ${isBreakAll ? "break-all" : ""}`}>{value}</span>
  </div>
);

export const GoMetadataSection = ({
  moduleName,
  goVersion,
  runnables,
}: {
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

export const CSharpMetadataSection = ({
  slnName,
  sdkVersion,
  projectFiles,
}: {
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
