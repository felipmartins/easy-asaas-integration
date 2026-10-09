# Architecture Invariants

These are constraints for implementation and AI-assisted changes. An accepted ADR may add detail. A change that would break an invariant must identify the impact and receive maintainer review before implementation.

## Product boundary

- The project integrates with Asaas only.
- Do not add a generic payment-provider abstraction without a real, approved requirement.
- Start with a Go modular monolith; keep API, worker, domain, application, and adapters cohesive and explicit.

## Package and persistence boundaries

- The domain models Asaas-specific use cases and must not depend on PostgreSQL drivers, HTTP frameworks, or frontend types.
- Application packages own narrow persistence interfaces they consume.
- PostgreSQL implements repositories in `internal/persistence/postgres`.
- Do not add a generic `Repository[T]`, repositories for every database table, or SQL/ORM types to domain APIs.
- Asaas API client interfaces and persistence repository interfaces are separate boundaries. `FakeAsaasClient` tests the external API boundary; it is not a repository.
- Persisted domain state and outbox work that must be atomic are committed in one PostgreSQL transaction.

## Payments and event processing

- Payment creation response is not proof of receipt. `PAYMENT_CONFIRMED` and `PAYMENT_RECEIVED` have distinct meanings.
- A payment attempt whose network outcome is unknown must be reconciled before another charge is created.
- Persist accepted Asaas webhook events before responding with HTTP `200`.
- Use the Asaas event ID for webhook idempotency. Database uniqueness constraints are required for cross-process correctness.
- Webhook delivery is at least once and may be out of order. History, retries, and state transitions must be explicit.
- PostgreSQL inbox/outbox records are durable. Go channels are only local wake-up hints.
- Worker concurrency is bounded; shutdown is graceful; cancellation and deadlines use `context.Context`.

## Security

- Never persist raw card numbers or CVV.
- Never expose the Asaas API key or webhook secret in browser code, logs, examples, tests, or committed configuration.
- Use environment variables or a secret abstraction for credentials.
- Tests and fixtures use synthetic data and Sandbox-only credentials.
- Minimize personal-data collection, retention, and logging.

## AI and architecture decisions

- No AI system participates in payment execution.
- A future AI-generated payment configuration requires validation and explicit user confirmation before payment creation.
- AI must not silently alter an invariant or accepted ADR.
- Any change crossing a module boundary must say whether an ADR is required.
