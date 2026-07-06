# ADR-004: CI pipeline design — trigger, internal coverage artifact, no coverage gate

**Date:** 2026-07-04
**Status:** Accepted

## Context
`.github/workflows/ci.yml` previously uploaded backend and frontend coverage
output to Codecov, an external SaaS reporting service, via
`codecov/codecov-action`. A decision was needed on whether test automation for
TP6 requires an external coverage-reporting integration, and separately,
whether this pipeline should enforce a minimum coverage percentage as a merge
gate.

## Decision
- **Trigger:** unchanged — runs on `push` and `pull_request` to
  `main`/`master`/`develop`.
- **Jobs:** `backend-tests` and `frontend-tests` run in parallel, each
  building and running its own test suite with coverage instrumentation
  (`go test ./tests/services/... -coverprofile=coverage.out
  -coverpkg=./internal/services/...`; `npm test -- --coverage`).
  `backend-build` and `frontend-build` run after their respective test job
  succeeds; `summary` reports final status.
- **Coverage reporting:** each test job uploads its coverage output
  (`backend/coverage.out`; the `frontend/coverage` directory) as an internal
  GitHub Actions artifact via `actions/upload-artifact`, replacing the
  Codecov integration. No external service is used.
- **No coverage gate:** the pipeline does not fail the build if coverage
  drops below any threshold.

## Rationale
- **A coverage gate is an explicit deliverable of `forum-app-qa-pipeline`
  (TP7), not TP6.** TP6's requirement is that tests run automatically on
  every push and that their results (pass/fail, coverage percentage) are
  visible — it does not require enforcing a minimum. Introducing a gate here
  would silently absorb a requirement that belongs to the next repo in the
  series, undermining the deliberate, incremental scope progression the
  series is built to demonstrate.
- **Removing the external reporting service is justified independently of any
  other repo's configuration.** What TP6's CI needs to prove is that the test
  suite runs automatically and that its coverage output exists and is
  retrievable after the run — both are satisfied by an artifact stored by the
  CI provider itself. An external SaaS integration adds an account/token
  dependency and a third-party surface for exactly zero additional
  capability TP6 requires; it would be justified once coverage needs to be
  tracked and trended over time across many runs, which is a `qa-pipeline`-
  scope concern, not a `ci-testing`-scope one.
- **Parallel test jobs with a build dependency is a correctness ordering, not
  an optimization.** `backend-build`/`frontend-build` depending on their
  test job (`needs: backend-tests` / `needs: frontend-tests`) ensures a
  build is never reported successful on top of a failing test suite —
  the CI's core promise for a testing-automation TP.
- **The backend coverage command must scope to the tested package.** The
  original command, `go test ./... -coverprofile=coverage.out`, has no
  `-coverpkg` flag. Since the test files live in `tests/services` (an
  external test package with no non-test statements of its own), Go
  instruments that package by default and reports `coverage: [no
  statements]` — the artifact this job uploaded was never a real number, in
  any run of this pipeline before this fix. Adding
  `-coverpkg=./internal/services/...` tells Go to attribute coverage to the
  actual package under test (`internal/services`, where `AuthService` and
  `PostService` live), matching the command already documented in
  `docs/COMMANDS.md` for local runs. Verified locally: **54.1%** of
  statements in `internal/services`.

## Alternatives considered
- **Keep Codecov:** rejected per the reasoning above — no TP6 requirement it
  satisfies that an internal artifact does not, at the cost of an external
  dependency.
- **Enforce a minimum coverage percentage now:** rejected — out of TP6's
  scope; see ADR-002 for the same reasoning applied to layer selection.
  Enforcing it here would also require deciding a threshold without the
  broader static-analysis context (`qa-pipeline` pairs its coverage gate with
  SonarCloud), which this repo does not have.
- **Single sequential job instead of parallel backend/frontend jobs:**
  rejected — backend and frontend test suites are fully independent (no
  shared state, no build-order dependency between them), so running them in
  parallel shortens pipeline time with no loss of correctness guarantees.

## Additional corrections
- **Go module rename: `tp06-testing` → `forum-app-ci-testing`.** The former
  name was an identifier assigned by the course for the assignment, not a
  project name. The module in `go.mod`, and every internal import path
  across the backend's `.go` files, was renamed to align with this
  repository's actual name in the portfolio.
- **Go version correction in `ci.yml`.** `go.mod` declared `go 1.24.1` as
  the module's minimum version, but `ci.yml` installed Go 1.21 via
  `actions/setup-go` in both the `backend-tests` and `backend-build` jobs.
  Corrected to `1.24` in both jobs so the CI environment matches what the
  module itself declares as its requirement. Go has handled
  `GOTOOLCHAIN=auto` since 1.21, so the mismatch did not actually break the
  pipeline — but it did contradict what `go.mod` itself states as the
  minimum required version.

## Consequences
- Coverage output is visible per-run as a downloadable artifact but is not
  tracked or trended across runs — that capability, if needed, belongs to
  `qa-pipeline`.
- A pull request can be merged even if coverage drops, as long as tests pass
  — this is a deliberate, scoped decision, not an oversight.
- The backend coverage command bug (uploading a meaningless `[no statements]`
  profile) had gone undetected because nothing in this pipeline previously
  asserted on the *content* of the uploaded artifact — only on whether the
  step succeeded, which it did regardless of what the profile contained.
  There is no coverage gate to catch this class of problem in this repo by
  design (see rationale above); catching a silently-broken measurement
  depends on periodically inspecting what the artifact actually contains,
  which is what surfaced this one.
