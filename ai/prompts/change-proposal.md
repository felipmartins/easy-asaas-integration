# Prompt: Prepare a Scoped Change Proposal

Read `AGENTS.md`, all required files in `ai/context/`, and the relevant architecture and decision documents before analyzing the task.

Given the requested change:

1. Restate the requested behavior and identify the smallest affected modules.
2. Check the Asaas capability matrix and current official documentation for any provider-specific behavior.
3. Identify domain, persistence, webhook, security, retry, idempotency, and observability impacts that apply.
4. State whether the change crosses a module boundary and whether a new or updated ADR is required.
5. Propose files, interfaces, failure behavior, and deterministic validation without introducing unrelated work.
6. Identify assumptions and unresolved decisions. Do not implement while a required architectural decision is unresolved.

Never add another payment provider, generic provider abstraction, raw card storage, Production credentials, or an AI payment execution path. Do not create or rewrite Git commits.
