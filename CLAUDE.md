# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Go backend for the Septimana app (companion Android app at https://github.com/BananaNosh/Septimanapp). It serves an OpenAPI-defined REST API (via Echo) backed by SQLite/GORM, exposing "events" (a schedule/horarium of gatherings) and "locations" (places, with multilingual titles/descriptions).

## Commands

Build and run:
```
go build ./...
go run main.go          # starts the REST server, auto-migrates and seeds data/septimana.db on startup
```

Test everything:
```
go test ./...
```

Run a single test:
```
go test ./database -run TestRepository_GetEvent
```

Regenerate OpenAPI-derived code (requires `oapi-codegen` installed at `~/go/bin/oapi-codegen`) after editing `openApi/definition/openapi.yaml` or `openApi/definition/types.yaml`:
```
make typeSpecs   # openApi/definition/types.yaml -> openApi/types/oapiSpecTypes.gen.go
make specs       # openApi/definition/openapi.yaml -> openApi/oapiSpec.gen.go (depends on typeSpecs)
make server      # openApi/definition/openapi.yaml -> openApi/oapiServer.gen.go
```
Always regenerate `typeSpecs` before `specs`, since `specs` imports the generated types.

Reset local dev data (deletes the sqlite DB, which is recreated/re-seeded on next run):
```
make cleanData
```

## Architecture

**Layering**: `main.go` → `database.GetRepository()` (opens/migrates/seeds SQLite) → `Openapi.StartRestApi()` (starts Echo server, registers OpenAPI-validated routes backed by the repository).

**OpenAPI-driven server**: `openApi/definition/openapi.yaml` and `types.yaml` are the source of truth for the API surface. `oapi-codegen` generates:
- `openApi/oapiServer.gen.go` — the `ServerInterface`, param structs, and route registration (`RegisterHandlers`) — **do not edit by hand**, regenerate via `make server`.
- `openApi/oapiSpec.gen.go` / `openApi/types/oapiSpecTypes.gen.go` — the embedded spec and shared types — regenerate via `make specs`/`make typeSpecs`.

`openApi/restapi.go` implements `ServerInterface` on `SeptimanappRestApi` (wraps a `database.Repository`) — this is where actual handler logic lives. `openApi/codegenHelper.go` wires up the API-key auth middleware (`appid` query param) that the generated code doesn't fully provide — POST/PUT/DELETE on events require a valid API key; GET requests are unauthenticated (see `Skipper` in `ProvideAuthMiddlewareModifyEvents`).

**Database layer** (`database/`): `Repository` wraps a `*gorm.DB`. `database.go` holds the CRUD methods used by the REST handlers. `horariumHelper.go` and `locationHelper.go` parse seed data from `data/*.json` (weekly schedule "horaria" per year/language, and locations per "overall location") into GORM models on `InitDatabaseFromPath`. `mockDatabase.go` provides `GetMockRepository`/`InitMockRepository` (in-memory sqlite seeded from `data/testData/...`) used by tests instead of the real DB file.

**Types** (`types/types.go`): core GORM models — `Event`, `Location`, `LocatedString` (a polymorphic multilingual string attached to an Event's `Names`, or a Location's `Titles`/`Descriptions` via GORM's `polymorphic:Parent`), and `ApiKeyInfo` (stores only a hash of issued API keys). `Event` and `Location` have custom `MarshalJSON`/`UnmarshalJSON` that flatten their `LocatedString` slices into `map[language]value` JSON objects, and convert `Event.Start`/`End` to/from unix timestamps — this custom (un)marshaling is central to how the API's JSON shape differs from the GORM struct shape.

**Multilingual data**: locations/events don't store German/Latin text in dedicated columns; instead each language variant is a row in `LocatedString`, joined back in application code (see `GetEvent`/`GetEvents`/`GetLocation`/`GetLocations` in `database/database.go`, which manually join `LocatedString` rows by `parent_id`/`parent_type` rather than relying on GORM preloading).

**Security** (`security/apikey.go`): API keys are random UUIDs; only their SHA-256 hash is persisted (`ApiKeyInfo.ApiKeyHash`). `security.StoreNewApiKey()` generates+prints+stores a new key. Never commit a raw key — `api_key.txt` (local key material) is gitignored.

**Tests**: table-driven, live alongside the code they test (`*_test.go` in `database/`, `openApi/`). `database/testUtil.go` builds shared fixture events; `database/mockDatabase.go` gives each test an isolated in-memory DB seeded from `data/testData/<helper>/`. `util.EqualEvent`/comparison helpers in `util/comparison.go` are used for fixture comparisons instead of `reflect.DeepEqual` (needed because of `gorm.Model` timestamp fields).
