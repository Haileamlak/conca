import { useState, useEffect } from "react";
import { Layout } from "@/components/Layout";
import { StatusBadge } from "@/components/StatusBadge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/components/ui/use-toast";
import { Plus, Twitter, Linkedin, Play, RefreshCw, X, Layers } from "lucide-react";

interface Brand {
  id: string;
  name: string;
  industry: string;
  voice: string;
  target_audience: string;
  topics: string[];
  anti_topics: string[];
  schedule_interval_hours: number;
  status?: string;
}

function BrandCard({ brand, onRun, onSync }: { brand: Brand; onRun: () => void; onSync: () => void }) {
  return (
    <div className="card-glass rounded-xl p-5 hover:border-primary/30 transition-all cursor-pointer group animate-fade-in" style={{ border: "1px solid hsl(var(--border))" }}>
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl flex items-center justify-center text-sm font-bold" style={{ background: "hsl(var(--primary) / 0.15)", color: "hsl(var(--primary))" }}>
            {brand.name[0]}
          </div>
          <div>
            <h3 className="font-semibold text-sm">{brand.name}</h3>
            <p className="text-xs font-mono mt-0.5" style={{ color: "hsl(var(--muted-foreground))" }}>{brand.industry} · {brand.voice}</p>
          </div>
        </div>
        <StatusBadge status={brand.status as any || "idle"} />
      </div>

      <p className="text-xs mb-3" style={{ color: "hsl(var(--muted-foreground))" }}>
        <span style={{ color: "hsl(var(--foreground))" }}>Audience:</span> {brand.target_audience}
      </p>

      <div className="flex flex-wrap gap-1.5 mb-4">
        {brand.topics?.map((t) => (
          <span key={t} className="px-2 py-0.5 rounded text-xs font-mono" style={{ background: "hsl(var(--primary) / 0.1)", color: "hsl(var(--primary))" }}>{t}</span>
        ))}
        {brand.anti_topics?.map((t) => (
          <span key={t} className="px-2 py-0.5 rounded text-xs font-mono" style={{ background: "hsl(var(--destructive) / 0.1)", color: "hsl(var(--destructive))" }}>✕ {t}</span>
        ))}
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex gap-1.5">
            <Twitter className="w-3.5 h-3.5" style={{ color: "hsl(200 100% 50%)" }} />
            <Linkedin className="w-3.5 h-3.5" style={{ color: "hsl(210 80% 55%)" }} />
          </div>
          <span className="text-xs font-mono" style={{ color: "hsl(var(--muted-foreground))" }}>Every {brand.schedule_interval_hours}h</span>
        </div>
      </div>

      <div className="flex gap-2 mt-4 pt-4" style={{ borderTop: "1px solid hsl(var(--border))" }}>
        <Button size="sm" variant="outline" className="flex-1 gap-1.5 text-xs h-8 border-border" onClick={onRun}>
          <Play className="w-3 h-3" /> Run Now
        </Button>
        <Button size="sm" variant="outline" className="flex-1 gap-1.5 text-xs h-8 border-border" onClick={onSync}>
          <RefreshCw className="w-3 h-3" /> Sync
        </Button>
      </div>
    </div>
  );
}

