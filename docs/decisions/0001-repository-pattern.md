# ADR 0001: Use Consumer-Owned Repositories for Persistence

- **Status:** Proposed
- **Date:** 2026-10-08

## Context

The backend is a Go modular monolith with domain behavior, application use cases, Asaas HTTP adapters, and PostgreSQL persistence. Payment creation, idempotency, webhook inbox processing, and outbox delivery need deterministic tests without making domain or application code depend on PostgreSQL details.

The project explicitly supports Asaas only. A persistence abstraction must not become a generic provider abstraction or a generic CRUD framework.

## Decision

Use narrow repository interfaces at the application consumer boundary and PostgreSQL implementations behind those interfaces.

- Define an interface in the package that consumes it, normally a cohesive `internal/application/<usecase>` package.
- Implement the interface in `internal/persistence/postgres`.
- Pass domain values through the interface. Do not expose SQL rows, `pgx` types, ORM models, or driver errors to the domain or application layer.
- Add repository methods when an application use case needs them. Do not create a repository for every table or a generic `Repository[T]` CRUD abstraction.
- Model repositories around meaningful domain aggregates or application operations, such as payments, customers, subscriptions, or checkout sessions. Add each only when its use case is implemented.
- Keep webhook inbox and outbox behavior explicit. They may have dedicated storage interfaces because deduplication, claiming work, retry tracking, and dispatch are different from saving a payment aggregate.
- Accept `context.Context` on persistence operations so callers can propagate deadlines and cancellation.
- Map database-specific errors to stable application/domain outcomes at the adapter boundary, with useful error context.

### Transaction boundaries

An application use case owns the decision that a group of operations must be atomic. When a use case must update payment state and insert an outbox message together, the PostgreSQL adapter must commit both in one database transaction. Webhook event processing must also coordinate the inbox processing result with its domain effects so a retry cannot duplicate those effects.

The exact Go transaction API is intentionally deferred until the first use case needs multi-repository atomicity. At that point, define the narrowest application-owned transaction runner or transaction-scoped repository contract that allows the concrete PostgreSQL implementation to provide the required atomicity. Do not leak `pgx.Tx` or `*sql.Tx` into application or domain packages.

Database constraints remain authoritative for cross-process uniqueness, including provider event IDs and local idempotency keys. Application checks may improve error messages but cannot replace those constraints.

## Intended package direction

```text
internal/application/payments
  - payment use case
  - PaymentRepository interface consumed by the use case

internal/domain
  - Asaas-specific payment and lifecycle types
  - no database driver or persistence dependency

internal/persistence/postgres
  - PostgreSQL implementation of application-owned storage interfaces
  - SQL, migrations, transaction handling, and database error mapping
```

This is a dependency direction, not a requirement to create all of these packages immediately. The package containing an interface should be selected by its consumer and the use case it supports.

## Testing implications

- Unit tests for an application use case can provide a small in-memory fake implementing only the interface that use case consumes.
- PostgreSQL adapter behavior, transactions, uniqueness constraints, and locking require integration tests against PostgreSQL.
- `FakeAsaasClient` remains a separate test double for the external Asaas API boundary; it is not a repository.
- Contract tests for the Asaas HTTP client validate provider behavior and do not replace repository integration tests.

## Consequences

### Benefits

- Application and domain behavior can be tested without a database driver.
- PostgreSQL details remain isolated in one adapter package.
- Interfaces describe actual use-case needs and stay smaller than a speculative persistence framework.
- Transaction and idempotency requirements are visible at the application boundary.

### Costs and risks

- A repository interface and a PostgreSQL adapter add code for each persisted behavior.
- Transactional operations across multiple repositories need a deliberate shared transaction boundary.
- Poorly scoped interfaces can become CRUD wrappers or duplicate domain rules.

Mitigate these costs by creating interfaces only when a consumer needs them, keeping them use-case-focused, and testing database invariants against PostgreSQL.

## Alternatives considered

### Generic repository interface

Rejected. Generic CRUD methods hide aggregate invariants, often expose persistence-shaped operations to the application, and make transaction semantics unclear.

### Application services call PostgreSQL directly

Rejected. This couples use cases to a database driver, makes deterministic unit tests harder, and blurs the boundary between application behavior and persistence.

### Repository interfaces declared in the domain package by default

Not selected as a default. Interfaces should live where they are consumed. A domain-owned interface may be appropriate if a domain service itself owns a persistence-dependent port, but it should not be created solely to follow a pattern.

## ADR review triggers

Review or supersede this ADR if a change:

- moves repository interfaces or persistence policy across module boundaries;
- introduces an ORM or a generic repository framework;
- changes how payment state and outbox records achieve atomicity;
- changes inbox claiming, event deduplication, or idempotency enforcement;
- introduces a second persistence technology or distributed transaction mechanism.
