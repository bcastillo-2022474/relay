# Relay API

Multi-tenant webhook delivery backend (see [`brief.md`](../../brief.md) for the product,
[`docs/transactional-outbox.md`](../../docs/transactional-outbox.md) for the ingest design).

```bash
go run ./cmd/api        # serves :8080, OpenAPI docs at /docs
```

---

## Architecture

The layout follows one dependency rule — **domain packages import nothing but other domain
packages; technology packages import the domain; nothing imports a technology package except
`main`** — and a naming rule: packages are named by what they are, which for adapters means
the technology itself. The import path is the architecture documentation; there are no
`adapters/`, `ports/`, or `usecases/` taxonomy folders.

```
cmd/
├── api/main.go            wiring ONLY: construct adapters, construct commands,
│                          hand everything to the server. No logic.
└── relayer/main.go        (future) the outbox relay loop as its own binary

internal/
├── domain/                DOMAIN packages, one per feature: types, invariants,
│   ├── application/       and the PORTS each consumes (e.g. application/repository.go).
│   ├── endpoint/          Interfaces live with their consumer, never in a central
│   ├── event_type/        ports/ directory. The domain/ grouping exists to make the
│   ├── message/           dependency rule lintable in one line: nothing under domain/
│   ├── organization/      may import anything outside domain/ and shared/.
│   └── user/
├── shared/types/          Cross-feature value types (IDs, Slug, Caller…)
├── shared/apperr/         Domain error kinds; the edge maps them to protocols.
│
├── http/                  DRIVING adapter: handlers, response structs, middleware,
│                          error→status mapping. One file per feature's endpoints.
├── grpc/                  (future) sibling of http/, same shape.
│
├── postgres/              DRIVEN adapter: all repository implementations.
│   └── queries/*.sql      SQL lives HERE, next to its only reader (sqlc).
├── nats/                  DRIVEN adapter: JetStream Broker implementation.
├── fakes/                 DRIVEN adapter: in-memory implementations ("memory" is
│                          the technology). Kept honest by contract tests that run
│                          against both the fake and the real adapter.
│
└── relayer/               Neither driving nor driven: application logic that
                           orchestrates two ports (message.Repository → Broker).
                           Lives in-process today; promotable to cmd/relayer.

migrations/                Consumed by tern against the DB, not by any package —
                           the one piece of SQL that does not live in postgres/.
```

### Placement rules (how to answer "where does X go" without asking)

1. **New domain concept** → its own package under `internal/`, containing its types and
   the interfaces it needs from the world.
2. **New way for the world to call us** (HTTP, gRPC, CLI) → a technology-named sibling of
   `http/`. Driving adapters depend on commands; they are constructed in `main`.
3. **New thing we call** (database, broker, external API) → a technology-named package
   implementing a domain port, plus a fake in `fakes/`, plus one contract test suite run
   against both.
4. **SQL queries** → inside `postgres/`, co-located with the adapter that executes them.
   Changing a query and its adapter is one directory, one diff.
5. **Long-running processes** (relayer, future dispatch workers) → their own package;
   they orchestrate ports and contain no technology-specific code.
6. **`main.go` only introduces things to each other.** If it contains an `if` statement
   that isn't config parsing, something is in the wrong place.

### Growth rules (add structure only when its trigger fires)

| Structure | Earned when |
|---|---|
| Service/command extraction | a second caller appears, or a multi-write transaction |
| A separate domain type vs. raw row | shapes genuinely diverge (invariants, state machine) |
| An interface | a real second implementation exists (incl. a fake worth contract-testing) |
| A response DTO decoupled from the domain type | the API contract must outlive the schema |
| A named error (vs. an `apperr` kind) | some caller branches on it with `errors.Is` |
| A new binary under `cmd/` | independent deploy/scaling need, not tidiness |

Corollaries: if a mapper copies fields 1:1, the boundary it crosses is fiction — delete it.
Adding one column should touch ~2 files; every extra touch is a layer that must justify
itself. Errors carry **kinds** (`apperr`), never HTTP statuses — the domain is also consumed
by workers that make retry decisions, not just by HTTP.

### Testing strategy

- **Logic** (validation, state transitions): pure functions, table-driven tests, no doubles.
- **Glue** (commands, repositories): tests against the real dependency — Postgres via a
  tx-rollback harness once `postgres/` lands.
- **Seams** (NATS, future external services): stateful fake + a contract suite executed
  against both the fake and the real thing, so the fake can never drift.
- Mocks that assert "method A was called, then method B" are banned: they restate the
  implementation and pass while production breaks.
