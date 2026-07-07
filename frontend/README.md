# Frontend

Web client for the ai-rag-platform: login/register, document upload with live status, a chat UI backed by `rag-service`, and a real-time notification bell over WebSocket.

Stack: React + TypeScript + Vite + Tailwind CSS v4 + React Router.

## Local development

```bash
cp .env.example .env   # VITE_API_BASE_URL, defaults to http://localhost:8080
npm install
npm run dev            # http://localhost:5173
```

Requires the backend stack (`docker compose up -d` from the repo root) to be running so `api-gateway` is reachable at `VITE_API_BASE_URL`.

## Build

```bash
npm run build     # type-checks then builds to dist/
npm run preview   # serve the production build locally
```

## Docker

Built and served via `docker-compose.yml` at the repo root (`frontend` service, nginx serving the static build). `VITE_API_BASE_URL` is baked in at build time via a Docker build arg since Vite env vars are compile-time.

## Structure

```
src/
  api/          fetch client + typed calls (auth, documents, rag, notifications)
  auth/         AuthContext (JWT in localStorage) + JWT decode helper
  components/   DocumentsPanel, ChatPanel, NotificationBell, form primitives
  pages/        LoginPage, RegisterPage, DashboardPage
```
