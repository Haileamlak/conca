import { Link, useLocation } from "react-router-dom";
import { 
  LayoutDashboard, Layers, FileText, BarChart3, Settings, 
  Zap, Activity, ChevronRight, Bot
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Dashboard", path: "/" },
  { icon: Layers, label: "Brands", path: "/brands" },
  { icon: FileText, label: "Posts", path: "/posts" },
  { icon: BarChart3, label: "Analytics", path: "/analytics" },
  { icon: Settings, label: "Settings", path: "/settings" },
];

export function Sidebar() {
  const location = useLocation();

  return (
    <aside className="w-60 flex-shrink-0 flex flex-col" style={{ background: "hsl(var(--sidebar-background))", borderRight: "1px solid hsl(var(--sidebar-border))" }}>
      {/* Logo */}
      <div className="p-5 border-b" style={{ borderColor: "hsl(var(--sidebar-border))" }}>
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: "linear-gradient(135deg, hsl(var(--primary)), hsl(190 100% 30%))", boxShadow: "0 0 16px hsl(183 100% 45% / 0.4)" }}>
            <Bot className="w-4 h-4" style={{ color: "hsl(var(--primary-foreground))" }} />
          </div>
          <div>
            <span className="font-bold text-base tracking-wide" style={{ color: "hsl(var(--sidebar-accent-foreground))" }}>Conca</span>
            <p className="text-xs font-mono" style={{ color: "hsl(var(--primary))" }}>v1.0 · Active</p>
          </div>
        </div>
      </div>

      {/* Agent status */}
      <div className="px-4 py-3 mx-3 mt-3 rounded-lg" style={{ background: "hsl(152 76% 40% / 0.1)", border: "1px solid hsl(152 76% 40% / 0.25)" }}>
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 rounded-full status-pulse" style={{ background: "hsl(var(--success))" }} />
          <span className="text-xs font-mono" style={{ color: "hsl(var(--success))" }}>Agent Loop · Running</span>
        </div>
        <p className="text-xs mt-1" style={{ color: "hsl(var(--muted-foreground))" }}>Next cycle in 2h 14m</p>
      </div>

      {/* Nav */}
      <nav className="flex-1 p-3 space-y-1 mt-2">
        {navItems.map((item) => {
          const isActive = location.pathname === item.path;
          return (
            <Link key={item.path} to={item.path} className={`sidebar-item ${isActive ? "active" : ""}`}>
              <item.icon className="w-4 h-4 flex-shrink-0" />
              <span>{item.label}</span>
              {isActive && <ChevronRight className="w-3 h-3 ml-auto" />}
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="p-4 border-t" style={{ borderColor: "hsl(var(--sidebar-border))" }}>
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold" style={{ background: "hsl(var(--primary) / 0.2)", color: "hsl(var(--primary))" }}>
            A
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-xs font-medium truncate" style={{ color: "hsl(var(--sidebar-accent-foreground))" }}>Admin</p>
            <p className="text-xs truncate" style={{ color: "hsl(var(--muted-foreground))" }}>admin@conca.ai</p>
          </div>
          <Activity className="w-3.5 h-3.5 flex-shrink-0" style={{ color: "hsl(var(--primary))" }} />
        </div>
      </div>
    </aside>
  );
}

interface LayoutProps {
  children: React.ReactNode;
  title: string;
  subtitle?: string;
  actions?: React.ReactNode;
}

export function Layout({ children, title, subtitle, actions }: LayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden" style={{ background: "hsl(var(--background))" }}>
      <Sidebar />
      <main className="flex-1 flex flex-col overflow-hidden">
        {/* Top bar */}
        <header className="flex items-center justify-between px-8 py-5 flex-shrink-0" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
          <div>
            <h1 className="text-xl font-semibold" style={{ color: "hsl(var(--foreground))" }}>{title}</h1>
            {subtitle && <p className="text-sm mt-0.5" style={{ color: "hsl(var(--muted-foreground))" }}>{subtitle}</p>}
          </div>
          {actions && <div className="flex items-center gap-3">{actions}</div>}
        </header>
        {/* Content */}
        <div className="flex-1 overflow-y-auto p-8">
          {children}
        </div>
      </main>
    </div>
  );
}
