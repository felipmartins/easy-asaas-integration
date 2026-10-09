# Asaas Capability Matrix

## Purpose and status

This matrix maps the product goals to capabilities documented by Asaas. It is an architecture baseline, not a promise that every capability is available to every Asaas account or enabled in every environment. Account permissions, payment-method eligibility, and endpoint behavior must be verified in Sandbox before implementation and again before production use.

Documentation reviewed on 2026-10-08. Asaas documentation and account rules can change; re-check the linked official references when implementing a capability.

### Status meanings

- **In scope:** the project intends to support the flow, subject to account configuration and verification.
- **Security gate:** the capability is documented but needs a security review before adoption.
- **Separate capability:** related Asaas product, not included in the initial subscription/payment commitment without a separate decision.
- **Project-owned:** the UI or behavior is built by this project rather than provided by an Asaas API feature.

## Capability matrix

| Product need | Documented Asaas capability | Project status and implementation constraint |
|---|---|---|
| Create customers | `POST /v3/customers` creates a customer and returns an Asaas ID. Asaas permits duplicate customer creation. | **In scope.** Persist the returned Asaas ID and define local idempotency/duplicate prevention. Do not assume Asaas enforces uniqueness by CPF/CNPJ or external reference. |
| One-off boleto | `POST /v3/payments` accepts `billingType: BOLETO`; `UNDEFINED` can also expose boleto when it is enabled for the account. The digitable line is available from `GET /v3/payments/{id}/identificationField`. | **In scope.** Retrieve the current line when needed; Asaas notes that it can change after a charge update. Verify which boleto details the reference UI needs. |
| One-off Pix | `POST /v3/payments` accepts `billingType: PIX`. `GET /v3/payments/{id}/pixQrCode` returns QR image data, a copy-and-paste payload, and expiration information. | **In scope.** Confirm account configuration and QR-code behavior in Sandbox. Do not treat charge creation as payment confirmation. Do not depend on the documented temporary QR-code path that works without a registered Pix key; Asaas says that path will be discontinued. |
| One-off credit card | `POST /v3/payments` accepts `billingType: CREDIT_CARD`. Asaas also offers a hosted Checkout flow with card payment. | **In scope.** Prefer hosted Checkout for the initial payer flow. Any direct card-data flow needs a separate security and PCI review; raw card data and CVV must never be stored. |
| Let the payer choose a one-off method | A payment with `billingType: UNDEFINED` can let the payer choose among methods enabled for the account. | **In scope, conditional.** Validate the configured methods and the payer experience; never assume all methods are enabled. |
| Installment payments | The payment endpoint documents `installmentCount` with one of `installmentValue` or `totalValue` for an installment charge. Checkout documents `chargeTypes: INSTALLMENT` and an installment configuration. | **In scope, conditional.** Do not send installment-only fields for a one-off single charge. Verify method combinations, account limits, rounding, and the applicable installment maximum in Sandbox; the Checkout guide lists its own `maxInstallmentCount` range, which must not be generalized to every endpoint. |
| Recurring subscriptions | `POST /v3/subscriptions` creates a charge schedule with a cycle, first due date, and optional end date or maximum payment count. Generated charges have their own payment lifecycle. The reference lists `BOLETO`, `PIX`, and `CREDIT_CARD` among billing types. | **In scope.** Model the subscription separately from its individual charges. A subscription creation response is not proof of payment. For PIX, distinguish a manually paid recurring Pix charge from automatic debit. |
| Pix Automatic | Asaas documents Pix Automatic as a separate authorization-based capability. It can support automatic recurring debit after payer authorization, unlike a normal subscription configured with `billingType: PIX`. | **Separate capability.** Do not treat ordinary Pix subscriptions as Pix Automatic. Defer implementation until its authorization, lifecycle, and webhook requirements receive a separate product and security review. |
| Asaas Checkout | `POST /v3/checkouts` creates an Asaas-hosted payment page. The current Checkout guide lists `PIX` and `CREDIT_CARD` as its payment methods and `DETACHED`, `RECURRENT`, and `INSTALLMENT` as charge types. It documents checkout expiration from 10 to 1,440 minutes. | **In scope, conditional.** Treat Checkout as hosted. The current documented method list does not include boleto. Validate combinations of payment method and charge type; do not infer every combination is accepted. A callback redirect is not financial confirmation. |
| Credit card tokenization | `POST /v3/creditCard/tokenizeCreditCard` accepts card and cardholder information and returns a token associated with a customer. The feature is available in Sandbox; production activation requires a request and is subject to Asaas review. | **Security gate.** Prefer hosted Checkout. Before using tokenization, review where card data is captured and transmitted, minimize its lifetime, ensure it is never persisted or logged, verify account approval, and document the resulting PCI responsibilities. |
| Webhooks and payment synchronization | Asaas sends event objects with an `id`; delivery can be repeated. The webhook auth token is sent in `asaas-access-token`. Asaas currently requires HTTP `200` for successful delivery. | **In scope.** Persist the event before acknowledging it, deduplicate by event ID, process asynchronously, and return HTTP `200` only after durable acceptance. Do not use the Asaas API key as the webhook token. Plan for retries, queue interruption, and the documented event-retention window. |
| Reusable React payment UI | Asaas provides hosted Checkout URLs; reusable components in this repository would be project-owned. | **Project-owned.** Components may start or display a hosted flow through the backend. They must not receive Asaas credentials or claim authoritative payment success from browser state. |

