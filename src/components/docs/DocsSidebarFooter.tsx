import React from "react";
import { Sun, Moon } from "lucide-react";
import { SidebarFooter } from "@/components/ui/sidebar";

interface DocsSidebarFooterProps {
  isDark: boolean;
  isCollapsed: boolean;
  onToggleTheme: () => void;
}

export const DocsSidebarFooter: React.FC<DocsSidebarFooterProps> = ({
  isDark,
  isCollapsed,
  onToggleTheme,
}) => {
  return (
    <SidebarFooter>
      <button
        onClick={onToggleTheme}
        className="flex w-full items-center gap-2 rounded-sm border border-sidebar-border px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
      >
        {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        {!isCollapsed && <span>{isDark ? "Light mode" : "Dark mode"}</span>}
      </button>
    </SidebarFooter>
  );
};
