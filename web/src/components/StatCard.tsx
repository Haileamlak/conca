import { cn } from "@/lib/utils";

interface StatCardProps {
  label: string;
  value: string | number;
  sub?: string;
  icon?: React.ReactNode;
  trend?: { value: string; up: boolean };
  glow?: "primary" | "success" | "warning";
  className?: string;
}

export function StatCard({ label, value, sub, icon, trend, glow, className }: StatCardProps) {
  const glowClass = glow === "primary" ? "shadow-glow-primary" : glow === "success" ? "shadow-glow-success" : glow === "warning" ? "shadow-glow-warning" : "";

  return (
    <div className={cn("card-glass rounded-xl p-5 animate-fade-in", glowClass, className)}>
      <div className="flex items-start justify-between mb-3">
        <p className="text-xs font-mono uppercase tracking-widest" style={{ color: "hsl(var(--muted-foreground))" }}>{label}</p>
        {icon && (
          <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: "hsl(var(--primary) / 0.1)" }}>
            <span style={{ color: "hsl(var(--primary))" }}>{icon}</span>
          </div>
        )}
      </div>
      <p className="text-3xl font-bold font-mono animate-count-up" style={{ color: "hsl(var(--foreground))" }}>{value}</p>
      <div className="flex items-center gap-3 mt-2">
        {sub && <p className="text-xs" style={{ color: "hsl(var(--muted-foreground))" }}>{sub}</p>}
        {trend && (
          <span className={cn("text-xs font-mono", trend.up ? "text-success" : "text-destructive")}>
            {trend.up ? "↑" : "↓"} {trend.value}
          </span>
        )}
      </div>
    </div>
  );
}
