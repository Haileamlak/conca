import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Layout } from "@/components/Layout";
import { StatCard } from "@/components/StatCard";
import { StatusBadge } from "@/components/StatusBadge";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/use-toast";
import {
  Layers, FileText, TrendingUp, Zap, Play, RefreshCw,
  Twitter, Linkedin, Brain, Search, BarChart2, Plus
} from "lucide-react";
import { api } from "@/lib/api";

interface Post {
  id: string;
  brand_id: string;
  topic: string;
  platform: string;
  status: string;
  analytics: { likes: number; shares: number };
  created_at: string;
}

interface BrandPerf {
  brand_id: string;
  brand_name: string;
  post_count: number;
  score: number;
}

export default function Dashboard() {
  const navigate = useNavigate();
  const [posts, setPosts] = useState<Post[]>([]);
  const [performance, setPerformance] = useState<BrandPerf[]>([]);
  const [stats, setStats] = useState({ impressions: 0, likes: 0, shares: 0 });
  const [logs, setLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();

  useEffect(() => {
    fetchDashboardData();
    fetchLogs();
    const interval = setInterval(() => {
      fetchDashboardData();
      fetchLogs();
    }, 10000);
    return () => clearInterval(interval);
  }, []);

  const fetchLogs = async () => {
    try {
      const res = await api.get("/logs");
      const data = await res.json();
      setLogs(data || []);
    } catch (err) {
      console.error("Failed to fetch logs", err);
    }
  };

  const fetchDashboardData = async () => {
    try {
      // Fetch global analytics (stats + brand performance)
      const analyticsRes = await api.get("/analytics");
      const analyticsData = await analyticsRes.json();
      if (analyticsData.global) {
        setStats({
          impressions: analyticsData.global.total_impressions,
          likes: analyticsData.global.total_likes,
          shares: analyticsData.global.total_shares,
        });
      }
      if (analyticsData && analyticsData.performance && Array.isArray(analyticsData.performance)) {
        setPerformance(analyticsData.performance);
      } else {
        setPerformance([]);
      }

      // Fetch recent posts
      const postsRes = await api.get("/posts");
      const postsData = await postsRes.json();
      if (Array.isArray(postsData)) {
        setPosts(postsData.slice(0, 5));
      } else {
        setPosts([]);
      }
    } catch (error) {
      console.error("Failed to fetch dashboard data", error);
    } finally {
      setLoading(false);
    }
  };

  const handleRunAll = async () => {
    toast({ title: "Agent Run Triggered", description: "Starting autonomous cycle for all brands..." });
    for (const b of performance) {
      await api.post(`/brands/${b.brand_id}/run`);
    }
  };

  const handleSyncAll = async () => {
    toast({ title: "Sync Triggered", description: "Fetching latest analytics from platforms..." });
    for (const b of performance) {
      await api.post(`/brands/${b.brand_id}/sync`);
    }
    setTimeout(fetchDashboardData, 2000);
  };

  const agentActivity = [
    { time: "Just now", event: "Dashboard updated", detail: `${posts.length} recent posts loaded`, type: "sync" },
    { time: "Today", event: "Trend research active", detail: "Monitoring NewsAPI + DuckDuckGo", type: "research" },
  ];

  return (
    <Layout
      title="Mission Control"
      subtitle="Autonomous agent dashboard · Real-time"
      actions={
        <div className="flex gap-2">
          <Button variant="outline" size="sm" className="gap-2 font-mono text-xs border-border" onClick={handleSyncAll}>
            <RefreshCw className="w-3.5 h-3.5" /> Sync All
          </Button>
          <Button size="sm" className="gap-2 font-mono text-xs" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }} onClick={handleRunAll}>
            <Play className="w-3.5 h-3.5" /> Run Agent (All)
          </Button>
        </div>
      }
    >
      {/* Welcome / Empty State for New Users */}
      {!loading && performance.length === 0 && (
        <div className="card-glass rounded-2xl p-8 mb-8 flex flex-col items-center text-center animate-fade-in" style={{ background: "linear-gradient(135deg, hsl(var(--primary) / 0.05), hsl(var(--secondary) / 0.1))", border: "1px solid hsl(var(--primary) / 0.2)" }}>
          <div className="w-16 h-16 rounded-2xl flex items-center justify-center mb-6" style={{ background: "hsl(var(--primary) / 0.15)", color: "hsl(var(--primary))", boxShadow: "0 0 20px hsl(var(--primary) / 0.2)" }}>
            <Layers className="w-8 h-8" />
          </div>
          <h2 className="text-2xl font-bold mb-2">Welcome to your Autonomous Content Hub</h2>
          <p className="text-sm text-muted-foreground max-w-md mb-8 leading-relaxed">
            You haven't created any brand identities yet. To start generating intelligent content and monitoring trends, define your first brand.
          </p>
          <div className="flex gap-4">
            <Button className="gap-2" style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }} onClick={() => navigate("/brands")}>
              <Plus className="w-4 h-4" /> Create First Brand
            </Button>
            <Button variant="outline" className="border-border" onClick={() => navigate("/settings")}>
              View Integrations
            </Button>
          </div>
        </div>
      )}

      {/* Stats */}
      {!loading && performance.length > 0 && (
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <StatCard label="Total Posts" value={performance.reduce((acc, b) => acc + b.post_count, 0).toString()} sub="All brands" icon={<FileText className="w-4 h-4" />} glow="primary" />
          <StatCard label="Active Brands" value={performance.length.toString()} sub="Personal Portfolio" icon={<Layers className="w-4 h-4" />} />
          <StatCard label="Total Impressions" value={stats.impressions.toLocaleString()} sub="Lifetime Reach" icon={<TrendingUp className="w-4 h-4" />} glow="success" />
          <StatCard label="Total Engagement" value={(stats.likes + stats.shares).toLocaleString()} sub="Likes + Shares" icon={<Zap className="w-4 h-4" />} />
        </div>
      )}

      <div className="grid grid-cols-12 gap-6">
        {/* Recent Posts */}
        <div className="col-span-7">
          <div className="card-glass rounded-xl overflow-hidden">
            <div className="flex items-center justify-between px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
              <h2 className="text-sm font-semibold">Recent Posts</h2>
              <a href="/posts" className="text-xs font-mono" style={{ color: "hsl(var(--primary))" }}>View all →</a>
            </div>
            <div className="divide-y" style={{ borderColor: "hsl(var(--border))" }}>
              {posts.map((post) => (
                <div key={post.id} className="flex items-center gap-4 px-5 py-3.5 hover:bg-secondary/40 transition-colors">
                  <div className="w-7 h-7 rounded-md flex items-center justify-center flex-shrink-0" style={{ background: post.platform === "twitter" ? "hsl(200 100% 50% / 0.1)" : "hsl(210 80% 55% / 0.1)" }}>
                    {post.platform === "twitter"
                      ? <Twitter className="w-3.5 h-3.5" style={{ color: "hsl(200 100% 50%)" }} />
                      : <Linkedin className="w-3.5 h-3.5" style={{ color: "hsl(210 80% 55%)" }} />
                    }
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{post.topic}</p>
                    <p className="text-xs font-mono mt-0.5" style={{ color: "hsl(var(--muted-foreground))" }}>{post.brand_id.split('_').pop()} · {new Date(post.created_at).toLocaleDateString()}</p>
                  </div>
                  <StatusBadge status={post.status as any} />
                  {post.analytics.likes > 0 && (
                    <div className="text-right hidden sm:block">
                      <p className="text-xs font-mono" style={{ color: "hsl(var(--foreground))" }}>{post.analytics.likes} <span style={{ color: "hsl(var(--muted-foreground))" }}>♥</span></p>
                      <p className="text-xs font-mono" style={{ color: "hsl(var(--muted-foreground))" }}>{post.analytics.shares} ↗</p>
                    </div>
                  )}
                </div>
              ))}
              {posts.length === 0 && (
                <div className="px-5 py-8 text-center text-sm text-muted-foreground">
                  {performance.length > 0
                    ? "No posts found. Start an agent run to generate content!"
                    : "Create a brand and start your first agent run to see history here."}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Right col */}
        <div className="col-span-5 space-y-5">
          {/* Agent Activity Feed */}
          <div className="card-glass rounded-xl overflow-hidden flex flex-col h-[350px]">
            <div className="flex items-center justify-between px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
              <h2 className="text-sm font-semibold flex items-center gap-2">
                <span className="w-2 h-2 rounded-full status-pulse" style={{ background: "hsl(var(--success))" }} />
                Agent Log
              </h2>
              <span className="text-xs font-mono" style={{ color: "hsl(var(--muted-foreground))" }}>Live</span>
            </div>
            <div className="flex-1 overflow-y-auto p-4 font-mono text-[10px] space-y-2">
              {logs.slice().reverse().map((log, i) => (
                <div key={i} className="flex gap-2 animate-fade-in border-b border-border/10 pb-1.5 last:border-0">
                  <span className="text-muted-foreground whitespace-nowrap">[{new Date(log.timestamp).toLocaleTimeString()}]</span>
                  <span className={log.level === "ERROR" ? "text-destructive" : log.level === "WARN" ? "text-warning" : "text-primary"} style={{ minWidth: "40px" }}>
                    {log.level}:
                  </span>
                  <span className="text-foreground/90 leading-normal">{log.message}</span>
                </div>
              ))}
              {logs.length === 0 && <div className="text-muted-foreground">Waiting for agent activity...</div>}
            </div>
          </div>

          {/* Brands Overview */}
          <div className="card-glass rounded-xl overflow-hidden">
            <div className="flex items-center justify-between px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
              <h2 className="text-sm font-semibold">Brand Performance</h2>
              <a href="/brands" className="text-xs font-mono" style={{ color: "hsl(var(--primary))" }}>Manage →</a>
            </div>
            <div className="p-4 space-y-3">
              {performance.map((brand) => (
                <div key={brand.brand_id} className="flex items-center gap-3">
                  <div className="w-7 h-7 rounded-md flex items-center justify-center text-xs font-bold flex-shrink-0" style={{ background: "hsl(var(--primary) / 0.15)", color: "hsl(var(--primary))" }}>
                    {brand.brand_name[0]}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between mb-1">
                      <p className="text-xs font-medium">{brand.brand_name}</p>
                      <span className="text-xs font-mono" style={{ color: "hsl(var(--primary))" }}>{brand.score}</span>
                    </div>
                    <div className="h-1.5 rounded-full overflow-hidden" style={{ background: "hsl(var(--secondary))" }}>
                      <div className="h-full rounded-full" style={{ width: `${brand.score}%`, background: "linear-gradient(90deg, hsl(var(--primary)), hsl(190 100% 35%))" }} />
                    </div>
                  </div>
                  <span className="text-xs font-mono flex-shrink-0" style={{ color: "hsl(var(--muted-foreground))" }}>{brand.post_count}p</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}
