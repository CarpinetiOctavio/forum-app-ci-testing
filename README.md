# Forum App — CI Testing

[![CI/CD Pipeline](https://github.com/CarpinetiOctavio/forum-app-ci-testing/actions/workflows/ci.yml/badge.svg)](https://github.com/CarpinetiOctavio/forum-app-ci-testing/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/github/license/CarpinetiOctavio/forum-app-ci-testing)](LICENSE)

This repository automates the unit tests of a small Go + React forum app with GitHub Actions. Every push to a `feature/**` branch and every pull request into `staging` or `main` runs both suites: 24 backend tests with 55.2% coverage of the Services layer, and 36 frontend tests ([run `43380208077`](https://github.com/CarpinetiOctavio/forum-app-ci-testing/actions/runs/43380208077); [screenshots](docs/screenshots/)). A ruleset blocks any merge into `staging` or `main` unless the `Test Summary` check passes, and that check fails whenever a test or build job fails ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md), [ADR-010](docs/decisions/ADR-010-branch-protection-mechanism-migration-classic-to-rulesets.md)).

> **Series.** This is the first of three repositories: **ci-testing** → [qa-pipeline](https://github.com/CarpinetiOctavio/forum-app-qa-pipeline) → [cloud-deploy](https://github.com/CarpinetiOctavio/forum-app-cloud-deploy). The split follows the course's assignments (TP6, TP7, TP8 of Software Engineering III, UCC); in practice it would be one repository with tags.

Reading paths: in 30 seconds, the opening paragraph; in five minutes, "What this repository shows", "Limitations" and "Errors I found and corrected"; in depth, the decision record each point and each error links. The rest is reference.

## Status

Closed. This repository continues the history of the TP6 submission (October 2025) and was rebuilt as a portfolio piece in 2026. Decisions are recorded in [`docs/decisions/`](docs/decisions/). `CLAUDE.md` and [`docs/rules/`](docs/rules/) are operating instructions for an AI assistant working in this repository, not project documentation.

## What this repository shows

1. **Only the Services layer is unit-tested.** It is the only layer that branches on business rules (validation, "only the author may delete"); testing Repository would mean testing SQL against a database, which is an integration test, not a unit test ([ADR-002](docs/decisions/ADR-002-testing-scope-services-layer.md)).
2. **Mocks stand exactly at the I/O boundary.** Repositories and `axios` are replaced; the service logic runs for real. Mocking the unit itself would make the test tautological; using a real database would make it an integration test ([ADR-003](docs/decisions/ADR-003-mocking-strategy.md)).
3. **The pipeline is a gate before merge, not a check after it.** It runs on pushes to `feature/**` and on pull requests into `staging` and `main`, never on a push to `main`: with a push trigger on `main`, code could reach it without a pull request and be checked only afterwards ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md)). A ruleset enforces it on GitHub's side ([ADR-010](docs/decisions/ADR-010-branch-protection-mechanism-migration-classic-to-rulesets.md)).
4. **Coverage is published, not enforced.** Each run uploads its coverage as an artifact; a minimum threshold is the next repository's scope (TP7), and an external service like Codecov adds an account and a token for nothing this scope needs ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md)).
5. **`npm ci`, not `npm install`.** It fails when `package.json` and `package-lock.json` disagree, instead of silently installing a different tree ([ADR-005](docs/decisions/ADR-005-package-lock-incident.md)).
6. **Go + React instead of the course's .NET/Angular example.** The testing concepts don't depend on the stack, and both sides have a first-class mocking library ([ADR-001](docs/decisions/ADR-001-stack-choice.md)).

## Limitations

- **No real authentication.** The app trusts a client-supplied `X-User-ID` header; CORS allows any origin; lists have no pagination ([ADR-008](docs/decisions/ADR-008-app-security-scope.md), accepted limitations). Out of scope for a unit-testing repository.
- **Unit tests don't verify the real SQL or the real HTTP contract** — only how the services behave given what the mocks return ([ADR-003](docs/decisions/ADR-003-mocking-strategy.md)). Repository, Handlers and three components have no unit tests, by scope ([ADR-002](docs/decisions/ADR-002-testing-scope-services-layer.md)).
- **No coverage gate, static analysis or E2E tests here** — they are [qa-pipeline](https://github.com/CarpinetiOctavio/forum-app-qa-pipeline)'s scope. The actions are pinned by version tag, not by commit SHA, and the workflow declares no `permissions:`; that hardening belongs to the later repositories in the series.
- **Dependency advisories** in Create React App's dev tooling and in the Go standard library are recorded, not fixed ([ADR-008](docs/decisions/ADR-008-app-security-scope.md)).

## Errors I found and corrected

- I first triggered the pipeline on a push to `main`, so code reached `main` before it was validated; I moved to `feature/** → staging → main` through pull requests ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md)).
- My backend coverage command uploaded `[no statements]` in every run until I added `-coverpkg` ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md)).
- `Test Summary` always passed, whatever its dependencies did, until I made it read their results ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md)).
- I bumped the actions' runtime and left `node-version` at 18, out of support since 2025-04-30; the jobs now run on Node.js 24 ([ADR-004](docs/decisions/ADR-004-ci-pipeline-design.md)).
- A lockfile out of sync broke `npm ci` in October 2025 ([ADR-005](docs/decisions/ADR-005-package-lock-incident.md)).

## Decisions

Every decision, with its alternatives and later amendments, is in [`docs/decisions/`](docs/decisions/).

## Evidence

![CI pipeline flow](docs/diagrams/ci-pipeline-flow.svg)

The figures above come from [run `43380208077`](https://github.com/CarpinetiOctavio/forum-app-ci-testing/actions/runs/43380208077). GitHub keeps Actions logs for 90 days; [`docs/screenshots/`](docs/screenshots/) is the record after that.

## Running it

Requires Go 1.24+ and Node.js 24. The test commands CI runs:

```bash
cd backend && go test ./tests/services/... -v -coverprofile=coverage.out -coverpkg=./internal/services/...
cd frontend && npm test -- --coverage --watchAll=false
```

Setup is in [`docs/SETUP.md`](docs/SETUP.md); the remaining commands, including how to run the app, are in [`docs/COMMANDS.md`](docs/COMMANDS.md).

## AI disclosure

Language models (Claude, in chat sessions and through Claude Code) were used to audit the repository, check claims against it, and write text and code. Every design choice was decided by Octavio Carpineti. Only some commits carry a `Co-Authored-By` trailer for that use; the record is irregular.

## License

MIT ([LICENSE](LICENSE)).
