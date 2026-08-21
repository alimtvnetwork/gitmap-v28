import { CategoryType, Category, CategoryFilterType, CATEGORY_META } from "@/data/troubleshootingIssues";

interface TroubleshootingFilterButtonsProps {
  activeCategory: Category | CategoryFilterType;
  categoryCounts: Record<string, number>;
  onSelectCategory: (category: Category | CategoryFilterType) => void;
}

export const TroubleshootingFilterButtons = ({
  activeCategory,
  categoryCounts,
  onSelectCategory,
}: TroubleshootingFilterButtonsProps) => {
  const isAllActive = activeCategory === CategoryFilterType.All;
  const allButtonClasses = isAllActive
    ? "bg-primary text-primary-foreground border-primary"
    : "bg-background text-muted-foreground border-border hover:text-foreground hover:border-foreground/30";

  return (
    <div className="flex flex-wrap gap-2 mt-4 mb-8">
      <button onClick={() => onSelectCategory(CategoryFilterType.All)} className={`px-3 py-1.5 rounded-md text-sm font-mono border transition-colors ${allButtonClasses}`}>
        all ({categoryCounts[CategoryFilterType.All]})
      </button>
      {(Object.keys(CATEGORY_META) as CategoryType[]).map((categoryKey) => {
        const Icon = CATEGORY_META[categoryKey].icon;
        const isCategoryActive = activeCategory === categoryKey;
        const categoryClasses = isCategoryActive
          ? "bg-primary text-primary-foreground border-primary"
          : "bg-background text-muted-foreground border-border hover:text-foreground hover:border-foreground/30";

        return (
          <button
            key={categoryKey}
            onClick={() => onSelectCategory(categoryKey)}
            className={`px-3 py-1.5 rounded-md text-sm font-mono border transition-colors flex items-center gap-1.5 ${categoryClasses}`}
          >
            <Icon className="h-3.5 w-3.5" />
            {CATEGORY_META[categoryKey].label} ({categoryCounts[categoryKey] ?? 0})
          </button>
        );
      })}
    </div>
  );
};

export default TroubleshootingFilterButtons;
