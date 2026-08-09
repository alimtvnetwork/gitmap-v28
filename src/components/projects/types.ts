import type { Code2 } from "lucide-react";

export enum ProjectType {
  Go = "go",
  Node = "node",
  React = "react",
  Cpp = "cpp",
  Csharp = "csharp",
}

export enum ProjectFilterType {
  All = "all",
  Go = "go",
  Node = "node",
  React = "react",
  Cpp = "cpp",
  Csharp = "csharp",
}
export type ProjectFilter = ProjectFilterType;

export interface GoRunnable {
  name: string;
  relativePath: string;
}

export interface GoMetadata {
  moduleName: string;
  goVersion: string;
  runnables: GoRunnable[];
}

export interface CSharpProjectFile {
  fileName: string;
  targetFramework: string;
  outputType: string;
}

export interface CSharpMetadata {
  slnName: string;
  sdkVersion: string;
  projectFiles: CSharpProjectFile[];
}

export interface DetectedProject {
  id: string;
  repoName: string;
  projectType: ProjectType;
  projectName: string;
  absolutePath: string;
  repoPath: string;
  relativePath: string;
  primaryIndicator: string;
  detectedAt: string;
  goMetadata?: GoMetadata;
  csharpMetadata?: CSharpMetadata;
}

export interface ProjectTypeConfig {
  label: string;
  color: string;
  icon: typeof Code2;
}
