import { useState, useEffect } from "react";
import { Layout } from "@/components/Layout";
import { StatusBadge } from "@/components/StatusBadge";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/use-toast";
import { Twitter, Linkedin, Heart, Share2, MessageCircle, Eye, Filter } from "lucide-react";
import { api } from "@/lib/api";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

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

function EditPostModal({ post, onClose, onSave }: { post: Post; onClose: () => void; onSave: (topic: string, content: string) => void }) {
  const [topic, setTopic] = useState(post.topic);
  const [content, setContent] = useState(post.content);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center font-sans" style={{ background: "hsl(var(--background) / 0.8)", backdropFilter: "blur(4px)" }}>
      <div className="w-full max-w-lg rounded-2xl p-6 shadow-2xl animate-fade-in" style={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))" }}>
        <h2 className="text-base font-semibold mb-4 text-foreground">Edit Generated Post</h2>
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label className="text-xs font-mono uppercase text-muted-foreground">Topic</Label>
            <Input value={topic} onChange={(e) => setTopic(e.target.value)} className="bg-secondary border-border h-9 text-sm" />
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs font-mono uppercase text-muted-foreground">Content</Label>
            <Textarea value={content} onChange={(e) => setContent(e.target.value)} rows={6} className="bg-secondary border-border text-sm leading-relaxed" />
          </div>
        </div>
        <div className="flex gap-3 mt-6">
          <Button variant="outline" className="flex-1 border-border" onClick={onClose}>Cancel</Button>
          <Button className="flex-1" onClick={() => onSave(topic, content)}>Save Changes</Button>
        </div>
      </div>
    </div>
  );
}

export default function Posts() {
  const [allPosts, setAllPosts] = useState<Post[]>([]);
  const [brands, setBrands] = useState<string[]>(["All Brands"]);
  const [brandFilter, setBrandFilter] = useState("All Brands");
  const [statusFilter, setStatusFilter] = useState("All Status");
  const [platformFilter, setPlatformFilter] = useState("All Platforms");
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [editingPost, setEditingPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();

  const statuses = ["All Status", "published", "draft", "pending_review", "scheduled", "failed"];
  const platforms = ["All Platforms", "twitter", "linkedin"];

  useEffect(() => {
    fetchPosts();
  }, []);

  const fetchPosts = async () => {
    try {
      const res = await api.get("/posts");
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

  const handleUpdateStatus = async (postID: string, brandID: string, status: string) => {
    try {
      const res = await api.patch(`/brands/${brandID}/calendar/status`, { id: postID, status });
      if (res.ok) {
        toast({ title: "Status Updated", description: `Post is now ${status.replace('_', ' ')}.` });
        fetchPosts();
      }
    } catch (error) {
      toast({ title: "Error", description: "Failed to update status.", variant: "destructive" });
    }
  };

  const handleSaveEdit = async (topic: string, content: string) => {
    if (!editingPost) return;
    try {
      const res = await api.put(`/brands/${editingPost.brand_id}/calendar/post`, {
        id: editingPost.id,
        topic,
        content
      });
      if (res.ok) {
        toast({ title: "Post Saved", description: "Your changes have been recorded." });
        setEditingPost(null);
        fetchPosts();
      }
    } catch (error) {
      toast({ title: "Error", description: "Failed to save changes.", variant: "destructive" });
    }
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
      {/* ... Filters UI ... */}
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
            {f.options.map((o: any) => <option key={o} value={o}>{o.split('_').pop()}</option>)}
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
                    {new Date(post.created_at).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                  </p>
                </div>

                <StatusBadge status={post.status as any} />
              </div>

              {expandedId === post.id && (
                <div className="px-4 pb-4 pt-0 animate-fade-in" style={{ borderTop: "1px solid hsl(var(--border))" }}>
                  <div className="mt-3 p-3 rounded-lg" style={{ background: "hsl(var(--secondary))" }}>
                    <p className="text-sm leading-relaxed whitespace-pre-wrap">{post.content || "(No content generated)"}</p>
                  </div>
                  {post.status === "pending_review" && (
                    <div className="flex gap-2 mt-3" onClick={(e) => e.stopPropagation()}>
                      <Button size="sm" className="gap-2 text-xs h-8" style={{ background: "hsl(var(--success))", color: "hsl(var(--success-foreground))" }} onClick={() => handleUpdateStatus(post.id, post.brand_id, "approved")}>Approve</Button>
                      <Button size="sm" variant="outline" className="gap-2 text-xs h-8 border-border" onClick={() => setEditingPost(post)}>Edit</Button>
                      <Button size="sm" variant="outline" className="gap-2 text-xs h-8 border-destructive/40 text-destructive" onClick={() => handleUpdateStatus(post.id, post.brand_id, "failed")}>Reject</Button>
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

      {editingPost && <EditPostModal post={editingPost} onClose={() => setEditingPost(null)} onSave={handleSaveEdit} />}
    </Layout>
  );
}
