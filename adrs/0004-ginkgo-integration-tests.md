# 4. Ginkgo/Gomega integration tests against a real binary

## Status

Accepted

## Context

Unit tests (`internal/**/*_test.go`) cover individual packages in isolation with mocked collaborators, enforced by a ≥70% coverage gate in CI. They don't exercise the full HTTP stack — Goa's generated transport, routing, middleware, and process startup/shutdown — so a class of bugs (wiring, serialization, config, graceful shutdown) can only be caught by hitting a running instance of the service.

## Decision

Add a separate integration test suite in `test/integration/` using [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/), run against the actual compiled binary rather than in-process handlers. `scripts/integration-test.sh` manages the full lifecycle — build, start the server, wait for `/health` readiness, run the suite, stop the server — and is invoked both locally (`make integration-test`) and in CI. The suite can also point at an already-running environment, skipping the local build/start step.

## Consequences

- Bugs in the generated transport layer, middleware ordering, and process lifecycle are caught before deployment, which unit tests with mocked layers cannot surface.
- Integration tests are slower and require a buildable binary and a free port, so they're kept as a separate `make` target and CI step rather than folded into `make test`.
- Ginkgo/Gomega is an additional testing framework alongside the standard library `testing` used for unit tests, which contributors need to learn; this was accepted to get BDD-style readable specs and built-in async/readiness assertions for HTTP-level tests.
