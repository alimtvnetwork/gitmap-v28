import SpecPage from "@/components/docs/SpecPage";
import md from "../../spec/01-app/102-scan-gd.md?raw";

const ScanGdSpecPage = () => (
  <SpecPage
    title="gitmap scan gd — Bulk GitHub Desktop registration"
    subtitle="Register every repo under the current scan root with GitHub Desktop. No re-walk. Planned for v3.35.0."
    sourcePath="spec/01-app/102-scan-gd.md"
    markdown={md}
    relatedLinks={[
      { label: "desktop-sync / gd", to: "/desktop-sync", description: "Single-repo and CWD-aware variant" },
      { label: "scan all (spec 100)", to: "/scan-all" },
      { label: "Spec index", to: "/spec" },
    ]}
  />
);

export default ScanGdSpecPage;
