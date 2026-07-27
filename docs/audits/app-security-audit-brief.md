# Brief — Full design/security sweep of the app in forum-app-ci-testing

## Context

While auditing `forum-app-qa-pipeline-legacy` (mirror of TP7), a specific coverage
gap appeared in `DeleteComment` that led to something more fundamental: the
application's authentication model is spoofable by design. It was checked
whether this was specific to that repo or also present in `forum-app-ci-testing`
(already closed, tag `v1.0.0`) — it's in both, with the same code. It's not a
problem introduced by qa-pipeline; it's in the app's source code that both
repos share, and which ci-testing inherited without fixing.

The app itself was not what the course evaluated — any application would have
served as a vehicle to demonstrate testing/CI/CD — which is why it makes sense
that no security criteria were applied to its design. But the goal now is a
portfolio, and the framework is already established: anything that doesn't
meet academic/professional standard gets fixed, whether or not it exceeds the
assignment's nominal scope, when exceeding that scope is a condition for the
guarantees the repo claims to demonstrate to be real.

## What's already confirmed — don't re-derive it, verify it if you want but don't start from zero

Reviewed directly against `backend/` of `forum-app-ci-testing` (not inferred):

1. **No real authentication**: the only middleware is CORS (`router.go`).
   `Login` (`auth_service.go`) doesn't issue a token/session. The 4 mutation
   endpoints (`post_handler.go`: `CreatePost`, `DeletePost`, `CreateComment`,
   `DeleteComment`) trust the `X-User-ID` header with no verification. The
   rule "only the author can delete" (`post.UserID != userID`) is real code
   and runs correctly, but `userID` was never authenticated — anyone can set
   that header by hand.
2. **Plaintext passwords**: stored and compared as-is in `auth_service.go`,
   with comments `// in production: hash with bcrypt` never implemented.
3. **Maximally permissive CORS**: `Access-Control-Allow-Origin: *` +
   `X-User-ID` explicitly allowed in headers — not the root cause, but it adds
   no additional barrier.
4. **Internal errors exposed to the client**: `GetAllPosts` and other paths
   return raw `err.Error()` in the HTTP response instead of logging
   server-side and responding with a generic message.

**Already verified as correct — no need to re-audit, but confirm it if your
sweep touches it in passing:**
- Parameterized SQL throughout the repository layer (`post_repository.go`,
  `user_repository.go`) — no injection risk.
- No `dangerouslySetInnerHTML` on the frontend — no stored XSS via rendering.
- Correct `.gitignore`, no committed secrets.

## What you need to sweep and hasn't been thoroughly reviewed from this chat

1. **`backend/internal/database/database.go`**: schema definition — foreign
   keys, cascades (`ON DELETE CASCADE` or otherwise), constraints, indexes. If
   something is missing there that causes data inconsistency (e.g., orphaned
   comments if a post is deleted without cascading), that's a design finding,
   not a security one, but of the same nature — app code that wasn't audited
   with professional criteria.
2. **`backend/internal/models/`** (`post.go`, `users.go`): check whether
   anything in the structs shouldn't be serialized in JSON responses (e.g.,
   the `Password` field of the `User` model — confirm whether it has a
   `json:"-"` tag or whether the hash/password is unintentionally being
   returned in some response to the client).
3. **Dependencies with known vulnerabilities**: `go.sum` (gorilla/mux, sqlite
   driver, testify) and the frontend's `package.json`/`package-lock.json` —
   run or review the equivalent of `npm audit` / `govulncheck` if the
   toolchain allows it. In `qa-pipeline-legacy`, an `npm audit` with 58
   vulnerabilities in transitive dependencies of `react-scripts` (Create React
   App, deprecated) had already been seen — confirm whether `ci-testing` has
   the same picture or whether something changed when it was rebuilt.
4. **Frontend — components** (`Login.tsx`, `CreatePost.tsx`, etc.): check
   whether anything sensitive is stored in `localStorage`/`sessionStorage`/React
   state for longer than necessary (e.g., password in memory after submit),
   and how the `X-User-ID` is persisted on the client between requests
   (`localStorage`? in-memory state that's lost on refresh?).
5. **No pagination on `GetAllPosts`**: returns everything with no limit. Low
   risk for an academic portfolio, but confirm whether any other endpoint has
   the same pattern with no limits (payloads with no maximum size cap, for
   example).
6. **CSRF**: since there are no session cookies (everything goes through an
   explicit header set by JS, not ambient), the classic CSRF vector doesn't
   apply the same way it would in apps with cookies — confirm that reasoning
   instead of assuming it, and flag if there's any other vector that does
   apply given this particular design.
7. **Anything else your sweep finds** that isn't on this list — this brief is
   a floor, not a ceiling. If something stands out for the same reason as the
   4 confirmed points (code that "works" but wasn't designed with professional
   judgment about what it exposes or allows), document it too even if it
   doesn't fit any of the categories above.

## Expected output format

An inventory, not applied changes. For each new finding: what it is, where it
is (file/line), concrete evidence (not "could be a problem" — show the code),
and severity relative to the 4 already confirmed (is it the same caliber as
authentication, or is it minor like the logging point?). If anything on the
"already verified as correct" list actually has a nuance that wasn't seen,
correct it explicitly — the same way the previous audit corrected the
incorrect generalization of the mocked-authorization pattern (`DeletePost` vs.
`DeleteComment`). That kind of self-correction is worth more than confirming
everything already said.

## Non-negotiable principle

Conceptual auditor, not a decision maker. How each thing gets resolved (bcrypt
vs. another alternative, JWT vs. server-side session, etc.) is decided by
Octavio in the ci-testing chat with the alternatives laid out there.
