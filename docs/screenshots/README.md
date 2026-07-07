# Screenshots

This directory holds the evidence screenshots referenced in the
README. Drop the real images at these exact filenames so existing
references resolve:

- **`01-pipeline-passing-main.png`** — GitHub Actions run for the
  PR staging→main showing all five jobs green (backend-tests,
  frontend-tests, backend-build, frontend-build, summary).
  Triggered by pull_request — evidence that the pipeline runs as
  a preventive gate, not post-hoc verification.

- **`02.1-branch-protection-gate-staging.png`** — GitHub PR page
  showing "merging is blocked" with required checks pending on
  staging. Evidence that direct push to staging is blocked and
  the gate is enforced.

- **`02.2-branch-protection-gate-main.png`** — GitHub PR page
  showing "merging is blocked" with required checks pending on
  main. Evidence that the same gate applies before anything
  reaches the stable branch.

- **`03-backend-tests-passing.png`** — terminal output of
  `go test ./tests/services/... -v` showing 23/23 tests passing.

- **`04-backend-coverage.png`** — terminal output of
  `go test ./tests/services/... -cover 
  -coverpkg=./internal/services/...` showing 54.1% coverage of
  internal/services — the declared scope of this repo (ADR-002).

- **`05-frontend-tests-passing.png`** — terminal output of
  `npm test -- --coverage --watchAll=false` showing 36/36 tests
  passing across 5 suites.

- **`06-git-history-before-branching.png`** — git log output
  showing the linear commit history on main before the branching
  model was adopted (all commits on a single line, no branches).

- **`07-git-history-after-branching.png`** — git log output
  showing the full branch tree after adopting
  feature→staging→main, with merge commits from all PRs visible.

Once all files exist here, they resolve automatically wherever
the README references them under docs/screenshots/.