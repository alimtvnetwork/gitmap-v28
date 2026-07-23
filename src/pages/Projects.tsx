import { useState } from "react";
import { motion } from "framer-motion";
import { FolderGit2, Search, Filter, ArrowRight } from "lucide-react";
import { Link } from "react-router-dom";
import DocsLayout from "@/components/docs/DocsLayout";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import ProjectDetailDialog from "@/components/projects/ProjectDetailDialog";
import RepoGroup from "@/components/projects/RepoGroup";
import { ProjectTypes } from "@/components/projects/TypeBadge";
import type { DetectedProject, ProjectType, ProjectFilter } from "@/components/projects/types";
import { FILTER_ALL } from "@/constants";

const SAMPLE_PROJECTS: DetectedProject[] = [
  {
    id: "1", repoName: "my-api", projectType: "go", projectName: "github.com/user/my-api",
    absolutePath: "/home/user/repos/my-api", repoPath: "/home/user/repos/my-api",
    relativePath: ".", primaryIndicator: "go.mod", detectedAt: "2026-03-11T09:54:00Z",
    goMetadata: { moduleName: "github.com/user/my-api", goVersion: "1.22", runnables: [
      { name: "server", relativePath: "cmd/server/main.go" },
      { name: "worker", relativePath: "cmd/worker/main.go" },
    ]},
  },
  {
    id: "2", repoName: "my-api", projectType: "react", projectName: "admin-dashboard",
    absolutePath: "/home/user/repos/my-api/web", repoPath: "/home/user/repos/my-api",
    relativePath: "web", primaryIndicator: "package.json", detectedAt: "2026-03-11T09:54:00Z",
  },
  {
    id: "3", repoName: "infra-tools", projectType: "go", projectName: "github.com/user/infra-tools",
    absolutePath: "/home/user/repos/infra-tools", repoPath: "/home/user/repos/infra-tools",
    relativePath: ".", primaryIndicator: "go.mod", detectedAt: "2026-03-11T09:55:00Z",
    goMetadata: { moduleName: "github.com/user/infra-tools", goVersion: "1.23", runnables: [
      { name: "infra-tools", relativePath: "main.go" },
    ]},
  },
  {
    id: "4", repoName: "web-platform", projectType: "react", projectName: "@platform/frontend",
    absolutePath: "/home/user/repos/web-platform", repoPath: "/home/user/repos/web-platform",
    relativePath: ".", primaryIndicator: "package.json", detectedAt: "2026-03-11T09:56:00Z",
  },
  {
    id: "5", repoName: "web-platform", projectType: "node", projectName: "@platform/api",
    absolutePath: "/home/user/repos/web-platform/api", repoPath: "/home/user/repos/web-platform",
    relativePath: "api", primaryIndicator: "package.json", detectedAt: "2026-03-11T09:56:00Z",
  },
  {
    id: "6", repoName: "signal-engine", projectType: "cpp", projectName: "signal-engine",
    absolutePath: "/home/user/repos/signal-engine", repoPath: "/home/user/repos/signal-engine",
    relativePath: ".", primaryIndicator: "CMakeLists.txt", detectedAt: "2026-03-11T09:57:00Z",
  },
  {
    id: "7", repoName: "enterprise-app", projectType: "csharp", projectName: "EnterpriseApp",
    absolutePath: "/home/user/repos/enterprise-app", repoPath: "/home/user/repos/enterprise-app",
    relativePath: ".", primaryIndicator: "EnterpriseApp.sln", detectedAt: "2026-03-11T09:58:00Z",
    csharpMetadata: { slnName: "EnterpriseApp.sln", sdkVersion: "8.0.100", projectFiles: [
      { fileName: "EnterpriseApp.Api.csproj", targetFramework: "net8.0", outputType: "Exe" },
      { fileName: "EnterpriseApp.Core.csproj", targetFramework: "net8.0", outputType: "Library" },
    ]},
  },
];

