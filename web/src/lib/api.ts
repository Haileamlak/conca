const BASE_URL = "/api";

export async function fetchWithAuth(path: string, options: RequestInit = {}) {
    // Ensure we don't try to use "null" or "undefined" as a token string
    const rawToken = localStorage.getItem("conca_token");
    const token = (rawToken && rawToken !== "null" && rawToken !== "undefined") ? rawToken : null;

    const headers = new Headers(options.headers || {});
    if (token) {
        headers.set("Authorization", `Bearer ${token}`);
    }
    if (options.body && !headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
    }

    const response = await fetch(`${BASE_URL}${path}`, {
        ...options,
        headers,
    });

    if (response.status === 401) {
        localStorage.removeItem("conca_token");
        window.location.href = "/auth";
        throw new Error("Unauthorized");
    }

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return response;
}

export const api = {
    get: (path: string) => fetchWithAuth(path, { method: "GET" }),
    post: (path: string, body?: any) => fetchWithAuth(path, {
        method: "POST",
        body: body ? JSON.stringify(body) : undefined
    }),
    put: (path: string, body?: any) => fetchWithAuth(path, {
        method: "PUT",
        body: body ? JSON.stringify(body) : undefined
    }),
    patch: (path: string, body?: any) => fetchWithAuth(path, {
        method: "PATCH",
        body: body ? JSON.stringify(body) : undefined
    }),
    delete: (path: string) => fetchWithAuth(path, { method: "DELETE" }),
};
