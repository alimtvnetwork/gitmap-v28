import SpecPage from "@/components/docs/SpecPage";
import md from "../../spec/01-app/10-github-desktop.md?raw";

const GitHubDesktopSpecPage = () => (
  <SpecPage
    title="gitmap github-desktop (gd)"
    subtitle="Register a repo with GitHub Desktop. Alias of desktop-sync (ds) as of v3.37.0."
    sourcePath="spec/01-app/10-github-desktop.md"
    markdown={md}
    relatedLinks={[
      { label: "Desktop Sync spec", to: "/desktop-sync", description: "Same command, different name" },
      { label: "scan gd (spec 102)", to: "/scan-gd", description: "Bulk register from DB without re-walking" },
      { label: "Spec index", to: "/spec" },
    ]}
  />
);

export default GitHubDesktopSpecPage;
