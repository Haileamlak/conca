import { useState, useEffect } from "react";
import { Layout } from "@/components/Layout";
import { StatCard } from "@/components/StatCard";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from "recharts";
import { TrendingUp, Twitter, Linkedin, Heart, Share2, Eye } from "lucide-react";

interface GlobalStats {
  total_impressions: number;
  total_likes: number;
  total_shares: number;
  total_comments: number;
}

interface BrandPerf {
  brand_id: string;
  brand_name: string;
  post_count: number;
  avg_likes: number;
  avg_shares: number;
  score: number;
}

import { api } from "@/lib/api";

export default function Analytics() {
  const [stats, setStats] = useState<GlobalStats>({ total_impressions: 0, total_likes: 0, total_shares: 0, total_comments: 0 });
  const [performance, setPerformance] = useState<BrandPerf[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchAnalytics();
  }, []);

  const fetchAnalytics = async () => {
    try {
      const res = await api.get("/analytics");
      const data = await res.json();
      if (data.global) setStats(data.global);
      if (data.performance) setPerformance(data.performance);
    } catch (error) {
      console.error("Failed to fetch analytics", error);
    } finally {
      setLoading(false);
    }
  };

  // Mock weekly data for the chart until we have historical snapshots in the DB
  const weeklyData = [
    { day: "Mon", impressions: stats.total_impressions * 0.1 },
    { day: "Tue", impressions: stats.total_impressions * 0.15 },
    { day: "Wed", impressions: stats.total_impressions * 0.2 },
    { day: "Thu", impressions: stats.total_impressions * 0.12 },
    { day: "Fri", impressions: stats.total_impressions * 0.25 },
    { day: "Sat", impressions: stats.total_impressions * 0.1 },
    { day: "Sun", impressions: stats.total_impressions * 0.08 },
  ];

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="rounded-lg p-3 text-xs font-mono" style={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))" }}>
          <p className="mb-2" style={{ color: "hsl(var(--foreground))" }}>{label}</p>
          {payload.map((p: any) => (
            <p key={p.name} style={{ color: p.color }}>{p.name}: {p.value.toFixed(0)}</p>
          ))}
        </div>
      );
    }
    return null;
  };

  return (
    <Layout title="Analytics" subtitle="Performance tracking across all brands and platforms">
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatCard label="Total Impressions" value={stats.total_impressions.toLocaleString()} sub="Lifetime" icon={<Eye className="w-4 h-4" />} glow="primary" />
        <StatCard label="Total Likes" value={stats.total_likes.toLocaleString()} sub="All platforms" icon={<Heart className="w-4 h-4" />} />
        <StatCard label="Total Shares" value={stats.total_shares.toLocaleString()} sub="All platforms" icon={<Share2 className="w-4 h-4" />} glow="success" />
        <StatCard label="Total Posts" value={performance.reduce((acc, b) => acc + b.post_count, 0).toString()} sub="Across all brands" icon={<TrendingUp className="w-4 h-4" />} />
      </div>

      <div className="grid grid-cols-12 gap-6 mb-6">
        {/* Weekly reach chart */}
        <div className="col-span-12 card-glass rounded-xl p-5">
          <h3 className="text-sm font-semibold mb-4">Estimated Reach Trend</h3>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={weeklyData}>
              <CartesianGrid strokeDasharray="3 3" stroke="hsl(220 18% 14%)" />
              <XAxis dataKey="day" tick={{ fill: "hsl(215 16% 50%)", fontSize: 11, fontFamily: "JetBrains Mono" }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fill: "hsl(215 16% 50%)", fontSize: 11, fontFamily: "JetBrains Mono" }} axisLine={false} tickLine={false} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="impressions" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} name="Impressions" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6">
        {/* Brand table */}
        <div className="card-glass rounded-xl overflow-hidden">
          <div className="px-5 py-4" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
            <h3 className="text-sm font-semibold">Brand Performance</h3>
          </div>
          {loading ? (
            <div className="p-10 text-center text-sm text-muted-foreground font-mono">Loading brand data...</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr style={{ borderBottom: "1px solid hsl(var(--border))" }}>
                    {["Brand", "Posts", "Avg Likes", "Avg Shares", "Score"].map((h) => (
                      <th key={h} className="text-left text-xs font-mono uppercase tracking-wide px-5 py-3" style={{ color: "hsl(var(--muted-foreground))" }}>{h}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {performance.map((b, i) => (
                    <tr key={b.brand_id} className="hover:bg-secondary/30 transition-colors" style={{ borderBottom: i < performance.length - 1 ? "1px solid hsl(var(--border))" : undefined }}>
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-6 h-6 rounded flex items-center justify-center text-xs font-bold" style={{ background: "hsl(var(--primary) / 0.15)", color: "hsl(var(--primary))" }}>{b.brand_name[0]}</div>
                          <span className="text-sm">{b.brand_name}</span>
                        </div>
                      </td>
                      <td className="px-5 py-3 text-sm font-mono">{b.post_count}</td>
                      <td className="px-5 py-3 text-sm font-mono">{b.avg_likes.toFixed(1)}</td>
                      <td className="px-5 py-3 text-sm font-mono">{b.avg_shares.toFixed(1)}</td>
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-2">
                          <div className="h-1.5 w-20 rounded-full overflow-hidden" style={{ background: "hsl(var(--secondary))" }}>
                            <div className="h-full rounded-full" style={{ width: `${b.score}%`, background: "hsl(var(--primary))" }} />
                          </div>
                          <span className="text-xs font-mono" style={{ color: "hsl(var(--primary))" }}>{b.score}</span>
                        </div>
                      </td>
                    </tr>
                  ))}
                  {performance.length === 0 && (
                    <tr>
                      <td colSpan={5} className="px-5 py-8 text-center text-sm text-muted-foreground">No brand data available yet. Start a run to see analytics.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
}
