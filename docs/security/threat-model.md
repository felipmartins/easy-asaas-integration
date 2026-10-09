# Security Threat Model

## Purpose and scope

This threat model covers the initial Asaas-specific Go API, worker, PostgreSQL persistence, webhook ingress, and future dashboard or reusable payment components. It identifies assets, trust boundaries, threats, and baseline controls before production code is written. It is not a certification or a substitute for a deployment-specific security review.

## Assets

- Asaas Production and Sandbox API keys.
- Asaas webhook authentication token.
- Customer identifiers, contact details, and payment metadata.
- Payment, subscription, checkout, refund, and chargeback state.
- Webhook inbox data, outbox messages, idempotency keys, and audit records.
- Checkout URLs and payer-facing payment instructions.
- Database credentials, backups, logs, metrics, and deployment configuration.
- Build, release, and AI-generated source artifacts.

Raw card numbers and CVV are prohibited storage assets: the application must never persist them.

## Trust boundaries

1. **Browser and payer device to Go API:** untrusted input; browser state and redirects are not payment authority.
2. **Asaas to webhook ingress:** external network input; authenticate the configured webhook token, validate the envelope, and treat delivery as at least once.
3. **Go services to Asaas API:** outbound authenticated traffic holding high-value credentials.
4. **API/worker to PostgreSQL:** privileged persistence boundary containing personal data and payment state.
5. **Runtime to logs, traces, and support artifacts:** secondary data stores that can leak secrets or personal information.
6. **Contributor/AI tooling to repository and release pipeline:** source changes can alter payment behavior or expose credentials.

## Threats and baseline controls

| Threat | Potential impact | Baseline controls |
|---|---|---|
| Asaas API key leaked in browser, source, logs, or support output | Unauthorized operations against the Asaas account | Load from environment variables or a secret abstraction; server-side only; redact headers; never use real keys in tests; rotate/revoke through operator procedures. |
| Webhook spoofing or token disclosure | Forged payment state changes or queue exhaustion | Validate the dedicated `asaas-access-token` secret using constant-time comparison; use HTTPS; restrict access where operationally practical; do not treat IP allowlisting as a replacement for authentication. |
| Duplicate or replayed webhook | Duplicate fulfillment, notifications, or ledger effects | Persist before acknowledgement; unique constraint on Asaas event ID; idempotent business effects and outbox consumers. |
| Out-of-order webhook delivery | Stale event overwrites a newer payment state | Preserve event history; apply explicit transitions; reconcile conflicts using occasional Asaas status lookup; do not rely only on arrival time. |
| Malformed, oversized, or unknown webhook input | Resource exhaustion, parser failure, or loss of future events | Request-size limit, method/content-type validation, robust JSON decoding, forward-compatible unknown fields, bounded worker pool, and quarantine for unknown event types. |
| Ambiguous timeout during payment creation | A second charge is created after Asaas accepted the first request | Persist local idempotency intent; record `outcome_unknown`; reconcile by external reference or provider lookup before retry. |
| Raw card information reaches logs, database, traces, or error reports | Severe customer data exposure and expanded compliance obligations | Prefer Asaas-hosted Checkout; never persist PAN/CVV; redact request bodies; do not log card data; require a separate security/PCI review before any tokenization path that transmits card data through this service. |
| Personal data copied into fixtures or documentation | Repository history permanently contains sensitive data | Use synthetic fixtures only; never include real customer data or production webhook payloads; review examples and generated artifacts. |
| SQL, transaction, or authorization errors leak implementation or secrets | Information disclosure | Map adapter errors to contextual but safe application errors; do not return driver details or credentials to clients; structured logs with redaction. |
| Worker retry storm or poison event | Database load, delayed payments, webhook queue interruption | Bounded concurrency, retry classification, backoff, parked state, queue-age metrics, alerting, and operator recovery procedures. |
| Database loss or unauthorized access | Loss or disclosure of customer/payment state and accepted work | Least-privilege database roles, TLS in deployed environments, protected backups, migration review, and tested recovery procedures. |
| AI-generated code silently changes payment invariants | Unauthorized or incorrect payment actions | Require architecture context, scoped diffs, mandatory validation gates, and explicit ADR review for boundary changes; no AI in payment execution. |

## Card handling boundary

The initial payer-facing preference is Asaas-hosted Checkout. Tokenized flows must be evaluated separately because tokenization requires card and cardholder information in its request and Production activation is account-reviewed. If a future design sends card data through any Easy Asaas Integration component, review the data path, browser capture, server memory, logs, telemetry, provider requirements, and PCI obligations before implementation. The system must never store raw card numbers or CVV.

## Webhook acceptance boundary

Webhook handling must authenticate the request, persist the event before HTTP `200`, deduplicate by the Asaas event ID, and process asynchronously. If durable persistence is unavailable, do not acknowledge acceptance. Keep webhook secrets separate from API keys. The exact payload retention period and stored fields must be selected according to processing and audit requirements, with access controls and deletion policy.

## Open security decisions

- Whether the first deployment has one Asaas account or multiple merchant accounts.
- Whether the dashboard is self-hosted only or offered as a hosted service.
- Which authentication and authorization protect the project’s own API and dashboard.
- Which webhook payload fields must be retained, for how long, and whether encryption at rest is deployment-provided or application-managed.
- Whether tokenized card flows are included in the first release and what PCI review is required.
- Which operational access, key rotation, backup, and incident-response procedures are supported in the initial release.

Resolve decisions that affect tenant isolation or credential ownership before designing configuration and persistence. Revisit this threat model when a trust boundary, data category, or payment flow changes.
