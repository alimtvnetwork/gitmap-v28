import { motion, AnimatePresence } from "framer-motion";
import { FileText, AlertTriangle, Compass, Terminal, ShieldCheck, Palette, ChevronDown, Link as LinkIcon } from "lucide-react";
import type { SpecSection } from "./specData";
import SpecEntryRow from "./SpecEntryRow";
import { copyToClipboard } from "@/lib/clipboard";

const iconMap = {
  "file-text": <FileText className="h-5 w-5" />,
  "alert-triangle": <AlertTriangle className="h-5 w-5" />,
  compass: <Compass className="h-5 w-5" />,
  terminal: <Terminal className="h-5 w-5" />,
  "shield-check": <ShieldCheck className="h-5 w-5" />,
  palette: <Palette className="h-5 w-5" />,
};

interface SpecSectionCardProps {
  section: SpecSection;
  isCollapsed: boolean;
  onToggle: () => void;
}

const SpecSectionCard = ({ section, isCollapsed, onToggle }: SpecSectionCardProps) => {
  const anchorId = section.folder;

  const handleCopyLink = (e: React.MouseEvent) => {
    e.stopPropagation();
    const url = `${window.location.origin}${window.location.pathname}#${anchorId}`;
    void copyToClipboard(url);
  };

  return (
  <div id={anchorId} className="border border-border rounded-lg overflow-hidden scroll-mt-20">
    <button
      onClick={onToggle}
      className="group/header w-full bg-muted/30 px-5 py-4 border-b border-border text-left hover:bg-muted/40 transition-colors"
    >
      <div className="flex items-center gap-3 mb-1">
        <ChevronDown
          className={`h-4 w-4 text-muted-foreground transition-transform duration-200 ${isCollapsed ? "-rotate-90" : ""}`}
        />
        <span className={section.color}>{iconMap[section.iconName]}</span>
        <h2 className="text-lg font-mono font-semibold text-foreground">
          <span className="text-muted-foreground">{section.folder}/</span> {section.title}
        </h2>
        <span
          className="ml-auto flex items-center gap-2"
          onClick={handleCopyLink}
          title="Copy link to section"
        >
          <LinkIcon className="h-3 w-3 text-muted-foreground opacity-0 group-hover/header:opacity-100 hover:text-primary transition-all" />
          <span className="text-xs font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded">
            {section.entries.length} docs
          </span>
        </span>
      </div>
      <p className="text-sm text-muted-foreground ml-9">{section.description}</p>
    </button>

    <AnimatePresence initial={false}>
      {!isCollapsed && (
        <motion.div
          initial={{ height: 0 }}
          animate={{ height: "auto" }}
          exit={{ height: 0 }}
          transition={{ duration: 0.2 }}
          className="overflow-hidden"
        >
          <div className="divide-y divide-border">
            {section.entries.map((entry) => (
              <SpecEntryRow key={entry.id} entry={entry} />
            ))}
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  </div>
  );
};

export default SpecSectionCard;
