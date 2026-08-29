import { useState, useMemo } from "react";
import DocsLayout from "@/components/docs/DocsLayout";
import SearchBar from "@/components/docs/SearchBar";
import { commands, CommandDef } from "@/data/commands";

export enum SortColType {
  Flag = "flag",
  Command = "command",
}

export interface FlagRow {
  flag: string;
  description: string;
  command: string;
  alias?: string;
}

interface FlagTableHeaderProps {
  sortCol: SortColType;
  isSortAsc: boolean;
  onSort: (column: SortColType) => void;
}

interface FlagTableProps {
  rows: FlagRow[];
  search: string;
  sortCol: SortColType;
  isSortAsc: boolean;
  onSort: (column: SortColType) => void;
}

function extractAllFlags(commandList: CommandDef[]): FlagRow[] {
  const rows: FlagRow[] = [];

  for (const commandItem of commandList) {
    const hasFlags = Boolean(commandItem.flags && commandItem.flags.length > 0);

    if (hasFlags === false) continue;

    for (const flagItem of commandItem.flags!) {
      rows.push({ flag: flagItem.flag, description: flagItem.description, command: commandItem.name, alias: commandItem.alias });
    }
  }

  return rows;
}

function matchesFlagSearch(rowItem: FlagRow, query: string): boolean {
  const isFlagMatch = rowItem.flag.toLowerCase().includes(query);
  const isDescMatch = rowItem.description.toLowerCase().includes(query);
  const isCmdMatch = rowItem.command.includes(query);

  return isFlagMatch || isDescMatch || isCmdMatch;
}

function compareFlagRows(rowA: FlagRow, rowB: FlagRow, sortCol: SortColType, isSortAsc: boolean): number {
  const isFlagCol = sortCol === SortColType.Flag;
  const valueA = isFlagCol ? rowA.flag : rowA.command;
  const valueB = isFlagCol ? rowB.flag : rowB.command;

  return isSortAsc ? valueA.localeCompare(valueB) : valueB.localeCompare(valueA);
}

function filterAndSortFlags(
  allFlags: FlagRow[],
  search: string,
  sortCol: SortColType,
  isSortAsc: boolean
): FlagRow[] {
  const query = search.toLowerCase();
  const searchApplied = search ? allFlags.filter((row) => matchesFlagSearch(row, query)) : allFlags;
  const sortedRows = [...searchApplied];
  sortedRows.sort((rowA, rowB) => compareFlagRows(rowA, rowB, sortCol, isSortAsc));

  return sortedRows;
}

function getSortIndicator(currentCol: SortColType, targetCol: SortColType, isSortAsc: boolean): string {
  const isCurrent = currentCol === targetCol;

  if (isCurrent === false) return "";

  return isSortAsc ? " ↑" : " ↓";
}

const FlagTableHeader = ({ sortCol, isSortAsc, onSort }: FlagTableHeaderProps) => {
  return (
    <thead>
      <tr className="bg-muted/30 border-b border-border">
        <th
          onClick={() => onSort(SortColType.Flag)}
          className="text-left px-4 py-2.5 font-mono font-semibold text-foreground cursor-pointer hover:text-primary transition-colors select-none"
        >
          Flag{getSortIndicator(sortCol, SortColType.Flag, isSortAsc)}
        </th>
        <th className="text-left px-4 py-2.5 font-mono font-semibold text-foreground">
          Description
        </th>
        <th
          onClick={() => onSort(SortColType.Command)}
          className="text-left px-4 py-2.5 font-mono font-semibold text-foreground cursor-pointer hover:text-primary transition-colors select-none"
        >
          Command{getSortIndicator(sortCol, SortColType.Command, isSortAsc)}
        </th>
      </tr>
    </thead>
  );
};

const FlagTableRow = ({ row }: { row: FlagRow }) => {
  return (
    <tr className="border-b border-border last:border-0 hover:bg-muted/20 transition-colors">
      <td className="px-4 py-2 font-mono text-primary whitespace-nowrap">
        {row.flag}
      </td>
      <td className="px-4 py-2 text-muted-foreground">{row.description}</td>
      <td className="px-4 py-2 font-mono whitespace-nowrap">
        {row.command}
        {row.alias && (
          <span className="text-muted-foreground ml-1">({row.alias})</span>
        )}
      </td>
    </tr>
  );
};

const FlagTable = ({ rows, search, sortCol, isSortAsc, onSort }: FlagTableProps) => {
  return (
    <div className="mt-6 rounded-lg border border-border overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <FlagTableHeader sortCol={sortCol} isSortAsc={isSortAsc} onSort={onSort} />
          <tbody>
            {rows.map((row, index) => (
              <FlagTableRow key={`${row.command}-${row.flag}-${index}`} row={row} />
            ))}
            {rows.length === 0 && (
              <tr>
                <td colSpan={3} className="px-4 py-8 text-center text-muted-foreground font-mono text-sm">
                  No flags matching "{search}"
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

const FlagReferencePage = () => {
  const [search, setSearch] = useState("");
  const [sortCol, setSortCol] = useState<SortColType>(SortColType.Flag);
  const [isSortAsc, setIsSortAsc] = useState(true);

  const allFlags = useMemo(() => extractAllFlags(commands), []);
  const filtered = useMemo(() => filterAndSortFlags(allFlags, search, sortCol, isSortAsc), [allFlags, search, sortCol, isSortAsc]);
  const commandsWithFlagsCount = useMemo(() => commands.filter((cmd) => (cmd.flags?.length ?? 0) > 0).length, []);

  const handleSort = (column: SortColType) => {
    const isSameColumn = sortCol === column;

    if (isSameColumn) setIsSortAsc((current) => current === false);
    else { setSortCol(column); setIsSortAsc(true); }
  };

  return (
    <DocsLayout>
      <h1 className="text-3xl font-heading font-bold mb-2 docs-h1">Flag Reference</h1>
      <p className="text-muted-foreground mb-6">
        {allFlags.length} flags across {commandsWithFlagsCount} commands.
      </p>
      <SearchBar value={search} onChange={setSearch} placeholder="Search flags..." />
      <FlagTable rows={filtered} search={search} sortCol={sortCol} isSortAsc={isSortAsc} onSort={handleSort} />
    </DocsLayout>
  );
};

export default FlagReferencePage;
