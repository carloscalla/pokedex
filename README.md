# Pokedex

A full-stack Pokedex app with a Go backend and a React/TypeScript frontend, backed
by data synced from [PokeAPI](https://pokeapi.co/).

Instead of proxying PokeAPI's own pagination on every request, the Go backend
**syncs the entire Pokemon catalog into its own Postgres database once**, then serves
all browsing/search/filter/sort from that local store. On top of that it has a
JWT-authenticated "My Team" feature (save up to 6 Pokemon), persisted in the same
database. See `backend/` for the Go service and `frontend/` for the React app.

## Prerequisites

- Go 1.24+
- Node 20+
- A [Supabase](https://supabase.com/) project (free tier is fine) for Postgres

## Setup

1. Create a Supabase project, then grab its Postgres connection string from
   **Project Settings -> Database -> Connection string** (use the "Session" or
   "Direct connection" string, not the "Transaction pooler" one - this is a
   long-running server, not serverless).
2. In `backend/`, copy `.env.example` to `.env` and fill in `DATABASE_URL`,
   `JWT_SECRET`, and `ADMIN_TOKEN`:
   ```
   cd backend
   cp .env.example .env
   ```
3. Run the backend. **The first run performs a full catalog sync from PokeAPI
   (~1300 Pokemon) and will take a few minutes** - this is expected, and only
   happens once (subsequent runs skip straight to serving):
   ```
   set -a; source .env; set +a
   go run ./cmd/server
   ```
4. In another terminal, run the frontend:
   ```
   cd frontend
   npm install
   npm run dev
   ```
5. Open http://localhost:5173.

To force a catalog refresh later without touching user accounts or teams:
```
go run ./cmd/server --sync
```
Or, with the server already running, `POST /api/admin/sync` with header
`X-Admin-Token: <your ADMIN_TOKEN>` (runs in the background, returns `202` immediately).

## Project layout

- `backend/` - Go service (stdlib `net/http`, `pgx`/Postgres, JWT auth). See
  `backend/internal/` for the sync job, catalog/auth/team packages, and HTTP layer.
- `frontend/` - Vite + React + TypeScript. Catalog browsing, search/filter/sort,
  and team management use [TanStack Query](https://tanstack.com/query) for
  server-state caching; `AuthContext` holds the client-side session.

## Running backend tests

```
cd backend
go test ./...
```
(Targeted tests for pagination math, team validation rules, JWT round-tripping,
and the auth middleware - see the plan this was built from for the full list.)
