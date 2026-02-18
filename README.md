<p align="center">
  <img src="logo_black.png" alt="Conca Logo" width="200px">
</p>

**Conca** is an autonomous, self-hosted AI content platform designed for modern creators and enterprises. It transforms raw trends into platform-optimized content across multiple brand identities using a sophisticated `Plan -> Generate -> Evaluate -> Publish` loop, all managed through a premium web dashboard.

---

## 🚀 Key Capabilities

- **Autonomous Agent Loop**: Continuous research and generation cycles powered by Gemini LLM.
- **Interactive Mission Control**: A stunning real-time dashboard to monitor agent activity and global performance metrics.
- **Content Calendar & Approval**: Visual weekly planner with "Auto-Plan" functionality to generate a week of content in one click.
- **Brand Identity Wizard**: Manage multiple "AI Personalities" with distinct industries, voices, and target audiences.
- **Analytics Dashboard**: Multi-brand performance tracking with automated scoring based on live engagement (likes, shares, views).
- **Enterprise-Ready Storage**: Flexible data layer supporting **PostgreSQL** or local JSON for self-hosting.
- **Semantic Memory (RAG)**: Learns from past successes to maintain brand consistency and viral potential.

---

## 🛠 High-Level Architecture

Conca is built with a pragmatic, full-stack architecture:

- **`/web`**: Modern React dashboard built with Vite, TypeScript, TailwindCSS, and shadcn/ui.
- **`/api`**: RESTful Go backend using Chi, featuring JWT auth and SQLite-backed job queues.
- **`/agent`**: The "Brain" – core autonomous logic and the planning loop.
- **`/tools`**: Integrations with Gemini, NewsAPI, NewsData.io, Twitter/X, and LinkedIn.
- **`/memory`**: Persistent storage (Postgres/JSON) and Vector Memory (Gemini Embeddings).

---

## ⚡️ Getting Started

### 1. Prerequisites
- **Go 1.21+**
- **Node.js 18+ & npm**
- **Gemini API Key** (Required for LLM & RAG)

### 2. Installation & Setup
```bash
# Clone the repository
git clone https://github.com/Haileamlak/conca.git
cd conca

# Install Backend dependencies
go mod download

# Install & Build Frontend
cd web
npm install
npm run build
cd ..
```

### 3. Configuration
Create a `.env` file in the root:
```bash
GEMINI_API_KEY="your-gemini-key"
JWT_SECRET="your-secure-secret"
# DATABASE_URL="postgres://user:pass@localhost:5432/dbname" # Optional: falls back to local JSON
```

### 4. Launch
```bash
# Start the full platform (Server + Dashboard + Worker)
go run cmd/server/main.go
```
The platform will be available at **`http://localhost:8080`**.

---

## 📡 Operational Endpoints

The system exposes a clean REST API for external integrations:

- `POST /api/auth/login` - Authenticate and receive JWT
- `GET  /api/analytics` - Retrieve global performance snapshots
- `GET  /api/brands` - List all configured agent identities
- `POST /api/brands/{id}/run` - Trigger an immediate autonomous cycle
- `GET  /api/brands/{id}/calendar/scheduled` - Access upcoming content queue

---

## 📜 License
MIT License - see [LICENSE](LICENSE) for details.
