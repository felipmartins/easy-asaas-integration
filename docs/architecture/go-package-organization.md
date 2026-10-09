# Go Package Organization

## Purpose

This document defines the initial Go package boundaries and dependency direction for the backend. It complements [ADR 0001: Consumer-Owned Repositories](../decisions/0001-repository-pattern.md). The layout is a starting point; create a package when it has a cohesive responsibility, not just because a directory appears in this diagram.

## Proposed layout

```text
cmd/
  api/
    main.go                 # API process composition and startup
  worker/
    main.go                 # Optional separate worker process composition
  cli/                      # Defer until a concrete operator use case exists

internal/
  application/
    customers/              # Customer use cases and consumed ports
    payments/               # Payment use cases and consumed ports
    checkouts/              # Checkout use cases and consumed ports
    subscriptions/          # Subscription use cases and consumed ports
    webhooks/                # Webhook event processing use cases
  domain/
    customer.go
    payment.go
    payment_attempt.go
    subscription.go
    checkout_session.go
    webhook_event.go
    money.go
  asaas/
    client/                  # Authenticated HTTP transport and API errors
    customers/               # Asaas customer adapter
    payments/                # Asaas payment adapter
    checkouts/               # Asaas Checkout adapter
    subscriptions/           # Asaas subscription adapter
  persistence/
    postgres/                # Repositories, inbox/outbox stores, migrations
  messaging/
    channels/                # In-process wake-up signaling
  jobs/                      # Bounded worker loop and job dispatch
  platform/
    config/                  # Environment parsing and validation
    httpserver/              # HTTP setup and middleware
    logging/                 # Structured logging and redaction

pkg/
  asaas/                     # Public Go API only after its compatibility promise is clear
  webhook/                   # Public helpers only if external consumers need them

web/
  dashboard/
  components/
  sdk/

docs/
  architecture/
  decisions/
  integrations/asaas/
  security/
  operations/
  contributing/

ai/
  prompts/
  context/
  evals/
```

Names and files are provisional. Avoid empty placeholder packages. A cohesive subpackage can be added when a real use case or boundary requires it.

## Dependency direction

```text
cmd -> application -> domain
cmd -> platform/config, HTTP server, PostgreSQL and Asaas adapters (composition only)
PostgreSQL adapter -> application-owned repository interfaces + domain types
Asaas adapter -> application-owned Asaas client interfaces + domain/application types
domain -> Go standard library only, unless a reviewed requirement justifies a dependency
```

- The domain models Asaas-specific concepts and lifecycle rules without importing HTTP, SQL, `pgx`, an ORM, or web framework packages.
- Application packages orchestrate a use case, validate its input, call the domain, and depend on the narrow ports they consume.
- PostgreSQL and Asaas HTTP implementations live in adapter packages and are wired in `cmd`.
- `cmd` packages construct dependencies and own process lifecycle; business rules do not belong in `main.go`.
- Avoid import cycles. A package should not import a concrete adapter to call a use case.

## Repository and client interfaces

Follow [ADR 0001](../decisions/0001-repository-pattern.md): application consumers own narrow persistence interfaces; PostgreSQL implements them. Do not define a generic `Repository[T]` or introduce one interface per database table.

The Asaas HTTP client is a different boundary from a repository. Application use cases may consume narrow Asaas operations; the `internal/asaas` adapters implement them. Do not name that interface `PaymentProvider` or make it multi-provider by implication. `FakeAsaasClient` is a test double for that external API boundary, not for PostgreSQL.

Use `context.Context` on I/O-bound interfaces. Return domain values or stable application errors, not database rows, HTTP response objects, driver errors, or concrete transport types.

## Public API policy

Keep backend implementation under `internal/` while the domain and API evolve. Add a package under `pkg/` only when an external Go consumer has a concrete need and the project is ready to support compatibility and versioning. Do not export internal persistence models as SDK types.

## Testing boundaries

- Unit tests live beside application and domain code and use small in-memory fakes.
- PostgreSQL transaction, uniqueness, and locking behavior requires integration tests against PostgreSQL.
- Asaas request, response, and error behavior requires contract tests using Sandbox-safe data or a deterministic HTTP test server.
- No test may use a real API key, real customer information, or a production webhook payload.

## ADR requirement

This layout documents the proposed dependency direction. An implementation change that crosses package boundaries must state whether it follows this layout or requires a new or updated ADR. Revisit this document if the module structure, persistence ownership, or external API boundary materially changes.
