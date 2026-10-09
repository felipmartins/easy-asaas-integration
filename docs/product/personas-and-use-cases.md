# Personas and Primary Use Cases

## Purpose

This document identifies the people the project serves and the primary payment journeys it should support. It complements the project scope in the repository README. It does not define API contracts or claim that a particular Asaas feature is available to every account; those details belong in the Asaas capability matrix.

## Personas

### Integrating Developer

The integrating developer adds Asaas payment collection to an existing application. They may use the Go toolkit directly or integrate through the reference application's API.

**Needs**

- Clear Go APIs and examples for customers, payments, checkout, subscriptions, and webhooks.
- Predictable errors, idempotent operations, and documented payment state changes.
- A way to test behavior without real Asaas credentials or real customer data.
- Payment components that can be adopted and customized independently of the reference dashboard.

**Success looks like**

- The application can request a payment without handling an Asaas API key in the browser.
- Retries do not create unintended duplicate business operations.
- The application can learn about payment outcomes from durable status updates.

### Merchant Operator

The merchant operator is a non-technical user responsible for configuring or monitoring payment collection for a business.

**Needs**

- Understandable payment status and a clear view of what action is needed.
- Payment flows that do not require handling raw card data.
- Clear distinction between a payment request being created and money being received.
- Safe recovery when an integration or external request has an uncertain outcome.

**Success looks like**

- The operator can tell whether a payment is pending, confirmed, or needs attention without interpreting raw webhook payloads.
- The operator is not asked to copy an API key into browser code or handle sensitive card details.

### Payer

The payer is the customer paying an invoice or purchase through an application that uses Easy Asaas Integration. The payer may interact with an Asaas-hosted page or a reusable payment component.

**Needs**

- A clear payment amount, payment instructions, and confirmation of the current payment state.
- Supported payment choices such as Pix, boleto, or a hosted/tokenized card flow, as applicable to the merchant's Asaas configuration.
- A usable experience on mobile and desktop.

**Success looks like**

- The payer can complete the selected flow without the application exposing provider credentials or storing raw card data.
- A browser redirect or client-side message is not presented as final proof that a payment was received.

### Deployment Operator

The deployment operator runs the self-hosted application and is responsible for its configuration, secrets, database, logs, and graceful shutdown.

**Needs**

- Explicit configuration and startup validation.
- Observable API and worker health, with enough context to investigate failures without logging secrets or unnecessary personal data.
- Durable recovery of accepted webhook work after process restarts.
- Documented database backup, migration, and operational procedures as those features are implemented.

**Success looks like**

- A process restart does not discard persisted webhook events or outbox work.
- Secrets are supplied through environment variables or a secret-management integration, not committed configuration files.

## Primary use cases

### UC-1: Create or find a customer

1. An integrating application submits the customer information required for its payment flow.
2. The backend validates the request and calls the Asaas-specific customer operation.
3. The backend returns a local result with the associated Asaas customer identifier, or a contextual error.

**Expected behavior:** retries and duplicate submissions are handled according to a documented idempotency policy. Customer data is minimized in logs and test fixtures.

### UC-2: Request a one-off Pix or boleto payment

1. The application submits the customer, amount, due date, and payment method information required by the selected flow.
2. The backend validates the command and records enough local state to identify the request and its Asaas attempt.
3. The backend returns the created payment details and its current known status.
4. Later Asaas events or an explicit reconciliation operation update the local status.

**Expected behavior:** successful creation is not reported as proof of receipt. Repeating a request with the same idempotency key does not create a second intended payment.

### UC-3: Request a credit card payment

1. The application selects an Asaas-hosted or tokenized card flow.
2. Sensitive card entry and tokenization occur through the supported secure flow.
3. The backend uses the resulting supported reference or token as required by Asaas.
4. Payment status is updated from confirmed Asaas information.

**Expected behavior:** Easy Asaas Integration does not store raw card numbers or CVV, and the Asaas API key never reaches the browser. Exact supported flows and restrictions must be confirmed in the capability matrix.

### UC-4: Request installment or recurring payments

1. The application supplies an installment or subscription configuration.
2. The backend validates the configuration against documented domain rules and Asaas capabilities.
3. The backend records the resulting Asaas identifiers and local lifecycle state.
4. Subsequent payment events update the related installment or subscription view.

**Expected behavior:** installment schedules and recurring subscription lifecycles are modeled explicitly; a successful setup response does not imply that every future charge will succeed.

### UC-5: Create an Asaas Checkout session

1. The application submits a validated checkout configuration.
2. The backend creates the session through the Asaas-specific checkout integration.
3. The payer completes payment using the returned supported checkout flow.
4. The backend updates local state from Asaas events or reconciliation, not from the browser redirect alone.

**Expected behavior:** checkout URLs and session identifiers are handled as sensitive operational data where appropriate, and no Asaas API key is exposed to the payer.

### UC-6: Receive and process an Asaas webhook

1. The webhook endpoint validates the request according to the documented Asaas security mechanism.
2. The endpoint persists the event before acknowledging receipt.
3. A unique Asaas event ID prevents duplicate processing.
4. A worker processes pending events with bounded concurrency and records success or a retryable failure.
5. The application updates payment-related state and writes any follow-up message to the outbox transactionally.

**Expected behavior:** delivery is treated as at least once. A duplicate event has no duplicate business effect, and persisted work remains recoverable after process restarts.

### UC-7: Synchronize payment status

1. The backend receives an Asaas event or runs an explicit reconciliation operation for a payment whose state is uncertain or stale.
2. The backend applies a valid state transition and records the source of the update.
3. Consumers read the current known payment status through the integration API or reference interface.

**Expected behavior:** state transitions are documented, invalid or stale updates do not silently overwrite newer state, and a payment creation response is not treated as settlement.

### UC-8: Reuse a payment UI component

1. A frontend integrator installs or imports a supported React component.
2. The host application supplies only the public configuration and payment-session reference required by the flow.
3. The component presents the selected payment experience and communicates user-visible progress.
4. The backend remains responsible for credentials, payment creation, and authoritative status updates.

**Expected behavior:** components are customizable and do not contain provider secrets or claim final payment success based only on client-side state.

## Cross-cutting requirements

- Validate payment commands before calling Asaas and return explicit, contextual errors.
- Use request idempotency for retried application commands and the Asaas event ID for webhook deduplication.
- Persist accepted webhook events before processing them.
- Keep durable state in PostgreSQL; Go channels are only an in-process wake-up mechanism.
- Bound worker concurrency and support cancellation and graceful shutdown.
- Exclude secrets, raw card data, and unnecessary personal data from logs, examples, and test fixtures.
- Verify exact Asaas fields, limits, authentication, and payment-method eligibility against official documentation before implementing each capability.

## Decisions still open

- Whether the first release supports one Asaas account per installation or multiple merchant accounts in one deployment.
- Whether the reference application is exclusively self-hosted initially or also operated as a hosted service.
- Which operator configuration and reconciliation workflows belong in the first release.

These decisions must be resolved before the configuration and persistence model commits that depend on them. This document records product needs and does not settle those architecture choices.
