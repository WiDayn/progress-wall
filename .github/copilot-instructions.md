# Progress Wall — AI agent instructions (concise)

This file is a short, actionable guide for AI coding agents working on this repo. Focus on the files and patterns below — they are the fastest way to be productive.

Backend (Go, Clean-ish architecture)
- Key dirs: `backend/models/`, `backend/dto/`, `backend/repository/`, `backend/services/`, `backend/router/`.
- Important files: `backend/main.go`, `backend/config/config.go` (uses `godotenv`), `backend/database/database.go` (uses GORM; supports `sqlite` and `mysql`).
- Routing: `backend/router/Router.go` registers routes under `/api/*`. Example: `userGroup.POST("/register", userhandler.Register)` and `UserHandler.Register` binds `dto.RegisterInput` via `c.ShouldBindJSON`.

Frontend (Vue 3 + TypeScript)
- Key dirs: `frontend/src/components/` (`features/`, `ui/`), `frontend/src/stores/` (Pinia), `frontend/src/lib/` (API helpers), `frontend/src/views/`.
- `frontend/src/lib/api.ts` creates an Axios instance with `baseURL` from `VITE_API_BASE_URL || '/api'`, attaches `Authorization: Bearer <token>` using `useUserStore().getToken()` and logs out on 401. Ensure the user store exposes `getToken()` or adapt `api.ts`.

Concrete patterns to follow (do not invent):
- DTO-first for handlers: handlers call `ShouldBindJSON` into `backend/dto/*.go` types. See `backend/router/UserHandler.go`.
- Business logic in services: handlers call `services.*` (example: `userService.Register(username, email, password)`).
- DB access in repository layer; services orchestrate repository calls.
- Config from env via `backend/config/config.go` — default DB type is `sqlite` (file `progress_wall.db`), or `mysql` if configured.

Developer workflows / commands (Windows PowerShell)
```powershell
# backend
Set-Location backend; cp config.env.example config.env; go run main.go

# frontend
Set-Location frontend; pnpm install; pnpm dev

# optional: docker compose
docker-compose up --build
```

Integration notes & gotchas discovered in code
- API base path: frontend expects API under `/api` (proxy or same-origin). Router registers `/api/users` in backend.
- `api.ts` assumes `useUserStore().getToken()` exists — current `frontend/src/stores/user.ts` in tree provides `logout()` but no `getToken()`; update either the store or `api.ts` to avoid runtime errors.
- Database: `backend/database/database.go` uses a singleton pattern (`once.Do`) and will fatal if `GetDB()` called before `InitDB()`.

When adding endpoints
1. Add request/response types in `backend/dto/` (name types clearly, e.g. `RegisterInput`).
2. Add handler in `backend/router/*` and register route in `Router.go` under the appropriate group (e.g., `/api/users`).
3. Implement service method in `backend/services/` and repository logic in `backend/repository/` if DB access is needed.

Files that are the fastest to open for context
- `backend/router/UserHandler.go` (handler patterns and error responses)
- `backend/config/config.go` (env keys and defaults)
- `backend/database/database.go` (DB initialization and supported drivers)
- `frontend/src/lib/api.ts` (axios, auth header, 401 behavior)
- `frontend/src/stores/user.ts` (local user state; may need token methods)

If anything in this file is unclear or you want more examples (sample DTOs, service signatures, or common test patterns), tell me which area to expand and I will update this doc.