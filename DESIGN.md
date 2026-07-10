# Loft — Design Doc

A personal focus/study workspace: Pomodoro-style timer + Spotify playback control + session logging, installable as a desktop-like PWA.

## Goal

Ship a small, polished, daily-usable tool — and a resume-defensible project built on stacks with real, recent hands-on experience (Go from work, React/OAuth from prior projects). Depth and honesty over feature count.

## v1 Scope (build this, nothing more)

1. **Spotify OAuth login** (Authorization Code flow) + basic playback control — play/pause, show currently-playing track.
2. **Timer** — user-adjustable work/break intervals (not fixed 25/5), auto-pauses Spotify on break, resumes on work.
3. **Per-session task input** — free-text "what am I working on."
4. **Session logging** to Postgres — task name, start time, duration, completed vs. abandoned.
5. **Clean single-page UI** — warm/cozy aesthetic (see below).
6. **PWA support** — installable to dock/desktop, no App Store.

## Explicitly out of scope for v1

Tracked in `docs/v2-ideas.md`, not built until v1 is shipped and used for a few weeks:

- ML/session-completion prediction or any "productivity scoring"
- Focus-pattern analytics/clustering dashboards
- Multi-device sync beyond what Postgres naturally gives us
- Ambient sound layering, non-Spotify audio sources
- Native Mac App Store app

## Aesthetic direction

**Warm & cozy.** Soft neutral tones (cream/warm-gray backgrounds, muted accent — terracotta or amber rather than a cold blue), rounded corners, generous padding. Should feel like a nice desk lamp, not a clinical dashboard. Now-playing album art can inform accent color later, but that's a nice-to-have, not v1-required. Light mode primary; dark mode is a fast-follow if time allows, not a blocker.

## Tech stack

| Layer | Choice | Why |
|---|---|---|
| Frontend | React + Tailwind (Vite) | Known, fast to a working UI |
| Backend | Go + `chi` | Real, recent work experience; static-binary deploy is simple |
| DB access | `pgx` + `sqlc` | Typed, compile-time-checked queries, no ORM magic |
| Migrations | `golang-migrate` | Standard, simple |
| OAuth | `golang.org/x/oauth2` (generic config, Spotify endpoints set manually) | Official package; Spotify's auth/token URLs are static, ~10 lines to wire up |
| Database | PostgreSQL | Matches session-logging needs, prior schema experience |
| Hosting | Vercel (frontend) + Render or Railway (Go backend + Postgres) | Free tiers, minimal setup |

## Data model

Multi-user-ready from the start (so a friend could log in with their own Spotify later), even though v1 usage is just you.

```sql
users
  id                uuid primary key
  spotify_id        text unique not null
  display_name      text
  refresh_token     text not null        -- encrypted at rest
  access_token      text
  access_expires_at timestamptz
  created_at        timestamptz default now()

sessions
  id                uuid primary key
  user_id           uuid references users(id)
  task_name         text
  work_minutes      int not null          -- interval config used for this session
  break_minutes     int not null
  start_time        timestamptz not null
  duration_minutes  int                   -- actual elapsed time
  completed         boolean not null      -- true = ran to completion, false = abandoned early
  created_at        timestamptz default now()
```

## API surface (v1)

- `GET /auth/login` — redirect to Spotify authorize URL
- `GET /auth/callback` — exchange code for tokens, upsert `users` row, set session cookie/JWT
- `POST /auth/refresh` — refresh access token when expired
- `GET /spotify/now-playing` — proxy to Spotify's currently-playing endpoint
- `POST /spotify/play` / `POST /spotify/pause` — playback control proxy
- `POST /sessions` — create a session record (called on session start)
- `PATCH /sessions/:id` — update on completion/abandonment (duration, completed flag)
- `GET /sessions` — list current user's session history

Access/refresh tokens never reach the frontend — the backend holds them and proxies all Spotify calls.

## Project structure

```
Loft/
├── DESIGN.md
├── README.md
├── .gitignore
├── .env.example
│
├── frontend/                 # React + Tailwind (Vite)
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── public/
│   │   ├── manifest.json     # PWA manifest
│   │   └── icons/
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       ├── api/
│       ├── components/       # Timer, TrackDisplay, TaskInput
│       ├── hooks/            # useTimer, useSpotify, useSession
│       └── pages/
│
├── backend/                  # Go + chi
│   ├── go.mod
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/           # env/settings
│   │   ├── db/               # sqlc-generated code + queries/*.sql
│   │   ├── auth/             # Spotify OAuth handlers, token refresh
│   │   ├── spotify/          # Spotify API client (proxy calls)
│   │   └── sessions/         # session CRUD handlers
│   └── migrations/           # golang-migrate .sql files
│
└── docs/
    └── v2-ideas.md
```

## Deployment (once ready, not tonight)

- Frontend → Vercel
- Backend + Postgres → Render or Railway
- Env vars documented in `.env.example`, real secrets never committed

## Tracking

No GitHub Projects board for v1 — this README/DESIGN checklist plus GitHub Issues for actual bugs is enough for a solo project this size.