function CreateBrandModal({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
  const [formData, setFormData] = useState({
    id: "",
    name: "",
    industry: "",
    voice: "",
    target_audience: "",
    topics: "",
    anti_topics: "",
    schedule_interval_hours: 4
  });
  const { toast } = useToast();

  const handleCreate = async () => {
    try {
      const res = await fetch("/api/brands", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          ...formData,
          topics: formData.topics.split(",").map(s => s.trim()).filter(s => s),
          anti_topics: formData.anti_topics.split(",").map(s => s.trim()).filter(s => s),
          schedule_interval_hours: Number(formData.schedule_interval_hours)
        })
      });

      if (res.ok) {
        toast({ title: "Brand Created", description: `Profile for ${formData.name} is now active.` });
        onSuccess();
        onClose();
      } else {
        toast({ title: "Error", description: "Failed to create brand profile.", variant: "destructive" });
      }
    } catch (error) {
      toast({ title: "Error", description: "Network error.", variant: "destructive" });
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center font-sans" style={{ background: "hsl(var(--background) / 0.8)", backdropFilter: "blur(4px)" }}>
      <div className="w-full max-w-lg rounded-2xl p-6 animate-fade-in shadow-2xl" style={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))" }}>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-base font-semibold">Create Brand Identity</h2>
          <Button variant="ghost" size="sm" className="w-7 h-7 p-0" onClick={onClose}><X className="w-4 h-4" /></Button>
        </div>
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Brand ID *</Label>
              <Input placeholder="tech_startup" className="font-mono text-sm h-9 bg-secondary border-border" value={formData.id} onChange={e => setFormData({ ...formData, id: e.target.value })} />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Name *</Label>
              <Input placeholder="TechVision" className="text-sm h-9 bg-secondary border-border" value={formData.name} onChange={e => setFormData({ ...formData, name: e.target.value })} />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Industry</Label>
              <Input placeholder="Technology" className="text-sm h-9 bg-secondary border-border" value={formData.industry} onChange={e => setFormData({ ...formData, industry: e.target.value })} />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Brand Voice</Label>
              <Input placeholder="Professional, Witty..." className="text-sm h-9 bg-secondary border-border" value={formData.voice} onChange={e => setFormData({ ...formData, voice: e.target.value })} />
            </div>
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Target Audience</Label>
            <Input placeholder="CTOs, Senior Engineers" className="text-sm h-9 bg-secondary border-border" value={formData.target_audience} onChange={e => setFormData({ ...formData, target_audience: e.target.value })} />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Topics (comma-sep)</Label>
              <Input placeholder="AI, Cloud, DevOps" className="text-sm h-9 bg-secondary border-border" value={formData.topics} onChange={e => setFormData({ ...formData, topics: e.target.value })} />
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Anti-Topics</Label>
              <Input placeholder="Politics, Sports" className="text-sm h-9 bg-secondary border-border" value={formData.anti_topics} onChange={e => setFormData({ ...formData, anti_topics: e.target.value })} />
            </div>
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs font-mono uppercase tracking-wide" style={{ color: "hsl(var(--muted-foreground))" }}>Schedule Interval (hours)</Label>
            <Input type="number" placeholder="4" value={formData.schedule_interval_hours} onChange={e => setFormData({ ...formData, schedule_interval_hours: Number(e.target.value) })} className="font-mono text-sm h-9 bg-secondary border-border" />
          </div>
        </div>
        <div className="flex gap-3 mt-6">
          <Button variant="outline" className="flex-1 border-border" onClick={onClose}>Cancel</Button>
          <Button className="flex-1" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }} onClick={handleCreate}>
            Create Brand Profile
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function Brands() {
  const [brands, setBrands] = useState<Brand[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();

  useEffect(() => {
    fetchBrands();
  }, []);

  const fetchBrands = async () => {
    try {
      const res = await fetch("/api/brands");
      const data = await res.json();
      setBrands(data || []);
    } catch (error) {
      console.error("Failed to fetch brands", error);
    } finally {
      setLoading(false);
    }
  };

  const handleRun = async (id: string) => {
    toast({ title: "Agent Run Triggered", description: "Generation cycle started in background." });
    await fetch(`/api/brands/${id}/run`, { method: "POST" });
  };

  const handleSync = async (id: string) => {
    toast({ title: "Sync Triggered", description: "Updating performance metrics..." });
    await fetch(`/api/brands/${id}/sync`, { method: "POST" });
  };

  return (
    <>
      <Layout
        title="Brand Identities"
        subtitle={`${brands.length} active agent profiles`}
        actions={
          <Button size="sm" className="gap-2 text-xs" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }} onClick={() => setShowCreate(true)}>
            <Plus className="w-3.5 h-3.5" /> New Brand
          </Button>
        }
      >
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-5">
          {brands.map((brand) => (
            <BrandCard key={brand.id} brand={brand} onRun={() => handleRun(brand.id)} onSync={() => handleSync(brand.id)} />
          ))}

          {!loading && (
            <button
              onClick={() => setShowCreate(true)}
              className="rounded-xl p-5 border-2 border-dashed flex flex-col items-center justify-center gap-3 min-h-[200px] transition-all hover:border-primary/40 hover:bg-primary/5 group"
              style={{ borderColor: "hsl(var(--border))", color: "hsl(var(--muted-foreground))" }}
            >
              <div className="w-10 h-10 rounded-full flex items-center justify-center transition-transform group-hover:scale-110" style={{ background: "hsl(var(--secondary))" }}>
                <Plus className="w-5 h-5" />
              </div>
              <p className="text-sm font-medium">Add New Brand</p>
            </button>
          )}
        </div>

        {/* API Reference */}
        <div className="mt-8 card-glass rounded-xl p-5">
          <h3 className="text-sm font-semibold mb-4 flex items-center gap-2">
            <Layers className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
            Operational Endpoints
          </h3>
          <div className="space-y-2 font-mono text-[10px] sm:text-xs">
            {[
              { method: "POST", path: "/api/brands", desc: "Define agent identity" },
              { method: "POST", path: "/api/brands/{id}/run", desc: "Force generation cycle" },
              { method: "POST", path: "/api/brands/{id}/sync", desc: "Refresh social metrics" },
              { method: "GET", path: "/api/brands/{id}/posts", desc: "View brand history" },
            ].map((ep) => (
              <div key={ep.path} className="flex items-center gap-3 py-2 px-3 rounded-lg" style={{ background: "hsl(var(--secondary))" }}>
                <span className="w-14 text-center rounded px-1 py-0.5 text-[10px] sm:text-xs" style={{ background: "hsl(var(--primary) / 0.15)", color: "hsl(var(--primary))" }}>{ep.method}</span>
                <span className="truncate" style={{ color: "hsl(var(--foreground))" }}>{ep.path}</span>
                <span style={{ color: "hsl(var(--muted-foreground))" }} className="ml-auto hidden sm:inline">{ep.desc}</span>
              </div>
            ))}
          </div>
        </div>
      </Layout>

      {showCreate && <CreateBrandModal onClose={() => setShowCreate(false)} onSuccess={fetchBrands} />}
    </>
  );
}