const ProjectsPage = () => {
  const [search, setSearch] = useState("");
  const [activeFilter, setActiveFilter] = useState<ProjectFilter>(FILTER_ALL);
  const [selectedProject, setSelectedProject] = useState<DetectedProject | null>(null);

  const filtered = SAMPLE_PROJECTS.filter((project) => {
    if (activeFilter !== FILTER_ALL && project.projectType !== activeFilter) return false;
    if (search.length > 0) {
      const searchLower = search.toLowerCase();

      return project.projectName.toLowerCase().includes(searchLower) || project.repoName.toLowerCase().includes(searchLower) || project.absolutePath.toLowerCase().includes(searchLower);
    }
    return true;
  });

  const grouped = filtered.reduce<Record<string, DetectedProject[]>>((acc, p) => {
    acc[p.repoName] = acc[p.repoName] || [];
    acc[p.repoName].push(p);
    return acc;
  }, {});

  const typeCounts = SAMPLE_PROJECTS.reduce<Record<string, number>>((acc, p) => {
    acc[p.projectType] = (acc[p.projectType] || 0) + 1;
    return acc;
  }, {});

  return (
    <DocsLayout>
      <div className="space-y-6">
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <h1 className="text-2xl sm:text-3xl font-mono font-bold text-foreground flex items-center gap-3">
            <FolderGit2 className="h-7 w-7 sm:h-8 sm:w-8 text-primary" />
            Detected Projects
          </h1>
          <p className="text-muted-foreground mt-2 text-sm sm:text-base">
            Projects discovered inside Git repositories during scan. Click any project to see full details.
          </p>
          <Link
            to="/project-detection"
            className="inline-flex items-center gap-1.5 mt-3 text-xs font-mono text-primary hover:underline"
          >
            How detection works <ArrowRight className="h-3 w-3" />
          </Link>
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3, delay: 0.1 }}
          className="grid grid-cols-3 sm:grid-cols-3 md:grid-cols-6 gap-2 sm:gap-3"
        >
          <Card
            className={`cursor-pointer transition-all ${activeFilter === FILTER_ALL ? "ring-2 ring-primary shadow-sm" : "hover:border-primary/40"}`}
            onClick={() => setActiveFilter(FILTER_ALL)}
          >
            <CardContent className="p-2 sm:p-3 text-center">
              <div className="text-xl sm:text-2xl font-mono font-bold text-foreground">{SAMPLE_PROJECTS.length}</div>
              <div className="text-[10px] sm:text-xs text-muted-foreground">All</div>
            </CardContent>
          </Card>
          {(Object.entries(ProjectTypes) as [ProjectType, typeof ProjectTypes[ProjectType]][]).map(([key, config]) => (
            <Card
              key={key}
              className={`cursor-pointer transition-all ${activeFilter === key ? "ring-2 ring-primary shadow-sm" : "hover:border-primary/40"}`}
              onClick={() => setActiveFilter(key)}
            >
              <CardContent className="p-2 sm:p-3 text-center">
                <div className="text-xl sm:text-2xl font-mono font-bold text-foreground">{typeCounts[key] || 0}</div>
                <div className="text-[10px] sm:text-xs text-muted-foreground">{config.label}</div>
              </CardContent>
            </Card>
          ))}
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3, delay: 0.15 }}
          className="relative"
        >
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Search by project name, repo, or path..." value={search} onChange={(e) => setSearch(e.target.value)} className="pl-10 font-mono text-sm" />
        </motion.div>

        <div className="space-y-6">
          {Object.entries(grouped).map(([repoName, projects]) => (
            <RepoGroup
              key={repoName}
              repoName={repoName}
              projects={projects}
              onSelectProject={setSelectedProject}
            />
          ))}
        </div>

        {filtered.length === 0 && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="text-center py-12 text-muted-foreground"
          >
            <Filter className="h-8 w-8 mx-auto mb-3 opacity-50" />
            <p className="font-mono">No projects match your filters.</p>
          </motion.div>
        )}
      </div>

      <ProjectDetailDialog
        project={selectedProject}
        open={!!selectedProject}
        onOpenChange={(open) => { if (!open) setSelectedProject(null); }}
      />
    </DocsLayout>
  );
};

export default ProjectsPage;
