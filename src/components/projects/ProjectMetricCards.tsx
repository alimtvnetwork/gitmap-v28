import { motion } from "framer-motion";
import { Card, CardContent } from "@/components/ui/card";
import { ProjectTypes } from "@/components/projects/TypeBadge";
import { ProjectType, type ProjectFilter } from "@/components/projects/types";
import { FILTER_ALL } from "@/constants";

interface ProjectMetricCardsProps {
  activeFilter: ProjectFilter;
  onSelectFilter: (filter: ProjectFilter) => void;
  totalProjects: number;
  typeCounts: Record<ProjectType, number>;
}

export const ProjectMetricCards = ({
  activeFilter,
  onSelectFilter,
  totalProjects,
  typeCounts,
}: ProjectMetricCardsProps) => {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.3, delay: 0.1 }}
      className="grid grid-cols-3 sm:grid-cols-3 md:grid-cols-6 gap-2 sm:gap-3"
    >
      <Card
        className={`cursor-pointer transition-all transform-gpu hover:-translate-y-0.5 active:scale-[0.98] ${
          activeFilter === FILTER_ALL
            ? "ring-2 ring-primary shadow-sm bg-accent/40"
            : "hover:border-primary/40 bg-card"
        }`}
        onClick={() => onSelectFilter(FILTER_ALL)}
      >
        <CardContent className="p-2 sm:p-3 text-center">
          <div className="text-xl sm:text-2xl font-mono font-bold text-foreground">
            {totalProjects}
          </div>
          <div className="text-[10px] sm:text-xs text-muted-foreground font-medium">All</div>
        </CardContent>
      </Card>
      {(Object.entries(ProjectTypes) as [ProjectType, (typeof ProjectTypes)[ProjectType]][]).map(
        ([key, config]) => (
          <Card
            key={key}
            className={`cursor-pointer transition-all transform-gpu hover:-translate-y-0.5 active:scale-[0.98] ${
              activeFilter === key
                ? "ring-2 ring-primary shadow-sm bg-accent/40"
                : "hover:border-primary/40 bg-card"
            }`}
            onClick={() => onSelectFilter(key)}
          >
            <CardContent className="p-2 sm:p-3 text-center">
              <div className="text-xl sm:text-2xl font-mono font-bold text-foreground">
                {typeCounts[key] || 0}
              </div>
              <div className="text-[10px] sm:text-xs text-muted-foreground font-medium">
                {config.label}
              </div>
            </CardContent>
          </Card>
        )
      )}
    </motion.div>
  );
};

export default ProjectMetricCards;
