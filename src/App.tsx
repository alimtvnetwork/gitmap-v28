import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import {
  DOCS_TOOLTIP_DELAY_MS,
  DOCS_TOOLTIP_SKIP_DELAY_MS,
} from "@/components/docs/docsTooltipConfig";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { MotionConfig } from "framer-motion";
import Index from "./pages/Index";
import Commands from "./pages/Commands";
import GoMod from "./pages/GoMod";
import Projects from "./pages/Projects";
import GettingStarted from "./pages/GettingStarted";
import Config from "./pages/Config";
import Architecture from "./pages/Architecture";
import Watch from "./pages/Watch";
import Release from "./pages/Release";
import ReleaseVersionPage from "./pages/ReleaseVersion";
import MakefilePage from "./pages/Makefile";
import HistoryPage from "./pages/History";
import StatsPage from "./pages/Stats";
import ProjectDetectionPage from "./pages/ProjectDetection";
import GenericCLIPage from "./pages/GenericCLI";
import ChangelogPage from "./pages/Changelog";
import ChangelogGeneratePage from "./pages/ChangelogGenerate";
import FlagReferencePage from "./pages/FlagReference";
import InteractiveExamplesPage from "./pages/InteractiveExamples";
import InteractiveTUIPage from "./pages/InteractiveTUI";
import BatchActionsPage from "./pages/BatchActions";
import ClearReleaseJSONPage from "./pages/ClearReleaseJSON";
import BookmarksPage from "./pages/Bookmarks";
import ExportPage from "./pages/Export";
import ImportPage from "./pages/Import";
import ProfilePage from "./pages/Profile";
import DiffProfilesPage from "./pages/DiffProfiles";
import NotFound from "./pages/NotFound";
import ZipGroupPage from "./pages/ZipGroup";
import AliasPage from "./pages/Alias";
import SSHPage from "./pages/SSH";
import PrunePage from "./pages/Prune";
import DoctorPage from "./pages/Doctor";
import TempReleasePage from "./pages/TempRelease";
import ReleaseSelfPage from "./pages/ReleaseSelf";
import DashboardPage from "./pages/Dashboard";
import CloneNextPage from "./pages/CloneNext";
import SpecIndexPage from "./pages/SpecIndex";
import CdPage from "./pages/Cd";
import SetupPage from "./pages/Setup";
import DesignSystemPage from "./pages/DesignSystem";
import InstallPage from "./pages/Install";
import InstallGitmapPage from "./pages/InstallGitmap";
import HelpDashboardPage from "./pages/HelpDashboard";
import HelpIndexPage from "./pages/HelpIndex";
import PostMortemsPage from "./pages/PostMortems";
import VersionHistoryPage from "./pages/VersionHistory";
import ScanAllSpecPage from "./pages/ScanAllSpec";
import DesktopSyncSpecPage from "./pages/DesktopSyncSpec";
import GitHubDesktopSpecPage from "./pages/GitHubDesktopSpec";
import ScanGdSpecPage from "./pages/ScanGdSpec";
import CloneMultiSpecPage from "./pages/CloneMultiSpec";
import ScanCommandPage from "./pages/ScanCommand";
import CloneCommandPage from "./pages/CloneCommand";
import CloneOverviewPage from "./pages/CloneOverview";
import CloneNextCommandPage from "./pages/CloneNextCommand";
import ScanCloneFlagsPage from "./pages/ScanCloneFlags";
import TroubleshootingPage from "./pages/Troubleshooting";
import DiffPage from "./pages/Diff";
import MvPage from "./pages/Mv";
import MergeBothPage from "./pages/MergeBoth";
import MergeLeftPage from "./pages/MergeLeft";
import MergeRightPage from "./pages/MergeRight";
import CommitLeftPage from "./pages/CommitLeft";
import CommitRightPage from "./pages/CommitRight";
import CommitBothPage from "./pages/CommitBoth";
import HistoryRewritePage from "./pages/HistoryRewrite";
import ReleaseAliasPage from "./pages/ReleaseAlias";
import ReplacePage from "./pages/Replace";
import FixRepoPage from "./pages/FixRepo";
import CloneFixRepoPage from "./pages/CloneFixRepo";
import MakePublicPage from "./pages/MakePublic";
import MakeAllPublicPage from "./pages/MakeAllPublic";
import MakeAllPrivatePage from "./pages/MakeAllPrivate";
import CommitInPage from "./pages/CommitIn";
import ChromeProfileSpecPage from "./pages/ChromeProfileSpec";
import CommandPalette from "./components/docs/CommandPalette";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <MotionConfig reducedMotion="user">
    <TooltipProvider
      delayDuration={DOCS_TOOLTIP_DELAY_MS}
      skipDelayDuration={DOCS_TOOLTIP_SKIP_DELAY_MS}
    >
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <CommandPalette />
        <Routes>
          <Route path="/" element={<Index />} />
          <Route path="/commands" element={<Commands />} />
          <Route path="/getting-started" element={<GettingStarted />} />
          <Route path="/config" element={<Config />} />
          <Route path="/architecture" element={<Architecture />} />
          <Route path="/watch" element={<Watch />} />
          <Route path="/release" element={<Release />} />
          <Route path="/release/:version" element={<ReleaseVersionPage />} />
          <Route path="/gomod" element={<GoMod />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/makefile" element={<MakefilePage />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/stats" element={<StatsPage />} />
          <Route path="/project-detection" element={<ProjectDetectionPage />} />
          <Route path="/generic-cli" element={<GenericCLIPage />} />
          <Route path="/changelog" element={<ChangelogPage />} />
          <Route path="/changelog-generate" element={<ChangelogGeneratePage />} />
          <Route path="/flags" element={<FlagReferencePage />} />
          <Route path="/examples" element={<InteractiveExamplesPage />} />
          <Route path="/interactive" element={<InteractiveTUIPage />} />
          <Route path="/batch-actions" element={<BatchActionsPage />} />
          <Route path="/clear-release-json" element={<ClearReleaseJSONPage />} />
          <Route path="/bookmarks" element={<BookmarksPage />} />
          <Route path="/export" element={<ExportPage />} />
          <Route path="/import" element={<ImportPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/diff-profiles" element={<DiffProfilesPage />} />
          <Route path="/zip-group" element={<ZipGroupPage />} />
          <Route path="/alias" element={<AliasPage />} />
          <Route path="/ssh" element={<SSHPage />} />
          <Route path="/prune" element={<PrunePage />} />
          <Route path="/doctor" element={<DoctorPage />} />
          <Route path="/temp-release" element={<TempReleasePage />} />
          <Route path="/release-self" element={<ReleaseSelfPage />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/clone-next" element={<CloneNextPage />} />
          <Route path="/version-history" element={<VersionHistoryPage />} />
          <Route path="/spec" element={<SpecIndexPage />} />
          <Route path="/cd" element={<CdPage />} />
          <Route path="/setup" element={<SetupPage />} />
          <Route path="/design-system" element={<DesignSystemPage />} />
          <Route path="/install" element={<InstallPage />} />
          <Route path="/install-gitmap" element={<InstallGitmapPage />} />
          <Route path="/help-dashboard" element={<HelpDashboardPage />} />
          <Route path="/help-index" element={<HelpIndexPage />} />
          <Route path="/post-mortems" element={<PostMortemsPage />} />
          <Route path="/scan-all" element={<ScanAllSpecPage />} />
          <Route path="/desktop-sync" element={<DesktopSyncSpecPage />} />
          <Route path="/github-desktop" element={<GitHubDesktopSpecPage />} />
          <Route path="/scan-gd" element={<ScanGdSpecPage />} />
          <Route path="/clone-multi" element={<CloneMultiSpecPage />} />
          <Route path="/scan-command" element={<ScanCommandPage />} />
          <Route path="/clone-command" element={<CloneCommandPage />} />
          <Route path="/clone" element={<CloneOverviewPage />} />
          <Route path="/clone-next-command" element={<CloneNextCommandPage />} />
          <Route path="/scan-clone-flags" element={<ScanCloneFlagsPage />} />
          <Route path="/troubleshooting" element={<TroubleshootingPage />} />
          <Route path="/diff" element={<DiffPage />} />
          <Route path="/mv" element={<MvPage />} />
          <Route path="/merge-both" element={<MergeBothPage />} />
          <Route path="/merge-left" element={<MergeLeftPage />} />
          <Route path="/merge-right" element={<MergeRightPage />} />
          <Route path="/commit-left" element={<CommitLeftPage />} />
          <Route path="/commit-right" element={<CommitRightPage />} />
          <Route path="/commit-both" element={<CommitBothPage />} />
          <Route path="/history-rewrite" element={<HistoryRewritePage />} />
          <Route path="/as" element={<ReleaseAliasPage />} />
          <Route path="/release-alias" element={<ReleaseAliasPage />} />
          <Route path="/release-alias-pull" element={<ReleaseAliasPage />} />
          <Route path="/replace" element={<ReplacePage />} />
          <Route path="/fix-repo" element={<FixRepoPage />} />
          <Route path="/clone-fix-repo" element={<CloneFixRepoPage />} />
          <Route path="/make-public" element={<MakePublicPage />} />
          <Route path="/make-all-public" element={<MakeAllPublicPage />} />
          <Route path="/mapub" element={<MakeAllPublicPage />} />
          <Route path="/make-all-private" element={<MakeAllPrivatePage />} />
          <Route path="/mapri" element={<MakeAllPrivatePage />} />
          <Route path="/commit-in" element={<CommitInPage />} />
          <Route path="/chrome-profile" element={<ChromeProfileSpecPage />} />
          <Route path="/cpc" element={<ChromeProfileSpecPage />} />
          <Route path="/cpe" element={<ChromeProfileSpecPage />} />
          <Route path="/cpi" element={<ChromeProfileSpecPage />} />
          <Route path="/cpl" element={<ChromeProfileSpecPage />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
    </MotionConfig>
  </QueryClientProvider>
);

export default App;
