import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type StatusType = "published" | "draft" | "pending_review" | "approved" | "scheduled" | "failed" | "running" | "idle";

const statusConfig: Record<StatusType, { label: string; className: string }> = {
  published: { label: "Published", className: "bg-success/15 text-success border-success/30" },
  draft: { label: "Draft", className: "bg-muted text-muted-foreground border-border" },
  pending_review: { label: "Pending Review", className: "bg-warning/15 text-warning border-warning/30" },
  approved: { label: "Approved", className: "bg-primary/15 text-primary border-primary/30" },
  scheduled: { label: "Scheduled", className: "bg-blue-500/15 text-blue-400 border-blue-500/30" },
  failed: { label: "Failed", className: "bg-destructive/15 text-destructive border-destructive/30" },
  running: { label: "Running", className: "bg-primary/15 text-primary border-primary/30" },
  idle: { label: "Idle", className: "bg-muted text-muted-foreground border-border" },
};

export function StatusBadge({ status }: { status: string }) {
  const config = statusConfig[status as StatusType] ?? { label: status, className: "bg-muted text-muted-foreground border-border" };
  return (
    <Badge variant="outline" className={cn("text-xs font-mono border", config.className)}>
      {status === "running" && (
        <span className="w-1.5 h-1.5 rounded-full bg-current mr-1.5 status-pulse inline-block" />
      )}
      {config.label}
    </Badge>
  );
}
