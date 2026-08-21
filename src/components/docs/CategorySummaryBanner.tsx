import { commands, Categories, CommandCategory } from "@/data/commands";

interface CategorySummaryButtonProps {
  category: CommandCategory;
  count: number;
  onClick: () => void;
}

interface CategorySummaryBannerProps {
  onSelectCategory: (key: string) => void;
}

const CategorySummaryButton = ({ category, count, onClick }: CategorySummaryButtonProps) => {
  return (
    <button
      onClick={onClick}
      className="rounded-lg border border-border bg-card px-3 py-2.5 text-left hover:bg-muted/50 hover:border-primary/40 transition-all duration-200 cursor-pointer group"
    >
      <div className="flex items-center gap-2 mb-1">
        {category.icon && <span className="text-base">{category.icon}</span>}
        <span className="text-lg font-mono font-bold text-primary">{count}</span>
      </div>
      <div className="text-[11px] text-muted-foreground font-mono leading-tight truncate group-hover:text-foreground transition-colors">
        {category.label}
      </div>
    </button>
  );
};

export const CategorySummaryBanner = ({ onSelectCategory }: CategorySummaryBannerProps) => {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2 mb-6">
      {Categories.map((category) => (
        <CategorySummaryButton
          key={category.key}
          category={category}
          count={commands.filter((cmd) => cmd.category === category.key).length}
          onClick={() => onSelectCategory(category.key)}
        />
      ))}
    </div>
  );
};

export default CategorySummaryBanner;
