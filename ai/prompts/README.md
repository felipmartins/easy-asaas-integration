# Reusable AI Prompts

Prompts in this directory are starting points, not permission to expand scope. Each prompt must require the AI tool to read `AGENTS.md`, the architecture invariants, the domain glossary, the Asaas capability context, and relevant ADRs.

Before using a prompt for implementation, provide a bounded task and identify the files or module boundary involved. The output must explain behavior, failure cases, tests, security effects, and whether an ADR is required. No prompt may authorize Production Asaas access or automatic payment execution.

## Prompts

- [`change-proposal.md`](change-proposal.md) — prepare a scoped change proposal before editing code.
- [`implementation-review.md`](implementation-review.md) — review a completed diff against the project's constraints.
