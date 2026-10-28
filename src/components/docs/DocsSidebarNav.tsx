import React from "react";
import { LucideIcon } from "lucide-react";
import { NavLink } from "@/components/NavLink";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

export interface NavItemDef {
  title: string;
  url: string;
  icon: LucideIcon;
}

interface DocsSidebarNavProps {
  items: NavItemDef[];
  isCollapsed: boolean;
}

export const DocsSidebarNav: React.FC<DocsSidebarNavProps> = ({ items, isCollapsed }) => {
  return (
    <SidebarGroup>
      {!isCollapsed && (
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
      {isCollapsed && (
        <div className="flex justify-center py-4">
          <span className="font-mono text-lg font-bold text-sidebar-primary">g</span>
        </div>
      )}
      <SidebarGroupLabel className="px-3 pt-3 text-[11px] font-mono uppercase tracking-[0.16em] text-muted-foreground">
        Explorer
      </SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          {items.map((item) => (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild>
                <NavLink
                  to={item.url}
                  end={item.url === "/"}
                  className="flex min-h-8 items-center rounded-sm border border-transparent px-2 text-sm text-sidebar-foreground/85 hover:border-sidebar-border hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                  activeClassName="border-sidebar-border bg-sidebar-accent text-sidebar-primary shadow-sm"
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {!isCollapsed && <span>{item.title}</span>}
                </NavLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
};
