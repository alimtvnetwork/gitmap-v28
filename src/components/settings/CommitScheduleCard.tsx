import { Clock, Globe, Calendar, ShieldCheck } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export interface CommitScheduleCardProps {
  timezone: string;
  onTimezoneChange: (timezone: string) => void;
  isScheduleActive: boolean;
  onToggleSchedule: (isActive: boolean) => void;
  weekdayCutoff?: string;
}

const COMMON_TIMEZONES = [
  { value: "Local", label: "Local System Time" },
  { value: "UTC", label: "UTC (Coordinated Universal Time)" },
  { value: "America/New_York", label: "America/New_York (US Eastern)" },
  { value: "America/Chicago", label: "America/Chicago (US Central)" },
  { value: "America/Denver", label: "America/Denver (US Mountain)" },
  { value: "America/Los_Angeles", label: "America/Los_Angeles (US Pacific)" },
  { value: "Europe/London", label: "Europe/London (GMT/BST)" },
  { value: "Europe/Paris", label: "Europe/Paris (CET/CEST)" },
  { value: "Europe/Berlin", label: "Europe/Berlin (CET/CEST)" },
  { value: "Asia/Dubai", label: "Asia/Dubai (GST)" },
  { value: "Asia/Shanghai", label: "Asia/Shanghai (CST)" },
  { value: "Asia/Singapore", label: "Asia/Singapore (SGT)" },
  { value: "Asia/Tokyo", label: "Asia/Tokyo (JST)" },
  { value: "Australia/Sydney", label: "Australia/Sydney (AEST/AEDT)" },
];

export function CommitScheduleCard({
  timezone,
  onTimezoneChange,
  isScheduleActive,
  onToggleSchedule,
  weekdayCutoff = "08:00 AM",
}: CommitScheduleCardProps) {
  return (
    <div className="rounded-xl bg-card border border-border/50 p-6 shadow-sm space-y-6">
      <div className="flex items-center justify-between border-b border-border/40 pb-4">
        <div className="flex items-center gap-2.5">
          <Clock className="h-5 w-5 text-primary" />
          <div>
            <h2 className="text-lg font-heading font-semibold text-foreground">
              Commit Schedule &amp; Timezone Restriction
            </h2>
            <p className="text-xs text-muted-foreground">
              Enforce weekday morning commit completion before cutoff in target timezone.
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="border-emerald-500/40 text-emerald-400 bg-emerald-500/10">
            Weekend Unconstrained
          </Badge>
          <span className="text-xs px-2.5 py-1 rounded-full bg-muted text-muted-foreground border border-border/60 font-medium">
            Schedule Engine
          </span>
        </div>
      </div>

      <div className="grid sm:grid-cols-2 gap-6">
        {/* Timezone Selection */}
        <div className="space-y-3">
          <div className="flex items-center gap-2">
            <Globe className="h-4 w-4 text-muted-foreground" />
            <Label htmlFor="commit-timezone-select" className="text-xs font-semibold text-foreground uppercase tracking-wider">
              Target Commit Timezone
            </Label>
          </div>
          <Select value={timezone || "Local"} onValueChange={onTimezoneChange}>
            <SelectTrigger id="commit-timezone-select" className="w-full bg-muted/20 border-border/60 font-mono text-xs">
              <SelectValue placeholder="Select target timezone" />
            </SelectTrigger>
            <SelectContent>
              {COMMON_TIMEZONES.map((tz) => (
                <SelectItem key={tz.value} value={tz.value} className="text-xs font-mono">
                  {tz.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <p className="text-xs text-muted-foreground">
            IANA timezone identifier used to evaluate local weekday and cutoff boundaries.
          </p>
        </div>

        {/* Cutoff Time & Schedule Toggle */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <Label htmlFor="schedule-toggle" className="text-xs font-semibold text-foreground uppercase tracking-wider">
                Weekday Schedule Enforcement
              </Label>
              <p className="text-xs text-muted-foreground">
                Restrict commits to before {weekdayCutoff} Mon&ndash;Fri
              </p>
            </div>
            <Switch
              id="schedule-toggle"
              checked={isScheduleActive}
              onCheckedChange={onToggleSchedule}
            />
          </div>

          <div className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-border/40">
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-primary" />
              <span className="text-xs font-medium text-foreground">Weekday Cutoff Time</span>
            </div>
            <span className="text-xs font-mono font-bold px-2.5 py-1 rounded bg-primary/10 text-primary border border-primary/20">
              {weekdayCutoff}
            </span>
          </div>
        </div>
      </div>

      {/* Policy Details / Rules */}
      <div className="rounded-lg bg-muted/20 border border-border/40 p-3.5 flex items-start gap-3 text-xs text-muted-foreground">
        <ShieldCheck className="h-4 w-4 text-primary shrink-0 mt-0.5" />
        <div className="space-y-1">
          <span className="font-semibold text-foreground">Schedule Guardrails:</span>
          <div>
            Weekday commits under <code className="text-primary font-mono font-semibold">--schedule</code> must be finished prior to <strong className="text-foreground">{weekdayCutoff}</strong> in the configured timezone. On Saturday and Sunday, commits are completely unconstrained around the clock.
          </div>
        </div>
      </div>
    </div>
  );
}

export default CommitScheduleCard;