## Cross-cutting constraints

### Authentication and environments

Asaas API requests use the `access_token` HTTP header rather than `Authorization: Bearer`. Sandbox and Production use different API keys and base URLs. The Go client should always send an identifying `User-Agent`, keep the API key server-side, and reject configurations that pair a key with the wrong environment.

### Customer identity

The Asaas customer ID is the reliable provider identifier for later calls. Because the API allows duplicate customer creation, this project must define its own retry and duplicate-prevention behavior before exposing customer creation as an idempotent operation.

### Payment truth and event delivery

Creating a charge, subscription, or Checkout session starts a flow; it does not prove that money was received. Payment state must be updated from verified Asaas events or an explicit reconciliation query. Event processing must tolerate duplicate delivery and delivery without ordering guarantees.

### Card data

Never store raw card numbers or CVV. Never expose the Asaas API key to the browser. Hosted Checkout is the preferred initial card flow. Tokenized flows remain behind the security gate above because tokenization itself accepts card details in its request.

## Verification gates before implementation

1. Run contract tests against Sandbox for each selected endpoint and supported payment-method combination.
2. Confirm account-level enablement and production eligibility, especially for card tokenization and methods exposed through `UNDEFINED`.
3. Verify the exact webhook event types, authentication configuration, retry behavior, and response requirements used by each flow.
4. Confirm whether a feature belongs in the first release or needs a separate architecture/security decision.
5. Update this matrix when the official API reference or account requirements change.

## Official references

All references below are Asaas documentation. They were accessed on 2026-10-08.

- [Create new customer](https://docs.asaas.com/reference/create-new-customer)
- [Create new payment](https://docs.asaas.com/reference/create-new-payment)
- [Get digitable bill line](https://docs.asaas.com/reference/get-digitable-bill-line)
- [Payments via Pix or dynamic QR Code](https://docs.asaas.com/docs/payments-via-pix-or-dynamic-qr-code)
- [Asaas Checkout guide](https://docs.asaas.com/docs/asaas-checkout)
- [Create new Checkout](https://docs.asaas.com/reference/create-new-checkout)
- [Create new subscription](https://docs.asaas.com/reference/create-new-subscription)
- [Difference between Pix Automatic and subscriptions](https://docs.asaas.com/docs/diferen%C3%A7a-entre-pix-autom%C3%A1tico-e-assinaturas-1)
- [Credit card tokenization](https://docs.asaas.com/reference/credit-card-tokenization)
- [Authentication](https://docs.asaas.com/docs/authentication)
- [Create a webhook through the API](https://docs.asaas.com/docs/create-new-webhook-via-api)
- [Receive Asaas events at your webhook endpoint](https://docs.asaas.com/docs/receive-asaas-events-at-your-webhook-endpoint)
- [Webhooks FAQ](https://docs.asaas.com/docs/webhooks-faq)
