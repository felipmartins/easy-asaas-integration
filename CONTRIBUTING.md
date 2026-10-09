# Contributing

Thank you for considering a contribution to Easy Asaas Integration. This is an unofficial, open-source community project and is not affiliated with Asaas.

## Before proposing a change

- Read the project [README](README.md), [`AGENTS.md`](AGENTS.md), and relevant architecture, security, and decision documents.
- Keep the proposal focused on one behavior or architectural concern.
- Use English for repository content, identifiers, comments, examples, and API errors.
- Check the Asaas capability matrix and current official documentation for provider-specific behavior.
- State whether the change crosses a module boundary and whether an ADR is required.
- Do not include credentials, raw card data, real customer data, or production webhook payloads.

## Code contributions

For code changes, explain the behavior and failure modes and add deterministic validation appropriate to the change. Before merge, run formatting, compilation, unit tests, PostgreSQL integration tests when persistence changes, Asaas contract tests when provider behavior changes, static analysis, and configured security checks. Report commands and results accurately.

Never test against Production Asaas. Use `FakeAsaasClient`, deterministic HTTP servers, or Sandbox-safe credentials and synthetic data.

## Security and payment changes

Payment, card, credential, webhook, and persistence changes need explicit review against the security threat model and architecture invariants. Never store raw card numbers or CVV. A payment-creation response or browser redirect is not proof that a payment was received.

## Review and commits

Keep pull requests reviewable, describe user-visible behavior and operational impact, and include relevant documentation. Conventional Commit messages are recommended. Do not combine unrelated behavior changes in one pull request or commit.

The project license has not yet been selected; maintainers must choose a license before a public release is made.
