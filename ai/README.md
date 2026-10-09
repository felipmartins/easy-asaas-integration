# AI-Assisted Development

## Purpose

This directory contains the project context, reusable prompts, and review cases used when AI tools assist with development. These materials help reviewers evaluate generated changes; they do not replace engineering review or validation.

## Required reading

Before proposing a change, an AI tool must read:

1. `../AGENTS.md` from the repository root.
2. `context/architecture-invariants.md`.
3. `context/domain-glossary.md`.
4. `context/asaas-capabilities.md`.
5. Relevant architecture documents and ADRs.

## Boundaries

- AI must not be part of the initial payment execution path.
- AI may assist with code, tests, documentation, and review, but it must not silently change an architectural decision.
- Changes crossing a module boundary must state whether an ADR is required.
- A future AI assistant may draft a payment configuration only. The configuration must be validated, shown to the user, and explicitly confirmed before any payment-creating operation can occur.
- Never provide AI tools with Production API keys, webhook secrets, real customer data, or production webhook payloads.

## Required review and validation

For generated or AI-modified code, reviewers must see the change scope, applicable invariants, error and retry behavior, and ADR impact. Before merge, run applicable formatting, compilation, unit, integration, and contract tests, static analysis, and security checks. Report each result accurately; do not infer that a check passed because a tool generated the code.

## Contents

- `context/` — architecture invariants, domain terms, and verified Asaas capability context.
- `prompts/` — reusable prompts with explicit scope and review expectations.
- `evals/` — cases that test whether AI-assisted work respects project invariants.
