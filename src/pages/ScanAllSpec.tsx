import SpecPage from "@/components/docs/SpecPage";
import scanAllMarkdown from "../../spec/01-app/100-scan-all.md?raw";

const ScanAllSpecPage = () => (
  <SpecPage
    title="gitmap scan all — Bulk Re-Scan"
    subtitle="Re-scan every previously-scanned root in parallel. Planned for v3.33.0."
    sourcePath="spec/01-app/100-scan-all.md"
    markdown={scanAllMarkdown}
    relatedLinks={[
      { label: "Commands reference", to: "/commands", description: "Live entry under the Scanning category" },
      { label: "Spec index", to: "/spec", description: "Browse all specs" },
    ]}
  />
);

export default ScanAllSpecPage;
