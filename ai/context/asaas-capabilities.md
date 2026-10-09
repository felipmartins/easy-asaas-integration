# Asaas Capability Context for AI Tools

## Source of truth

The maintained capability matrix is [`docs/integrations/asaas/capability-matrix.md`](../../docs/integrations/asaas/capability-matrix.md). It links official Asaas API and guide references and records the review date. Read it before proposing or implementing a payment-method capability.

This context is a short index, not a substitute for the matrix or the live official API reference. Asaas capabilities, account eligibility, and request contracts can change.

## Verified baseline

- The Asaas API has customer, payment, subscription, checkout, and webhook operations.
- One-off payments document `BOLETO`, `PIX`, `CREDIT_CARD`, and `UNDEFINED` billing types; the payer's choices depend on account configuration.
- Asaas Checkout is hosted. The current Checkout guide documents Pix and credit card and distinguishes one-off, recurring, and installment charge types.
- Subscriptions schedule individual payment charges. A regular subscription with `billingType: PIX` is not the same as Pix Automatic authorization/debit.
- Tokenization has separate account activation requirements for Production and must pass a security review before adoption.
- Webhook events can be repeated; use the Asaas event ID for idempotency and verify the configured `asaas-access-token` header.
- Payment creation or Checkout session creation is not evidence that a payment was received.

## Instructions to AI tools

- Do not invent endpoint fields, supported method combinations, installment limits, status transitions, webhook guarantees, or account permissions.
- Check the current official Asaas documentation for the exact endpoint before changing its client or domain mapping.
- If official sources conflict or leave a behavior unclear, record the uncertainty and request Sandbox or maintainer verification; do not silently choose an interpretation.
- Do not add another payment provider or a generic provider interface.
- Do not add real credentials, customer data, or production webhook payloads to code, tests, prompts, or evaluations.
