import {
  Home,
  BookOpen,
  Rocket,
  Settings,
  Boxes,
  FolderOpen,
  Monitor,
  Hammer,
  GitBranch,
  GitCommit,
  Tag,
  Sun,
  Moon,
  FolderGit2,
  Clock,
  BarChart3,
  Search,
  Terminal,
  FileText,
  Flag,
  PlayCircle,
  Zap,
  Archive,
  Link2,
  KeyRound,
  Stethoscope,
  Layers,
  Palette,
  MonitorPlay,
  AlertTriangle,
  Move,
  GitMerge,
  GitCompare,
  Download,
  Trash2,
  Replace as ReplaceIcon,
  Wrench,
  Globe,
} from "lucide-react";
import { NavLink } from "@/components/NavLink";

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarFooter,
  useSidebar,
} from "@/components/ui/sidebar";
import { useTheme } from "@/hooks/useTheme";

const navItems = [
  { title: "Home", url: "/", icon: Home },
  { title: "Commands", url: "/commands", icon: BookOpen },
  { title: "scan (cmd)", url: "/scan-command", icon: Search },
  { title: "Clone Overview", url: "/clone", icon: GitBranch },
  { title: "clone (cmd)", url: "/clone-command", icon: GitBranch },
  { title: "clone-next (cmd)", url: "/clone-next-command", icon: GitBranch },
  { title: "scan/clone flags", url: "/scan-clone-flags", icon: Flag },
  { title: "Getting Started", url: "/getting-started", icon: Rocket },
  { title: "Install / Update gitmap", url: "/install-gitmap", icon: Download },
  { title: "Configuration", url: "/config", icon: Settings },
  { title: "Architecture", url: "/architecture", icon: Boxes },
  { title: "Watch", url: "/watch", icon: Monitor },
  { title: "Release", url: "/release", icon: Tag },
  { title: "GoMod", url: "/gomod", icon: GitBranch },
  { title: "Projects", url: "/projects", icon: FolderGit2 },
  { title: "Makefile", url: "/makefile", icon: Hammer },
  { title: "History", url: "/history", icon: Clock },
  { title: "Version History", url: "/version-history", icon: Clock },
  { title: "Stats", url: "/stats", icon: BarChart3 },
  { title: "Detection", url: "/project-detection", icon: Search },
  { title: "Generic CLI", url: "/generic-cli", icon: Terminal },
  { title: "Changelog", url: "/changelog", icon: FileText },
  { title: "Changelog Generate", url: "/changelog-generate", icon: GitCommit },
  { title: "Flags", url: "/flags", icon: Flag },
  { title: "Examples", url: "/examples", icon: PlayCircle },
  { title: "Interactive TUI", url: "/interactive", icon: Terminal },
  { title: "Batch Actions", url: "/batch-actions", icon: Zap },
  { title: "Clear Release JSON", url: "/clear-release-json", icon: FileText },
  { title: "Bookmarks", url: "/bookmarks", icon: BookOpen },
  { title: "Export", url: "/export", icon: FileText },
  { title: "Import", url: "/import", icon: FileText },
  { title: "Profile", url: "/profile", icon: FileText },
  { title: "Diff Profiles", url: "/diff-profiles", icon: FileText },
  { title: "Zip Groups", url: "/zip-group", icon: Archive },
  { title: "Aliases", url: "/alias", icon: Link2 },
  { title: "CD / Navigate", url: "/cd", icon: FolderOpen },
  { title: "SSH Keys", url: "/ssh", icon: KeyRound },
  { title: "Prune", url: "/prune", icon: Archive },
  { title: "Temp Release", url: "/temp-release", icon: Layers },
  { title: "Release Self", url: "/release-self", icon: Tag },
  { title: "Doctor", url: "/doctor", icon: Stethoscope },
  { title: "Troubleshooting", url: "/troubleshooting", icon: AlertTriangle },
  { title: "Clone Next", url: "/clone-next", icon: GitBranch },
  { title: "Dashboard", url: "/dashboard", icon: BarChart3 },
  { title: "Setup", url: "/setup", icon: Settings },
  { title: "Install", url: "/install", icon: Hammer },
  { title: "Uninstall", url: "/commands#uninstall", icon: Hammer },
  { title: "Help Dashboard", url: "/help-dashboard", icon: MonitorPlay },
  { title: "Help Index", url: "/help-index", icon: BookOpen },
  { title: "Diff", url: "/diff", icon: GitCompare },
  { title: "Move (mv)", url: "/mv", icon: Move },
  { title: "Merge Both", url: "/merge-both", icon: GitMerge },
  { title: "Merge Left", url: "/merge-left", icon: GitMerge },
  { title: "Merge Right", url: "/merge-right", icon: GitMerge },
  { title: "Commit In (cin)", url: "/commit-in", icon: GitCommit },
  { title: "Commit Left (planned)", url: "/commit-left", icon: GitCommit },
  { title: "Commit Right (planned)", url: "/commit-right", icon: GitCommit },
  { title: "Commit Both (planned)", url: "/commit-both", icon: GitCommit },
  { title: "Replace (rpl)", url: "/replace", icon: ReplaceIcon },
  { title: "Fix Repo (fr)", url: "/fix-repo", icon: Wrench },
  { title: "Clone + Fix Repo (cfr)", url: "/clone-fix-repo", icon: Wrench },
  { title: "Make Public Repo", url: "/make-public", icon: Globe },
  { title: "Make All Public (MAPUB)", url: "/make-all-public", icon: Globe },
  { title: "Make All Private (MAPRI)", url: "/make-all-private", icon: Globe },
  { title: "History Rewrite (purge/pin)", url: "/history-rewrite", icon: Trash2 },
  { title: "Register (as)", url: "/as", icon: Tag },
  { title: "Release Alias", url: "/release-alias", icon: Rocket },
  { title: "Release Alias Pull", url: "/release-alias-pull", icon: Rocket },
  { title: "Release Undo (ru)", url: "/commands#release-undo", icon: Rocket },
  { title: "Spec Index", url: "/spec", icon: BookOpen },
  { title: "Post-Mortems", url: "/post-mortems", icon: AlertTriangle },
  { title: "Design System", url: "/design-system", icon: Palette },
  { title: "Migration Guide", url: "/changelog#migration-guide-v2x--v300", icon: FileText },
];

