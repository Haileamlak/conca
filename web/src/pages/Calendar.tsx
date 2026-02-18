import { useState, useEffect } from "react";
import { format, startOfWeek, addDays, isSameDay } from "date-fns";
import { Calendar as CalendarIcon, ChevronLeft, ChevronRight, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { api } from "@/lib/api";
import { Layout } from "@/components/Layout";

interface ScheduledPost {
    id: string;
    brand_id: string;
    topic: string;
    content: string;
    platform: string;
    status: "pending_review" | "approved" | "scheduled" | "published" | "failed";
    scheduled_at: string;
}

interface Brand {
    id: string;
    name: string;
}

export default function Calendar() {
    const [currentDate, setCurrentDate] = useState(new Date());
    const [selectedBrand, setSelectedBrand] = useState<string>("");
    const [brands, setBrands] = useState<Brand[]>([]);
    const [posts, setPosts] = useState<ScheduledPost[]>([]);
    const { toast } = useToast();

    useEffect(() => {
        fetchBrands();
    }, []);

    useEffect(() => {
        if (selectedBrand) {
            fetchScheduledPosts();
        }
    }, [selectedBrand]);

    const fetchBrands = async () => {
        try {
            const res = await api.get("/brands");
            const data = await res.json();
            setBrands(data || []);
            if (data?.length > 0) setSelectedBrand(data[0].id);
        } catch (error) {
            console.error("Failed to fetch brands", error);
        }
    };

    const fetchScheduledPosts = async () => {
        try {
            const res = await api.get(`/brands/${selectedBrand}/calendar/scheduled`);
            const data = await res.json();
            setPosts(data || []);
        } catch (error) {
            console.error("Failed to fetch posts", error);
        }
    };

    const handlePlanBatch = async () => {
        try {
            const res = await api.post(`/brands/${selectedBrand}/calendar/plan`);
            if (res.ok) {
                toast({
                    title: "Planning Started",
                    description: "Generating a week of content in the background...",
                });
                setTimeout(fetchScheduledPosts, 5000); // Poll for updates
            }
        } catch (error) {
            toast({
                title: "Error",
                description: "Failed to start planning job.",
                variant: "destructive",
            });
        }
    };

    const renderWeekView = () => {
        const start = startOfWeek(currentDate, { weekStartsOn: 1 });
        const days = Array.from({ length: 7 }).map((_, i) => addDays(start, i));

        return (
            <div className="grid grid-cols-7 gap-4 h-[600px]">
                {days.map((day) => {
                    const dayPosts = posts.filter((p) => isSameDay(new Date(p.scheduled_at), day));
                    return (
                        <div key={day.toString()} className="border rounded-lg p-2 bg-card/50">
                            <div className="text-sm font-medium mb-2 text-muted-foreground">
                                {format(day, "EEE d")}
                            </div>
                            <div className="space-y-2">
                                {dayPosts.map((post) => (
                                    <div
                                        key={post.id}
                                        className={`p-2 rounded text-xs border ${post.status === "published"
                                            ? "bg-green-500/10 border-green-500/20 text-green-500"
                                            : post.status === "approved"
                                                ? "bg-blue-500/10 border-blue-500/20 text-blue-500"
                                                : "bg-yellow-500/10 border-yellow-500/20 text-yellow-500"
                                            }`}
                                    >
                                        <div className="font-semibold truncate">{post.topic}</div>
                                        <div className="text-[10px] opacity-70 capitalize mt-1">
                                            {post.status.replace("_", " ")}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    );
                })}
            </div>
        );
    };

    return (
        <Layout
            title="Content Calendar"
            subtitle="Plan and review your upcoming content strategy."
            actions={
                <div className="flex items-center gap-4">
                    <Select value={selectedBrand} onValueChange={setSelectedBrand}>
                        <SelectTrigger className="w-[200px] h-9 bg-secondary border-border font-mono text-xs">
                            <SelectValue placeholder="Select Brand" />
                        </SelectTrigger>
                        <SelectContent className="bg-popover border-border">
                            {brands.map((b) => (
                                <SelectItem key={b.id} value={b.id} className="text-xs font-mono">
                                    {b.name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                    <Button size="sm" className="gap-2 font-mono text-xs" onClick={handlePlanBatch} style={{ background: "hsl(var(--primary))", color: "hsl(var(--primary-foreground))" }}>
                        <Plus className="w-3.5 h-3.5" />
                        Auto-Plan Week
                    </Button>
                </div>
            }
        >

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-lg font-medium">
                        {format(currentDate, "MMMM yyyy")}
                    </CardTitle>
                    <div className="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="icon"
                            onClick={() => setCurrentDate(addDays(currentDate, -7))}
                        >
                            <ChevronLeft className="h-4 w-4" />
                        </Button>
                        <Button
                            variant="outline"
                            size="icon"
                            onClick={() => setCurrentDate(new Date())}
                        >
                            Today
                        </Button>
                        <Button
                            variant="outline"
                            size="icon"
                            onClick={() => setCurrentDate(addDays(currentDate, 7))}
                        >
                            <ChevronRight className="h-4 w-4" />
                        </Button>
                    </div>
                </CardHeader>
                <CardContent className="pt-6">{renderWeekView()}</CardContent>
            </Card>
        </Layout>
    );
}
