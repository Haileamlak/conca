import { useState, useEffect } from "react";
import { Layout } from "@/components/Layout";
import { StatusBadge } from "@/components/StatusBadge";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/use-toast";
import { Twitter, Linkedin, Heart, Share2, MessageCircle, Eye, Filter } from "lucide-react";

interface Post {
  id: string;
  brand_id: string;
  topic: string;
  platform: string;
  status: string;
  content: string;
  analytics: { views: number; likes: number; shares: number; comments: number };
  created_at: string;
}

export default function Posts() {
  const [allPosts, setAllPosts] = useState<Post[]>([]);
  const [brands, setBrands] = useState<string[]>(["All Brands"]);
  const [brandFilter, setBrandFilter] = useState("All Brands");
  const [statusFilter, setStatusFilter] = useState("All Status");
  const [platformFilter, setPlatformFilter] = useState("All Platforms");
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();

  const statuses = ["All Status", "published", "draft", "pending_review", "scheduled", "failed"];
  const platforms = ["All Platforms", "twitter", "linkedin"];

  useEffect(() => {
    fetchPosts();
  }, []);

  const fetchPosts = async () => {
    try {
      const res = await fetch("/api/posts");
      const data = await res.json();
      setAllPosts(data || []);

      // Extract unique brands
      const uniqueBrands = Array.from(new Set(data.map((p: Post) => p.brand_id)));
      setBrands(["All Brands", ...(uniqueBrands as string[])]);
    } catch (error) {
      console.error("Failed to fetch posts", error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (postID: string) => {
    // In a real app, we'd have a specific endpoint for post approval
    // For now we'll just toast
    toast({ title: "Approved", description: "Post has been moved to scheduled queue." });
  };

  const filtered = allPosts.filter((p) => {
    if (brandFilter !== "All Brands" && p.brand_id !== brandFilter) return false;
    if (statusFilter !== "All Status" && p.status !== statusFilter) return false;
    if (platformFilter !== "All Platforms" && p.platform !== platformFilter) return false;
    return true;
  });

  return (
    <Layout
      title="Post History"
      subtitle={`${filtered.length} posts across all brands`}
      actions={
        <div className="flex items-center gap-2">
          <Filter className="w-4 h-4" style={{ color: "hsl(var(--muted-foreground))" }} />
          <span className="text-xs" style={{ color: "hsl(var(--muted-foreground))" }}>Filters:</span>
        </div>
      }
    >
      {/* Filters */}
      <div className="flex flex-wrap gap-2 mb-6">
        {[
          { label: "Brand", options: brands, value: brandFilter, onChange: setBrandFilter },
          { label: "Status", options: statuses, value: statusFilter, onChange: setStatusFilter },
          { label: "Platform", options: platforms, value: platformFilter, onChange: setPlatformFilter },
        ].map((f) => (
          <select
            key={f.label}
            value={f.value}
            onChange={(e) => f.onChange(e.target.value)}
            className="text-xs font-mono px-3 py-1.5 rounded-lg border outline-none cursor-pointer"
            style={{ background: "hsl(var(--secondary))", borderColor: "hsl(var(--border))", color: "hsl(var(--foreground))" }}
          >
            {f.options.map((o) => <option key={o} value={o}>{o.split('_').pop()}</option>)}
          </select>
        ))}
      </div>

      {loading ? (
        <div className="p-20 text-center text-sm font-mono text-muted-foreground">Loading history...</div>
      ) : (
        <div className="space-y-3">
          {filtered.map((post) => (
            <div
              key={post.id}
              className="card-glass rounded-xl overflow-hidden cursor-pointer hover:border-primary/20 transition-all animate-fade-in"
              style={{ border: "1px solid hsl(var(--border))" }}
              onClick={() => setExpandedId(expandedId === post.id ? null : post.id)}
            >
              <div className="flex items-center gap-4 p-4">
                <div className="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0" style={{ background: post.platform === "twitter" ? "hsl(200 100% 50% / 0.1)" : "hsl(210 80% 55% / 0.1)" }}>
                  {post.platform === "twitter"
                    ? <Twitter className="w-4 h-4" style={{ color: "hsl(200 100% 50%)" }} />
                    : <Linkedin className="w-4 h-4" style={{ color: "hsl(210 80% 55%)" }} />
                  }
                </div>

                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-0.5">
                    <span className="text-xs font-mono px-1.5 py-0.5 rounded" style={{ background: "hsl(var(--primary) / 0.1)", color: "hsl(var(--primary))" }}>{post.brand_id.split('_').pop()}</span>
                  </div>
                  <p className="text-sm font-medium truncate">{post.topic}</p>
                  <p className="text-xs mt-0.5 font-mono" style={{ color: "hsl(var(--muted-foreground))" }}>
                    {new Date(post.created_at).toLocaleDateString("en-US", { month: "short", day: "numeric", hour: "2-digit", minute: "2-digit" })}
                  </p>
                </div>

                {post.status === "published" && (
                  <div className="flex items-center gap-4 text-xs font-mono flex-shrink-0 hidden sm:flex">
                    <span className="flex items-center gap-1" style={{ color: "hsl(var(--muted-foreground))" }}>
                      <Heart className="w-3 h-3" style={{ color: "hsl(var(--destructive))" }} /> {post.analytics.likes}
                    </span>
                    <span className="flex items-center gap-1" style={{ color: "hsl(var(--muted-foreground))" }}>
                      <Share2 className="w-3 h-3" style={{ color: "hsl(var(--primary))" }} /> {post.analytics.shares}
                    </span>
                    <span className="flex items-center gap-1" style={{ color: "hsl(var(--muted-foreground))" }}>
                      <Eye className="w-3 h-3" /> {post.analytics.views.toLocaleString()}
                    </span>
                  </div>
                )}

                <StatusBadge status={post.status as any} />
              </div>

              {expandedId === post.id && (
                <div className="px-4 pb-4 pt-0 animate-fade-in" style={{ borderTop: "1px solid hsl(var(--border))" }}>
                  <div className="mt-3 p-3 rounded-lg" style={{ background: "hsl(var(--secondary))" }}>
                    <p className="text-sm leading-relaxed whitespace-pre-wrap">{post.content || "(No content generated)"}</p>
                  </div>
                  {post.status === "published" && (
                    <div className="grid grid-cols-4 gap-3 mt-3">
                      {[
                        { icon: Heart, label: "Likes", value: post.analytics.likes, color: "hsl(var(--destructive))" },
                        { icon: Share2, label: "Shares", value: post.analytics.shares, color: "hsl(var(--primary))" },
                        { icon: MessageCircle, label: "Comments", value: post.analytics.comments, color: "hsl(var(--success))" },
                        { icon: Eye, label: "Views", value: post.analytics.views.toLocaleString(), color: "hsl(var(--warning))" },
                      ].map((m) => (
                        <div key={m.label} className="rounded-lg p-3 text-center" style={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))" }}>
                          <m.icon className="w-4 h-4 mx-auto mb-1" style={{ color: m.color }} />
                          <p className="text-sm font-mono font-bold">{m.value}</p>
                          <p className="text-xs mt-0.5" style={{ color: "hsl(var(--muted-foreground))" }}>{m.label}</p>
                        </div>
                      ))}
                    </div>
                  )}
                  {post.status === "pending_review" && (
                    <div className="flex gap-2 mt-3">
                      <Button size="sm" className="gap-2 text-xs h-8" style={{ background: "hsl(var(--success))", color: "hsl(var(--success-foreground))" }} onClick={() => handleApprove(post.id)}>Approve & Publish</Button>
                      <Button size="sm" variant="outline" className="gap-2 text-xs h-8 border-border">Edit</Button>
                      <Button size="sm" variant="outline" className="gap-2 text-xs h-8 border-destructive/40 text-destructive">Reject</Button>
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
          {filtered.length === 0 && (
            <div className="p-10 text-center text-sm text-muted-foreground">No posts matching filters.</div>
          )}
        </div>
      )}
    </Layout>
  );
}
