# Decisions

Each Architecture Decision Record (ADR) states one decision: its context, the alternatives considered and the reasons for the choice. This index only says what each record decides; the reasoning lives in the record itself. Later corrections are added as dated amendments at the end of the record; the original text is not rewritten.

| ADR | Area | What it decides |
|---|---|---|
| [000](ADR-000-resolving-forward-not-mirroring-backward.md) | Method | Each decision rests on this repository's own scope and propagates forward through the series, never copied back from a later repository. Amended 2026-11-12. |
| [001](ADR-001-stack-choice.md) | Stack | Go + React, with testify and Jest, instead of the course's .NET/Angular example. |
| [002](ADR-002-testing-scope-services-layer.md) | Testing | Unit tests cover only the units that branch on a business rule: the backend Services layer, the frontend HTTP services and three components. |
| [003](ADR-003-mocking-strategy.md) | Testing | Mocks replace exactly the I/O boundary (repositories, `axios`); the logic under test runs for real. |
| [004](ADR-004-ci-pipeline-design.md) | Pipeline | Job layout, coverage as an internal artifact without a gate, the `feature/** → staging → main` trigger model, and a `Test Summary` job that fails when any job fails. Amended 2026-11-12 (Node.js 24, required check) and 2026-11-17. |
| [005](ADR-005-package-lock-incident.md) | Pipeline | Retroactive record of the October 2025 lockfile desync, and why CI installs with `npm ci`. |
| [006](ADR-006-test-name-translation.md) | Language | Test names are translated to English. Amended 2026-11-12 and 2026-11-17. |
| [007](ADR-007-application-language.md) | Language | The application's UI text and API error messages are in English; user-authored fixture data stays in Spanish. Amended 2026-11-17. |
| [008](ADR-008-app-security-scope.md) | Security | Hardening applied after `v1.0.0` (bcrypt, 1 MB body limit, no raw internal errors) and the limitations accepted as out of scope. Amended 2026-11-17. |
| [009](ADR-009-comment-language-and-aaa-convention.md) | Language | Residual Spanish comments in tests are translated, and Arrange-Act-Assert markers use a single form. |
| [010](ADR-010-branch-protection-mechanism-migration-classic-to-rulesets.md) | Pipeline | Branch protection on `staging` and `main` moves from classic rules to a ruleset. |
