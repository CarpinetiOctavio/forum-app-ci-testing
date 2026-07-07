# Screenshots — pending capture

This repo's README references three screenshots that are not yet captured. Drop
the real images at these exact filenames so existing references resolve:

- **`pipeline-passing.png`** — the GitHub Actions run for this repo showing all
  five jobs green (`backend-tests`, `frontend-tests`, `backend-build`,
  `frontend-build`, `summary`).
- **`backend-coverage.png`** — terminal output of the backend coverage command
  (`go test ./tests/services/... -coverprofile=coverage.out
  -coverpkg=./internal/services/...` followed by `go tool cover
  -func=coverage.out`), showing the 54.1% figure referenced in ADR-004 and the
  README.
- **`frontend-tests-passing.png`** — terminal output of
  `npm test -- --watchAll=false`, showing all 34 frontend tests passing.

Once all three files exist here, they resolve automatically wherever the
README references `docs/screenshots/pipeline-passing.png`,
`docs/screenshots/backend-coverage.png`, and
`docs/screenshots/frontend-tests-passing.png`.
