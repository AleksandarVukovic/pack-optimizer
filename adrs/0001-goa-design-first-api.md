# 1. Use Goa for a design-first API

## Status

Accepted

## Context

The service exposes an HTTP API (pack size configuration, pack calculation, health check) that needs a documented contract (OpenAPI/Swagger), consistent request validation, and generated transport/marshalling code so the handler layer can focus on business logic rather than boilerplate routing and encoding.

## Decision

Define the API contract in Go using [Goa v3](https://goa.design/) (`design/optimizer.go`) as the single source of truth, and generate the transport layer, request/response types, validation, and OpenAPI 3.0 spec from it via `make generate` (`goa gen`). Generated code lives under `gen/` and is never edited by hand; service logic lives separately under `internal/api/`.

## Consequences

- The API contract, validation, and documentation (Swagger UI, `/openapi.json`) stay in sync automatically — there is no hand-written spec to drift from the code.
- Any contract change requires editing `design/optimizer.go` and re-running `make generate`; editing `gen/` directly is a no-op on the next generation and must be avoided.
- Contributors need to understand Goa's DSL, which is an extra concept compared to a plain `net/http` or `chi` router, but it removes an entire class of manual marshalling/routing bugs.
