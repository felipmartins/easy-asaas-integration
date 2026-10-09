<<<<<<< HEAD
# easy-asaas-integration
An open-source Go toolkit and reference application that simplifies payment collection with Asaas.
=======
# Easy Asaas Integration

An open-source Go toolkit and reference application for collecting payments through Asaas.

> This is an unofficial, open-source community project and is not affiliated with Asaas.

## Project status

The project is in its architecture-first planning phase. The initial work establishes product scope, domain behavior, security boundaries, and package responsibilities before production code is added.

## Project charter

Easy Asaas Integration aims to make Asaas payment collection easier for developers integrating payments into their applications and for operators configuring payment flows through a reference application.

The project will provide an Asaas-specific Go backend and, as the architecture matures, a reference dashboard, JavaScript SDK, and reusable React payment components. Integrators should be able to adopt the backend toolkit independently of the reference user interface.

## Product goals

- Create and manage Asaas customers.
- Create one-off payments using Pix, boleto, and credit card flows supported by Asaas.
- Prefer Asaas-hosted checkout or tokenized card flows.
- Support installment payments, recurring subscriptions, and Asaas Checkout.
- Receive Asaas webhooks, persist events before processing, and synchronize payment status safely under at-least-once delivery.
- Provide reusable payment UI components that can be customized by other applications.
- Make payment operations understandable through explicit errors, observable processing, and documented behavior.

## Scope boundaries

This project is intentionally specific to Asaas. It will not integrate with other payment providers or introduce a provider-neutral abstraction without a concrete requirement in this project.

The initial system is a modular monolith. Its planned backend consists of a Go API and worker, PostgreSQL for application state and durable inbox/outbox records, and Go channels for in-process notifications. RabbitMQ is an optional future transport and will only be considered when distributed processing, independent worker scaling, retry isolation, or operational reliability establishes a clear need.

TypeScript and React may be used for the dashboard, documentation website, JavaScript SDK, and reusable UI components. The backend remains Go.

## Security commitments

- Never store raw card numbers, CVV, or other sensitive card data.
- Never expose an Asaas API key to a browser.
- Load credentials from environment variables or a secret abstraction.
- Never use real Asaas credentials, customer data, or production webhook payloads in tests or examples.
- Persist webhook events before processing and use the Asaas event ID for idempotency.
- Treat webhook delivery as at least once; duplicate delivery must not duplicate business effects.
- Treat payment creation responses as acknowledgements, not proof that a payment was received.
- Do not put an AI system in the payment execution path. Any future AI-assisted configuration must require explicit user confirmation before a payment can be created.

## Engineering principles

- Keep package boundaries small and cohesive, with explicit dependencies and interfaces at consumer boundaries.
- Use `context.Context` for cancellation and deadlines, wrapped errors with useful context, bounded worker pools, and graceful shutdown.
- Prefer deterministic tests and explicit handling of retries, idempotency, and external API failures.
- Keep repository content, identifiers, comments, documentation, examples, and API errors in English.
- Document architectural decisions and revisit them only through an explicit, reviewed change.

## License

The project license has not yet been selected. A license will be added after the maintainers choose one.
>>>>>>> origin/master
