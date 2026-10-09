# Initial Contribution and Release Strategy

## Contribution workflow

1. Open or identify a focused issue describing the user need, scope, and acceptance conditions.
2. Check the project charter, architecture documents, Asaas capability matrix, security threat model, and relevant ADRs.
3. For a change crossing a module boundary, state whether an ADR is required before implementation.
4. Keep the pull request limited to one cohesive behavior or decision. Include documentation for changed behavior and operational effects.
5. Run the applicable formatting, build, unit, integration, contract, static-analysis, and security checks; report results.
6. Review secret handling, personal data, card-data boundaries, webhook idempotency, and payment-state evidence.
7. Obtain maintainer review before merging. Do not merge a change with an unresolved architectural decision or failed required validation.

## Versioning and releases

- Select and publish the project license before the first public release. Until then, the README must continue to state that the license is pending.
- Begin with pre-1.0 versions while public Go APIs and domain contracts are still evolving.
- Adopt Semantic Versioning once public package compatibility guarantees are documented.
- Tag releases with `vMAJOR.MINOR.PATCH`, publish concise release notes, and identify breaking API or database-migration changes.
- Release only from a reviewed, validated commit. Record the source revision and attach generated artifacts or checksums when binary artifacts are distributed.
- Document supported Go versions, database migration sequence, configuration changes, and rollback or recovery notes in release notes.
- Do not publish Production credentials, customer data, real webhook payloads, or secrets in release artifacts.

## Release validation gates

Before a release, verify:

- formatting, compilation, unit tests, and race checks where applicable;
- PostgreSQL integration tests for persistence and transaction behavior;
- Asaas contract tests using Sandbox-safe inputs or deterministic fixtures;
- static analysis and dependency/security checks;
- migration and backup/restore notes for database changes;
- credential and personal-data review for logs, examples, binaries, and generated files;
- architecture decision review for changes crossing module boundaries;
- documented supported Asaas capabilities and known account-level restrictions.

## Governance decisions still open

- The exact open-source license and copyright notice.
- Maintainer and security-reporting contact details.
- Whether releases include binaries, Go packages only, frontend packages, or a combination.
- The support window and policy for security fixes to prior release lines.

These decisions require maintainer approval before the first tagged public release.
