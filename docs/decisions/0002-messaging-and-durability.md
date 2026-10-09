# ADR 0002: Use PostgreSQL Durability and In-Process Signaling First

- **Status:** Proposed
- **Date:** 2026-10-08

## Context

The initial deployment is a Go modular monolith with an API and worker, at-least-once Asaas webhooks, PostgreSQL, and a need to publish durable follow-up work. Webhook work must survive process restarts. Go channels are local to a process and cannot serve as a durable queue or transport between independently running API and worker processes.

RabbitMQ is available as a possible later transport, but a broker adds deployment, monitoring, retry, and recovery responsibilities. No initial distributed-processing requirement has been established.

## Decision

- Persist incoming webhook events in a PostgreSQL inbox before returning HTTP `200` to Asaas.
- Use a Go channel only as an in-process wake-up hint after an inbox transaction commits. The channel carries no unique source of truth; a periodic or startup database scan recovers pending work.
- Run the API and bounded worker in the same process initially if that keeps deployment and lifecycle simple. A separate worker process may be introduced later; it must discover and claim durable work from PostgreSQL rather than depend on a channel from the API process.
- Persist outbox messages in PostgreSQL in the same transaction as the domain change that requires the outgoing effect.
- Make inbox processing and outbox delivery idempotent. Assume both work items can be attempted more than once.
- Do not add RabbitMQ to the initial release.

## Recovery and concurrency requirements

Workers must use bounded concurrency and durable claim/lease semantics so a crash does not strand work and two workers do not apply the same event concurrently. The implementation should use database transactions and constraints for correctness; the exact PostgreSQL query or claim mechanism is deferred to the worker implementation and its integration tests.

Retries must be bounded and observable. Persist retry state, next-attempt time, and a reviewable parked state for work that cannot proceed automatically. An in-memory channel must never be the only record that work exists.

## Criteria for reconsidering RabbitMQ

Create or update an ADR before adding a broker. Reconsider RabbitMQ only when an observed requirement demonstrates one or more of the following:

- API and worker must run on separate hosts and PostgreSQL polling/claiming is insufficient.
- Workers need independent horizontal scaling or queues need isolation by workload.
- Retry routing, backpressure, or broker-level delivery guarantees are needed beyond the PostgreSQL inbox/outbox design.
- Operational evidence shows database queue contention or recovery limitations that a broker can address.
- A reliability target cannot be met with the current process and PostgreSQL design.

The proposal must state the measured or concrete limitation, the topology, failure modes, operational ownership, and migration/recovery behavior. A broker must not be added only for demonstration.

## Consequences

### Benefits

- Durable webhook acceptance uses the existing database.
- A local channel can provide low-latency wake-up without affecting correctness.
- The initial deployment avoids an additional service and its failure modes.
- The same durable records support restart recovery, idempotency, and operational inspection.

### Costs and risks

- PostgreSQL is responsible for both application state and durable work queues; indexing, retention, and claim contention require attention.
- Separate worker scaling depends on correct database claim semantics.
- Outbox dispatch destinations and delivery policy must be defined when a real consumer is introduced.

## Alternatives considered

### Go channel as the queue

Rejected. Channel contents disappear on process termination and channels do not communicate across processes.

### RabbitMQ from the start

Rejected. No current deployment requirement justifies the additional broker, and PostgreSQL durability is already required for inbox/outbox records.

### PostgreSQL polling only, without a local signal

Valid fallback and recovery mechanism. The initial co-located API/worker may use a channel to reduce pickup latency while retaining database scans for correctness.

## ADR review triggers

Review this ADR before splitting runtime processes, changing inbox/outbox atomicity, adding a broker, or changing the durable retry/claim model.
