import { Folder } from "lucide-react";

interface StorageTempCardProps {
  tempDir: string;
  onTempDirChange: (value: string) => void;
}

export function StorageTempCard({ tempDir, onTempDirChange }: StorageTempCardProps) {
  return (
    <div className="rounded-xl border border-border bg-card p-6 shadow-sm space-y-4">
      <div className="flex items-center gap-2 text-foreground font-semibold text-lg">
        <Folder className="h-5 w-5 text-primary" />
        <h2>Storage & Temp Directory</h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Directory path used for storing error logs with <code className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded text-foreground">--tempfile</code> and automated CI artifacts.
      </p>
      <div className="space-y-2">
        <label htmlFor="temp-dir-input" className="text-xs font-mono uppercase text-muted-foreground tracking-wider">
          Temp Directory Path
        </label>
        <input
          id="temp-dir-input"
          type="text"
          value={tempDir}
          onChange={(e) => onTempDirChange(e.target.value)}
          className="w-full px-3 py-2 rounded-md border border-input bg-background text-foreground text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
          placeholder=".lovable/temp"
        />
      </div>
    </div>
  );
}

export default StorageTempCard;
