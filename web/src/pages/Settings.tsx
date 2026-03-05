import { useState, useEffect } from "react";
import { Layout } from "@/components/Layout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { User, LogOut, Bell, Shield, CreditCard, Sparkles } from "lucide-react";
import { api } from "@/lib/api";

export default function Settings() {
  const [user, setUser] = useState<{ id: string, email: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await api.get("/auth/me");
        if (res.ok) {
          const data = await res.json();
          setUser(data);
        }
      } catch (err) {
        console.error("Failed to fetch user data", err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  return (
    <Layout title="Settings" subtitle="Manage your account preferences and profile">
      <div className="max-w-3xl space-y-6">
        {/* Profile Section */}
        <div className="card-glass rounded-xl overflow-hidden animate-fade-in">
          <div className="px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <User className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
              Profile Settings
            </h3>
          </div>
          <div className="p-6">
            <div className="flex items-center gap-6 mb-8">
              <div className="w-20 h-20 rounded-2xl flex items-center justify-center text-3xl font-bold relative group overflow-hidden" style={{ background: "hsl(var(--primary) / 0.1)", color: "hsl(var(--primary))", border: "1px solid hsl(var(--primary) / 0.2)" }}>
                {user?.email?.[0].toUpperCase() || "A"}
                <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center cursor-pointer">
                  <p className="text-[10px] text-white font-bold uppercase tracking-widest">Change</p>
                </div>
              </div>
              <div className="flex-1">
                <h4 className="text-lg font-bold text-foreground mb-1">{user?.email?.split('@')[0] || "User Name"}</h4>
                <p className="text-sm text-muted-foreground mb-3">{user?.email}</p>
                <div className="flex gap-2">
                  <div className="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider bg-primary/10 text-primary border border-primary/20">
                    Pro Plan
                  </div>
                  <div className="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider bg-secondary text-muted-foreground border border-border">
                    Verified
                  </div>
                </div>
              </div>
              <Button variant="outline" size="sm" className="gap-2 border-border text-xs rounded-lg px-4" onClick={() => { localStorage.removeItem("conca_token"); window.location.href = "/auth"; }}>
                <LogOut className="w-3.5 h-3.5" /> Sign Out
              </Button>
            </div>

            <div className="grid grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label className="text-xs font-semibold text-muted-foreground">Email Address</Label>
                <Input value={user?.email || ""} readOnly className="bg-secondary/50 border-border h-10" />
              </div>
              <div className="space-y-2">
                <Label className="text-xs font-semibold text-muted-foreground">Account Status</Label>
                <Input value="Active Premium" readOnly className="bg-secondary/50 border-border h-10 text-success font-bold" />
              </div>
            </div>
          </div>
        </div>

        {/* Global Preferences */}
        <div className="card-glass rounded-xl overflow-hidden">
          <div className="px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Sparkles className="w-4 h-4" style={{ color: "hsl(var(--primary))" }} />
              Content Processing
            </h3>
          </div>
          <div className="p-6 space-y-4">
            {[
              { icon: Bell, title: "Smart Notifications", desc: "Get alerted when top trends align with your brand voice" },
              { icon: Shield, title: "Brand Protection", desc: "Rigorous filtering for sensitive topics and controversial trends" },
              { icon: CreditCard, title: "Usage & Billing", desc: "Manage your subscription plan and view monthly usage" }
            ].map((item) => (
              <div key={item.title} className="flex items-center gap-4 p-4 rounded-xl hover:bg-secondary/30 transition-colors cursor-pointer border border-transparent hover:border-border">
                <div className="w-10 h-10 rounded-lg flex items-center justify-center bg-secondary">
                  <item.icon className="w-5 h-5 text-muted-foreground" />
                </div>
                <div className="flex-1">
                  <p className="text-sm font-bold">{item.title}</p>
                  <p className="text-xs text-muted-foreground">{item.desc}</p>
                </div>
                <div className="w-10 h-6 rounded-full bg-primary/20 p-1 flex items-center justify-end">
                  <div className="w-4 h-4 rounded-full bg-primary" />
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Newsletter / Updates */}
        <div className="card-glass rounded-xl p-6 relative overflow-hidden">
          <div className="absolute top-0 right-0 p-8 opacity-10 pointer-events-none">
            <Sparkles className="w-24 h-24 text-primary" />
          </div>
          <h3 className="text-sm font-bold mb-2">Conca AI Intelligence</h3>
          <p className="text-xs text-muted-foreground mb-4 max-w-md leading-relaxed">
            Get weekly reports on how our AI is evolving to better serve your brand goals. We constantly update our underlying semantic models to improve trend accuracy by 40%+.
          </p>
          <Button variant="outline" className="text-xs h-9 px-6 border-primary/30 text-primary hover:bg-primary/5">
            View Release Notes
          </Button>
        </div>
      </div>
    </Layout>
  );
}
