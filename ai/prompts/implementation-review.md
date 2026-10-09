# Prompt: Review an Implementation Diff

Read `AGENTS.md`, `ai/context/architecture-invariants.md`, `ai/context/domain-glossary.md`, `ai/context/asaas-capabilities.md`, and every relevant ADR before reviewing the supplied diff.

Review only the supplied changes. Report findings first, ordered by severity, with file and line references. Check:

- requested behavior and scope;
- package direction and consumer-owned repository interfaces;
- context cancellation, explicit errors, transaction boundaries, and graceful worker lifecycle;
- webhook persistence before acknowledgement, Asaas event-ID idempotency, duplicate and out-of-order behavior;
- distinction between payment creation, `PAYMENT_CONFIRMED`, and `PAYMENT_RECEIVED`;
- secret handling, personal-data minimization, and prohibition on storing card number/CVV;
- deterministic unit, PostgreSQL integration, and Asaas contract tests as applicable;
- formatting, build, static analysis, and security validation evidence;
- whether an ADR is required for any module-boundary change.

Do not recommend another payment provider or a generic multi-provider abstraction. Do not claim checks passed unless their command results are included.
