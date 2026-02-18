import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/components/ui/use-toast";
import { Layout } from "@/components/Layout";
import { Brain } from "lucide-react";

export default function Auth() {
    const [isLogin, setIsLogin] = useState(true);
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const { toast } = useToast();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);

        const path = isLogin ? "/api/auth/login" : "/api/auth/register";

        try {
            const res = await fetch(path, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            });

            const data = await res.json();

            if (res.ok) {
                localStorage.setItem("conca_token", data.token);
                toast({ title: isLogin ? "Welcome Back" : "Account Created", description: "Successfully authenticated." });
                navigate("/");
            } else {
                toast({ title: "Auth Failed", description: data.error || "Please check your credentials.", variant: "destructive" });
            }
        } catch (error) {
            toast({ title: "Network Error", description: "Failed to connect to the server.", variant: "destructive" });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center p-4 bg-background font-sans">
            <div className="w-full max-w-md space-y-8 animate-fade-in">
                <div className="text-center">
                    <div className="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto mb-4">
                        <Brain className="w-8 h-8 text-primary" />
                    </div>
                    <h1 className="text-2xl font-bold tracking-tight">Conca</h1>
                    <p className="text-sm text-muted-foreground mt-2">
                        {isLogin ? "Sign in to your autonomous agent dashboard" : "Create your autonomous agent platform"}
                    </p>
                </div>

                <div className="card-glass p-8 rounded-2xl border border-border shadow-xl">
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="email">Email</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="name@example.com"
                                required
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                className="bg-secondary/50 border-border"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                required
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                className="bg-secondary/50 border-border"
                            />
                        </div>
                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? "Authenticating..." : isLogin ? "Sign In" : "Register"}
                        </Button>
                    </form>

                    <div className="mt-6 text-center">
                        <button
                            onClick={() => setIsLogin(!isLogin)}
                            className="text-xs text-muted-foreground hover:text-primary transition-colors"
                        >
                            {isLogin ? "Don't have an account? Register" : "Already have an account? Sign In"}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