export function DocsSidebar() {
  const { state } = useSidebar();
  const collapsed = state === "collapsed";
  const { isDark: dark, toggleTheme } = useTheme();

  return (
    <Sidebar collapsible="icon" className="border-r border-sidebar-border bg-sidebar">
      <SidebarContent>
        <SidebarGroup>
          {!collapsed && (
            <div className="border-b border-sidebar-border px-3 py-4">
              <div className="mb-1 text-[11px] font-mono uppercase tracking-[0.16em] text-muted-foreground">
                Workspace
              </div>
              <div className="flex items-center gap-2">
                <span className="font-mono text-lg font-bold text-sidebar-primary">gitmap</span>
                <span className="text-xs font-mono text-sidebar-foreground/70">docs</span>
              </div>
            </div>
          )}
          {collapsed && (
            <div className="flex justify-center py-4">
              <span className="font-mono text-lg font-bold text-sidebar-primary">g</span>
            </div>
          )}
          <SidebarGroupLabel className="px-3 pt-3 text-[11px] font-mono uppercase tracking-[0.16em] text-muted-foreground">
            Explorer
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <NavLink
                      to={item.url}
                      end={item.url === "/"}
                      className="flex min-h-8 items-center rounded-sm border border-transparent px-2 text-sm text-sidebar-foreground/85 hover:border-sidebar-border hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                      activeClassName="border-sidebar-border bg-sidebar-accent text-sidebar-primary shadow-sm"
                    >
                      <item.icon className="mr-2 h-4 w-4" />
                      {!collapsed && <span>{item.title}</span>}
                    </NavLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <button
          onClick={toggleTheme}
          className="flex w-full items-center gap-2 rounded-sm border border-sidebar-border px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
        >
          {dark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          {!collapsed && <span>{dark ? "Light mode" : "Dark mode"}</span>}
        </button>
      </SidebarFooter>
    </Sidebar>
  );
}
