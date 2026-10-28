export interface ClassifiedCell {
  cellType: "null" | "bool" | "blob" | "json" | "multiline" | "text" | "plain";
  badgeText: string;
  rawValue: unknown;
  sizeBytes: number;
}

export interface InspectTarget {
  columnName: string;
  classified: ClassifiedCell;
}
