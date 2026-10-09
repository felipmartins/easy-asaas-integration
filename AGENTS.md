# Repository Guidance for AI Agents

## Required context

Before changing code or architecture, read:

- `README.md`
- `ai/context/architecture-invariants.md`
- `ai/context/domain-glossary.md`
- `ai/context/asaas-capabilities.md`
- the relevant documents under `docs/architecture/`, `docs/security/`, and `docs/decisions/`

Use `AGENTS.md` as repository-wide guidance. More specific `AGENTS.md` files may add local instructions but must not silently weaken these invariants.

## Language and scope

- Keep source code, identifiers, comments, documentation, examples, API error messages, and generated repository content in English.
- Build an Asaas-specific toolkit and reference application. Do not add Stripe, Mercado Pago, PayPal, another provider, or a generic multi-provider payment abstraction without an explicitly approved project requirement and ADR.
- Keep a modular monolith first. Do not introduce microservices or RabbitMQ for demonstration.

## Payment and security invariants

- Never store raw card numbers or CVV. Prefer Asaas-hosted Checkout; review any tokenized flow and its data path before implementing it.
- Never expose an Asaas API key or webhook authentication token to browser code.
- Use only Sandbox credentials and synthetic data in tests. Never add production webhook payloads or real customer data to the repository.
- Treat payment creation as charge creation, not proof of receipt. Use verified Asaas webhook events or an explicit reconciliation operation for payment status.
- Treat webhooks as at-least-once delivery. Persist before acknowledging, deduplicate by Asaas event ID, and make effects idempotent.
- Do not put AI in payment execution. Any future AI-generated payment configuration must be validated and explicitly confirmed by a user before it can create a payment.

## Architecture and persistence

- Keep domain concepts Asaas-specific and keep persistence, HTTP, and framework details outside the domain package.
- Follow `docs/decisions/0001-repository-pattern.md`: consumer-owned, use-case-focused repository interfaces; PostgreSQL implementations in the persistence adapter; no generic CRUD repository.
- Use explicit dependencies, `context.Context` for I/O cancellation and deadlines, wrapped errors, bounded workers, and graceful shutdown.
- PostgreSQL is the durable inbox/outbox store. A Go channel is only an in-process wake-up hint.
- Any change that crosses a module boundary must state whether an ADR is required. Do not silently change an accepted architectural decision; update or add an ADR and request maintainer review.

## Validation required for generated or modified code

Before presenting an implementation as complete, run the applicable checks and report their results:

- formatting;
- compilation/build;
- unit tests;
- PostgreSQL integration tests when persistence behavior changes;
- Asaas contract tests using Sandbox-safe fixtures or a deterministic test server;
- static analysis;
- security checks, including secret and dependency scanning when configured.

Do not claim a check passed unless it was actually run. Keep tests deterministic and do not contact Production Asaas from automated tests.

## Change and commit discipline

- Keep each change scoped and reviewable. Explain behavior, tests, and ADR impact.
- Do not create, amend, or rewrite Git commits unless a maintainer explicitly asks for that action.
- Do not use destructive Git operations to make a branch or diff appear clean.
