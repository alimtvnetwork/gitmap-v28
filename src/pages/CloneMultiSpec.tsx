import SpecPage from "@/components/docs/SpecPage";
import md from "../../spec/01-app/104-clone-multi.md?raw";

const CloneMultiSpecPage = () => (
  <SpecPage
    title="gitmap clone — Multiple URLs"
    subtitle="Clone many repos in one invocation. Space-separated, comma-separated, or both. Planned for v3.38.0."
    sourcePath="spec/01-app/104-clone-multi.md"
    markdown={md}
    relatedLinks={[
      { label: "Single-URL clone", to: "/commands", description: "Existing direct-URL behaviour this extends" },
      { label: "desktop-sync / gd", to: "/desktop-sync", description: "Used by --github-desktop after each clone" },
      { label: "Spec index", to: "/spec" },
    ]}
  />
);

export default CloneMultiSpecPage;
