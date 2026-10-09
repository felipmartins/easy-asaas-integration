# Asaas Webhook Processing Strategy

## Purpose

This document defines the initial durable, at-least-once processing strategy for incoming Asaas webhook events. Asaas event delivery is external input and may be repeated, delayed, or delivered out of order depending on the configured delivery type.

## Intake and acknowledgement

1. Accept only the configured HTTP method and JSON content type; enforce a request-size limit.
2. Validate the configured webhook authentication token from the `asaas-access-token` header using constant-time comparison. Keep this token separate from the Asaas API key.
3. Decode the event envelope without rejecting unknown fields. Require the Asaas event ID and enough event metadata to persist and route it.
4. Insert the event into the PostgreSQL inbox with a uniqueness constraint on the Asaas event ID (scoped to the configured Asaas account if multi-account support is later approved).
5. Commit the inbox row before acknowledging it to Asaas.
6. Return HTTP `200` after durable acceptance. Asaas currently documents `200` as the successful response for webhook delivery.
7. Signal a local worker through a Go channel as a latency optimization. If the channel is full or the process restarts, the database row remains discoverable by the worker.

An already-persisted event ID is a duplicate: do not repeat its business effects, and return HTTP `200` after confirming the existing durable record. If authentication fails, do not persist or process the event. If the database is unavailable, do not acknowledge durable acceptance; return an error so Asaas can retry.

An authenticated event with a valid envelope but an unknown event type should be persisted for inspection and handled without crashing the worker. A future event type must not be silently treated as a known payment transition.

## Inbox state and worker behavior

The inbox should distinguish at least:

- accepted and pending;
- claimed or processing, with a recovery mechanism for an abandoned claim;
- processed;
- retryable failure with a next-attempt time;
- parked for operator review after a permanent or exhausted failure policy.

The exact retry limits and lease duration are operational configuration to define before implementation. Processing must use a bounded worker pool, respect context cancellation, and release or expire claims safely during shutdown or crashes. Multiple workers must claim rows without concurrently applying the same event effects.

For every event, persist the processing result and any payment-state changes transactionally where they share the PostgreSQL database. If processing causes a follow-up side effect, insert an outbox message in the same transaction. Outbox delivery is at least once and its consumer must also be idempotent.

## Idempotency and ordering

- Use the Asaas event `id` as the inbox deduplication key.
- Use database uniqueness constraints for effects that must not occur twice; an in-memory check is insufficient across restarts or replicas.
- Do not assume the arrival order is the event order. Asaas offers sequential and non-sequential delivery modes; the latter does not guarantee ordering.
- Keep the event record and observed payment status so handlers can detect stale or conflicting information.
- When an event cannot safely update the local projection, reconcile that payment with an occasional Asaas status request and apply a documented rule. Do not implement continuous polling as a replacement for webhooks.
- Be forward-compatible with additional event fields. Keep unknown event types visible for review rather than failing the entire worker loop.

## Retry and queue health

Classify processing errors as retryable or permanent. Retry transient database or network failures with bounded backoff and observable attempt counts. Validation failures, unsupported transitions, or missing local references should be parked for review rather than retried forever without a changed condition.

Expose metrics and logs for accepted events, duplicate deliveries, queue age, processing duration, retry count, parked events, and outbox backlog. Redact credentials and minimize personal data in payload logs. Alert on growing backlog and events approaching the retention window documented by Asaas.

## Security requirements

- Store the webhook authentication token as a secret; never reuse or expose the Asaas API key.
- Accept events only over HTTPS in deployed environments.
- Apply request-size limits, JSON validation, rate limiting, and constant-time token comparison.
- Persist only the payload data required to process and audit the event; define retention and access controls in the security and operations documents.
- Do not include real customer data or production payloads in examples, tests, logs, or this repository.
- Do not treat a callback URL or client-side redirect as authoritative payment evidence.

## Asaas delivery facts to account for

The official documentation says Asaas retries failed deliveries, may interrupt a queue after consecutive failures, retains events for a limited period, and recommends deduplicating with the event ID. The current Webhooks FAQ says HTTP `200` is the successful response and that only `200` is accepted even though other `2xx` responses exist. Recheck this contract during implementation.

## Official references

All references below are Asaas documentation, accessed on 2026-10-08.

- [Create a webhook through the API](https://docs.asaas.com/docs/create-new-webhook-via-api)
- [Receive Asaas events at your webhook endpoint](https://docs.asaas.com/docs/receive-asaas-events-at-your-webhook-endpoint)
- [Webhook events](https://docs.asaas.com/docs/webhooks-events)
- [Webhooks FAQ](https://docs.asaas.com/docs/webhooks-faq)
- [Events for Payments](https://docs.asaas.com/docs/payment-events)
