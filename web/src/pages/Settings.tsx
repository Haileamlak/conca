import { Layout } from "@/components/Layout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Database, Key, Globe, Brain, BarChart2, CheckCircle2, Circle } from "lucide-react";

const integrations = [
  { name: "Gemini API", desc: "LLM + Embeddings for content generation and semantic memory", icon: Brain, status: "connected", key: "AIza••••••••••••••••••••Xk9" },
  { name: "NewsAPI", desc: "Real-time news and trend discovery", icon: Globe, status: "connected", key: "••••••••••••••••••••a2f1" },
  { name: "Twitter / X API v2", desc: "Post publishing and analytics sync", icon: Globe, status: "connected", key: "••••••••••••••••••••AAAB" },
  { name: "LinkedIn API", desc: "ugcPosts publishing endpoint", icon: Globe, status: "disconnected", key: "" },
  { name: "NewsData.io", desc: "Secondary trend research source", icon: Globe, status: "disconnected", key: "" },
  { name: "PostgreSQL", desc: "Primary relational store (optional, falls back to JSON)", icon: Database, status: "disconnected", key: "" },
];

export default function Settings() {
  return (
    <Layout title="Settings" subtitle="Configure integrations, API keys, and agent behavior">
      <div className="max-w-3xl space-y-6">
        {/* API Config */}
        <div className="card-glass rounded-xl overflow-hidden">
          <div className="px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Key className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
              API Server Configuration
            </h3>
          </div>
          <div className="p-5 grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Server Port</Label>
              <Input defaultValue="8080" className="font-mono text-sm h-9 bg-secondary border-border" />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>JWT Secret</Label>
              <Input type="password" defaultValue="supersecret" className="font-mono text-sm h-9 bg-secondary border-border" />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Data Directory</Label>
              <Input defaultValue="./data" className="font-mono text-sm h-9 bg-secondary border-border" />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Default Interval (hours)</Label>
              <Input type="number" defaultValue={4} className="font-mono text-sm h-9 bg-secondary border-border" />
            </div>
          </div>
          <div className="px-5 pb-5">
            <Button className="text-xs h-8 px-4" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }}>
              Save Configuration
            </Button>
          </div>
        </div>

        {/* Integrations */}
        <div className="card-glass rounded-xl overflow-hidden">
          <div className="px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Globe className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
              Integrations
            </h3>
          </div>
          <div className="divide-y" style={{ borderColor: "hsl(var(--border))" }}>
            {integrations.map((intg) => (
              <div key={intg.name} className="flex items-center gap-4 px-5 py-4">
                <div className="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0" style={{ background: "hsl(var(--secondary))" }}>
                  <intg.icon className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-0.5">
                    <p className="text-sm font-medium">{intg.name}</p>
                    {intg.status === "connected"
                      ? <CheckCircle2 className="w-3.5 h-3.5" style={{ color: "hsl(var(--success))" }} />
                      : <Circle className="w-3.5 h-3.5" style={{ color: "hsl(var(--muted-foreground))" }} />
                    }
                  </div>
                  <p className="text-xs" style={{ color: "hsl(var(--muted-foreground))" }}>{intg.desc}</p>
                  {intg.key && <p className="text-xs font-mono mt-1" style={{ color: "hsl(var(--primary) / 0.7)" }}>{intg.key}</p>}
                </div>
                {intg.status === "connected"
                  ? <Button size="sm" variant="outline" className="text-xs h-7 px-3 border-border">Update Key</Button>
                  : <Button size="sm" className="text-xs h-7 px-3" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }}>Connect</Button>
                }
              </div>
            ))}
          </div>
        </div>

        {/* Agent settings */}
        <div className="card-glass rounded-xl overflow-hidden">
          <div className="px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Brain className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
              Agent Behavior
            </h3>
          </div>
          <div className="p-5 grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Content Model</Label>
              <select className="w-full text-xs font-mono px-3 py-2 rounded-lg border h-9" style={{ background: "hsl(var(--secondary))", borderColor: "hsl(var(--border))", color: "hsl(var(--foreground))" }}>
                <option>gemini-1.5-pro</option>
                <option>gemini-1.5-flash</option>
                <option>gemini-2.0-pro</option>
              </select>
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Embedding Model</Label>
              <select className="w-full text-xs font-mono px-3 py-2 rounded-lg border h-9" style={{ background: "hsl(var(--secondary))", borderColor: "hsl(var(--border))", color: "hsl(var(--foreground))" }}>
                <option>text-embedding-004</option>
                <option>text-embedding-005</option>
              </select>
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Max Trends Per Cycle</Label>
              <Input type="number" defaultValue={10} className="font-mono text-sm h-9 bg-secondary border-border" />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Content Variants Per Trend</Label>
              <Input type="number" defaultValue={3} className="font-mono text-sm h-9 bg-secondary border-border" />
            </div>
          </div>
          <div className="px-5 pb-5">
            <Button className="text-xs h-8 px-4" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }}>
              Save Agent Config
            </Button>
          </div>
        </div>

        {/* CLI Quick Reference */}
        <div className="card-glass rounded-xl p-5">
          <h3 className="text-sm font-semibold mb-4 flex items-center gap-2">
            <BarChart2 className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
            CLI Quick Reference
          </h3>
          <div className="space-y-2">
            {[
              { cmd: "go run cmd/server/main.go --port 8080", desc: "Start SaaS API server" },
              { cmd: "go run cmd/main.go --config config/brand.json --daemon --interval 4h", desc: "Run daemon mode" },
              { cmd: "go run cmd/main.go --sync", desc: "Sync analytics from all platforms" },
            ].map((c) => (
              <div key={c.cmd} className="p-3 rounded-lg" style={{ background: "hsl(var(--secondary))" }}>
                <p className="text-xs font-mono" style={{ color: "hsl(var(--primary))" }}>$ {c.cmd}</p>
                <p className="text-xs mt-0.5" style={{ color: "hsl(var(--muted-foreground))" }}>{c.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </Layout>
  );
}
