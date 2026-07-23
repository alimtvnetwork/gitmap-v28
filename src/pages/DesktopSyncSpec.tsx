import SpecPage from "@/components/docs/SpecPage";
import md from "../../spec/01-app/11-desktop-sync.md?raw";

const DesktopSyncSpecPage = () => (
  <SpecPage
    title="gitmap desktop-sync (ds = gd)"
    subtitle="Register repos with GitHub Desktop. As of v3.37.0 ds is an alias of github-desktop — no scan required."
    sourcePath="spec/01-app/11-desktop-sync.md"
    markdown={md}
    relatedLinks={[
      { label: "GitHub Desktop spec", to: "/github-desktop", description: "Same command, different name" },
      { label: "scan gd (spec 102)", to: "/scan-gd", description: "Bulk register from DB" },
      { label: "Spec index", to: "/spec" },
    ]}
  />
);

export default DesktopSyncSpecPage;
