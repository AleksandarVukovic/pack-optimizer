# 2. In-memory pack size store

## Status

Accepted

## Context

Pack sizes must be configurable at runtime via `PUT /packs/sizes` without a restart, and read on every calculation request (`GET`-style lookup inside `CalculateOptimalPacks`). The original requirement only called for runtime-configurable pack sizes, not durable storage — an in-memory store was the simplest thing that satisfied it.

Beyond satisfying that original requirement, `main` is intended to stay the simple, storage-agnostic core of the project. The plan is to branch off variants that swap in real persistence (Postgres, MongoDB) to demonstrate the same service with different tech stacks, so the core needs to stay easy to fork from and free of any one storage tech's dependencies.

## Decision

Store pack sizes in process memory behind the `PackSvc` interface (`internal/pack/types.go`), implemented by `InMemomorySvc` (`internal/pack/memorypack.go`). Reads and writes are guarded by a `sync.Mutex`, sizes are validated (non-empty, all positive) and sorted on write, and `GetSizes` returns a defensive copy so callers can't mutate internal state.

## Consequences

- Zero external dependencies (no database, no schema migrations) and negligible latency for reads/writes — appropriate for a single small config value.
- Pack size changes do not survive a process restart and are not shared across multiple instances of the service; running more than one replica would give each its own independent configuration.
- `PackSvc` is deliberately kept as the only integration point with storage, so `main` stays a clean base to branch from: a `postgres` or `mongo` branch can add a new `PackSvc` implementation without touching `internal/calculator` or `internal/api`.
