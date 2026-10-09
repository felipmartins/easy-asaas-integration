# Domain Glossary

Use these terms consistently in code, documentation, API errors, tests, and AI-generated changes.

| Term | Meaning in this project |
|---|---|
| **Customer** | A local customer record linked to the identifier returned by Asaas. Customer records contain personal data and must be minimized and protected. |
| **Payment** | A specific Asaas charge and its latest known provider state. It is not the same as a subscription or a request to create a charge. |
| **Payment command** | The application request to create or change a payment-related operation. It has local validation and idempotency behavior. |
| **PaymentAttempt** | One attempt to submit a payment command to Asaas, including accepted, definitively rejected, or outcome-unknown results. |
| **Payment status** | The latest known Asaas state for a charge, updated from verified webhook events or explicit reconciliation. Keep provider facts distinct from user-interface labels. |
| **Confirmed** | Asaas has confirmed the payment, but its documentation distinguishes this from funds being available. |
| **Received** | Asaas reports that funds are available in the Asaas account. |
| **Overdue** | The payment passed its due date without payment. It may still later be confirmed or received. |
| **Refund** | A refund operation associated with a payment. A payment can have multiple refund entries, including partial refunds; refund completion is tracked separately. |
| **Chargeback** | A card-payment dispute lifecycle that can affect a previously confirmed or received payment. It is not interchangeable with a refund request. |
| **Subscription** | An Asaas schedule that generates separate payment charges. The subscription and each generated payment have distinct lifecycles. |
| **CheckoutSession** | An Asaas-hosted payment session and its returned link/status. Creating a session is not payment confirmation. |
| **WebhookEvent** | An Asaas event envelope with a provider event ID and event type. Delivery is at least once and can be repeated. |
| **Inbox** | Durable PostgreSQL records for accepted external events awaiting or undergoing processing. |
| **Outbox** | Durable PostgreSQL records for application effects that must be delivered after a domain transaction commits. |
| **IdempotencyKey** | A local key used to recognize a repeated application command. It is distinct from the Asaas webhook event ID. |
| **Outcome unknown** | A network result where Asaas may have accepted a request but the application did not receive a definitive response. Reconcile before retrying. |
| **Repository** | A consumer-owned persistence interface for a cohesive domain aggregate or use-case operation, implemented by an adapter. It is not a generic CRUD abstraction. |
| **Asaas client** | The outbound HTTP adapter for Asaas APIs. It is distinct from a repository and has no provider-neutral meaning. |
| **FakeAsaasClient** | A deterministic test double for the outbound Asaas API boundary. It is not a persistence fake. |
