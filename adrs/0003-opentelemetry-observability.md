# 3. OpenTelemetry for observability

## Status

Accepted

## Context

Running as a deployed HTTP service (Docker/Heroku), the API needs request tracing and operational metrics (calculate calls, pack size update successes/failures) without locking observability to one specific vendor's proprietary SDK or backend.

## Decision

Instrument the service with [OpenTelemetry](https://opentelemetry.io/) (`internal/telemetry/`), exporting both traces and metrics over OTLP/HTTP:

- Every request is wrapped in a span via `otelhttp` middleware, with the Goa-generated request ID attached as a span attribute (`internal/telemetry/middleware.go`).
- Service-level operations (calculate, update pack sizes) are traced explicitly through a `Tracer` interface, and counted through a `Meter` interface (`optimizer.calculate.count`, `optimizer.update_pack_sizes.count`, `optimizer.update_pack_sizes.failed.count`).
- The OTLP collector endpoint and metrics export interval are configurable via flags/env vars, and both providers are flushed on graceful shutdown within the same 30s window as the HTTP server.

## Consequences

- Traces and metrics can be sent to any OTLP-compatible backend (Jaeger, Grafana Tempo/Mimir, vendor SaaS, etc.) without changing application code — only the endpoint configuration.
- `Tracer` and `Meter` are defined as interfaces (`internal/telemetry/types.go`), so call sites and tests don't depend on the concrete OTel SDK, and the collector can be swapped or mocked without touching `internal/api` or `internal/calculator`.
- Adds a startup/shutdown dependency (the OTLP exporter) and a runtime dependency (a reachable collector) to the deployment; if no collector is configured, exports are effectively a no-op sink rather than a hard failure.
